package desktop

import (
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// dialogStart devolve a pasta onde o diálogo abre: a guardada, ou a pasta-mãe mais próxima que ainda existe. Vazio
// se nenhuma existe. O diálogo do Windows falha com uma pasta inicial que não existe (pasta do servidor movida ou
// apagada), e o usuário ficava sem conseguir escolher outra.
func dialogStart(dir string) string {
	if dir == "" {
		return ""
	}
	for d := filepath.Clean(dir); ; {
		if info, err := os.Stat(d); err == nil && info.IsDir() {
			return d
		}
		parent := filepath.Dir(d)
		if parent == d {
			return ""
		}
		d = parent
	}
}

// chooseDir abre o diálogo de pasta a partir de start (ver dialogStart). Se o Windows ainda recusar a pasta inicial,
// abre de novo sem ela.
func (a *app) chooseDir(start string) (string, error) {
	opts := runtime.OpenDialogOptions{DefaultDirectory: dialogStart(start), CanCreateDirectories: true}
	dir, err := runtime.OpenDirectoryDialog(a.ctx, opts)
	if err != nil && opts.DefaultDirectory != "" {
		opts.DefaultDirectory = ""
		dir, err = runtime.OpenDirectoryDialog(a.ctx, opts)
	}
	return dir, err
}
