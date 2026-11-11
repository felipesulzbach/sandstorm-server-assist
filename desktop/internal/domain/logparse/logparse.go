// Package logparse traduz linhas do Insurgency.log em eventos (ADR-004) e classifica o nível de cada linha.
// As regras foram tiradas de logs reais do servidor (09/11/2026).
package logparse

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"sandstorm-server-assist/desktop/internal/domain/modcheck"
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
	// Jogadores (mensagens de entrada e saída, ADR-020).
	PlayerLogin  Kind = "player-login"
	PlayerJoined Kind = "player-joined"
	PlayerLeft   Kind = "player-left"
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
	// Object é o objeto do jogo citado num erro de mod ("/Game/..."), para achar o mod que o traz (ADR-039).
	Object string `json:"object,omitempty"`
	// Player é o nome do jogador (PlayerLogin, PlayerJoined); PlayerID é o id da plataforma (PlayerLogin, PlayerLeft).
	Player   string `json:"player,omitempty"`
	PlayerID string `json:"playerId,omitempty"`
	// Time é a hora da linha (UTC); zero se a linha não tiver.
	Time time.Time `json:"-"`
}

var (
	reModioReady  = regexp.MustCompile(`LogModioGame: ModSubsystem: mod\.io initialization complete`)
	reListening   = regexp.MustCompile(`IpNetDriver listening on port (\d+)`)
	reModsEnabled = regexp.MustCompile(`LogModioGame: ModSubsystem: ModManager enabled`)
	reActivated   = regexp.MustCompile(`OnModActivatedFromCloud: (.+)\[(\d+)\] bFromGameServer`)
	reModEvent    = regexp.MustCompile(`ModManager Event received\. ModId: (\d+), ModEvent: (\w+)`)
	rePending     = regexp.MustCompile(`Mod added, id: (\d+) to PendingChanges with state (\w+)`)
	// Mapa oficial vem pelo nome ("Border?Scenario=..."); mapa de mod pelo caminho ("/Estates/Maps/Estates?Scenario=...",
	// visto em 11/11/2026 com travelscenario). O cenário pode vir depois de outras opções.
	reTravel     = regexp.MustCompile(`ProcessServerTravel: /?(?:[^?\s]*/)?([A-Za-z0-9_]+)(?:\?[^\s]*?Scenario=([A-Za-z0-9_]+))?`)
	reLoadMap    = regexp.MustCompile(`LogLoad: LoadMap: /?(?:[^?\s]*/)?([A-Za-z0-9_]+)(?:\?[^\s]*Scenario=([A-Za-z0-9_]+))?`)
	reMatchState = regexp.MustCompile(`Match State Changed from (\w+) to (\w+)`)
	reGameStats  = regexp.MustCompile(`GameStatsServer INIT: No authentication provided`)

	// Jogadores, como aparecem no log (10/11/2026): o pedido de entrada traz o nome completo e o id; a entrada
	// confirmada traz o nome (às vezes cortado); a saída traz só o id, na conexão do jogo com um controle de jogador.
	reLogin  = regexp.MustCompile(`LogNet: Login request: \S*?\?Name=(.+?) userId: ([A-Za-z]+:[0-9A-Za-z_-]{1,64}) platform:`)
	reJoined = regexp.MustCompile(`LogNet: Join succeeded: (.+)$`)
	reLeft   = regexp.MustCompile(`LogNet: UNetConnection::Close: .*Driver: GameNetDriver .*PC: INSPlayerController_\d+.*UniqueId: ([A-Za-z]+:[0-9A-Za-z_-]{1,64}),`)
	reStamp  = regexp.MustCompile(`^\[(\d{4})\.(\d{2})\.(\d{2})-(\d{2})\.(\d{2})\.(\d{2}):(\d{3})\]`)

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
		`LogStreaming: Error: Could not find outer \S+ to create|` +
		// Mapa sem a distância máxima de raycast definida: o jogo usa 1 km e segue normal (Border, 11/11/2026).
		`LogWorldSettings: Error: No "Max Raycast Distance" defined`)
)

// Parse devolve o evento de uma linha, se ela significar algo para o app.
func Parse(line string) (Event, bool) {
	// Atalho: quase todas as linhas não interessam.
	if !strings.Contains(line, "Log") {
		return Event{}, false
	}
	for _, r := range problemRules {
		if r.re.MatchString(line) {
			ev := Event{Kind: Problem, Problem: r.code}
			if r.code == problems.ModConflict {
				if objs := modcheck.Assets(line); len(objs) > 0 {
					ev.Object = objs[0]
				}
			}
			return ev, true
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
	if m := reLogin.FindStringSubmatch(line); m != nil {
		return Event{Kind: PlayerLogin, Player: strings.TrimSpace(m[1]), PlayerID: m[2], Time: Stamp(line)}, true
	}
	if m := reJoined.FindStringSubmatch(line); m != nil {
		return Event{Kind: PlayerJoined, Player: strings.TrimSpace(m[1]), Time: Stamp(line)}, true
	}
	if m := reLeft.FindStringSubmatch(line); m != nil {
		return Event{Kind: PlayerLeft, PlayerID: m[1], Time: Stamp(line)}, true
	}
	return Event{}, false
}

// Stamp lê a hora do começo da linha ([2026.11.10-02.20.43:923], em UTC). Zero se a linha não tiver.
func Stamp(line string) time.Time {
	m := reStamp.FindStringSubmatch(line)
	if m == nil {
		return time.Time{}
	}
	var n [7]int
	for i := range n {
		n[i], _ = strconv.Atoi(m[i+1])
	}
	return time.Date(n[0], time.Month(n[1]), n[2], n[3], n[4], n[5], n[6]*int(time.Millisecond), time.UTC)
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
