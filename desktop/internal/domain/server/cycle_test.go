package server

import (
	"reflect"
	"testing"
)

var cycleLines = []string{
	"Scenario_Farmhouse_Checkpoint_Security",
	`(Scenario="Scenario_Farmhouse_Checkpoint_Insurgents",Lighting="Night")`,
	"Scenario_Outback_Checkpoint_Security",
	"something odd",
	"Scenario_Farmhouse_Push_Security",
}

func TestGroupMapCycleKeepsTheOrder(t *testing.T) {
	cat := CycleCatalog([]string{"Scenario_Farmhouse_Checkpoint_Security"},
		[]CycleSource{{Scenario: "Scenario_Outback_Checkpoint_Security", ModID: 7, ModName: "Outback"}})
	g := GroupMapCycle(cycleLines, []bool{false, false, false, false, true}, cat)
	if len(g) != 4 {
		t.Fatalf("groups %+v", g)
	}
	if g[0].Map != "Farmhouse" || len(g[0].Entries) != 2 || g[0].Entries[1].Lighting != LightingNight ||
		g[0].Entries[0].Lighting != LightingDay ||
		g[0].Entries[1].Mode != "Checkpoint_Insurgents" {
		t.Errorf("first %+v", g[0])
	}
	if g[1].ModName != "Outback" || g[1].Label != "Outback" {
		t.Errorf("mod group %+v", g[1])
	}
	if g[2].Map != "" || g[2].Entries[0].Raw != "something odd" {
		t.Errorf("odd line %+v", g[2])
	}
	if g[3].Map != "Farmhouse" || !g[3].Entries[0].Missing || g[3].Entries[0].Index != 4 {
		t.Errorf("farmhouse again %+v", g[3])
	}
}

func TestCycleCatalog(t *testing.T) {
	cat := CycleCatalog(
		[]string{"Scenario_Farmhouse_Push_Security", "Scenario_Farmhouse_Checkpoint_Insurgents",
			"Scenario_Farmhouse_Checkpoint_Security", "not a scenario", "Scenario_DEV_TestingGamemode_AllLevels"},
		[]CycleSource{{Scenario: "Scenario_Outback_Checkpoint_Security", ModID: 7, ModName: "Outback"},
			{Scenario: "Scenario_Farmhouse_Checkpoint_Security", ModID: 9, ModName: "Copy"}})
	if len(cat) != 2 || cat[0].ModID != 7 || cat[1].Label != "Farmhouse" {
		t.Fatalf("catalog %+v", cat)
	}
	got := []string{}
	for _, s := range cat[1].Scenarios {
		got = append(got, s.Mode)
	}
	if !reflect.DeepEqual(got, []string{"Checkpoint_Security", "Checkpoint_Insurgents", "Push_Security"}) {
		t.Errorf("modes %v", got)
	}
	if fallback := CycleCatalog(nil, nil); len(fallback) != len(Maps) || len(fallback[0].Scenarios) != len(Modes) {
		t.Errorf("fallback %d", len(fallback))
	}
}

func TestAddToMapCycle(t *testing.T) {
	// Farmhouse Checkpoint Security já está (linha sem iluminação = Day): só entra de noite. Push entra nas duas.
	out, err := AddToMapCycle(cycleLines, 1, []string{"Scenario_Farmhouse_Push_Insurgents",
		"Scenario_Farmhouse_Checkpoint_Security", "Scenario_Farmhouse_Push_Insurgents"},
		[]string{LightingDay, LightingNight})
	if err != nil {
		t.Fatal(err)
	}
	added := []string{
		`(Scenario="Scenario_Farmhouse_Push_Insurgents",Lighting="Day")`,
		`(Scenario="Scenario_Farmhouse_Push_Insurgents",Lighting="Night")`,
		`(Scenario="Scenario_Farmhouse_Checkpoint_Security",Lighting="Night")`,
	}
	want := append(append(append([]string{}, cycleLines[:2]...), added...), cycleLines[2:]...)
	if !reflect.DeepEqual(out, want) {
		t.Fatalf("got %q", out)
	}
	// Farmhouse Checkpoint Insurgents já está de noite: pedir de noite não repete.
	if out, _ := AddToMapCycle(cycleLines, -1, []string{"Scenario_Farmhouse_Checkpoint_Insurgents"},
		[]string{LightingNight}); len(out) != len(cycleLines) {
		t.Errorf("repeated %q", out)
	}
	if out, _ := AddToMapCycle(nil, -1, []string{"Scenario_Alps_Checkpoint_Security"}, nil); len(out) != 1 ||
		out[0] != `(Scenario="Scenario_Alps_Checkpoint_Security",Lighting="Day")` {
		t.Errorf("empty cycle %q", out)
	}
	if _, err := AddToMapCycle(cycleLines, -1, []string{"Scenario_X?game=evil"}, nil); err == nil {
		t.Error("accepted a free text")
	}
	if _, err := AddToMapCycle(cycleLines, -1, []string{"Scenario_Alps_Checkpoint_Security"}, []string{"Dusk"}); err == nil {
		t.Error("accepted another lighting")
	}
	big := make([]string, MaxCycleLines)
	for i := range big {
		big[i] = "x"
	}
	if _, err := AddToMapCycle(big, -1, []string{"Scenario_Alps_Checkpoint_Security"}, nil); err == nil {
		t.Error("accepted a cycle over the limit")
	}
}

func TestWithDayAndNight(t *testing.T) {
	ids := []string{"Scenario_Bab_Checkpoint_Security", "Scenario_Bab_Checkpoint_Insurgents", "Scenario_Hold_Push_Security"}
	want := []string{
		`(Scenario="Scenario_Bab_Checkpoint_Security",Lighting="Day")`,
		`(Scenario="Scenario_Bab_Checkpoint_Insurgents",Lighting="Day")`,
		`(Scenario="Scenario_Bab_Checkpoint_Security",Lighting="Night")`,
		`(Scenario="Scenario_Bab_Checkpoint_Insurgents",Lighting="Night")`,
		`(Scenario="Scenario_Hold_Push_Security",Lighting="Day")`,
		`(Scenario="Scenario_Hold_Push_Security",Lighting="Night")`,
	}
	if got := WithDayAndNight(ids); !reflect.DeepEqual(got, want) {
		t.Fatalf("got %q", got)
	}
	if n := CountMaps(ids); n != 2 {
		t.Errorf("maps %d", n)
	}
	if g := GroupMapCycle(want, nil, nil); len(g) != 2 || len(g[0].Entries) != 4 {
		t.Errorf("groups %+v", g)
	}
}

func TestRemoveAndMoveMapCycle(t *testing.T) {
	if out := RemoveFromMapCycle(cycleLines, []int{0, 1, 99, -1}); len(out) != 3 || out[0] != cycleLines[2] {
		t.Errorf("remove %q", out)
	}
	down := MoveMapCycleGroup(cycleLines, 1, false)
	if !reflect.DeepEqual(down, []string{cycleLines[2], cycleLines[0], cycleLines[1], cycleLines[3], cycleLines[4]}) {
		t.Errorf("down %q", down)
	}
	up := MoveMapCycleGroup(cycleLines, 2, true)
	if !reflect.DeepEqual(up, down) {
		t.Errorf("up %q", up)
	}
	if out := MoveMapCycleGroup(cycleLines, 0, true); !reflect.DeepEqual(out, cycleLines) {
		t.Errorf("first up %q", out)
	}
	if out := MoveMapCycleGroup(cycleLines, 4, false); !reflect.DeepEqual(out, cycleLines) {
		t.Errorf("last down %q", out)
	}
}
