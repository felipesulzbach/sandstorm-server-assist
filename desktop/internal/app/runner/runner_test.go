package runner

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"sandstorm-server-assist/desktop/internal/adapters/proc"
	"sandstorm-server-assist/desktop/internal/app/logbuf"
	"sandstorm-server-assist/desktop/internal/domain/logparse"
	"sandstorm-server-assist/desktop/internal/domain/problems"
	"sandstorm-server-assist/desktop/internal/domain/server"
)

type fakeOS struct {
	mu      sync.Mutex
	running bool
	cmdline string
	expiry  time.Time
	started time.Time // hora de início informada pelo Windows
	starts  int
}

func setup(t *testing.T) (*Service, *fakeOS, string) {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ServerExe), []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(logPath(dir)), 0o700); err != nil {
		t.Fatal(err)
	}
	f := &fakeOS{}
	p := server.Default()
	p.ServerDir = dir
	p.TravelMap, p.TravelScenario = "Containment", "Scenario_Containment_Checkpoint_Security"
	s := New(Deps{
		Find: func(string) ([]proc.Process, error) {
			f.mu.Lock()
			defer f.mu.Unlock()
			if f.running {
				return []proc.Process{{PID: 4242, Started: f.started}}, nil
			}
			return nil, nil
		},
		Kill: func(uint32) error {
			f.mu.Lock()
			f.running = false
			f.mu.Unlock()
			return nil
		},
		Start: func(_, cmdline, _ string) error {
			f.mu.Lock()
			f.running, f.cmdline = true, cmdline
			f.starts++
			f.mu.Unlock()
			return nil
		},
		Config: func() (server.Profile, server.Secrets, error) {
			return p, server.Secrets{RconPassword: "segredo#1"}, nil
		},
		SessionExpiry: func() time.Time {
			f.mu.Lock()
			defer f.mu.Unlock()
			return f.expiry
		},
		Log: logbuf.New(1000, nil),
	})
	return s, f, dir
}

func appendLog(t *testing.T, dir string, lines ...string) {
	t.Helper()
	fh, err := os.OpenFile(logPath(dir), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()
	for _, l := range lines {
		if _, err := fh.WriteString(l + "\r\n"); err != nil {
			t.Fatal(err)
		}
	}
}

func TestPhasesFromRealLog(t *testing.T) {
	s, f, dir := setup(t)
	appendLog(t, dir, "old log line from previous run")
	if err := s.Start(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(f.cmdline, "-SecurityCode=none") {
		t.Errorf("cmdline: %s", f.cmdline)
	}
	if s.Status().Phase != Starting {
		t.Fatalf("phase %s", s.Status().Phase)
	}
	steps := []struct {
		line string
		want Phase
	}{
		{"[t][  0]LogModioGame: ModSubsystem: mod.io initialization complete with result: ok", Modio},
		{"[t][  0]LogNet: GameNetDriver IpNetDriver_1 IpNetDriver listening on port 27102", Modio},
		{"[t][  0]LogModioGame: ModSubsystem: ModManager enabled", Mods},
		{"[t][  1]LogModioGame: ModSubsystem: OnModActivatedFromCloud: Containment[2123731] bFromGameServer = true", Mods},
		// No mapa inicial a partida passa por LoadingAssets; não é "pronto" porque falta viajar ao Containment.
		{"[t][  1]LogGameState: Match State Changed from EnteringMap to LoadingAssets", Mods},
		{"[t][  1]LogGameMode: ProcessServerTravel: Containment?Scenario=Scenario_Containment_Checkpoint_Security", Traveling},
		{"[t][ 58]LogGameState: Match State Changed from EnteringMap to LoadingAssets", Traveling},
		{"[t][ 70]LogGameState: Match State Changed from LoadingAssets to WaitingToStart", Ready},
	}
	for _, st := range steps {
		appendLog(t, dir, st.line)
		s.Tick()
		if got := s.Status().Phase; got != st.want {
			t.Fatalf("after %q: phase %s, want %s", st.line, got, st.want)
		}
	}
	st := s.Status()
	if st.Port != 27102 || st.Map != "Containment" || len(st.Mods) != 1 || st.Mods[0].State != "active" {
		t.Errorf("%+v", st)
	}
	if strings.Contains(st.CmdLine, "segredo#1") {
		t.Error("segredo na linha mostrada")
	}

	if err := s.Stop(); err != nil {
		t.Fatal(err)
	}
	s.Tick()
	if got := s.Status().Phase; got != Stopped {
		t.Errorf("after stop: %s", got)
	}
}

func TestCrashDetected(t *testing.T) {
	s, f, dir := setup(t)
	if err := s.Start(); err != nil {
		t.Fatal(err)
	}
	time.Sleep(10 * time.Millisecond)
	appendLog(t, dir, "[t][ 21]LogWindows: Error: StaticMesh /Game/X: Serial size mismatch. Got 1, Expected 2",
		"[t][ 21]LogWindows: Error: Fatal error: [File:Unknown] [Line: 3237]")
	s.Tick()
	s.mu.Lock()
	s.st.Since = time.Now().Add(-time.Minute)
	s.mu.Unlock()
	f.mu.Lock()
	f.running = false
	f.mu.Unlock()
	s.Tick()
	st := s.Status()
	if st.Phase != Crashed || st.Problem == nil || st.Problem.Code != "server-crashed" && st.Problem.Code != "mod-conflict" {
		t.Errorf("%+v", st)
	}
}

func TestAuthorizeRestartsWithNone(t *testing.T) {
	s, f, _ := setup(t)
	if err := s.Authorize("12"); err == nil {
		t.Fatal("código inválido aceito")
	}
	if err := s.Authorize("PBRXD"); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(f.cmdline, "-SecurityCode=PBRXD") || !s.Status().Authorizing {
		t.Fatalf("cmdline %s", f.cmdline)
	}
	f.mu.Lock()
	f.expiry = time.Now().Add(90 * 24 * time.Hour) // o servidor gravou a sessão nova
	f.mu.Unlock()
	s.Tick()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		f.mu.Lock()
		cmd := f.cmdline
		f.mu.Unlock()
		if strings.Contains(cmd, "-SecurityCode=none") {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Error("não reiniciou com none")
}

func TestStartRefusesWhenRunningOrMissing(t *testing.T) {
	s, f, dir := setup(t)
	f.running = true
	if err := s.Start(); err == nil {
		t.Error("deveria recusar com servidor rodando")
	}
	f.running = false
	_ = os.Remove(filepath.Join(dir, ServerExe))
	if err := s.Start(); err == nil {
		t.Error("deveria recusar sem executável")
	}
}

// Servidor que já estava rodando antes de o app abrir: "em execução desde" é a hora real de início do processo.
func TestExternalServerKeepsRealStartTime(t *testing.T) {
	s, f, _ := setup(t)
	began := time.Now().Add(-10 * time.Hour).Truncate(time.Second)
	f.mu.Lock()
	f.running, f.started = true, began
	f.mu.Unlock()
	s.Tick()
	st := s.Status()
	if !st.External || !st.Since.Equal(began) {
		t.Errorf("since %v, esperado %v (external=%v)", st.Since, began, st.External)
	}

	// Sem hora do Windows (ou hora no futuro): usa a hora atual.
	if got := startedAt([]proc.Process{{PID: 1}, {PID: 2, Started: time.Now().Add(time.Hour)}}, began); !got.Equal(began) {
		t.Errorf("sem hora: %v", got)
	}
	early := began.Add(-time.Minute)
	if got := startedAt([]proc.Process{{PID: 1, Started: began}, {PID: 2, Started: early}}, time.Now()); !got.Equal(early) {
		t.Errorf("mais antigo: %v", got)
	}
}

// Linhas reais (09/11/2026): com inscrições novas, o servidor monta só os mods que baixou e viaja sem os que já
// estavam instalados. O app reinicia uma vez; na partida seguinte (sem nada novo) ele monta todos.
var newModsRun = []string{
	"[t][  0]LogModioGame: ModSubsystem: mod.io initialization complete with result: ok",
	"[t][  0]LogModioGame: ModSubsystem: Mod added, id: 4422455 to PendingChanges with state Added",
	"[t][  0]LogModioGame: ModSubsystem: Mod added, id: 4774366 to PendingChanges with state Removed",
	"[t][  0]LogModioGame: ModSubsystem: ModManager enabled",
	"[t][823]LogModioGame: ModSubsystem: ModManager Event received. ModId: 4422455, ModEvent: Installed, Message: ok",
	"[t][808]LogModioGame: ModSubsystem: OnModActivatedFromCloud: Border[4422455] bFromGameServer = true",
	"[t][808]LogGameMode: ProcessServerTravel: Containment?Scenario=Scenario_Containment_Checkpoint_Security",
}

func waitStarts(t *testing.T, f *fakeOS, n int) {
	t.Helper()
	for i := 0; i < 100; i++ {
		f.mu.Lock()
		got := f.starts
		f.mu.Unlock()
		if got >= n {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("server started fewer than %d times", n)
}

func TestRestartsOnceWhenInstalledModsWereNotMounted(t *testing.T) {
	s, f, dir := setup(t)
	s.d.InstalledMods = func() []int { return []int{2123731, 98145, 4422455} }
	if err := s.Start(); err != nil {
		t.Fatal(err)
	}
	appendLog(t, dir, newModsRun...)
	s.Tick()
	waitStarts(t, f, 2)
	if st := s.Status(); !st.Remount || st.Phase != Starting {
		t.Fatalf("after restart: %+v", st)
	}

	// Segunda partida: nada novo, monta todos e fica pronta. Mesmo que faltasse algo, não reinicia de novo.
	appendLog(t, dir,
		"[t][  0]LogModioGame: ModSubsystem: ModManager enabled",
		"[t][  1]LogModioGame: ModSubsystem: OnModActivatedFromCloud: Containment[2123731] bFromGameServer = true",
		"[t][  1]LogModioGame: ModSubsystem: OnModActivatedFromCloud: No Restricted Area Mutator [1.15][98145] bFromGameServer = true",
		"[t][  1]LogModioGame: ModSubsystem: OnModActivatedFromCloud: Border[4422455] bFromGameServer = true",
		"[t][  1]LogGameMode: ProcessServerTravel: Containment?Scenario=Scenario_Containment_Checkpoint_Security",
		"[t][ 70]LogGameState: Match State Changed from LoadingAssets to WaitingToStart",
	)
	s.Tick()
	if st := s.Status(); st.Phase != Ready || st.Remount || st.Problem != nil {
		t.Fatalf("second run: %+v", st)
	}
	time.Sleep(100 * time.Millisecond)
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.starts != 2 {
		t.Fatalf("starts = %d", f.starts)
	}
}

func TestNoRestartWhenAllInstalledModsMounted(t *testing.T) {
	s, f, dir := setup(t)
	s.d.InstalledMods = func() []int { return []int{4422455} }
	if err := s.Start(); err != nil {
		t.Fatal(err)
	}
	appendLog(t, dir, newModsRun...)
	s.Tick()
	time.Sleep(100 * time.Millisecond)
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.starts != 1 {
		t.Fatalf("restarted without need: starts = %d", f.starts)
	}
}

func TestExternalServerOnlyWarnsAboutUnmountedMods(t *testing.T) {
	s, f, dir := setup(t)
	s.d.InstalledMods = func() []int { return []int{2123731, 4422455} }
	f.running = true // aberto fora do app
	s.Tick()
	appendLog(t, dir, newModsRun...)
	s.Tick()
	st := s.Status()
	if st.Problem == nil || st.Problem.Code != "mods-not-mounted" || f.starts != 0 {
		t.Fatalf("%+v starts=%d", st, f.starts)
	}
}

func TestEachStartOpensARunAndCleansOnce(t *testing.T) {
	s, _, dir := setup(t)
	cleaned := 0
	s.d.CleanLogs = func(string) { cleaned++ }
	if err := s.Start(); err != nil {
		t.Fatal(err)
	}
	s.Tick()
	if cleaned != 0 {
		t.Fatal("cleaned before the server rotated its log")
	}
	appendLog(t, dir, "[t][  0]LogModioGame: ModSubsystem: mod.io initialization complete with result: ok")
	s.Tick()
	s.Tick()
	if cleaned != 1 {
		t.Fatalf("cleaned %d times", cleaned)
	}
	lines := s.d.Log.Recent(0)
	runs := 0
	for _, l := range lines {
		if l.Kind == "run" {
			runs++
		}
	}
	last := lines[len(lines)-1]
	if runs != 1 || last.Run != 1 {
		t.Fatalf("runs=%d last=%+v", runs, last)
	}
}

// App fechado e aberto de novo com o servidor que ele iniciou (ADR-026): continua como servidor do app, com a linha
// de comando, sem o aviso de "iniciado fora". Quando a partida acaba, o registro é apagado.
func TestReopenedAppResumesItsOwnServer(t *testing.T) {
	s, f, dir := setup(t)
	var saved *Run
	var mu sync.Mutex
	s.d.SaveRun = func(r *Run) { mu.Lock(); saved = r; mu.Unlock() }
	s.d.LoadRun = func() *Run { mu.Lock(); defer mu.Unlock(); return saved }
	if err := s.Start(); err != nil {
		t.Fatal(err)
	}
	mu.Lock()
	run := saved
	mu.Unlock()
	if run == nil || run.Dir != dir || run.CmdLine == "" || strings.Contains(run.CmdLine, "segredo#1") {
		t.Fatalf("registro: %+v", run)
	}

	// O app fecha; o processo começou logo depois do registro. O app reabre.
	f.mu.Lock()
	f.started = run.Started.Add(3 * time.Second)
	f.mu.Unlock()
	again := New(s.d)
	again.Tick()
	st := again.Status()
	if st.External || st.CmdLine != run.CmdLine || !again.ownStart {
		t.Errorf("deveria continuar como servidor do app: %+v", st)
	}

	// O servidor para: o registro some.
	_ = again.Stop()
	again.Tick()
	mu.Lock()
	defer mu.Unlock()
	if saved != nil {
		t.Errorf("registro não apagado: %+v", saved)
	}
}

// Registro antigo, de outra pasta ou processo sem hora do Windows: o servidor é tratado como aberto fora do app.
func TestRunRecordOnlyMatchesItsOwnStart(t *testing.T) {
	at := time.Date(2026, 11, 10, 18, 4, 0, 0, time.UTC)
	run := &Run{Started: at, Dir: `C:\srv`}
	cases := []struct {
		name string
		run  *Run
		dir  string
		p    time.Time
		want bool
	}{
		{"logo depois", run, `C:\srv`, at.Add(2 * time.Second), true},
		{"maiúsculas na pasta", run, `c:\SRV\`, at.Add(time.Second), true},
		{"iniciado bem depois (por fora)", run, `C:\srv`, at.Add(time.Hour), false},
		{"iniciado antes do registro", run, `C:\srv`, at.Add(-time.Minute), false},
		{"outra pasta", run, `C:\outro`, at.Add(time.Second), false},
		{"sem hora do Windows", run, `C:\srv`, time.Time{}, false},
		{"sem registro", nil, `C:\srv`, at, false},
	}
	for _, c := range cases {
		if got := c.run.matches(c.dir, []proc.Process{{PID: 1, Started: c.p}}); got != c.want {
			t.Errorf("%s: %v, esperado %v", c.name, got, c.want)
		}
	}
}

// Linhas reais (10/11/2026): depois de pronto no mapa de mod, troca para Border. O app ficava preso em "a carregar
// o mapa", esperando o mapa da inicialização (Containment), e bloqueava as ações.
func TestMapChangeAfterStartBecomesReady(t *testing.T) {
	s, _, dir := setup(t)
	if err := s.Start(); err != nil {
		t.Fatal(err)
	}
	steps := []struct {
		line string
		want Phase
	}{
		{"[t][  0]LogLoad: LoadMap: /Game/Maps/Farmhouse/Farmhouse?Name=Player?Scenario=Scenario_Farmhouse_Checkpoint_Security?MaxPlayers=10?Lighting=Day", Starting},
		{"[t][  0]LogGameMode: ProcessServerTravel: Containment?Scenario=Scenario_Containment_Checkpoint_Security", Traveling},
		{"[t][ 56]LogGameState: Match State Changed from EnteringMap to LoadingAssets", Traveling},
		{"[t][ 60]LogGameState: Match State Changed from LoadingAssets to WaitingToStart", Ready},
		{"[t][684]LogGameState: Match State Changed from RoundActive to LeavingMap", Ready},
		{"[t][684]LogGameMode: ProcessServerTravel: Border?Scenario=Scenario_Border_Checkpoint_Insurgents?Lighting=Day?MaxPlayers=10", Traveling},
		{"[t][859]LogGameState: Match State Changed from EnteringMap to LoadingAssets", Traveling},
		{"[t][874]LogGameState: Match State Changed from LoadingAssets to WaitingToStart", Ready},
		// Linhas reais (11/11/2026): troca para o mapa de mod Estates; o painel seguia mostrando Border.
		{"[t][324]LogGameState: Match State Changed from RoundActive to LeavingMap", Ready},
		{"[t][324]LogGameMode: ProcessServerTravel: /Estates/Maps/Estates?Scenario=Scenario_Estates_Checkpoint_Insurgency", Traveling},
		{"[t][390]LogGameState: Match State Changed from EnteringMap to LoadingAssets", Traveling},
		{"[t][402]LogGameState: Match State Changed from LoadingAssets to WaitingToStart", Ready},
	}
	for _, st := range steps {
		appendLog(t, dir, st.line)
		s.Tick()
		if got := s.Status().Phase; got != st.want {
			t.Fatalf("after %q: phase %s, want %s", st.line, got, st.want)
		}
	}
	if st := s.Status(); st.Map != "Estates" || st.Scenario != "Scenario_Estates_Checkpoint_Insurgency" {
		t.Errorf("mapa: %s %s", st.Map, st.Scenario)
	}
}

// Servidor aberto fora do app, sem viagem para o mapa de mod: fica pronto no mapa em que abriu.
func TestExternalServerWithoutTravelBecomesReady(t *testing.T) {
	s, f, dir := setup(t)
	appendLog(t, dir,
		"[t][  0]LogLoad: LoadMap: /Game/Maps/Farmhouse/Farmhouse?Name=Player?Scenario=Scenario_Farmhouse_Checkpoint_Security?MaxPlayers=10",
		"[t][  0]LogGameState: Match State Changed from EnteringMap to LoadingAssets",
		"[t][  4]LogGameState: Match State Changed from LoadingAssets to WaitingToStart",
	)
	f.mu.Lock()
	f.running, f.started = true, time.Now().Add(-time.Minute)
	f.mu.Unlock()
	s.Tick()
	if st := s.Status(); st.Phase != Ready || !st.External || st.Map != "Farmhouse" {
		t.Errorf("%+v", st)
	}
}

// O erro de mods em conflito leva o objeto e o mod que o traz (ADR-039).
func TestModConflictNamesTheMod(t *testing.T) {
	s := &Service{d: Deps{ModOwners: func(obj string) []string {
		if obj == "/Game/Env/SM_X" {
			return []string{"Daybreak (#5748299)"}
		}
		return nil
	}}}
	s.apply(logparse.Event{Kind: logparse.Problem, Problem: problems.ModConflict, Object: "/Game/Env/SM_X"})
	p := s.st.Problem
	if p == nil || p.Code != problems.ModConflict || p.Params["object"] != "/Game/Env/SM_X" ||
		p.Params["mods"] != "Daybreak (#5748299)" {
		t.Fatalf("problem %+v", p)
	}
	s.apply(logparse.Event{Kind: logparse.Problem, Problem: problems.ModConflict, Object: "/Game/Other"})
	if _, ok := s.st.Problem.Params["mods"]; ok {
		t.Fatal("mods for an object no mod has")
	}
}
