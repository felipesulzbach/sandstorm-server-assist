package server

import (
	"strings"
	"testing"

	"sandstorm-server-assist/desktop/internal/domain/problems"
)

func msgOf(err error) string {
	if pe, ok := err.(*problems.Error); ok {
		return pe.Field + ":" + pe.Msg
	}
	return ""
}

func TestValidateMessages(t *testing.T) {
	p := Default()
	p.Motd = "Welcome!\nBe nice.\r\nHave fun."
	if err := p.Validate(); err != nil {
		t.Fatalf("multi-line motd: %v", err)
	}
	p.Motd = strings.Repeat("a", MaxMotd+1)
	if got := msgOf(p.Validate()); got != "motd:validation.motd" {
		t.Fatalf("long motd: %s", got)
	}
	p.Motd = "bell\a"
	if got := msgOf(p.Validate()); got != "motd:validation.motd" {
		t.Fatalf("control motd: %s", got)
	}

	p = Default()
	p.JoinMessage = "Welcome {name}"
	if got := msgOf(p.Validate()); got != "joinMessage:validation.chatNeedsRcon" {
		t.Fatalf("needs rcon: %s", got)
	}
	p.Rcon = true
	if err := p.Validate(); err != nil {
		t.Fatalf("join ok: %v", err)
	}
	p.LeaveMessage = "two\nlines"
	if got := msgOf(p.Validate()); got != "leaveMessage:validation.chatMessage" {
		t.Fatalf("multi-line chat: %s", got)
	}
	p.LeaveMessage = strings.Repeat("x", MaxChat+1)
	if got := msgOf(p.Validate()); got != "leaveMessage:validation.chatMessage" {
		t.Fatalf("long chat: %s", got)
	}
}

func TestMotdFile(t *testing.T) {
	if got := MotdFile("  Hello\nWorld \n"); got != "Hello\r\nWorld\r\n" {
		t.Fatalf("file = %q", got)
	}
	if MotdFile(" \n ") != "" {
		t.Fatal("empty motd must not write a file")
	}
	if got, ok := ParseMotd("Hello\r\nWorld\r\n"); !ok || got != "Hello\nWorld" {
		t.Fatalf("parse = %q %v", got, ok)
	}
	if _, ok := ParseMotd(""); ok {
		t.Fatal("empty file must not count")
	}
}

func TestChatCommand(t *testing.T) {
	for _, c := range []struct{ tpl, name, want string }{
		{"Welcome, {name}!", "Alpha", "say Welcome, Alpha!"},
		{"{name} joined. {name}, have fun", "Bravo", "say Bravo joined. Bravo, have fun"},
		// Nome com quebra de linha e controle não vira outro comando.
		{"Hi {name}", "Evil\nquit\x00", "say Hi Evil quit"},
		{"Hi {name}", strings.Repeat("n", 60), "say Hi " + strings.Repeat("n", 40)},
		{"No name here", "Alpha", "say No name here"},
		{"", "Alpha", ""},
		{"   ", "Alpha", ""},
	} {
		if got := ChatCommand(c.tpl, c.name); got != c.want {
			t.Fatalf("%q/%q = %q", c.tpl, c.name, got)
		}
	}
}
