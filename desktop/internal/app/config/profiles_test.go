package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"sandstorm-server-assist/desktop/internal/domain/ini"
	"sandstorm-server-assist/desktop/internal/domain/server"
)

func TestProfiles(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	s := Open(path, &memVault{})
	p := s.Profile()
	p.Hostname = "Coop server"
	if _, err := s.SaveProfile(p); err != nil {
		t.Fatal(err)
	}
	list, err := s.CreateProfile("Versus", false)
	if err != nil || len(list) != 2 || !list[0].Active || list[1].Name != "Versus" {
		t.Fatal(list, err)
	}
	if _, err := s.CreateProfile("  ", false); err == nil {
		t.Fatal("empty name accepted")
	}
	versus := list[1].ID
	got, err := s.SwitchProfile(versus)
	if err != nil || got.Hostname != server.Default().Hostname {
		t.Fatal(got.Hostname, err)
	}
	// Reabre: o perfil em uso e o guardado continuam.
	again := Open(path, &memVault{})
	refs := again.Profiles()
	if refs[0].ID != versus || refs[1].Hostname != "Coop server" {
		t.Fatalf("%+v", refs)
	}
	if _, err := again.DeleteProfile(versus); err == nil {
		t.Fatal("deleted the active profile")
	}
	if refs, err = again.DeleteProfile(refs[1].ID); err != nil || len(refs) != 1 {
		t.Fatal(refs, err)
	}
	if refs, err = again.RenameProfile(versus, "PvP"); err != nil || refs[0].Name != "PvP" {
		t.Fatal(refs, err)
	}
}

func TestPrepareWritesRules(t *testing.T) {
	dir := t.TempDir()
	game := filepath.Join(dir, GameIni)
	if err := os.MkdirAll(filepath.Dir(game), 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(game, []byte("[Rcon]\r\nbEnabled=True\r\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	p := server.Default()
	p.ServerDir = dir
	// Sem regras (nil): não mexe no Game.ini nem cria o Engine.ini.
	if err := Prepare(p); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, EngineIni)); err == nil {
		t.Fatal("Engine.ini created without rules")
	}
	p.Rules = map[string]string{"killFeed": "True", "tickRate": "60"}
	if err := Prepare(p); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(game)
	if v, _ := ini.Get(string(b), "/Script/Insurgency.INSGameMode", "bKillFeed"); v != "True" || !strings.Contains(string(b), "[Rcon]") {
		t.Fatal(string(b))
	}
	e, _ := os.ReadFile(filepath.Join(dir, EngineIni))
	if v, _ := ini.Get(string(e), "/Script/OnlineSubsystemUtils.IpNetDriver", "NetServerMaxTickRate"); v != "60" {
		t.Fatal(string(e))
	}
	// Trazer do servidor lê as regras de volta.
	s := Open(filepath.Join(dir, "config.json"), &memVault{})
	if _, err := s.SetServerDir(dir); err != nil {
		t.Fatal(err)
	}
	d, _, err := s.Detect("")
	if err != nil || d.Profile.Rules["killFeed"] != "True" || !d.Found["rules"] {
		t.Fatal(d.Profile.Rules, err)
	}
}
