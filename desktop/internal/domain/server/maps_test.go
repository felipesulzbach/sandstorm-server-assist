package server

import (
	"strings"
	"testing"
)

func TestDefaultMapCycle(t *testing.T) {
	// Sem os .pak: a lista Maps, Checkpoint dos dois lados.
	got := DefaultMapCycle("Scenario_Farmhouse_Checkpoint_Security", nil)
	if len(got) != 2*len(Maps) || got[0] != "Scenario_Bab_Checkpoint_Security" || got[1] != "Scenario_Bab_Checkpoint_Insurgents" {
		t.Fatalf("fallback %v", got[:4])
	}
	// Com os .pak: só o modo do mapa inicial, sem repetir, e mapas que a lista não tem.
	base := []string{
		"Scenario_Trainyard_Push_Security", "Scenario_Trainyard_Checkpoint_Insurgents",
		"Scenario_Hold_Checkpoint_Security", "Scenario_Hold_Checkpoint_Security", "Scenario_Farmhouse_Firefight_East",
	}
	got = DefaultMapCycle("Scenario_Farmhouse_Checkpoint_Security", base)
	if strings.Join(got, ",") != "Scenario_Hold_Checkpoint_Security,Scenario_Trainyard_Checkpoint_Insurgents" {
		t.Fatalf("from paks %v", got)
	}
	got = DefaultMapCycle("Scenario_Farmhouse_Push_Security", base)
	if strings.Join(got, ",") != "Scenario_Trainyard_Push_Security" {
		t.Fatalf("push %v", got)
	}
}
