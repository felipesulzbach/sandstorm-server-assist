package netdiag

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"sandstorm-server-assist/desktop/internal/adapters/netx"
	"sandstorm-server-assist/desktop/internal/domain/netcheck"
)

// Consulta do Steam (A2S_INFO), a mesma que a lista de servidores do jogo usa na porta de consulta.
// https://developer.valvesoftware.com/wiki/Server_queries
var a2sRequest = append([]byte{0xFF, 0xFF, 0xFF, 0xFF, 'T'}, append([]byte("Source Engine Query"), 0)...)

const (
	a2sChallenge = 0x41 // 'A': o servidor pede o desafio de volta
	a2sInfo      = 0x49 // 'I': resposta com as informações
)

// Answers diz se o servidor em ip:port responde à consulta (com o desafio, quando o servidor pede).
func Answers(ctx context.Context, ip net.IP, port int) bool {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "udp4", net.JoinHostPort(ip.String(), strconv.Itoa(port)))
	if err != nil {
		return false
	}
	defer func() { _ = conn.Close() }()
	req := a2sRequest
	buf := make([]byte, 1400)
	for attempt := 0; attempt < 4; attempt++ {
		if err := conn.SetDeadline(time.Now().Add(1500 * time.Millisecond)); err != nil {
			return false
		}
		if _, err := conn.Write(req); err != nil {
			return false
		}
		n, err := conn.Read(buf)
		if err != nil {
			if ctx.Err() != nil {
				return false
			}
			continue
		}
		reply := buf[:n]
		if len(reply) < 5 || !bytes.Equal(reply[:4], []byte{0xFF, 0xFF, 0xFF, 0xFF}) {
			return false
		}
		switch reply[4] {
		case a2sChallenge:
			if len(reply) < 9 {
				return false
			}
			req = append(append([]byte(nil), a2sRequest...), reply[5:9]...)
		case a2sInfo:
			return true
		default:
			return false
		}
	}
	return false
}

// steamAPI lista os servidores que a Steam conhece num IP (API pública, sem chave).
const steamAPI = "https://api.steampowered.com/ISteamApps/GetServersAtAddress/v1/"

type steamReply struct {
	Response struct {
		Success bool `json:"success"`
		Servers []struct {
			Addr string `json:"addr"`
		} `json:"servers"`
	} `json:"response"`
}

// SteamListed pergunta à Steam se há um servidor neste IP com a porta de consulta.
func SteamListed(ctx context.Context, ip net.IP, queryPort int) netcheck.Tri {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	q := url.Values{"addr": {ip.String()}, "format": {"json"}}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, steamAPI+"?"+q.Encode(), nil)
	if err != nil {
		return netcheck.Unknown
	}
	resp, err := netx.Do(req)
	if err != nil {
		return netcheck.Unknown
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return netcheck.Unknown
	}
	var r steamReply
	if err := json.NewDecoder(io.LimitReader(resp.Body, 256<<10)).Decode(&r); err != nil || !r.Response.Success {
		return netcheck.Unknown
	}
	want := net.JoinHostPort(ip.String(), strconv.Itoa(queryPort))
	for _, s := range r.Response.Servers {
		if s.Addr == want {
			return netcheck.Yes
		}
	}
	return netcheck.No
}
