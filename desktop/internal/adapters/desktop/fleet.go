package desktop

import (
	"errors"
	"fmt"
	"net"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"sandstorm-server-assist/desktop/internal/adapters/discord"
	"sandstorm-server-assist/desktop/internal/adapters/modio"
	"sandstorm-server-assist/desktop/internal/adapters/proc"
	"sandstorm-server-assist/desktop/internal/adapters/rcon"
	"sandstorm-server-assist/desktop/internal/app/chat"
	"sandstorm-server-assist/desktop/internal/app/config"
	"sandstorm-server-assist/desktop/internal/app/logbuf"
	"sandstorm-server-assist/desktop/internal/app/mapchange"
	"sandstorm-server-assist/desktop/internal/app/notify"
	"sandstorm-server-assist/desktop/internal/app/players"
	"sandstorm-server-assist/desktop/internal/app/runner"
	"sandstorm-server-assist/desktop/internal/app/watchdog"
	"sandstorm-server-assist/desktop/internal/domain/logparse"
	"sandstorm-server-assist/desktop/internal/domain/problems"
	"sandstorm-server-assist/desktop/internal/domain/server"
)

// Vários servidores ao mesmo tempo (ADR-048): cada servidor tem o próprio runner, log, chat, jogadores, troca de
// mapa, Discord e reinício automático. As telas usam o servidor escolhido (config.SelectedID); o Painel lista
// todos.

// srv são os serviços de um servidor.
type srv struct {
	id      string
	log     *logbuf.Buffer
	runner  *runner.Service
	maps    *mapchange.Service
	chat    *chat.Service
	players *players.Service
	notify  *notify.Service
	watch   *watchdog.Service
	// lastPhase é a fase vista no último tique (avisos do Discord quando muda).
	lastPhase runner.Phase
}

// fleet guarda os servidores pelo id.
type fleet struct {
	mu   sync.Mutex
	byID map[string]*srv
	// startLock: um servidor por vez entre conferir os mods e começar a iniciar.
	startLock sync.Mutex
}

// modsBusyLimit é quanto tempo um servidor iniciando segura os outros que esperam para sincronizar os mods.
const modsBusyLimit = 10 * time.Minute

// server devolve os serviços de um servidor (criados na primeira vez). nil se o id não existe na configuração.
func (a *app) server(id string) *srv {
	a.servers.mu.Lock()
	s := a.servers.byID[id]
	a.servers.mu.Unlock()
	if s != nil {
		return s
	}
	if a.cfg.Name(id) == "" {
		return nil
	}
	s = a.newServer(id)
	a.servers.mu.Lock()
	if old := a.servers.byID[id]; old != nil {
		s = old
	} else {
		a.servers.byID[id] = s
	}
	a.servers.mu.Unlock()
	return s
}

// sel devolve os serviços do servidor escolhido na tela (sempre existe: a configuração tem ao menos um).
func (a *app) sel() *srv {
	if a.cfg == nil { // testes da ponte sem configuração: o único servidor montado
		for _, s := range a.servers.byID {
			return s
		}
	}
	if s := a.server(a.cfg.SelectedID()); s != nil {
		return s
	}
	return a.all()[0]
}

// all devolve os servidores da configuração, em ordem de instância. Cria os novos; esquece os apagados.
func (a *app) all() []*srv {
	refs := a.cfg.Profiles()
	keep := map[string]bool{}
	out := make([]*srv, 0, len(refs))
	for _, r := range refs {
		keep[r.ID] = true
		if s := a.server(r.ID); s != nil {
			out = append(out, s)
		}
	}
	a.servers.mu.Lock()
	for id := range a.servers.byID {
		if !keep[id] {
			delete(a.servers.byID, id)
		}
	}
	a.servers.mu.Unlock()
	return out
}

// running devolve os servidores com processo rodando (menos except).
func (a *app) running(except string) []*srv {
	var out []*srv
	for _, s := range a.all() {
		if s.id != except && len(s.runner.Status().PIDs) > 0 {
			out = append(out, s)
		}
	}
	return out
}

// newServer monta os serviços de um servidor.
func (a *app) newServer(id string) *srv {
	cfg := a.cfg
	s := &srv{id: id, log: logbuf.New(6000, cfg.SecretValues)}
	conf := func() (server.Profile, server.Secrets, error) { return cfg.CurrentOf(id) }
	logf := func(text, level string) { s.log.Add("app", text, level) }
	runs := newRunFile(a.cfgPath, id, a.legacyRunFile(id), s.log)
	s.runner = runner.New(runner.Deps{
		ID:   id,
		Find: func(dir string) ([]proc.Process, error) { return proc.FindInDir(dir, "InsurgencyServer") },
		Kill: proc.Kill,
		Start: func(exe, cmdline, dir string) error {
			_, err := proc.Start(exe, cmdline, dir, !cfg.Settings().ShowConsole)
			return err
		},
		Config:     conf,
		Prepare:    config.Prepare,
		CheckStart: func(p server.Profile, sec server.Secrets) error { return a.checkStart(id, p, sec) },
		ModsBusy:   func() bool { return a.modsBusy(id) },
		StartLock:  &a.servers.startLock,
		SessionExpiry: func() time.Time {
			uf, err := modio.UserFile()
			if err != nil {
				return time.Time{}
			}
			sess, _, _ := modio.ReadSession(uf)
			return sess.Expires
		},
		ModInfo:   a.mods.info,
		ModOwners: modOwnerNames,
		CleanLogs: func(logsDir string) {
			st := cfg.Settings()
			if !st.LogCleanup {
				return
			}
			n, freed, err := config.DeleteOldLogs(logsDir, st.LogKeep)
			if n > 0 {
				logf(fmt.Sprintf("Old server logs deleted: %d files (%d KB), keeping the last %d runs", n, freed/1024, st.LogKeep), "")
			}
			if err != nil {
				logf("Some old server logs could not be deleted: "+err.Error(), "warning")
			}
		},
		InstalledMods: installedModIDs,
		OnPlayer:      func(ev logparse.Event) { s.chat.Observe(ev) },
		SaveRun:       runs.Save,
		LoadRun:       runs.Load,
		Log:           s.log,
	})
	ready := func() bool { return s.runner.Status().Phase == runner.Ready }
	s.maps = mapchange.New(mapchange.Deps{
		Config:       conf,
		Ready:        ready,
		ModScenarios: func() []string { return loadedModScenarios(s.runner.Status().Mods) },
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
		Log: logf,
	})
	s.notify = notify.New(notify.Deps{Config: conf, Post: discord.Post, Log: logf}, true)
	s.players = players.New(players.Deps{
		Config: conf,
		Ready:  ready,
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
		Log: logf,
	})
	s.watch = watchdog.New(watchdog.Deps{
		Config:      conf,
		Phase:       func() string { return string(s.runner.Status().Phase) },
		Start:       s.runner.Start,
		Restart:     s.runner.Restart,
		StopAndWait: s.runner.StopAndWait,
		Update:      func() error { return a.updateServer(id) },
		Say:         func(text string) { s.chat.Broadcast(text) },
		Log:         logf,
	})
	s.chat = chat.New(chat.Deps{
		OnPlayer: func(kind, name string) { s.notify.Event(kind, "", name) },
		Config:   conf,
		Dial: func(port int, password string) (chat.Conn, error) {
			c, err := rcon.Dial(port, password)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		Log: logf,
	}, true)
	a.wireServer(s)
	return s
}

// wireServer liga os avisos do servidor à tela: o estado de cada um vai no "fleet:status"; o do escolhido também no
// "server:status", e só as linhas de log dele vão para a gaveta.
func (a *app) wireServer(s *srv) {
	if a.ctx == nil {
		return
	}
	ctx := a.ctx
	s.log.SetNotify(func(lines []logbuf.Line) {
		if a.cfg.SelectedID() == s.id {
			runtime.EventsEmit(ctx, "log:lines", lines)
		}
	})
	s.runner.SetNotify(func(st runner.Status) {
		// Servidor parou: a conexão das mensagens de chat fecha e a lista de jogadores recomeça.
		if st.Phase == runner.Stopped || st.Phase == runner.Crashed {
			s.chat.Reset()
		}
		runtime.EventsEmit(ctx, "fleet:status", a.summary(s, st))
		if a.cfg.SelectedID() == s.id {
			runtime.EventsEmit(ctx, "server:status", st)
		}
	})
}

// checkStart confere portas e GSLT contra os outros servidores rodando e as portas ocupadas no Windows (ADR-048).
func (a *app) checkStart(id string, p server.Profile, sec server.Secrets) error {
	var peers []server.Peer
	for _, peer := range a.cfg.Peers(id) {
		if o := a.server(peer.ID); o != nil && len(o.runner.Status().PIDs) > 0 {
			peers = append(peers, peer)
		}
	}
	if err := server.CheckStart(p, sec, peers); err != nil {
		return err
	}
	if busy := server.BusyPorts(p, portBusy); len(busy) > 0 {
		e := problems.New(problems.PortTaken, "")
		e.Field = busy[0].Field
		e.Params = map[string]string{"port": strconv.Itoa(busy[0].Port), "portField": busy[0].Field}
		return e
	}
	return nil
}

// portBusy diz se outro programa já usa a porta nesta máquina: abre e fecha a porta em 0.0.0.0 e em 127.0.0.1. Só
// 0.0.0.0 não basta: o Windows deixa abrir ali uma porta que outro programa usa só em 127.0.0.1, e as conexões
// locais vão para o outro programa (o MongoDB na 27017 tomava o RCON do terceiro servidor).
func portBusy(port int, tcp bool) bool {
	for _, host := range []string{"0.0.0.0", "127.0.0.1"} {
		addr := net.JoinHostPort(host, strconv.Itoa(port))
		if tcp {
			l, err := net.Listen("tcp4", addr) // #nosec G102 -- só confere se a porta está livre (fecha na hora), como o servidor a abriria
			if err != nil {
				return true
			}
			_ = l.Close()
			continue
		}
		c, err := net.ListenPacket("udp4", addr) // #nosec G102 -- idem, em UDP
		if err != nil {
			return true
		}
		_ = c.Close()
	}
	return false
}

// modsBusy diz se outro servidor (que não id) está iniciando ou sincronizando os mods agora.
func (a *app) modsBusy(id string) bool {
	for _, s := range a.all() {
		if s.id == id {
			continue
		}
		p, _, err := a.cfg.CurrentOf(s.id)
		if err == nil && s.runner.ModsBusy(p.UseMods, modsBusyLimit) {
			return true
		}
	}
	return false
}

// ServerSummary é uma linha da lista de servidores do Painel.
type ServerSummary struct {
	ID       string          `json:"id"`
	Name     string          `json:"name"`
	Hostname string          `json:"hostname"`
	Instance int             `json:"instance"`
	Port     int             `json:"port"`
	Selected bool            `json:"selected"`
	Phase    runner.Phase    `json:"phase"`
	Map      string          `json:"map"`
	Since    time.Time       `json:"since"`
	External bool            `json:"external"`
	Problem  *problems.Error `json:"problem,omitempty"`
	// Conflicts são as portas deste servidor que outro servidor também usa (eles não podem rodar juntos).
	Conflicts []server.PortConflict `json:"conflicts"`
}

func (a *app) summary(s *srv, st runner.Status) ServerSummary {
	p, _, _ := a.cfg.CurrentOf(s.id)
	return ServerSummary{
		ID: s.id, Name: a.cfg.Name(s.id), Hostname: p.Hostname, Instance: server.InstanceNumber(p.Instance),
		Port: p.Port, Selected: a.cfg.SelectedID() == s.id, Phase: st.Phase, Map: st.Map, Since: st.Since,
		External: st.External, Problem: st.Problem, Conflicts: server.PortConflicts(p, a.cfg.Peers(s.id)),
	}
}

// summaries devolve a lista de servidores para o Painel.
func (a *app) summaries() []ServerSummary {
	list := a.all()
	out := make([]ServerSummary, 0, len(list))
	for _, s := range list {
		out = append(out, a.summary(s, s.runner.Status()))
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Instance < out[j].Instance })
	return out
}

// logAll escreve uma linha no log de todos os servidores (a instalação é compartilhada).
func (a *app) logAll(source, text, level string) {
	for _, s := range a.all() {
		s.log.Add(source, text, level)
	}
}

// fanLog é o log da instalação: vai para todos os servidores.
type fanLog struct{ a *app }

func (f fanLog) Add(source, text, level string) { f.a.logAll(source, text, level) }

// installedModIDs devolve os ids dos mods instalados (state.json do mod.io).
func installedModIDs() []int {
	list, err := modio.ServerInstalled(modio.RootDir())
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
}

// othersRunning devolve o problema "outros servidores rodando" com os nomes, ou nil.
func (a *app) othersRunning(except string) error {
	list := a.running(except)
	if len(list) == 0 {
		return nil
	}
	names := ""
	for i, s := range list {
		if i > 0 {
			names += ", "
		}
		names += a.cfg.Name(s.id)
	}
	e := problems.New(problems.OthersRunning, "")
	e.Params = map[string]string{"servers": names}
	return e
}

// stopAll para os servidores rodando e devolve quais eram (para iniciar de novo depois).
func (a *app) stopAll() ([]string, error) {
	var ids []string
	for _, s := range a.running("") {
		ids = append(ids, s.id)
		if err := s.runner.StopAndWait(); err != nil {
			return ids, err
		}
	}
	return ids, nil
}

// startIDs inicia os servidores da lista (um de cada vez quando usam mods: o runner espera a vez).
func (a *app) startIDs(ids []string) error {
	var errs []error
	for _, id := range ids {
		if s := a.server(id); s != nil {
			if err := s.runner.Start(); err != nil {
				errs = append(errs, err)
			}
		}
	}
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}
