package server

import (
	"strings"
	"testing"

	"sandstorm-server-assist/desktop/internal/domain/problems"
)

func TestInstanceFiles(t *testing.T) {
	main := InstanceFiles(0)
	if main != InstanceFiles(1) {
		t.Fatal("instance 0 (old profiles) must be the main one")
	}
	if main.GameIni != "Insurgency/Saved/Config/WindowsServer/Game.ini" || main.LogsDir != "Insurgency/Saved/Logs" ||
		main.Admins != "Insurgency/Config/Server/Admins.txt" || main.MapCycleName != "MapCycle" || main.MotdName != "Motd" {
		t.Errorf("main files changed: %+v", main)
	}
	// O jogo junta "_" ao sufixo (teste real: -saveddirsuffix=_ssatest criou Saved__ssatest).
	f := InstanceFiles(3)
	if f.Saved != "Insurgency/Saved_3" || f.EngineIni != "Insurgency/Saved_3/Config/WindowsServer/Engine.ini" ||
		f.Motd != "Insurgency/Config/Server/Motd_3.txt" || f.AdminsName != "Admins_3" || f.MapCycle != "Insurgency/Config/Server/MapCycle_3.txt" {
		t.Errorf("instance 3: %+v", f)
	}
}

func TestBuildLaunchInstance(t *testing.T) {
	p := Default()
	p.UseMods = false
	l, _ := BuildLaunch(p, Secrets{}, "")
	if strings.Contains(l.CmdLine, "saveddirsuffix") || strings.Contains(l.CmdLine, "-Motd") ||
		strings.Contains(l.CmdLine, "-AdminList") || strings.Contains(l.CmdLine, "-MapCycle") {
		t.Errorf("main server without lists: %s", l.CmdLine)
	}
	p.Admins = []Admin{{SteamID: "76561198000000001"}}
	p.MapCycle = []string{CycleLine("Scenario_Hideout_Checkpoint_Security", LightingDay)}
	l, _ = BuildLaunch(p, Secrets{}, "")
	if !strings.Contains(l.CmdLine, " -AdminList=Admins -MapCycle=MapCycle") {
		t.Errorf("main server with lists: %s", l.CmdLine)
	}

	p = Default()
	p.UseMods = false
	p.Instance = 2
	l, _ = BuildLaunch(p, Secrets{}, "")
	if !strings.Contains(l.CmdLine, " -saveddirsuffix=2 -AdminList=Admins_2 -Motd=Motd_2") || strings.Contains(l.CmdLine, "-MapCycle") {
		t.Errorf("instance 2 without lists: %s", l.CmdLine)
	}
	p.MapCycle = []string{CycleLine("Scenario_Hideout_Checkpoint_Security", LightingDay)}
	l, _ = BuildLaunch(p, Secrets{}, "")
	if !strings.Contains(l.CmdLine, " -MapCycle=MapCycle_2") {
		t.Errorf("instance 2 with a rotation: %s", l.CmdLine)
	}
	if InstanceOf(l.CmdLine) != 2 {
		t.Errorf("InstanceOf(own command line) = %d", InstanceOf(l.CmdLine))
	}
	p.Instance = MaxInstance + 1
	if _, err := BuildLaunch(p, Secrets{}, ""); err == nil {
		t.Error("instance out of range accepted")
	}
}

func TestInstanceOf(t *testing.T) {
	for cmd, want := range map[string]int{
		``: 1,
		`Insurgency Hideout?Scenario=Scenario_Hideout_Checkpoint_Security -Port=27102`:        1,
		`Hideout?Scenario=X -Port=27103 -saveddirsuffix=2 -AdminList=Admins_2`:                2,
		`Insurgency Hideout -SavedDirSuffix="12" -Port=1`:                                     12,
		`Hideout -saveddirsuffix=_ssatest`:                                                    0,
		`Hideout -saveddirsuffix=02`:                                                          0,
		`Hideout -saveddirsuffix=1`:                                                           0,
		`Hideout -saveddirsuffix=100`:                                                         0,
		`Hideout -saveddirsuffix=`:                                                            0,
		`Hideout -hostname="-saveddirsuffix=5"`:                                               1,
		`Hideout?Scenario=X -Port=27104 -hostname="My server" -saveddirsuffix=4 -Motd=Motd_4`: 4,
	} {
		if got := InstanceOf(cmd); got != want {
			t.Errorf("InstanceOf(%q) = %d, want %d", cmd, got, want)
		}
	}
}

func TestNextInstance(t *testing.T) {
	if n := NextInstance(nil); n != 1 {
		t.Errorf("first = %d", n)
	}
	if n := NextInstance([]int{0, 2, 4}); n != 3 {
		t.Errorf("0 counts as 1: got %d", n)
	}
}

func TestPortConflictsAndFreePorts(t *testing.T) {
	a := Default()
	b := Default()
	b.Port, b.QueryPort, b.RconPort = 27103, 27131, 27016
	b.Rcon = true
	c := PortConflicts(a, []Peer{{ID: "b", Name: "Versus", Profile: b}})
	if len(c) != 1 || c[0].Field != "queryPort" || c[0].Port != 27131 || c[0].Server != "Versus" {
		t.Errorf("conflicts = %+v", c)
	}
	// RCON é TCP: a mesma porta em UDP de outro servidor não conflita.
	a.Rcon, a.RconPort = true, 27103
	a.QueryPort = 27132
	if c := PortConflicts(a, []Peer{{Name: "Versus", Profile: b}}); len(c) != 0 {
		t.Errorf("tcp vs udp: %+v", c)
	}
	// RCON desligado no outro servidor não ocupa a porta.
	b.Rcon = false
	a.RconPort = 27016
	if c := PortConflicts(a, []Peer{{Name: "Versus", Profile: b}}); len(c) != 0 {
		t.Errorf("rcon off: %+v", c)
	}

	port, query, rcon := FreePorts(nil, nil)
	if port != 27102 || query != 27131 || rcon != 27015 {
		t.Errorf("first free ports = %d %d %d", port, query, rcon)
	}
	first := Default()
	first.Rcon = true
	port, query, rcon = FreePorts([]Profile{first}, nil)
	if port != 27103 || query != 27132 || rcon != 27016 {
		t.Errorf("second server ports = %d %d %d", port, query, rcon)
	}
	n := Default()
	n.Port, n.QueryPort, n.RconPort = port, query, rcon
	n.Rcon = true
	if c := PortConflicts(n, []Peer{{Profile: first}}); len(c) != 0 {
		t.Errorf("free ports conflict: %+v", c)
	}
}

// 14/11/2026: o MongoDB do PC do mantenedor escuta em 127.0.0.1:27017 e o terceiro servidor ganhou essa porta de
// RCON; o RCON dele não respondia (as conexões iam para o MongoDB).
func TestFreePortsSkipsBusyPorts(t *testing.T) {
	a, b := Default(), Default()
	b.Port, b.QueryPort, b.RconPort = 27103, 27132, 27016
	mongo := func(port int, tcp bool) bool { return tcp && port == 27017 }
	port, query, rcon := FreePorts([]Profile{a, b}, mongo)
	if port != 27105 || query != 27134 || rcon != 27018 {
		t.Errorf("ports = %d %d %d", port, query, rcon)
	}
	c := Default()
	c.Rcon, c.RconPort = true, 27017
	got := BusyPorts(c, mongo)
	if len(got) != 1 || got[0].Field != "rconPort" || got[0].Port != 27017 || got[0].Server != "" {
		t.Errorf("busy = %+v", got)
	}
	if len(BusyPorts(c, nil)) != 0 {
		t.Error("nil check")
	}
}

func TestCheckStart(t *testing.T) {
	a := Default()
	b := Default()
	b.Port, b.QueryPort = 27103, 27132
	if err := CheckStart(a, Secrets{GSLTToken: "token-aaaa"}, []Peer{{Name: "B", Profile: b, Secrets: Secrets{GSLTToken: "token-bbbb"}}}); err != nil {
		t.Fatalf("distinct servers refused: %v", err)
	}
	err := CheckStart(a, Secrets{}, []Peer{{Name: "B", Profile: a}})
	var pe *problems.Error
	if !asProblem(err, &pe) || pe.Code != problems.PortsShared || pe.Params["server"] != "B" || pe.Params["port"] != "27102" || pe.Params["portField"] != "port" {
		t.Errorf("same ports: %v", err)
	}
	err = CheckStart(a, Secrets{GSLTToken: "token-aaaa"}, []Peer{{Name: "B", Profile: b, Secrets: Secrets{GSLTToken: "token-aaaa"}}})
	if !asProblem(err, &pe) || pe.Code != problems.GSLTShared || pe.Params["server"] != "B" {
		t.Errorf("same GSLT: %v", err)
	}
}

func TestSecretsSplit(t *testing.T) {
	s := Secrets{ServerPassword: "pw1234", ModioEmail: "a@example.com", ModioAPIKey: strings.Repeat("a", 32)}
	if s.ServerOnly().ModioEmail != "" || s.ServerOnly().ServerPassword != "pw1234" {
		t.Error("ServerOnly")
	}
	other := Secrets{GSLTToken: "token-1234"}.WithAccount(s)
	if other.ModioEmail != s.ModioEmail || other.ModioAPIKey != s.ModioAPIKey || other.GSLTToken != "token-1234" || other.ServerPassword != "" {
		t.Errorf("WithAccount = %+v", other)
	}
	for _, k := range SecretKeys {
		want := k == "modioEmail" || k == "modioApiKey"
		if IsAccountKey(k) != want {
			t.Errorf("IsAccountKey(%s)", k)
		}
	}
}

func asProblem(err error, pe **problems.Error) bool {
	p, ok := err.(*problems.Error)
	*pe = p
	return ok
}
