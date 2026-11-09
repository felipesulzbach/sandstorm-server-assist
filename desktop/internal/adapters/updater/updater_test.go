package updater

import (
	"context"
	"testing"
)

func TestFetcherRefusesOtherHostsAndNames(t *testing.T) {
	for _, base := range []string{
		"http://game-in.duckdns.org/sandstorm/downloads/", // sem HTTPS
		"https://example.com/downloads/",                  // host fora da lista
		"https://game-in.duckdns.org/sandstorm/downloads", // sem a barra do fim
	} {
		if _, err := NewHTTPFetcher(base); err == nil {
			t.Errorf("accepted %s", base)
		}
	}
	f, err := NewHTTPFetcher(DefaultBaseURL)
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"../x.exe", "a/b.exe", `a\b.exe`, "", ".hidden"} {
		if _, _, err := f.Download(context.Background(), name, t.TempDir()+"/x", nil); err == nil {
			t.Errorf("accepted file name %q", name)
		}
	}
	if u, _ := f.url("latest.json"); u != DefaultBaseURL+"latest.json" {
		t.Fatalf("url %s", u)
	}
}

func TestPublicKeyIsValid(t *testing.T) {
	if len(PublicKey()) != 32 {
		t.Fatal("public key size")
	}
}
