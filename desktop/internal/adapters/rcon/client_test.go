package rcon

import (
	"errors"
	"net"
	"strings"
	"testing"

	proto "sandstorm-server-assist/desktop/internal/domain/rcon"
)

// fakeServer imita o RCON do servidor: confere a senha e responde os comandos (a resposta pode vir em pedaços).
func fakeServer(t *testing.T, password string, replies map[string][]string) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = ln.Close() })
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer func() { _ = c.Close() }()
				write := func(p proto.Packet) {
					b, _ := proto.Encode(p)
					_, _ = c.Write(b)
				}
				for {
					p, err := proto.Decode(c)
					if err != nil {
						return
					}
					switch p.Type {
					case proto.TypeAuth:
						write(proto.Packet{ID: p.ID, Type: proto.TypeResponse}) // pacote vazio antes, como o Source
						id := p.ID
						if p.Body != password {
							id = proto.AuthFailedID
						}
						write(proto.Packet{ID: id, Type: proto.TypeAuthResponse})
					case proto.TypeExec:
						for _, part := range replies[p.Body] {
							write(proto.Packet{ID: p.ID, Type: proto.TypeResponse, Body: part})
						}
					}
				}
			}(conn)
		}
	}()
	return ln.Addr().(*net.TCPAddr).Port
}

func TestDialExec(t *testing.T) {
	port := fakeServer(t, "s3cret", map[string][]string{
		"scenarios": {"Scenario_Farmhouse_Checkpoint_Security\n", "Scenario_Containment_Checkpoint_Security\n"},
		"help":      {"commands"},
	})
	c, err := Dial(port, "s3cret")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = c.Close() }()
	got, err := c.Exec("scenarios")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "Farmhouse") || !strings.Contains(got, "Containment") {
		t.Fatalf("reply = %q", got)
	}
	// Segundo comando na mesma conexão.
	if got, err := c.Exec("help"); err != nil || got != "commands" {
		t.Fatalf("help = %q %v", got, err)
	}
}

func TestDialWrongPassword(t *testing.T) {
	port := fakeServer(t, "s3cret", nil)
	if _, err := Dial(port, "nope"); !errors.Is(err, ErrAuth) {
		t.Fatalf("err = %v", err)
	}
}

func TestDialClosedPort(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	_ = ln.Close()
	if _, err := Dial(port, "x"); err == nil {
		t.Fatal("expected an error on a closed port")
	}
	if _, err := Dial(0, "x"); err == nil {
		t.Fatal("expected an error on port 0")
	}
}
