package desktop

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"sandstorm-server-assist/desktop/internal/adapters/fsx"
	"sandstorm-server-assist/desktop/internal/adapters/modio"
	"sandstorm-server-assist/desktop/internal/adapters/winsys"
	"sandstorm-server-assist/desktop/internal/app/config"
	"sandstorm-server-assist/desktop/internal/app/notify"
	"sandstorm-server-assist/desktop/internal/app/players"
	"sandstorm-server-assist/desktop/internal/app/runner"
	"sandstorm-server-assist/desktop/internal/domain/problems"
	"sandstorm-server-assist/desktop/internal/domain/server"
)

// ---- automação (ADR-022) ----

// tickAutomation roda a cada segundo para cada servidor: avisos do Discord quando a fase muda, mensagens
// programadas e reinício automático.
func (a *app) tickAutomation(s *srv) {
	st := s.runner.Status()
	prev := s.lastPhase
	s.lastPhase = st.Phase
	if prev != "" && prev != st.Phase {
		switch {
		case st.Phase == runner.Ready && prev != runner.Ready:
			s.notify.Event(notify.Ready, st.Map, "")
		case st.Phase == runner.Crashed:
			s.notify.Event(notify.Crashed, st.Map, "")
		case st.Phase == runner.Stopped && prev != runner.NotInstalled && prev != runner.Waiting:
			s.notify.Event(notify.Stopped, st.Map, "")
		}
	}
	s.chat.Tick(st.Phase == runner.Ready)
	s.watch.Tick()
}

// updateServer atualiza a instalação do servidor id pelo SteamCMD (reinício diário) e espera terminar. Com outro
// servidor da mesma instalação rodando, não atualiza (os arquivos estão em uso, ADR-048).
func (a *app) updateServer(id string) error {
	p, _, err := a.cfg.CurrentOf(id)
	if err != nil {
		return err
	}
	if p.SteamCMDDir == "" || p.ServerDir == "" {
		return problems.New(problems.ServerUpdateFailed, "SteamCMD folder not set")
	}
	if err := a.installRunning(p.ServerDir, id); err != nil {
		return err
	}
	if err := a.install.Install(p.SteamCMDDir, p.ServerDir, false); err != nil {
		return err
	}
	deadline := time.Now().Add(3 * time.Hour)
	for a.install.Busy() {
		if time.Now().After(deadline) {
			a.install.Cancel()
			return problems.New(problems.ServerUpdateFailed, "timeout")
		}
		time.Sleep(2 * time.Second)
	}
	if st := a.install.State(); st.Error != nil {
		return st.Error
	}
	a.logAll("app", "Server updated by SteamCMD", "info")
	return nil
}

// installRunning diz (com o problema "outros servidores rodando") se algum servidor da instalação dir, menos
// except, está rodando: atualizar a instalação troca arquivos em uso.
func (a *app) installRunning(dir, except string) error {
	var names []string
	for _, s := range a.running(except) {
		if p, _, err := a.cfg.CurrentOf(s.id); err == nil && sameFolder(p.ServerDir, dir) {
			names = append(names, a.cfg.Name(s.id))
		}
	}
	if len(names) == 0 {
		return nil
	}
	e := problems.New(problems.OthersRunning, "")
	e.Params = map[string]string{"servers": strings.Join(names, ", ")}
	return e
}

// sameFolder compara duas pastas do Windows.
func sameFolder(a, b string) bool {
	return a != "" && b != "" && strings.EqualFold(filepath.Clean(a), filepath.Clean(b))
}

// autoStartServer inicia os servidores quando o Windows abriu o app e o usuário pediu (ADR-022). Com vários, inicia
// todos em ordem; os que usam mods esperam a vez (ADR-048), e um com as portas de outro já rodando fica parado.
func (a *app) autoStartServer() {
	if !a.autostart || !a.cfg.Settings().AutoStartServer {
		return
	}
	go func() {
		time.Sleep(5 * time.Second)
		if a.install.Busy() {
			return
		}
		for _, s := range a.all() {
			if s.runner.Status().Phase != runner.Stopped {
				continue
			}
			if err := s.runner.Start(); err != nil {
				s.log.Add("app", "Automatic start with Windows failed: "+err.Error(), "warning")
			}
		}
	}()
}

// syncAutostart refaz a entrada do registro se o app mudou de lugar (atualização, pasta nova).
func (a *app) syncAutostart() {
	if a.cfg.Settings().StartWithWindows {
		_ = winsys.SetAutostart(true)
	}
}

// FirewallHelper roda no app aberto como administrador (winsys.FirewallFlag): lê a pasta do servidor da
// configuração do usuário e cria as regras. Devolve o código de saída.
func FirewallHelper() int {
	path, err := config.DefaultPath()
	if err != nil {
		return 1
	}
	dir := config.Open(path, memoryVault{}).Profile().ServerDir
	if dir == "" || !fsx.SafeDir(dir) || !fsx.Exists(filepath.Join(dir, config.ServerExe)) {
		return 2
	}
	if err := winsys.AddFirewallRules(dir); err != nil {
		_, _ = os.Stderr.WriteString("firewall: " + err.Error() + "\n")
		return 3
	}
	return 0
}

// ---- jogadores (RCON) ----

// Players pede a lista de jogadores e de banidos ao servidor escolhido.
func (s *ServerAPI) Players() (players.View, error) { return s.a.sel().players.Refresh() }

// Kick expulsa um jogador da última lista (a tela manda o net ID da lista e o motivo).
func (s *ServerAPI) Kick(netID, reason string) error { return s.a.sel().players.Kick(netID, reason) }

// Ban bane um jogador da última lista (minutes 0 = permanente).
func (s *ServerAPI) Ban(netID string, minutes int, reason string) error {
	return s.a.sel().players.Ban(netID, minutes, reason)
}

// Unban tira um banimento da última lista de banidos.
func (s *ServerAPI) Unban(netID string) error { return s.a.sel().players.Unban(netID) }

// Say manda uma mensagem a todos no chat.
func (s *ServerAPI) Say(text string) error { return s.a.sel().players.Say(text) }

// AutomationInfo é o que a aba Automação mostra.
type AutomationInfo struct {
	NextRestart time.Time `json:"nextRestart"`
	HasWebhook  bool      `json:"hasWebhook"`
	HasSteamCMD bool      `json:"hasSteamCmd"`
}

// Automation devolve o próximo reinício diário e o que falta para os avisos e a atualização.
func (s *ServerAPI) Automation() AutomationInfo {
	sec, _ := s.a.cfg.SecretsState()
	return AutomationInfo{NextRestart: s.a.sel().watch.Next(), HasWebhook: sec["discordWebhook"],
		HasSteamCMD: s.a.cfg.Profile().SteamCMDDir != ""}
}

// TestDiscord manda um aviso de teste para o webhook.
func (s *ServerAPI) TestDiscord() error {
	sv := s.a.sel()
	return sv.notify.Test(sv.runner.Status().Map)
}

// ---- catálogos e perfis ----

// Mutators devolve os mutators do jogo.
func (c *ConfigAPI) Mutators() []server.Mutator { return server.Mutators }

// ModMutators devolve os mutators achados nos mods instalados (lidos dos .pak, ADR-023).
func (c *ConfigAPI) ModMutators() []modio.ModMutator { return modio.Mutators(modio.RootDir()) }

// ModScenarios devolve os mapas e cenários dos mods instalados (lidos dos .pak, ADR-024), para escolher o mapa de
// mod depois de carregar.
func (c *ConfigAPI) ModScenarios() []modio.ModScenario { return modio.Scenarios(modio.RootDir()) }

// RulesInfo é o catálogo da aba Regras.
type RulesInfo struct {
	Groups []string         `json:"groups"`
	Rules  []server.RuleDef `json:"rules"`
}

// Rules devolve as opções do Game.ini e do Engine.ini que o app sabe gravar.
func (c *ConfigAPI) Rules() RulesInfo {
	return RulesInfo{Groups: server.RuleGroups, Rules: server.Rules}
}

// ServerRules devolve as regras do jogo que o Game.ini e o Engine.ini do servidor já têm (vazio se nenhuma), para a
// aba Regras oferecer usá-las (ADR-040).
func (c *ConfigAPI) ServerRules() map[string]string {
	if r := c.a.cfg.ServerRules(); r != nil {
		return r
	}
	return map[string]string{}
}

// SetMapCycleLighting troca a iluminação de uma linha da rotação (index) ou de todas (index -1).
func (c *ConfigAPI) SetMapCycleLighting(lines []string, index int, lighting string) ([]string, error) {
	if !server.ValidLighting(lighting) || len(lines) > 500 || index < -1 || index >= len(lines) {
		return lines, problems.Field("mapCycle", "validation.lighting")
	}
	out := append([]string(nil), lines...)
	for i := range out {
		if index == -1 || i == index {
			out[i] = server.SetLineLighting(out[i], lighting)
		}
	}
	return out, nil
}

// CycleView é a rotação em blocos por mapa e os mapas que podem entrar nela (ADR-046).
type CycleView struct {
	Groups []server.CycleGroup `json:"groups"`
	Maps   []server.CycleMap   `json:"maps"`
}

// MapCycleView divide a rotação em blocos por mapa para a tela (o MapCycle.txt continua uma lista simples).
func (c *ConfigAPI) MapCycleView(lines []string) CycleView {
	if len(lines) > server.MaxCycleLines {
		lines = lines[:server.MaxCycleLines]
	}
	p, _, _ := c.a.cfg.Current()
	cat := c.a.scenarios.catalog(p.ServerDir)
	missing := server.MissingInCycle(lines, c.a.scenarios.known(p.ServerDir))
	return CycleView{Groups: server.GroupMapCycle(lines, missing, cat), Maps: cat}
}

// AddToMapCycle põe cenários na rotação depois da linha after (-1: no fim), de dia, de noite ou nos dois, sem
// repetir.
func (c *ConfigAPI) AddToMapCycle(lines []string, after int, ids, lightings []string) ([]string, error) {
	if len(lines) > server.MaxCycleLines || len(ids) > server.MaxCycleLines || len(lightings) > 2 {
		return lines, problems.Field("mapCycle", "validation.mapCycleTooLong")
	}
	return server.AddToMapCycle(lines, after, ids, lightings)
}

// RemoveFromMapCycle tira linhas da rotação.
func (c *ConfigAPI) RemoveFromMapCycle(lines []string, indexes []int) []string {
	if len(lines) > server.MaxCycleLines {
		return lines
	}
	return server.RemoveFromMapCycle(lines, indexes)
}

// MoveMapCycleGroup sobe ou desce na rotação o bloco do mapa da linha index.
func (c *ConfigAPI) MoveMapCycleGroup(lines []string, index int, up bool) []string {
	if len(lines) > server.MaxCycleLines {
		return lines
	}
	return server.MoveMapCycleGroup(lines, index, up)
}

// LocalAddresses lista os IPv4 desta máquina (escolha da placa de rede do servidor, -MultiHome).
func (c *ConfigAPI) LocalAddresses() []winsys.Address { return winsys.LocalAddresses() }

// Profiles lista os servidores (ADR-048; antes perfis).
func (c *ConfigAPI) Profiles() []config.ProfileRef { return c.a.cfg.Profiles() }

// CreateProfile guarda um servidor novo (cópia do escolhido ou padrão), com portas livres, e passa a editá-lo.
func (c *ConfigAPI) CreateProfile(name string, copyCurrent bool) ([]config.ProfileRef, error) {
	refs, id, err := c.a.cfg.CreateProfile(name, copyCurrent, portBusy)
	if err != nil {
		return refs, err
	}
	c.a.server(id)
	if _, err := c.a.cfg.SwitchProfile(id); err != nil {
		return refs, err
	}
	return c.a.cfg.Profiles(), nil
}

// RenameProfile troca o nome de um perfil.
func (c *ConfigAPI) RenameProfile(id, name string) ([]config.ProfileRef, error) {
	return c.a.cfg.RenameProfile(id, name)
}

// DeleteProfile apaga um servidor que não é o escolhido, parado (os arquivos dele na pasta do servidor ficam).
func (c *ConfigAPI) DeleteProfile(id string) ([]config.ProfileRef, error) {
	if s := c.a.server(id); s != nil {
		if st := s.runner.Status(); len(st.PIDs) > 0 || st.Phase == runner.Waiting {
			return c.a.cfg.Profiles(), problems.New(problems.ServerRunning, "")
		}
	}
	refs, err := c.a.cfg.DeleteProfile(id)
	c.a.all()
	return refs, err
}

// SwitchProfile escolhe o servidor que as telas editam (também com servidores rodando, ADR-048).
func (c *ConfigAPI) SwitchProfile(id string) (server.Profile, error) {
	p, err := c.a.cfg.SwitchProfile(id)
	if err == nil {
		if s := c.a.server(id); s != nil && c.a.ctx != nil {
			runtime.EventsEmit(c.a.ctx, "server:status", s.runner.Status())
		}
	}
	return p, err
}

// PortConflicts lista as portas do servidor escolhido que outro servidor também usa e, com ele parado, as que outro
// programa do PC ocupa (Server vazio).
func (c *ConfigAPI) PortConflicts() []server.PortConflict {
	p := c.a.cfg.Profile()
	out := server.PortConflicts(p, c.a.cfg.Peers(c.a.cfg.SelectedID()))
	if st := c.a.sel().runner.Status(); len(st.PIDs) == 0 && st.Phase != runner.Waiting {
		taken := map[string]bool{}
		for _, x := range out {
			taken[x.Field] = true
		}
		for _, b := range server.BusyPorts(p, portBusy) {
			if !taken[b.Field] {
				out = append(out, b)
			}
		}
	}
	return out
}

// FreePorts sugere portas de jogo, consulta e RCON que nenhum outro servidor usa (a tela põe nos campos; o usuário
// salva).
func (c *ConfigAPI) FreePorts() []int {
	var others []server.Profile
	for _, p := range c.a.cfg.Peers(c.a.cfg.SelectedID()) {
		others = append(others, p.Profile)
	}
	port, query, rcon := server.FreePorts(others, portBusy)
	return []int{port, query, rcon}
}

// ---- Windows ----

// SetAutostart liga ou desliga abrir com o Windows e iniciar o servidor junto.
func (s *SystemAPI) SetAutostart(withWindows, startServer bool) (config.Settings, error) {
	if err := winsys.SetAutostart(withWindows); err != nil {
		return s.a.cfg.Settings(), problems.Wrap(problems.AutostartFailed, err)
	}
	return s.a.cfg.SetAutostart(withWindows, startServer)
}

// FirewallInfo diz se as regras do Firewall existem.
type FirewallInfo struct {
	// Ready: há servidor instalado para criar a regra.
	Ready bool `json:"ready"`
	Rules bool `json:"rules"`
}

// Firewall confere as regras do Firewall do Windows para o servidor (sem administrador).
func (s *SystemAPI) Firewall() FirewallInfo {
	dir := s.a.cfg.Profile().ServerDir
	info := FirewallInfo{Ready: dir != "" && fsx.Exists(filepath.Join(dir, config.ServerExe))}
	if !info.Ready {
		return info
	}
	info.Rules = true
	for _, prog := range winsys.ServerPrograms(dir) {
		if fsx.Exists(prog) && !winsys.FirewallRule(prog) {
			info.Rules = false
		}
	}
	return info
}

// CreateFirewallRule cria as regras do Firewall: o Windows pede a confirmação de administrador.
func (s *SystemAPI) CreateFirewallRule() (FirewallInfo, error) {
	if !s.Firewall().Ready {
		return FirewallInfo{}, problems.New(problems.ServerNotInstalled, "")
	}
	if err := winsys.RunElevated(winsys.FirewallFlag); err != nil {
		if errors.Is(err, winsys.ErrCanceled) {
			return s.Firewall(), problems.New(problems.Canceled, "")
		}
		return s.Firewall(), problems.Wrap(problems.FirewallFailed, err)
	}
	s.a.sel().log.Add("app", "Windows Firewall rules created for the server", "info")
	return s.Firewall(), nil
}

// CycleDefault é a rotação padrão com os mapas do jogo (ADR-043).
type CycleDefault struct {
	Lines []string `json:"lines"`
	// Maps: quantos mapas a rotação padrão tem (cada um com os cenários de dia e de noite).
	Maps int `json:"maps"`
	// ServerHasFile: o servidor já tem um MapCycle.txt (com a rotação vazia no app, o app o mantém).
	ServerHasFile bool `json:"serverHasFile"`
}

// DefaultMapCycle monta a rotação padrão no modo do mapa inicial, com os cenários lidos dos .pak do servidor (ou a
// lista de mapas do app, sem o servidor instalado).
func (c *ConfigAPI) DefaultMapCycle() CycleDefault {
	p := c.a.cfg.Profile()
	var base []string
	if p.ServerDir != "" {
		base = modio.BaseScenarios(filepath.Join(p.ServerDir, "Insurgency", "Content", "Paks"))
	}
	ids := server.DefaultMapCycle(p.StartScenario, base)
	out := CycleDefault{Lines: server.WithDayAndNight(ids), Maps: server.CountMaps(ids)}
	if out.Lines == nil {
		out.Lines = []string{}
	}
	out.ServerHasFile = p.ServerDir != "" && fsx.Exists(config.FilePath(p, config.Files(p).MapCycle))

	return out
}
