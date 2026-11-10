package config

import (
	"os"
	"path/filepath"
	"testing"

	"sandstorm-server-assist/desktop/internal/domain/server"
)

func TestPrepareWritesMotd(t *testing.T) {
	dir := t.TempDir()
	p := server.Default()
	p.ServerDir = dir
	if err := Prepare(p); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, MotdTxt)); !os.IsNotExist(err) {
		t.Fatalf("empty motd must not create the file: %v", err)
	}
	p.Motd = "Welcome!\nHave fun."
	if err := Prepare(p); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(dir, MotdTxt))
	if err != nil || string(b) != "Welcome!\r\nHave fun.\r\n" {
		t.Fatalf("Motd.txt = %q %v", b, err)
	}
}
