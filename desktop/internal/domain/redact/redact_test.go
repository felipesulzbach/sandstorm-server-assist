package redact

import (
	"strings"
	"testing"
)

func TestText(t *testing.T) {
	in := `GET /v1/me/subscribed?_limit=100&api_key=0123456789abcdef0123456789abcdef -RconPassword=abc123 -GSLTToken=ABCDEF0123456789 -SecurityCode=none -SecurityCode=PBRXD Password=segredo1 Authorization: Bearer xyz.abc`
	out := Text(in)
	for _, leak := range []string{"0123456789abcdef", "abc123", "ABCDEF0123", "PBRXD", "segredo1", "xyz.abc"} {
		if strings.Contains(out, leak) {
			t.Errorf("vazou %q em %q", leak, out)
		}
	}
	if !strings.Contains(out, "-SecurityCode=none") {
		t.Errorf("none deveria continuar: %q", out)
	}
}

func TestValues(t *testing.T) {
	if got := Values("senha minhaSenha e abc", "minhaSenha", "abc"); got != "senha ••• e abc" {
		t.Errorf("%q", got)
	}
}
