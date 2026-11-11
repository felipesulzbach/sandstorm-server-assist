package server

import (
	"regexp"
	"sort"
	"strconv"
	"strings"

	"sandstorm-server-assist/desktop/internal/domain/problems"
)

// Troca de mapa com o servidor ligado (ADR-019): o app pede ao servidor, pelo RCON, a lista de cenários que ele
// conhece (comando "scenarios", que inclui os mapas dos mods montados) e troca com "travelscenario <cenário>". O
// travelscenario mantém o resto da partida (senha, vagas, mutators), ao contrário de "travel <url>".

// Origens de um cenário na lista.
const (
	// FromServer: o servidor listou (está carregado e pronto para trocar).
	FromServer = "server"
	// FromConfig: está no MapCycle ou é o mapa inicial/de mod do perfil, mas o servidor não listou. Só usado quando
	// o app não sabe quais mods estão carregados.
	FromConfig = "config"
	// FromMod: cenário de um mod carregado nesta partida (lido dos .pak do mod, ADR-029).
	FromMod = "mod"
	// FromMissing: está na rotação (ou é o mapa de mod do perfil), mas o mod dele não está carregado (ADR-030).
	// Aparece na lista para o usuário entender, sem poder ser escolhido.
	FromMissing = "missing"
)

// Scenario é uma opção de troca de mapa.
type Scenario struct {
	ID string `json:"id"`
	// Map é a parte do nome antes do modo (Scenario_Containment_Checkpoint_Security → "Containment").
	Map string `json:"map"`
	// Mode é o resto ("Checkpoint_Security"); vazio se o nome não tem um modo conhecido.
	Mode   string `json:"mode"`
	Source string `json:"source"`
}

// scenarioToken acha nomes de cenário num texto (resposta do servidor ou linha do MapCycle).
var scenarioToken = regexp.MustCompile(`\bScenario_[A-Za-z0-9_]{1,80}`)

// modeWords são as palavras que começam o modo no nome do cenário (modos do jogo base; mods usam os mesmos).
var modeWords = map[string]bool{
	"Checkpoint": true, "Push": true, "Firefight": true, "Skirmish": true, "Domination": true,
	"Frontline": true, "Ambush": true, "Outpost": true, "Survival": true, "Defusal": true,
	"TeamDeathmatch": true, "Elimination": true, "Hardcore": true, "Competitive": true, "Conquer": true,
	"FFA": true, "TDM": true, "Team": true, "Range": true, "Tutorial": true, "Interception": true,
}

// ParseScenarios tira os nomes de cenário de um texto, sem repetir, em ordem.
func ParseScenarios(text string) []string {
	seen := map[string]bool{}
	var out []string
	for _, id := range scenarioToken.FindAllString(text, -1) {
		if !seen[id] && scenarioRe.MatchString(id) {
			seen[id] = true
			out = append(out, id)
		}
	}
	sort.Strings(out)
	return out
}

// SplitScenario separa o mapa e o modo do nome do cenário.
func SplitScenario(id string) (mapName, mode string) {
	rest := strings.TrimPrefix(id, "Scenario_")
	parts := strings.Split(rest, "_")
	for i := 1; i < len(parts); i++ {
		if modeWords[parts[i]] {
			return strings.Join(parts[:i], "_"), strings.Join(parts[i:], "_")
		}
	}
	return rest, ""
}

// MapChoices junta o que o servidor listou com os mapas de mod. O "scenarios" do servidor só lista os mapas do jogo
// base (visto em 10/11/2026 com o Containment carregado), então os mapas de mod entram por aqui:
//   - fromMods são os cenários dos mods carregados nesta partida (ADR-029). Um mapa da rotação cujo mod não está
//     carregado (inscrição retirada, por exemplo) fica de fora: o servidor recusaria a troca.
//   - fromMods nil: o app não sabe quais mods estão carregados; usa os cenários do perfil (MapCycle, mapa inicial e
//     mapa de mod), como antes.
func MapChoices(fromServer []string, p Profile, fromMods []string) []Scenario {
	seen := map[string]bool{}
	out := []Scenario{}
	add := func(id, source string) {
		if seen[id] || !scenarioRe.MatchString(id) {
			return
		}
		seen[id] = true
		m, mode := SplitScenario(id)
		out = append(out, Scenario{ID: id, Map: m, Mode: mode, Source: source})
	}
	for _, id := range fromServer {
		add(id, FromServer)
	}
	// Do perfil: sem saber os mods carregados, entram como opção; sabendo, o que sobrou não tem mod carregado.
	fromProfile := FromConfig
	if fromMods != nil {
		for _, id := range fromMods {
			add(id, FromMod)
		}
		fromProfile = FromMissing
	}
	for _, line := range p.MapCycle {
		if id := scenarioToken.FindString(FixMapCycleLine(line)); id != "" {
			add(id, fromProfile)
		}
	}
	add(p.TravelScenario, fromProfile)
	add(p.StartScenario, fromProfile)
	sort.Slice(out, func(i, j int) bool {
		a, b := out[i], out[j]
		if !strings.EqualFold(a.Map, b.Map) {
			return strings.ToLower(a.Map) < strings.ToLower(b.Map)
		}
		return a.Mode < b.Mode
	})
	return out
}

// TravelCommand é o comando de troca para um cenário da lista. Só aceita um id que está na lista (a tela manda
// só o id; nada de texto livre vai para o servidor).
func TravelCommand(id string, choices []Scenario) (string, error) {
	if !scenarioRe.MatchString(id) {
		return "", problems.Field("scenario", "validation.scenario")
	}
	for _, c := range choices {
		if c.ID == id {
			if c.Source == FromMissing {
				return "", problems.Field("scenario", "validation.scenarioNotLoaded")
			}
			return "travelscenario " + id, nil
		}
	}
	return "", problems.Field("scenario", "validation.scenarioUnknown")
}

// MissingInCycle diz, linha a linha, se o cenário da rotação não existe no jogo base nem nos mods instalados
// (ADR-030): quando a rotação chega nele, o servidor não tem para onde ir. known são os cenários que existem; vazio
// quer dizer que o app não conseguiu ler o jogo base, e nada é marcado. Linhas que o app não entende não são marcadas.
func MissingInCycle(lines []string, known []string) []bool {
	out := make([]bool, len(lines))
	if len(known) == 0 {
		return out
	}
	set := make(map[string]bool, len(known))
	for _, k := range known {
		set[strings.ToLower(k)] = true
	}
	for i, l := range lines {
		if id := LineScenario(l); id != "" && !set[strings.ToLower(id)] {
			out[i] = true
		}
	}
	return out
}

// TravelWithLighting monta a troca de mapa com dia ou noite (ADR-022). O travelscenario não aceita a iluminação,
// então vai "travel <mapa>?Scenario=<cenário>?Lighting=<Day|Night>" com as vagas e a senha da partida (o travel
// troca a URL inteira). Sem iluminação, é o travelscenario de sempre.
func TravelWithLighting(id, lighting string, choices []Scenario, maxPlayers int, password string) (string, error) {
	cmd, err := TravelCommand(id, choices)
	if err != nil || lighting == "" {
		return cmd, err
	}
	if lighting != LightingDay && lighting != LightingNight {
		return "", problems.Field("lighting", "validation.lighting")
	}
	m := ""
	for _, c := range choices {
		if c.ID == id {
			m = c.Map
		}
	}
	if !identRe.MatchString(m) {
		return "", problems.Field("lighting", "validation.lightingNoMap")
	}
	url := m + "?Scenario=" + id + "?Lighting=" + lighting + "?MaxPlayers=" + strconv.Itoa(maxPlayers)
	if password != "" {
		url += "?Password=" + password
	}
	return "travel " + url, nil
}

// reTravelRefused reconhece a resposta do servidor quando ele não aceita a troca.
var reTravelRefused = regexp.MustCompile(`(?i)\b(unknown|invalid|not found|failed|could ?n[o']t|cannot|no such|usage)\b`)

// TravelRefused diz se a resposta do servidor ao travelscenario é uma recusa.
func TravelRefused(resp string) bool { return reTravelRefused.MatchString(resp) }
