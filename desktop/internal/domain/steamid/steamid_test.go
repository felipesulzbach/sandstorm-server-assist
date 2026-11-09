package steamid

import "testing"

func TestValidAndExtract(t *testing.T) {
	if !Valid("76561198000000001") {
		t.Error("válido recusado")
	}
	for _, bad := range []string{"", "123", "76561197960265728", "7656119800000000a", "86561198000000001"} {
		if Valid(bad) {
			t.Errorf("inválido aceito: %q", bad)
		}
	}
	if id, ok := Extract("https://steamcommunity.com/profiles/76561198000000001/"); !ok || id != "76561198000000001" {
		t.Errorf("Extract link = %q %v", id, ok)
	}
	if _, ok := Extract("https://steamcommunity.com/id/meunome/"); ok {
		t.Error("link com nome personalizado não tem número")
	}
}
