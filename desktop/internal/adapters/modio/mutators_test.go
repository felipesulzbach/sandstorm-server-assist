package modio

import (
	"os"
	"testing"
)

// Teste opcional com os mods instalados nesta máquina (SSA_MODIO_TEST=1): só lê.
func TestMutatorsReal(t *testing.T) {
	if os.Getenv("SSA_MODIO_TEST") == "" {
		t.Skip("set SSA_MODIO_TEST=1 to read the installed mods")
	}
	for _, m := range Mutators(RootDir()) {
		t.Logf("%s (mod %d)", m.ID, m.ModID)
	}
}
