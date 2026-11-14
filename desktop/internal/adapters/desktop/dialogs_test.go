package desktop

import (
	"os"
	"path/filepath"
	"testing"
)

// 14/11/2026: a pasta do servidor foi movida e o "Escolher..." dava erro desconhecido, sem abrir o diálogo.
func TestDialogStart(t *testing.T) {
	root := t.TempDir()
	kept := filepath.Join(root, "kept")
	if err := os.MkdirAll(kept, 0o700); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(root, "file.txt")
	if err := os.WriteFile(file, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	cases := map[string]string{
		kept: kept,
		filepath.Join(kept, "moved", "sandstorm_server"): kept,
		filepath.Join(root, "gone"):                      root,
		file:                                             root,
		"":                                               "",
	}
	for in, want := range cases {
		if got := dialogStart(in); got != want {
			t.Errorf("dialogStart(%q) = %q, want %q", in, got, want)
		}
	}
	// Disco que não existe: nenhuma pasta-mãe existe.
	if got := dialogStart(`Q:\nada\aqui`); got != "" && got != `Q:\` {
		t.Errorf("disco inexistente: %q", got)
	}
}
