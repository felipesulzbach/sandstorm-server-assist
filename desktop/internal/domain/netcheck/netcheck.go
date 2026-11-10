// Package netcheck são as regras do teste de portas pela internet (ADR-021): o que o app manda ao verificador e
// como lê a resposta. O verificador (pasta checker do repositório) testa, de fora, o IP de quem pediu.
package netcheck

import (
	"encoding/json"
	"net"
	"unicode/utf8"

	"sandstorm-server-assist/desktop/internal/domain/problems"
)

// Request é o que vai ao verificador: só as portas.
type Request struct {
	QueryPort int `json:"queryPort"`
	RconPort  int `json:"rconPort,omitempty"`
}

// reply é a resposta do verificador, como chega.
type reply struct {
	IP    string `json:"ip"`
	Error string `json:"error"`
	Query struct {
		Port int  `json:"port"`
		Open bool `json:"open"`
		Info *struct {
			Name       string `json:"name"`
			Map        string `json:"map"`
			Players    int    `json:"players"`
			MaxPlayers int    `json:"maxPlayers"`
		} `json:"info"`
	} `json:"query"`
	Rcon struct {
		Port    int  `json:"port"`
		Tested  bool `json:"tested"`
		Exposed bool `json:"exposed"`
	} `json:"rcon"`
}

// Result é o que a tela mostra.
type Result struct {
	// PublicIP é o IP que a internet vê (o que os jogadores usam sem VPN).
	PublicIP   string `json:"publicIp"`
	GamePort   int    `json:"gamePort"`
	QueryPort  int    `json:"queryPort"`
	QueryOpen  bool   `json:"queryOpen"`
	ServerName string `json:"serverName"`
	Map        string `json:"map"`
	Players    int    `json:"players"`
	MaxPlayers int    `json:"maxPlayers"`
	RconPort   int    `json:"rconPort"`
	RconTested bool   `json:"rconTested"`
	// RconExposed: a porta do RCON aceita conexões da internet (risco: deve ser fechada no roteador).
	RconExposed bool `json:"rconExposed"`
}

// MaxReply é o maior corpo de resposta aceito.
const MaxReply = 16 << 10

// errorCodes liga os erros do verificador ao catálogo.
var errorCodes = map[string]string{
	"rate-limited": problems.PortCheckLimit,
	"busy":         problems.PortCheckLimit,
	"ipv6":         problems.PortCheckIPv6,
}

// Parse lê a resposta do verificador (status HTTP e corpo) e confere com o que foi pedido.
func Parse(status int, body []byte, req Request, gamePort int) (Result, error) {
	var r reply
	if err := json.Unmarshal(body, &r); err != nil {
		return Result{}, problems.New(problems.PortCheckFailed, "invalid reply")
	}
	if status != 200 || r.Error != "" {
		code, ok := errorCodes[r.Error]
		if !ok {
			code = problems.PortCheckFailed
		}
		return Result{}, problems.New(code, r.Error)
	}
	ip := net.ParseIP(r.IP)
	if ip == nil || r.Query.Port != req.QueryPort || (r.Rcon.Tested && r.Rcon.Port != req.RconPort) {
		return Result{}, problems.New(problems.PortCheckFailed, "reply does not match the request")
	}
	res := Result{
		PublicIP: ip.String(), GamePort: gamePort, QueryPort: req.QueryPort, QueryOpen: r.Query.Open,
		RconPort: req.RconPort, RconTested: r.Rcon.Tested, RconExposed: r.Rcon.Tested && r.Rcon.Exposed,
	}
	if i := r.Query.Info; r.Query.Open && i != nil {
		res.ServerName, res.Map = clip(i.Name), clip(i.Map)
		res.Players, res.MaxPlayers = bound(i.Players), bound(i.MaxPlayers)
	}
	return res, nil
}

func clip(s string) string {
	if !utf8.ValidString(s) {
		return ""
	}
	if r := []rune(s); len(r) > 64 {
		return string(r[:64])
	}
	return s
}

func bound(n int) int {
	if n < 0 || n > 255 {
		return 0
	}
	return n
}
