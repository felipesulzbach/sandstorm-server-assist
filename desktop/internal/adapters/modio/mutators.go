package modio

import (
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"sandstorm-server-assist/desktop/internal/domain/pak"
)

// ModMutator é um mutator achado num mod instalado (ADR-023).
type ModMutator struct {
	ID      string `json:"id"`
	ModID   int    `json:"modId"`
	ModName string `json:"modName"`
}

// ModScenario é um cenário de mapa achado num mod instalado (para o mapa de mod depois de carregar, ADR-024).
type ModScenario struct {
	Map      string `json:"map"`
	Scenario string `json:"scenario"`
	ModID    int    `json:"modId"`
	ModName  string `json:"modName"`
}

// maxPaks limita quantos .pak de um mod são lidos.
const maxPaks = 32

// installedMod é um mod instalado com o que os .pak dele trazem.
type installedMod struct {
	mod      Installed
	contents pak.Contents
	readable bool
}

// scanMods lê o índice dos .pak de cada mod instalado (só o rodapé e os índices, nunca o arquivo inteiro).
func scanMods(root string) []installedMod {
	list, err := ReadInstalled(root)
	if err != nil {
		return nil
	}
	modsDir := filepath.Join(root, "254", "mods")
	var out []installedMod
	for _, m := range list {
		if m.State != StateInstalled {
			continue
		}
		// A pasta vem do id do mod: só vale dentro da pasta de mods do mod.io.
		dir := ModDir(root, m.ID)
		if rel, err := filepath.Rel(modsDir, dir); err != nil || strings.HasPrefix(rel, "..") {
			continue
		}
		c, readable := scanDir(dir)
		out = append(out, installedMod{mod: m, contents: c, readable: readable})
	}
	return out
}

// Mutators lista os mutators dos mods instalados (ADR-023). Se nenhum .pak do mod puder ser lido e o mod tem a tag
// "Mutator" no mod.io, usa o nome do plugin.
func Mutators(root string) []ModMutator {
	out := []ModMutator{}
	seen := map[string]bool{}
	for _, im := range scanMods(root) {
		names := im.contents.Mutators
		if !im.readable && hasTag(im.mod.Tags, "Mutator") {
			if n := pak.PluginName(im.mod.Blob); n != "" {
				names = []string{n}
			}
		}
		for _, n := range names {
			if !seen[n] {
				seen[n] = true
				out = append(out, ModMutator{ID: n, ModID: im.mod.ID, ModName: im.mod.Name})
			}
		}
	}
	sort.Slice(out, func(i, j int) bool { return strings.ToLower(out[i].ID) < strings.ToLower(out[j].ID) })
	return out
}

// Scenarios lista os cenários dos mods instalados com o mapa de cada um (cenário sem mapa no mesmo mod fica de
// fora: o servidor não teria para onde viajar).
func Scenarios(root string) []ModScenario {
	out := []ModScenario{}
	seen := map[string]bool{}
	for _, im := range scanMods(root) {
		for _, sc := range im.contents.Scenarios {
			m := pak.ScenarioMap(sc, im.contents.Maps)
			if m == "" || seen[sc] {
				continue
			}
			seen[sc] = true
			out = append(out, ModScenario{Map: m, Scenario: sc, ModID: im.mod.ID, ModName: im.mod.Name})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Map != out[j].Map {
			return strings.ToLower(out[i].Map) < strings.ToLower(out[j].Map)
		}
		return out[i].Scenario < out[j].Scenario
	})
	return out
}

func hasTag(tags []string, tag string) bool {
	for _, t := range tags {
		if strings.EqualFold(t, tag) {
			return true
		}
	}
	return false
}

// scanDir lê os .pak da pasta de um mod e junta o conteúdo. readable diz se pelo menos um índice foi lido.
func scanDir(dir string) (all pak.Contents, readable bool) {
	var paks []string
	_ = filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
		if err != nil || len(paks) >= maxPaks {
			return filepath.SkipDir
		}
		if !d.IsDir() && strings.EqualFold(filepath.Ext(p), ".pak") && d.Type().IsRegular() {
			paks = append(paks, p)
		}
		return nil
	})
	for _, p := range paks {
		c, err := scanPak(p)
		if err == nil {
			readable = true
			all.Mutators = append(all.Mutators, c.Mutators...)
			all.Maps = append(all.Maps, c.Maps...)
			all.Scenarios = append(all.Scenarios, c.Scenarios...)
		}
	}
	return all, readable
}

// scanPak lê o que interessa de um .pak: no formato novo (versões 10 e 11), pelo índice de pastas; no antigo (8 e
// 9), pela lista de arquivos.
func scanPak(path string) (pak.Contents, error) {
	f, err := os.Open(path) // #nosec G304 -- .pak dentro da pasta de mods do mod.io ou do servidor
	if err != nil {
		return pak.Contents{}, err
	}
	defer f.Close()
	p, err := openPak(f)
	if err != nil {
		return pak.Contents{}, err
	}
	if p.legacy != nil {
		return pak.ContentsOf(p.legacy), nil
	}
	return pak.Scan(p.dirIndex)
}

// pakIndex são os índices lidos de um .pak: dirIndex e primary no formato novo, legacy (a lista de arquivos) no
// antigo.
type pakIndex struct {
	primary  []byte
	dirIndex []byte
	legacy   []pak.File
}

// openPak lê só o rodapé e os índices (nunca o arquivo inteiro).
func openPak(f *os.File) (pakIndex, error) {
	info, err := f.Stat()
	if err != nil {
		return pakIndex{}, err
	}
	size := info.Size()
	read := func(off, n int64) ([]byte, error) {
		b := make([]byte, n)
		_, err := f.ReadAt(b, off)
		if err != nil && err != io.EOF {
			return nil, err
		}
		return b, nil
	}
	tailLen := min(int64(pak.TailSize), size)
	tail, err := read(size-tailLen, tailLen)
	if err != nil {
		return pakIndex{}, err
	}
	ft, err := pak.FooterAny(tail, size)
	if err != nil {
		return pakIndex{}, err
	}
	primary, err := read(ft.IndexOffset, ft.IndexSize)
	if err != nil {
		return pakIndex{}, err
	}
	if ft.Version < 10 {
		files, err := pak.LegacyFiles(primary)
		return pakIndex{legacy: files}, err
	}
	off, n, err := pak.DirIndex(primary, size)
	if err != nil {
		return pakIndex{}, err
	}
	dirIndex, err := read(off, n)
	if err != nil {
		return pakIndex{}, err
	}
	return pakIndex{primary: primary, dirIndex: dirIndex}, nil
}

// pakFiles lista os arquivos de um .pak com o tamanho de cada um.
func pakFiles(path string) ([]pak.File, error) {
	f, err := os.Open(path) // #nosec G304 -- .pak dentro da pasta de mods do mod.io
	if err != nil {
		return nil, err
	}
	defer f.Close()
	p, err := openPak(f)
	if err != nil {
		return nil, err
	}
	if p.legacy != nil {
		return p.legacy, nil
	}
	return pak.ModernFiles(p.primary, p.dirIndex)
}

// BaseScenarios lista os cenários do jogo base, lidos dos .pak do servidor (pasta Insurgency/Content/Paks, só os
// índices). Serve para saber se um mapa da rotação existe (ADR-030): a lista fixa de mapas do app não tem todos.
func BaseScenarios(paksDir string) []string {
	c, _ := scanDir(paksDir)
	return c.Scenarios
}

// ScenarioIDs lista todos os cenários dos mods instalados, mesmo os sem mapa identificado (para conferir se um
// cenário da rotação existe, ADR-030).
func ScenarioIDs(root string) []string {
	out := []string{}
	for _, im := range scanMods(root) {
		out = append(out, im.contents.Scenarios...)
	}
	return out
}
