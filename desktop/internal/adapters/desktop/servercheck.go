package desktop

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"sandstorm-server-assist/desktop/internal/app/installer"
	"sandstorm-server-assist/desktop/internal/domain/steamcmd"
)

// Procura de versão nova do servidor dedicado: a primeira 2 min depois de abrir (depois do início automático e da
// procura de atualização do app), depois a cada 6 horas. Com o SteamCMD ocupado (instalação), tenta de novo em
// 10 min. O SteamCMD não precisa de aviso: ele se atualiza sozinho cada vez que roda, inclusive nesta procura.
const (
	serverCheckFirst = 2 * time.Minute
	serverCheckEvery = 6 * time.Hour
	serverCheckRetry = 10 * time.Minute
)

// ServerUpdate é o aviso de versão nova do servidor dedicado (barra de título).
type ServerUpdate struct {
	// Installed é o buildid instalado (appmanifest), Latest o publicado na Steam (vazio antes da primeira procura).
	Installed string    `json:"installed"`
	Latest    string    `json:"latest"`
	Available bool      `json:"available"`
	CheckedAt time.Time `json:"checkedAt"`
}

// serverCheck guarda a última versão publicada que o SteamCMD informou.
type serverCheck struct {
	mu     sync.Mutex
	latest string
	at     time.Time
}

func (c *serverCheck) set(latest string, at time.Time) {
	c.mu.Lock()
	c.latest, c.at = latest, at
	c.mu.Unlock()
}

func (c *serverCheck) get() (string, time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.latest, c.at
}

// serverUpdate compara a versão publicada com a instalada agora: depois de atualizar, o aviso some sozinho.
func (a *app) serverUpdate() ServerUpdate {
	in := installer.Inspect(a.cfg.Profile())
	latest, at := a.srvCheck.get()
	u := ServerUpdate{Latest: latest, CheckedAt: at}
	if in.ServerInstalled {
		u.Installed = in.BuildID
		u.Available = steamcmd.NewerBuild(in.BuildID, latest)
	}
	return u
}

// checkServerVersion pergunta ao SteamCMD a versão publicada. Volta true se o SteamCMD estava ocupado.
func (a *app) checkServerVersion(ctx context.Context) (busy bool) {
	p := a.cfg.Profile()
	in := installer.Inspect(p)
	if !in.SteamCMDFound || !in.ServerInstalled {
		return false
	}
	latest, err := a.install.LatestBuild(ctx, p.SteamCMDDir)
	if errors.Is(err, installer.ErrSteamBusy) {
		return true
	}
	if err != nil {
		return false
	}
	a.srvCheck.set(latest, time.Now())
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "server:update", a.serverUpdate())
	}
	return false
}

// runServerCheck procura versão nova do servidor até o app fechar.
func (a *app) runServerCheck(ctx context.Context) {
	wait := serverCheckFirst
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(wait):
		}
		wait = serverCheckEvery
		if a.checkServerVersion(ctx) {
			wait = serverCheckRetry
		}
	}
}

// ServerUpdate diz se há versão nova do servidor dedicado (evento "server:update" a cada procura).
func (i *InstallAPI) ServerUpdate() ServerUpdate { return i.a.serverUpdate() }
