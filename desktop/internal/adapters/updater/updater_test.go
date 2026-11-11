package updater

import (
	"context"
	"testing"
)

func TestFetcherRefusesOtherHostsAndNames(t *testing.T) {
	for _, base := range []string{
		"http://github.com/felipesulzbach/sandstorm-server-assist/releases/", // sem HTTPS
		"https://example.com/downloads/",                                     // host fora da lista
		"https://github.com/felipesulzbach/sandstorm-server-assist/releases", // sem a barra do fim
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
	for _, name := range []string{"other.exe", "sandstorm-server-assist-0.2-windows-setup.exe"} {
		if _, err := f.url(name); err == nil {
			t.Errorf("accepted file name %q", name)
		}
	}
	if u, _ := f.url("latest.json"); u != DefaultBaseURL+"latest/download/latest.json" {
		t.Fatalf("manifest url %s", u)
	}
	want := DefaultBaseURL + "download/v0.2.0/sandstorm-server-assist-0.2.0-windows-setup.exe"
	if u, _ := f.url("sandstorm-server-assist-0.2.0-windows-setup.exe"); u != want {
		t.Fatalf("installer url %s", u)
	}
}

func TestPublicKeyIsValid(t *testing.T) {
	if len(PublicKey()) != 32 {
		t.Fatal("public key size")
	}
}
