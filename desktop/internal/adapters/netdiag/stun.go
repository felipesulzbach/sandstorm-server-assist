// Package netdiag mede, do próprio PC, o que o diagnóstico de portas precisa (ADR-033): o IP público (STUN), o
// roteador (UPnP), a consulta ao servidor (A2S) e a lista da Steam. As regras ficam em domain/netcheck.
package netdiag

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"net"
	"time"

	"sandstorm-server-assist/desktop/internal/domain/netcheck"
)

// stunServers são servidores STUN públicos e gratuitos. Eles só devolvem o IP e a porta de onde o pedido veio.
var stunServers = []string{"stun.cloudflare.com:3478", "stun.l.google.com:19302"}

const (
	stunCookie       = 0x2112A442
	stunBindRequest  = 0x0001
	stunBindResponse = 0x0101
	attrMapped       = 0x0001
	attrXorMapped    = 0x0020
)

var errSTUN = errors.New("stun: no valid answer")

// PublicIP pergunta aos servidores STUN qual IPv4 a internet vê. Devolve nil se nenhum responder.
func PublicIP(ctx context.Context) net.IP {
	for _, server := range stunServers {
		if ip, err := stunQuery(ctx, server); err == nil {
			return ip
		}
	}
	return nil
}

func stunQuery(ctx context.Context, server string) (net.IP, error) {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "udp4", server)
	if err != nil {
		return nil, err
	}
	defer func() { _ = conn.Close() }()

	var id [12]byte
	if _, err := rand.Read(id[:]); err != nil {
		return nil, err
	}
	req := make([]byte, 20)
	binary.BigEndian.PutUint16(req[0:], stunBindRequest)
	binary.BigEndian.PutUint32(req[4:], stunCookie)
	copy(req[8:], id[:])

	buf := make([]byte, 1024)
	for attempt := 0; attempt < 2; attempt++ {
		if err := conn.SetDeadline(time.Now().Add(1500 * time.Millisecond)); err != nil {
			return nil, err
		}
		if _, err := conn.Write(req); err != nil {
			return nil, err
		}
		n, err := conn.Read(buf)
		if err != nil {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			continue
		}
		if ip := parseSTUN(buf[:n], id); ip != nil {
			return ip, nil
		}
	}
	return nil, errSTUN
}

// parseSTUN lê a resposta de Binding e devolve o IPv4 público (só se a resposta for deste pedido e o IP for da
// internet).
func parseSTUN(b []byte, id [12]byte) net.IP {
	if len(b) < 20 || binary.BigEndian.Uint16(b[0:]) != stunBindResponse ||
		binary.BigEndian.Uint32(b[4:]) != stunCookie || !bytes.Equal(b[8:20], id[:]) {
		return nil
	}
	length := int(binary.BigEndian.Uint16(b[2:]))
	if 20+length > len(b) {
		return nil
	}
	attrs := b[20 : 20+length]
	var found net.IP
	for len(attrs) >= 4 {
		typ := binary.BigEndian.Uint16(attrs[0:])
		alen := int(binary.BigEndian.Uint16(attrs[2:]))
		if 4+alen > len(attrs) {
			return nil
		}
		v := attrs[4 : 4+alen]
		if len(v) >= 8 && v[1] == 0x01 { // família IPv4
			ip := net.IPv4(v[4], v[5], v[6], v[7]).To4()
			switch typ {
			case attrXorMapped:
				var cookie [4]byte
				binary.BigEndian.PutUint32(cookie[:], stunCookie)
				for i := range ip {
					ip[i] ^= cookie[i]
				}
				found = ip
			case attrMapped:
				if found == nil {
					found = ip
				}
			}
		}
		step := 4 + (alen+3)/4*4
		if step > len(attrs) {
			break
		}
		attrs = attrs[step:]
	}
	if found == nil || netcheck.NotPublic(found) {
		return nil
	}
	return found
}
