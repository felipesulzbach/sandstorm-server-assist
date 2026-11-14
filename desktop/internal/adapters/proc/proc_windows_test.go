//go:build windows

package proc

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// O próprio binário do teste aparece com a linha de comando (é por ela que o app separa as instâncias, ADR-048).
func TestFindInDirReadsCommandLine(t *testing.T) {
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	name := strings.TrimSuffix(filepath.Base(exe), ".exe")
	list, err := FindInDir(filepath.Dir(exe), name)
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range list {
		if IsSelf(p.PID) {
			if !strings.Contains(p.CmdLine, "-test.") {
				t.Errorf("command line not read: %q", p.CmdLine)
			}
			return
		}
	}
	t.Fatalf("own process not found in %d processes", len(list))
}
