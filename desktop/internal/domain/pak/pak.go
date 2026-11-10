// Package pak lê o índice de arquivos de um .pak do Unreal (versões 10 e 11, as do Insurgency: Sandstorm) para achar
// os mutators que um mod traz (ADR-023). Não faz I/O: recebe os pedaços do arquivo já lidos. Arquivo de mod é
// conteúdo de terceiros, então todo tamanho e posição é conferido antes de usar.
//
// Formato (UE 4.26+): no fim do arquivo fica o rodapé com a marca 0x5A6F12E1, a versão, a posição e o tamanho do
// índice principal. O índice principal aponta o "índice de pastas", um mapa pasta → arquivos. Um mutator de mod é
// um arquivo <Plugin>/Content/Mutators/<Nome>.uasset, e <Nome> é o que vai em -mutators=.
package pak

import (
	"encoding/binary"
	"errors"
	"regexp"
	"strings"
	"unicode/utf16"
)

// Magic é a marca do rodapé.
const Magic = 0x5A6F12E1

// Limites (um índice de pastas de mod grande tem poucos MB).
const (
	TailSize     = 512
	MaxIndex     = 16 << 20
	MaxDirIndex  = 32 << 20
	maxEntries   = 1 << 20
	maxStringLen = 1024
)

// ErrFormat: o arquivo não é um .pak que o app sabe ler (versão antiga, índice cifrado ou dados inválidos).
var ErrFormat = errors.New("pak: unsupported or invalid format")

// Footer é o que interessa do rodapé.
type Footer struct {
	Version     int32
	IndexOffset int64
	IndexSize   int64
}

// ParseFooter acha o rodapé nos últimos bytes do arquivo (tail) e confere com o tamanho do arquivo.
func ParseFooter(tail []byte, fileSize int64) (Footer, error) {
	var magic [4]byte
	binary.LittleEndian.PutUint32(magic[:], Magic)
	i := strings.LastIndex(string(tail), string(magic[:]))
	if i < 1 || i+24 > len(tail) {
		return Footer{}, ErrFormat
	}
	if tail[i-1] != 0 { // índice cifrado
		return Footer{}, ErrFormat
	}
	f := Footer{
		Version:     int32(binary.LittleEndian.Uint32(tail[i+4:])),  // #nosec G115 -- campo int32 do formato
		IndexOffset: int64(binary.LittleEndian.Uint64(tail[i+8:])),  // #nosec G115 -- campo int64 do formato
		IndexSize:   int64(binary.LittleEndian.Uint64(tail[i+16:])), // #nosec G115 -- campo int64 do formato
	}
	if f.Version < 10 || f.IndexOffset < 0 || f.IndexSize <= 0 || f.IndexSize > MaxIndex ||
		f.IndexOffset > fileSize-f.IndexSize {
		return Footer{}, ErrFormat
	}
	return f, nil
}

// reader lê campos com limite (nunca passa do fim).
type reader struct {
	b   []byte
	off int
	err error
}

func (r *reader) need(n int) bool {
	if r.err == nil && (n < 0 || r.off+n > len(r.b)) {
		r.err = ErrFormat
	}
	return r.err == nil
}

func (r *reader) i32() int32 {
	if !r.need(4) {
		return 0
	}
	v := int32(binary.LittleEndian.Uint32(r.b[r.off:])) // #nosec G115 -- campo int32 do formato
	r.off += 4
	return v
}

func (r *reader) i64() int64 {
	if !r.need(8) {
		return 0
	}
	v := int64(binary.LittleEndian.Uint64(r.b[r.off:])) // #nosec G115 -- campo int64 do formato
	r.off += 8
	return v
}

func (r *reader) skip(n int) {
	if r.need(n) {
		r.off += n
	}
}

// str lê uma FString do Unreal: tamanho com o zero final; negativo = UTF-16.
func (r *reader) str() string {
	n := int(r.i32())
	switch {
	case r.err != nil:
		return ""
	case n == 0:
		return ""
	case n > 0 && n <= maxStringLen:
		if !r.need(n) {
			return ""
		}
		s := string(r.b[r.off : r.off+n-1])
		r.off += n
		return s
	case n < 0 && -n <= maxStringLen:
		if !r.need(-2 * n) {
			return ""
		}
		u := make([]uint16, -n-1)
		for i := range u {
			u[i] = binary.LittleEndian.Uint16(r.b[r.off+2*i:])
		}
		r.off += -2 * n
		return string(utf16.Decode(u))
	}
	r.err = ErrFormat
	return ""
}

// DirIndex lê o índice principal e devolve onde está o índice de pastas.
func DirIndex(primary []byte, fileSize int64) (offset, size int64, err error) {
	r := &reader{b: primary}
	r.str()   // ponto de montagem
	r.i32()   // número de arquivos
	r.skip(8) // semente do hash dos caminhos
	if r.i32() != 0 {
		r.skip(8 + 8 + 20) // índice por hash: posição, tamanho, SHA1
	}
	if r.i32() == 0 {
		return 0, 0, ErrFormat
	}
	offset, size = r.i64(), r.i64()
	if r.err != nil || offset < 0 || size <= 0 || size > MaxDirIndex || offset > fileSize-size {
		return 0, 0, ErrFormat
	}
	return offset, size, nil
}

// nameRe é o formato aceito para um nome de mutator (o mesmo da linha de comando).
var nameRe = regexp.MustCompile(`^[A-Za-z0-9_]{1,64}$`)

// Mutators lê o índice de pastas e devolve os nomes dos mutators (arquivos .uasset em pastas "Mutators").
func Mutators(dirIndex []byte) ([]string, error) {
	r := &reader{b: dirIndex}
	dirs := int(r.i32())
	if r.err != nil || dirs < 0 || dirs > maxEntries {
		return nil, ErrFormat
	}
	seen := map[string]bool{}
	var out []string
	for d := 0; d < dirs && r.err == nil; d++ {
		dir := strings.TrimSuffix(r.str(), "/")
		files := int(r.i32())
		if files < 0 || files > maxEntries {
			return nil, ErrFormat
		}
		isMut := strings.EqualFold(dir[strings.LastIndex(dir, "/")+1:], "Mutators")
		for i := 0; i < files && r.err == nil; i++ {
			name := r.str()
			r.i32() // posição do arquivo na lista
			base, ok := strings.CutSuffix(name, ".uasset")
			if isMut && ok && nameRe.MatchString(base) && !seen[base] {
				seen[base] = true
				out = append(out, base)
			}
		}
	}
	if r.err != nil {
		return nil, r.err
	}
	return out, nil
}

// PluginName lê o nome do plugin no metadata_blob do mod.io ("PluginName:X,ContentDirectory:..."). Serve quando o
// .pak não pode ser lido e o mod tem a tag "Mutator" (o mutator costuma ter o nome do plugin).
func PluginName(blob string) string {
	for _, part := range strings.Split(blob, ",") {
		if v, ok := strings.CutPrefix(strings.TrimSpace(part), "PluginName:"); ok && nameRe.MatchString(v) {
			return v
		}
	}
	return ""
}
