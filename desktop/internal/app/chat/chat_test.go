package chat

import (
	"errors"
	"testing"
	"time"

	"sandstorm-server-assist/desktop/internal/domain/logparse"
	"sandstorm-server-assist/desktop/internal/domain/server"
)

var now = time.Date(2026, 11, 10, 3, 0, 0, 0, time.UTC)

type fakeConn struct {
	sent   *[]string
	fail   *int
	closed *int
}

func (f fakeConn) Exec(cmd string) (string, error) {
	if *f.fail > 0 {
		*f.fail--
		return "", errors.New("broken pipe")
	}
	*f.sent = append(*f.sent, cmd)
	return "", nil
}
func (f fakeConn) Close() error { *f.closed++; return nil }

type env struct {
	p             server.Profile
	sec           server.Secrets
	sent          []string
	dials, closed int
	fail          int
	logs          []string
}

func newEnv() *env {
	p := server.Default()
	p.Rcon = true
	p.JoinMessage, p.LeaveMessage = "Welcome, {name}!", "{name} left."
	return &env{p: p, sec: server.Secrets{RconPassword: "s3cret"}}
}

func (e *env) service() *Service {
	return New(Deps{
		Config: func() (server.Profile, server.Secrets, error) { return e.p, e.sec, nil },
		Dial: func(int, string) (Conn, error) {
			e.dials++
			return fakeConn{sent: &e.sent, fail: &e.fail, closed: &e.closed}, nil
		},
		Log: func(t, _ string) { e.logs = append(e.logs, t) },
		Now: func() time.Time { return now },
	}, false)
}

func ev(kind logparse.Kind, name, id string, at time.Time) logparse.Event {
	return logparse.Event{Kind: kind, Player: name, PlayerID: id, Time: at}
}

func TestJoinAndLeave(t *testing.T) {
	e := newEnv()
	s := e.service()
	at := now.Add(-5 * time.Second)
	if cmd := s.Handle(ev(logparse.PlayerLogin, "[TAG] Long Player Name", "SteamNWI:1", at)); cmd != "" {
		t.Fatalf("login must not send: %q", cmd)
	}
	// A entrada confirmada traz o nome cortado: a mensagem usa o nome completo do pedido.
	if cmd := s.Handle(ev(logparse.PlayerJoined, "[TAG] Long ", "", at)); cmd != "say Welcome, [TAG] Long Player Name!" {
		t.Fatalf("join = %q", cmd)
	}
	if cmd := s.Handle(ev(logparse.PlayerLeft, "", "SteamNWI:1", at)); cmd != "say [TAG] Long Player Name left." {
		t.Fatalf("leave = %q", cmd)
	}
	// Saída de quem não entrou (conexão que caiu no meio) não gera mensagem.
	if cmd := s.Handle(ev(logparse.PlayerLeft, "", "SteamNWI:2", at)); cmd != "" {
		t.Fatalf("unknown leave = %q", cmd)
	}
}

func TestOldLinesAndSettings(t *testing.T) {
	e := newEnv()
	s := e.service()
	old := now.Add(-10 * time.Minute)
	s.Handle(ev(logparse.PlayerLogin, "Alpha", "SteamNWI:1", old))
	if cmd := s.Handle(ev(logparse.PlayerJoined, "Alpha", "", old)); cmd != "" {
		t.Fatalf("old join must not send: %q", cmd)
	}
	// Mas quem entrou antes ainda gera a mensagem de saída agora.
	if cmd := s.Handle(ev(logparse.PlayerLeft, "", "SteamNWI:1", now)); cmd != "say Alpha left." {
		t.Fatalf("leave after old join = %q", cmd)
	}

	e.p.JoinMessage = ""
	if cmd := s.Handle(ev(logparse.PlayerJoined, "Bravo", "", now)); cmd != "" {
		t.Fatalf("empty message must not send: %q", cmd)
	}
	e.p.JoinMessage, e.p.Rcon = "Hi {name}", false
	if cmd := s.Handle(ev(logparse.PlayerJoined, "Bravo", "", now)); cmd != "" {
		t.Fatalf("rcon off must not send: %q", cmd)
	}
}

func TestSendReusesConnection(t *testing.T) {
	e := newEnv()
	s := e.service()
	s.Send("say one")
	s.Send("say two")
	if e.dials != 1 || len(e.sent) != 2 {
		t.Fatalf("dials=%d sent=%v", e.dials, e.sent)
	}
	// Conexão caiu: fecha, abre outra e manda.
	e.fail = 1
	s.Send("say three")
	if e.dials != 2 || e.closed != 1 || e.sent[len(e.sent)-1] != "say three" {
		t.Fatalf("retry: dials=%d closed=%d sent=%v", e.dials, e.closed, e.sent)
	}
	s.Reset()
	if e.closed != 2 {
		t.Fatalf("reset must close the connection, closed=%d", e.closed)
	}
}

func TestReset(t *testing.T) {
	e := newEnv()
	s := e.service()
	s.Handle(ev(logparse.PlayerLogin, "Alpha", "SteamNWI:1", now))
	s.Handle(ev(logparse.PlayerJoined, "Alpha", "", now))
	s.Reset()
	if cmd := s.Handle(ev(logparse.PlayerLeft, "", "SteamNWI:1", now)); cmd != "" {
		t.Fatalf("after reset = %q", cmd)
	}
}
