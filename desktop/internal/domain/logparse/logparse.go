// Package logparse traduz linhas do Insurgency.log em eventos (ADR-004) e classifica o nível de cada linha.
// As regras foram tiradas de logs reais do servidor (09/11/2026).
package logparse

import (
	"regexp"
	"strconv"
	"strings"

	"sandstorm-server-assist/desktop/internal/domain/problems"
)

// Kind é o tipo de evento.
type Kind string

// Tipos de evento reconhecidos.
const (
	ModioReady     Kind = "modio-ready"
	Listening      Kind = "listening"
	ModsEnabled    Kind = "mods-enabled"
	ModActivated   Kind = "mod-activated"
	ModEvent       Kind = "mod-event"
	ModPending     Kind = "mod-pending"
	Travel         Kind = "travel"
	LoadMap        Kind = "load-map"
	MatchState     Kind = "match-state"
	Problem        Kind = "problem"
	GameStatsNoKey Kind = "gamestats-no-token"
)

// Event é o que uma linha do log significa.
type Event struct {
	Kind Kind `json:"kind"`
	// ModID e ModName para eventos de mod; Action para ModEvent (BeginInstall, Installed, BeginUninstall...).
	ModID   int    `json:"modId,omitempty"`
	ModName string `json:"modName,omitempty"`
	Action  string `json:"action,omitempty"`
	// Map para Travel e LoadMap ("Containment", "Farmhouse"); Scenario quando houver.
	Map      string `json:"map,omitempty"`
	Scenario string `json:"scenario,omitempty"`
	// From/To para MatchState.
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`
	Port int    `json:"port,omitempty"`
	// Problem é o código do catálogo quando a linha indica um erro conhecido.
	Problem string `json:"problem,omitempty"`
}

var (
	reModioReady  = regexp.MustCompile(`LogModioGame: ModSubsystem: mod\.io initialization complete`)
	reListening   = regexp.MustCompile(`IpNetDriver listening on port (\d+)`)
	reModsEnabled = regexp.MustCompile(`LogModioGame: ModSubsystem: ModManager enabled`)
	reActivated   = regexp.MustCompile(`OnModActivatedFromCloud: (.+)\[(\d+)\] bFromGameServer`)
	reModEvent    = regexp.MustCompile(`ModManager Event received\. ModId: (\d+), ModEvent: (\w+)`)
	rePending     = regexp.MustCompile(`Mod added, id: (\d+) to PendingChanges with state (\w+)`)
	reTravel      = regexp.MustCompile(`ProcessServerTravel: ([A-Za-z0-9_]+)(?:\?Scenario=([A-Za-z0-9_]+))?`)
	reLoadMap     = regexp.MustCompile(`LogLoad: LoadMap: /?(?:[^?\s]*/)?([A-Za-z0-9_]+)(?:\?[^\s]*Scenario=([A-Za-z0-9_]+))?`)
	reMatchState  = regexp.MustCompile(`Match State Changed from (\w+) to (\w+)`)
	reGameStats   = regexp.MustCompile(`GameStatsServer INIT: No authentication provided`)

	// Erros conhecidos (a ordem importa: o mais específico primeiro).
	problemRules = []struct {
		re   *regexp.Regexp
		code string
	}{
		{regexp.MustCompile(`(?i)Serial size mismatch|Corrupt data found`), problems.ModConflict},
		{regexp.MustCompile(`(?i)Invalid security code`), problems.ModioCodeInvalid},
		{regexp.MustCompile(`(?i)(could not|unable to|failed to) (bind|create socket)|address already in use|WSAEADDRINUSE`), problems.PortInUse},
		{regexp.MustCompile(`(?i)not enough space on the disk|disk (is )?full`), problems.DiskFull},
		{regexp.MustCompile(`LogWindows: Error: Fatal error`), problems.ServerCrashed},
	}

	reError   = regexp.MustCompile(`(?:^|\]|\s)Log\w+: (?:Error|Fatal):`)
	reWarning = regexp.MustCompile(`(?:^|\]|\s)Log\w+: Warning:`)

	// harmless são erros que o jogo grava sempre, mesmo com o servidor funcionando (vistos em logs reais de
	// 09/11/2026): cosméticos dos jogadores, replay e voz da NWI, pacotes do jogo base procurados fora do .pak.
	// Contam como aviso, para o número de erros mostrar só o que importa.
	harmless = regexp.MustCompile(`LogCosmetics: Error: GetSelectedCosmeticFromItemId|LogRedstone\w*: Error:|` +
		`LogHttpReplay: Error:|LogDemo: Error: UDemoNetDriver::TickFlush: ReplayStreamer|` +
		`LogStreaming: Error: (?:Couldn't find file for package /(?:Game|StatsCollector|Engine)/|This will hitch streaming|Found \d+ dependent packages)|` +
		`LogHydra: Error:|LogStreamableManager: Error: RequestAsyncLoad called with empty|LogAI: Error: \[INSAISoldierController|` +
		// Moldura do resumo de estatísticas, voz da Steam e recursos que dois mods de mapa trazem iguais (sem
		// "Serial size mismatch", o servidor segue normal).
		`LogGameStatsServer: Error: (?:-{5,}|--- GAME STATISTICS)|LogOnline: Error: STEAM \(NWI\): Failed initializing server voice|` +
		`LogAssetRegistry: Error: AddAssetData called with ObjectPath .* which already exists|` +
		// Mapas de mod empacotados sem alguns objetos internos (paisagem, partículas). O mapa carrega normal: na
		// Containment são centenas destas linhas por partida.
		`LogStreaming: Error: Could not find outer \S+ to create`)
)

// Parse devolve o evento de uma linha, se ela significar algo para o app.
func Parse(line string) (Event, bool) {
	// Atalho: quase todas as linhas não interessam.
	if !strings.Contains(line, "Log") {
		return Event{}, false
	}
	for _, r := range problemRules {
		if r.re.MatchString(line) {
			return Event{Kind: Problem, Problem: r.code}, true
		}
	}
	if reModioReady.MatchString(line) {
		return Event{Kind: ModioReady}, true
	}
	if m := reListening.FindStringSubmatch(line); m != nil {
		port, _ := strconv.Atoi(m[1])
		return Event{Kind: Listening, Port: port}, true
	}
	if reModsEnabled.MatchString(line) {
		return Event{Kind: ModsEnabled}, true
	}
	if m := reActivated.FindStringSubmatch(line); m != nil {
		id, _ := strconv.Atoi(m[2])
		return Event{Kind: ModActivated, ModID: id, ModName: strings.TrimSpace(m[1])}, true
	}
	if m := reModEvent.FindStringSubmatch(line); m != nil {
		id, _ := strconv.Atoi(m[1])
		return Event{Kind: ModEvent, ModID: id, Action: m[2]}, true
	}
	if m := rePending.FindStringSubmatch(line); m != nil {
		id, _ := strconv.Atoi(m[1])
		return Event{Kind: ModPending, ModID: id, Action: m[2]}, true
	}
	if m := reTravel.FindStringSubmatch(line); m != nil {
		return Event{Kind: Travel, Map: m[1], Scenario: m[2]}, true
	}
	if m := reLoadMap.FindStringSubmatch(line); m != nil {
		return Event{Kind: LoadMap, Map: m[1], Scenario: m[2]}, true
	}
	if m := reMatchState.FindStringSubmatch(line); m != nil {
		return Event{Kind: MatchState, From: m[1], To: m[2]}, true
	}
	if reGameStats.MatchString(line) {
		return Event{Kind: GameStatsNoKey}, true
	}
	return Event{}, false
}

// Level classifica a linha para o filtro da tela: "error", "warning" ou "info".
func Level(line string) string {
	switch {
	case harmless.MatchString(line):
		return "warning"
	case reError.MatchString(line):
		return "error"
	case reWarning.MatchString(line):
		return "warning"
	}
	return "info"
}

// Noise diz se a linha é ruído de inicialização (centenas de linhas de montagem dos .pak do jogo base). A tela
// esconde essas linhas no filtro "Importantes".
func Noise(line string) bool {
	return strings.Contains(line, "LogPakFile:") || strings.Contains(line, "LogShaderLibrary:") ||
		strings.Contains(line, "LogStreaming: Display") || strings.Contains(line, "LogConfig:")
}
