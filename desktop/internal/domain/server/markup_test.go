package server

import (
	"testing"

	"sandstorm-server-assist/desktop/internal/domain/problems"
)

func TestParseMarkup(t *testing.T) {
	text := "<Header1>Welcome</>\nPlay <Bold>fair</> and have fun.\n<extimg src=\"https://example.com/logo.png\"/>"
	segs, err := ParseMarkup(text, false, "motd")
	if err != nil {
		t.Fatal(err)
	}
	want := []Segment{
		{Style: "Header1", Text: "Welcome"}, {Text: "\nPlay "}, {Style: "Bold", Text: "fair"},
		{Text: " and have fun.\n"}, {Src: "https://example.com/logo.png"},
	}
	if len(segs) != len(want) {
		t.Fatalf("%+v", segs)
	}
	for i := range want {
		if segs[i] != want[i] {
			t.Fatalf("%d: %+v", i, segs[i])
		}
	}
	// Texto comum com < e > não é tag.
	if segs, err := ParseMarkup("a < b and c > d", true, "x"); err != nil || len(segs) != 1 {
		t.Fatal(segs, err)
	}
}

func TestParseMarkupErrors(t *testing.T) {
	cases := []struct {
		text string
		chat bool
		msg  string
	}{
		{"<Red>hi</>", false, "validation.markupTag"},
		{"<Bold>hi", false, "validation.markupUnclosed"},
		{"hi</>", false, "validation.markupClose"},
		{"<Bold><Header1>x</></>", false, "validation.markupNested"},
		{`<extimg src="http://x.com/a.png"/>`, false, "validation.markupImage"},
		{`<extimg src="https://x.com/a.png"/>`, true, "validation.markupTag"},
		{"<Header1>x</>", true, "validation.markupTag"},
	}
	for _, c := range cases {
		_, err := ParseMarkup(c.text, c.chat, "f")
		pe, ok := err.(*problems.Error)
		if !ok || pe.Msg != c.msg {
			t.Fatalf("%q: %v", c.text, err)
		}
	}
	p := Default()
	p.Motd = "<Bold>Welcome"
	if err := p.Validate(); err == nil {
		t.Fatal("unclosed MOTD tag accepted")
	}
	p.Motd, p.Rcon, p.JoinMessage = "<Bold>Welcome</>", true, "<Bold>{name}</> joined"
	if err := p.Validate(); err != nil {
		t.Fatal(err)
	}
}
