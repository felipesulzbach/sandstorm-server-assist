package config

import (
	"os"
	"path/filepath"
	"testing"

	"sandstorm-server-assist/desktop/internal/domain/server"
)

// 14/11/2026: a pasta do servidor foi movida para começar do zero, e o app continuava com tudo do servidor antigo.
func TestFreshStart(t *testing.T) {
	dir := t.TempDir()
	v := &memVault{}
	s := Open(filepath.Join(dir, "config.json"), v)

	// App novo: nada a oferecer.
	if f := s.FreshState(); f.Offer || f.Missing || f.HasSecrets {
		t.Fatalf("app novo: %+v", f)
	}

	// Servidor instalado e configurado: não oferece.
	srv := filepath.Join(dir, "sandstorm_server")
	writeFile(t, filepath.Join(srv, serverExe), "x")
	if _, err := s.SetServerDir(srv); err != nil {
		t.Fatal(err)
	}
	p := s.Profile()
	p.Hostname, p.Admins = "Old server", []server.Admin{{SteamID: "76561198000000001", Name: "Admin"}}
	if _, err := s.SaveProfile(p); err != nil {
		t.Fatal(err)
	}
	if err := s.SetSecret("rconPassword", "rcon#Senha1"); err != nil {
		t.Fatal(err)
	}
	if f := s.FreshState(); f.Offer || f.Missing || !f.HasSecrets {
		t.Fatalf("servidor instalado: %+v", f)
	}

	// A pasta foi movida: oferece, e diz que ela sumiu.
	if err := os.RemoveAll(srv); err != nil {
		t.Fatal(err)
	}
	if f := s.FreshState(); !f.Offer || !f.Missing {
		t.Fatalf("pasta movida: %+v", f)
	}

	// Começar do zero sem apagar o cofre: perfil padrão, pastas mantidas, senha continua.
	steam := s.Profile().SteamCMDDir
	got, err := s.StartFresh(false)
	if err != nil {
		t.Fatal(err)
	}
	if got.ServerDir != srv || got.SteamCMDDir != steam || got.Hostname == "Old server" || len(got.Admins) != 0 {
		t.Fatalf("perfil: %+v", got)
	}
	if v.s.RconPassword == "" {
		t.Fatal("cofre apagado sem pedir")
	}
	// Ainda oferece por causa do cofre; apagando o cofre também, não oferece mais.
	if f := s.FreshState(); !f.Offer || !f.HasSecrets {
		t.Fatalf("só cofre: %+v", f)
	}
	if _, err := s.StartFresh(true); err != nil {
		t.Fatal(err)
	}
	if v.s != (server.Secrets{}) {
		t.Fatalf("cofre: %+v", v.s)
	}
	if f := s.FreshState(); f.Offer || f.HasSecrets {
		t.Fatalf("depois de começar do zero: %+v", f)
	}
	// Gravado: abrir de novo mantém o perfil limpo.
	if again := Open(filepath.Join(dir, "config.json"), v).Profile(); again.Hostname == "Old server" {
		t.Fatal("perfil antigo voltou")
	}
}
