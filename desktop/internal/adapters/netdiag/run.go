package netdiag

import (
	"context"
	"net"
	"sync"
	"time"

	"sandstorm-server-assist/desktop/internal/domain/netcheck"
)

// Input é o que o diagnóstico precisa saber do perfil e do app.
type Input struct {
	GamePort, QueryPort, RconPort int
	Rcon                          bool
	ServerRunning                 bool
	LocalIPs                      []string
	Firewall                      netcheck.Tri
}

// timeout cobre o diagnóstico inteiro (roteador, STUN, consultas e Steam).
const timeout = 15 * time.Second

// Run faz as medidas (em paralelo quando dá) e devolve o resultado das regras.
func Run(ctx context.Context, in Input) netcheck.Result {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	f := netcheck.Facts{
		GamePort: in.GamePort, QueryPort: in.QueryPort, RconPort: in.RconPort, Rcon: in.Rcon,
		ServerRunning: in.ServerRunning, LocalIPs: in.LocalIPs, Firewall: in.Firewall,
	}

	var wg sync.WaitGroup
	var public net.IP
	wg.Add(3)
	go func() {
		defer wg.Done()
		public = PublicIP(ctx)
	}()
	go func() {
		defer wg.Done()
		f.Router = router(ctx, in)
	}()
	go func() {
		defer wg.Done()
		if in.ServerRunning {
			f.LocalAnswer = Answers(ctx, net.IPv4(127, 0, 0, 1), in.QueryPort)
		}
	}()
	wg.Wait()

	if public != nil {
		f.PublicIP = public.String()
		if in.ServerRunning {
			wg.Add(2)
			go func() {
				defer wg.Done()
				f.Loopback = Answers(ctx, public, in.QueryPort)
			}()
			go func() {
				defer wg.Done()
				f.Steam = SteamListed(ctx, public, in.QueryPort)
			}()
			wg.Wait()
		}
	}
	return netcheck.Evaluate(f)
}

// router lê do roteador o IP externo e os redirecionamentos das portas do servidor.
func router(ctx context.Context, in Input) netcheck.Router {
	g := Discover(ctx)
	if g == nil {
		return netcheck.Router{}
	}
	r := netcheck.Router{Found: true, ExternalIP: g.ExternalIP(ctx)}
	r.Game = g.Mapping(ctx, in.GamePort, "UDP")
	r.Query = g.Mapping(ctx, in.QueryPort, "UDP")
	if in.Rcon {
		r.Rcon = g.Mapping(ctx, in.RconPort, "TCP")
	}
	return r
}
