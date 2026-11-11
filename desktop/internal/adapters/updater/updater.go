// Package updater é o lado de fora da atualização automática (ADR-014): HTTPS pelo netx (só hosts permitidos),
// disco e o lançador do instalador.
package updater

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"sandstorm-server-assist/desktop/internal/adapters/netx"
	"sandstorm-server-assist/desktop/internal/domain/update"
)

// DefaultBaseURL é a pasta das Releases do projeto no GitHub (ADR-032). O latest.json vem da Release mais recente
// ("latest/download/latest.json") e o instalador vem da Release da própria versão ("download/v<versão>/<arquivo>"):
// se outra Release sair no meio do download, o instalador continua o mesmo que o manifesto assinou.
const DefaultBaseURL = "https://github.com/felipesulzbach/sandstorm-server-assist/releases/"

// PageURL é a página da Release mais recente (para a versão portátil baixar a nova à mão).
const PageURL = DefaultBaseURL + "latest"

const maxManifestBytes = 64 << 10

// HTTPFetcher busca o latest.json e o instalador por HTTPS.
type HTTPFetcher struct{ base *url.URL }

// NewHTTPFetcher cria o fetcher. Só aceita HTTPS para um host da lista do netx, terminando em "/".
func NewHTTPFetcher(base string) (*HTTPFetcher, error) {
	u, err := url.Parse(base)
	if err != nil || !strings.HasSuffix(u.Path, "/") {
		return nil, errors.New("invalid update address")
	}
	if err := netx.CheckURL(u); err != nil {
		return nil, err
	}
	return &HTTPFetcher{base: u}, nil
}

// fileName aceita só nomes simples de arquivo (sem pastas nem ".."): o nome vem do manifesto assinado, mas a
// checagem fica aqui também.
var fileName = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)

// installerVersion tira a versão do nome do instalador (o manifesto já exige o nome da própria versão).
var installerVersion = regexp.MustCompile(`^sandstorm-server-assist-(\d{1,4}\.\d{1,4}\.\d{1,4})-windows-setup\.exe$`)

const manifestName = "latest.json"

func (f *HTTPFetcher) url(name string) (string, error) {
	if !fileName.MatchString(name) || strings.Contains(name, "..") {
		return "", errors.New("invalid update file name")
	}
	if name == manifestName {
		return f.base.JoinPath("latest", "download", manifestName).String(), nil
	}
	m := installerVersion.FindStringSubmatch(name)
	if m == nil {
		return "", errors.New("invalid update file name")
	}
	// #nosec G125 -- nome e versão conferidos acima (sem barras nem "..")
	return f.base.JoinPath("download", "v"+m[1], name).String(), nil
}

// Manifest busca o latest.json (no máximo 64 KB).
func (f *HTTPFetcher) Manifest(ctx context.Context) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	u, err := f.url(manifestName)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Cache-Control", "no-cache")
	resp, err := netx.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	return io.ReadAll(io.LimitReader(resp.Body, maxManifestBytes))
}

// Download grava o arquivo em dst (no máximo MaxInstallerBytes) e devolve o SHA-256 e o tamanho.
func (f *HTTPFetcher) Download(ctx context.Context, file, dst string, onProgress func(int64)) (string, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()
	u, err := f.url(file)
	if err != nil {
		return "", 0, err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o700); err != nil {
		return "", 0, err
	}
	err = netx.Download(ctx, u, dst, update.MaxInstallerBytes, func(done, _ int64) {
		if onProgress != nil {
			onProgress(done)
		}
	})
	if err != nil {
		return "", 0, err
	}
	return hashFile(dst)
}

func hashFile(path string) (string, int64, error) {
	in, err := os.Open(path) // #nosec G304 -- arquivo baixado na pasta de atualizações do app
	if err != nil {
		return "", 0, err
	}
	defer in.Close()
	h := sha256.New()
	n, err := io.Copy(h, in)
	if err != nil {
		return "", 0, err
	}
	return hex.EncodeToString(h.Sum(nil)), n, nil
}

// OSFiles implementa o pouco de disco da atualização.
type OSFiles struct{}

func (OSFiles) Rename(from, to string) error { return os.Rename(from, to) }

func (OSFiles) Remove(p string) error {
	if err := os.Remove(p); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

// Dir é a pasta das atualizações baixadas. Começa vazia a cada execução (apaga instaladores antigos).
func Dir() (string, error) {
	base, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	d := filepath.Join(base, "sandstorm-server-assist", "updates")
	_ = os.RemoveAll(d)
	return d, os.MkdirAll(d, 0o700)
}

// Installed diz se o app veio do instalador: o desinstalador fica ao lado do executável. A versão portátil
// não tem, e então só avisa da versão nova.
func Installed() bool {
	exe, err := os.Executable()
	if err != nil {
		return false
	}
	st, err := os.Stat(filepath.Join(filepath.Dir(exe), "uninstall.exe"))
	return err == nil && st.Mode().IsRegular()
}
