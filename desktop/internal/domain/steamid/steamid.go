// Package steamid valida e extrai o SteamID64 (usado na lista de administradores).
package steamid

import (
	"regexp"
	"strconv"
)

// base é o menor SteamID64 de conta individual (76561197960265728).
const base uint64 = 76561197960265728

var (
	digitsRe  = regexp.MustCompile(`^7656119\d{10}$`)
	profileRe = regexp.MustCompile(`steamcommunity\.com/profiles/(7656119\d{10})`)
	anyRe     = regexp.MustCompile(`7656119\d{10}`)
)

// Valid diz se o texto é um SteamID64 de conta individual (17 números, começando com 7656119).
func Valid(id string) bool {
	if !digitsRe.MatchString(id) {
		return false
	}
	n, err := strconv.ParseUint(id, 10, 64)
	return err == nil && n > base
}

// Extract acha o SteamID64 num texto colado: o número sozinho ou o link do perfil
// (https://steamcommunity.com/profiles/7656119...). Link com nome personalizado (/id/nome) não tem o número.
func Extract(text string) (string, bool) {
	if m := profileRe.FindStringSubmatch(text); m != nil && Valid(m[1]) {
		return m[1], true
	}
	if m := anyRe.FindString(text); m != "" && Valid(m) {
		return m, true
	}
	return "", false
}
