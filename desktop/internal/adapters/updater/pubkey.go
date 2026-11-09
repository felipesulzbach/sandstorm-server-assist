package updater

import (
	"crypto/ed25519"
	"encoding/base64"
)

// publicKeyB64 é a chave pública da publicação (ADR-014). A privada fica só na máquina de quem publica
// (~/.sandstorm-server-assist/update-signing-key.pem, "npm run update:keygen"). Trocar esta chave exige instalar
// a versão nova à mão.
const publicKeyB64 = "/w/LOXQ5WKZCCc6UKUW5ovA+MeQ+IokjiTaBZ8HmjTE="

// PublicKey devolve a chave pública embutida (nil se estiver inválida: a atualização fica desligada).
func PublicKey() ed25519.PublicKey {
	b, err := base64.StdEncoding.DecodeString(publicKeyB64)
	if err != nil || len(b) != ed25519.PublicKeySize {
		return nil
	}
	return ed25519.PublicKey(b)
}
