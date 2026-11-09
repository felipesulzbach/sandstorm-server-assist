package backup

import (
	"bytes"
	"crypto/rand"
	"errors"
	"testing"

	"sandstorm-server-assist/desktop/internal/domain/problems"
	"sandstorm-server-assist/desktop/internal/domain/server"
)

func TestRoundTrip(t *testing.T) {
	c := Contents{
		Manifest: Manifest{AppVersion: "0.1.0", Machine: "PC"},
		Profile:  server.Default(),
		Secrets:  server.Secrets{RconPassword: "rcon#Senha1"},
		Files: map[string][]byte{
			"server/Game.ini": []byte("[Rcon]\r\nbEnabled=True\r\n"),
			ModioSessionName:  []byte(`{"OAuth":{}}`),
		},
	}
	plain, err := Pack(c)
	if err != nil {
		t.Fatal(err)
	}
	sealed, err := Seal("senha-forte-123", plain, rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(sealed, []byte("rcon#Senha1")) {
		t.Fatal("segredo em texto puro no arquivo")
	}
	opened, err := Open("senha-forte-123", sealed)
	if err != nil {
		t.Fatal(err)
	}
	got, err := Unpack(opened)
	if err != nil {
		t.Fatal(err)
	}
	if got.Secrets.RconPassword != "rcon#Senha1" || got.Profile.Port != 27102 || !got.Manifest.HasModio {
		t.Errorf("%+v", got)
	}
	if string(got.Files["server/Game.ini"]) != "[Rcon]\r\nbEnabled=True\r\n" {
		t.Error("Game.ini mudou")
	}

	_, err = Open("senha-errada-123", sealed)
	var pe *problems.Error
	if !errors.As(err, &pe) || pe.Code != problems.BackupPassword {
		t.Errorf("senha errada: %v", err)
	}
	sealed[len(sealed)-1] ^= 1
	if _, err := Open("senha-forte-123", sealed); err == nil {
		t.Error("arquivo alterado deveria falhar")
	}
	if _, err := Seal("curta", plain, rand.Reader); err == nil {
		t.Error("senha curta aceita")
	}
}

func TestPackRejectsUnknownNames(t *testing.T) {
	if _, err := Pack(Contents{Files: map[string][]byte{"../../Windows/x.dll": nil}}); err == nil {
		t.Error("nome fora da lista aceito")
	}
}
