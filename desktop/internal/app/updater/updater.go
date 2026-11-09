// Package updater é a atualização automática (ADR-014, mesmo modelo do midia-multi-tools): procura uma versão
// nova no servidor, baixa o instalador em segundo plano conferindo a assinatura do manifesto e o SHA-256, e
// instala quando o usuário pede (botão "Atualizar") ou ao fechar o app. O servidor do jogo não é afetado: ele
// roda em outro processo e continua no ar enquanto o app se atualiza.
package updater

import (
	"context"
	"crypto/ed25519"
	"encoding/json"
	"errors"
	"path/filepath"
	"sync"
	"time"

	"sandstorm-server-assist/desktop/internal/domain/problems"
	"sandstorm-server-assist/desktop/internal/domain/update"
)

// Status da atualização, para a tela.
type Status string

const (
	StatusIdle        Status = "idle"        // ainda não procurou
	StatusChecking    Status = "checking"    // procurando
	StatusUpToDate    Status = "uptodate"    // já está na versão mais nova
	StatusDownloading Status = "downloading" // baixando a versão nova
	StatusReady       Status = "ready"       // baixada e conferida: pronta para instalar
	StatusManual      Status = "manual"      // versão nova, mas o app é portátil: baixar pela página
	StatusDev         Status = "dev"         // build de desenvolvimento: sem atualização automática
	StatusError       Status = "error"       // falhou (tenta de novo na próxima vez)
)

// Stage diz onde a atualização falhou (a tela escolhe o texto).
const (
	StageCheck    = "check"
	StageDownload = "download"
)

// State é o que a tela mostra. O Go não tem texto: a tela traduz Status e Stage (ADR-009).
type State struct {
	Status   Status  `json:"status"`
	Current  string  `json:"current"`
	Version  string  `json:"version"`  // versão nova (quando há)
	Progress float64 `json:"progress"` // 0 a 1, durante o download
	Stage    string  `json:"stage"`    // em erro: check ou download
}

// Fetcher busca o manifesto e o instalador no servidor.
type Fetcher interface {
	// Manifest devolve o conteúdo do latest.json.
	Manifest(ctx context.Context) ([]byte, error)
	// Download baixa o arquivo para dst e devolve o SHA-256 (hex) e o tamanho.
	Download(ctx context.Context, file, dst string, onProgress func(done int64)) (sha string, size int64, err error)
}

// Files é o pouco de disco que a atualização usa.
type Files interface {
	Rename(from, to string) error
	Remove(path string) error
}

// Launcher inicia o instalador num processo separado (que continua depois que o app fecha).
type Launcher interface {
	Start(path string, args ...string) error
}

// Config junta as dependências.
type Config struct {
	Current   string            // versão deste app
	Installed bool              // true se veio do instalador (há como se atualizar sozinho)
	Dir       string            // pasta dos downloads de atualização
	PublicKey ed25519.PublicKey // chave pública da publicação
	Fetcher   Fetcher
	Files     Files
	Launcher  Launcher
	Notify    func(State)
}

// Service controla a atualização.
type Service struct {
	cfg Config

	mu        sync.Mutex
	state     State
	installer string // caminho do instalador pronto
	busy      bool
	launched  bool
}

// New cria o serviço. Sem fetcher (endereço ou pasta indisponível) ou em build de desenvolvimento, fica parado.
func New(cfg Config) *Service {
	if cfg.Notify == nil {
		cfg.Notify = func(State) {}
	}
	st := State{Status: StatusIdle, Current: cfg.Current}
	if !update.IsRelease(cfg.Current) || cfg.Fetcher == nil {
		st.Status = StatusDev
	}
	return &Service{cfg: cfg, state: st}
}

// SetNotify liga o serviço à tela (eventos "update:state").
func (s *Service) SetNotify(f func(State)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if f == nil {
		f = func(State) {}
	}
	s.cfg.Notify = f
}

// State devolve a situação atual.
func (s *Service) State() State {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.state
}

func (s *Service) set(st State) {
	s.mu.Lock()
	st.Current = s.cfg.Current
	s.state = st
	notify := s.cfg.Notify
	s.mu.Unlock()
	notify(st)
}

// Run procura atualizações depois de delay e depois a cada interval, até o ctx acabar.
func (s *Service) Run(ctx context.Context, delay, interval time.Duration) {
	if s.State().Status == StatusDev {
		return
	}
	t := time.NewTimer(delay)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			s.Check(ctx)
			t.Reset(interval)
		}
	}
}

// Check procura (e, se houver, baixa) uma versão nova. Devolve a situação ao terminar.
func (s *Service) Check(ctx context.Context) State {
	s.mu.Lock()
	if s.state.Status == StatusDev || s.busy || s.state.Status == StatusReady {
		st := s.state
		s.mu.Unlock()
		return st
	}
	s.busy = true
	s.mu.Unlock()
	defer func() { s.mu.Lock(); s.busy = false; s.mu.Unlock() }()

	s.set(State{Status: StatusChecking})
	raw, err := s.cfg.Fetcher.Manifest(ctx)
	if err != nil {
		s.set(State{Status: StatusError, Stage: StageCheck})
		return s.State()
	}
	var m update.Manifest
	if err := json.Unmarshal(raw, &m); err != nil || m.Verify(s.cfg.PublicKey) != nil {
		s.set(State{Status: StatusError, Stage: StageCheck})
		return s.State()
	}
	if !update.Newer(s.cfg.Current, m.Version) {
		s.set(State{Status: StatusUpToDate})
		return s.State()
	}
	if !s.cfg.Installed {
		s.set(State{Status: StatusManual, Version: m.Version})
		return s.State()
	}

	final := filepath.Join(s.cfg.Dir, m.File)
	part := final + ".part"
	_ = s.cfg.Files.Remove(part)
	s.set(State{Status: StatusDownloading, Version: m.Version})
	var last time.Time
	sha, size, err := s.cfg.Fetcher.Download(ctx, m.File, part, func(done int64) {
		if time.Since(last) < 300*time.Millisecond {
			return
		}
		last = time.Now()
		s.set(State{Status: StatusDownloading, Version: m.Version, Progress: float64(done) / float64(m.Size)})
	})
	if err == nil && (sha != m.SHA256 || size != m.Size) {
		err = errors.New("downloaded file does not match the manifest")
	}
	if err == nil {
		err = s.cfg.Files.Rename(part, final)
	}
	if err != nil {
		_ = s.cfg.Files.Remove(part)
		s.set(State{Status: StatusError, Version: m.Version, Stage: StageDownload})
		return s.State()
	}
	s.mu.Lock()
	s.installer = final
	s.mu.Unlock()
	s.set(State{Status: StatusReady, Version: m.Version, Progress: 1})
	return s.State()
}

// Install inicia o instalador em modo silencioso, que reabre o app no fim ("/relaunch").
// Quem chama fecha o app logo depois (o instalador espera o app sair).
func (s *Service) Install() error { return s.launch(true) }

// InstallOnExit instala a versão pronta enquanto o app fecha (sem reabrir), como no game-in.
func (s *Service) InstallOnExit() { _ = s.launch(false) }

func (s *Service) launch(relaunch bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.state.Status != StatusReady || s.installer == "" {
		return problems.New(problems.UpdateFailed, "no update ready to install")
	}
	if s.launched {
		return nil
	}
	args := []string{"/S"}
	if relaunch {
		args = append(args, "/relaunch")
	}
	if err := s.cfg.Launcher.Start(s.installer, args...); err != nil {
		return problems.Wrap(problems.UpdateFailed, err)
	}
	s.launched = true
	return nil
}
