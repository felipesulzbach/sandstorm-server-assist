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

// maxPaks limita quantos .pak de um mod são lidos.
const maxPaks = 32

// Mutators lista os mutators dos mods instalados: lê o índice de cada .pak (só o rodapé e os índices, nunca o
// arquivo inteiro). Se nenhum .pak puder ser lido e o mod tem a tag "Mutator" no mod.io, usa o nome do plugin.
func Mutators(root string) []ModMutator {
	list, err := ReadInstalled(root)
	if err != nil {
		return nil
	}
	modsDir := filepath.Join(root, "254", "mods")
	var out []ModMutator
	seen := map[string]bool{}
	for _, m := range list {
		if m.State != StateInstalled {
			continue
		}
		// A pasta vem do state.json: só vale dentro da pasta de mods do mod.io.
		dir := ModDir(root, m.ID)
		if rel, err := filepath.Rel(modsDir, dir); err != nil || strings.HasPrefix(rel, "..") {
			continue
		}
		names, readable := scanDir(dir)
		if !readable && hasTag(m.Tags, "Mutator") {
			if n := pak.PluginName(m.Blob); n != "" {
				names = []string{n}
			}
		}
		for _, n := range names {
			if !seen[n] {
				seen[n] = true
				out = append(out, ModMutator{ID: n, ModID: m.ID, ModName: m.Name})
			}
		}
	}
	sort.Slice(out, func(i, j int) bool { return strings.ToLower(out[i].ID) < strings.ToLower(out[j].ID) })
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

// scanDir lê os .pak da pasta de um mod. readable diz se pelo menos um índice foi lido.
func scanDir(dir string) (names []string, readable bool) {
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
		n, err := scanPak(p)
		if err == nil {
			readable = true
			names = append(names, n...)
		}
	}
	return names, readable
}

func scanPak(path string) ([]string, error) {
	f, err := os.Open(path) // #nosec G304 -- .pak dentro da pasta de mods do mod.io
	if err != nil {
		return nil, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return nil, err
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
		return nil, err
	}
	ft, err := pak.ParseFooter(tail, size)
	if err != nil {
		return nil, err
	}
	primary, err := read(ft.IndexOffset, ft.IndexSize)
	if err != nil {
		return nil, err
	}
	off, n, err := pak.DirIndex(primary, size)
	if err != nil {
		return nil, err
	}
	dirIndex, err := read(off, n)
	if err != nil {
		return nil, err
	}
	return pak.Mutators(dirIndex)
}
