package server

import (
	"slices"
	"testing"
)

func TestMutatorConflictsAreSymmetric(t *testing.T) {
	byID := map[string]Mutator{}
	for _, m := range Mutators {
		byID[m.ID] = m
	}
	for _, m := range Mutators {
		for _, c := range m.Conflicts {
			if c == m.ID || !slices.Contains(byID[c].Conflicts, m.ID) {
				t.Errorf("%s ↔ %s", m.ID, c)
			}
		}
	}
	supply := byID["AllYouCanEat"].Conflicts
	for _, id := range []string{"Broke", "Guerrillas", "SpecialOperations", "Strapped", "Warlords"} {
		if !slices.Contains(supply, id) {
			t.Errorf("AllYouCanEat should conflict with %s", id)
		}
	}
	if slices.Contains(supply, "SoldierOfFortune") {
		t.Error("Soldier of Fortune does not set the starting supply")
	}
	hard := byID["Hardcore"].Conflicts
	for _, id := range []string{"FastMovement", "SlowMovement", "SlowCaptureTimes", "Competitive"} {
		if !slices.Contains(hard, id) {
			t.Errorf("Hardcore should conflict with %s", id)
		}
	}
	if len(byID["BulletSponge"].Conflicts) != 0 || byID["BulletSponge"].Changes == nil {
		t.Errorf("BulletSponge %+v", byID["BulletSponge"])
	}
}

func TestCheckMutatorConflicts(t *testing.T) {
	for _, ok := range [][]string{
		nil,
		{"AllYouCanEat", "Hardcore", "PistolsOnly"},
		{"Strapped", "SoldierOfFortune", "MyModMutator", "OtherModMutator"},
		{"AllYouCanEat", "AllYouCanEat"},
	} {
		if err := CheckMutatorConflicts(ok); err != nil {
			t.Errorf("%v: %v", ok, err)
		}
	}
	for _, bad := range [][]string{
		{"AllYouCanEat", "Strapped"},
		{"PistolsOnly", "CompetitiveLoadouts"},
		{"Hardcore", "FastMovement"},
		{"SlowCaptureTimes", "Competitive"},
	} {
		if err := CheckMutatorConflicts(bad); err == nil {
			t.Errorf("%v accepted", bad)
		}
	}
}
