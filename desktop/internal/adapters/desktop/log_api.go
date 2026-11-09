package desktop

import (
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"sandstorm-server-assist/desktop/internal/adapters/fsx"
	"sandstorm-server-assist/desktop/internal/app/config"
	"sandstorm-server-assist/desktop/internal/app/logbuf"
	"sandstorm-server-assist/desktop/internal/app/runner"
	"sandstorm-server-assist/desktop/internal/domain/problems"
	"sandstorm-server-assist/desktop/internal/domain/serverlogs"
)

// LogAPI: linhas recentes (as novas chegam pelo evento "log:lines"), partidas anteriores lidas dos arquivos do
// servidor, limpeza dos logs antigos (ADR-018) e abrir a pasta de logs.
type LogAPI struct{ a *app }

// Recent devolve as últimas linhas (já mascaradas).
func (l *LogAPI) Recent(limit int) []logbuf.Line {
	if limit <= 0 || limit > 6000 {
		limit = 2000
	}
	return l.a.log.Recent(limit)
}

// Clear limpa a lista da tela (não apaga os arquivos de log do servidor).
func (l *LogAPI) Clear() { l.a.log.Clear() }

// OpenFolder abre a pasta de logs do servidor no Explorador de Arquivos.
func (l *LogAPI) OpenFolder() error {
	p := l.a.cfg.Profile()
	if p.ServerDir == "" {
		return problems.New(problems.ServerNotInstalled, "")
	}
	dir := filepath.Dir(runner.LogPath(p.ServerDir))
	if !fsx.Exists(dir) {
		return problems.New(problems.ServerNotInstalled, dir)
	}
	// Programa fixo do Windows; a pasta vem do perfil (não da tela).
	return exec.Command("explorer.exe", dir).Start() // #nosec G204 -- explorer.exe com a pasta de logs do perfil
}

// LogRun é uma partida guardada pelo servidor (um arquivo de log).
type LogRun struct {
	ID      string    `json:"id"` // nome do arquivo (a tela só devolve este id)
	Started time.Time `json:"started"`
	Ended   time.Time `json:"ended"`
	Size    int64     `json:"size"`
	Current bool      `json:"current"`
}

// reFirstTime acha o horário (UTC) da primeira linha com data do log: [2026.11.09-22.56.07:123].
var reFirstTime = regexp.MustCompile(`\[(\d{4}\.\d{2}\.\d{2}-\d{2}\.\d{2}\.\d{2}):\d{3}\]`)

// Runs lista as partidas que o servidor guardou, da mais nova para a mais velha.
func (l *LogAPI) Runs() []LogRun {
	dir := l.a.cfg.Profile().ServerDir
	files := config.ServerLogs(dir)
	out := make([]LogRun, 0, len(files))
	for _, f := range files {
		r := LogRun{ID: f.Name, Ended: f.ModTime, Size: f.Size, Current: f.Name == serverlogs.Current, Started: f.ModTime}
		if head, err := fsx.ReadHead(filepath.Join(dir, config.LogsDir, f.Name), 256<<10); err == nil {
			if m := reFirstTime.FindSubmatch(head); m != nil {
				if t, err := time.ParseInLocation("2006.01.02-15.04.05", string(m[1]), time.UTC); err == nil {
					r.Started = t.Local()
				}
			}
		}
		out = append(out, r)
	}
	return out
}

// maxRunLines é quanto de uma partida antiga a tela recebe (as últimas linhas).
const maxRunLines = 6000

// ReadRun devolve as últimas linhas de uma partida guardada, já mascaradas e classificadas. id precisa ser um
// dos nomes de Runs (a tela nunca manda caminho).
func (l *LogAPI) ReadRun(id string) ([]logbuf.Line, error) {
	dir := l.a.cfg.Profile().ServerDir
	found := false
	for _, f := range config.ServerLogs(dir) {
		found = found || f.Name == id
	}
	if !found || !serverlogs.IsServerLog(id) {
		return nil, problems.New(problems.Unknown, "unknown log")
	}
	b, err := fsx.ReadLimited(filepath.Join(dir, config.LogsDir, id), 128<<20)
	if err != nil {
		return nil, problems.Wrap(problems.Unknown, err)
	}
	lines := strings.Split(strings.ReplaceAll(string(b), "\r\n", "\n"), "\n")
	if len(lines) > maxRunLines {
		lines = lines[len(lines)-maxRunLines:]
	}
	out := make([]logbuf.Line, 0, len(lines))
	for i, text := range lines {
		if strings.TrimSpace(text) == "" {
			continue
		}
		line := l.a.log.Format("server", text, "")
		line.Seq = int64(i + 1)
		out = append(out, line)
	}
	return out, nil
}

// LogCleanup é a situação da limpeza dos logs antigos do servidor.
type LogCleanup struct {
	Enabled bool  `json:"enabled"`
	Keep    int   `json:"keep"`
	Min     int   `json:"min"`
	Max     int   `json:"max"`
	Runs    int   `json:"runs"`    // partidas guardadas agora
	Old     int   `json:"old"`     // quantas a limpeza apagaria
	OldSize int64 `json:"oldSize"` // bytes que a limpeza liberaria
	Deleted int   `json:"deleted"` // na resposta de CleanOld
	Freed   int64 `json:"freed"`
}

func (l *LogAPI) cleanup() LogCleanup {
	st := l.a.cfg.Settings()
	dir := l.a.cfg.Profile().ServerDir
	c := LogCleanup{Enabled: st.LogCleanup, Keep: st.LogKeep, Min: serverlogs.MinKeep, Max: serverlogs.MaxKeep,
		Runs: len(config.ServerLogs(dir))}
	for _, f := range config.OldLogs(dir, st.LogKeep) {
		c.Old++
		c.OldSize += f.Size
	}
	return c
}

// Cleanup devolve a configuração da limpeza e quanto ela apagaria agora.
func (l *LogAPI) Cleanup() LogCleanup { return l.cleanup() }

// SetCleanup liga ou desliga a limpeza automática (a cada início do servidor) e define quantas partidas manter.
func (l *LogAPI) SetCleanup(enabled bool, keep int) (LogCleanup, error) {
	if _, err := l.a.cfg.SetLogRetention(enabled, keep); err != nil {
		return l.cleanup(), err
	}
	return l.cleanup(), nil
}

// CleanOld apaga agora os logs antigos, mantendo as partidas definidas. Nunca apaga o log da partida atual.
func (l *LogAPI) CleanOld() (LogCleanup, error) {
	dir := l.a.cfg.Profile().ServerDir
	if dir == "" {
		return l.cleanup(), problems.New(problems.ServerNotInstalled, "")
	}
	n, freed, err := config.DeleteOldLogs(dir, l.a.cfg.Settings().LogKeep)
	c := l.cleanup()
	c.Deleted, c.Freed = n, freed
	if err != nil {
		return c, problems.Wrap(problems.LogDelete, err)
	}
	return c, nil
}
