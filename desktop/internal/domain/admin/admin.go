// Package admin tem as regras do painel de jogadores (ADR-022): ler a lista de jogadores e de banidos que o RCON
// do servidor devolve e montar os comandos de expulsar, banir e tirar o banimento. Não faz I/O.
//
// Comandos (help do servidor, conferido em 10/11/2026):
//
//	listplayers
//	kick <id/netid/name> [reason]
//	banid <netid> [duration in minutes] [reason]
//	permban <id/netid/name> [reason]
//	unban <netid>
//	listbans
package admin

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"sandstorm-server-assist/desktop/internal/domain/problems"
)

// Player é uma linha do listplayers.
type Player struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	NetID string `json:"netId"`
	IP    string `json:"ip"`
	Score int    `json:"score"`
	// Bot: sem net ID (bots do coop). Não dá para expulsar nem banir.
	Bot bool `json:"bot"`
}

// Ban é um item do listbans: o net ID e a linha como o servidor mostrou (data, motivo).
type Ban struct {
	NetID string `json:"netId"`
	Text  string `json:"text"`
}

// Limites.
const (
	MaxReason = 100
	// MaxBanMinutes é um ano. Banimento permanente usa permban.
	MaxBanMinutes = 525600
	maxRows       = 256
)

// netIDRe: "SteamNWI:76561198000000000" (como aparece no log e no listplayers) ou só o SteamID64.
var (
	netIDRe   = regexp.MustCompile(`^(?:[A-Za-z]{2,16}:)?[0-9]{5,25}$`)
	netIDFind = regexp.MustCompile(`(?:[A-Za-z]{2,16}:)?[0-9]{15,25}`)
)

// ValidNetID diz se o texto é um net ID aceito nos comandos.
func ValidNetID(id string) bool { return netIDRe.MatchString(id) }

// ParsePlayers lê a tabela do listplayers:
//
//	ID	 | Name				 | NetID			 | IP			 | Score		 |
//	===============================================================================
//	0	 | Fulano	 | SteamNWI:765...	 | 203.0.113.5	 | 120	 |
func ParsePlayers(out string) []Player {
	var list []Player
	for _, line := range strings.Split(strings.ReplaceAll(out, "\r\n", "\n"), "\n") {
		cols := strings.Split(line, "|")
		if len(cols) < 4 {
			continue
		}
		for i := range cols {
			cols[i] = strings.TrimSpace(cols[i])
		}
		if strings.EqualFold(cols[0], "ID") || strings.Trim(line, "=- \t") == "" {
			continue
		}
		if _, err := strconv.Atoi(cols[0]); err != nil {
			continue
		}
		p := Player{ID: cols[0], Name: cleanText(cols[1], 64), NetID: cols[2]}
		if len(cols) > 3 {
			p.IP = cols[3]
		}
		if len(cols) > 4 {
			p.Score, _ = strconv.Atoi(cols[4])
		}
		if !ValidNetID(p.NetID) {
			p.NetID, p.Bot = "", true
		}
		list = append(list, p)
		if len(list) >= maxRows {
			break
		}
	}
	return list
}

// ParseBans lê o listbans: cada linha com um net ID vira um banimento (o resto da linha fica como texto).
func ParseBans(out string) []Ban {
	var list []Ban
	seen := map[string]bool{}
	for _, line := range strings.Split(strings.ReplaceAll(out, "\r\n", "\n"), "\n") {
		id := netIDFind.FindString(line)
		if id == "" || seen[id] {
			continue
		}
		seen[id] = true
		list = append(list, Ban{NetID: id, Text: cleanText(strings.TrimSpace(line), 200)})
		if len(list) >= maxRows {
			break
		}
	}
	return list
}

// cleanText tira caracteres de controle e corta o tamanho.
func cleanText(s string, max int) string {
	s = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return ' '
		}
		return r
	}, s)
	s = strings.Join(strings.Fields(s), " ")
	if utf8.RuneCountInString(s) > max {
		s = string([]rune(s)[:max])
	}
	return s
}

// reason confere o motivo digitado (opcional).
func reason(r string) (string, error) {
	r = strings.TrimSpace(r)
	if utf8.RuneCountInString(r) > MaxReason {
		return "", problems.Field("reason", "validation.reason", "max", strconv.Itoa(MaxReason))
	}
	for _, c := range r {
		if unicode.IsControl(c) || c == '"' {
			return "", problems.Field("reason", "validation.reason", "max", strconv.Itoa(MaxReason))
		}
	}
	return strings.Join(strings.Fields(r), " "), nil
}

func withReason(cmd, r string) string {
	if r == "" {
		return cmd
	}
	return cmd + " " + r
}

// KickCommand monta o "kick" de um jogador da lista.
func KickCommand(netID, why string) (string, error) {
	if !ValidNetID(netID) {
		return "", problems.Field("netId", "validation.netId")
	}
	r, err := reason(why)
	if err != nil {
		return "", err
	}
	return withReason("kick "+netID, r), nil
}

// BanCommand monta o banimento: minutes > 0 é temporário (banid), 0 é permanente (permban).
func BanCommand(netID string, minutes int, why string) (string, error) {
	if !ValidNetID(netID) {
		return "", problems.Field("netId", "validation.netId")
	}
	if minutes < 0 || minutes > MaxBanMinutes {
		return "", problems.Field("minutes", "validation.banMinutes", "max", strconv.Itoa(MaxBanMinutes))
	}
	r, err := reason(why)
	if err != nil {
		return "", err
	}
	if minutes == 0 {
		return withReason("permban "+netID, r), nil
	}
	return withReason("banid "+netID+" "+strconv.Itoa(minutes), r), nil
}

// UnbanCommand monta o "unban".
func UnbanCommand(netID string) (string, error) {
	if !ValidNetID(netID) {
		return "", problems.Field("netId", "validation.netId")
	}
	return "unban " + netID, nil
}
