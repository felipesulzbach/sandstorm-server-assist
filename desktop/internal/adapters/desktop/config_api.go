package desktop

import (
	"os"

	"sandstorm-server-assist/desktop/internal/adapters/modio"
	"sandstorm-server-assist/desktop/internal/app/config"
	"sandstorm-server-assist/desktop/internal/app/installer"
	"sandstorm-server-assist/desktop/internal/domain/problems"
	"sandstorm-server-assist/desktop/internal/domain/server"
	"sandstorm-server-assist/desktop/internal/domain/steamid"
)

// ConfigAPI: perfil do servidor, segredos e pastas.
type ConfigAPI struct{ a *app }

// Profile devolve o perfil.
func (c *ConfigAPI) Profile() server.Profile { return c.a.cfg.Profile() }

// SaveProfile valida e grava o perfil (as pastas só mudam pelos diálogos).
func (c *ConfigAPI) SaveProfile(p server.Profile) (server.Profile, error) {
	return c.a.cfg.SaveProfile(p)
}

// DefaultProfile devolve os valores recomendados (botão "Restaurar padrões").
func (c *ConfigAPI) DefaultProfile() server.Profile { return server.Default() }

// MapsInfo são os mapas do jogo base e os modos de cenário.
type MapsInfo struct {
	Maps  []server.Map  `json:"maps"`
	Modes []server.Mode `json:"modes"`
}

// Maps devolve os mapas conhecidos.
func (c *ConfigAPI) Maps() MapsInfo {
	return MapsInfo{Maps: server.Maps, Modes: server.Modes}
}

// SecretsState diz quais segredos estão definidos (sem os valores).
func (c *ConfigAPI) SecretsState() (map[string]bool, error) { return c.a.cfg.SecretsState() }

// SetSecret grava um segredo no cofre (vazio apaga).
func (c *ConfigAPI) SetSecret(key, value string) (map[string]bool, error) {
	if err := c.a.cfg.SetSecret(key, value); err != nil {
		return nil, err
	}
	return c.a.cfg.SecretsState()
}

// RevealSecret devolve o valor de um segredo quando o usuário clica em "Mostrar".
func (c *ConfigAPI) RevealSecret(key string) (string, error) { return c.a.cfg.Secret(key) }

// ExtractSteamID acha o SteamID64 num texto colado (número ou link do perfil).
func (c *ConfigAPI) ExtractSteamID(text string) (string, error) {
	if len(text) > 512 {
		return "", problems.Field("admins", "validation.steamId", "value", "")
	}
	id, ok := steamid.Extract(text)
	if !ok {
		return "", problems.Field("admins", "validation.steamIdNotFound")
	}
	return id, nil
}

// FreshInfo é o que a tela precisa para oferecer começar do zero (ADR-047). Campos copiados de config.Fresh (o
// gerador do Wails não cria os tipos de uma struct embutida) mais os dados do mod.io no Windows.
type FreshInfo struct {
	Offer      bool `json:"offer"`
	Missing    bool `json:"missing"`
	HasSecrets bool `json:"hasSecrets"`
	Others     int  `json:"others"`
	// Mods é quantos mods baixados há na pasta do mod.io; ModioSession diz se o servidor tem autorização do mod.io.
	Mods         int  `json:"mods"`
	ModioSession bool `json:"modioSession"`
}

// FreshState diz se a tela deve oferecer começar do zero: a pasta do servidor não tem servidor (ou sumiu) e o
// perfil, o cofre ou o mod.io do servidor ainda têm dados (ADR-047).
func (c *ConfigAPI) FreshState() FreshInfo {
	f := c.a.cfg.FreshState()
	info := FreshInfo{Missing: f.Missing, HasSecrets: f.HasSecrets, Others: f.Others}
	if list, err := modio.ReadInstalled(modio.RootDir()); err == nil {
		info.Mods = len(list)
	}
	if uf, err := modio.UserFile(); err == nil {
		_, statErr := os.Stat(uf)
		info.ModioSession = statErr == nil
	}
	info.Offer = f.Offer || (f.NoServer && (info.Mods > 0 || info.ModioSession))
	return info
}

// StartFresh volta o perfil em uso ao padrão, mantendo as pastas; clearSecrets também esvazia o cofre (ADR-047).
func (c *ConfigAPI) StartFresh(clearSecrets bool) (server.Profile, error) {
	if len(c.a.runner.Status().PIDs) > 0 {
		return c.a.cfg.Profile(), problems.New(problems.ServerRunning, "")
	}
	if c.a.install.Busy() {
		return c.a.cfg.Profile(), problems.New(problems.Busy, "")
	}
	return c.a.cfg.StartFresh(clearSecrets)
}

// ChooseServerDir abre o diálogo de pasta para a pasta do servidor.
func (c *ConfigAPI) ChooseServerDir() (server.Profile, error) {
	cur := c.a.cfg.Profile()
	dir, err := c.a.chooseDir(cur.ServerDir)
	if err != nil || dir == "" {
		return cur, err
	}
	p, err := c.a.cfg.SetServerDir(dir)
	if err == nil && p.SteamCMDDir == "" {
		// Instalação existente: o appmanifest diz qual SteamCMD foi usado.
		if sc := installer.LauncherDir(p.ServerDir); sc != "" {
			p, _ = c.a.cfg.SetSteamCMDDir(sc)
		}
	}
	return p, err
}

// ChooseSteamCMDDir abre o diálogo de pasta para o SteamCMD.
func (c *ConfigAPI) ChooseSteamCMDDir() (server.Profile, error) {
	cur := c.a.cfg.Profile()
	dir, err := c.a.chooseDir(cur.SteamCMDDir)
	if err != nil || dir == "" {
		return cur, err
	}
	return c.a.cfg.SetSteamCMDDir(dir)
}

// ExistingScan é a comparação com o servidor existente e se ele está rodando agora. Campos copiados de
// config.Scan (o gerador do Wails não cria os tipos de uma struct embutida).
type ExistingScan struct {
	Found     bool          `json:"found"`
	LogFile   string        `json:"logFile"`
	Diffs     []server.Diff `json:"diffs"`
	RconInIni bool          `json:"rconInIni"`
	Running   bool          `json:"running"`
}

// ScanExisting compara o perfil do app com o servidor já configurado na pasta escolhida (rodando ou parado).
func (c *ConfigAPI) ScanExisting() (ExistingScan, error) {
	s, err := c.a.cfg.ScanExisting(findAPIKey(c.a.cfg.Profile().ServerDir))
	return ExistingScan{
		Found: s.Found, LogFile: s.LogFile, Diffs: s.Diffs, RconInIni: s.RconInIni,
		Running: len(c.a.runner.Status().PIDs) > 0,
	}, err
}

// ApplyExisting traz os dados do servidor existente. keepMine são os campos em conflito em que o usuário quer
// manter o que digitou no app.
func (c *ConfigAPI) ApplyExisting(keepMine []string) (config.ApplyResult, error) {
	if len(keepMine) > 64 {
		return config.ApplyResult{}, problems.Field("keepMine", "validation.unknownField")
	}
	dir := c.a.cfg.Profile().ServerDir
	return c.a.cfg.ApplyExisting(installer.LauncherDir(dir), findAPIKey(dir), keepMine)
}

// RconInIni diz se o Game.ini ainda tem a senha de RCON em texto puro.
func (c *ConfigAPI) RconInIni() bool { return c.a.cfg.RconInIni() }

// MoveRconToVault tira a senha de RCON do Game.ini e guarda no cofre (servidor parado).
func (c *ConfigAPI) MoveRconToVault() error {
	if len(c.a.runner.Status().PIDs) > 0 {
		return problems.New(problems.ServerRunning, "")
	}
	return c.a.cfg.MoveRconToVault()
}
