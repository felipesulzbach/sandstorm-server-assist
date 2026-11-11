package desktop

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"sandstorm-server-assist/desktop/internal/adapters/fsx"
	"sandstorm-server-assist/desktop/internal/app/logbuf"
	"sandstorm-server-assist/desktop/internal/app/runner"
)

// runFile guarda o registro da partida iniciada pelo app (ADR-026) ao lado do config.json. Sem segredos: a linha
// de comando já vem mascarada.
type runFile struct {
	path string
	log  *logbuf.Buffer
}

func newRunFile(configPath string, log *logbuf.Buffer) *runFile {
	return &runFile{path: filepath.Join(filepath.Dir(configPath), "server-run.json"), log: log}
}

// Save grava o registro; nil apaga. Falha só fica no log: sem o registro, o app reaberto mostra o servidor como
// iniciado fora dele (o comportamento antigo).
func (f *runFile) Save(run *runner.Run) {
	var err error
	if run == nil {
		if err = os.Remove(f.path); errors.Is(err, fs.ErrNotExist) {
			err = nil
		}
	} else {
		var b []byte
		if b, err = json.Marshal(run); err == nil {
			err = fsx.WriteAtomic(f.path, b, false)
		}
	}
	if err != nil {
		f.log.Add("app", "Could not save the server run record: "+err.Error(), "warning")
	}
}

// Load lê o registro (nil se não houver ou estiver ilegível).
func (f *runFile) Load() *runner.Run {
	b, err := fsx.ReadLimited(f.path, 64*1024)
	if err != nil {
		return nil
	}
	var run runner.Run
	if json.Unmarshal(b, &run) != nil {
		return nil
	}
	return &run
}
