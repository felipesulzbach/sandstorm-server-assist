package serverlogs

import (
	"testing"
	"time"
)

func TestIsServerLog(t *testing.T) {
	for _, ok := range []string{"Insurgency.log", "Insurgency_2.log", "Insurgency-backup-2026.11.09-22.56.07.log"} {
		if !IsServerLog(ok) {
			t.Errorf("should accept %s", ok)
		}
	}
	for _, bad := range []string{"Game.ini", "Insurgency.log.bak", "../Insurgency.log", "Insurgency-backup-x.log", "other.log", "Insurgency_2.log/x"} {
		if IsServerLog(bad) {
			t.Errorf("should refuse %s", bad)
		}
	}
}

func TestToDeleteKeepsNewestAndCurrent(t *testing.T) {
	base := time.Date(2026, 11, 9, 0, 0, 0, 0, time.UTC)
	files := []File{
		{Name: "Insurgency.log", ModTime: base.Add(10 * time.Hour)},
		{Name: "Insurgency-backup-2026.11.09-01.00.00.log", ModTime: base.Add(1 * time.Hour)},
		{Name: "Insurgency-backup-2026.11.09-05.00.00.log", ModTime: base.Add(5 * time.Hour)},
		{Name: "Insurgency-backup-2026.11.09-03.00.00.log", ModTime: base.Add(3 * time.Hour)},
		{Name: "Insurgency_2.log", ModTime: base.Add(2 * time.Hour)},
		{Name: "notes.txt", ModTime: base},
	}
	got := ToDelete(files, 3)
	// Mantém a atual + as 2 mais novas (05h e 03h); apaga 02h e 01h.
	if len(got) != 2 || got[0].Name != "Insurgency_2.log" || got[1].Name != "Insurgency-backup-2026.11.09-01.00.00.log" {
		t.Fatalf("got %+v", got)
	}
	// Abaixo do mínimo vira o mínimo: atual + 1.
	if n := len(ToDelete(files, 0)); n != 3 {
		t.Fatalf("keep 0 deleted %d", n)
	}
	// Nunca apaga o log atual, mesmo sem outros.
	if len(ToDelete([]File{{Name: Current}}, MinKeep)) != 0 {
		t.Fatal("deleted the current log")
	}
	if !ValidKeep(MinKeep) || !ValidKeep(MaxKeep) || ValidKeep(1) || ValidKeep(MaxKeep+1) {
		t.Fatal("ValidKeep limits")
	}
}
