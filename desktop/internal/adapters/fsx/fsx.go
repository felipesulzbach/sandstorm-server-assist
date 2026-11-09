// Package fsx tem as operações de arquivo do app: leitura com limite, gravação atômica (com cópia .bak antes da
// primeira troca de um arquivo do usuário) e tamanho de pastas.
package fsx

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// ReadLimited lê um arquivo inteiro, recusando se passar de limit bytes.
func ReadLimited(path string, limit int64) ([]byte, error) {
	f, err := os.Open(path) // #nosec G304 -- caminhos montados pelo app (pasta do servidor/config), não pela tela
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := io.ReadAll(io.LimitReader(f, limit+1))
	if err != nil {
		return nil, err
	}
	if int64(len(b)) > limit {
		return nil, errors.New("file too large: " + filepath.Base(path))
	}
	return b, nil
}

// ReadHead lê no máximo n bytes do começo de um arquivo (logs grandes, ainda abertos pelo servidor).
func ReadHead(path string, n int64) ([]byte, error) {
	f, err := os.Open(path) // #nosec G304 -- caminhos montados pelo app (pasta de logs do servidor), não pela tela
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(io.LimitReader(f, n))
}

// WriteAtomic grava num temporário ao lado e troca pelo nome final. Com backup=true e o arquivo já existindo sem
// um .bak, guarda uma cópia <nome>.bak antes (a versão original do usuário nunca se perde).
func WriteAtomic(path string, data []byte, backup bool) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	if backup {
		if _, err := os.Stat(path + ".bak"); errors.Is(err, fs.ErrNotExist) {
			if orig, err := os.ReadFile(path); err == nil { // #nosec G304 -- arquivo de configuração conhecido
				if err := os.WriteFile(path+".bak", orig, 0o600); err != nil { // #nosec G703 -- caminho montado pelo app (config do servidor), não pela tela
					return err
				}
			}
		}
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), "."+filepath.Base(path)+".*.tmp")
	if err != nil {
		return err
	}
	name := tmp.Name()
	_, werr := tmp.Write(data)
	if serr := tmp.Sync(); werr == nil {
		werr = serr
	}
	if cerr := tmp.Close(); werr == nil {
		werr = cerr
	}
	if werr != nil {
		_ = os.Remove(name)
		return werr
	}
	if err := os.Rename(name, path); err != nil {
		_ = os.Remove(name)
		return err
	}
	return nil
}

// DirSize soma o tamanho dos arquivos de uma pasta (com subpastas). Pasta ausente = 0.
func DirSize(dir string) int64 {
	var total int64
	_ = filepath.WalkDir(dir, func(_ string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() {
			if info, err := d.Info(); err == nil {
				total += info.Size()
			}
		}
		return nil
	})
	return total
}

// Exists diz se o caminho existe.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// SafeDir confere uma pasta escolhida pelo usuário para o servidor ou o SteamCMD: caminho completo e fora das
// pastas do sistema.
func SafeDir(dir string) bool {
	if !filepath.IsAbs(dir) {
		return false
	}
	clean := strings.ToLower(filepath.Clean(dir))
	for _, env := range []string{"SystemRoot", "ProgramFiles", "ProgramFiles(x86)", "ProgramData"} {
		if v := os.Getenv(env); v != "" {
			sys := strings.ToLower(filepath.Clean(v))
			if clean == sys || strings.HasPrefix(clean, sys+string(filepath.Separator)) {
				return false
			}
		}
	}
	// Raiz de disco sozinha (D:\) também não: o SteamCMD espalharia arquivos na raiz.
	return filepath.Dir(clean) != clean
}
