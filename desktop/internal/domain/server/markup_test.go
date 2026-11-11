package server

import (
	"testing"

	"sandstorm-server-assist/desktop/internal/domain/problems"
)

// O jogo não interpreta tags na tela de carregamento nem no chat (11/11/2026): mensagem com tag é recusada.
func TestMessagesRejectTags(t *testing.T) {
	cases := []struct{ field, text string }{
		{"motd", "<Bold>Welcome</>"},
		{"motd", `Hi <extimg src="https://x.com/a.png"/>`},
		{"joinMessage", "<Bold>{name}</> joined"},
		{"leaveMessage", "bye </>"},
		{"announcements", "<Header1>Discord</>"},
	}
	for _, c := range cases {
		p := Default()
		p.Rcon = true
		switch c.field {
		case "motd":
			p.Motd = c.text
		case "joinMessage":
			p.JoinMessage = c.text
		case "leaveMessage":
			p.LeaveMessage = c.text
		case "announcements":
			p.Announcements = []string{c.text}
		}
		err := p.Validate()
		pe, ok := err.(*problems.Error)
		if !ok || pe.Msg != "validation.markupPlain" || pe.Field != c.field {
			t.Errorf("%s %q: %v", c.field, c.text, err)
		}
	}
	// Texto comum com < e > não é tag.
	p := Default()
	p.Rcon, p.Motd, p.JoinMessage = true, "a < b and <c> d\nsecond line", "{name} <3"
	if err := p.Validate(); err != nil {
		t.Fatal(err)
	}
}

// Mensagens salvas com formatação antes da correção saem sem as tags, no Motd.txt e no chat.
func TestSavedMarkupIsStripped(t *testing.T) {
	if got := MotdFile("<Header1>Welcome</>\nPlay <Bold>fair</>"); got != "Welcome\r\nPlay fair\r\n" {
		t.Errorf("motd %q", got)
	}
	for _, c := range []struct{ tpl, want string }{
		{"<Bold>{name}</> joined", "say Alpha joined"},
		{`Welcome to <Header1>our</> server <extimg src="https://x.com/a.png"/>`, "say Welcome to our server"},
		{"a < b and <c> stays", "say a < b and <c> stays"},
	} {
		if got := ChatCommand(c.tpl, "Alpha"); got != c.want {
			t.Errorf("%q = %q, want %q", c.tpl, got, c.want)
		}
	}
}
