package desktop

import (
	"errors"
	"os"
	"path/filepath"
	"time"

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

// tickAutomation roda a cada segundo: avisos do Discord quando a fase muda, mensagens programadas e reinício
// automático.
func (a *app) tickAutomation() {
	st := a.runner.Status()
	prev := a.lastPhase
	a.lastPhase = st.Phase
	if prev != "" && prev != st.Phase {
		switch {
		case st.Phase == runner.Ready && prev != runner.Ready:
			a.notify.Event(notify.Ready, st.Map, "")
		case st.Phase == runner.Crashed:
			a.notify.Event(notify.Crashed, st.Map, "")
		case st.Phase == runner.Stopped && prev != runner.NotInstalled:
			a.notify.Event(notify.Stopped, st.Map, "")
		}
	}
	a.chat.Tick(st.Phase == runner.Ready)
	a.watch.Tick()
}

// updateServer atualiza o servidor pelo SteamCMD (reinício diário) e espera terminar.
func (a *app) updateServer() error {
	p := a.cfg.Profile()
	if p.SteamCMDDir == "" || p.ServerDir == "" {
		return problems.New(problems.ServerUpdateFailed, "SteamCMD folder not set")
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
	a.log.Add("app", "Server updated by SteamCMD", "info")
	return nil
}

// autoStartServer inicia o servidor quando o Windows abriu o app e o usuário pediu (ADR-022).
func (a *app) autoStartServer() {
	if !a.autostart || !a.cfg.Settings().AutoStartServer {
		return
	}
	go func() {
		time.Sleep(5 * time.Second)
		if a.runner.Status().Phase != runner.Stopped || a.install.Busy() {
			return
		}
		if err := a.runner.Start(); err != nil {
			a.log.Add("app", "Automatic start with Windows failed: "+err.Error(), "warning")
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

// Players pede a lista de jogadores e de banidos ao servidor ligado.
func (s *ServerAPI) Players() (players.View, error) { return s.a.players.Refresh() }

// Kick expulsa um jogador da última lista (a tela manda o net ID da lista e o motivo).
func (s *ServerAPI) Kick(netID, reason string) error { return s.a.players.Kick(netID, reason) }

// Ban bane um jogador da última lista (minutes 0 = permanente).
func (s *ServerAPI) Ban(netID string, minutes int, reason string) error {
	return s.a.players.Ban(netID, minutes, reason)
}

// Unban tira um banimento da última lista de banidos.
func (s *ServerAPI) Unban(netID string) error { return s.a.players.Unban(netID) }

// Say manda uma mensagem a todos no chat.
func (s *ServerAPI) Say(text string) error { return s.a.players.Say(text) }

// AutomationInfo é o que a aba Automação mostra.
type AutomationInfo struct {
	NextRestart time.Time `json:"nextRestart"`
	HasWebhook  bool      `json:"hasWebhook"`
	HasSteamCMD bool      `json:"hasSteamCmd"`
}

// Automation devolve o próximo reinício diário e o que falta para os avisos e a atualização.
func (s *ServerAPI) Automation() AutomationInfo {
	sec, _ := s.a.cfg.SecretsState()
	return AutomationInfo{NextRestart: s.a.watch.Next(), HasWebhook: sec["discordWebhook"],
		HasSteamCMD: s.a.cfg.Profile().SteamCMDDir != ""}
}

// TestDiscord manda um aviso de teste para o webhook.
func (s *ServerAPI) TestDiscord() error { return s.a.notify.Test(s.a.runner.Status().Map) }

// ---- catálogos e perfis ----

// Mutators devolve os mutators do jogo.
func (c *ConfigAPI) Mutators() []server.Mutator { return server.Mutators }

// ModMutators devolve os mutators achados nos mods instalados (lidos dos .pak, ADR-023).
func (c *ConfigAPI) ModMutators() []modio.ModMutator { return modio.Mutators(modio.RootDir()) }

// RulesInfo é o catálogo da aba Regras.
type RulesInfo struct {
	Groups []string         `json:"groups"`
	Rules  []server.RuleDef `json:"rules"`
}

// Rules devolve as opções do Game.ini e do Engine.ini que o app sabe gravar.
func (c *ConfigAPI) Rules() RulesInfo {
	return RulesInfo{Groups: server.RuleGroups, Rules: server.Rules}
}

// CycleLine é uma linha da rotação com o cenário e a iluminação.
type CycleLine struct {
	Scenario string `json:"scenario"`
	Lighting string `json:"lighting"`
}

// MapCycleLines lê o cenário e a iluminação de cada linha da rotação (para a tela mostrar).
func (c *ConfigAPI) MapCycleLines(lines []string) []CycleLine {
	if len(lines) > 500 {
		lines = lines[:500]
	}
	out := make([]CycleLine, len(lines))
	for i, l := range lines {
		out[i] = CycleLine{Scenario: server.LineScenario(l), Lighting: server.LineLighting(l)}
	}
	return out
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

// LocalAddresses lista os IPv4 desta máquina (escolha da placa de rede do servidor, -MultiHome).
func (c *ConfigAPI) LocalAddresses() []winsys.Address { return winsys.LocalAddresses() }

// Profiles lista os perfis guardados.
func (c *ConfigAPI) Profiles() []config.ProfileRef { return c.a.cfg.Profiles() }

// CreateProfile guarda um perfil novo (cópia do atual ou padrão).
func (c *ConfigAPI) CreateProfile(name string, copyCurrent bool) ([]config.ProfileRef, error) {
	return c.a.cfg.CreateProfile(name, copyCurrent)
}

// RenameProfile troca o nome de um perfil.
func (c *ConfigAPI) RenameProfile(id, name string) ([]config.ProfileRef, error) {
	return c.a.cfg.RenameProfile(id, name)
}

// DeleteProfile apaga um perfil guardado.
func (c *ConfigAPI) DeleteProfile(id string) ([]config.ProfileRef, error) {
	return c.a.cfg.DeleteProfile(id)
}

// SwitchProfile passa a usar outro perfil (só com o servidor parado).
func (c *ConfigAPI) SwitchProfile(id string) (server.Profile, error) {
	if len(c.a.runner.Status().PIDs) > 0 {
		return c.a.cfg.Profile(), problems.New(problems.ServerRunning, "")
	}
	if c.a.install.Busy() {
		return c.a.cfg.Profile(), problems.New(problems.Busy, "")
	}
	return c.a.cfg.SwitchProfile(id)
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
	s.a.log.Add("app", "Windows Firewall rules created for the server", "info")
	return s.Firewall(), nil
}
