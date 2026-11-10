// Package rcon é o formato dos pacotes do RCON (protocolo Source RCON, o mesmo do Insurgency: Sandstorm). Só
// monta e lê pacotes: a conexão fica em adapters/rcon (ADR-019).
//
// Pacote: tamanho (int32, sem contar ele mesmo), id (int32), tipo (int32), corpo (texto) e dois bytes zero.
// Tudo em little-endian.
package rcon

import (
	"encoding/binary"
	"errors"
	"io"
	"strings"
)

// Tipos de pacote. Exec e AuthResponse têm o mesmo número no protocolo (o sentido diz qual é).
const (
	TypeResponse     int32 = 0
	TypeExec         int32 = 2
	TypeAuthResponse int32 = 2
	TypeAuth         int32 = 3
)

// AuthFailedID é o id que o servidor devolve quando a senha está errada.
const AuthFailedID int32 = -1

// MaxBody é o maior corpo aceito num pacote (a lista de cenários do servidor cabe com folga).
const MaxBody = 1 << 20

// header é id + tipo; trailer são os dois zeros do fim.
const (
	header  = 8
	trailer = 2
)

// Packet é um pacote do RCON.
type Packet struct {
	ID   int32
	Type int32
	Body string
}

// Erros de formato.
var (
	ErrBodyTooLarge = errors.New("rcon: body too large")
	ErrBodyNul      = errors.New("rcon: body has a NUL byte")
	ErrBadSize      = errors.New("rcon: invalid packet size")
)

// Encode monta os bytes de um pacote.
func Encode(p Packet) ([]byte, error) {
	if len(p.Body) > MaxBody {
		return nil, ErrBodyTooLarge
	}
	if strings.IndexByte(p.Body, 0) >= 0 {
		return nil, ErrBodyNul
	}
	size := header + len(p.Body) + trailer
	b := make([]byte, 4+size)
	binary.LittleEndian.PutUint32(b[0:], uint32(size))   // #nosec G115 -- size <= MaxBody+10
	binary.LittleEndian.PutUint32(b[4:], uint32(p.ID))   // #nosec G115 -- bits do int32 como estão
	binary.LittleEndian.PutUint32(b[8:], uint32(p.Type)) // #nosec G115 -- bits do int32 como estão
	copy(b[12:], p.Body)
	return b, nil
}

// Decode lê um pacote. Recusa tamanhos fora do protocolo antes de alocar.
func Decode(r io.Reader) (Packet, error) {
	var sz [4]byte
	if _, err := io.ReadFull(r, sz[:]); err != nil {
		return Packet{}, err
	}
	size := int64(binary.LittleEndian.Uint32(sz[:]))
	if size < header+trailer || size > header+MaxBody+trailer {
		return Packet{}, ErrBadSize
	}
	b := make([]byte, size)
	if _, err := io.ReadFull(r, b); err != nil {
		return Packet{}, err
	}
	body := b[header : size-trailer]
	// Alguns servidores mandam só um zero no fim: corta no primeiro zero.
	if i := strings.IndexByte(string(body), 0); i >= 0 {
		body = body[:i]
	}
	return Packet{
		ID:   int32(binary.LittleEndian.Uint32(b[0:])), // #nosec G115 -- bits do int32 como estão
		Type: int32(binary.LittleEndian.Uint32(b[4:])), // #nosec G115 -- bits do int32 como estão
		Body: string(body),
	}, nil
}
