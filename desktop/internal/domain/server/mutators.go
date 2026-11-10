package server

import (
	"regexp"
	"strings"

	"sandstorm-server-assist/desktop/internal/domain/problems"
)

// Mutator é um mutator oficial do jogo: o id vai na linha de comando (-mutators=); Label é o nome no jogo. A
// descrição da tela vem das traduções (mutators.<id>).
type Mutator struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

// Mutators são os mutators que vêm com o jogo (Server Admin Guide). Mutators de mods continuam aceitos pelo nome.
var Mutators = []Mutator{
	{"AllYouCanEat", "All You Can Eat"},
	{"AntiMaterielRiflesOnly", "Anti-Materiel Only"},
	{"BoltActionsOnly", "Bolt-Actions Only"},
	{"Broke", "Broke"},
	{"BulletSponge", "Bullet Sponge"},
	{"Competitive", "Competitive"},
	{"CompetitiveLoadouts", "Competitive Loadouts"},
	{"FastMovement", "Fast Movement"},
	{"Frenzy", "Frenzy"},
	{"Guerrillas", "Guerrillas"},
	{"Hardcore", "Hardcore"},
	{"HeadshotOnly", "Headshots Only"},
	{"HotPotato", "Hot Potato"},
	{"LockedAim", "Locked Aim"},
	{"NoAim", "No Aim Down Sights"},
	{"PistolsOnly", "Pistols Only"},
	{"ShotgunsOnly", "Shotguns Only"},
	{"SlowCaptureTimes", "Slow Capture Times"},
	{"SlowMovement", "Slow Movement"},
	{"SoldierOfFortune", "Soldier of Fortune"},
	{"SpecialOperations", "Special Operations"},
	{"Strapped", "Strapped"},
	{"Ultralethal", "Ultralethal"},
	{"Vampirism", "Vampirism"},
	{"Warlords", "Warlords"},
}

// IsOfficialMutator diz se o nome é de um mutator do jogo.
func IsOfficialMutator(id string) bool {
	for _, m := range Mutators {
		if m.ID == id {
			return true
		}
	}
	return false
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
