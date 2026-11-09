package vault

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"sandstorm-server-assist/desktop/internal/domain/server"
)

func TestDPAPIRoundTrip(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("DPAPI só no Windows")
	}
	v := Vault{Path: filepath.Join(t.TempDir(), "vault.bin")}
	empty, err := v.Load()
	if err != nil || empty != (server.Secrets{}) {
		t.Fatalf("cofre vazio: %+v %v", empty, err)
	}
	in := server.Secrets{RconPassword: "rcon#Senha1", GSLTToken: "ABCDEF0123456789"}
	if err := v.Save(in); err != nil {
		t.Fatal(err)
	}
	raw, _ := os.ReadFile(v.Path)
	if bytes.Contains(raw, []byte("rcon#Senha1")) {
		t.Fatal("segredo em texto puro")
	}
	out, err := v.Load()
	if err != nil || out != in {
		t.Fatalf("%+v %v", out, err)
	}
	// Arquivo alterado não abre (e não derruba o app: vira erro do catálogo).
	raw[len(raw)-1] ^= 0xff
	_ = os.WriteFile(v.Path, raw, 0o600)
	if _, err := v.Load(); err == nil {
		t.Error("cofre alterado abriu")
	}
}
