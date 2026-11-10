package main

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

type probeLog struct {
	queried, tcp []string
	queryOK      bool
	tcpOpen      bool
}

func (p *probeLog) probes() Probes {
	return Probes{
		Query: func(_ context.Context, ip net.IP, port int) (Info, error) {
			p.queried = append(p.queried, net.JoinHostPort(ip.String(), itoa(port)))
			if !p.queryOK {
				return Info{}, errors.New("timeout")
			}
			return Info{Name: "Test", Map: "Farmhouse", Players: 1, MaxPlayers: 8}, nil
		},
		TCP: func(_ context.Context, ip net.IP, port int) bool {
			p.tcp = append(p.tcp, net.JoinHostPort(ip.String(), itoa(port)))
			return p.tcpOpen
		},
	}
}

func itoa(n int) string { return strconv.Itoa(n) }

func post(s *Server, body, xff string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(http.MethodPost, "/check", strings.NewReader(body))
	r.RemoteAddr = "172.18.0.2:5555" // o Caddy, na rede do Docker
	if xff != "" {
		r.Header.Set("X-Forwarded-For", xff)
	}
	w := httptest.NewRecorder()
	s.ServeHTTP(w, r)
	return w
}

func TestCheckTestsOnlyTheCaller(t *testing.T) {
	p := &probeLog{queryOK: true}
	s := NewServer(p.probes(), NewLimiter(10, 100))
	w := post(s, `{"queryPort":27131,"rconPort":27015}`, "203.0.113.7")
	if w.Code != http.StatusOK {
		t.Fatalf("status %d %s", w.Code, w.Body)
	}
	var res Response
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	if res.IP != "203.0.113.7" || !res.Query.Open || res.Query.Info.Map != "Farmhouse" || !res.Rcon.Tested || res.Rcon.Exposed {
		t.Fatalf("res = %+v", res)
	}
	if len(p.queried) != 1 || p.queried[0] != "203.0.113.7:27131" || p.tcp[0] != "203.0.113.7:27015" {
		t.Fatalf("probed %v %v", p.queried, p.tcp)
	}
}

func TestCheckRefuses(t *testing.T) {
	p := &probeLog{}
	s := NewServer(p.probes(), NewLimiter(10, 100))
	for _, c := range []struct{ body, xff, want string }{
		{`{"queryPort":27131}`, "192.168.0.10", "private-ip"},
		{`{"queryPort":27131}`, "127.0.0.1", "private-ip"},
		{`{"queryPort":27131}`, "2001:db8::1", "ipv6"},
		{`{"queryPort":80}`, "203.0.113.7", "bad-request"},
		{`{"queryPort":27131,"rconPort":70000}`, "203.0.113.7", "bad-request"},
		// Não dá para escolher outro alvo: campos a mais são recusados.
		{`{"queryPort":27131,"ip":"198.51.100.1"}`, "203.0.113.7", "bad-request"},
		{`not json`, "203.0.113.7", "bad-request"},
	} {
		w := post(s, c.body, c.xff)
		if w.Code != http.StatusBadRequest || !strings.Contains(w.Body.String(), c.want) {
			t.Fatalf("%s from %s: %d %s", c.body, c.xff, w.Code, w.Body)
		}
	}
	if len(p.queried) != 0 {
		t.Fatalf("refused requests must not probe: %v", p.queried)
	}
}

func TestCallerIPComesFromProxy(t *testing.T) {
	p := &probeLog{}
	s := NewServer(p.probes(), NewLimiter(10, 100))
	// O último valor é o que o Caddy pôs; o primeiro poderia ter vindo do cliente.
	post(s, `{"queryPort":27131}`, "198.51.100.1, 203.0.113.9")
	if len(p.queried) != 1 || p.queried[0] != "203.0.113.9:27131" {
		t.Fatalf("probed %v", p.queried)
	}
}

func TestRateLimit(t *testing.T) {
	p := &probeLog{}
	s := NewServer(p.probes(), NewLimiter(2, 100))
	for i := 0; i < 2; i++ {
		if w := post(s, `{"queryPort":27131}`, "203.0.113.7"); w.Code != http.StatusOK {
			t.Fatalf("request %d: %d", i, w.Code)
		}
	}
	if w := post(s, `{"queryPort":27131}`, "203.0.113.7"); w.Code != http.StatusTooManyRequests {
		t.Fatalf("third request: %d", w.Code)
	}
	// Outro IP segue livre.
	if w := post(s, `{"queryPort":27131}`, "203.0.113.8"); w.Code != http.StatusOK {
		t.Fatalf("other ip: %d", w.Code)
	}
	l := NewLimiter(1, 1)
	now := time.Date(2026, 11, 10, 3, 0, 10, 0, time.UTC)
	if !l.Allow("a", now) || l.Allow("b", now) || !l.Allow("b", now.Add(time.Minute)) {
		t.Fatal("global limit and window reset")
	}
}

func TestHealthAndRoutes(t *testing.T) {
	s := NewServer((&probeLog{}).probes(), NewLimiter(1, 1))
	w := httptest.NewRecorder()
	s.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("healthz %d", w.Code)
	}
	w = httptest.NewRecorder()
	s.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/check", nil))
	if w.Code != http.StatusNotFound {
		t.Fatalf("GET /check %d", w.Code)
	}
}

// Servidor do jogo falso na máquina: responde a consulta do Steam com desafio, como o servidor real.
func TestQueryInfoWithChallenge(t *testing.T) {
	pc, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = pc.Close() }()
	go func() {
		buf := make([]byte, 1400)
		for {
			n, addr, err := pc.ReadFrom(buf)
			if err != nil {
				return
			}
			req := buf[:n]
			if len(req) == len(a2sRequest) { // sem desafio: pede o desafio
				_, _ = pc.WriteTo([]byte{0xFF, 0xFF, 0xFF, 0xFF, a2sChallenge, 1, 2, 3, 4}, addr)
				continue
			}
			reply := []byte{0xFF, 0xFF, 0xFF, 0xFF, a2sInfo, 17}
			for _, s := range []string{"My\x01 Server", "Containment", "Insurgency", "Insurgency: Sandstorm"} {
				reply = append(append(reply, s...), 0)
			}
			reply = append(reply, 0, 0, 3, 8)
			_, _ = pc.WriteTo(reply, addr)
		}
	}()
	port := pc.LocalAddr().(*net.UDPAddr).Port
	info, err := QueryInfo(context.Background(), net.ParseIP("127.0.0.1"), port, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if info.Name != "My Server" || info.Map != "Containment" || info.Players != 3 || info.MaxPlayers != 8 {
		t.Fatalf("info = %+v", info)
	}
	if _, err := parseInfo([]byte{17, 'a', 0}); err == nil {
		t.Fatal("truncated reply accepted")
	}
}

// Teste com um servidor real (só leitura): SSA_A2S_TEST=127.0.0.1:27131 go test -run Real ./...
func TestRealServerQuery(t *testing.T) {
	target := os.Getenv("SSA_A2S_TEST")
	if target == "" {
		t.Skip("SSA_A2S_TEST não definido")
	}
	host, port, err := net.SplitHostPort(target)
	if err != nil {
		t.Fatal(err)
	}
	p, _ := strconv.Atoi(port)
	info, err := QueryInfo(context.Background(), net.ParseIP(host), p, 3*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("map=%s players=%d/%d nameLen=%d", info.Map, info.Players, info.MaxPlayers, len(info.Name))
}
