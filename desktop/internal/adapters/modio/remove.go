package modio

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"sandstorm-server-assist/desktop/internal/adapters/fsx"
)

// RemoveInstalled apaga um mod baixado (ADR-037): tira o mod do state.json do mod.io e apaga a pasta dele. Sem o
// registro, o servidor vê na próxima partida que o mod inscrito não está instalado e baixa de novo. Só apagar a
// pasta não serve: o state.json continuaria dizendo "instalado" e o servidor tentaria montar arquivos que não
// existem. O servidor precisa estar parado (ele regrava o state.json ao fechar).
//
// O state.json é reescrito mantendo os outros mods e campos como estavam; a primeira troca guarda uma cópia .bak.
// A pasta vem só do id (nunca do PathOnDisk do arquivo) e precisa ficar dentro da pasta de mods do mod.io.
func RemoveInstalled(root string, id int) error {
	if id <= 0 {
		return errors.New("invalid mod id")
	}
	modsDir := filepath.Join(root, "254", "mods")
	dir := ModDir(root, id)
	if rel, err := filepath.Rel(modsDir, dir); err != nil || rel == "." || strings.HasPrefix(rel, "..") {
		return errors.New("mod folder outside the mod.io folder")
	}
	found, err := dropFromState(filepath.Join(root, "254", "metadata", "state.json"), id)
	if err != nil {
		return err
	}
	_, statErr := os.Stat(dir)
	if !found && errors.Is(statErr, os.ErrNotExist) {
		return fmt.Errorf("mod %d is not installed", id)
	}
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	return nil
}

// dropFromState tira o mod da lista "Mods" do state.json. found diz se ele estava lá.
func dropFromState(path string, id int) (found bool, err error) {
	b, err := fsx.ReadLimited(path, 64<<20)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	var top map[string]json.RawMessage
	if err := json.Unmarshal(b, &top); err != nil {
		return false, err
	}
	var mods []json.RawMessage
	if raw, ok := top["Mods"]; ok {
		if err := json.Unmarshal(raw, &mods); err != nil {
			return false, err
		}
	}
	kept := make([]json.RawMessage, 0, len(mods))
	for _, m := range mods {
		var head struct {
			ID int `json:"ID"`
		}
		if err := json.Unmarshal(m, &head); err != nil {
			return false, err
		}
		if head.ID == id {
			found = true
			continue
		}
		kept = append(kept, m)
	}
	if !found {
		return false, nil
	}
	list, err := json.Marshal(kept)
	if err != nil {
		return false, err
	}
	top["Mods"] = list
	out, err := json.Marshal(top)
	if err != nil {
		return false, err
	}
	return true, fsx.WriteAtomic(path, out, true)
}
