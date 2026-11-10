package watchdog

import (
	"testing"
	"time"

	"sandstorm-server-assist/desktop/internal/domain/server"
)

type rig struct {
	now     time.Time
	phase   string
	profile server.Profile
	calls   []string
	said    []string
}

func (r *rig) service() *Service {
	return New(Deps{
		Config:      func() (server.Profile, server.Secrets, error) { return r.profile, server.Secrets{}, nil },
		Phase:       func() string { return r.phase },
		Start:       func() error { r.calls = append(r.calls, "start"); return nil },
		Restart:     func() error { r.calls = append(r.calls, "restart"); return nil },
		StopAndWait: func() error { r.calls = append(r.calls, "stop"); return nil },
		Update:      func() error { r.calls = append(r.calls, "update"); return nil },
		Say:         func(text string) { r.said = append(r.said, text) },
		Log:         func(string, string) {},
		Now:         func() time.Time { return r.now },
		Go:          func(fn func()) { fn() },
	})
}

func TestDailyRestart(t *testing.T) {
	r := &rig{now: time.Date(2026, 11, 10, 4, 50, 0, 0, time.Local), phase: PhaseReady, profile: server.Default()}
	r.profile.AutoRestart = server.AutoRestart{Daily: "05:00", WarnMessage: "Restart in {minutes} min", Update: true}
	s := r.service()
	s.Tick()
	if len(r.calls) != 0 || len(r.said) != 0 {
		t.Fatal(r.calls, r.said)
	}
	r.now = r.now.Add(5 * time.Minute)
	s.Tick()
	s.Tick()
	if len(r.said) != 1 || r.said[0] != "Restart in 5 min" {
		t.Fatal(r.said)
	}
	r.now = r.now.Add(4 * time.Minute)
	s.Tick()
	if len(r.said) != 2 || r.said[1] != "Restart in 1 min" {
		t.Fatal(r.said)
	}
	r.now = r.now.Add(time.Minute)
	s.Tick()
	want := []string{"stop", "update", "start"}
	if len(r.calls) != 3 || r.calls[0] != want[0] || r.calls[1] != want[1] || r.calls[2] != want[2] {
		t.Fatal(r.calls)
	}
	if next := s.Next(); next.Day() != 11 || next.Hour() != 5 {
		t.Fatal(next)
	}
	// Servidor parado na hora: não inicia nada.
	r.calls, r.phase = nil, PhaseStopped
	r.now = r.now.Add(24 * time.Hour)
	s.Tick()
	if len(r.calls) != 0 {
		t.Fatal(r.calls)
	}
}

func TestCrashRestart(t *testing.T) {
	r := &rig{now: time.Date(2026, 11, 10, 12, 0, 0, 0, time.UTC), phase: PhaseCrashed, profile: server.Default()}
	r.profile.AutoRestart.OnCrash = true
	s := r.service()
	for i := 0; i < CrashMax; i++ {
		r.phase = PhaseCrashed
		s.Tick()
		r.now = r.now.Add(CrashDelay)
		s.Tick()
		r.phase = "starting"
		s.Tick()
		r.now = r.now.Add(time.Minute)
	}
	if len(r.calls) != CrashMax {
		t.Fatal(r.calls)
	}
	// Quarta queda em 30 minutos: para de tentar.
	r.phase = PhaseCrashed
	s.Tick()
	r.now = r.now.Add(CrashDelay)
	s.Tick()
	if len(r.calls) != CrashMax {
		t.Fatal("restarted after the limit", r.calls)
	}
	// Desligado: não reinicia.
	r2 := &rig{now: r.now, phase: PhaseCrashed, profile: server.Default()}
	s2 := r2.service()
	s2.Tick()
	r2.now = r2.now.Add(time.Minute)
	s2.Tick()
	if len(r2.calls) != 0 {
		t.Fatal(r2.calls)
	}
}
