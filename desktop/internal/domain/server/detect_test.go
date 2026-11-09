package server

import (
	"slices"
	"testing"
)

// Linha no formato que o servidor grava no log (LogInit: Command Line:), com valores de exemplo.
const exampleCmd = ` Farmhouse?Scenario=Scenario_Farmhouse_Checkpoint_Security?MaxPlayers=12?Password=abc12345 -SecurityCode=none -mods ` +
	`-ModDownloadTravelTo=Containment?Scenario=Scenario_Containment_Checkpoint_Security -Port=27200 -QueryPort=27231 ` +
	`-hostname="[XX] Example Server" -mutators=AllYouCanEat,HotPotato -Rcon -RconPassword=rc0nPass -RconListenPort=27300 ` +
	`-GSLTToken=ABCDEF0123456789 -GameStats -GameStatsToken=FEDCBA9876543210 -LogCmds="LogModio Verbose, LogModioGame Verbose"`

func TestParseCommandLine(t *testing.T) {
	d := NewDetected()
	d.ParseCommandLine(exampleCmd)
	p := d.Profile
	if p.Hostname != "[XX] Example Server" || p.Port != 27200 || p.QueryPort != 27231 || p.MaxPlayers != 12 {
		t.Errorf("geral: %+v", p)
	}
	if p.StartMap != "Farmhouse" || p.StartScenario != "Scenario_Farmhouse_Checkpoint_Security" {
		t.Errorf("mapa: %+v", p)
	}
	if p.TravelMap != "Containment" || p.TravelScenario != "Scenario_Containment_Checkpoint_Security" || !p.UseMods {
		t.Errorf("viagem: %+v", p)
	}
	if !slices.Equal(p.Mutators, []string{"AllYouCanEat", "HotPotato"}) || !p.Rcon || p.RconPort != 27300 {
		t.Errorf("mutators/rcon: %+v", p)
	}
	s := d.Secrets
	if s.ServerPassword != "abc12345" || s.RconPassword != "rc0nPass" || s.GSLTToken != "ABCDEF0123456789" ||
		s.GameStatsToken != "FEDCBA9876543210" {
		t.Errorf("segredos: %+v", s)
	}

	// Valores inválidos ficam de fora; o que falta na linha está desligado.
	bad := NewDetected()
	bad.ParseCommandLine(`Containment?Scenario_Containment_Checkpoint_Security -Port=GAME_PORT -QueryPort=27131`)
	if bad.Found["startMap"] || bad.Found["port"] || !bad.Found["queryPort"] || bad.Profile.UseMods || !bad.Found["useMods"] {
		t.Errorf("linha ruim: %+v %v", bad.Profile, bad.Found)
	}
	empty := NewDetected()
	empty.ParseCommandLine("   ")
	if len(empty.Found) != 0 {
		t.Errorf("linha vazia achou %v", empty.Found)
	}
}

func diffByKey(ds []Diff) map[string]Diff {
	m := map[string]Diff{}
	for _, d := range ds {
		m[d.Key] = d
	}
	return m
}

func TestCompareAndMerge(t *testing.T) {
	d := NewDetected()
	d.ParseCommandLine(exampleCmd)

	// O usuário mudou o nome e a porta; o resto está no padrão.
	cur := Default()
	cur.Hostname = "Nome que eu digitei"
	cur.Port = 27500
	cur.MaxPlayers = 12 // igual ao do servidor
	sec := Secrets{RconPassword: "outraSenha1"}

	diffs := diffByKey(Compare(cur, sec, d))
	want := map[string]string{
		"hostname": SyncConflict, "port": SyncConflict, "maxPlayers": SyncSame, "queryPort": SyncFill,
		"travelMap": SyncFill, "mutators": SyncFill, "rcon": SyncFill, "rconPort": SyncFill,
		"rconPassword": SyncConflict, "gsltToken": SyncFill, "serverPassword": SyncFill,
	}
	for k, st := range want {
		if diffs[k].State != st {
			t.Errorf("%s: estado %q, esperado %q", k, diffs[k].State, st)
		}
	}
	if diffs["rconPassword"].Mine != "" || diffs["rconPassword"].Server != "" {
		t.Error("segredo vazou na comparação")
	}
	if _, ok := diffs["admins"]; ok {
		t.Error("admins não achados não devem aparecer")
	}

	// Mantém o nome digitado; o resto vem do servidor.
	p, s, applied, skipped, err := Merge(cur, sec, d, map[string]bool{"hostname": true})
	if err != nil {
		t.Fatal(err)
	}
	if p.Hostname != "Nome que eu digitei" || p.Port != 27200 || p.QueryPort != 27231 || p.TravelMap != "Containment" {
		t.Errorf("perfil: %+v", p)
	}
	if s.RconPassword != "rc0nPass" || s.GSLTToken != "ABCDEF0123456789" {
		t.Errorf("segredos: %+v", s)
	}
	if slices.Contains(applied, "hostname") || !slices.Contains(applied, "port") || len(skipped) != 0 {
		t.Errorf("applied=%v skipped=%v", applied, skipped)
	}

	// Mantém tudo do usuário: só os campos "fill" mudam.
	p, s, _, _, _ = Merge(cur, sec, d, map[string]bool{"hostname": true, "port": true, "rconPassword": true})
	if p.Port != 27500 || s.RconPassword != "outraSenha1" || p.QueryPort != 27231 {
		t.Errorf("manter: %+v %+v", p, s)
	}
}

func TestMergeRevertsCollidingPorts(t *testing.T) {
	d := NewDetected()
	d.ParseCommandLine(`Farmhouse?Scenario=Scenario_Farmhouse_Checkpoint_Security -QueryPort=27102`)
	cur := Default() // Port 27102: a QueryPort do servidor bate com a porta do app
	p, _, applied, skipped, err := Merge(cur, Secrets{}, d, nil)
	if err != nil {
		t.Fatal(err)
	}
	if p.QueryPort != cur.QueryPort || !slices.Contains(skipped, "queryPort") || slices.Contains(applied, "queryPort") {
		t.Errorf("p=%+v applied=%v skipped=%v", p, applied, skipped)
	}
}
