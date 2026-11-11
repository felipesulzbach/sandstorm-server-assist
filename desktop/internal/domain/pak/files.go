package pak

import (
	"path"
	"strings"
)

// File é um arquivo dentro do .pak: caminho a partir da pasta do jogo (com "/", sem "../"; compare com
// strings.EqualFold ou em minúsculas) e tamanho sem compressão. Dois mods com o mesmo caminho e tamanhos diferentes trazem versões diferentes do mesmo arquivo
// (ADR-039).
type File struct {
	Path string
	Size int64
}

// FooterAny é ParseFooter aceitando também as versões 8 e 9 (mods feitos para versões antigas do jogo).
func FooterAny(tail []byte, fileSize int64) (Footer, error) {
	f, err := ParseFooter(tail, fileSize)
	if err == nil {
		return f, nil
	}
	var magic [4]byte
	magic[0], magic[1], magic[2], magic[3] = 0xE1, 0x12, 0x6F, 0x5A
	i := strings.LastIndex(string(tail), string(magic[:]))
	if i < 1 || i+24 > len(tail) || tail[i-1] != 0 {
		return Footer{}, ErrFormat
	}
	r := &reader{b: tail, off: i + 4}
	f = Footer{Version: r.i32(), IndexOffset: r.i64(), IndexSize: r.i64()}
	if r.err != nil || f.Version < 8 || f.Version > 9 || f.IndexOffset < 0 || f.IndexSize <= 0 ||
		f.IndexSize > MaxIndex || f.IndexOffset > fileSize-f.IndexSize {
		return Footer{}, ErrFormat
	}
	return f, nil
}

// cleanPath junta ponto de montagem, pasta e nome e devolve o caminho a partir da pasta do jogo.
func cleanPath(mount, dir, name string) string {
	p := path.Clean(strings.ReplaceAll(mount+dir+name, `\`, "/"))
	for strings.HasPrefix(p, "../") {
		p = p[3:]
	}
	return strings.TrimPrefix(p, "/")
}

// legacyEntry lê um FPakEntry no formato antigo (índice das versões 8 e 9 e arquivos não codificados das 10 e 11)
// e devolve o tamanho sem compressão. methodSize é 4 (padrão) ou 1 (pacotes do UE 4.22, versão 8).
func (r *reader) legacyEntry(methodSize int) int64 {
	r.skip(8) // posição
	r.skip(8) // tamanho comprimido
	size := r.i64()
	var method uint32
	if methodSize == 1 {
		if r.need(1) {
			method = uint32(r.b[r.off])
			r.off++
		}
	} else {
		method = uint32(r.i32()) // #nosec G115 -- índice de compressão, só comparado com zero
	}
	r.skip(20) // SHA1
	if method != 0 {
		blocks := int(r.i32())
		if blocks < 0 || blocks > maxEntries {
			r.err = ErrFormat
			return 0
		}
		r.skip(16 * blocks)
	}
	r.skip(1 + 4) // flags, tamanho do bloco de compressão
	if size < 0 {
		r.err = ErrFormat
	}
	return size
}

// LegacyFiles lê o índice das versões 8 e 9: ponto de montagem, número de arquivos e, para cada um, o nome e o
// FPakEntry. Tenta o índice de compressão com 4 bytes e, se não fechar, com 1 byte (UE 4.22).
func LegacyFiles(index []byte) ([]File, error) {
	var lastErr error
	for _, methodSize := range []int{4, 1} {
		files, err := legacyFiles(index, methodSize)
		if err == nil {
			return files, nil
		}
		lastErr = err
	}
	return nil, lastErr
}

func legacyFiles(index []byte, methodSize int) ([]File, error) {
	r := &reader{b: index}
	mount := r.str()
	n := int(r.i32())
	if r.err != nil || n < 0 || n > maxEntries {
		return nil, ErrFormat
	}
	files := make([]File, 0, n)
	for i := 0; i < n && r.err == nil; i++ {
		name := r.str()
		size := r.legacyEntry(methodSize)
		if r.err == nil {
			files = append(files, File{Path: cleanPath(mount, "", name), Size: size})
		}
	}
	if r.err != nil || r.off != len(r.b) {
		return nil, ErrFormat
	}
	return files, nil
}

// ModernFiles lê os arquivos das versões 10 e 11: o índice principal (ponto de montagem e entradas codificadas) e o
// índice de pastas (pasta → arquivos → posição da entrada).
func ModernFiles(primary, dirIndex []byte) ([]File, error) {
	r := &reader{b: primary}
	mount := r.str()
	r.i32()   // número de arquivos
	r.skip(8) // semente do hash dos caminhos
	if r.i32() != 0 {
		r.skip(8 + 8 + 20)
	}
	if r.i32() != 0 {
		r.skip(8 + 8 + 20)
	}
	encLen := int(r.i32())
	if r.err != nil || encLen < 0 || !r.need(encLen) {
		return nil, ErrFormat
	}
	encoded := r.b[r.off : r.off+encLen]
	r.off += encLen
	// Entradas que não couberam no formato codificado, no formato antigo.
	var plain []int64
	if r.off < len(r.b) {
		n := int(r.i32())
		if n < 0 || n > maxEntries {
			return nil, ErrFormat
		}
		for i := 0; i < n && r.err == nil; i++ {
			plain = append(plain, r.legacyEntry(4))
		}
	}
	if r.err != nil {
		return nil, ErrFormat
	}

	d := &reader{b: dirIndex}
	dirs := int(d.i32())
	if d.err != nil || dirs < 0 || dirs > maxEntries {
		return nil, ErrFormat
	}
	var files []File
	for i := 0; i < dirs && d.err == nil; i++ {
		dir := d.str()
		n := int(d.i32())
		if n < 0 || n > maxEntries {
			return nil, ErrFormat
		}
		for j := 0; j < n && d.err == nil; j++ {
			name := d.str()
			loc := int(d.i32())
			size := int64(-1)
			switch {
			case loc >= 0:
				size = encodedSize(encoded, loc)
			case -loc-1 < len(plain):
				size = plain[-loc-1]
			}
			if size < 0 {
				return nil, ErrFormat
			}
			if len(files) >= maxEntries {
				return nil, ErrFormat
			}
			files = append(files, File{Path: cleanPath(mount, dir, name), Size: size})
		}
	}
	if d.err != nil {
		return nil, ErrFormat
	}
	return files, nil
}

// encodedSize lê o tamanho sem compressão de uma entrada codificada (UE 4.26+): um campo de bits de 32 bits, o
// tamanho do bloco (se não couber nos bits), a posição (32 ou 64 bits) e o tamanho (32 ou 64 bits). -1 se inválida.
func encodedSize(enc []byte, at int) int64 {
	r := &reader{b: enc, off: at}
	v := uint32(r.i32()) // #nosec G115 -- campo de bits do formato
	if v&0x3f == 0x3f {
		r.skip(4) // tamanho do bloco de compressão
	}
	if v&(1<<31) != 0 {
		r.skip(4)
	} else {
		r.skip(8)
	}
	var size int64
	if v&(1<<30) != 0 {
		size = int64(uint32(r.i32())) // #nosec G115 -- tamanho de 32 bits sem sinal
	} else {
		size = r.i64()
	}
	if r.err != nil || size < 0 {
		return -1
	}
	return size
}

// ContentsOf separa mutators, mapas e cenários de uma lista de arquivos (o mesmo que Scan faz no índice de pastas;
// usado para os .pak antigos).
func ContentsOf(files []File) Contents {
	var c Contents
	seen := map[string]bool{}
	add := func(list *[]string, kind, name string) {
		if !seen[kind+name] {
			seen[kind+name] = true
			*list = append(*list, name)
		}
	}
	for _, f := range files {
		dir, file := path.Split(f.Path)
		if base, ok := strings.CutSuffix(file, ".umap"); ok && nameRe.MatchString(base) {
			add(&c.Maps, "m", base)
			continue
		}
		base, ok := strings.CutSuffix(file, ".uasset")
		if !ok {
			continue
		}
		switch {
		case scenarioRe.MatchString(base):
			add(&c.Scenarios, "s", base)
		case strings.EqualFold(path.Base(strings.TrimSuffix(dir, "/")), "Mutators") && nameRe.MatchString(base):
			add(&c.Mutators, "u", base)
		}
	}
	return c
}
