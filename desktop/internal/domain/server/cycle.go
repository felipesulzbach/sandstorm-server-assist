package server

import (
	"sort"
	"strings"

	"sandstorm-server-assist/desktop/internal/domain/problems"
)

// Rotação por mapa (ADR-046). O MapCycle.txt continua uma lista simples, uma entrada por linha; a tela mostra as
// linhas em blocos por mapa e muda a lista por estas funções, que só mexem nas linhas pedidas.

// MaxCycleLines limita a rotação que a tela mostra e muda.
const MaxCycleLines = 500

// CycleScenario é um cenário que pode entrar na rotação.
type CycleScenario struct {
	ID string `json:"id"`
	// Mode é o resto do nome depois do mapa ("Checkpoint_Security"); vazio se o nome não tem um modo conhecido.
	Mode string `json:"mode"`
}

// CycleMap é um mapa que pode entrar na rotação, com os cenários dele.
type CycleMap struct {
	// Map é a parte do nome do cenário antes do modo (Scenario_Farmhouse_Checkpoint_Security → "Farmhouse").
	Map   string `json:"map"`
	Label string `json:"label"`
	// ModID e ModName: mapa de um mod instalado (0 e vazio no jogo base).
	ModID     int             `json:"modId"`
	ModName   string          `json:"modName"`
	Scenarios []CycleScenario `json:"scenarios"`
}

// CycleSource é um cenário de um mod instalado.
type CycleSource struct {
	Scenario string
	ModID    int
	ModName  string
}

// CycleCatalog junta os mapas do jogo base e dos mods instalados que podem entrar na rotação. base são os cenários
// lidos dos .pak do servidor; sem eles (servidor não instalado), usa a lista Maps nos modos de Checkpoint.
func CycleCatalog(base []string, mods []CycleSource) []CycleMap {
	if len(base) == 0 {
		for _, m := range Maps {
			for _, md := range Modes {
				base = append(base, "Scenario_"+m.Scenario+"_"+md.ID)
			}
		}
	}
	byMap := map[string]*CycleMap{}
	var order []string
	seen := map[string]bool{}
	add := func(id string, modID int, modName string) {
		key := strings.ToLower(id)
		if seen[key] || !scenarioRe.MatchString(id) {
			return
		}
		seen[key] = true
		name, mode := SplitScenario(id)
		if modID == 0 && strings.HasPrefix(name, "DEV_") {
			return // cenários de teste dos desenvolvedores do jogo (Scenario_DEV_...), não são para jogar
		}
		k := strings.ToLower(name)
		if modID != 0 {
			k = strings.ToLower(name) + "#mod"
		}
		cm, ok := byMap[k]
		if !ok {
			cm = &CycleMap{Map: name, Label: mapLabel(name), ModID: modID, ModName: modName}
			byMap[k] = cm
			order = append(order, k)
		}
		cm.Scenarios = append(cm.Scenarios, CycleScenario{ID: id, Mode: mode})
	}
	for _, id := range base {
		add(id, 0, "")
	}
	for _, s := range mods {
		add(s.Scenario, s.ModID, s.ModName)
	}
	out := make([]CycleMap, 0, len(order))
	for _, k := range order {
		cm := *byMap[k]
		sort.Slice(cm.Scenarios, func(i, j int) bool { return modeOrder(cm.Scenarios[i].Mode, cm.Scenarios[j].Mode) })
		out = append(out, cm)
	}
	sort.SliceStable(out, func(i, j int) bool {
		if (out[i].ModID == 0) != (out[j].ModID == 0) {
			return out[i].ModID != 0 // mapas de mod primeiro: é o que a pessoa costuma procurar
		}
		return strings.ToLower(out[i].Label) < strings.ToLower(out[j].Label)
	})
	return out
}

// modeOrder põe Checkpoint Security antes de Insurgents e o resto em ordem alfabética.
func modeOrder(a, b string) bool {
	ra, rb := modeRank(a), modeRank(b)
	if ra != rb {
		return ra < rb
	}
	return a < b
}

func modeRank(m string) int {
	switch {
	case m == "Checkpoint_Security":
		return 0
	case m == "Checkpoint_Insurgents":
		return 1
	case strings.HasPrefix(m, "Checkpoint"):
		return 2
	}
	return 3
}

// mapLabel é o nome do mapa para a tela: o da lista Maps, ou o nome do cenário com espaços.
func mapLabel(name string) string {
	for _, m := range Maps {
		if strings.EqualFold(m.Scenario, name) {
			return m.Label
		}
	}
	return strings.ReplaceAll(name, "_", " ")
}

// CycleEntry é uma linha da rotação.
type CycleEntry struct {
	// Index é a posição da linha na lista (para mudar).
	Index    int    `json:"index"`
	Scenario string `json:"scenario"`
	Mode     string `json:"mode"`
	Lighting string `json:"lighting"`
	// Missing: o cenário não está no jogo base nem nos mods instalados (ADR-030).
	Missing bool `json:"missing"`
	// Raw é a linha como está, quando o app não acha um cenário nela.
	Raw string `json:"raw"`
}

// CycleGroup é um bloco de linhas seguidas do mesmo mapa. Um mapa pode aparecer em mais de um bloco se a rotação
// voltar a ele depois de outros (a ordem da rotação é a da lista). Map vazio: linhas que o app não entende.
type CycleGroup struct {
	Map     string       `json:"map"`
	Label   string       `json:"label"`
	ModName string       `json:"modName"`
	Entries []CycleEntry `json:"entries"`
}

// GroupMapCycle divide a rotação em blocos por mapa, na ordem da lista. missing vem de MissingInCycle; catalog dá o
// nome do mapa e o mod.
func GroupMapCycle(lines []string, missing []bool, catalog []CycleMap) []CycleGroup {
	if len(lines) > MaxCycleLines {
		lines = lines[:MaxCycleLines]
	}
	modOf := map[string]string{}
	for _, cm := range catalog {
		for _, s := range cm.Scenarios {
			modOf[strings.ToLower(s.ID)] = cm.ModName
		}
	}
	out := []CycleGroup{}
	for i, l := range lines {
		id := LineScenario(l)
		e := CycleEntry{Index: i, Scenario: id, Lighting: lineLightingOrDay(l), Missing: i < len(missing) && missing[i]}
		name := ""
		if id == "" {
			e.Raw = strings.TrimSpace(l)
		} else {
			name, e.Mode = SplitScenario(id)
		}
		if n := len(out); n > 0 && strings.EqualFold(out[n-1].Map, name) {
			out[n-1].Entries = append(out[n-1].Entries, e)
			continue
		}
		g := CycleGroup{Map: name, Entries: []CycleEntry{e}}
		if name != "" {
			g.Label, g.ModName = mapLabel(name), modOf[strings.ToLower(id)]
		}
		out = append(out, g)
	}
	return out
}

// WithDayAndNight monta as linhas da rotação com cada cenário de dia e de noite: em cada mapa, primeiro os cenários
// de dia e depois os mesmos de noite (ADR-046). ids são cenários em ordem de mapa (DefaultMapCycle).
func WithDayAndNight(ids []string) []string {
	out := make([]string, 0, 2*len(ids))
	for i := 0; i < len(ids); {
		m, _ := SplitScenario(ids[i])
		j := i
		for j < len(ids) {
			if mj, _ := SplitScenario(ids[j]); !strings.EqualFold(mj, m) {
				break
			}
			j++
		}
		for _, light := range []string{LightingDay, LightingNight} {
			for _, id := range ids[i:j] {
				out = append(out, CycleLine(id, light))
			}
		}
		i = j
	}
	return out
}

// CountMaps diz quantos mapas diferentes os cenários têm.
func CountMaps(ids []string) int {
	seen := map[string]bool{}
	for _, id := range ids {
		m, _ := SplitScenario(id)
		seen[strings.ToLower(m)] = true
	}
	return len(seen)
}

// lineLightingOrDay é a iluminação da linha; sem ela, Day (a padrão dos cenários, ver FixMapCycleLine).
func lineLightingOrDay(line string) string {
	if l := LineLighting(line); l != "" {
		return l
	}
	return LightingDay
}

// cycleKey identifica uma linha da rotação: cenário e iluminação.
func cycleKey(scenario, lighting string) string {
	return strings.ToLower(scenario) + "|" + lighting
}

// groupBounds devolve o bloco (início e fim, inclusivos) da linha index.
func groupBounds(lines []string, index int) (int, int) {
	key := func(i int) string {
		id := LineScenario(lines[i])
		if id == "" {
			return ""
		}
		m, _ := SplitScenario(id)
		return strings.ToLower(m)
	}
	k := key(index)
	a, b := index, index
	for a > 0 && key(a-1) == k {
		a--
	}
	for b < len(lines)-1 && key(b+1) == k {
		b++
	}
	return a, b
}

// AddToMapCycle põe os cenários depois da linha after (-1: no fim), cada um em cada iluminação pedida (Day, Night
// ou as duas, nessa ordem), sem repetir um cenário que já está na rotação com a mesma iluminação. Recusa nomes que
// não são de cenário e iluminação que não seja Day ou Night.
func AddToMapCycle(lines []string, after int, ids, lightings []string) ([]string, error) {
	if len(lightings) == 0 {
		lightings = []string{LightingDay}
	}
	for _, l := range lightings {
		if l != LightingDay && l != LightingNight {
			return lines, problems.Field("mapCycle", "validation.lighting")
		}
	}
	have := map[string]bool{}
	for _, l := range lines {
		if id := LineScenario(l); id != "" {
			have[cycleKey(id, lineLightingOrDay(l))] = true
		}
	}
	var add []string
	for _, light := range lightings {
		for _, id := range ids {
			if !scenarioRe.MatchString(id) {
				return lines, problems.Field("mapCycle", "validation.scenario")
			}
			if k := cycleKey(id, light); !have[k] {
				have[k] = true
				add = append(add, CycleLine(id, light))
			}
		}
	}
	if len(lines)+len(add) > MaxCycleLines {
		return lines, problems.Field("mapCycle", "validation.mapCycleTooLong")
	}
	if after < 0 || after >= len(lines) {
		after = len(lines) - 1
	}
	out := make([]string, 0, len(lines)+len(add))
	out = append(out, lines[:after+1]...)
	out = append(out, add...)
	return append(out, lines[after+1:]...), nil
}

// RemoveFromMapCycle tira as linhas pedidas (posições fora da lista são ignoradas).
func RemoveFromMapCycle(lines []string, indexes []int) []string {
	drop := map[int]bool{}
	for _, i := range indexes {
		drop[i] = true
	}
	out := make([]string, 0, len(lines))
	for i, l := range lines {
		if !drop[i] {
			out = append(out, l)
		}
	}
	return out
}

// MoveMapCycleGroup troca o bloco da linha index com o bloco de cima (up) ou de baixo. As linhas de dentro dos
// blocos ficam na mesma ordem.
func MoveMapCycleGroup(lines []string, index int, up bool) []string {
	if index < 0 || index >= len(lines) {
		return lines
	}
	a, b := groupBounds(lines, index)
	var c, d int // o outro bloco
	if up {
		if a == 0 {
			return lines
		}
		c, _ = groupBounds(lines, a-1)
		d = a - 1
		a, b, c, d = c, d, a, b
	} else {
		if b == len(lines)-1 {
			return lines
		}
		c = b + 1
		_, d = groupBounds(lines, c)
	}
	// [a..b] e [c..d] são seguidos: [a..b][c..d] → [c..d][a..b]
	out := append([]string{}, lines[:a]...)
	out = append(out, lines[c:d+1]...)
	out = append(out, lines[a:b+1]...)
	return append(out, lines[d+1:]...)
}
