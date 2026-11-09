// Package backupsvc cria e restaura o backup portátil (ADR-006): junta perfil, segredos e arquivos do servidor,
// cifra com a senha do usuário e grava; ao restaurar, confere tudo antes de gravar na máquina nova.
package backupsvc

import (
	"crypto/rand"
	"os"
	"path/filepath"
	"sort"
	"time"

	"sandstorm-server-assist/desktop/internal/adapters/fsx"
	"sandstorm-server-assist/desktop/internal/adapters/modio"
	"sandstorm-server-assist/desktop/internal/domain/backup"
	"sandstorm-server-assist/desktop/internal/domain/problems"
	"sandstorm-server-assist/desktop/internal/domain/server"
)

// Source é de onde vêm perfil e segredos (config.Service).
type Source interface {
	Current() (server.Profile, server.Secrets, error)
	ReplaceAll(p server.Profile, s server.Secrets, serverDir string) error
}

// Summary é o que a tela mostra de um backup antes de restaurar.
type Summary struct {
	CreatedAt  time.Time `json:"createdAt"`
	AppVersion string    `json:"appVersion"`
	Machine    string    `json:"machine"`
	Hostname   string    `json:"hostname"`
	OldDir     string    `json:"oldDir"`
	Files      []string  `json:"files"`
	Secrets    []string  `json:"secrets"`
	HasModio   bool      `json:"hasModio"`
	Admins     int       `json:"admins"`
	Mods       int       `json:"mods"`
}

// Options do backup.
type Options struct {
	IncludeModio bool `json:"includeModio"`
}

// Create monta e grava o backup em dest.
func Create(src Source, appVersion, dest, password string, opt Options) (Summary, error) {
	if err := backup.ValidatePassword(password); err != nil {
		return Summary{}, err
	}
	p, sec, err := src.Current()
	if err != nil {
		return Summary{}, err
	}
	host, _ := os.Hostname()
	c := backup.Contents{
		Manifest: backup.Manifest{AppVersion: appVersion, CreatedAt: time.Now().UTC(), Machine: host},
		Profile:  p,
		Secrets:  sec,
		Files:    map[string][]byte{},
	}
	if p.ServerDir != "" {
		for name, rel := range backup.ConfigFiles {
			if b, err := fsx.ReadLimited(filepath.Join(p.ServerDir, filepath.FromSlash(rel)), 4<<20); err == nil {
				c.Files[name] = b
			}
		}
	}
	if opt.IncludeModio {
		if uf, err := modio.UserFile(); err == nil {
			if b, err := fsx.ReadLimited(uf, 1<<20); err == nil {
				c.Files[backup.ModioSessionName] = b
			}
		}
	}
	plain, err := backup.Pack(c)
	if err != nil {
		return Summary{}, problems.Wrap(problems.BackupWrite, err)
	}
	sealed, err := backup.Seal(password, plain, rand.Reader)
	if err != nil {
		return Summary{}, err
	}
	if err := fsx.WriteAtomic(dest, sealed, false); err != nil {
		return Summary{}, problems.Wrap(problems.BackupWrite, err)
	}
	return summarize(c), nil
}

// Open lê e decifra um backup.
func Open(path, password string) (backup.Contents, error) {
	data, err := fsx.ReadLimited(path, 64<<20)
	if err != nil {
		return backup.Contents{}, problems.Wrap(problems.BackupCorrupt, err)
	}
	plain, err := backup.Open(password, data)
	if err != nil {
		return backup.Contents{}, err
	}
	return backup.Unpack(plain)
}

// Inspect devolve o resumo de um backup (pede a senha).
func Inspect(path, password string) (Summary, error) {
	c, err := Open(path, password)
	if err != nil {
		return Summary{}, err
	}
	return summarize(c), nil
}

// Restore grava o backup nesta máquina: perfil e segredos no cofre local, arquivos do servidor em serverDir e,
// se o usuário marcar, a sessão do mod.io.
func Restore(src Source, path, password, serverDir string, restoreModio bool) (Summary, error) {
	if !fsx.SafeDir(serverDir) {
		return Summary{}, problems.New(problems.ServerDirInvalid, serverDir)
	}
	c, err := Open(path, password)
	if err != nil {
		return Summary{}, err
	}
	if err := src.ReplaceAll(c.Profile, c.Secrets, filepath.Clean(serverDir)); err != nil {
		return Summary{}, err
	}
	for name, data := range c.Files {
		if name == backup.ModioSessionName {
			continue
		}
		rel, ok := backup.ConfigFiles[name]
		if !ok {
			continue
		}
		if err := fsx.WriteAtomic(filepath.Join(serverDir, filepath.FromSlash(rel)), data, true); err != nil {
			return Summary{}, problems.Wrap(problems.ConfigWrite, err)
		}
	}
	if data, ok := c.Files[backup.ModioSessionName]; ok && restoreModio {
		uf, err := modio.UserFile()
		if err != nil {
			return Summary{}, problems.Wrap(problems.ConfigWrite, err)
		}
		if err := fsx.WriteAtomic(uf, data, true); err != nil {
			return Summary{}, problems.Wrap(problems.ConfigWrite, err)
		}
	}
	return summarize(c), nil
}

func summarize(c backup.Contents) Summary {
	s := Summary{
		CreatedAt: c.Manifest.CreatedAt, AppVersion: c.Manifest.AppVersion, Machine: c.Manifest.Machine,
		Hostname: c.Profile.Hostname, OldDir: c.Profile.ServerDir, HasModio: c.Manifest.HasModio,
		Admins: len(c.Profile.Admins), Mods: len(c.Profile.ModIDs),
	}
	for name := range c.Files {
		if name != backup.ModioSessionName {
			s.Files = append(s.Files, filepath.Base(name))
		}
	}
	sort.Strings(s.Files)
	for _, k := range server.SecretKeys {
		if v, _ := c.Secrets.Get(k); v != "" {
			s.Secrets = append(s.Secrets, k)
		}
	}
	return s
}
