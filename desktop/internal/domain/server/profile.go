// Package server tem as regras do servidor dedicado: o perfil (o que o usuário configura), a validação de cada
// campo e a linha de comando do InsurgencyServer.exe (ADR-003). Não faz I/O.
// Mensagens de validação são chaves de tradução ("validation.*", ADR-009).
package server

import (
	"net"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"sandstorm-server-assist/desktop/internal/domain/problems"
	"sandstorm-server-assist/desktop/internal/domain/steamid"
)

// AppID é o app do servidor dedicado no SteamCMD; GameAppID é o jogo (usado no GSLT).
const (
	AppID     = 581330
	GameAppID = 581320
)

// Admin é um administrador do servidor (Admins.txt).
type Admin struct {
	SteamID string `json:"steamId"`
	Name    string `json:"name"`
}

// Profile é tudo que o usuário configura, sem segredos (os segredos ficam em Secrets, no cofre).
type Profile struct {
	ServerDir   string `json:"serverDir"`
	SteamCMDDir string `json:"steamcmdDir"`

	Hostname   string `json:"hostname"`
	Port       int    `json:"port"`
	QueryPort  int    `json:"queryPort"`
	MaxPlayers int    `json:"maxPlayers"`

	StartMap      string `json:"startMap"`
	StartScenario string `json:"startScenario"`
	// StartLighting é a iluminação do mapa inicial (?Lighting=Day|Night). Vazio = padrão do cenário.
	StartLighting string `json:"startLighting"`
	// Mapa de mod para onde o servidor viaja depois de montar os mods (-ModDownloadTravelTo). Vazio = não viaja.
	TravelMap      string `json:"travelMap"`
	TravelScenario string `json:"travelScenario"`

	UseMods  bool     `json:"useMods"`
	Mutators []string `json:"mutators"`

	Rcon     bool `json:"rcon"`
	RconPort int  `json:"rconPort"`

	// BindIP é a placa de rede (ou a rede de um programa de VPN) em que o servidor escuta (-MultiHome). Vazio = todas.
	BindIP string `json:"bindIp"`

	// ConnectIP é o IP que os jogadores usam (IP público ou o IP de um programa de VPN). Só para mostrar na tela.
	ConnectIP string `json:"connectIp"`

	Admins   []Admin  `json:"admins"`
	MapCycle []string `json:"mapCycle"`
	// ModIDs são os mods esperados (informativo: quem decide é a inscrição da conta do servidor no mod.io).
	ModIDs []int `json:"modIds"`
	// CmdModList manda os ModIDs ao servidor (-CmdModList). Desligado = o servidor usa as inscrições da conta.
	CmdModList bool `json:"cmdModList"`

	// Mensagens (ADR-020). Motd aparece na tela de carregamento (Motd.txt). JoinMessage e LeaveMessage vão para o
	// chat pelo RCON quando um jogador entra ou sai; {name} vira o nome do jogador. Vazio = desligado.
	Motd         string `json:"motd"`
	JoinMessage  string `json:"joinMessage"`
	LeaveMessage string `json:"leaveMessage"`
	// Announcements vão para o chat, uma por vez, a cada AnnounceEvery minutos (ADR-022).
	Announcements []string `json:"announcements"`
	AnnounceEvery int      `json:"announceEvery"`

	// Rules são as opções do Game.ini e do Engine.ini (rules.go). nil = o app não mexe nelas.
	Rules map[string]string `json:"rules"`

	// AutoRestart: reinício diário, ao travar e atualização pelo SteamCMD (ADR-022).
	AutoRestart AutoRestart `json:"autoRestart"`
	// Discord: textos dos avisos por webhook (vazio = aviso desligado). O endereço fica no cofre.
	Discord Discord `json:"discord"`
}

// AutoRestart é o reinício automático do servidor.
type AutoRestart struct {
	// Daily é a hora do reinício diário ("HH:MM", hora local). Vazio = desligado.
	Daily string `json:"daily"`
	// WarnMessage vai para o chat 5 e 1 minuto antes ({minutes} vira o número). Vazio = sem aviso. Precisa do RCON.
	WarnMessage string `json:"warnMessage"`
	// Update atualiza o servidor pelo SteamCMD no reinício diário.
	Update bool `json:"update"`
	// OnCrash inicia de novo quando o servidor fecha sozinho (até 3 vezes em 30 minutos).
	OnCrash bool `json:"onCrash"`
}

// Discord são os textos dos avisos. {server} vira o nome do servidor, {map} o mapa e {name} o jogador.
type Discord struct {
	Ready   string `json:"ready"`
	Stopped string `json:"stopped"`
	Crashed string `json:"crashed"`
	Join    string `json:"join"`
	Leave   string `json:"leave"`
}

// Secrets ficam só no cofre DPAPI e no backup cifrado (ADR-002).
type Secrets struct {
	ServerPassword string `json:"serverPassword"`
	RconPassword   string `json:"rconPassword"`
	GSLTToken      string `json:"gsltToken"`
	GameStatsToken string `json:"gameStatsToken"`
	ModioEmail     string `json:"modioEmail"`
	// ModioAPIKey é a chave de API do mod.io (não vem no código: o projeto é público, ADR-010).
	ModioAPIKey string `json:"modioApiKey"`
	// DiscordWebhook é o endereço do webhook do Discord (quem tem o endereço consegue postar no canal).
	DiscordWebhook string `json:"discordWebhook"`
}

// SecretKeys são os nomes aceitos para ler, mostrar e trocar um segredo (a tela manda só a chave).
var SecretKeys = []string{"serverPassword", "rconPassword", "gsltToken", "gameStatsToken", "modioEmail", "modioApiKey", "discordWebhook"}

func (s *Secrets) field(key string) *string {
	switch key {
	case "serverPassword":
		return &s.ServerPassword
	case "rconPassword":
		return &s.RconPassword
	case "gsltToken":
		return &s.GSLTToken
	case "gameStatsToken":
		return &s.GameStatsToken
	case "modioEmail":
		return &s.ModioEmail
	case "modioApiKey":
		return &s.ModioAPIKey
	case "discordWebhook":
		return &s.DiscordWebhook
	}
	return nil
}

// Get devolve um segredo pela chave.
func (s Secrets) Get(key string) (string, bool) {
	if f := s.field(key); f != nil {
		return *f, true
	}
	return "", false
}

// With devolve uma cópia com um segredo trocado (valor validado; vazio apaga).
func (s Secrets) With(key, value string) (Secrets, error) {
	value = strings.TrimSpace(value)
	f := s.field(key)
	if f == nil {
		return s, problems.Field(key, "validation.unknownField")
	}
	if value != "" {
		if err := ValidateSecret(key, value); err != nil {
			return s, err
		}
	}
	*f = value
	return s, nil
}

// Values devolve os segredos preenchidos (para mascarar em log).
func (s Secrets) Values() []string {
	return []string{s.ServerPassword, s.RconPassword, s.GSLTToken, s.GameStatsToken, s.ModioAPIKey, s.DiscordWebhook}
}

// Default é o perfil inicial, igual ao que funcionou no servidor configurado à mão (resumo_sessao.md).
func Default() Profile {
	return Profile{
		Hostname:      "My Sandstorm Server",
		Port:          27102,
		QueryPort:     27131,
		MaxPlayers:    8,
		StartMap:      "Farmhouse",
		StartScenario: "Scenario_Farmhouse_Checkpoint_Security",
		UseMods:       true,
		RconPort:      27015,
		AnnounceEvery: 15,
	}
}

var (
	identRe    = regexp.MustCompile(`^[A-Za-z0-9_]{1,80}$`)
	mutatorRe  = regexp.MustCompile(`^[A-Za-z0-9_]{1,64}$`)
	scenarioRe = regexp.MustCompile(`^Scenario_[A-Za-z0-9_]{1,80}$`)
	// Senhas vão na URL (?Password=) e na linha de comando: sem espaço, aspas, ?, & e =.
	passwordRe = regexp.MustCompile(`^[A-Za-z0-9!@#$%*._+\-]{4,64}$`)
	tokenRe    = regexp.MustCompile(`^[A-Za-z0-9_\-]{8,128}$`)
	apiKeyRe   = regexp.MustCompile(`^[0-9a-fA-F]{32}$`)
	emailRe    = regexp.MustCompile(`^[^\s@"<>]{1,64}@[A-Za-z0-9.\-]{1,190}\.[A-Za-z]{2,24}$`)
	codeRe     = regexp.MustCompile(`^[A-Za-z0-9]{5}$`)
	// Webhook do Discord: só os endereços oficiais.
	webhookRe = regexp.MustCompile(`^https://(?:(?:ptb|canary)\.)?discord(?:app)?\.com/api/webhooks/[0-9]{5,25}/[A-Za-z0-9_\-]{20,120}$`)
)

// Erros de digitação comuns em domínios de e-mail.
var emailTypos = map[string]string{
	"gamil.com": "gmail.com", "gmial.com": "gmail.com", "gmail.con": "gmail.com", "gmal.com": "gmail.com",
	"hotmial.com": "hotmail.com", "hotmail.con": "hotmail.com", "outlok.com": "outlook.com",
	"yahoo.con": "yahoo.com", "yaho.com": "yahoo.com",
}

// ValidateSecret confere o formato de um segredo antes de guardar.
func ValidateSecret(key, v string) error {
	switch key {
	case "serverPassword", "rconPassword":
		if !passwordRe.MatchString(v) {
			return problems.Field(key, "validation.password")
		}
	case "gsltToken", "gameStatsToken":
		if !tokenRe.MatchString(v) {
			return problems.Field(key, "validation.token")
		}
	case "modioApiKey":
		if !apiKeyRe.MatchString(v) {
			return problems.Field(key, "validation.apiKey")
		}
	case "modioEmail":
		return ValidateEmail(v)
	case "discordWebhook":
		if !webhookRe.MatchString(v) {
			return problems.Field(key, "validation.webhook")
		}
	default:
		return problems.Field(key, "validation.unknownField")
	}
	return nil
}

// ValidateEmail confere o formato e avisa erros comuns de digitação no domínio.
func ValidateEmail(v string) error {
	if !emailRe.MatchString(v) {
		return &problems.Error{Code: problems.ModioEmailInvalid, Field: "modioEmail", Msg: "validation.email"}
	}
	domain := strings.ToLower(v[strings.LastIndex(v, "@")+1:])
	if fix, ok := emailTypos[domain]; ok {
		return &problems.Error{Code: problems.ModioEmailInvalid, Field: "modioEmail", Msg: "validation.emailTypo",
			Params: map[string]string{"suggestion": fix}}
	}
	return nil
}

// ValidCode diz se o código do mod.io tem o formato certo (5 letras ou números).
func ValidCode(c string) bool { return codeRe.MatchString(c) }

// Validate confere o perfil inteiro. Devolve o primeiro campo inválido.
func (p Profile) Validate() error {
	if p.ServerDir != "" && !filepath.IsAbs(p.ServerDir) {
		return problems.Field("serverDir", "validation.absolutePath")
	}
	if err := validHostname(p.Hostname); err != nil {
		return err
	}
	for _, f := range []struct {
		name string
		v    int
	}{{"port", p.Port}, {"queryPort", p.QueryPort}, {"rconPort", p.RconPort}} {
		if f.v < 1024 || f.v > 65535 {
			return problems.Field(f.name, "validation.port")
		}
	}
	if p.Port == p.QueryPort || (p.Rcon && (p.RconPort == p.Port || p.RconPort == p.QueryPort)) {
		return problems.Field("queryPort", "validation.portsDistinct")
	}
	if p.MaxPlayers < 1 || p.MaxPlayers > 64 {
		return problems.Field("maxPlayers", "validation.maxPlayers")
	}
	if !identRe.MatchString(p.StartMap) {
		return problems.Field("startMap", "validation.map")
	}
	if !scenarioRe.MatchString(p.StartScenario) {
		return problems.Field("startScenario", "validation.scenario")
	}
	if p.TravelMap != "" || p.TravelScenario != "" {
		if !identRe.MatchString(p.TravelMap) {
			return problems.Field("travelMap", "validation.map")
		}
		if !scenarioRe.MatchString(p.TravelScenario) {
			return problems.Field("travelScenario", "validation.scenario")
		}
		if !p.UseMods {
			return problems.Field("travelMap", "validation.travelNeedsMods")
		}
	}
	for _, m := range p.Mutators {
		if !mutatorRe.MatchString(m) {
			return problems.Field("mutators", "validation.mutator", "value", m)
		}
	}
	if p.ConnectIP != "" && net.ParseIP(p.ConnectIP) == nil {
		return problems.Field("connectIp", "validation.ip")
	}
	if p.BindIP != "" {
		if ip := net.ParseIP(p.BindIP); ip == nil || ip.To4() == nil || ip.IsUnspecified() {
			return problems.Field("bindIp", "validation.bindIp")
		}
	}
	if err := validateLighting("startLighting", p.StartLighting); err != nil {
		return err
	}
	for _, a := range p.Admins {
		if !steamid.Valid(a.SteamID) {
			return problems.Field("admins", "validation.steamId", "value", a.SteamID)
		}
		if strings.ContainsAny(a.Name, "\r\n") || utf8.RuneCountInString(a.Name) > 40 {
			return problems.Field("admins", "validation.adminName")
		}
	}
	for _, line := range p.MapCycle {
		if strings.ContainsAny(line, "\r\n") || len(line) > 200 {
			return problems.Field("mapCycle", "validation.mapCycleLine")
		}
	}
	for _, id := range p.ModIDs {
		if id <= 0 {
			return problems.Field("modIds", "validation.modId", "value", strconv.Itoa(id))
		}
	}
	if len(p.ModIDs) > 200 {
		return problems.Field("modIds", "validation.modIdsMax", "max", "200")
	}
	if p.CmdModList && len(p.ModIDs) == 0 {
		return problems.Field("modIds", "validation.modListEmpty")
	}
	if err := p.validateMessages(); err != nil {
		return err
	}
	if err := p.validateMarkup(); err != nil {
		return err
	}
	if err := p.validateAutomation(); err != nil {
		return err
	}
	return p.validateRules()
}

func validHostname(h string) error {
	n := utf8.RuneCountInString(h)
	if n == 0 || n > 64 {
		return problems.Field("hostname", "validation.hostnameLength")
	}
	// Vai entre aspas na linha de comando: sem aspas, sem caracteres de controle.
	for _, r := range h {
		if r == '"' || r < 32 || r == 127 {
			return problems.Field("hostname", "validation.hostnameQuotes")
		}
	}
	return nil
}
