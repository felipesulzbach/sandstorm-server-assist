package server

import (
	"strings"
	"testing"
	"time"

	"sandstorm-server-assist/desktop/internal/domain/ini"
)

func TestNormalizeRule(t *testing.T) {
	cases := map[[2]string]string{
		{"killFeed", "true"}: "True", {"killFeed", "0"}: "False", {"aiDifficulty", "0,75"}: "0.75",
		{"teamKillLimit", " 5 "}: "5", {"tickRate", ""}: "",
	}
	for in, want := range cases {
		got, err := NormalizeRule(in[0], in[1])
		if err != nil || got != want {
			t.Fatalf("%v: %q %v", in, got, err)
		}
	}
	for _, bad := range [][2]string{{"aiDifficulty", "2"}, {"killFeed", "yes"}, {"tickRate", "500"}, {"nope", "1"}} {
		if _, err := NormalizeRule(bad[0], bad[1]); err == nil {
			t.Fatalf("accepted %v", bad)
		}
	}
}

func TestApplyAndReadRules(t *testing.T) {
	game := "[Rcon]\r\nbEnabled=True\r\n\r\n[/Script/Insurgency.INSGameMode]\r\nbKillFeed=False\r\nTeamKillLimit=9\r\n"
	rules := map[string]string{"killFeed": "True", "aiDifficulty": "0.8", "voteKick": "True", "tickRate": "64"}
	g, e := ApplyRules(game, "", rules)
	if v, _ := ini.Get(g, "/Script/Insurgency.INSGameMode", "bKillFeed"); v != "True" {
		t.Fatalf("kill feed: %s", g)
	}
	if _, ok := ini.Get(g, "/Script/Insurgency.INSGameMode", "TeamKillLimit"); ok {
		t.Fatal("empty rule should be removed")
	}
	if v, _ := ini.Get(g, "/Script/Insurgency.TeamInfo", "TeamVoteIssues"); v != "/Script/Insurgency.VoteIssueKick" {
		t.Fatal("vote issue missing")
	}
	if v, _ := ini.Get(g, "Rcon", "bEnabled"); v != "True" {
		t.Fatal("other sections must stay")
	}
	if v, _ := ini.Get(e, "/Script/OnlineSubsystemUtils.IpNetDriver", "NetServerMaxTickRate"); v != "64" {
		t.Fatal("engine tick rate")
	}
	back := ReadRules(g, e)
	for k, v := range rules {
		if back[k] != v {
			t.Fatalf("read %s=%q", k, back[k])
		}
	}
	// Desligar a votação tira a lista de assuntos que o app pôs.
	g2, _ := ApplyRules(g, e, map[string]string{"voteKick": "False"})
	if _, ok := ini.Get(g2, "/Script/Insurgency.TeamInfo", "TeamVoteIssues"); ok {
		t.Fatal("vote issue should be removed")
	}
	// nil: não mexe.
	if g3, _ := ApplyRules(game, "", nil); g3 != game {
		t.Fatal("nil rules changed the file")
	}
}

func TestRulesValidate(t *testing.T) {
	p := Default()
	p.Rules = map[string]string{"minimumEnemies": "20", "maximumEnemies": "10"}
	if err := p.Validate(); err == nil {
		t.Fatal("min > max accepted")
	}
	p.Rules = map[string]string{"minimumEnemies": "5", "maximumEnemies": "10"}
	if err := p.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestLighting(t *testing.T) {
	if got := SetLineLighting("Scenario_Farmhouse_Checkpoint_Security", LightingNight); got != `(Scenario="Scenario_Farmhouse_Checkpoint_Security",Lighting="Night")` {
		t.Fatal(got)
	}
	line := `(Scenario="Scenario_Town_Checkpoint_Security",Lighting="Day")`
	if got := SetLineLighting(line, LightingNight); LineLighting(got) != LightingNight || strings.Count(got, "Lighting") != 1 {
		t.Fatal(got)
	}
	if got := SetLineLighting(line, ""); LineLighting(got) != "" || LineScenario(got) != "Scenario_Town_Checkpoint_Security" {
		t.Fatal(got)
	}
	if got := SetLineLighting("whatever", LightingNight); got != "whatever" {
		t.Fatal(got)
	}
	p := Default()
	p.StartLighting = LightingNight
	l, err := BuildLaunch(p, Secrets{}, "")
	if err != nil || !strings.Contains(l.Args[0], "?Lighting=Night") {
		t.Fatal(l.Args, err)
	}
	p.StartLighting = "Dusk"
	if p.Validate() == nil {
		t.Fatal("bad lighting accepted")
	}
}

func TestLaunchExtras(t *testing.T) {
	p := Default()
	p.BindIP = "10.0.0.5"
	p.ModIDs, p.CmdModList = []int{11, 22}, true
	l, err := BuildLaunch(p, Secrets{}, "")
	if err != nil || !strings.Contains(l.CmdLine, "-MultiHome=10.0.0.5") || !strings.Contains(l.CmdLine, `-CmdModList="11,22"`) {
		t.Fatal(l.CmdLine, err)
	}
	d := NewDetected()
	d.ParseCommandLine(l.CmdLine)
	if d.Profile.BindIP != "10.0.0.5" || !d.Profile.CmdModList || len(d.Profile.ModIDs) != 2 {
		t.Fatalf("%+v", d.Profile)
	}
	p.BindIP = "0.0.0.0"
	if p.Validate() == nil {
		t.Fatal("0.0.0.0 accepted")
	}
	p.BindIP, p.ModIDs = "", nil
	if p.Validate() == nil {
		t.Fatal("empty mod list accepted")
	}
}

func TestAutomation(t *testing.T) {
	p := Default()
	p.Announcements = []string{"Rules: no team kill"}
	if p.Validate() == nil {
		t.Fatal("announcements without rcon accepted")
	}
	p.Rcon = true
	if err := p.Validate(); err != nil {
		t.Fatal(err)
	}
	p.AnnounceEvery = 1
	if p.Validate() == nil {
		t.Fatal("interval too short accepted")
	}
	p.AnnounceEvery = 15
	p.AutoRestart.Daily = "25:00"
	if p.Validate() == nil {
		t.Fatal("bad clock accepted")
	}
	loc := time.FixedZone("x", -3*3600)
	now := time.Date(2026, 11, 10, 5, 0, 0, 0, loc)
	if next, _ := NextDaily(now, "04:30"); !next.Equal(time.Date(2026, 11, 11, 4, 30, 0, 0, loc)) {
		t.Fatal(next)
	}
	if next, _ := NextDaily(now, "06:00"); !next.Equal(time.Date(2026, 11, 10, 6, 0, 0, 0, loc)) {
		t.Fatal(next)
	}
	got := DiscordText("{server} online: {map} @everyone", "My *Server*", "Farmhouse", "")
	if !strings.Contains(got, `My \*Server\*`) || strings.Contains(got, "{map}") {
		t.Fatal(got)
	}
	if got := DiscordText("{name} joined", "", "", "@everyone"); strings.Contains(got, "@everyone") {
		t.Fatal(got)
	}
	var s Secrets
	if _, err := s.With("discordWebhook", "https://evil.example/api/webhooks/1/2"); err == nil {
		t.Fatal("foreign webhook accepted")
	}
	if _, err := s.With("discordWebhook", "https://discord.com/api/webhooks/123456789012/abcdefghijklmnopqrstuvwxyz_-ABC"); err != nil {
		t.Fatal(err)
	}
}
