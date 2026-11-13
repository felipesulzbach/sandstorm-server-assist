package admin

import "testing"

// Formato real do listplayers (cabeçalho conferido no servidor em 10/11/2026); ids fictícios.
const sample = "ID\t | Name\t\t\t\t | NetID\t\t\t | IP\t\t\t | Score\t\t |\n" +
	"===============================================================================\n" +
	"0\t | Fulano\t | SteamNWI:76561190000000001\t | 203.0.113.5\t | 120\t |\n" +
	"1\t | Bot Marksman\t | INVALID\t | \t | 0\t |\n"

func TestParsePlayers(t *testing.T) {
	list := ParsePlayers(sample)
	if len(list) != 2 {
		t.Fatalf("got %d players", len(list))
	}
	p := list[0]
	if p.Name != "Fulano" || p.NetID != "SteamNWI:76561190000000001" || p.IP != "203.0.113.5" || p.Score != 120 || p.Bot {
		t.Fatalf("player: %+v", p)
	}
	if !list[1].Bot || list[1].NetID != "" {
		t.Fatalf("bot: %+v", list[1])
	}
	if len(ParsePlayers("ID\t | Name |\n=====\n")) != 0 {
		t.Fatal("empty table")
	}
}

// Formato real de 12/11/2026: todos os jogadores numa linha só (nomes, ids e IPs fictícios). O app mostrava só o
// primeiro.
const sampleOneLine = "ID\t | Name\t\t\t\t | NetID\t\t\t | IP\t\t\t | Score\t\t |\n" +
	"===============================================================================\n" +
	"257\t | Alpha Bravo\t | SteamNWI:76561190000000001\t | 203.0.113.5\t | 625\t\t | " +
	"258\t | [CL] Charlie \t | SteamNWI:76561190000000002\t | 203.0.113.6\t | 120\t\t | " +
	"1\t | Marksman\t | None:INVALID\t | \t | 3\t\t | " +
	"2\t | \t\t\t\t | None:INVALID\t | \t | 0\t\t | " +
	"259\t | Delta\t\t | SteamNWI:76561190000000003\t | 203.0.113.7\t | 0\t\t | "

func TestParsePlayersOneLine(t *testing.T) {
	list := ParsePlayers(sampleOneLine)
	if len(list) != 5 {
		t.Fatalf("got %d players: %+v", len(list), list)
	}
	want := []struct {
		id, name, netID string
		score           int
		bot             bool
	}{
		{"257", "Alpha Bravo", "SteamNWI:76561190000000001", 625, false},
		{"258", "[CL] Charlie", "SteamNWI:76561190000000002", 120, false},
		{"1", "Marksman", "", 3, true},
		{"2", "", "", 0, true},
		{"259", "Delta", "SteamNWI:76561190000000003", 0, false},
	}
	for i, w := range want {
		p := list[i]
		if p.ID != w.id || p.Name != w.name || p.NetID != w.netID || p.Score != w.score || p.Bot != w.bot {
			t.Errorf("%d: %+v, want %+v", i, p, w)
		}
	}
	if list[4].IP != "203.0.113.7" {
		t.Errorf("ip: %q", list[4].IP)
	}
}

func TestParseBans(t *testing.T) {
	out := "SteamNWI:76561190000000001 - 2026.11.10 - cheating\n76561190000000002 permanent\nnothing here\n"
	bans := ParseBans(out)
	if len(bans) != 2 || bans[0].NetID != "SteamNWI:76561190000000001" || bans[1].NetID != "76561190000000002" {
		t.Fatalf("bans: %+v", bans)
	}
}

func TestCommands(t *testing.T) {
	id := "SteamNWI:76561190000000001"
	if c, _ := KickCommand(id, "  team  killing "); c != "kick "+id+" team killing" {
		t.Fatal(c)
	}
	if c, _ := BanCommand(id, 60, ""); c != "banid "+id+" 60" {
		t.Fatal(c)
	}
	if c, _ := BanCommand(id, 0, "x"); c != "permban "+id+" x" {
		t.Fatal(c)
	}
	if c, _ := UnbanCommand(id); c != "unban "+id {
		t.Fatal(c)
	}
	for _, bad := range []string{"", "Fulano", "kick all", "SteamNWI:1;quit"} {
		if _, err := KickCommand(bad, ""); err == nil {
			t.Fatalf("accepted %q", bad)
		}
	}
	if _, err := KickCommand(id, "line\nbreak"); err == nil {
		t.Fatal("accepted control char")
	}
	if _, err := BanCommand(id, -1, ""); err == nil {
		t.Fatal("accepted negative minutes")
	}
}
