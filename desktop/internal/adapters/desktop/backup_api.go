package desktop

import (
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"sandstorm-server-assist/desktop/internal/app/backupsvc"
	"sandstorm-server-assist/desktop/internal/domain/backup"
	"sandstorm-server-assist/desktop/internal/domain/problems"
)

// BackupAPI: criar e restaurar o backup portátil (ADR-006). Os caminhos vêm só dos diálogos abertos aqui.
type BackupAPI struct {
	a  *app
	mu sync.Mutex
	// chosen é o backup escolhido para restaurar; restoreDir, a pasta do servidor nesta máquina.
	chosen     string
	restoreDir string
}

// CreateResult diz onde o backup foi salvo.
type CreateResult struct {
	Saved   bool              `json:"saved"`
	File    string            `json:"file"`
	Summary backupsvc.Summary `json:"summary"`
}

// Create pede onde salvar e grava o backup cifrado com a senha.
func (b *BackupAPI) Create(password string, includeModio bool) (CreateResult, error) {
	if err := backup.ValidatePassword(password); err != nil {
		return CreateResult{}, err
	}
	name := "sandstorm-server-" + time.Now().Format("2006-01-02") + backup.Extension
	dest, err := runtime.SaveFileDialog(b.a.ctx, runtime.SaveDialogOptions{
		DefaultFilename: name,
		Filters:         []runtime.FileFilter{{DisplayName: "Sandstorm Server Assist backup (*.ssbackup)", Pattern: "*.ssbackup"}},
	})
	if err != nil || dest == "" {
		return CreateResult{}, err
	}
	if !strings.EqualFold(filepath.Ext(dest), backup.Extension) {
		dest += backup.Extension
	}
	sum, err := backupsvc.Create(b.a.cfg, AppVersion, dest, password, backupsvc.Options{IncludeModio: includeModio})
	if err != nil {
		return CreateResult{}, err
	}
	b.a.log.Add("app", "Backup saved: "+filepath.Base(dest), "info")
	return CreateResult{Saved: true, File: filepath.Base(dest), Summary: sum}, nil
}

// ChooseFile abre o diálogo para escolher um backup. Devolve só o nome do arquivo.
func (b *BackupAPI) ChooseFile() (string, error) {
	path, err := runtime.OpenFileDialog(b.a.ctx, runtime.OpenDialogOptions{
		Filters: []runtime.FileFilter{{DisplayName: "Sandstorm Server Assist backup (*.ssbackup)", Pattern: "*.ssbackup"}},
	})
	if err != nil || path == "" {
		return "", err
	}
	b.mu.Lock()
	b.chosen = path
	b.mu.Unlock()
	return filepath.Base(path), nil
}

// Inspect abre o backup escolhido com a senha e mostra o resumo.
func (b *BackupAPI) Inspect(password string) (backupsvc.Summary, error) {
	b.mu.Lock()
	path := b.chosen
	b.mu.Unlock()
	if path == "" {
		return backupsvc.Summary{}, problems.New(problems.BackupCorrupt, "no file chosen")
	}
	return backupsvc.Inspect(path, password)
}

// ChooseRestoreDir escolhe a pasta do servidor nesta máquina.
func (b *BackupAPI) ChooseRestoreDir() (string, error) {
	dir, err := runtime.OpenDirectoryDialog(b.a.ctx, runtime.OpenDialogOptions{CanCreateDirectories: true})
	if err != nil || dir == "" {
		return "", err
	}
	b.mu.Lock()
	b.restoreDir = dir
	b.mu.Unlock()
	return dir, nil
}

// Restore grava o backup nesta máquina (servidor parado).
func (b *BackupAPI) Restore(password string, restoreModio bool) (backupsvc.Summary, error) {
	if len(b.a.runner.Status().PIDs) > 0 {
		return backupsvc.Summary{}, problems.New(problems.ServerRunning, "")
	}
	b.mu.Lock()
	path, dir := b.chosen, b.restoreDir
	b.mu.Unlock()
	if path == "" {
		return backupsvc.Summary{}, problems.New(problems.BackupCorrupt, "no file chosen")
	}
	if dir == "" {
		return backupsvc.Summary{}, problems.New(problems.ServerDirInvalid, "")
	}
	sum, err := backupsvc.Restore(b.a.cfg, path, password, dir, restoreModio)
	if err != nil {
		return sum, err
	}
	b.a.log.Add("app", "Backup restored into "+dir, "info")
	return sum, nil
}
