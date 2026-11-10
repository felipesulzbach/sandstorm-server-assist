// Verificador de portas do Sandstorm Server Assist (ADR-021).
//
// Roda no servidor da página de download, atrás do Caddy (rota /sandstorm/api/). O app pede "teste minhas portas"
// e este serviço testa, de fora, o IP de quem pediu:
//   - porta de consulta (UDP): manda a consulta do Steam (A2S_INFO) e espera a resposta do servidor do jogo;
//   - porta do RCON (TCP), se informada: tenta conectar. Aberta para a internet é um risco, então o app avisa.
//
// Só testa o IP de origem do pedido (nunca um endereço escolhido por quem pede), um pedido por vez por IP, com
// limite por minuto. Não grava nada e não registra IPs.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	defaultAddr  = ":8095"
	probeTimeout = 3 * time.Second
	perIPPerMin  = 6
	allPerMin    = 120
	maxBody      = 1 << 10
)

// Request é o que o app manda: só números de porta.
type Request struct {
	QueryPort int `json:"queryPort"`
	RconPort  int `json:"rconPort"`
}

// Response é o resultado do teste.
type Response struct {
	IP    string      `json:"ip"`
	Query QueryResult `json:"query"`
	Rcon  RconResult  `json:"rcon"`
}

// QueryResult diz se a porta de consulta respondeu de fora.
type QueryResult struct {
	Port int   `json:"port"`
	Open bool  `json:"open"`
	Info *Info `json:"info,omitempty"`
}

// RconResult diz se a porta do RCON aceita conexões de fora (não deveria).
type RconResult struct {
	Port    int  `json:"port"`
	Tested  bool `json:"tested"`
	Exposed bool `json:"exposed"`
}

// Probes são os testes de rede (dublês nos testes).
type Probes struct {
	Query func(ctx context.Context, ip net.IP, port int) (Info, error)
	TCP   func(ctx context.Context, ip net.IP, port int) bool
}

func realProbes() Probes {
	return Probes{
		Query: func(ctx context.Context, ip net.IP, port int) (Info, error) {
			return QueryInfo(ctx, ip, port, probeTimeout)
		},
		TCP: func(ctx context.Context, ip net.IP, port int) bool {
			d := net.Dialer{Timeout: probeTimeout}
			c, err := d.DialContext(ctx, "tcp", net.JoinHostPort(ip.String(), strconv.Itoa(port)))
			if err != nil {
				return false
			}
			_ = c.Close()
			return true
		},
	}
}

// Server é o serviço HTTP.
type Server struct {
	probes Probes
	limit  *Limiter
	// busy: um teste por vez para cada IP.
	mu   sync.Mutex
	busy map[string]bool
}

// NewServer cria o serviço.
func NewServer(p Probes, l *Limiter) *Server {
	return &Server{probes: p, limit: l, busy: map[string]bool{}}
}

type apiError struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// clientIP é o IP de quem pediu. Atrás do Caddy, é o último do X-Forwarded-For (o Caddy troca o que o cliente
// mandou); sem o cabeçalho, o endereço da conexão.
func clientIP(r *http.Request) net.IP {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return net.ParseIP(strings.TrimSpace(parts[len(parts)-1]))
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return nil
	}
	return net.ParseIP(host)
}

// publicIP diz se dá para testar o IP pela internet.
func publicIP(ip net.IP) bool {
	return ip != nil && ip.IsGlobalUnicast() && !ip.IsPrivate()
}

func validPort(p int) bool { return p >= 1024 && p <= 65535 }

// ServeHTTP: POST /check (JSON) e GET /healthz.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/healthz" && r.Method == http.MethodGet:
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	case r.URL.Path == "/check" && r.Method == http.MethodPost:
		s.check(w, r)
	default:
		writeJSON(w, http.StatusNotFound, apiError{"not-found"})
	}
}

func (s *Server) check(w http.ResponseWriter, r *http.Request) {
	ip := clientIP(r)
	if ip == nil {
		writeJSON(w, http.StatusBadRequest, apiError{"no-ip"})
		return
	}
	if !publicIP(ip) {
		writeJSON(w, http.StatusBadRequest, apiError{"private-ip"})
		return
	}
	if ip.To4() == nil {
		// O servidor do jogo escuta em IPv4: um teste por IPv6 daria "fechada" sem ser verdade.
		writeJSON(w, http.StatusBadRequest, apiError{"ipv6"})
		return
	}
	var req Request
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBody))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil || !validPort(req.QueryPort) || (req.RconPort != 0 && !validPort(req.RconPort)) {
		writeJSON(w, http.StatusBadRequest, apiError{"bad-request"})
		return
	}
	key := ip.String()
	if !s.limit.Allow(key, time.Now()) {
		writeJSON(w, http.StatusTooManyRequests, apiError{"rate-limited"})
		return
	}
	s.mu.Lock()
	if s.busy[key] {
		s.mu.Unlock()
		writeJSON(w, http.StatusTooManyRequests, apiError{"busy"})
		return
	}
	s.busy[key] = true
	s.mu.Unlock()
	defer func() {
		s.mu.Lock()
		delete(s.busy, key)
		s.mu.Unlock()
	}()

	ctx, cancel := context.WithTimeout(r.Context(), 4*probeTimeout)
	defer cancel()
	res := Response{IP: key, Query: QueryResult{Port: req.QueryPort}, Rcon: RconResult{Port: req.RconPort}}
	if info, err := s.probes.Query(ctx, ip, req.QueryPort); err == nil {
		res.Query.Open, res.Query.Info = true, &info
	}
	if req.RconPort != 0 {
		res.Rcon.Tested = true
		res.Rcon.Exposed = s.probes.TCP(ctx, ip, req.RconPort)
	}
	writeJSON(w, http.StatusOK, res)
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		os.Exit(healthcheck())
	}
	addr := os.Getenv("SSA_CHECK_ADDR")
	if addr == "" {
		addr = defaultAddr
	}
	srv := &http.Server{
		Addr:              addr,
		Handler:           NewServer(realProbes(), NewLimiter(perIPPerMin, allPerMin)),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		_, _ = os.Stderr.WriteString("checker: " + err.Error() + "\n")
		os.Exit(1)
	}
}

// healthcheck é usado pelo Docker (a imagem não tem curl).
func healthcheck() int {
	c := http.Client{Timeout: 3 * time.Second}
	resp, err := c.Get("http://127.0.0.1" + defaultAddr + "/healthz")
	if err != nil {
		return 1
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 1
	}
	return 0
}
