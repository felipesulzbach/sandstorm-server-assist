// Package rcon conecta no RCON do servidor desta máquina (127.0.0.1) para pedir a lista de mapas e trocar o mapa
// (ADR-019). Uma conexão por ação do usuário: o RCON do Sandstorm deixa uma thread presa no servidor a cada
// conexão nova, então o app nunca consulta em intervalos.
package rcon

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	proto "sandstorm-server-assist/desktop/internal/domain/rcon"
)

// Tempos: conectar, esperar a primeira resposta e, depois dela, quanto esperar por mais pedaços.
const (
	dialTimeout  = 3 * time.Second
	replyTimeout = 8 * time.Second
	idleTimeout  = 400 * time.Millisecond
	// maxReply limita o total lido de uma resposta.
	maxReply = 4 << 20
)

// ErrAuth: o servidor recusou a senha.
var ErrAuth = errors.New("rcon: password refused")

// Client é uma conexão autenticada.
type Client struct {
	conn net.Conn
	id   int32
}

// Dial conecta em 127.0.0.1:port e autentica.
func Dial(port int, password string) (*Client, error) {
	if port < 1 || port > 65535 {
		return nil, fmt.Errorf("rcon: invalid port %d", port)
	}
	conn, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(port)), dialTimeout)
	if err != nil {
		return nil, err
	}
	c := &Client{conn: conn, id: 1}
	if err := c.auth(password); err != nil {
		_ = conn.Close()
		return nil, err
	}
	return c, nil
}

func (c *Client) send(p proto.Packet) error {
	b, err := proto.Encode(p)
	if err != nil {
		return err
	}
	if err := c.conn.SetWriteDeadline(time.Now().Add(replyTimeout)); err != nil {
		return err
	}
	_, err = c.conn.Write(b)
	return err
}

func (c *Client) read(timeout time.Duration) (proto.Packet, error) {
	if err := c.conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		return proto.Packet{}, err
	}
	return proto.Decode(c.conn)
}

// auth manda a senha. O servidor pode mandar um pacote vazio antes da resposta da autenticação (como o Source).
func (c *Client) auth(password string) error {
	if err := c.send(proto.Packet{ID: c.id, Type: proto.TypeAuth, Body: password}); err != nil {
		return err
	}
	deadline := time.Now().Add(replyTimeout)
	for time.Now().Before(deadline) {
		p, err := c.read(time.Until(deadline))
		if err != nil {
			return err
		}
		if p.Type != proto.TypeAuthResponse {
			continue
		}
		if p.ID == proto.AuthFailedID {
			return ErrAuth
		}
		return nil
	}
	return errors.New("rcon: no answer to the password")
}

// Exec roda um comando e devolve a resposta inteira (o servidor pode dividir em vários pacotes).
func (c *Client) Exec(cmd string) (string, error) {
	c.id++
	if err := c.send(proto.Packet{ID: c.id, Type: proto.TypeExec, Body: cmd}); err != nil {
		return "", err
	}
	var b strings.Builder
	timeout := replyTimeout
	got := false
	for {
		p, err := c.read(timeout)
		if err != nil {
			var ne net.Error
			if got && errors.As(err, &ne) && ne.Timeout() {
				return b.String(), nil // acabaram os pedaços
			}
			return b.String(), err
		}
		if p.Type != proto.TypeResponse {
			continue
		}
		if p.ID == proto.AuthFailedID {
			return "", ErrAuth
		}
		got = true
		timeout = idleTimeout
		if b.Len()+len(p.Body) > maxReply {
			return b.String(), errors.New("rcon: reply too large")
		}
		b.WriteString(p.Body)
	}
}

// Close fecha a conexão.
func (c *Client) Close() error { return c.conn.Close() }
