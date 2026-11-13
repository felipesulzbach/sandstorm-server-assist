package modio

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"sandstorm-server-assist/desktop/internal/domain/pak"
)

// gameKeyCache guarda a chave lida por pasta do servidor e data do .pak (o índice do pakchunk0 é grande).
var gameKeyCache struct {
	mu   sync.Mutex
	path string
	mod  time.Time
	key  string
}

// GameAPIKey lê a chave de API do mod.io que o próprio jogo usa, de Insurgency/Config/DefaultGame.ini dentro do
// pakchunk0 do servidor instalado (ADR-044). É com ela que o servidor troca o código de segurança, então o pedido
// do código precisa usar a mesma. "" se o servidor não está instalado ou a chave não foi achada. A chave nunca vai
// para a tela nem para o log.
func GameAPIKey(serverDir string) string {
	if serverDir == "" {
		return ""
	}
	paks, _ := filepath.Glob(filepath.Join(serverDir, "Insurgency", "Content", "Paks", "pakchunk0-*.pak"))
	for _, p := range paks {
		info, err := os.Stat(p)
		if err != nil {
			continue
		}
		gameKeyCache.mu.Lock()
		if gameKeyCache.path == p && gameKeyCache.mod.Equal(info.ModTime()) && gameKeyCache.key != "" {
			k := gameKeyCache.key
			gameKeyCache.mu.Unlock()
			return k
		}
		gameKeyCache.mu.Unlock()
		if k, err := gameKeyFromPak(p); err == nil && k != "" {
			gameKeyCache.mu.Lock()
			gameKeyCache.path, gameKeyCache.mod, gameKeyCache.key = p, info.ModTime(), k
			gameKeyCache.mu.Unlock()
			return k
		}
	}
	return ""
}

func gameKeyFromPak(path string) (string, error) {
	files, err := pakFiles(path)
	if err != nil {
		return "", err
	}
	var entry *pak.File
	for i := range files {
		if strings.EqualFold(files[i].Path, pak.GameConfigFile) {
			entry = &files[i]
			break
		}
	}
	if entry == nil {
		return "", errors.New("game config not found")
	}
	f, err := os.Open(path) // #nosec G304 -- .pak do servidor instalado (pasta do perfil)
	if err != nil {
		return "", err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return "", err
	}
	size := info.Size()
	read := func(off, n int64) ([]byte, error) {
		if off < 0 || n < 0 || n > pak.MaxFileData || off > size-n {
			return nil, pak.ErrFormat
		}
		b := make([]byte, n)
		if _, err := f.ReadAt(b, off); err != nil && err != io.EOF {
			return nil, err
		}
		return b, nil
	}
	tail, err := read(size-min(int64(pak.TailSize), size), min(int64(pak.TailSize), size))
	if err != nil {
		return "", err
	}
	header, err := read(entry.Offset, min(int64(pak.MaxEntryHeader), size-entry.Offset))
	if err != nil {
		return "", err
	}
	e, err := pak.ParseEntry(header)
	if err != nil {
		return "", err
	}
	data, err := pak.ReadData(e, pak.Methods(tail), func(pos, n int64) ([]byte, error) {
		return read(entry.Offset+pos, n)
	})
	if err != nil {
		return "", err
	}
	key, ok := pak.ModioAPIKey(string(data))
	if !ok {
		return "", errors.New("no mod.io key in the game config")
	}
	return key, nil
}
