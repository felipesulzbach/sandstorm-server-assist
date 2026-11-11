package desktop

import (
	"path/filepath"
	"sync"
	"time"

	"sandstorm-server-assist/desktop/internal/adapters/modio"
	"sandstorm-server-assist/desktop/internal/app/runner"
)

// loadedModScenarios devolve os cenários dos mods carregados nesta partida, lidos dos .pak (ADR-029). Sem nenhum
// mod ativo no log, devolve nil: a troca de mapa usa os cenários do perfil, como antes.
func loadedModScenarios(mods []runner.ModState) []string {
	active := map[int]bool{}
	for _, m := range mods {
		if m.State == "active" {
			active[m.ID] = true
		}
	}
	if len(active) == 0 {
		return nil
	}
	out := []string{}
	for _, sc := range modio.Scenarios(modio.RootDir()) {
		if active[sc.ModID] {
			out = append(out, sc.Scenario)
		}
	}
	return out
}

// scenarioCache guarda os cenários que existem nesta máquina: os do jogo base (lidos dos .pak do servidor) e os
// dos mods instalados (ADR-030). A tela da rotação pergunta a cada mudança; a leitura vale por 15 s.
type scenarioCache struct {
	mu   sync.Mutex
	at   time.Time
	dir  string
	list []string
}

const scenarioCacheFor = 15 * time.Second

// known devolve os cenários que existem. Sem o jogo base lido (servidor não instalado), devolve nil: sem saber
// quais são os mapas oficiais, nada pode ser marcado como inexistente.
func (c *scenarioCache) known(serverDir string) []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	if serverDir == c.dir && time.Since(c.at) < scenarioCacheFor {
		return c.list
	}
	c.dir, c.at, c.list = serverDir, time.Now(), nil
	if serverDir == "" {
		return nil
	}
	base := modio.BaseScenarios(filepath.Join(serverDir, "Insurgency", "Content", "Paks"))
	if len(base) == 0 {
		return nil
	}
	c.list = append(append([]string{}, base...), modio.ScenarioIDs(modio.RootDir())...)
	return c.list
}
