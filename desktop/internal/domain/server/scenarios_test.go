package server

import (
	"reflect"
	"testing"

	"sandstorm-server-assist/desktop/internal/domain/problems"
)

func TestParseScenarios(t *testing.T) {
	resp := "Scenarios:\nScenario_Farmhouse_Checkpoint_Security\nScenario_Containment_Checkpoint_Security (Containment)\n" +
		"Scenario_Farmhouse_Checkpoint_Security\n  Scenario_Gap_Push_Insurgents\nnot a scenario\nScenario_\n"
	want := []string{
		"Scenario_Containment_Checkpoint_Security",
		"Scenario_Farmhouse_Checkpoint_Security",
		"Scenario_Gap_Push_Insurgents",
	}
	if got := ParseScenarios(resp); !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v", got)
	}
	if got := ParseScenarios(""); got != nil {
		t.Fatalf("empty: %v", got)
	}
}

func TestSplitScenario(t *testing.T) {
	for _, c := range []struct{ id, m, mode string }{
		{"Scenario_Farmhouse_Checkpoint_Security", "Farmhouse", "Checkpoint_Security"},
		{"Scenario_LastLight_Push_Insurgents", "LastLight", "Push_Insurgents"},
		{"Scenario_Old_World_Checkpoint_Security", "Old_World", "Checkpoint_Security"},
		{"Scenario_Tell_Firefight_East", "Tell", "Firefight_East"},
		{"Scenario_Crossing_Team_Deathmatch", "Crossing", "Team_Deathmatch"},
		{"Scenario_Prison_FFA", "Prison", "FFA"},
		{"Scenario_RedwoodReborn_Checkpoint_2", "RedwoodReborn", "Checkpoint_2"},
		{"Scenario_Custom_Thing", "Custom_Thing", ""},
	} {
		m, mode := SplitScenario(c.id)
		if m != c.m || mode != c.mode {
			t.Fatalf("%s → %q %q", c.id, m, mode)
		}
	}
}

func TestMapChoices(t *testing.T) {
	p := Default()
	p.TravelMap, p.TravelScenario = "Containment", "Scenario_Containment_Checkpoint_Security"
	p.MapCycle = []string{
		`(Scenario="Scenario=Scenario_Outback_Checkpoint_Security",Lighting="Day")`,
		"Scenario_Farmhouse_Checkpoint_Security",
	}
	got := MapChoices([]string{"Scenario_Farmhouse_Checkpoint_Security", "Scenario_Gap_Push_Insurgents"}, p, nil)
	want := []Scenario{
		{"Scenario_Containment_Checkpoint_Security", "Containment", "Checkpoint_Security", FromConfig},
		{"Scenario_Farmhouse_Checkpoint_Security", "Farmhouse", "Checkpoint_Security", FromServer},
		{"Scenario_Gap_Push_Insurgents", "Gap", "Push_Insurgents", FromServer},
		{"Scenario_Outback_Checkpoint_Security", "Outback", "Checkpoint_Security", FromConfig},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v", got)
	}
}

// Com os mods carregados conhecidos, a lista traz os mapas deles e deixa de fora os da rotação sem mod (11/11/2026:
// um mapa de mod com a inscrição retirada continuava na rotação e aparecia na lista).
func TestMapChoicesUseLoadedMods(t *testing.T) {
	p := Default()
	p.TravelMap, p.TravelScenario = "Containment", "Scenario_Containment_Checkpoint_Security"
	p.MapCycle = []string{"Scenario_Redwood_Checkpoint_Security", "Scenario_Farmhouse_Checkpoint_Security"}
	got := MapChoices([]string{"Scenario_Farmhouse_Checkpoint_Security"}, p,
		[]string{"Scenario_Containment_Checkpoint_Security", "Scenario_Estates_Checkpoint_Insurgency", "bad id; quit"})
	want := []Scenario{
		{"Scenario_Containment_Checkpoint_Security", "Containment", "Checkpoint_Security", FromMod},
		{"Scenario_Estates_Checkpoint_Insurgency", "Estates", "Checkpoint_Insurgency", FromMod},
		{"Scenario_Farmhouse_Checkpoint_Security", "Farmhouse", "Checkpoint_Security", FromServer},
		{"Scenario_Redwood_Checkpoint_Security", "Redwood", "Checkpoint_Security", FromMissing},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v", got)
	}
	// O mapa sem mod aparece, mas a troca para ele é recusada antes de ir ao servidor.
	_, err := TravelCommand("Scenario_Redwood_Checkpoint_Security", got)
	if pe, ok := err.(*problems.Error); !ok || pe.Msg != "validation.scenarioNotLoaded" {
		t.Fatalf("troca para mapa sem mod: %v", err)
	}
	// Nenhum mod carregado (lista vazia, não nil): jogo base e o que sobrou do perfil, marcado.
	got = MapChoices([]string{"Scenario_Farmhouse_Checkpoint_Security"}, p, []string{})
	if len(got) != 3 || got[0].Source != FromMissing || got[2].Source != FromMissing {
		t.Fatalf("sem mods: %+v", got)
	}
}

// Linhas reais da rotação (11/11/2026): mapas de mods que não estão instalados ficam marcados.
func TestMissingInCycle(t *testing.T) {
	lines := []string{
		`(Scenario="Scenario_Farmhouse_Checkpoint_Security",Lighting="Day")`,
		`(Scenario="Scenario_Outback_Checkpoint_Security",Lighting="Night")`,
		`(Scenario="Scenario_RedwoodReborn_Checkpoint_2",Lighting="Winter")`,
		"Scenario_Estates_Checkpoint_Insurgency",
		"linha que o app não entende",
	}
	known := []string{"Scenario_Farmhouse_Checkpoint_Security", "scenario_estates_checkpoint_insurgency"}
	if got := MissingInCycle(lines, known); !reflect.DeepEqual(got, []bool{false, true, true, false, false}) {
		t.Fatalf("got %v", got)
	}
	// Sem o jogo base lido, nada é marcado.
	if got := MissingInCycle(lines, nil); !reflect.DeepEqual(got, make([]bool, len(lines))) {
		t.Fatalf("sem lista: %v", got)
	}
}

func TestTravelCommand(t *testing.T) {
	list := MapChoices([]string{"Scenario_Gap_Push_Insurgents"}, Default(), nil)
	cmd, err := TravelCommand("Scenario_Gap_Push_Insurgents", list)
	if err != nil || cmd != "travelscenario Scenario_Gap_Push_Insurgents" {
		t.Fatalf("%q %v", cmd, err)
	}
	for _, c := range []struct{ id, msg string }{
		{"Scenario_Gap_Push_Insurgents; quit", "validation.scenario"},
		{"Gap?Password=x", "validation.scenario"},
		{"Scenario_Tell_Checkpoint_Security", "validation.scenarioUnknown"},
	} {
		_, err := TravelCommand(c.id, list)
		pe, ok := err.(*problems.Error)
		if !ok || pe.Msg != c.msg {
			t.Fatalf("%q: %v", c.id, err)
		}
	}
}

func TestTravelRefused(t *testing.T) {
	for _, r := range []string{"Unknown scenario", "Invalid scenario name", "Scenario not found", "Usage: travelscenario <scenario>"} {
		if !TravelRefused(r) {
			t.Fatalf("not refused: %q", r)
		}
	}
	for _, r := range []string{"", "Travelling to Scenario_Gap_Push_Insurgents", "OK"} {
		if TravelRefused(r) {
			t.Fatalf("refused: %q", r)
		}
	}
}
