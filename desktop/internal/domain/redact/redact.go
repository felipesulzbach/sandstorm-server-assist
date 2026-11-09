// Package redact mascara segredos em textos que vão para a tela ou para log (ADR-002).
package redact

import (
	"regexp"
	"strings"
)

const mask = "•••"

var patterns = []*regexp.Regexp{
	// chave=valor em URL, linha de comando ou INI: api_key, Password, RconPassword, GSLTToken, GameStatsToken,
	// SecurityCode, AccessToken, token. O valor vai até espaço, aspas, ?, & ou fim.
	regexp.MustCompile(`(?i)((?:api_key|password|rconpassword|gslttoken|gamestatstoken|securitycode|accesstoken|access_token|token)\s*[=:]\s*"?)([^\s"?&,;]+)`),
	// Cabeçalho Authorization: Bearer xxx
	regexp.MustCompile(`(?i)(bearer\s+)([A-Za-z0-9._~+/=-]+)`),
	// JWT solto (três partes base64url com pontos).
	regexp.MustCompile(`()(eyJ[A-Za-z0-9_-]{8,}\.[A-Za-z0-9_-]{8,}\.[A-Za-z0-9_-]{8,})`),
}

// Text troca valores secretos conhecidos por "•••". "none" no SecurityCode é mantido (não é segredo).
func Text(s string) string {
	for _, re := range patterns {
		s = re.ReplaceAllStringFunc(s, func(m string) string {
			sub := re.FindStringSubmatch(m)
			if len(sub) < 3 || strings.EqualFold(sub[2], "none") {
				return m
			}
			return sub[1] + mask
		})
	}
	return s
}

// Values troca valores exatos (segredos do cofre) por "•••". Valores curtos (menos de 4 caracteres) não
// são trocados, para não mascarar palavras comuns.
func Values(s string, secrets ...string) string {
	for _, v := range secrets {
		if len(v) >= 4 {
			s = strings.ReplaceAll(s, v, mask)
		}
	}
	return s
}
