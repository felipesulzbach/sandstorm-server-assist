package pak

import (
	"bytes"
	"compress/zlib"
	"errors"
	"io"
	"regexp"
	"strings"

	"sandstorm-server-assist/desktop/internal/domain/ini"
)

// Limites da leitura do conteúdo de um arquivo (só arquivos pequenos de configuração).
const (
	MaxEntryHeader = 64 << 10
	MaxFileData    = 4 << 20
	maxBlocks      = 1024
)

// ErrUnsupported: o arquivo está cifrado ou usa uma compressão que o app não lê.
var ErrUnsupported = errors.New("pak: encrypted or unsupported compression")

// Entry é o cabeçalho que o .pak grava na posição do arquivo, antes dos dados (FPakEntry, formato 8+).
type Entry struct {
	Size             int64
	UncompressedSize int64
	Method           uint32 // 0 = sem compressão; senão, índice na lista de métodos do rodapé
	Blocks           [][2]int64
	Encrypted        bool
	HeaderSize       int64
}

// ParseEntry lê o cabeçalho na posição do arquivo. Os blocos são posições a partir do início do cabeçalho.
func ParseEntry(b []byte) (Entry, error) {
	r := &reader{b: b}
	r.skip(8) // posição (zero no cabeçalho)
	e := Entry{Size: r.i64(), UncompressedSize: r.i64()}
	e.Method = uint32(r.i32()) // #nosec G115 -- índice de compressão do formato
	r.skip(20)                 // SHA1
	if e.Method != 0 {
		n := int(r.i32())
		if n < 0 || n > maxBlocks {
			return Entry{}, ErrFormat
		}
		for i := 0; i < n && r.err == nil; i++ {
			e.Blocks = append(e.Blocks, [2]int64{r.i64(), r.i64()})
		}
	}
	if r.need(1) {
		e.Encrypted = r.b[r.off]&1 != 0
		r.off++
	}
	r.skip(4) // tamanho do bloco de compressão
	if r.err != nil || e.Size < 0 || e.UncompressedSize < 0 || e.UncompressedSize > MaxFileData {
		return Entry{}, ErrFormat
	}
	e.HeaderSize = int64(r.off)
	for _, bl := range e.Blocks {
		if bl[0] < e.HeaderSize || bl[1] < bl[0] || bl[1]-bl[0] > MaxFileData {
			return Entry{}, ErrFormat
		}
	}
	return e, nil
}

// ReadData monta o conteúdo do arquivo a partir do cabeçalho. read(pos, n) lê n bytes a partir do início do
// cabeçalho. methods são os nomes de compressão do rodapé (só "Zlib" é aceito). O resultado tem no máximo
// MaxFileData bytes e precisa bater com o tamanho do cabeçalho.
func ReadData(e Entry, methods []string, read func(pos, n int64) ([]byte, error)) ([]byte, error) {
	if e.Encrypted {
		return nil, ErrUnsupported
	}
	if e.Method == 0 {
		if e.Size != e.UncompressedSize {
			return nil, ErrFormat
		}
		return read(e.HeaderSize, e.Size)
	}
	if int(e.Method) > len(methods) || !strings.EqualFold(methods[e.Method-1], "Zlib") {
		return nil, ErrUnsupported
	}
	var out bytes.Buffer
	for _, bl := range e.Blocks {
		raw, err := read(bl[0], bl[1]-bl[0])
		if err != nil {
			return nil, err
		}
		zr, err := zlib.NewReader(bytes.NewReader(raw))
		if err != nil {
			return nil, ErrFormat
		}
		_, err = io.Copy(&out, io.LimitReader(zr, MaxFileData+1-int64(out.Len())))
		_ = zr.Close()
		if err != nil || out.Len() > MaxFileData {
			return nil, ErrFormat
		}
	}
	if int64(out.Len()) != e.UncompressedSize {
		return nil, ErrFormat
	}
	return out.Bytes(), nil
}

// Methods lê os nomes dos métodos de compressão do rodapé (5 nomes de 32 bytes depois do hash, versões 8+).
func Methods(tail []byte) []string {
	var magic [4]byte
	magic[0], magic[1], magic[2], magic[3] = 0xE1, 0x12, 0x6F, 0x5A
	i := strings.LastIndex(string(tail), string(magic[:]))
	start := i + 4 + 4 + 8 + 8 + 20
	if i < 0 || start > len(tail) {
		return nil
	}
	var out []string
	for k := 0; k < 5 && start+32*(k+1) <= len(tail); k++ {
		name := tail[start+32*k : start+32*(k+1)]
		if j := bytes.IndexByte(name, 0); j >= 0 {
			name = name[:j]
		}
		out = append(out, string(name))
	}
	return out
}

// GameConfigFile é o arquivo do jogo com a configuração do plugin do mod.io.
const GameConfigFile = "Insurgency/Config/DefaultGame.ini"

var apiKeyRe = regexp.MustCompile(`^[0-9A-Fa-f]{32}$`)

// ModioAPIKey lê a chave de API do mod.io que o jogo usa ([/Script/Modio.ModioSettings] ApiKey, com GameId 254).
// É a chave com que o servidor troca o código de segurança: o pedido do código precisa usar a mesma (ADR-044).
func ModioAPIKey(defaultGameIni string) (string, bool) {
	const section = "/Script/Modio.ModioSettings"
	if id, ok := ini.Get(defaultGameIni, section, "GameId"); !ok || id != "254" {
		return "", false
	}
	key, ok := ini.Get(defaultGameIni, section, "ApiKey")
	if !ok || !apiKeyRe.MatchString(key) {
		return "", false
	}
	return key, true
}
