package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"sandstorm-server-assist/desktop/internal/domain/server"
)

// Configuração da versão 0.2.9: dois perfis sem instância e o cofre antigo (tudo no nível de cima).
func TestMigrateProfilesToServers(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	old := `{"version":1,"profile":{"hostname":"Coop","port":27102,"queryPort":27131,"maxPlayers":8,"startMap":"Farmhouse",` +
		`"startScenario":"Scenario_Farmhouse_Checkpoint_Security","rconPort":27015},"activeId":"aaa","activeName":"Coop",` +
		`"others":[{"id":"bbb","name":"Versus","profile":{"hostname":"Versus","port":27102,"queryPort":27131,"maxPlayers":8,` +
		`"startMap":"Farmhouse","startScenario":"Scenario_Farmhouse_Checkpoint_Security","rconPort":27015}}],"settings":{}}`
	if err := os.WriteFile(path, []byte(old), 0o600); err != nil {
		t.Fatal(err)
	}
	v := &memVault{s: server.Vault{Secrets: server.Secrets{RconPassword: "rcon#Senha1", GSLTToken: "token-1234",
		ModioEmail: "a@example.com"}}}
	s := Open(path, v)
	refs := s.Profiles()
	if len(refs) != 2 || refs[0].ID != "aaa" || refs[0].Instance != 1 || refs[1].ID != "bbb" || refs[1].Instance != 2 {
		t.Fatalf("servers: %+v", refs)
	}
	// Os segredos de antes valiam para todos os perfis: cada servidor fica com uma cópia; a conta fica em cima.
	if len(v.s.Servers) != 2 || v.s.Servers["bbb"].RconPassword != "rcon#Senha1" || v.s.Servers["aaa"].ModioEmail != "" ||
		v.s.ModioEmail != "a@example.com" || v.s.Version != server.VaultVersion {
		t.Fatalf("vault: %+v", v.s)
	}
	// Uma versão antiga do app lê só o nível de cima: ele leva os segredos do servidor escolhido.
	if v.s.RconPassword != "rcon#Senha1" || v.s.GSLTToken != "token-1234" {
		t.Fatalf("legacy copy: %+v", v.s.Secrets)
	}
	_, sec, err := s.CurrentOf("bbb")
	if err != nil || sec.GSLTToken != "token-1234" || sec.ModioEmail != "a@example.com" {
		t.Fatalf("CurrentOf: %+v %v", sec, err)
	}
	// A instância fica gravada: reabrir não troca os números.
	if again := Open(path, v).Profiles(); again[1].Instance != 2 {
		t.Fatalf("reopened: %+v", again)
	}
	b, _ := os.ReadFile(path)
	if !strings.Contains(string(b), `"instance": 2`) || !strings.Contains(string(b), `"others"`) {
		t.Errorf("file kept in the old shape with instances: %s", b)
	}
}

func TestServersSecretsPortsAndFolders(t *testing.T) {
	root := t.TempDir()
	install := filepath.Join(root, "sandstorm")
	moved := filepath.Join(root, "moved")
	for _, d := range []string{install, moved} {
		if err := os.MkdirAll(d, 0o750); err != nil {
			t.Fatal(err)
		}
	}
	v := &memVault{}
	s := Open(filepath.Join(root, "config.json"), v)
	if _, err := s.SetServerDir(install); err != nil {
		t.Fatal(err)
	}
	p := s.Profile()
	p.Rcon = true
	if _, err := s.SaveProfile(p); err != nil {
		t.Fatal(err)
	}
	for k, val := range map[string]string{"rconPassword": "rcon#Senha1", "gsltToken": "token-main1", "modioEmail": "a@example.com"} {
		if err := s.SetSecret(k, val); err != nil {
			t.Fatal(err)
		}
	}
	mainID := s.SelectedID()

	// Cópia: mesmas pastas, instância e portas novas, senhas copiadas, sem GSLT.
	_, copyID, err := s.CreateProfile("Copy", true, nil)
	if err != nil {
		t.Fatal(err)
	}
	cp, csec, _ := s.CurrentOf(copyID)
	if cp.ServerDir != install || cp.Instance != 2 || cp.Port != 27103 || cp.QueryPort != 27132 || cp.RconPort != 27016 {
		t.Fatalf("copy: %+v", cp)
	}
	if csec.RconPassword != "rcon#Senha1" || csec.GSLTToken != "" || csec.ModioEmail != "a@example.com" {
		t.Fatalf("copy secrets: %+v", csec)
	}
	if len(server.PortConflicts(cp, s.Peers(copyID))) != 0 {
		t.Fatal("new server ports conflict")
	}

	// Novo padrão: terceira instância e portas.
	_, newID, _ := s.CreateProfile("Fresh", false, nil)
	np, nsec, _ := s.CurrentOf(newID)
	if np.Instance != 3 || np.Port != 27104 || np.ServerDir != install || nsec.RconPassword != "" {
		t.Fatalf("default: %+v %+v", np, nsec)
	}

	// Segredo do servidor só no escolhido; o da conta em todos.
	if _, err := s.SwitchProfile(copyID); err != nil {
		t.Fatal(err)
	}
	if err := s.SetSecret("gsltToken", "token-copy2"); err != nil {
		t.Fatal(err)
	}
	if err := s.SetSecret("modioEmail", "b@example.com"); err != nil {
		t.Fatal(err)
	}
	_, msec, _ := s.CurrentOf(mainID)
	if msec.GSLTToken != "token-main1" || msec.ModioEmail != "b@example.com" {
		t.Fatalf("main secrets: %+v", msec)
	}
	if v.s.GSLTToken != "token-copy2" {
		t.Fatalf("legacy copy is not the selected server's: %+v", v.s.Secrets)
	}
	// Reaberto, o cofre já migrado não copia a cópia antiga para os outros servidores.
	again := Open(filepath.Join(root, "config.json"), v)
	if _, sec, _ := again.CurrentOf(newID); sec.GSLTToken != "" || sec.RconPassword != "" {
		t.Fatalf("migrated twice: %+v", sec)
	}
	// Escolher não muda a instância do perfil salvo.
	cur := s.Profile()
	cur.Instance = 1
	if saved, _ := s.SaveProfile(cur); saved.Instance != 2 {
		t.Fatalf("SaveProfile changed the instance: %d", saved.Instance)
	}

	// Trocar a pasta da instalação leva todos os servidores que a usavam.
	if _, err := s.SetServerDir(moved); err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{mainID, copyID, newID} {
		if p, _, _ := s.CurrentOf(id); p.ServerDir != moved {
			t.Errorf("%s not moved: %s", id, p.ServerDir)
		}
	}

	// Apagar leva os segredos do servidor.
	if _, err := s.DeleteProfile(mainID); err != nil {
		t.Fatal(err)
	}
	if _, ok := v.s.Servers[mainID]; ok {
		t.Error("deleted server secrets kept")
	}
	for _, val := range s.SecretValues() {
		if val == "token-main1" {
			t.Error("deleted server secret still masked/kept")
		}
	}
}

func TestPrepareInstanceFiles(t *testing.T) {
	dir := t.TempDir()
	p := server.Default()
	p.ServerDir, p.Instance = dir, 2
	if err := Prepare(p); err != nil {
		t.Fatal(err)
	}
	// Instância 2: Admins_2 e Motd_2 sempre (a linha de comando aponta para eles), sem tocar nos da principal.
	for _, f := range []string{`Insurgency\Config\Server\Admins_2.txt`, `Insurgency\Config\Server\Motd_2.txt`} {
		if _, err := os.Stat(filepath.Join(dir, f)); err != nil {
			t.Errorf("%s: %v", f, err)
		}
	}
	for _, f := range []string{AdminsTxt, MotdTxt, MapCycleTxt, `Insurgency\Config\Server\MapCycle_2.txt`} {
		if _, err := os.Stat(filepath.Join(dir, f)); err == nil {
			t.Errorf("%s written", f)
		}
	}
	p.MapCycle = []string{server.CycleLine("Scenario_Hideout_Checkpoint_Security", server.LightingNight)}
	p.Rules = map[string]string{}
	if err := Prepare(p); err != nil {
		t.Fatal(err)
	}
	if b, err := os.ReadFile(filepath.Join(dir, `Insurgency\Config\Server\MapCycle_2.txt`)); err != nil || !strings.Contains(string(b), "Hideout") {
		t.Errorf("MapCycle_2: %q %v", b, err)
	}
	if LogsPath(p) != filepath.Join(dir, `Insurgency\Saved_2\Logs`) {
		t.Errorf("logs: %s", LogsPath(p))
	}
}

// Cofre gravado no meio do caminho (servidores já separados, sem a versão): os segredos de cada servidor ficam.
func TestMigrateVaultKeepsServerSecrets(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	first := Open(path, &memVault{})
	id := first.SelectedID()
	v := &memVault{s: server.Vault{Secrets: server.Secrets{ModioEmail: "a@example.com"},
		Servers: map[string]server.Secrets{id: {GSLTToken: "token-keep1", RconPassword: "rcon#Senha1"}}}}
	s := Open(path, v)
	if _, sec, _ := s.CurrentOf(id); sec.GSLTToken != "token-keep1" || sec.RconPassword != "rcon#Senha1" {
		t.Fatalf("server secrets lost: %+v", sec)
	}
	if v.s.Version != server.VaultVersion {
		t.Errorf("version not set: %d", v.s.Version)
	}
}
