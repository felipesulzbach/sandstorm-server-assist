// Package vault guarda os segredos (senhas, tokens, chave do mod.io) cifrados com o DPAPI do Windows, no escopo do
// usuário (ADR-002). Só o mesmo usuário, na mesma máquina, consegue abrir.
package vault

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"sandstorm-server-assist/desktop/internal/adapters/fsx"
	"sandstorm-server-assist/desktop/internal/domain/problems"
	"sandstorm-server-assist/desktop/internal/domain/server"
)

// entropy é misturada à chave do DPAPI: outro programa do mesmo usuário que chame CryptUnprotectData sem ela
// não abre o cofre.
var entropy = []byte("sandstorm-server-assist/vault/v1")

// Vault é o arquivo do cofre.
type Vault struct{ Path string }

// Default usa %APPDATA%\sandstorm-server-assist\vault.bin.
func Default() (Vault, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return Vault{}, err
	}
	return Vault{Path: filepath.Join(base, "sandstorm-server-assist", "vault.bin")}, nil
}

// Load abre o cofre. Sem arquivo, devolve segredos vazios.
func (v Vault) Load() (server.Secrets, error) {
	var s server.Secrets
	data, err := fsx.ReadLimited(v.Path, 1<<20)
	if errors.Is(err, fs.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return s, problems.Wrap(problems.VaultUnavailable, err)
	}
	plain, err := unprotect(data)
	if err != nil {
		return s, problems.Wrap(problems.VaultUnavailable, err)
	}
	defer wipe(plain)
	if err := json.Unmarshal(plain, &s); err != nil {
		return s, problems.New(problems.VaultUnavailable, "invalid vault content")
	}
	return s, nil
}

// Save grava o cofre (troca atômica).
func (v Vault) Save(s server.Secrets) error {
	plain, err := json.Marshal(s)
	if err != nil {
		return err
	}
	defer wipe(plain)
	data, err := protect(plain)
	if err != nil {
		return problems.Wrap(problems.VaultUnavailable, err)
	}
	return fsx.WriteAtomic(v.Path, data, false)
}

// wipe zera um buffer com segredo depois de usar (reduz o tempo em memória; o GC pode ter feito cópias).
func wipe(b []byte) {
	for i := range b {
		b[i] = 0
	}
}
