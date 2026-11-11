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
