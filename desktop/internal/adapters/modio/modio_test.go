package modio

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestReadInstalledAndSession(t *testing.T) {
	root := t.TempDir()
	meta := filepath.Join(root, "254", "metadata")
	if err := os.MkdirAll(meta, 0o700); err != nil {
		t.Fatal(err)
	}
	state := `{"Mods":[{"ID":2123731,"PathOnDisk":"C:\\x","SizeOnDisk":100,"State":1,"Profile":{"id":2123731,"name":"Containment","profile_url":"https://mod.io/x","modfile":{"filesize":50,"filesize_uncompressed":5000}}}],"version":1}`
	if err := os.WriteFile(filepath.Join(meta, "state.json"), []byte(state), 0o600); err != nil {
		t.Fatal(err)
	}
	list, err := ReadInstalled(root)
	if err != nil || len(list) != 1 || list[0].Name != "Containment" || list[0].Size != 5000 {
		t.Fatalf("%+v %v", list, err)
	}

	uf := filepath.Join(root, "user.json")
	exp := time.Now().Add(24 * time.Hour).Unix()
	body := `{"OAuth":{"expiry":` + itoa(exp) + `,"status":0,"token":"tok"},"Profile":{"username":"servidor"}}`
	if err := os.WriteFile(uf, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	s, tok, err := ReadSession(uf)
	if err != nil || !s.Valid || tok != "tok" || s.Username != "servidor" {
		t.Fatalf("%+v %q %v", s, tok, err)
	}
	if s, _, _ := ReadSession(filepath.Join(root, "nao-existe.json")); s.Found {
		t.Error("sessão ausente marcada como encontrada")
	}
}

func itoa(n int64) string {
	b := []byte{}
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	return string(b)
}

func TestFindAPIKey(t *testing.T) {
	line := `LogModio: GET request for g-254.modapi.io /v1/me/subscribed?_limit=100&game_id=254&api_key=0123456789abcdef0123456789abcdef`
	if FindAPIKey(line) != "0123456789abcdef0123456789abcdef" {
		t.Error("chave não encontrada")
	}
	if FindAPIKey("nada aqui") != "" {
		t.Error("chave falsa")
	}
}

// O cliente recusa qualquer endereço fora da lista (mesmo que o código peça outro BaseURL).
func TestClientBlocksOtherHosts(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(200) }))
	defer srv.Close()
	c := Client{BaseURL: srv.URL}
	if err := c.RequestCode(context.Background(), "k", "a@b.com"); err == nil {
		t.Error("host fora da lista aceito")
	}
}
