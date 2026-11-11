// Package netcheck são as regras do diagnóstico de portas feito no próprio PC (ADR-033). O app mede o que dá para
// medir daqui (o servidor responde? o Firewall tem a regra? o roteador redireciona as portas? há CGNAT? a Steam
// lista o servidor?) e este pacote transforma as medidas em itens com código e situação. A tela traduz os códigos.
// Só um teste de fora da rede prova que a porta abre pela internet: a tela ensina a fazê-lo.
package netcheck

import (
	"net"
	"strconv"
)

// Tri é uma medida que pode não ter sido feita.
type Tri int

const (
	Unknown Tri = iota
	Yes
	No
)

// Mapping é um redirecionamento de porta lido do roteador (UPnP).
type Mapping struct {
	// Known: o roteador respondeu sobre esta porta (com ou sem redirecionamento).
	Known        bool
	Exists       bool
	Enabled      bool
	Client       string // IP da rede local que recebe a porta
	InternalPort int
}

// Router é o que o roteador contou por UPnP.
type Router struct {
	Found      bool   // respondeu por UPnP
	ExternalIP string // IP do lado da internet, segundo o roteador
	Game       Mapping
	Query      Mapping
	Rcon       Mapping
}

// Facts são as medidas feitas no PC.
type Facts struct {
	GamePort, QueryPort, RconPort int
	Rcon                          bool // RCON ligado no perfil
	ServerRunning                 bool
	LocalAnswer                   bool     // o servidor respondeu à consulta em 127.0.0.1
	PublicIP                      string   // IP que a internet vê (STUN); "" = não se sabe
	LocalIPs                      []string // IPv4 das placas deste PC
	Router                        Router
	Loopback                      bool // o servidor respondeu pela consulta enviada ao IP público
	Steam                         Tri  // a Steam lista um servidor neste IP com a porta de consulta
	Firewall                      Tri  // regras do app no Firewall do Windows
}

// Situação de um item.
const (
	OK   = "ok"
	Bad  = "bad"
	Warn = "warn"
	Info = "info"
)

// Conclusão geral.
const (
	VerdictOpen    = "open"    // o servidor respondeu pelo IP público: o roteador redireciona a porta
	VerdictLikely  = "likely"  // tudo o que dá para ver daqui está certo
	VerdictProblem = "problem" // achou um problema
	VerdictUnsure  = "unsure"  // não deu para confirmar daqui
)

// Item é uma linha do diagnóstico. Key é a chave do texto (server.portCheck.items.<key>).
type Item struct {
	Key    string            `json:"key"`
	Status string            `json:"status"`
	Params map[string]string `json:"params,omitempty"`
}

// Result é o que a tela mostra.
type Result struct {
	PublicIP  string `json:"publicIp"`
	GamePort  int    `json:"gamePort"`
	QueryPort int    `json:"queryPort"`
	Verdict   string `json:"verdict"`
	Items     []Item `json:"items"`
}

// Evaluate transforma as medidas em itens e na conclusão.
func Evaluate(f Facts) Result {
	r := Result{GamePort: f.GamePort, QueryPort: f.QueryPort}
	if ip := net.ParseIP(f.PublicIP); ip != nil && ip.To4() != nil {
		r.PublicIP = ip.String()
	}
	add := func(key, status string, params map[string]string) {
		r.Items = append(r.Items, Item{Key: key, Status: status, Params: params})
	}
	port := func(p int) map[string]string { return map[string]string{"port": strconv.Itoa(p)} }

	// Servidor ligado e respondendo nesta máquina.
	if !f.ServerRunning {
		add("serverOff", Info, nil)
	} else if f.LocalAnswer {
		add("localOk", OK, port(f.QueryPort))
	} else {
		add("localNo", Bad, port(f.QueryPort))
	}

	// Firewall do Windows (só as regras que o app cria; uma regra feita à mão não aparece aqui).
	switch f.Firewall {
	case Yes:
		add("firewallOk", OK, nil)
	case No:
		add("firewallMissing", Warn, nil)
	}

	// Roteador.
	mappedOK := false
	if !f.Router.Found {
		add("routerUnknown", Info, nil)
	} else {
		game := mappingItem(f.Router.Game, f.GamePort, f.LocalIPs)
		query := mappingItem(f.Router.Query, f.QueryPort, f.LocalIPs)
		for _, it := range []Item{game, query} {
			if it.Key != "" {
				r.Items = append(r.Items, it)
			}
		}
		mappedOK = game.Status == OK && query.Status == OK
		if f.Rcon && f.Router.Rcon.Exists {
			add("rconForwarded", Bad, port(f.RconPort))
		}
	}

	// CGNAT: o roteador não tem o IP que a internet vê.
	cgnat := false
	if ext := net.ParseIP(f.Router.ExternalIP); f.Router.Found && ext != nil {
		switch {
		case notPublic(ext):
			cgnat = true
			add("cgnat", Bad, map[string]string{"ip": ext.String()})
		case r.PublicIP != "" && ext.String() != r.PublicIP:
			cgnat = true
			add("ipMismatch", Warn, map[string]string{"router": ext.String(), "public": r.PublicIP})
		case r.PublicIP != "":
			add("publicOk", OK, map[string]string{"ip": r.PublicIP})
		}
	} else if r.PublicIP != "" {
		add("publicIp", Info, map[string]string{"ip": r.PublicIP})
	} else {
		add("publicUnknown", Info, nil)
	}

	// Testes que precisam do servidor ligado e do IP público.
	if f.ServerRunning && r.PublicIP != "" {
		if f.Loopback {
			add("loopbackOk", OK, map[string]string{"ip": r.PublicIP, "port": strconv.Itoa(f.QueryPort)})
		} else {
			add("loopbackNo", Info, nil)
		}
		switch f.Steam {
		case Yes:
			add("steamListed", OK, nil)
		case No:
			add("steamNotListed", Warn, nil)
		}
	}

	r.Verdict = verdict(r.Items, f, mappedOK, cgnat)
	return r
}

func verdict(items []Item, f Facts, mappedOK, cgnat bool) string {
	for _, it := range items {
		if it.Status == Bad {
			return VerdictProblem
		}
	}
	if f.ServerRunning && f.Loopback {
		return VerdictOpen
	}
	if f.ServerRunning && f.LocalAnswer && mappedOK && !cgnat && f.Firewall != No {
		return VerdictLikely
	}
	return VerdictUnsure
}

// mappingItem diz se a porta (UDP) está redirecionada para este PC.
func mappingItem(m Mapping, port int, local []string) Item {
	params := map[string]string{"port": strconv.Itoa(port)}
	switch {
	case !m.Known:
		return Item{}
	case !m.Exists:
		// Muitos roteadores não mostram por UPnP os redirecionamentos feitos à mão: é aviso, não erro.
		return Item{Key: "routerNoMap", Status: Warn, Params: params}
	case !contains(local, m.Client):
		params["client"] = m.Client
		return Item{Key: "routerOtherPc", Status: Bad, Params: params}
	case m.InternalPort != port:
		params["internal"] = strconv.Itoa(m.InternalPort)
		return Item{Key: "routerWrongPort", Status: Bad, Params: params}
	case !m.Enabled:
		return Item{Key: "routerDisabled", Status: Bad, Params: params}
	}
	return Item{Key: "routerMapped", Status: OK, Params: params}
}

// cgnatRange é a faixa 100.64.0.0/10, usada pelas operadoras no CGNAT.
var cgnatRange = &net.IPNet{IP: net.IPv4(100, 64, 0, 0), Mask: net.CIDRMask(10, 32)}

// notPublic diz se o IP não é da internet (rede local, CGNAT, loopback...).
func notPublic(ip net.IP) bool {
	return ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsUnspecified() || cgnatRange.Contains(ip)
}

// NotPublic é notPublic para os adaptadores (conferir o IP que o STUN devolveu).
func NotPublic(ip net.IP) bool { return notPublic(ip) }

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
