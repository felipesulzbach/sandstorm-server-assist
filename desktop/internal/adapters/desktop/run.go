// Package desktop é a ponte com a tela (Wails). Única parte do Go que conhece o Wails. Os métodos exportados das
// *API são chamáveis pelo JavaScript: recebem ids e opções, e caminhos só vêm de diálogos abertos aqui.
package desktop

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"sandstorm-server-assist/desktop/internal/adapters/discord"
	"sandstorm-server-assist/desktop/internal/adapters/fsx"
	"sandstorm-server-assist/desktop/internal/adapters/modio"
	"sandstorm-server-assist/desktop/internal/adapters/proc"
	"sandstorm-server-assist/desktop/internal/adapters/rcon"
	upd "sandstorm-server-assist/desktop/internal/adapters/updater"
	"sandstorm-server-assist/desktop/internal/adapters/vault"
	"sandstorm-server-assist/desktop/internal/adapters/winsys"
	"sandstorm-server-assist/desktop/internal/app/chat"
	"sandstorm-server-assist/desktop/internal/app/config"
	"sandstorm-server-assist/desktop/internal/app/installer"
	"sandstorm-server-assist/desktop/internal/app/logbuf"
	"sandstorm-server-assist/desktop/internal/app/mapchange"
	"sandstorm-server-assist/desktop/internal/app/notify"
	"sandstorm-server-assist/desktop/internal/app/players"
	"sandstorm-server-assist/desktop/internal/app/runner"
	"sandstorm-server-assist/desktop/internal/app/updater"
	"sandstorm-server-assist/desktop/internal/app/watchdog"
	"sandstorm-server-assist/desktop/internal/domain/legal"
	"sandstorm-server-assist/desktop/internal/domain/logparse"
	"sandstorm-server-assist/desktop/internal/domain/server"
)

// AppVersion é trocada no build (-ldflags "-X .../desktop.AppVersion=x.y.z").
var AppVersion = "0.1.0-dev"

// app junta os serviços compartilhados pelas APIs.
type app struct {
	ctx     context.Context
	cfg     *config.Service
	log     *logbuf.Buffer
	runner  *runner.Service
	install *installer.Service
	mods    *modCache
	updates *updater.Service
	maps    *mapchange.Service
	chat    *chat.Service
	players *players.Service
	notify  *notify.Service
	watch   *watchdog.Service
	// lastPhase é a fase vista no último tique (avisos do Discord quando muda).
	lastPhase runner.Phase
	// autostart: o Windows abriu o app ao entrar (--autostart).
	autostart bool
	// updatesOnce liga a procura de atualizações uma vez, só depois do aceite dos termos (ADR-015).
	updatesOnce sync.Once
	// scenarios guarda os cenários que existem (jogo base e mods instalados) para conferir a rotação (ADR-030).
	scenarios *scenarioCache
}

// Run monta a aplicação e abre a janela. Só volta quando a janela fecha.
func Run(assets fs.FS) error {
	a := newApp()
	a.autostart = slices.Contains(os.Args[1:], winsys.AutostartFlag)
	startState := options.Normal
	if a.autostart {
		startState = options.Minimised
	}
	sys := &SystemAPI{a: a}
	cfgAPI := &ConfigAPI{a: a}
	inst := &InstallAPI{a: a}
	srv := &ServerAPI{a: a}
	logs := &LogAPI{a: a}
	reports := &ReportAPI{a: a}
	bak := &BackupAPI{a: a}

	return wails.Run(&options.App{
		Title:            "Sandstorm Server Assist",
		Width:            1200,
		Height:           800,
		MinWidth:         960,
		MinHeight:        620,
		BackgroundColour: &options.RGBA{R: 21, G: 20, B: 18, A: 255},
		// Sem a barra do Windows: a tela desenha a própria barra de título (como no midia-multi-tools).
		Frameless:        true,
		WindowStartState: startState,
		AssetServer:      &assetserver.Options{Assets: assets},
		OnStartup: func(ctx context.Context) {
			a.ctx = ctx
			a.log.SetNotify(func(lines []logbuf.Line) { runtime.EventsEmit(ctx, "log:lines", lines) })
			a.runner.SetNotify(func(st runner.Status) {
				// Servidor parou: a conexão das mensagens de chat fecha e a lista de jogadores recomeça.
				if st.Phase == runner.Stopped || st.Phase == runner.Crashed {
					a.chat.Reset()
				}
				runtime.EventsEmit(ctx, "server:status", st)
			})
			a.install.SetNotify(func(st installer.State) { runtime.EventsEmit(ctx, "install:state", st) })
			a.updates.SetNotify(func(st updater.State) { runtime.EventsEmit(ctx, "update:state", st) })
			if legal.Accepted(a.cfg.Settings().TermsAccepted) {
				a.startUpdates()
				a.autoStartServer()
			}
			a.syncAutostart()
			go a.loop(ctx)
		},
		OnBeforeClose: func(ctx context.Context) bool {
			if !a.install.Busy() {
				return false
			}
			answer, err := runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
				Type:    runtime.QuestionDialog,
				Title:   "Sandstorm Server Assist",
				Message: "An installation is in progress. Close and cancel it?",
			})
			return err != nil || (answer != "Yes" && answer != "Sim")
		},
		OnShutdown: func(context.Context) {
			a.install.Cancel()
			// Versão nova já baixada e o usuário fechou o app: instala agora (sem reabrir), como no game-in.
			a.updates.InstallOnExit()
		},
		Bind:    []any{sys, cfgAPI, inst, srv, logs, bak, reports},
		Windows: &windows.Options{Theme: windows.SystemDefault},
	})
}

func newApp() *app {
	v, err := vault.Default()
	var store config.SecretStore = memoryVault{}
	if err == nil {
		store = v
	}
	path, err := config.DefaultPath()
	if err != nil {
		path = filepath.Join(os.TempDir(), "sandstorm-server-assist", "config.json")
	}
	cfg := config.Open(path, store)
	log := logbuf.New(6000, cfg.SecretValues)
	mods := &modCache{}
	a := &app{cfg: cfg, log: log, install: installer.New(log), mods: mods, scenarios: &scenarioCache{}}
	a.updates = newUpdater()
	runs := newRunFile(path, log)
	a.runner = runner.New(runner.Deps{
		Find: func(dir string) ([]proc.Process, error) { return proc.FindInDir(dir, "InsurgencyServer") },
		Kill: proc.Kill,
		Start: func(exe, cmdline, dir string) error {
			_, err := proc.Start(exe, cmdline, dir, !cfg.Settings().ShowConsole)
			return err
		},
		Config:  cfg.Current,
		Prepare: config.Prepare,
		SessionExpiry: func() time.Time {
			uf, err := modio.UserFile()
			if err != nil {
				return time.Time{}
			}
			s, _, _ := modio.ReadSession(uf)
			return s.Expires
		},
		ModInfo: mods.info,
		CleanLogs: func(dir string) {
			st := cfg.Settings()
			if !st.LogCleanup {
				return
			}
			n, freed, err := config.DeleteOldLogs(dir, st.LogKeep)
			if n > 0 {
				log.Add("app", fmt.Sprintf("Old server logs deleted: %d files (%d KB), keeping the last %d runs", n, freed/1024, st.LogKeep), "")
			}
			if err != nil {
				log.Add("app", "Some old server logs could not be deleted: "+err.Error(), "warning")
			}
		},
		InstalledMods: func() []int {
			list, err := modio.ReadInstalled(modio.RootDir())
			if err != nil {
				return nil
			}
			ids := make([]int, 0, len(list))
			for _, m := range list {
				if m.State == modio.StateInstalled {
					ids = append(ids, m.ID)
				}
			}
			return ids
		},
		OnPlayer: func(ev logparse.Event) { a.chat.Observe(ev) },
		SaveRun:  runs.Save,
		LoadRun:  runs.Load,
		Log:      log,
	})
	a.maps = mapchange.New(mapchange.Deps{
		Config:       cfg.Current,
		Ready:        func() bool { return a.runner.Status().Phase == runner.Ready },
		ModScenarios: func() []string { return loadedModScenarios(a.runner.Status().Mods) },
		Dial: func(port int, password string) (mapchange.Conn, error) {
			c, err := rcon.Dial(port, password)
			if errors.Is(err, rcon.ErrAuth) {
				return nil, mapchange.ErrAuth
			}
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		Log: func(text, level string) { log.Add("app", text, level) },
	})
	a.notify = notify.New(notify.Deps{
		Config: cfg.Current,
		Post:   discord.Post,
		Log:    func(text, level string) { log.Add("app", text, level) },
	}, true)
	a.players = players.New(players.Deps{
		Config: cfg.Current,
		Ready:  func() bool { return a.runner.Status().Phase == runner.Ready },
		Dial: func(port int, password string) (players.Conn, error) {
			c, err := rcon.Dial(port, password)
			if errors.Is(err, rcon.ErrAuth) {
				return nil, players.ErrAuth
			}
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		Log: func(text, level string) { log.Add("app", text, level) },
	})
	a.watch = watchdog.New(watchdog.Deps{
		Config:      cfg.Current,
		Phase:       func() string { return string(a.runner.Status().Phase) },
		Start:       a.runner.Start,
		Restart:     a.runner.Restart,
		StopAndWait: a.runner.StopAndWait,
		Update:      a.updateServer,
		Say:         func(text string) { a.chat.Broadcast(text) },
		Log:         func(text, level string) { log.Add("app", text, level) },
	})
	a.chat = chat.New(chat.Deps{
		OnPlayer: func(kind, name string) { a.notify.Event(kind, "", name) },
		Config:   cfg.Current,
		Dial: func(port int, password string) (chat.Conn, error) {
			c, err := rcon.Dial(port, password)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		Log: func(text, level string) { log.Add("app", text, level) },
	}, true)
	return a
}

// Atualização automática: primeira procura 10 s depois de abrir, depois a cada 4 horas.
const (
	updateFirstCheck = 10 * time.Second
	updateInterval   = 4 * time.Hour
)

// startUpdates liga a procura de atualizações (uma vez).
func (a *app) startUpdates() {
	a.updatesOnce.Do(func() { go a.updates.Run(a.ctx, updateFirstCheck, updateInterval) })
}

// newUpdater monta a atualização automática. Sem pasta de cache, chave ou endereço válidos, ela fica desligada
// (o app funciona normalmente).
func newUpdater() *updater.Service {
	cfg := updater.Config{Current: AppVersion, Installed: upd.Installed(), PublicKey: upd.PublicKey(),
		Files: upd.OSFiles{}, Launcher: upd.Launcher{}}
	dir, err := upd.Dir()
	fetcher, ferr := upd.NewHTTPFetcher(upd.DefaultBaseURL)
	if err == nil && ferr == nil && cfg.PublicKey != nil {
		cfg.Dir, cfg.Fetcher = dir, fetcher
	}
	return updater.New(cfg)
}

// loop: estado do servidor a cada segundo e logs para a tela a cada 300 ms.
func (a *app) loop(ctx context.Context) {
	tick := time.NewTicker(time.Second)
	flush := time.NewTicker(300 * time.Millisecond)
	defer tick.Stop()
	defer flush.Stop()
	a.runner.Tick()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			a.runner.Tick()
			a.tickAutomation()
		case <-flush.C:
			a.log.Flush()
		}
	}
}

// memoryVault é usado se o Windows não der a pasta de configuração (segredos só na memória).
type memoryVault struct{}

func (memoryVault) Load() (server.Secrets, error) { return server.Secrets{}, nil }
func (memoryVault) Save(server.Secrets) error     { return nil }

// modCache guarda os mods do state.json (relido a cada 5 s) para nome e tamanho no progresso de download.
type modCache struct {
	mu   sync.Mutex
	at   time.Time
	byID map[int]modio.Installed
}

func (c *modCache) info(id int) (string, int64, int64) {
	root := modio.RootDir()
	c.mu.Lock()
	if time.Since(c.at) > 5*time.Second {
		c.byID = map[int]modio.Installed{}
		if list, err := modio.ReadInstalled(root); err == nil {
			for _, m := range list {
				c.byID[m.ID] = m
			}
		}
		c.at = time.Now()
	}
	m := c.byID[id]
	c.mu.Unlock()
	return m.Name, m.Size, fsx.DirSize(modio.ModDir(root, id))
}

// instanceID identifica a cópia única do app. A versão de desenvolvimento usa outro id, para abrir junto da
// instalada sem uma bloquear a outra.
func instanceID() string {
	const id = "sandstorm-server-assist-7d1f3c2a"
	if strings.Contains(AppVersion, "-") {
		return id + "-dev"
	}
	return id
}
