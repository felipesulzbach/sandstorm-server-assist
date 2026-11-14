package desktop

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	"sandstorm-server-assist/desktop/internal/adapters/fsx"
	"sandstorm-server-assist/desktop/internal/app/logbuf"
	"sandstorm-server-assist/desktop/internal/app/runner"
	"sandstorm-server-assist/desktop/internal/domain/server"
)

// runFile guarda o registro da partida iniciada pelo app (ADR-026) ao lado do config.json, um por servidor
// (server-run-<id>.json, ADR-048). Sem segredos: a linha de comando já vem mascarada. legacy é o arquivo único das
// versões antigas (server-run.json), lido só pelo servidor principal, para um servidor iniciado antes da
// atualização do app continuar sendo dele.
type runFile struct {
	path   string
	legacy string
	log    *logbuf.Buffer
}

// reRunID: o id vai no nome do arquivo; só letras e números (os ids do app são hexadecimais).
var reRunID = regexp.MustCompile(`^[A-Za-z0-9]{1,32}$`)

func newRunFile(configPath, id string, legacy bool, log *logbuf.Buffer) *runFile {
	dir := filepath.Dir(configPath)
	name := "server-run.json"
	if reRunID.MatchString(id) {
		name = "server-run-" + id + ".json"
	}
	f := &runFile{path: filepath.Join(dir, name), log: log}
	if legacy && name != "server-run.json" {
		f.legacy = filepath.Join(dir, "server-run.json")
	}
	return f
}

// legacyRunFile diz se o servidor id é o principal (instância 1), o único que lê o registro antigo.
func (a *app) legacyRunFile(id string) bool {
	p, _, err := a.cfg.CurrentOf(id)
	return err == nil && server.InstanceNumber(p.Instance) == 1
}

// Save grava o registro; nil apaga. Falha só fica no log: sem o registro, o app reaberto mostra o servidor como
// iniciado fora dele (o comportamento antigo).
func (f *runFile) Save(run *runner.Run) {
	var err error
	if run == nil {
		if err = os.Remove(f.path); errors.Is(err, fs.ErrNotExist) {
			err = nil
		}
		if f.legacy != "" {
			_ = os.Remove(f.legacy)
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
	if err != nil && f.legacy != "" {
		b, err = fsx.ReadLimited(f.legacy, 64*1024)
	}
	if err != nil {
		return nil
	}
	var run runner.Run
	if json.Unmarshal(b, &run) != nil {
		return nil
	}
	return &run
}
