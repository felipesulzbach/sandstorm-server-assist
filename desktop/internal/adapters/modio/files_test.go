package modio

import (
	"os"
	"testing"
	"time"

	"sandstorm-server-assist/desktop/internal/domain/modcheck"
)

// TestInstalledFilesReal lê os mods instalados nesta máquina. Só com SSA_MODS_E2E=1 (não altera nada).
func TestInstalledFilesReal(t *testing.T) {
	if os.Getenv("SSA_MODS_E2E") != "1" {
		t.Skip("set SSA_MODS_E2E=1 to read the real installed mods")
	}
	start := time.Now()
	root := RootDir()
	if r := os.Getenv("SSA_MODS_ROOT"); r != "" {
		root = r
	}
	mf := InstalledFiles(root)
	t.Logf("mods=%d unreadable=%d in %v", len(mf.Mods), len(mf.Unreadable), time.Since(start))
	for _, u := range mf.Unreadable {
		t.Logf("  unreadable: %s #%d", u.Name, u.ID)
	}
	start = time.Now()
	c := modcheck.Conflicts(mf.Mods)
	t.Logf("conflicting pairs=%d in %v", len(c), time.Since(start))
	for i, x := range c {
		if i == 15 {
			break
		}
		t.Logf("  %d: %s <-> %s (e.g. %s)", x.Count, x.A.Name, x.B.Name, x.Example)
	}
	if obj := os.Getenv("SSA_MODS_OBJECT"); obj != "" {
		t.Logf("owners of %s: %+v", obj, modcheck.Owners(obj, mf.Mods))
	}
}

// Só os mods da conta do servidor (ADR-043): a pasta é dividida com o jogo da mesma máquina.
func TestForServer(t *testing.T) {
	list := []Installed{
		{Mod: Mod{ID: 1}, Subscribers: []int{111}},      // só a conta do jogador
		{Mod: Mod{ID: 2}, Subscribers: []int{111, 222}}, // jogador e servidor
		{Mod: Mod{ID: 3}}, // sem a lista (mod.io antigo)
	}
	ids := func(l []Installed) (out []int) {
		for _, m := range l {
			out = append(out, m.ID)
		}
		return
	}
	if got := ids(ForServer(list, 222)); len(got) != 2 || got[0] != 2 || got[1] != 3 {
		t.Fatalf("server 222: %v", got)
	}
	if got := ids(ForServer(list, 0)); len(got) != 1 || got[0] != 3 {
		t.Fatalf("no server account: %v", got)
	}
}

func TestReadInstalledSubscribers(t *testing.T) {
	root := t.TempDir()
	dir := root + `\254\metadata`
	if err := os.MkdirAll(dir, 0o750); err != nil {
		t.Fatal(err)
	}
	state := `{"Mods":[{"ID":7,"State":1,"SubscriptionCount":[4,5],"Profile":{"name":"X"}}],"version":1}`
	if err := os.WriteFile(dir+`\state.json`, []byte(state), 0o600); err != nil {
		t.Fatal(err)
	}
	list, err := ReadInstalled(root)
	if err != nil || len(list) != 1 || len(list[0].Subscribers) != 2 || list[0].Subscribers[1] != 5 {
		t.Fatalf("list %+v err %v", list, err)
	}
}

// TestGameAPIKeyReal lê a chave do mod.io do servidor instalado. Só com SSA_SERVER_DIR (não mostra a chave).
func TestGameAPIKeyReal(t *testing.T) {
	dir := os.Getenv("SSA_SERVER_DIR")
	if dir == "" {
		t.Skip("set SSA_SERVER_DIR to read the real server files")
	}
	start := time.Now()
	key := GameAPIKey(dir)
	t.Logf("found=%v length=%d in %v", key != "", len(key), time.Since(start))
	start = time.Now()
	_ = GameAPIKey(dir)
	t.Logf("cached in %v", time.Since(start))
	if key == "" {
		t.Fatal("no key")
	}
}
