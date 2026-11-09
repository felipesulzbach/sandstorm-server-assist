// Package update guarda as regras puras da atualização automática (ADR-014): o manifesto publicado ao lado do
// instalador, a mensagem que a chave de publicação assina e a comparação de versões. Nada aqui faz I/O.
package update

import (
	"crypto/ed25519"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// MaxInstallerBytes limita o tamanho do instalador baixado (o atual tem poucos MB).
const MaxInstallerBytes = 200 << 20

// Prefixo do nome do instalador publicado e da mensagem assinada.
const (
	appName       = "sandstorm-server-assist"
	messageHeader = appName + "-update-v1"
)

// Manifest é o latest.json publicado ao lado do instalador.
type Manifest struct {
	Version   string `json:"version"`   // ex.: 0.2.0
	File      string `json:"file"`      // nome do instalador na mesma pasta do manifesto
	SHA256    string `json:"sha256"`    // hex do instalador
	Size      int64  `json:"size"`      // bytes do instalador
	Date      string `json:"date"`      // AAAA-MM-DD, só informativo
	Signature string `json:"signature"` // Ed25519 (base64) de Message()
}

var (
	versionRe = regexp.MustCompile(`^\d{1,4}\.\d{1,4}\.\d{1,4}$`)
	shaRe     = regexp.MustCompile(`^[0-9a-f]{64}$`)

	ErrBadManifest  = errors.New("invalid update manifest")
	ErrBadSignature = errors.New("update signature does not match")
)

// InstallerName é o nome do instalador de uma versão (o mesmo que o build gera em build/dist).
func InstallerName(version string) string { return appName + "-" + version + "-windows-setup.exe" }

// Message é exatamente o que a chave de publicação assina (scripts/update-manifest.mjs monta a mesma coisa).
func (m Manifest) Message() []byte {
	return []byte(fmt.Sprintf("%s\n%s\n%s\n%s\n%d\n", messageHeader, m.Version, m.File, m.SHA256, m.Size))
}

// Validate confere o formato dos campos. O nome do arquivo precisa ter a mesma versão, para um manifesto antigo
// não apontar para outro instalador.
func (m Manifest) Validate() error {
	if !versionRe.MatchString(m.Version) || !shaRe.MatchString(m.SHA256) || m.Size <= 0 ||
		m.Size > MaxInstallerBytes || m.File != InstallerName(m.Version) {
		return ErrBadManifest
	}
	return nil
}

// Verify confere o formato e a assinatura com a chave pública embutida no app.
func (m Manifest) Verify(pub ed25519.PublicKey) error {
	if err := m.Validate(); err != nil {
		return err
	}
	sig, err := base64.StdEncoding.DecodeString(m.Signature)
	if err != nil || len(pub) != ed25519.PublicKeySize || !ed25519.Verify(pub, m.Message(), sig) {
		return ErrBadSignature
	}
	return nil
}

// Newer diz se candidate é maior que current. Versões de desenvolvimento ("0.1.0-dev") nunca se atualizam.
func Newer(current, candidate string) bool {
	c, ok1 := parse(current)
	n, ok2 := parse(candidate)
	if !ok1 || !ok2 {
		return false
	}
	for i := range c {
		if n[i] != c[i] {
			return n[i] > c[i]
		}
	}
	return false
}

// IsRelease diz se a versão é de uma build publicada (sem sufixo de desenvolvimento).
func IsRelease(v string) bool { _, ok := parse(v); return ok }

func parse(v string) ([3]int, bool) {
	var out [3]int
	if !versionRe.MatchString(v) {
		return out, false
	}
	for i, p := range strings.Split(v, ".") {
		out[i], _ = strconv.Atoi(p)
	}
	return out, true
}
