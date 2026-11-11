package modcheck

import (
	"strings"
	"testing"
)

const crash = `Fatal error: [File:Unknown] [Line: 3237]
StaticMesh /Game/Environment/Props/Modular/Stairs/SM_Railing_Wooden_04a_200x70x20.SM_Railing_Wooden_04a_200x70x20: Serial size mismatch: Got 92331, Expected 92328

0x00007ffed263569a KERNELBASE.dll!UnknownFunction []
Crash in runnable thread FAsyncLoadingThread`

func mods() []Mod {
	return []Mod{
		{ID: 5748299, Name: "Daybreak", Files: []File{
			{"Insurgency/Content/Environment/Props/Modular/Stairs/SM_Railing_Wooden_04a_200x70x20.uasset", 1200},
			{"Insurgency/Content/Environment/Props/Modular/Stairs/SM_Railing_Wooden_04a_200x70x20.uexp", 92331},
			{"Insurgency/Content/Shared/T_Rock.uasset", 500},
			{"Insurgency/Content/Shared/T_Rock.uexp", 9000},
		}},
		{ID: 2123731, Name: "Containment", Files: []File{
			{"Insurgency/Content/Shared/T_Rock.uasset", 500},
			{"Insurgency/Content/Shared/T_Rock.uexp", 9000}, // cópia idêntica: não conta
			{"Insurgency/Content/Subway/M_Wall.uasset", 100},
			{"Insurgency/Content/Subway/M_Wall.uexp", 4000},
		}},
		{ID: 2853082, Name: "Containment Storm", Files: []File{
			{"Insurgency/Content/Subway/M_Wall.uasset", 100},
			{"Insurgency/Content/Subway/M_Wall.uexp", 4100}, // versão diferente
			{"insurgency/content/subway/M_Floor.uexp", 1},
		}},
		{ID: 98145, Name: "Mutator", Files: []File{
			{"Insurgency/Plugins/NoRestrictedArea/Content/Mutators/NoRestrictedArea.uasset", 300},
		}},
	}
}

func TestCheckFindsTheModOfACrash(t *testing.T) {
	findings, mismatch := Check(crash, mods())
	if !mismatch || len(findings) != 1 {
		t.Fatalf("mismatch=%v findings=%+v", mismatch, findings)
	}
	f := findings[0]
	if f.Object != "/Game/Environment/Props/Modular/Stairs/SM_Railing_Wooden_04a_200x70x20" ||
		len(f.Mods) != 1 || f.Mods[0].Name != "Daybreak" {
		t.Fatalf("finding %+v", f)
	}
}

func TestOwnersOfPluginAndUnknown(t *testing.T) {
	if o := Owners("/NoRestrictedArea/Mutators/NoRestrictedArea", mods()); len(o) != 1 || o[0].ID != 98145 {
		t.Fatalf("plugin owners %+v", o)
	}
	findings, mismatch := Check("Serial size mismatch on /Game/Not/There.There", mods())
	if !mismatch || len(findings) != 1 || len(findings[0].Mods) != 0 {
		t.Fatalf("findings %+v", findings)
	}
	if f, _ := Check("nothing here", mods()); f == nil || len(f) != 0 {
		t.Fatalf("no objects should be an empty list, got %+v", f)
	}
}

func TestAssetsIgnoresScriptAndRepeats(t *testing.T) {
	got := Assets("/Script/Engine.StaticMesh /Game/A/B.B and again /game/a/b.B (/Engine/X/Y.Y)")
	if strings.Join(got, ",") != "/Game/A/B,/Engine/X/Y" {
		t.Fatalf("assets %v", got)
	}
	long := strings.Repeat("/Game/A/B", 100000)
	_ = Assets(long) // limitado, sem travar
}

func TestConflicts(t *testing.T) {
	c := Conflicts(mods())
	if len(c) != 1 {
		t.Fatalf("conflicts %+v", c)
	}
	if c[0].A.ID != 2123731 || c[0].B.ID != 2853082 || c[0].Count != 1 || c[0].Example != "/Game/Subway/M_Wall" {
		t.Fatalf("conflict %+v", c[0])
	}
	if c := Conflicts(nil); len(c) != 0 {
		t.Fatal("no mods, no conflicts")
	}
}

func TestGamePath(t *testing.T) {
	if GamePath("Insurgency/Content/A/B") != "/Game/A/B" || GamePath("Engine/Content/X") != "/Engine/X" ||
		GamePath("Insurgency/Plugins/P/Content/Z") != "/Insurgency/Plugins/P/Content/Z" {
		t.Fatal("game path")
	}
}
