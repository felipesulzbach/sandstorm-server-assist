// Package desktop é a ponte com a tela (Wails). Única parte do Go que conhece o Wails. Os métodos exportados das
// *API são chamáveis pelo JavaScript: recebem ids e opções, e caminhos só vêm de diálogos abertos aqui.
package desktop

import (
	"context"
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

	"sandstorm-server-assist/desktop/internal/adapters/fsx"
	"sandstorm-server-assist/desktop/internal/adapters/modio"
	upd "sandstorm-server-assist/desktop/internal/adapters/updater"
	"sandstorm-server-assist/desktop/internal/adapters/vault"
	"sandstorm-server-assist/desktop/internal/adapters/winsys"
	"sandstorm-server-assist/desktop/internal/app/config"
	"sandstorm-server-assist/desktop/internal/app/installer"
	"sandstorm-server-assist/desktop/internal/app/modfiles"
	"sandstorm-server-assist/desktop/internal/app/updater"
	"sandstorm-server-assist/desktop/internal/domain/legal"
	"sandstorm-server-assist/desktop/internal/domain/server"
)

// AppVersion é trocada no build (-ldflags "-X .../desktop.AppVersion=x.y.z").
var AppVersion = "0.1.0-dev"

// AppReleased é a data de lançamento da versão (AAAA-MM-DD, da seção dela no CHANGELOG), trocada no build. Vazia
// na versão de desenvolvimento.
var AppReleased = ""

// app junta os serviços compartilhados pelas APIs. Os de cada servidor (runner, log, chat...) ficam em servers
// (ADR-048).
type app struct {
	ctx     context.Context
	cfg     *config.Service
	cfgPath string
	servers *fleet
	install *installer.Service
	mods    *modCache
	updates *updater.Service
	// autostart: o Windows abriu o app ao entrar (--autostart).
	autostart bool
	// updatesOnce liga a procura de atualizações uma vez, só depois do aceite dos termos (ADR-015).
	updatesOnce sync.Once
	// modFiles apaga mods baixados (ADR-037, ADR-047) com todos os servidores parados; modNew baixa os mods novos
	// reiniciando só o servidor escolhido (ADR-038).
	modFiles *modfiles.Service
	modNew   *modfiles.Service
	// scenarios guarda os cenários que existem (jogo base e mods instalados) para conferir a rotação (ADR-030).
	scenarios *scenarioCache
	// srvCheck guarda a versão publicada do servidor dedicado (aviso na barra de título).
	srvCheck serverCheck
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
			for _, sv := range a.all() {
				a.wireServer(sv)
			}
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
	a := &app{cfg: cfg, cfgPath: path, servers: &fleet{byID: map[string]*srv{}}, mods: &modCache{},
		scenarios: &scenarioCache{}}
	a.install = installer.New(fanLog{a})
	a.updates = newUpdater()
	a.all()
	a.modFiles = modfiles.New(a.modDeps(true))
	a.modNew = modfiles.New(a.modDeps(false))
	return a
}

// modDeps liga o serviço de mods aos servidores: fleet=true para e volta todos os que estavam rodando (apagar
// mods mexe na pasta de todos); fleet=false usa só o servidor escolhido.
func (a *app) modDeps(fleet bool) modfiles.Deps {
	logf := func(text, level string) { a.logAll("app", text, level) }
	d := modfiles.Deps{
		Busy:      a.install.Busy,
		Remove:    func(id int) error { return modio.RemoveInstalled(modio.RootDir(), id) },
		RemoveAll: func() error { return modio.RemoveAllInstalled(modio.RootDir()) },
		Logout:    modio.Logout,
		Log:       logf,
	}
	if !fleet {
		d.Running = func() bool { return len(a.sel().runner.Status().PIDs) > 0 }
		d.StopAndWait = func() error { return a.sel().runner.StopAndWait() }
		d.Start = func() error { return a.sel().runner.Start() }
		return d
	}
	var mu sync.Mutex
	var stopped []string
	d.Running = func() bool { return len(a.running("")) > 0 }
	d.StopAndWait = func() error {
		ids, err := a.stopAll()
		mu.Lock()
		stopped = ids
		mu.Unlock()
		return err
	}
	d.Start = func() error {
		mu.Lock()
		ids := stopped
		stopped = nil
		mu.Unlock()
		if len(ids) == 0 {
			ids = []string{a.cfg.SelectedID()}
		}
		return a.startIDs(ids)
	}
	return d
}

// Atualização automática: primeira procura 10 s depois de abrir, depois a cada 4 horas.
const (
	updateFirstCheck = 10 * time.Second
	updateInterval   = 4 * time.Hour
)

// startUpdates liga a procura de atualizações (uma vez).
func (a *app) startUpdates() {
	a.updatesOnce.Do(func() {
		go a.updates.Run(a.ctx, updateFirstCheck, updateInterval)
		go a.runServerCheck(a.ctx)
	})
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

// loop: estado dos servidores a cada segundo e logs para a tela a cada 300 ms.
func (a *app) loop(ctx context.Context) {
	tick := time.NewTicker(time.Second)
	flush := time.NewTicker(300 * time.Millisecond)
	defer tick.Stop()
	defer flush.Stop()
	for _, s := range a.all() {
		s.runner.Tick()
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			for _, s := range a.all() {
				s.runner.Tick()
				a.tickAutomation(s)
			}
		case <-flush.C:
			for _, s := range a.all() {
				s.log.Flush()
			}
		}
	}
}

// memoryVault é usado se o Windows não der a pasta de configuração (segredos só na memória).
type memoryVault struct{}

func (memoryVault) Load() (server.Vault, error) { return server.Vault{}, nil }
func (memoryVault) Save(server.Vault) error     { return nil }

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
		if list, err := modio.ServerInstalled(root); err == nil {
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
