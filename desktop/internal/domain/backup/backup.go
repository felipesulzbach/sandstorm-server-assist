// Package backup define o arquivo de backup portátil (ADR-006): o que vai dentro e a criptografia com senha
// (Argon2id + AES-256-GCM). Não lê nem grava arquivos: recebe e devolve bytes.
package backup

import (
	"archive/zip"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/crypto/argon2"

	"sandstorm-server-assist/desktop/internal/domain/problems"
	"sandstorm-server-assist/desktop/internal/domain/server"
)

// Extension é a extensão do arquivo de backup.
const Extension = ".ssbackup"

// FormatVersion sobe quando o conteúdo muda de forma incompatível.
const FormatVersion = 1

// MinPassword é o tamanho mínimo da senha do backup.
const MinPassword = 10

// Limites ao abrir (arquivo de outra pessoa pode ser malicioso).
const (
	maxFileSize  = 4 << 20  // cada arquivo dentro do backup
	maxTotalSize = 32 << 20 // backup inteiro, já aberto
	maxFiles     = 32
)

// magic identifica o arquivo: "SSAB" + versão do cabeçalho.
var magic = []byte{'S', 'S', 'A', 'B', 1}

// Parâmetros do Argon2id (gravados no cabeçalho para poder mudar no futuro).
type kdfParams struct {
	Time    uint32
	Memory  uint32 // KiB
	Threads uint8
}

var defaultKDF = kdfParams{Time: 3, Memory: 64 * 1024, Threads: 4}

// ConfigFiles são os arquivos de configuração do servidor que entram no backup (nome no backup → caminho
// relativo à pasta do servidor). Só esses nomes são aceitos ao restaurar.
var ConfigFiles = map[string]string{
	"server/Game.ini":     "Insurgency/Saved/Config/WindowsServer/Game.ini",
	"server/Engine.ini":   "Insurgency/Saved/Config/WindowsServer/Engine.ini",
	"server/Admins.txt":   "Insurgency/Config/Server/Admins.txt",
	"server/MapCycle.txt": "Insurgency/Config/Server/MapCycle.txt",
	"server/Bans.json":    "Insurgency/Config/Server/Bans.json",
	"server/Motd.txt":     "Insurgency/Config/Server/Motd.txt",
	"server/Mods.txt":     "Insurgency/Config/Server/Mods.txt",
	"server/Notes.txt":    "Insurgency/Config/Server/Notes.txt",
}

// ModioSessionName é o nome da sessão do mod.io (user.json) dentro do backup.
const ModioSessionName = "modio/user.json"

// Manifest descreve o backup.
type Manifest struct {
	Format     int       `json:"format"`
	AppVersion string    `json:"appVersion"`
	CreatedAt  time.Time `json:"createdAt"`
	Machine    string    `json:"machine"`
	Note       string    `json:"note,omitempty"`
	Files      []string  `json:"files"`
	HasModio   bool      `json:"hasModio"`
}

// Contents é tudo que o backup guarda.
type Contents struct {
	Manifest Manifest
	Profile  server.Profile
	Secrets  server.Secrets
	// Files: nome no backup → conteúdo (ConfigFiles e, se houver, ModioSessionName).
	Files map[string][]byte
}

// Pack monta o zip com o conteúdo (ainda sem criptografia).
func Pack(c Contents) ([]byte, error) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	c.Manifest.Format = FormatVersion
	c.Manifest.Files = c.Manifest.Files[:0]
	for name := range c.Files {
		if !allowedName(name) {
			return nil, fmt.Errorf("file not allowed in backup: %s", name)
		}
		c.Manifest.Files = append(c.Manifest.Files, name)
	}
	sort.Strings(c.Manifest.Files)
	_, c.Manifest.HasModio = c.Files[ModioSessionName]
	put := func(name string, data []byte) error {
		w, err := zw.Create(name)
		if err != nil {
			return err
		}
		_, err = w.Write(data)
		return err
	}
	for name, v := range map[string]any{"manifest.json": c.Manifest, "profile.json": c.Profile, "secrets.json": c.Secrets} {
		b, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			return nil, err
		}
		if err := put(name, b); err != nil {
			return nil, err
		}
	}
	for name, data := range c.Files {
		if err := put(name, data); err != nil {
			return nil, err
		}
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Unpack lê o zip, com limites de tamanho e só nomes conhecidos.
func Unpack(data []byte) (Contents, error) {
	var c Contents
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return c, problems.New(problems.BackupCorrupt, "not a valid zip")
	}
	if len(zr.File) > maxFiles+3 {
		return c, problems.New(problems.BackupCorrupt, "too many files")
	}
	c.Files = map[string][]byte{}
	var total int64
	seen := map[string]bool{}
	for _, f := range zr.File {
		name := f.Name
		if seen[name] {
			return c, problems.New(problems.BackupCorrupt, "duplicate name: "+name)
		}
		seen[name] = true
		meta := name == "manifest.json" || name == "profile.json" || name == "secrets.json"
		if !meta && !allowedName(name) {
			return c, problems.New(problems.BackupCorrupt, "unexpected file: "+name)
		}
		b, err := readLimited(f, maxFileSize)
		if err != nil {
			return c, err
		}
		total += int64(len(b))
		if total > maxTotalSize {
			return c, problems.New(problems.BackupCorrupt, "backup too large")
		}
		switch name {
		case "manifest.json":
			err = json.Unmarshal(b, &c.Manifest)
		case "profile.json":
			err = json.Unmarshal(b, &c.Profile)
		case "secrets.json":
			err = json.Unmarshal(b, &c.Secrets)
		default:
			c.Files[name] = b
		}
		if err != nil {
			return c, problems.New(problems.BackupCorrupt, "invalid content in "+name)
		}
	}
	if c.Manifest.Format == 0 {
		return c, problems.New(problems.BackupCorrupt, "missing manifest")
	}
	if c.Manifest.Format > FormatVersion {
		return c, problems.New(problems.BackupCorrupt, "made by a newer app version")
	}
	return c, nil
}

func readLimited(f *zip.File, limit int64) ([]byte, error) {
	if f.UncompressedSize64 > uint64(limit) { // #nosec G115 -- limit é constante positiva (maxFileSize)
		return nil, problems.New(problems.BackupCorrupt, "file too large: "+f.Name)
	}
	rc, err := f.Open()
	if err != nil {
		return nil, problems.New(problems.BackupCorrupt, "cannot read "+f.Name)
	}
	defer rc.Close()
	b, err := io.ReadAll(io.LimitReader(rc, limit+1))
	if err != nil || int64(len(b)) > limit {
		return nil, problems.New(problems.BackupCorrupt, "cannot read "+f.Name)
	}
	return b, nil
}

func allowedName(name string) bool {
	if name == ModioSessionName {
		return true
	}
	_, ok := ConfigFiles[name]
	return ok && path.Clean(name) == name && !strings.Contains(name, "..")
}

// ValidatePassword confere a senha do backup.
func ValidatePassword(pw string) error {
	if utf8.RuneCountInString(pw) < MinPassword {
		return problems.Field("password", "validation.backupPassword", "min", strconv.Itoa(MinPassword))
	}
	return nil
}

// Seal cifra o conteúdo com a senha. rnd é a fonte de aleatoriedade (crypto/rand.Reader no app).
// Formato: magic | time u32 | memory u32 | threads u8 | sal 16 | nonce 12 | dados cifrados (com tag).
func Seal(password string, plain []byte, rnd io.Reader) ([]byte, error) {
	if err := ValidatePassword(password); err != nil {
		return nil, err
	}
	p := defaultKDF
	salt := make([]byte, 16)
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rnd, salt); err != nil {
		return nil, err
	}
	if _, err := io.ReadFull(rnd, nonce); err != nil {
		return nil, err
	}
	header := new(bytes.Buffer)
	header.Write(magic)
	_ = binary.Write(header, binary.BigEndian, p.Time)
	_ = binary.Write(header, binary.BigEndian, p.Memory)
	header.WriteByte(p.Threads)
	header.Write(salt)
	header.Write(nonce)
	aead, err := newAEAD(password, salt, p)
	if err != nil {
		return nil, err
	}
	// O cabeçalho entra como dado autenticado: mexer nos parâmetros invalida o arquivo.
	hdr := header.Bytes()
	return aead.Seal(append([]byte{}, hdr...), nonce, plain, hdr), nil
}

// Open decifra. Senha errada e arquivo alterado dão o mesmo erro (a criptografia não distingue).
func Open(password string, data []byte) ([]byte, error) {
	hl := len(magic) + 4 + 4 + 1 + 16 + 12
	if len(data) < hl+16 || !bytes.Equal(data[:len(magic)], magic) {
		return nil, problems.New(problems.BackupCorrupt, "not a backup file")
	}
	r := bytes.NewReader(data[len(magic):hl])
	var p kdfParams
	_ = binary.Read(r, binary.BigEndian, &p.Time)
	_ = binary.Read(r, binary.BigEndian, &p.Memory)
	p.Threads, _ = r.ReadByte()
	// Parâmetros absurdos travariam o app (memória): aceita só uma faixa razoável.
	if p.Time < 1 || p.Time > 10 || p.Memory < 8*1024 || p.Memory > 512*1024 || p.Threads < 1 || p.Threads > 16 {
		return nil, problems.New(problems.BackupCorrupt, "invalid header")
	}
	salt := make([]byte, 16)
	nonce := make([]byte, 12)
	_, _ = io.ReadFull(r, salt)
	_, _ = io.ReadFull(r, nonce)
	aead, err := newAEAD(password, salt, p)
	if err != nil {
		return nil, err
	}
	plain, err := aead.Open(nil, nonce, data[hl:], data[:hl])
	if err != nil {
		return nil, problems.New(problems.BackupPassword, "")
	}
	return plain, nil
}

func newAEAD(password string, salt []byte, p kdfParams) (cipher.AEAD, error) {
	key := argon2.IDKey([]byte(password), salt, p.Time, p.Memory, p.Threads, 32)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}
