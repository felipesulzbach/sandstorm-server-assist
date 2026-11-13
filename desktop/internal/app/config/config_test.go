package config

import (
	"errors"
	"path/filepath"
	"testing"
	"time"

	"sandstorm-server-assist/desktop/internal/domain/legal"
	"sandstorm-server-assist/desktop/internal/domain/problems"
	"sandstorm-server-assist/desktop/internal/domain/server"
)

type memVault struct{ s server.Secrets }

func (m *memVault) Load() (server.Secrets, error) { return m.s, nil }
func (m *memVault) Save(s server.Secrets) error   { m.s = s; return nil }

func TestProfileSettingsAndSecrets(t *testing.T) {
	dir := t.TempDir()
	v := &memVault{}
	s := Open(filepath.Join(dir, "config.json"), v)
	if _, err := s.SetServerDir(filepath.Join(dir, "srv")); err != nil {
		t.Fatal(err)
	}

	// A tela não troca as pastas pelo SaveProfile.
	p := s.Profile()
	p.Hostname = "Novo nome"
	p.ServerDir = `C:\Windows`
	saved, err := s.SaveProfile(p)
	if err != nil || saved.Hostname != "Novo nome" || saved.ServerDir != filepath.Join(dir, "srv") {
		t.Fatalf("%+v %v", saved, err)
	}
	p.Port = 1
	if _, err := s.SaveProfile(p); err == nil {
		t.Error("porta inválida aceita")
	}

	if _, err := s.SetLanguage("de-DE"); err == nil {
		t.Error("idioma fora da lista aceito")
	}
	if st, err := s.SetLanguage("uk"); err != nil || st.Language != "uk" {
		t.Errorf("%+v %v", st, err)
	}

	if err := s.SetSecret("rconPassword", "a b"); err == nil {
		t.Error("senha com espaço aceita")
	}
	if err := s.SetSecret("rconPassword", "boa#Senha1"); err != nil {
		t.Fatal(err)
	}
	state, _ := s.SecretsState()
	if !state["rconPassword"] || state["gsltToken"] || v.s.RconPassword != "boa#Senha1" {
		t.Errorf("%v %+v", state, v.s)
	}

	// Reabre: perfil e preferências voltam do arquivo; segredos do cofre.
	again := Open(filepath.Join(dir, "config.json"), v)
	if again.Profile().Hostname != "Novo nome" || again.Settings().Language != "uk" {
		t.Errorf("%+v %+v", again.Profile(), again.Settings())
	}
	if pw, _ := again.Secret("rconPassword"); pw != "boa#Senha1" {
		t.Error("segredo perdido")
	}
	if _, err := again.SetServerDir(`C:\Windows\System32`); err == nil {
		t.Error("pasta do sistema aceita")
	}
}

func TestAcceptTerms(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	s := Open(path, &memVault{})
	if s.Settings().TermsAccepted != "" {
		t.Fatal("new config should not have terms accepted")
	}
	now := time.Date(2026, 11, 9, 12, 0, 0, 0, time.UTC)
	if _, err := s.AcceptTerms("1999-01-01", now); err == nil {
		t.Fatal("old terms version accepted")
	}
	st, err := s.AcceptTerms(legal.Version, now)
	if err != nil || st.TermsAccepted != legal.Version || st.TermsAcceptedAt != "2026-11-09T12:00:00Z" {
		t.Fatalf("%+v %v", st, err)
	}
	// Fica gravado.
	if got := Open(path, &memVault{}).Settings().TermsAccepted; got != legal.Version {
		t.Fatalf("not saved: %q", got)
	}
}

// A pasta do SteamCMD não pode ficar dentro da do servidor, nem o contrário (ADR-042).
func TestRefusesNestedSteamCMDDir(t *testing.T) {
	dir := t.TempDir()
	s := Open(filepath.Join(dir, "config.json"), &memVault{})
	srv := filepath.Join(dir, "servers")
	if _, err := s.SetServerDir(srv); err != nil {
		t.Fatal(err)
	}
	_, err := s.SetSteamCMDDir(filepath.Join(srv, "steamcmd"))
	var pe *problems.Error
	if !errors.As(err, &pe) || pe.Code != problems.FoldersNested {
		t.Fatalf("nested steamcmd accepted: %v", err)
	}
	if _, err := s.SetSteamCMDDir(filepath.Join(dir, "steamcmd")); err != nil {
		t.Fatalf("sibling refused: %v", err)
	}
	if _, err := s.SetServerDir(filepath.Join(dir, "steamcmd", "server")); !errors.As(err, &pe) {
		t.Fatalf("server inside steamcmd accepted: %v", err)
	}
	if s.Profile().ServerDir != srv {
		t.Fatalf("server dir changed to %q", s.Profile().ServerDir)
	}
}
