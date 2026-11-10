package players

import (
	"testing"

	"sandstorm-server-assist/desktop/internal/domain/problems"
	"sandstorm-server-assist/desktop/internal/domain/server"
)

type fakeConn struct {
	sent  *[]string
	reply map[string]string
}

func (c fakeConn) Exec(cmd string) (string, error) {
	*c.sent = append(*c.sent, cmd)
	return c.reply[cmd], nil
}
func (fakeConn) Close() error { return nil }

const id = "SteamNWI:76561190000000001"

func newService(sent *[]string, dials *int, ready bool) *Service {
	return New(Deps{
		Config: func() (server.Profile, server.Secrets, error) {
			p := server.Default()
			p.Rcon = true
			return p, server.Secrets{RconPassword: "secret1"}, nil
		},
		Ready: func() bool { return ready },
		Dial: func(int, string) (Conn, error) {
			*dials++
			return fakeConn{sent: sent, reply: map[string]string{
				"listplayers": "ID | Name | NetID | IP | Score |\n=====\n0 | Fulano | " + id + " | 203.0.113.5 | 10 |\n",
				"listbans":    id + " banned\n",
			}}, nil
		},
		Log: func(string, string) {},
	})
}

func TestRefreshAndActions(t *testing.T) {
	var sent []string
	dials := 0
	s := newService(&sent, &dials, true)
	if err := s.Kick(id, ""); err == nil {
		t.Fatal("kick before the list was loaded")
	}
	v, err := s.Refresh()
	if err != nil || len(v.Players) != 1 || len(v.Bans) != 1 {
		t.Fatal(v, err)
	}
	if dials != 1 {
		t.Fatalf("refresh should use one connection, used %d", dials)
	}
	if err := s.Kick(id, "afk"); err != nil {
		t.Fatal(err)
	}
	if err := s.Ban(id, 30, ""); err != nil {
		t.Fatal(err)
	}
	if err := s.Unban(id); err != nil {
		t.Fatal(err)
	}
	if err := s.Say("hello all"); err != nil {
		t.Fatal(err)
	}
	if err := s.Kick("SteamNWI:76561190000000009", ""); err == nil {
		t.Fatal("kick of a player not in the list")
	}
	want := []string{"listplayers", "listbans", "kick " + id + " afk", "banid " + id + " 30", "unban " + id, "say hello all"}
	if len(sent) != len(want) {
		t.Fatal(sent)
	}
	for i := range want {
		if sent[i] != want[i] {
			t.Fatalf("%d: %q", i, sent[i])
		}
	}
}

func TestNeedsRunningServer(t *testing.T) {
	var sent []string
	dials := 0
	s := newService(&sent, &dials, false)
	if _, err := s.Refresh(); problems.CodeOf(err) != problems.ServerNotRunning {
		t.Fatal(err)
	}
}
