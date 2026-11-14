package modio

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// Mesma forma do state.json real (11/11/2026): o mod some da lista e da pasta; os outros ficam iguais.
func TestRemoveInstalled(t *testing.T) {
	root := t.TempDir()
	meta := filepath.Join(root, "254", "metadata")
	if err := os.MkdirAll(meta, 0o700); err != nil {
		t.Fatal(err)
	}
	state := `{"Mods":[` +
		`{"ID":111,"NeverRetryCategory":0,"PathOnDisk":"x","Profile":{"name":"Keep","extra":[1,2]},"SizeOnDisk":10,"State":1},` +
		`{"ID":222,"PathOnDisk":"y","Profile":{"name":"Drop"},"SizeOnDisk":20,"State":1}` +
		`],"version":3}`
	if err := os.WriteFile(filepath.Join(meta, "state.json"), []byte(state), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, id := range []int{111, 222} {
		if err := os.MkdirAll(filepath.Join(ModDir(root, id), "sub"), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(ModDir(root, id), "sub", "a.pak"), []byte("pak"), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	if err := RemoveInstalled(root, 222); err != nil {
		t.Fatal(err)
	}
	list, err := ReadInstalled(root)
	if err != nil || len(list) != 1 || list[0].ID != 111 || list[0].Name != "Keep" {
		t.Fatalf("state depois: %+v %v", list, err)
	}
	var top map[string]json.RawMessage
	b, _ := os.ReadFile(filepath.Join(meta, "state.json"))
	if json.Unmarshal(b, &top) != nil || string(top["version"]) != "3" {
		t.Fatalf("outros campos: %s", b)
	}
	if _, err := os.Stat(ModDir(root, 222)); !os.IsNotExist(err) {
		t.Error("pasta do mod continua")
	}
	if _, err := os.Stat(filepath.Join(ModDir(root, 111), "sub", "a.pak")); err != nil {
		t.Error("outro mod foi apagado")
	}
	if _, err := os.Stat(filepath.Join(meta, "state.json.bak")); err != nil {
		t.Error("sem cópia do state.json")
	}

	// Mod que não existe, id inválido: erro e nada muda.
	if err := RemoveInstalled(root, 222); err == nil {
		t.Error("mod já apagado aceito")
	}
	for _, id := range []int{0, -1} {
		if err := RemoveInstalled(root, id); err == nil {
			t.Errorf("id %d aceito", id)
		}
	}
	if _, err := os.Stat(ModDir(root, 111)); err != nil {
		t.Error("mod 111 sumiu")
	}
}

// Começar do zero com os mods (ADR-047): todos os mods do state.json e as pastas sem registro saem.
func TestRemoveAllInstalled(t *testing.T) {
	root := t.TempDir()
	meta := filepath.Join(root, "254", "metadata")
	if err := os.MkdirAll(meta, 0o700); err != nil {
		t.Fatal(err)
	}
	state := `{"Mods":[{"ID":111,"Profile":{"name":"A"},"State":1},{"ID":222,"Profile":{"name":"B"},"State":1}],"version":3}`
	if err := os.WriteFile(filepath.Join(meta, "state.json"), []byte(state), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, id := range []int{111, 222, 333} {
		if err := os.MkdirAll(ModDir(root, id), 0o700); err != nil {
			t.Fatal(err)
		}
	}
	if err := RemoveAllInstalled(root); err != nil {
		t.Fatal(err)
	}
	if list, err := ReadInstalled(root); err != nil || len(list) != 0 {
		t.Fatalf("state: %v %v", list, err)
	}
	entries, _ := os.ReadDir(filepath.Join(root, "254", "mods"))
	if len(entries) != 0 {
		t.Fatalf("pastas: %v", entries)
	}
	// Sem nada instalado: nada a fazer, sem erro.
	if err := RemoveAllInstalled(t.TempDir()); err != nil {
		t.Fatal(err)
	}
}
