package logbuf

import (
	"strings"
	"testing"
)

func TestMaskRingAndFlush(t *testing.T) {
	b := New(3, func() []string { return []string{"minhaSenha#1"} })
	var got []Line
	b.SetNotify(func(l []Line) { got = append(got, l...) })

	b.Add("server", "LogInit: Command Line: x?Password=minhaSenha#1 -GSLTToken=ABCDEF0123456789", "")
	b.Add("app", "texto com minhaSenha#1 solto", "info")
	b.Add("server", "[x]LogWindows: Error: Fatal error", "")
	b.Add("server", "LogModioGame: ModSubsystem: ModManager enabled", "")

	recent := b.Recent(10)
	if len(recent) != 3 {
		t.Fatalf("anel: %d linhas", len(recent))
	}
	for _, l := range recent {
		if strings.Contains(l.Text, "minhaSenha") || strings.Contains(l.Text, "ABCDEF0123") {
			t.Errorf("segredo visível: %q", l.Text)
		}
	}
	if recent[1].Level != "error" || !recent[2].Event {
		t.Errorf("%+v", recent)
	}
	b.Flush()
	if len(got) != 4 {
		t.Errorf("flush enviou %d", len(got))
	}
	b.Flush()
	if len(got) != 4 {
		t.Error("flush repetiu linhas")
	}
}
