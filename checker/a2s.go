package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"net"
	"strconv"
	"time"
)

// Consulta do Steam (A2S_INFO), a mesma que a lista de servidores do jogo usa na porta de consulta.
// https://developer.valvesoftware.com/wiki/Server_queries

var a2sRequest = append([]byte{0xFF, 0xFF, 0xFF, 0xFF, 'T'}, append([]byte("Source Engine Query"), 0)...)

const (
	a2sChallenge = 0x41 // 'A': o servidor pede o desafio de volta
	a2sInfo      = 0x49 // 'I': resposta com as informações
)

// Info é o que a resposta do servidor traz de útil para a tela.
type Info struct {
	Name       string `json:"name"`
	Map        string `json:"map"`
	Players    int    `json:"players"`
	MaxPlayers int    `json:"maxPlayers"`
}

var errBadReply = errors.New("a2s: invalid reply")

// QueryInfo manda a consulta para ip:port (UDP) e espera a resposta, com o desafio quando o servidor pede.
func QueryInfo(ctx context.Context, ip net.IP, port int, timeout time.Duration) (Info, error) {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "udp", net.JoinHostPort(ip.String(), strconv.Itoa(port)))
	if err != nil {
		return Info{}, err
	}
	defer func() { _ = conn.Close() }()
	req := a2sRequest
	buf := make([]byte, 1400)
	for attempt := 0; attempt < 3; attempt++ {
		if err := conn.SetDeadline(time.Now().Add(timeout)); err != nil {
			return Info{}, err
		}
		if _, err := conn.Write(req); err != nil {
			return Info{}, err
		}
		n, err := conn.Read(buf)
		if err != nil {
			return Info{}, err
		}
		reply := buf[:n]
		if len(reply) < 5 || !bytes.Equal(reply[:4], []byte{0xFF, 0xFF, 0xFF, 0xFF}) {
			return Info{}, errBadReply
		}
		switch reply[4] {
		case a2sChallenge:
			if len(reply) < 9 {
				return Info{}, errBadReply
			}
			req = append(append([]byte(nil), a2sRequest...), reply[5:9]...)
		case a2sInfo:
			return parseInfo(reply[5:])
		default:
			return Info{}, errBadReply
		}
	}
	return Info{}, errBadReply
}

// parseInfo lê: protocolo, nome, mapa, pasta, jogo (textos terminados em zero), id (2 bytes), jogadores, máximo.
func parseInfo(b []byte) (Info, error) {
	if len(b) < 1 {
		return Info{}, errBadReply
	}
	b = b[1:]
	var fields [4]string
	for i := range fields {
		end := bytes.IndexByte(b, 0)
		if end < 0 {
			return Info{}, errBadReply
		}
		fields[i] = string(b[:end])
		b = b[end+1:]
	}
	if len(b) < 4 {
		return Info{}, errBadReply
	}
	_ = binary.LittleEndian.Uint16(b) // id do jogo
	return Info{Name: clip(fields[0], 64), Map: clip(fields[1], 64), Players: int(b[2]), MaxPlayers: int(b[3])}, nil
}

// clip corta o texto e tira caracteres de controle (o nome vem do servidor testado).
func clip(s string, n int) string {
	out := make([]rune, 0, n)
	for _, r := range s {
		if len(out) == n {
			break
		}
		if r < 32 || r == 127 {
			continue
		}
		out = append(out, r)
	}
	return string(out)
}
