package desktop

import (
	"runtime/debug"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"sandstorm-server-assist/desktop/internal/adapters/modio"
	upd "sandstorm-server-assist/desktop/internal/adapters/updater"
	"sandstorm-server-assist/desktop/internal/app/config"
	"sandstorm-server-assist/desktop/internal/app/installer"
	"sandstorm-server-assist/desktop/internal/app/updater"
	"sandstorm-server-assist/desktop/internal/domain/legal"
	"sandstorm-server-assist/desktop/internal/domain/problems"
	"sandstorm-server-assist/desktop/internal/domain/server"
)

// SystemAPI: versão, links, preferências e a lista de primeiros passos.
type SystemAPI struct{ a *app }

// links são os únicos endereços que a tela pode pedir para abrir (pela chave, nunca pela URL).
var links = map[string]string{
	"steam-gslt":          "https://steamcommunity.com/dev/managegameservers",
	"gamestats":           "https://gamestats.sandstorm.game/",
	"steamid-io":          "https://steamid.io/",
	"steam-profile":       "https://steamcommunity.com/my/profile",
	"modio":               "https://mod.io/g/insurgencysandstorm",
	"modio-subscriptions": "https://mod.io/g/insurgencysandstorm",
	"modio-api-access":    "https://mod.io/me/access",
	"modio-status":        "https://status.mod.io/",
	"modio-admin-guide":   "https://mod.io/g/insurgencysandstorm/r/server-admin-guide",
	"steamcmd-docs":       "https://developer.valvesoftware.com/wiki/SteamCMD",
	"radmin":              "https://www.radmin-vpn.com/",
	"server-guide":        "https://steamcommunity.com/sharedfiles/filedetails/?id=2102023759",
	"apache-license":      "https://www.apache.org/licenses/LICENSE-2.0",
	"wails":               "https://wails.io",
	"download-page":       upd.PageURL,
}

// Info é o que a tela mostra sobre o app.
type Info struct {
	AppVersion string   `json:"appVersion"`
	GoVersion  string   `json:"goVersion"`
	Languages  []string `json:"languages"`
	Problems   []string `json:"problems"`
}

// Info devolve versão e listas fixas.
func (s *SystemAPI) Info() Info {
	gv := ""
	if bi, ok := debug.ReadBuildInfo(); ok {
		gv = bi.GoVersion
	}
	return Info{AppVersion: AppVersion, GoVersion: gv, Languages: config.Languages, Problems: problems.Codes}
}

// OpenLink abre no navegador um dos endereços conhecidos (pela chave).
func (s *SystemAPI) OpenLink(key string) error {
	url, ok := links[key]
	if !ok {
		return problems.New(problems.Unknown, "unknown link")
	}
	runtime.BrowserOpenURL(s.a.ctx, url)
	return nil
}

// ProblemLinks devolve os links de ajuda de um erro do catálogo.
func (s *SystemAPI) ProblemLinks(code string) []string { return problems.Links(code) }

// Settings devolve as preferências do app.
func (s *SystemAPI) Settings() config.Settings { return s.a.cfg.Settings() }

// SetLanguage troca o idioma (en-US, pt-PT, es-ES, ru, uk).
func (s *SystemAPI) SetLanguage(lang string) (config.Settings, error) {
	return s.a.cfg.SetLanguage(lang)
}

// SetShowConsole mostra ou esconde a janela de console do servidor (vale no próximo início).
func (s *SystemAPI) SetShowConsole(on bool) (config.Settings, error) {
	return s.a.cfg.SetShowConsole(on)
}

// CopyText copia um texto para a área de transferência.
func (s *SystemAPI) CopyText(text string) error {
	if len(text) > 4096 {
		return problems.New(problems.Unknown, "text too long")
	}
	return runtime.ClipboardSetText(s.a.ctx, text)
}

// Step é um item dos primeiros passos.
type Step struct {
	ID       string `json:"id"`
	Done     bool   `json:"done"`
	Optional bool   `json:"optional"`
	// Page é a tela onde o passo se resolve.
	Page string `json:"page"`
}

// Checklist diz o que já está pronto para o servidor subir.
func (s *SystemAPI) Checklist() []Step {
	p := s.a.cfg.Profile()
	in := installer.Inspect(p)
	sec, _ := s.a.cfg.SecretsState()
	session := modioSession()
	return []Step{
		{ID: "serverDir", Done: p.ServerDir != "", Page: "install"},
		{ID: "steamcmd", Done: in.SteamCMDFound || in.ServerInstalled, Page: "install"},
		{ID: "serverInstalled", Done: in.ServerInstalled, Page: "install"},
		{ID: "hostname", Done: p.Hostname != server.Default().Hostname, Page: "server"},
		{ID: "modioAuth", Done: !p.UseMods || session.Valid, Optional: !p.UseMods, Page: "mods"},
		{ID: "admins", Done: len(p.Admins) > 0, Optional: true, Page: "server"},
		{ID: "connectIp", Done: p.ConnectIP != "", Optional: true, Page: "server"},
		{ID: "secretsSafe", Done: !s.a.cfg.RconInIni(), Page: "server"},
		{ID: "gslt", Done: sec["gsltToken"], Optional: true, Page: "server"},
	}
}

func modioSession() modio.Session {
	uf, err := modio.UserFile()
	if err != nil {
		return modio.Session{}
	}
	s, _, _ := modio.ReadSession(uf)
	return s
}

// Terms é a situação do aceite dos termos de uso e do aviso de dados (ADR-015).
type Terms struct {
	Version    string `json:"version"`
	Accepted   bool   `json:"accepted"`
	AcceptedAt string `json:"acceptedAt"`
}

// Terms diz qual versão está em vigor e se o usuário já a aceitou.
func (s *SystemAPI) Terms() Terms {
	st := s.a.cfg.Settings()
	t := Terms{Version: legal.Version, Accepted: legal.Accepted(st.TermsAccepted)}
	if t.Accepted {
		t.AcceptedAt = st.TermsAcceptedAt
	}
	return t
}

// AcceptTerms grava o aceite da versão mostrada e liga a procura de atualizações.
func (s *SystemAPI) AcceptTerms(version string) (Terms, error) {
	if _, err := s.a.cfg.AcceptTerms(version, time.Now()); err != nil {
		return s.Terms(), err
	}
	s.a.startUpdates()
	return s.Terms(), nil
}

// UpdateState devolve a situação da atualização automática.
func (s *SystemAPI) UpdateState() updater.State { return s.a.updates.State() }

// CheckUpdates procura uma versão nova agora (botão "Procurar atualizações").
func (s *SystemAPI) CheckUpdates() updater.State { return s.a.updates.Check(s.a.ctx) }

// InstallUpdate instala a versão já baixada e reabre o app nela. Recusa com instalação do servidor em andamento.
// O servidor do jogo continua rodando (é outro processo).
func (s *SystemAPI) InstallUpdate() error {
	if s.a.install.Busy() {
		return problems.New(problems.Busy, "installation in progress")
	}
	if err := s.a.updates.Install(); err != nil {
		return err
	}
	runtime.Quit(s.a.ctx) // o instalador espera o app fechar
	return nil
}
