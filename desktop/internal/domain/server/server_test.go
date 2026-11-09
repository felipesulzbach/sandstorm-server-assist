package server

import (
	"errors"
	"strings"
	"testing"

	"sandstorm-server-assist/desktop/internal/domain/problems"
)

// O perfil que funcionou à mão em 09/11/2026 (resumo_sessao.md) tem de gerar a mesma linha do start_server.bat.
func TestBuildLaunchMatchesWorkingBat(t *testing.T) {
	p := Default()
	p.Hostname = "[BR] Example Server"
	p.MaxPlayers = 8
	p.TravelMap = "Containment"
	p.TravelScenario = "Scenario_Containment_Checkpoint_Security"
	l, err := BuildLaunch(p, Secrets{}, "")
	if err != nil {
		t.Fatal(err)
	}
	want := `Farmhouse?Scenario=Scenario_Farmhouse_Checkpoint_Security?MaxPlayers=8 -Port=27102 -QueryPort=27131 -hostname="[BR] Example Server" -SecurityCode=none -mods -ModDownloadTravelTo=Containment?Scenario=Scenario_Containment_Checkpoint_Security`
	if l.CmdLine != want {
		t.Errorf("\n got %s\nwant %s", l.CmdLine, want)
	}
}

func TestBuildLaunchSecretsAndMask(t *testing.T) {
	p := Default()
	p.Rcon = true
	p.Admins = []Admin{{SteamID: "76561198000000001", Name: "Eu"}}
	s := Secrets{ServerPassword: "entrar123", RconPassword: "rcon#Senha1", GSLTToken: "ABCDEF0123456789ABCDEF0123456789", GameStatsToken: "GS0123456789"}
	l, err := BuildLaunch(p, s, "PBRXD")
	if err != nil {
		t.Fatal(err)
	}
	for _, part := range []string{"?Password=entrar123", "-SecurityCode=PBRXD", "-Rcon -RconPassword=rcon#Senha1 -RconListenPort=27015", "-GSLTToken=ABCDEF", "-GameStats -GameStatsToken=GS0123456789", "-AdminList=Admins"} {
		if !strings.Contains(l.CmdLine, part) {
			t.Errorf("falta %q em %s", part, l.CmdLine)
		}
	}
	shown := l.Shown(s)
	for _, leak := range []string{"entrar123", "rcon#Senha1", "ABCDEF0123", "GS0123456789", "PBRXD"} {
		if strings.Contains(shown, leak) {
			t.Errorf("segredo %q visível: %s", leak, shown)
		}
	}
}

func TestValidateRejectsInjection(t *testing.T) {
	cases := []func(*Profile){
		func(p *Profile) { p.Hostname = `x" -Port=1` },
		func(p *Profile) { p.StartMap = "Farmhouse -mods" },
		func(p *Profile) { p.StartScenario = "Scenario_X?Password=1" },
		func(p *Profile) { p.TravelMap = "Containment"; p.TravelScenario = "Scenario_A B" },
		func(p *Profile) { p.Port = 80 },
		func(p *Profile) { p.QueryPort = p.Port },
		func(p *Profile) { p.Mutators = []string{"a,b"} },
		func(p *Profile) { p.Admins = []Admin{{SteamID: "123"}} },
		func(p *Profile) { p.ConnectIP = "abc" },
		func(p *Profile) { p.ServerDir = `relativo\pasta` },
	}
	for i, mut := range cases {
		p := Default()
		mut(&p)
		if p.Validate() == nil {
			t.Errorf("caso %d deveria ser recusado: %+v", i, p)
		}
	}
	if _, err := BuildLaunch(Default(), Secrets{ServerPassword: "a b?c"}, ""); err == nil {
		t.Error("senha com espaço e ? deveria ser recusada")
	}
	if _, err := BuildLaunch(Default(), Secrets{}, "12"); err == nil {
		t.Error("código curto deveria ser recusado")
	}
}

func TestEmailTypo(t *testing.T) {
	err := ValidateEmail("someone@gamil.com")
	var pe *problems.Error
	if !errors.As(err, &pe) || pe.Params["suggestion"] != "gmail.com" {
		t.Errorf("deveria sugerir gmail.com: %v", err)
	}
	if ValidateEmail("alguem@gmail.com") != nil {
		t.Error("e-mail válido recusado")
	}
}

func TestAdminsRoundTrip(t *testing.T) {
	in := ";Alice\r\n76561198000000001\r\n;Bob\r\n76561198000000002\r\n"
	admins := ParseAdmins(in)
	if len(admins) != 2 || admins[0].Name != "Alice" || admins[1].SteamID != "76561198000000002" {
		t.Fatalf("%+v", admins)
	}
	if AdminsFile(admins) != in {
		t.Errorf("ida e volta mudou: %q", AdminsFile(admins))
	}
	if got := ParseMapCycle("//Default\r\n(Scenario=\"A\")\r\n\r\n//x\r\n(Scenario=\"B\")"); len(got) != 2 {
		t.Errorf("%v", got)
	}
}

func TestSecretsWith(t *testing.T) {
	s, err := Secrets{}.With("rconPassword", "  boa#Senha1 ")
	if err != nil || s.RconPassword != "boa#Senha1" {
		t.Fatalf("%v %q", err, s.RconPassword)
	}
	if _, err := s.With("outra", "x"); err == nil {
		t.Error("chave desconhecida aceita")
	}
	if s2, _ := s.With("rconPassword", ""); s2.RconPassword != "" {
		t.Error("vazio deveria apagar")
	}
}

func TestFixMapCycleLine(t *testing.T) {
	cases := map[string]string{
		`(Scenario="Scenario=Scenario_Outback_Checkpoint_Security",Lighting="Day")`:        `(Scenario="Scenario_Outback_Checkpoint_Security",Lighting="Day")`,
		`(Scenario = "Scenario = Scenario_Border_Checkpoint_Insurgents",Lighting="Night")`: `(Scenario = "Scenario_Border_Checkpoint_Insurgents",Lighting="Night")`,
		`Scenario=Scenario_West_Checkpoint_Security`:                                       `Scenario_West_Checkpoint_Security`,
		`  (Scenario="Scenario_Containment_Checkpoint_Security",Lighting="Day")  `:         `(Scenario="Scenario_Containment_Checkpoint_Security",Lighting="Day")`,
		`Scenario_Farmhouse_Checkpoint_Security`:                                           `Scenario_Farmhouse_Checkpoint_Security`,
	}
	for in, want := range cases {
		if got := FixMapCycleLine(in); got != want {
			t.Errorf("FixMapCycleLine(%q) = %q, want %q", in, got, want)
		}
	}
	file := MapCycleFile([]string{`(Scenario="Scenario=Scenario_Outback_Checkpoint_Security",Lighting="Day")`})
	if file != "(Scenario=\"Scenario_Outback_Checkpoint_Security\",Lighting=\"Day\")\r\n" {
		t.Errorf("file %q", file)
	}
}
