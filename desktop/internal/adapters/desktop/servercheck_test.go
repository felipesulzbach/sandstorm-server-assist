package desktop

import (
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"sandstorm-server-assist/desktop/internal/app/config"
)

// Aviso de versão nova do servidor: compara a versão publicada com a instalada agora (depois de atualizar, some).
func TestServerUpdate(t *testing.T) {
	dir := t.TempDir()
	cfg := config.Open(filepath.Join(dir, "config.json"), memoryVault{})
	srvDir := filepath.Join(dir, "server")
	if err := os.MkdirAll(filepath.Join(srvDir, "steamapps"), 0o750); err != nil {
		t.Fatal(err)
	}
	if _, err := cfg.SetServerDir(srvDir); err != nil {
		t.Fatal(err)
	}
	a := &app{cfg: cfg}
	if u := a.serverUpdate(); u.Available || u.Installed != "" {
		t.Fatalf("not installed: %+v", u)
	}
	write := func(build string) {
		t.Helper()
		acf := "\"AppState\"\n{\n\t\"appid\"\t\t\"581330\"\n\t\"buildid\"\t\t\"" + build + "\"\n}\n"
		if err := os.WriteFile(filepath.Join(srvDir, "steamapps", "appmanifest_581330.acf"), []byte(acf), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(srvDir, "InsurgencyServer.exe"), nil, 0o600); err != nil {
		t.Fatal(err)
	}
	write("100")
	if u := a.serverUpdate(); u.Available || u.Installed != "100" {
		t.Fatalf("before any check: %+v", u)
	}
	a.srvCheck.set("101", time.Now())
	if u := a.serverUpdate(); !u.Available || u.Latest != "101" {
		t.Fatalf("newer build: %+v", u)
	}
	write("101")
	if u := a.serverUpdate(); u.Available {
		t.Fatalf("after updating: %+v", u)
	}
}

// Links novos de "Ligações úteis": lojas onde o jogo é vendido para PC e página oficial, só https.
func TestGameLinks(t *testing.T) {
	for _, k := range []string{"steam-store", "epic-store", "official-site", "modio"} {
		u, err := url.Parse(links[k])
		if err != nil || u.Scheme != "https" || u.Host == "" {
			t.Errorf("%s: %q", k, links[k])
		}
	}
}
