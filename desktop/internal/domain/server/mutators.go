package server

import (
	"regexp"
	"strings"

	"sandstorm-server-assist/desktop/internal/domain/problems"
)

// Mutator é um mutator oficial do jogo: o id vai na linha de comando (-mutators=); Label é o nome no jogo. A
// descrição da tela vem das traduções (mutators.<id>). Changes são os ajustes do jogo que ele troca; dois mutators
// que trocam o mesmo ajuste não somam: um sobrescreve o outro (ADR-045). Conflicts são os mutators da lista que
// trocam algum ajuste em comum com este.
type Mutator struct {
	ID        string   `json:"id"`
	Label     string   `json:"label"`
	Changes   []string `json:"changes"`
	Conflicts []string `json:"conflicts"`
}

// Ajustes que os mutators do jogo trocam, conferidos nos arquivos do jogo (Insurgency/Content/Game/Mutators no
// pakchunk0): pontos de abastecimento iniciais (InitialSupply), classes e armas (troca do SquadLayout), velocidade
// dos jogadores (SpeedMultiplier) e tempo de captura (ObjectiveCaptureTime).
const (
	ChangeSupply      = "supply"
	ChangeLoadout     = "loadout"
	ChangeSpeed       = "speed"
	ChangeCaptureTime = "captureTime"
)

// Mutators são os mutators que vêm com o jogo (Server Admin Guide). Mutators de mods continuam aceitos pelo nome.
var Mutators = withConflicts([]Mutator{
	{ID: "AllYouCanEat", Label: "All You Can Eat", Changes: []string{ChangeSupply}},
	{ID: "AntiMaterielRiflesOnly", Label: "Anti-Materiel Only", Changes: []string{ChangeLoadout}},
	{ID: "BoltActionsOnly", Label: "Bolt-Actions Only", Changes: []string{ChangeLoadout}},
	{ID: "Broke", Label: "Broke", Changes: []string{ChangeSupply}},
	{ID: "BulletSponge", Label: "Bullet Sponge"},
	{ID: "Competitive", Label: "Competitive", Changes: []string{ChangeCaptureTime}},
	{ID: "CompetitiveLoadouts", Label: "Competitive Loadouts", Changes: []string{ChangeLoadout}},
	{ID: "FastMovement", Label: "Fast Movement", Changes: []string{ChangeSpeed}},
	{ID: "Frenzy", Label: "Frenzy"},
	{ID: "Guerrillas", Label: "Guerrillas", Changes: []string{ChangeSupply}},
	{ID: "Hardcore", Label: "Hardcore", Changes: []string{ChangeSpeed, ChangeCaptureTime}},
	{ID: "HeadshotOnly", Label: "Headshots Only"},
	{ID: "HotPotato", Label: "Hot Potato"},
	{ID: "LockedAim", Label: "Locked Aim"},
	{ID: "NoAim", Label: "No Aim Down Sights"},
	{ID: "PistolsOnly", Label: "Pistols Only", Changes: []string{ChangeLoadout}},
	{ID: "ShotgunsOnly", Label: "Shotguns Only", Changes: []string{ChangeLoadout}},
	{ID: "SlowCaptureTimes", Label: "Slow Capture Times", Changes: []string{ChangeCaptureTime}},
	{ID: "SlowMovement", Label: "Slow Movement", Changes: []string{ChangeSpeed}},
	{ID: "SoldierOfFortune", Label: "Soldier of Fortune"},
	{ID: "SpecialOperations", Label: "Special Operations", Changes: []string{ChangeSupply}},
	{ID: "Strapped", Label: "Strapped", Changes: []string{ChangeSupply}},
	{ID: "Ultralethal", Label: "Ultralethal"},
	{ID: "Vampirism", Label: "Vampirism"},
	{ID: "Warlords", Label: "Warlords", Changes: []string{ChangeSupply}},
})

// withConflicts preenche Conflicts: os outros mutators que trocam algum ajuste em comum.
func withConflicts(list []Mutator) []Mutator {
	for i := range list {
		list[i].Conflicts = []string{}
		if list[i].Changes == nil {
			list[i].Changes = []string{}
		}
		for j := range list {
			if i != j && shareChange(list[i], list[j]) {
				list[i].Conflicts = append(list[i].Conflicts, list[j].ID)
			}
		}
	}
	return list
}

func shareChange(a, b Mutator) bool {
	for _, x := range a.Changes {
		for _, y := range b.Changes {
			if x == y {
				return true
			}
		}
	}
	return false
}

func officialMutator(id string) (Mutator, bool) {
	for _, m := range Mutators {
		if m.ID == id {
			return m, true
		}
	}
	return Mutator{}, false
}

// CheckMutatorConflicts recusa dois mutators do jogo que trocam o mesmo ajuste (um sobrescreveria o outro).
// Mutators de mods não entram: o app não sabe o que eles trocam.
func CheckMutatorConflicts(ids []string) error {
	for i, a := range ids {
		ma, ok := officialMutator(a)
		if !ok {
			continue
		}
		for _, b := range ids[i+1:] {
			if mb, ok := officialMutator(b); ok && a != b && shareChange(ma, mb) {
				return problems.Field("mutators", "validation.mutatorConflict", "a", ma.Label, "b", mb.Label)
			}
		}
	}
	return nil
}

// IsOfficialMutator diz se o nome é de um mutator do jogo.
func IsOfficialMutator(id string) bool {
	_, ok := officialMutator(id)
	return ok
}

// Iluminação do mapa (parâmetro ?Lighting= da URL de viagem e Lighting="..." na rotação).
const (
	LightingDay   = "Day"
	LightingNight = "Night"
)

// ValidLighting aceita vazio (padrão do cenário), Day e Night.
func ValidLighting(v string) bool { return v == "" || v == LightingDay || v == LightingNight }

func validateLighting(field, v string) error {
	if !ValidLighting(v) {
		return problems.Field(field, "validation.lighting")
	}
	return nil
}

var (
	reLightingAttr = regexp.MustCompile(`(?i),\s*Lighting\s*=\s*"?(Day|Night)"?`)
	reScenarioAttr = regexp.MustCompile(`(?i)Scenario\s*=\s*"?(Scenario_[A-Za-z0-9_]+)"?`)
)

// LineLighting devolve a iluminação de uma linha da rotação ("" se não tem).
func LineLighting(line string) string {
	m := reLightingAttr.FindStringSubmatch(line)
	if m == nil {
		return ""
	}
	if strings.EqualFold(m[1], LightingNight) {
		return LightingNight
	}
	return LightingDay
}

// LineScenario devolve o cenário de uma linha da rotação (simples ou entre parênteses).
func LineScenario(line string) string {
	line = strings.TrimSpace(line)
	if m := reScenarioAttr.FindStringSubmatch(line); m != nil {
		return m[1]
	}
	if scenarioRe.MatchString(line) {
		return line
	}
	return ""
}

// SetLineLighting troca a iluminação de uma linha da rotação. Linha simples vira o formato entre parênteses;
// vazio tira o Lighting. Linhas que o app não entende voltam iguais.
func SetLineLighting(line, lighting string) string {
	line = FixMapCycleLine(line)
	if !ValidLighting(lighting) {
		return line
	}
	if scenarioRe.MatchString(line) {
		if lighting == "" {
			return line
		}
		return `(Scenario="` + line + `",Lighting="` + lighting + `")`
	}
	if !strings.HasPrefix(line, "(") || !strings.HasSuffix(line, ")") || LineScenario(line) == "" {
		return line
	}
	inner := reLightingAttr.ReplaceAllString(line[1:len(line)-1], "")
	if lighting != "" {
		inner += `,Lighting="` + lighting + `"`
	}
	return "(" + inner + ")"
}
