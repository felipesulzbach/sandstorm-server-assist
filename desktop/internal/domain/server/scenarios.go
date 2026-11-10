package server

import (
	"regexp"
	"sort"
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
	// FromConfig: está no MapCycle ou é o mapa inicial/de mod do perfil, mas o servidor não listou (mapas de mod).
	FromConfig = "config"
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

// MapChoices junta o que o servidor listou com os cenários do perfil (MapCycle, mapa inicial e mapa de mod) que
// ele não listou. O "scenarios" do servidor só lista os mapas do jogo base (visto em 10/11/2026 com o Containment
// carregado), então os mapas de mod entram por aqui.
func MapChoices(fromServer []string, p Profile) []Scenario {
	seen := map[string]bool{}
	var out []Scenario
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
	for _, line := range p.MapCycle {
		if id := scenarioToken.FindString(FixMapCycleLine(line)); id != "" {
			add(id, FromConfig)
		}
	}
	add(p.TravelScenario, FromConfig)
	add(p.StartScenario, FromConfig)
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
			return "travelscenario " + id, nil
		}
	}
	return "", problems.Field("scenario", "validation.scenarioUnknown")
}

// reTravelRefused reconhece a resposta do servidor quando ele não aceita a troca.
var reTravelRefused = regexp.MustCompile(`(?i)\b(unknown|invalid|not found|failed|could ?n[o']t|cannot|no such|usage)\b`)

// TravelRefused diz se a resposta do servidor ao travelscenario é uma recusa.
func TravelRefused(resp string) bool { return reTravelRefused.MatchString(resp) }
