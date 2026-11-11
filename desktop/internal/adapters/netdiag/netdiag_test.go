package netdiag

import (
	"context"
	"encoding/binary"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

func stunReply(id [12]byte, ip net.IP, xor bool) []byte {
	b := make([]byte, 20+12)
	binary.BigEndian.PutUint16(b[0:], stunBindResponse)
	binary.BigEndian.PutUint16(b[2:], 12)
	binary.BigEndian.PutUint32(b[4:], stunCookie)
	copy(b[8:], id[:])
	typ := uint16(attrMapped)
	addr := append(net.IP(nil), ip.To4()...)
	if xor {
		typ = attrXorMapped
		var c [4]byte
		binary.BigEndian.PutUint32(c[:], stunCookie)
		for i := range addr {
			addr[i] ^= c[i]
		}
	}
	binary.BigEndian.PutUint16(b[20:], typ)
	binary.BigEndian.PutUint16(b[22:], 8)
	b[25] = 0x01
	copy(b[28:], addr)
	return b
}

func TestParseSTUN(t *testing.T) {
	id := [12]byte{1, 2, 3}
	if ip := parseSTUN(stunReply(id, net.IPv4(203, 0, 113, 7), true), id); ip.String() != "203.0.113.7" {
		t.Fatalf("xor: %v", ip)
	}
	if ip := parseSTUN(stunReply(id, net.IPv4(203, 0, 113, 8), false), id); ip.String() != "203.0.113.8" {
		t.Fatalf("plain: %v", ip)
	}
	other := [12]byte{9}
	if ip := parseSTUN(stunReply(id, net.IPv4(203, 0, 113, 7), true), other); ip != nil {
		t.Fatal("accepted another transaction")
	}
	if ip := parseSTUN(stunReply(id, net.IPv4(192, 168, 0, 2), true), id); ip != nil {
		t.Fatal("accepted a private address")
	}
	if ip := parseSTUN([]byte{1, 1, 0, 0}, id); ip != nil {
		t.Fatal("accepted a short reply")
	}
}

func TestLanURL(t *testing.T) {
	from := net.IPv4(192, 168, 0, 1)
	if _, err := lanURL("http://192.168.0.1:5000/desc.xml", from); err != nil {
		t.Fatal(err)
	}
	for _, raw := range []string{
		"https://192.168.0.1/desc.xml",      // só http na rede local
		"http://8.8.8.8/desc.xml",           // fora da rede local
		"http://192.168.0.2/desc.xml",       // outro aparelho que não o que respondeu
		"http://router.local/desc.xml",      // nome, não IP
		"http://u:p@192.168.0.1/desc.xml",   // com usuário
		"http://192.168.0.1:99999/desc.xml", // porta inválida
	} {
		if _, err := lanURL(raw, from); err == nil {
			t.Errorf("accepted %s", raw)
		}
	}
}

func TestSSDPLocation(t *testing.T) {
	msg := "HTTP/1.1 200 OK\r\nCACHE-CONTROL: max-age=120\r\nLocation: http://192.168.0.1:1900/igd.xml\r\n\r\n"
	if got := ssdpLocation([]byte(msg)); got != "http://192.168.0.1:1900/igd.xml" {
		t.Fatalf("got %q", got)
	}
}

const igdXML = `<?xml version="1.0"?>
<root xmlns="urn:schemas-upnp-org:device-1-0"><device>
 <deviceType>urn:schemas-upnp-org:device:InternetGatewayDevice:1</deviceType>
 <deviceList><device><deviceList><device>
  <serviceList><service>
   <serviceType>urn:schemas-upnp-org:service:WANIPConnection:1</serviceType>
   <controlURL>/ctl/IPConn</controlURL>
  </service></serviceList>
 </device></deviceList></device></deviceList>
</device></root>`

func fakeRouter(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			_, _ = w.Write([]byte(igdXML))
			return
		}
		action := r.Header.Get("SOAPAction")
		body := make([]byte, 4096)
		n, _ := r.Body.Read(body)
		switch {
		case strings.Contains(action, "GetExternalIPAddress"):
			_, _ = w.Write([]byte(`<s:Envelope><s:Body><u:R><NewExternalIPAddress>203.0.113.7</NewExternalIPAddress></u:R></s:Body></s:Envelope>`))
		case strings.Contains(string(body[:n]), "<NewExternalPort>27131</NewExternalPort>"):
			_, _ = w.Write([]byte(`<s:Envelope><s:Body><u:R><NewInternalPort>27131</NewInternalPort>` +
				`<NewInternalClient>192.168.0.10</NewInternalClient><NewEnabled>1</NewEnabled></u:R></s:Body></s:Envelope>`))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`<s:Envelope><s:Body><s:Fault><detail><UPnPError><errorCode>714</errorCode>` +
				`</UPnPError></detail></s:Fault></s:Body></s:Envelope>`))
		}
	}))
}

func TestGatewayAgainstFakeRouter(t *testing.T) {
	srv := fakeRouter(t)
	defer srv.Close()
	// O servidor de teste fica em 127.0.0.1 (não é rede privada): chama describe direto, sem lanURL.
	loc, _ := url.Parse(srv.URL + "/igd.xml")
	g, err := describe(context.Background(), loc)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(g.control.Path, "/ctl/IPConn") || g.service != wanServices[1] {
		t.Fatalf("gateway %+v", g)
	}
	if ip := g.ExternalIP(context.Background()); ip != "203.0.113.7" {
		t.Fatalf("external ip %q", ip)
	}
	m := g.Mapping(context.Background(), 27131, "UDP")
	if !m.Known || !m.Exists || !m.Enabled || m.Client != "192.168.0.10" || m.InternalPort != 27131 {
		t.Fatalf("mapping %+v", m)
	}
	if m := g.Mapping(context.Background(), 27102, "UDP"); !m.Known || m.Exists {
		t.Fatalf("missing mapping %+v", m)
	}
}

func TestAnswersFakeServer(t *testing.T) {
	pc, err := net.ListenPacket("udp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer pc.Close()
	go func() {
		buf := make([]byte, 1400)
		for {
			n, addr, err := pc.ReadFrom(buf)
			if err != nil {
				return
			}
			if n == len(a2sRequest) { // primeiro pedido: pede o desafio
				_, _ = pc.WriteTo([]byte{0xFF, 0xFF, 0xFF, 0xFF, a2sChallenge, 1, 2, 3, 4}, addr)
				continue
			}
			_, _ = pc.WriteTo([]byte{0xFF, 0xFF, 0xFF, 0xFF, a2sInfo, 17}, addr)
		}
	}()
	port := pc.LocalAddr().(*net.UDPAddr).Port
	if !Answers(context.Background(), net.IPv4(127, 0, 0, 1), port) {
		t.Fatal("no answer")
	}
}

// TestRealNetwork roda o diagnóstico de verdade nesta rede (STUN, UPnP, Steam). Só com SSA_NETDIAG_E2E=1. Mostra só
// os códigos e as situações, nunca IPs.
func TestRealNetwork(t *testing.T) {
	if os.Getenv("SSA_NETDIAG_E2E") != "1" {
		t.Skip("set SSA_NETDIAG_E2E=1 to test on the real network")
	}
	r := Run(context.Background(), Input{GamePort: 27102, QueryPort: 27131, RconPort: 27015})
	t.Logf("verdict=%s publicIP=%v", r.Verdict, r.PublicIP != "")
	for _, it := range r.Items {
		t.Logf("  %s %s", it.Status, it.Key)
	}
}
