package config

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"

	"sandstorm-server-assist/desktop/internal/domain/problems"
	"sandstorm-server-assist/desktop/internal/domain/serverlogs"
)

// SetLogRetention liga ou desliga a limpeza automática dos logs do servidor e define quantas partidas manter.
func (s *Service) SetLogRetention(on bool, keep int) (Settings, error) {
	if !serverlogs.ValidKeep(keep) {
		return s.Settings(), problems.Field("logKeep", "validation.logKeep",
			"min", strconv.Itoa(serverlogs.MinKeep), "max", strconv.Itoa(serverlogs.MaxKeep))
	}
	s.mu.Lock()
	s.data.Settings.LogCleanup, s.data.Settings.LogKeep = on, keep
	err := s.save()
	s.mu.Unlock()
	return s.Settings(), err
}

// ServerLogs lista os logs da pasta de logs de uma instância (LogsPath; só os nomes que o servidor cria), do mais
// novo para o mais velho.
func ServerLogs(logsDir string) []serverlogs.File {
	if logsDir == "" {
		return nil
	}
	entries, err := os.ReadDir(logsDir)
	if err != nil {
		return nil
	}
	var out []serverlogs.File
	for _, e := range entries {
		if e.IsDir() || !serverlogs.IsServerLog(e.Name()) {
			continue
		}
		if info, err := e.Info(); err == nil {
			out = append(out, serverlogs.File{Name: e.Name(), ModTime: info.ModTime(), Size: info.Size()})
		}
	}
	serverlogs.Sort(out)
	return out
}

// OldLogs devolve os logs que a limpeza apagaria mantendo keep partidas.
func OldLogs(logsDir string, keep int) []serverlogs.File {
	return serverlogs.ToDelete(ServerLogs(logsDir), keep)
}

// DeleteOldLogs apaga os logs antigos mantendo keep partidas. O log atual nunca é apagado. Devolve quantos
// arquivos e quantos bytes saíram; um arquivo que não pôde ser apagado (em uso) não impede os outros.
func DeleteOldLogs(logsDir string, keep int) (int, int64, error) {
	dir := logsDir
	var n int
	var freed int64
	var errs []error
	for _, f := range OldLogs(logsDir, keep) {
		if !serverlogs.IsServerLog(f.Name) || f.Name == serverlogs.Current {
			continue
		}
		if err := os.Remove(filepath.Join(dir, f.Name)); err != nil && !errors.Is(err, os.ErrNotExist) {
			errs = append(errs, err)
			continue
		}
		n++
		freed += f.Size
	}
	return n, freed, errors.Join(errs...)
}
