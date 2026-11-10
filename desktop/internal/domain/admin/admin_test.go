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
