package desktop

import (
	"path/filepath"
	"testing"

	"sandstorm-server-assist/desktop/internal/app/config"
	"sandstorm-server-assist/desktop/internal/app/installer"
	"sandstorm-server-assist/desktop/internal/domain/problems"
)

// Iniciar os marcados (ADR-048): só ids de servidores que existem; um que não inicia não para os outros.
func TestStartServers(t *testing.T) {
	cfg := config.Open(filepath.Join(t.TempDir(), "config.json"), memoryVault{})
	if _, _, err := cfg.CreateProfile("Versus", false, nil); err != nil {
		t.Fatal(err)
	}
	a := &app{cfg: cfg, servers: &fleet{byID: map[string]*srv{}}, mods: &modCache{}}
	a.install = installer.New(fanLog{a})
	api := &ServerAPI{a: a}
	refs := cfg.Profiles()
	if err := api.StartServers(nil); problems.CodeOf(err) != problems.InvalidField {
		t.Errorf("empty list: %v", err)
	}
	if err := api.StartServers([]string{refs[0].ID, "unknown"}); problems.CodeOf(err) != problems.InvalidField {
		t.Errorf("unknown id: %v", err)
	}
	// Sem pasta do servidor nenhum inicia: volta o erro, e cada um tentou (fica no log de cada um).
	err := api.StartServers([]string{refs[1].ID, refs[0].ID})
	if problems.CodeOf(err) != problems.ServerNotInstalled {
		t.Fatalf("not installed: %v", err)
	}
	if err := api.RestartServers([]string{"unknown"}); problems.CodeOf(err) != problems.InvalidField {
		t.Errorf("restart unknown id: %v", err)
	}
	// Nenhum roda: reiniciar os marcados não inicia nenhum (parados ficam parados).
	if err := api.RestartServers([]string{refs[0].ID, refs[1].ID}); err != nil {
		t.Errorf("restart stopped servers: %v", err)
	}
	if err := api.StopServers([]string{"unknown"}); problems.CodeOf(err) != problems.InvalidField {
		t.Errorf("stop unknown id: %v", err)
	}
	// Nenhum roda: parar os marcados não faz nada e não dá erro.
	if err := api.StopServers([]string{refs[0].ID, refs[1].ID}); err != nil {
		t.Errorf("stop stopped servers: %v", err)
	}
	for _, r := range refs {
		lines := a.server(r.ID).log.Recent(10)
		if len(lines) == 0 || lines[len(lines)-1].Level != "warning" {
			t.Errorf("%s: no failure logged: %+v", r.Name, lines)
		}
	}
}
