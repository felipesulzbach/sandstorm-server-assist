package backupsvc

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"sandstorm-server-assist/desktop/internal/app/config"
	"sandstorm-server-assist/desktop/internal/domain/server"
)

type memVault struct{ s server.Vault }

func (m *memVault) Load() (server.Vault, error) { return m.s, nil }
func (m *memVault) Save(s server.Vault) error   { m.s = s; return nil }

// gameIni é o Game.ini do servidor principal.
const gameIni = `Insurgency\Saved\Config\WindowsServer\Game.ini`

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
	write(t, filepath.Join(serverA, gameIni), "[Rcon]\r\nbEnabled=True\r\nPassword=rcon#Senha1\r\nListenPort=27015\r\n")
	write(t, filepath.Join(serverA, `Insurgency\Config\Server\Admins.txt`), ";Alice\r\n76561198000000001\r\n")
	write(t, filepath.Join(serverA, `Insurgency\Config\Server\MapCycle.txt`), "//x\r\n(Scenario=\"Scenario_Farmhouse_Checkpoint_Security\")\r\n")

	vaultA := &memVault{}
	cfgA := config.Open(filepath.Join(a, "config.json"), vaultA)
	if _, err := cfgA.SetServerDir(serverA); err != nil {
		t.Fatal(err)
	}
	r, err := cfgA.ApplyExisting("", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	if p := r.Profile; len(p.Admins) != 1 || len(p.MapCycle) != 1 || !p.Rcon || !r.RconInIni || one(vaultA).RconPassword != "rcon#Senha1" {
		t.Fatalf("import %+v vault %+v", r, vaultA.s)
	}
	if err := cfgA.MoveRconToVault(); err != nil {
		t.Fatal(err)
	}
	if cfgA.RconInIni() {
		t.Fatal("senha continua no Game.ini")
	}
	if _, err := os.Stat(filepath.Join(serverA, gameIni+".bak")); err != nil {
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
	if one(vaultB).GSLTToken != "ABCDEF0123456789" || one(vaultB).RconPassword != "rcon#Senha1" {
		t.Errorf("cofre B: %+v", vaultB.s)
	}
	if b, err := os.ReadFile(filepath.Join(serverB, `Insurgency\Config\Server\Admins.txt`)); err != nil || !strings.Contains(string(b), "76561198000000001") {
		t.Errorf("Admins.txt: %v %q", err, b)
	}
}

// one devolve os segredos do único servidor do cofre, com os da conta.
func one(m *memVault) server.Secrets {
	for _, sec := range m.s.Servers {
		return sec.WithAccount(m.s.Secrets)
	}
	return m.s.Secrets
}
