package installer

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"sandstorm-server-assist/desktop/internal/adapters/trust"
)

// Baixa o SteamCMD de verdade (cerca de 3 MB da Valve) numa pasta temporária e confere a assinatura.
// Só roda com SSA_STEAMCMD_E2E=1 (precisa de internet).
func TestRealSteamCMDDownload(t *testing.T) {
	if os.Getenv("SSA_STEAMCMD_E2E") != "1" {
		t.Skip("defina SSA_STEAMCMD_E2E=1 para baixar o SteamCMD de verdade")
	}
	dir := filepath.Join(t.TempDir(), "steamcmd")
	s := New(nil)
	var states []State
	s.SetNotify(func(st State) { states = append(states, st) })
	if err := s.FetchSteamCMD(dir); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(2 * time.Minute)
	for s.Busy() && time.Now().Before(deadline) {
		time.Sleep(200 * time.Millisecond)
	}
	st := s.State()
	if st.Error != nil || st.Finished != FetchSteamCMD {
		t.Fatalf("estado final: %+v", st)
	}
	if err := trust.VerifySigner(SteamCMDExe(dir), "Valve"); err != nil {
		t.Fatalf("assinatura: %v", err)
	}
	sawBytes := false
	for _, x := range states {
		sawBytes = sawBytes || (x.Done > 0 && x.Total > 0)
	}
	if !sawBytes {
		t.Error("nenhum progresso em bytes foi informado")
	}
	if _, err := os.Stat(filepath.Join(dir, ".steamcmd-download.zip")); err == nil {
		t.Error("zip temporário não foi apagado")
	}
}
