// Package runner inicia, para e acompanha o servidor (ADR-003, ADR-004). O estado vem dos processos encontrados
// na pasta do servidor e dos eventos do Insurgency.log.
package runner

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"sandstorm-server-assist/desktop/internal/adapters/proc"
	"sandstorm-server-assist/desktop/internal/app/logbuf"
	"sandstorm-server-assist/desktop/internal/domain/logparse"
	"sandstorm-server-assist/desktop/internal/domain/problems"
	"sandstorm-server-assist/desktop/internal/domain/server"
)

// Phase é a fase mostrada no painel.
type Phase string

// Fases.
const (
	NotInstalled Phase = "not-installed"
	Stopped      Phase = "stopped"
	Starting     Phase = "starting"
	Modio        Phase = "modio"
	Mods         Phase = "mods"
	Traveling    Phase = "traveling"
	Ready        Phase = "ready"
	Stopping     Phase = "stopping"
	Crashed      Phase = "crashed"
)

// ModState é a situação de um mod nesta partida.
type ModState struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	State string `json:"state"` // active | installing | installed | updating | uninstalling | removed | pending
	Done  int64  `json:"done"`
	Total int64  `json:"total"`
}

// Status é o que o painel mostra.
type Status struct {
	Phase            Phase           `json:"phase"`
	PIDs             []uint32        `json:"pids"`
	Since            time.Time       `json:"since"`
	Map              string          `json:"map"`
	Scenario         string          `json:"scenario"`
	MatchState       string          `json:"matchState"`
	Port             int             `json:"port"`
	Mods             []ModState      `json:"mods"`
	Problem          *problems.Error `json:"problem,omitempty"`
	CmdLine          string          `json:"cmdLine"`
	External         bool            `json:"external"`
	Authorizing      bool            `json:"authorizing"`
	GameStatsMissing bool            `json:"gameStatsMissing"`
	// Remount: o app está reiniciando o servidor uma vez para montar os mods que ficaram de fora.
	Remount bool `json:"remount"`
}

// Run é o registro da partida iniciada pelo app, gravado em disco (ADR-026). Com ele, o app reaberto reconhece o
// servidor que ele mesmo iniciou e continua de onde parou, em vez de tratá-lo como aberto fora do app.
type Run struct {
	Started time.Time `json:"started"`
	Dir     string    `json:"dir"`
	CmdLine string    `json:"cmdLine"` // já com os segredos mascarados (Launch.Shown)
}

// Janela em que a hora de início do processo precisa cair para ser a partida registrada.
const (
	runBefore = 5 * time.Second
	runAfter  = 2 * time.Minute
)

// matches diz se os processos encontrados são os da partida registrada: mesma pasta e processo mais antigo
// iniciado logo depois do registro. Sem hora do Windows, não reconhece (fica como servidor de fora).
func (r *Run) matches(dir string, list []proc.Process) bool {
	if r == nil || r.Started.IsZero() || !strings.EqualFold(filepath.Clean(r.Dir), filepath.Clean(dir)) {
		return false
	}
	first := time.Time{}
	for _, p := range list {
		if !p.Started.IsZero() && (first.IsZero() || p.Started.Before(first)) {
			first = p.Started
		}
	}
	return !first.IsZero() && !first.Before(r.Started.Add(-runBefore)) && !first.After(r.Started.Add(runAfter))
}

// ServerExe é o executável do servidor dentro da pasta.
const ServerExe = "InsurgencyServer.exe"

// Deps são as partes do sistema que o runner usa (dublês nos testes).
type Deps struct {
	Find  func(dir string) ([]proc.Process, error)
	Kill  func(pid uint32) error
	Start func(exe, cmdline, dir string) error
	// Config devolve o perfil e os segredos atuais.
	Config func() (server.Profile, server.Secrets, error)
	// Prepare grava os arquivos de configuração (Admins.txt, MapCycle.txt) antes de iniciar.
	Prepare func(server.Profile) error
	// SessionExpiry devolve o vencimento da sessão do mod.io (zero se não houver).
	SessionExpiry func() time.Time
	// ModInfo devolve nome e tamanho final de um mod (state.json) e o tamanho atual da pasta.
	ModInfo func(id int) (name string, total, done int64)
	// InstalledMods devolve os ids dos mods instalados (state.json do mod.io).
	InstalledMods func() []int
	// CleanLogs apaga os logs antigos antes de iniciar, se o usuário ligou a limpeza (ADR-018). Pode ser nil.
	CleanLogs func(serverDir string)
	// OnPlayer recebe as entradas e saídas de jogadores lidas no log (mensagens no chat, ADR-020). Pode ser nil.
	OnPlayer func(logparse.Event)
	// SaveRun grava (nil apaga) o registro da partida iniciada pelo app; LoadRun o lê (ADR-026). Podem ser nil.
	SaveRun func(*Run)
	LoadRun func() *Run
	Log     *logbuf.Buffer
	Now     func() time.Time
}

// Service é o controle do servidor.
type Service struct {
	mu       sync.Mutex
	d        Deps
	st       Status
	tail     Tail
	tailInfo os.FileInfo
	notify   func(Status)
	stopping bool
	// Autorização do mod.io (ADR-003): partida com o código, depois volta sozinho para "none".
	authCode    string
	authBefore  time.Time
	authStarted time.Time
	travelMap   string
	mods        map[int]*ModState
	// ownStart diz se o servidor desta partida foi iniciado pelo app (só então ele reinicia sozinho).
	ownStart bool
	// added são os mods novos desta partida (PendingChanges Added). remountCheck pede a conferência dos mods
	// montados; remounted evita reiniciar mais de uma vez seguida.
	added        map[int]bool
	remountCheck bool
	remounted    bool
	// cleanPending: limpar os logs antigos assim que o servidor criar o log desta partida.
	cleanPending bool
}

// authTimeout é quanto o app espera a sessão nova do mod.io aparecer.
const authTimeout = 2 * time.Minute

// New cria o serviço.
func New(d Deps) *Service {
	if d.Now == nil {
		d.Now = time.Now
	}
	return &Service{d: d, st: Status{Phase: Stopped}, mods: map[int]*ModState{}, added: map[int]bool{}}
}

// SetNotify define quem recebe o estado novo.
func (s *Service) SetNotify(fn func(Status)) {
	s.mu.Lock()
	s.notify = fn
	s.mu.Unlock()
}

// Status devolve o estado atual.
func (s *Service) Status() Status {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.snapshot()
}

func (s *Service) snapshot() Status {
	st := s.st
	st.PIDs = append([]uint32(nil), s.st.PIDs...)
	st.Mods = make([]ModState, 0, len(s.mods))
	for _, m := range s.mods {
		st.Mods = append(st.Mods, *m)
	}
	sort.Slice(st.Mods, func(i, j int) bool { return st.Mods[i].Name < st.Mods[j].Name })
	st.Authorizing = s.authCode != ""
	return st
}

func (s *Service) emit() {
	st, fn := s.snapshot(), s.notify
	if fn != nil {
		go fn(st)
	}
}

// Start inicia o servidor no modo normal (-SecurityCode=none).
func (s *Service) Start() error {
	s.mu.Lock()
	s.remounted = false
	s.mu.Unlock()
	return s.start("")
}

// Authorize inicia o servidor com um código novo do mod.io. Quando a sessão nova aparece, o servidor é
// reiniciado sozinho com "none" (com código novo os mods já baixados ficam parados na fila).
func (s *Service) Authorize(code string) error {
	code = strings.TrimSpace(code)
	if !server.ValidCode(code) {
		return problems.Field("securityCode", "validation.code")
	}
	s.mu.Lock()
	s.remounted = false
	s.mu.Unlock()
	return s.start(code)
}

func (s *Service) start(code string) error {
	p, sec, err := s.d.Config()
	if err != nil {
		return err
	}
	exe := filepath.Join(p.ServerDir, ServerExe)
	if p.ServerDir == "" || !fileExists(exe) {
		return problems.New(problems.ServerNotInstalled, exe)
	}
	if running, _ := s.d.Find(p.ServerDir); len(running) > 0 {
		return problems.New(problems.ServerRunning, "")
	}
	launch, err := server.BuildLaunch(p, sec, codeOrNone(code))
	if err != nil {
		return err
	}
	if s.d.Prepare != nil {
		if err := s.d.Prepare(p); err != nil {
			return problems.Wrap(problems.ConfigWrite, err)
		}
	}

	s.mu.Lock()
	s.tail = Tail{Path: logPath(p.ServerDir)}
	s.tail.SkipToEnd(0) // o arquivo antigo vira backup; o novo é detectado pela troca de arquivo
	s.tailInfo, _ = os.Stat(s.tail.Path)
	s.mods = map[int]*ModState{}
	s.travelMap = p.TravelMap
	s.stopping = false
	s.ownStart = true
	s.cleanPending = true
	s.added = map[int]bool{}
	s.remountCheck = false
	s.st = Status{Phase: Starting, Since: s.d.Now(), CmdLine: launch.Shown(sec)}
	if code != "" {
		s.authCode = code
		s.authStarted = s.d.Now()
		if s.d.SessionExpiry != nil {
			s.authBefore = s.d.SessionExpiry()
		}
	} else {
		s.authCode = ""
	}
	s.mu.Unlock()

	if s.d.Log != nil {
		s.d.Log.StartRun(s.d.Now())
	}
	s.log("app", "Starting server: "+launch.Shown(sec))
	if err := s.d.Start(exe, launch.CmdLine, p.ServerDir); err != nil {
		s.mu.Lock()
		s.st.Phase = Stopped
		s.authCode = ""
		s.mu.Unlock()
		s.emitLocked()
		return problems.Wrap(problems.ServerStartFailed, err)
	}
	if s.d.SaveRun != nil {
		s.mu.Lock()
		run := &Run{Started: s.st.Since, Dir: p.ServerDir, CmdLine: s.st.CmdLine}
		s.mu.Unlock()
		s.d.SaveRun(run)
	}
	s.emitLocked()
	return nil
}

func codeOrNone(code string) string {
	if code == "" {
		return server.SecurityNone
	}
	return code
}

// Stop encerra os processos do servidor desta pasta.
func (s *Service) Stop() error {
	p, _, err := s.d.Config()
	if err != nil {
		return err
	}
	list, err := s.d.Find(p.ServerDir)
	if err != nil {
		return problems.Wrap(problems.ServerStopFailed, err)
	}
	if len(list) == 0 {
		s.mu.Lock()
		s.st.Phase = Stopped
		s.mu.Unlock()
		s.emitLocked()
		return nil
	}
	s.mu.Lock()
	s.stopping = true
	s.authCode = ""
	s.st.Phase = Stopping
	s.st.Remount = false
	s.mu.Unlock()
	s.emitLocked()
	s.log("app", "Stopping server")
	var errs []error
	for _, pr := range list {
		if proc.IsSelf(pr.PID) {
			continue
		}
		if err := s.d.Kill(pr.PID); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return problems.Wrap(problems.ServerStopFailed, errors.Join(errs...))
	}
	return nil
}

// Restart para, espera os processos sumirem e inicia de novo.
func (s *Service) Restart() error {
	if err := s.StopAndWait(); err != nil {
		return err
	}
	return s.Start()
}

// StopAndWait para o servidor e espera os processos sumirem (reinício com atualização, ADR-022).
func (s *Service) StopAndWait() error {
	if err := s.Stop(); err != nil {
		return err
	}
	return s.waitGone(30 * time.Second)
}

func (s *Service) waitGone(limit time.Duration) error {
	p, _, err := s.d.Config()
	if err != nil {
		return err
	}
	deadline := s.d.Now().Add(limit)
	for s.d.Now().Before(deadline) {
		if list, _ := s.d.Find(p.ServerDir); len(list) == 0 {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return problems.New(problems.ServerStopFailed, "timeout")
}

func (s *Service) emitLocked() {
	s.mu.Lock()
	s.emit()
	s.mu.Unlock()
}

func (s *Service) log(source, text string) {
	if s.d.Log != nil {
		s.d.Log.Add(source, text, "")
	}
}

// Tick confere os processos, lê o log novo e atualiza o estado. A ponte chama a cada segundo.
func (s *Service) Tick() {
	p, _, err := s.d.Config()
	if err != nil {
		return
	}
	if p.ServerDir == "" {
		s.mu.Lock()
		changed := s.st.Phase != NotInstalled
		s.st.Phase = NotInstalled
		s.mu.Unlock()
		if changed {
			s.emitLocked()
		}
		return
	}
	installed := fileExists(filepath.Join(p.ServerDir, ServerExe))
	list, _ := s.d.Find(p.ServerDir)

	s.mu.Lock()
	before := s.st.Phase
	running := len(list) > 0
	pids := make([]uint32, 0, len(list))
	for _, pr := range list {
		pids = append(pids, pr.PID)
	}
	s.st.PIDs = pids
	clearRun := false

	switch {
	case !installed && !running:
		s.st.Phase = NotInstalled
	case running && (s.st.Phase == Stopped || s.st.Phase == NotInstalled || s.st.Phase == Crashed):
		// Servidor que já estava rodando quando o app abriu: acompanha lendo o log atual desde o início.
		// "Em execução desde" é a hora em que o processo começou, não a hora em que o app o viu.
		// Se é a partida que o próprio app iniciou antes de fechar (registro em disco, ADR-026), continua como dele.
		var run *Run
		if s.d.LoadRun != nil {
			run = s.d.LoadRun()
		}
		own := run.matches(p.ServerDir, list)
		s.st = Status{Phase: Starting, Since: startedAt(list, s.d.Now()), External: !own, PIDs: pids}
		if own {
			s.st.CmdLine = run.CmdLine
		}
		if s.d.Log != nil {
			s.d.Log.StartRun(s.st.Since)
		}
		s.mods = map[int]*ModState{}
		s.ownStart = own
		s.added = map[int]bool{}
		s.remountCheck = false
		s.travelMap = p.TravelMap
		s.tail = Tail{Path: logPath(p.ServerDir)}
		s.tailInfo, _ = os.Stat(s.tail.Path)
	case !running && s.st.Phase != Stopped && s.st.Phase != NotInstalled && s.st.Phase != Crashed:
		if s.stopping || s.st.Phase == Starting && s.authCode == "" && s.d.Now().Sub(s.st.Since) < 2*time.Second {
			s.st.Phase = Stopped
		} else if s.st.Phase != Stopping {
			s.st.Phase = Crashed
			if s.st.Problem == nil {
				s.st.Problem = problems.New(problems.ServerCrashed, "")
			}
		} else {
			s.st.Phase = Stopped
		}
		s.stopping = false
		s.authCode = ""
		clearRun = true // a partida acabou: o registro não vale mais
	case !running && installed && s.st.Phase == NotInstalled:
		s.st.Phase = Stopped
	}
	active := running || s.st.Phase == Crashed
	s.mu.Unlock()
	if clearRun && s.d.SaveRun != nil {
		s.d.SaveRun(nil)
	}

	if active {
		s.readLog()
	}
	s.updateModProgress()
	s.checkAuthorization()
	s.checkRemount()
	s.cleanLogs(p.ServerDir)

	s.mu.Lock()
	changed := before != s.st.Phase
	s.mu.Unlock()
	// Sempre emite enquanto roda (mapa, mods e progresso mudam sem trocar de fase).
	if changed || active {
		s.emitLocked()
	}
}

func (s *Service) readLog() {
	s.mu.Lock()
	// Arquivo trocado (servidor reiniciou e renomeou o antigo): recomeça do início.
	if info, err := os.Stat(s.tail.Path); err == nil {
		if s.tailInfo != nil && !os.SameFile(s.tailInfo, info) {
			s.tail.Reset()
		}
		s.tailInfo = info
	}
	tail := &s.tail
	s.mu.Unlock()

	for i := 0; i < 8; i++ { // até 32 MB por tique
		lines, err := tail.Read()
		if err != nil || len(lines) == 0 {
			return
		}
		for _, line := range lines {
			s.log("server", line)
			if ev, ok := logparse.Parse(line); ok {
				s.apply(ev)
				if isPlayer(ev.Kind) && s.d.OnPlayer != nil && s.Status().Phase != Stopping {
					s.d.OnPlayer(ev)
				}
			}
		}
	}
}

// isPlayer diz se o evento é de entrada ou saída de jogador.
func isPlayer(k logparse.Kind) bool {
	return k == logparse.PlayerLogin || k == logparse.PlayerJoined || k == logparse.PlayerLeft
}

// readyStates são estados da partida que não são de carregamento.
func readyState(st string) bool {
	switch st {
	case "EnteringMap", "LoadingAssets", "LeavingMap", "":
		return false
	}
	return true
}

func (s *Service) apply(ev logparse.Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.st.Phase == Stopping || s.st.Phase == Stopped {
		return
	}
	switch ev.Kind {
	case logparse.ModioReady:
		if s.st.Phase == Starting {
			s.st.Phase = Modio
		}
	case logparse.Listening:
		s.st.Port = ev.Port
	case logparse.ModsEnabled:
		s.st.Phase = Mods
	case logparse.ModActivated:
		m := s.mod(ev.ModID)
		m.Name, m.State = ev.ModName, "active"
	case logparse.ModEvent:
		m := s.mod(ev.ModID)
		switch ev.Action {
		case "BeginInstall":
			m.State = "installing"
		case "Installed":
			m.State = "installed"
		case "BeginUpdate":
			m.State = "updating"
		case "Updated":
			m.State = "installed"
		case "BeginUninstall":
			m.State = "uninstalling"
		case "Uninstalled":
			m.State = "removed"
		}
	case logparse.ModPending:
		m := s.mod(ev.ModID)
		if ev.Action == "Removed" {
			m.State = "removed"
		} else if m.State == "" {
			m.State = "pending"
		}
		if ev.Action == "Added" {
			s.added[ev.ModID] = true
		}
	case logparse.Travel:
		s.st.Phase = Traveling
		s.st.Map, s.st.Scenario = ev.Map, ev.Scenario
		// Com mods novos nesta partida, o servidor monta só os que baixou e viaja sem os já instalados.
		if len(s.added) > 0 {
			s.remountCheck = true
		}
	case logparse.LoadMap:
		s.st.Map, s.st.Scenario = ev.Map, ev.Scenario
	case logparse.MatchState:
		s.st.MatchState = ev.To
		waitingTravel := s.travelMap != "" && !strings.EqualFold(s.st.Map, s.travelMap) && s.st.Phase != Ready
		if readyState(ev.To) && !waitingTravel {
			s.st.Phase = Ready
			s.st.Remount = false
		}
	case logparse.Problem:
		s.st.Problem = problems.New(ev.Problem, "")
	case logparse.GameStatsNoKey:
		s.st.GameStatsMissing = true
	}
}

func (s *Service) mod(id int) *ModState {
	m, ok := s.mods[id]
	if !ok {
		m = &ModState{ID: id}
		s.mods[id] = m
	}
	return m
}

// updateModProgress mede a pasta dos mods sendo baixados.
func (s *Service) updateModProgress() {
	if s.d.ModInfo == nil {
		return
	}
	s.mu.Lock()
	var ids []int
	for id, m := range s.mods {
		if m.State == "installing" || m.State == "updating" || m.Name == "" {
			ids = append(ids, id)
		}
	}
	s.mu.Unlock()
	for _, id := range ids {
		name, total, done := s.d.ModInfo(id)
		s.mu.Lock()
		if m := s.mods[id]; m != nil {
			if m.Name == "" {
				m.Name = name
			}
			m.Total, m.Done = total, done
		}
		s.mu.Unlock()
	}
}

// checkRemount: quando a partida baixou mods novos, o servidor monta só esses e viaja sem os que já estavam
// instalados (visto no log real de 09/11/2026: o mapa do mod e o mutator ficaram de fora). Na partida seguinte,
// sem nada pendente, ele monta todos. Se faltou algum, o app reinicia uma vez (servidor dele) ou avisa (servidor
// aberto fora do app).
func (s *Service) checkRemount() {
	s.mu.Lock()
	if !s.remountCheck || s.authCode != "" || s.st.Phase == Stopping || s.st.Phase == Stopped {
		s.mu.Unlock()
		return
	}
	s.remountCheck = false
	mounted := map[int]bool{}
	for id, m := range s.mods {
		if m.State == "active" {
			mounted[id] = true
		}
	}
	s.mu.Unlock()

	if s.d.InstalledMods == nil {
		return
	}
	missing := 0
	for _, id := range s.d.InstalledMods() {
		if !mounted[id] {
			missing++
		}
	}
	if missing == 0 {
		return
	}

	s.mu.Lock()
	own, again := s.ownStart, !s.remounted
	if !own || !again {
		s.st.Problem = problems.New(problems.ModsNotMounted, "")
		s.mu.Unlock()
		return
	}
	s.remounted = true
	s.st.Remount = true
	s.mu.Unlock()
	s.log("app", "New mods were installed and the server did not mount the others; restarting once to mount all mods")
	go func() {
		err := s.Stop()
		if err == nil {
			err = s.waitGone(30 * time.Second)
		}
		if err == nil {
			err = s.start("")
		}
		s.mu.Lock()
		s.st.Remount = err == nil
		if err != nil {
			s.st.Problem = problems.New(problems.ModsNotMounted, err.Error())
		}
		s.mu.Unlock()
		s.emitLocked()
		if err != nil {
			s.log("app", "restart to mount mods failed: "+err.Error())
		}
	}()
}

// checkAuthorization: com código novo, espera a sessão do mod.io mudar e reinicia em "none".
func (s *Service) checkAuthorization() {
	s.mu.Lock()
	code, before, started := s.authCode, s.authBefore, s.authStarted
	problem := s.st.Problem
	s.mu.Unlock()
	if code == "" || s.d.SessionExpiry == nil {
		return
	}
	if problem != nil && problem.Code == problems.ModioCodeInvalid {
		s.mu.Lock()
		s.authCode = ""
		s.mu.Unlock()
		s.log("app", "mod.io rejected the security code")
		_ = s.Stop()
		return
	}
	exp := s.d.SessionExpiry()
	if exp.After(before) && exp.After(s.d.Now()) {
		s.log("app", "mod.io authorized; restarting with -SecurityCode=none")
		s.mu.Lock()
		s.authCode = ""
		s.mu.Unlock()
		go func() {
			if err := s.Restart(); err != nil {
				s.log("app", "restart after authorization failed: "+err.Error())
			}
		}()
		return
	}
	if s.d.Now().Sub(started) > authTimeout {
		s.mu.Lock()
		s.authCode = ""
		s.st.Problem = problems.New(problems.ModioNotAuthorized, "timeout")
		s.mu.Unlock()
		_ = s.Stop()
	}
}

func logPath(serverDir string) string {
	return filepath.Join(serverDir, "Insurgency", "Saved", "Logs", "Insurgency.log")
}

// LogPath é o caminho do log principal do servidor.
func LogPath(serverDir string) string { return logPath(serverDir) }

// startedAt é a hora de início mais antiga entre os processos do servidor (o lançador abre antes do jogo).
// Sem hora informada pelo Windows, ou com hora no futuro, usa now.
func startedAt(list []proc.Process, now time.Time) time.Time {
	var first time.Time
	for _, p := range list {
		if !p.Started.IsZero() && !p.Started.After(now) && (first.IsZero() || p.Started.Before(first)) {
			first = p.Started
		}
	}
	if first.IsZero() {
		return now
	}
	return first
}

func fileExists(p string) bool {
	info, err := os.Stat(p)
	return err == nil && !info.IsDir()
}

// cleanLogs apaga os logs antigos uma vez por partida iniciada pelo app, depois que o servidor já trocou o log
// (fase mod.io em diante): assim a partida atual conta entre as mantidas.
func (s *Service) cleanLogs(serverDir string) {
	s.mu.Lock()
	ready := s.cleanPending && s.st.Phase != Starting && s.st.Phase != Stopped && s.st.Phase != Stopping
	if ready {
		s.cleanPending = false
	}
	s.mu.Unlock()
	if ready && s.d.CleanLogs != nil {
		s.d.CleanLogs(serverDir)
	}
}
