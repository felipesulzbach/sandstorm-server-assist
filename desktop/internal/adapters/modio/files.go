package modio

import (
	"os"
	"path/filepath"
	"strings"

	"sandstorm-server-assist/desktop/internal/domain/modcheck"
)

// ModFiles é o resultado da leitura dos arquivos dos mods instalados (ADR-039).
type ModFiles struct {
	Mods []modcheck.Mod
	// Unreadable são os mods com algum .pak que o app não conseguiu ler (formato desconhecido ou índice cifrado):
	// esses mods podem ter conflitos que o app não vê.
	Unreadable []modcheck.ModRef
}

// InstalledFiles lê os índices dos .pak de cada mod instalado (só rodapé e índices) e lista os arquivos com o
// tamanho de cada um.
func InstalledFiles(root string) ModFiles {
	out := ModFiles{Mods: []modcheck.Mod{}, Unreadable: []modcheck.ModRef{}}
	list, err := ReadInstalled(root)
	if err != nil {
		return out
	}
	modsDir := filepath.Join(root, "254", "mods")
	for _, m := range list {
		if m.State != StateInstalled {
			continue
		}
		dir := ModDir(root, m.ID)
		if rel, err := filepath.Rel(modsDir, dir); err != nil || strings.HasPrefix(rel, "..") {
			continue
		}
		mod := modcheck.Mod{ID: m.ID, Name: m.Name}
		bad := false
		count := 0
		_ = filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
			if err != nil || count >= maxPaks {
				return filepath.SkipDir
			}
			if d.IsDir() || !strings.EqualFold(filepath.Ext(p), ".pak") || !d.Type().IsRegular() {
				return nil
			}
			count++
			files, err := pakFiles(p)
			if err != nil {
				bad = true
				return nil
			}
			for _, f := range files {
				mod.Files = append(mod.Files, modcheck.File{Path: f.Path, Size: f.Size})
			}
			return nil
		})
		out.Mods = append(out.Mods, mod)
		if bad {
			out.Unreadable = append(out.Unreadable, modcheck.ModRef{ID: m.ID, Name: m.Name})
		}
	}
	return out
}
