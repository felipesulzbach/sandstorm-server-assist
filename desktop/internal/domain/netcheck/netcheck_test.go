package netcheck

import (
	"testing"

	"sandstorm-server-assist/desktop/internal/domain/problems"
)

func TestParseOK(t *testing.T) {
	req := Request{QueryPort: 27131, RconPort: 27015}
	body := `{"ip":"203.0.113.7","query":{"port":27131,"open":true,"info":{"name":"My Server","map":"Farmhouse","players":2,"maxPlayers":8}},"rcon":{"port":27015,"tested":true,"exposed":true}}`
	res, err := Parse(200, []byte(body), req, 27102)
	if err != nil {
		t.Fatal(err)
	}
	want := Result{PublicIP: "203.0.113.7", GamePort: 27102, QueryPort: 27131, QueryOpen: true, ServerName: "My Server",
		Map: "Farmhouse", Players: 2, MaxPlayers: 8, RconPort: 27015, RconTested: true, RconExposed: true}
	if res != want {
		t.Fatalf("got %+v", res)
	}
}

func TestParseClosed(t *testing.T) {
	body := `{"ip":"203.0.113.7","query":{"port":27131,"open":false},"rcon":{"port":0,"tested":false,"exposed":false}}`
	res, err := Parse(200, []byte(body), Request{QueryPort: 27131}, 27102)
	if err != nil || res.QueryOpen || res.RconTested || res.ServerName != "" {
		t.Fatalf("got %+v %v", res, err)
	}
}

func TestParseErrors(t *testing.T) {
	req := Request{QueryPort: 27131}
	for _, c := range []struct {
		status int
		body   string
		code   string
	}{
		{429, `{"error":"rate-limited"}`, problems.PortCheckLimit},
		{429, `{"error":"busy"}`, problems.PortCheckLimit},
		{400, `{"error":"ipv6"}`, problems.PortCheckIPv6},
		{400, `{"error":"private-ip"}`, problems.PortCheckFailed},
		{502, `<html>Bad gateway</html>`, problems.PortCheckFailed},
		// Resposta de outro pedido ou com IP inválido não vale.
		{200, `{"ip":"203.0.113.7","query":{"port":9999,"open":true}}`, problems.PortCheckFailed},
		{200, `{"ip":"not-an-ip","query":{"port":27131,"open":true}}`, problems.PortCheckFailed},
	} {
		_, err := Parse(c.status, []byte(c.body), req, 27102)
		if problems.CodeOf(err) != c.code {
			t.Fatalf("%d %s: %v", c.status, c.body, err)
		}
	}
}
