package runner

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"sandstorm-server-assist/desktop/internal/adapters/proc"
	"sandstorm-server-assist/desktop/internal/domain/problems"
	"sandstorm-server-assist/desktop/internal/domain/server"
)

// Duas instâncias na mesma pasta (ADR-048): cada runner só vê, para e acompanha os processos da sua.
func TestRunnerOnlyOwnInstance(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ServerExe), []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	var mu sync.Mutex
	procs := map[uint32]string{
		10: `Hideout -Port=27102`,
		11: `Insurgency Hideout -Port=27102`,
		20: `Hideout -Port=27103 -saveddirsuffix=2`,
		30: `Hideout -saveddirsuffix=_other`,
	}
	var killed []uint32
	p := server.Default()
	p.ServerDir, p.Instance, p.UseMods = dir, 2, false
	s := New(Deps{
		ID: "two",
		Find: func(string) ([]proc.Process, error) {
			mu.Lock()
			defer mu.Unlock()
			var out []proc.Process
			for pid, cmd := range procs {
				out = append(out, proc.Process{PID: pid, CmdLine: cmd})
			}
			return out, nil
		},
		Kill: func(pid uint32) error {
			mu.Lock()
			defer mu.Unlock()
			killed = append(killed, pid)
			delete(procs, pid)
			return nil
		},
		Start:  func(string, string, string) error { return nil },
		Config: func() (server.Profile, server.Secrets, error) { return p, server.Secrets{}, nil },
	})
	s.Tick()
	st := s.Status()
	if st.Server != "two" || len(st.PIDs) != 1 || st.PIDs[0] != 20 || !st.External {
		t.Fatalf("status: %+v", st)
	}
	if err := s.Start(); problems.CodeOf(err) != problems.ServerRunning {
		t.Fatalf("start with own instance running: %v", err)
	}
	if err := s.Stop(); err != nil {
		t.Fatal(err)
	}
	if len(killed) != 1 || killed[0] != 20 {
		t.Fatalf("killed %v (only instance 2 may stop)", killed)
	}
	// Com a instância 2 parada, o servidor principal rodando não impede o início.
	if err := s.Start(); err != nil {
		t.Fatalf("start next to the main server: %v", err)
	}
}

func TestRunnerCheckStartAndWaitForMods(t *testing.T) {
	s, f, _ := setup(t)
	s.d.CheckStart = func(server.Profile, server.Secrets) error {
		return problems.New(problems.PortsShared, "")
	}
	if err := s.Start(); problems.CodeOf(err) != problems.PortsShared || f.starts != 0 {
		t.Fatalf("check refused nothing: %v (%d starts)", err, f.starts)
	}
	s.d.CheckStart = nil

	busy := true
	s.d.ModsBusy = func() bool { return busy }
	if err := s.Start(); err != nil {
		t.Fatal(err)
	}
	if st := s.Status(); st.Phase != Waiting || f.starts != 0 {
		t.Fatalf("not waiting: %s (%d starts)", st.Phase, f.starts)
	}
	s.Tick()
	if s.Status().Phase != Waiting || f.starts != 0 {
		t.Fatal("started while the other server syncs mods")
	}
	busy = false
	s.Tick()
	if f.starts != 1 || s.Status().Phase != Starting {
		t.Fatalf("did not start after waiting: %s (%d starts)", s.Status().Phase, f.starts)
	}
	if !strings.Contains(f.cmdline, "-mods") {
		t.Errorf("command line: %s", f.cmdline)
	}
	// Enquanto inicia com mods, este servidor segura os outros (até o limite).
	if !s.ModsBusy(true, time.Minute) || s.ModsBusy(false, time.Minute) || s.ModsBusy(true, 0) {
		t.Error("ModsBusy while starting")
	}

	// Parar enquanto espera cancela a espera.
	_ = s.Stop()
	f.mu.Lock()
	f.running = false
	f.mu.Unlock()
	s.Tick()
	busy = true
	_ = s.Start()
	if s.Status().Phase != Waiting {
		t.Fatal("not waiting the second time")
	}
	if err := s.Stop(); err != nil || s.Status().Phase != Stopped {
		t.Fatalf("stop while waiting: %v %s", err, s.Status().Phase)
	}
	busy = false
	s.Tick()
	if f.starts != 1 {
		t.Error("started after the wait was canceled")
	}
}
