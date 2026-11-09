package backupsvc

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"sandstorm-server-assist/desktop/internal/app/config"
	"sandstorm-server-assist/desktop/internal/domain/server"
)

type memVault struct{ s server.Secrets }

func (m *memVault) Load() (server.Secrets, error) { return m.s, nil }
func (m *memVault) Save(s server.Secrets) error   { m.s = s; return nil }

func write(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

// Máquina A: importa um servidor existente, move a senha de RCON para o cofre, faz backup.
// Máquina B: restaura numa pasta nova e confere perfil, segredos e arquivos.
func TestImportBackupRestore(t *testing.T) {
	a := t.TempDir()
	serverA := filepath.Join(a, "sandstorm_server")
	write(t, filepath.Join(serverA, config.GameIni), "[Rcon]\r\nbEnabled=True\r\nPassword=rcon#Senha1\r\nListenPort=27015\r\n")
	write(t, filepath.Join(serverA, config.AdminsTxt), ";Alice\r\n76561198000000001\r\n")
	write(t, filepath.Join(serverA, config.MapCycleTxt), "//x\r\n(Scenario=\"Scenario_Farmhouse_Checkpoint_Security\")\r\n")

	vaultA := &memVault{}
	cfgA := config.Open(filepath.Join(a, "config.json"), vaultA)
	if _, err := cfgA.SetServerDir(serverA); err != nil {
		t.Fatal(err)
	}
	r, err := cfgA.ApplyExisting("", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	if p := r.Profile; len(p.Admins) != 1 || len(p.MapCycle) != 1 || !p.Rcon || !r.RconInIni || vaultA.s.RconPassword != "rcon#Senha1" {
		t.Fatalf("import %+v vault %+v", r, vaultA.s)
	}
	if err := cfgA.MoveRconToVault(); err != nil {
		t.Fatal(err)
	}
	if cfgA.RconInIni() {
		t.Fatal("senha continua no Game.ini")
	}
	if _, err := os.Stat(filepath.Join(serverA, config.GameIni+".bak")); err != nil {
		t.Fatal("sem cópia .bak do Game.ini")
	}
	if err := cfgA.SetSecret("gsltToken", "ABCDEF0123456789"); err != nil {
		t.Fatal(err)
	}

	file := filepath.Join(a, "b"+".ssbackup")
	if _, err := Create(cfgA, "0.1.0", file, "curta", Options{}); err == nil {
		t.Fatal("senha curta aceita")
	}
	sum, err := Create(cfgA, "0.1.0", file, "senha-do-backup-1", Options{})
	if err != nil {
		t.Fatal(err)
	}
	if len(sum.Files) < 3 {
		t.Errorf("arquivos: %v", sum.Files)
	}
	raw, _ := os.ReadFile(file)
	if strings.Contains(string(raw), "ABCDEF0123456789") || strings.Contains(string(raw), "Alice") {
		t.Fatal("conteúdo legível no backup")
	}

	b := t.TempDir()
	serverB := filepath.Join(b, "novo_servidor")
	vaultB := &memVault{}
	cfgB := config.Open(filepath.Join(b, "config.json"), vaultB)
	if _, err := Inspect(file, "senha-errada-99"); err == nil {
		t.Fatal("senha errada abriu")
	}
	if _, err := Restore(cfgB, file, "senha-do-backup-1", serverB, false); err != nil {
		t.Fatal(err)
	}
	p := cfgB.Profile()
	if p.ServerDir != serverB || len(p.Admins) != 1 || !p.Rcon {
		t.Errorf("perfil: %+v", p)
	}
	if vaultB.s.GSLTToken != "ABCDEF0123456789" || vaultB.s.RconPassword != "rcon#Senha1" {
		t.Errorf("cofre B: %+v", vaultB.s)
	}
	if b, err := os.ReadFile(filepath.Join(serverB, config.AdminsTxt)); err != nil || !strings.Contains(string(b), "76561198000000001") {
		t.Errorf("Admins.txt: %v %q", err, b)
	}
}
