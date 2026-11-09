package trust

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// Usa o notepad.exe (assinado pela Microsoft) e um arquivo sem assinatura.
func TestVerifySigner(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("só no Windows")
	}
	notepad := filepath.Join(os.Getenv("SystemRoot"), "System32", "notepad.exe")
	if _, err := os.Stat(notepad); err != nil {
		t.Skip("sem notepad.exe")
	}
	if err := VerifySigner(notepad, "Microsoft"); err != nil {
		// No Windows 11 o notepad pode vir assinado por catálogo (sem assinatura embutida).
		t.Logf("notepad: %v", err)
	}
	if err := VerifySigner(notepad, "Valve"); err == nil {
		t.Error("signatário errado aceito")
	}
	unsigned := filepath.Join(t.TempDir(), "x.exe")
	if err := os.WriteFile(unsigned, []byte("MZ"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := VerifySigner(unsigned, "Valve"); !errors.Is(err, ErrUntrusted) {
		t.Errorf("sem assinatura: %v", err)
	}
	// SteamCMD real, se existir nesta máquina (D:\_server\tools\steamcmd).
	if p := os.Getenv("SSA_STEAMCMD"); p != "" {
		if err := VerifySigner(p, "Valve"); err != nil {
			t.Errorf("steamcmd: %v", err)
		}
	}
}
