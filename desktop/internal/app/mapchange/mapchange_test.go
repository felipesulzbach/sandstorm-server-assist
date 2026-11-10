package mapchange

import (
	"errors"
	"testing"

	"sandstorm-server-assist/desktop/internal/domain/problems"
	"sandstorm-server-assist/desktop/internal/domain/server"
)

type fakeConn struct {
	replies map[string]string
	sent    *[]string
}

func (f fakeConn) Exec(cmd string) (string, error) {
	*f.sent = append(*f.sent, cmd)
	return f.replies[cmd], nil
}
func (fakeConn) Close() error { return nil }

type env struct {
	profile server.Profile
	secrets server.Secrets
	ready   bool
	dialErr error
	replies map[string]string
	sent    []string
	dials   int
	logs    []string
}

func newEnv() *env {
	p := server.Default()
	p.Rcon = true
	p.TravelMap, p.TravelScenario = "Containment", "Scenario_Containment_Checkpoint_Security"
	return &env{
		profile: p,
		secrets: server.Secrets{RconPassword: "s3cret"},
		ready:   true,
		replies: map[string]string{"scenarios": "Scenario_Farmhouse_Checkpoint_Security\nScenario_Gap_Push_Insurgents\n"},
	}
}

func (e *env) service() *Service {
	return New(Deps{
		Config: func() (server.Profile, server.Secrets, error) { return e.profile, e.secrets, nil },
		Ready:  func() bool { return e.ready },
		Dial: func(port int, password string) (Conn, error) {
			e.dials++
			if e.dialErr != nil {
				return nil, e.dialErr
			}
			return fakeConn{replies: e.replies, sent: &e.sent}, nil
		},
		Log: func(text, _ string) { e.logs = append(e.logs, text) },
	})
}

func code(err error) string { return problems.CodeOf(err) }

func TestScenariosAndChange(t *testing.T) {
	e := newEnv()
	s := e.service()
	if !s.Info().Rcon {
		t.Fatal("rcon should be available")
	}
	list, err := s.Scenarios()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 3 || list[0].ID != "Scenario_Containment_Checkpoint_Security" || list[0].Source != server.FromConfig {
		t.Fatalf("list = %+v", list)
	}
	// Mapa de mod do perfil: entra mesmo sem o servidor listar.
	if err := s.Change("Scenario_Containment_Checkpoint_Security"); err != nil {
		t.Fatal(err)
	}
	if last := e.sent[len(e.sent)-1]; last != "travelscenario Scenario_Containment_Checkpoint_Security" {
		t.Fatalf("sent %q", last)
	}
	if e.dials != 2 {
		t.Fatalf("one connection per action, got %d", e.dials)
	}
}

func TestChangeNeedsListedScenario(t *testing.T) {
	e := newEnv()
	s := e.service()
	if err := s.Change("Scenario_Gap_Push_Insurgents"); code(err) != problems.InvalidField {
		t.Fatalf("without a list: %v", err)
	}
	if _, err := s.Scenarios(); err != nil {
		t.Fatal(err)
	}
	if err := s.Change("Scenario_Tell_Checkpoint_Security"); code(err) != problems.InvalidField {
		t.Fatalf("not listed: %v", err)
	}
	if e.dials != 1 {
		t.Fatalf("refused ids must not connect, dials = %d", e.dials)
	}
}

func TestConditions(t *testing.T) {
	e := newEnv()
	e.profile.Rcon = false
	if _, err := e.service().Scenarios(); code(err) != problems.RconOff {
		t.Fatalf("rcon off: %v", err)
	}
	e = newEnv()
	e.secrets.RconPassword = ""
	if e.service().Info().Rcon {
		t.Fatal("no password: rcon must be unavailable")
	}
	e = newEnv()
	e.ready = false
	if _, err := e.service().Scenarios(); code(err) != problems.ServerNotRunning {
		t.Fatalf("not ready: %v", err)
	}
	e = newEnv()
	e.dialErr = ErrAuth
	if _, err := e.service().Scenarios(); code(err) != problems.RconAuth {
		t.Fatalf("auth: %v", err)
	}
	e = newEnv()
	e.dialErr = errors.New("connection refused")
	if _, err := e.service().Scenarios(); code(err) != problems.RconFailed {
		t.Fatalf("dial: %v", err)
	}
}

func TestChangeRefusedByServer(t *testing.T) {
	e := newEnv()
	e.replies["travelscenario Scenario_Gap_Push_Insurgents"] = "Unknown scenario"
	s := e.service()
	if _, err := s.Scenarios(); err != nil {
		t.Fatal(err)
	}
	if err := s.Change("Scenario_Gap_Push_Insurgents"); code(err) != problems.MapChangeFailed {
		t.Fatalf("refused: %v", err)
	}
}
