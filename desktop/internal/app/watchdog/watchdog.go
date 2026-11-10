// Package watchdog é o reinício automático do servidor (ADR-022): todo dia num horário (com aviso no chat e,
// se o usuário quiser, atualização pelo SteamCMD) e quando o servidor fecha sozinho.
package watchdog

import (
	"sync"
	"time"

	"sandstorm-server-assist/desktop/internal/domain/server"
)

// Fases do runner que o serviço usa.
const (
	PhaseReady   = "ready"
	PhaseCrashed = "crashed"
	PhaseStopped = "stopped"
)

// Limites do reinício depois de travar.
const (
	CrashDelay  = 10 * time.Second
	CrashWindow = 30 * time.Minute
	CrashMax    = 3
)

// Deps são as partes do sistema que o serviço usa (dublês nos testes).
type Deps struct {
	Config func() (server.Profile, server.Secrets, error)
	// Phase é a fase atual do servidor (runner).
	Phase func() string
	Start func() error
	// Restart para, espera e inicia de novo.
	Restart func() error
	// StopAndWait para e espera os processos sumirem.
	StopAndWait func() error
	// Update atualiza o servidor pelo SteamCMD e só volta quando termina.
	Update func() error
	// Say manda um texto para o chat (fila do serviço de chat).
	Say func(text string)
	Log func(text, level string)
	Now func() time.Time
	// Go roda a tarefa longa (nos testes, na hora).
	Go func(func())
}

// Service guarda o próximo reinício e as quedas recentes.
type Service struct {
	d       Deps
	mu      sync.Mutex
	clock   string
	dailyAt time.Time
	warned  map[int]bool
	busy    bool
	// Quedas: quando foram, se a atual já foi tratada e quando iniciar de novo.
	crashes []time.Time
	handled bool
	crashAt time.Time
	gaveUp  bool
}

// New cria o serviço.
func New(d Deps) *Service {
	if d.Now == nil {
		d.Now = time.Now
	}
	if d.Go == nil {
		d.Go = func(fn func()) { go fn() }
	}
	return &Service{d: d, warned: map[int]bool{}}
}

// Next devolve o próximo reinício diário (zero se desligado), para a tela.
func (s *Service) Next() time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.dailyAt
}

// Tick confere o horário e as quedas. A ponte chama a cada segundo.
func (s *Service) Tick() {
	p, _, err := s.d.Config()
	if err != nil {
		return
	}
	now, phase := s.d.Now(), s.d.Phase()
	s.crash(p, phase, now)
	s.daily(p, phase, now)
}

func (s *Service) crash(p server.Profile, phase string, now time.Time) {
	s.mu.Lock()
	if phase != PhaseCrashed || !p.AutoRestart.OnCrash {
		s.handled, s.crashAt = false, time.Time{}
		if phase == PhaseReady {
			s.gaveUp = false
		}
		s.mu.Unlock()
		return
	}
	if !s.handled {
		s.handled = true
		kept := s.crashes[:0]
		for _, t := range s.crashes {
			if now.Sub(t) < CrashWindow {
				kept = append(kept, t)
			}
		}
		s.crashes = kept
		if len(s.crashes) >= CrashMax {
			if !s.gaveUp {
				s.gaveUp = true
				s.mu.Unlock()
				s.d.Log("The server closed by itself 3 times in 30 minutes: automatic restart paused", "warning")
				return
			}
			s.mu.Unlock()
			return
		}
		s.crashAt = now.Add(CrashDelay)
		s.mu.Unlock()
		s.d.Log("The server closed by itself; starting it again in 10 seconds", "warning")
		return
	}
	due := !s.crashAt.IsZero() && !now.Before(s.crashAt) && !s.busy
	if due {
		s.crashAt = time.Time{}
		s.crashes = append(s.crashes, now)
		s.busy = true
	}
	s.mu.Unlock()
	if due {
		s.d.Go(func() {
			if err := s.d.Start(); err != nil {
				s.d.Log("Automatic restart failed: "+err.Error(), "warning")
			}
			s.done()
		})
	}
}

func (s *Service) done() {
	s.mu.Lock()
	s.busy = false
	s.mu.Unlock()
}

func (s *Service) daily(p server.Profile, phase string, now time.Time) {
	ar := p.AutoRestart
	s.mu.Lock()
	if ar.Daily == "" {
		s.clock, s.dailyAt = "", time.Time{}
		s.mu.Unlock()
		return
	}
	if ar.Daily != s.clock || s.dailyAt.IsZero() {
		s.clock = ar.Daily
		s.dailyAt, _ = server.NextDaily(now, ar.Daily)
		clear(s.warned)
	}
	at := s.dailyAt
	running := phase != PhaseStopped && phase != PhaseCrashed && phase != "not-installed"
	var say []string
	if phase == PhaseReady && ar.WarnMessage != "" {
		for _, m := range []int{5, 1} {
			if !s.warned[m] && !now.Before(at.Add(-time.Duration(m)*time.Minute)) && now.Before(at) {
				s.warned[m] = true
				say = append(say, server.WarningText(ar.WarnMessage, m))
			}
		}
	}
	due := !now.Before(at)
	if due {
		s.dailyAt, _ = server.NextDaily(now, ar.Daily)
		clear(s.warned)
	}
	start := due && running && !s.busy
	if start {
		s.busy = true
	}
	s.mu.Unlock()

	for _, text := range say {
		s.d.Say(text)
	}
	if !start {
		return
	}
	s.d.Log("Daily restart", "info")
	s.d.Go(func() {
		defer s.done()
		if !ar.Update {
			if err := s.d.Restart(); err != nil {
				s.d.Log("Daily restart failed: "+err.Error(), "warning")
			}
			return
		}
		if err := s.d.StopAndWait(); err != nil {
			s.d.Log("Daily restart failed: "+err.Error(), "warning")
			return
		}
		if err := s.d.Update(); err != nil {
			s.d.Log("Server update failed; starting the installed version: "+err.Error(), "warning")
		}
		if err := s.d.Start(); err != nil {
			s.d.Log("Daily restart failed: "+err.Error(), "warning")
		}
	})
}
