package server

import (
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"sandstorm-server-assist/desktop/internal/domain/problems"
)

// Automação (ADR-022): mensagens programadas no chat, reinício automático e avisos no Discord.

// Limites.
const (
	MaxAnnouncements = 10
	MinAnnounceEvery = 2
	MaxAnnounceEvery = 240
	// MaxDiscord é o tamanho de cada texto de aviso do Discord (o Discord aceita até 2000 por mensagem).
	MaxDiscord = 300
	// Marcadores dos avisos do Discord.
	ServerTag  = "{server}"
	MapTag     = "{map}"
	MinutesTag = "{minutes}"
)

// validateAutomation confere as mensagens programadas, o reinício automático e os textos do Discord.
func (p Profile) validateAutomation() error {
	if len(p.Announcements) > MaxAnnouncements {
		return problems.Field("announcements", "validation.announcementsMax", "max", strconv.Itoa(MaxAnnouncements))
	}
	for _, a := range p.Announcements {
		if strings.TrimSpace(a) == "" || utf8.RuneCountInString(a) > MaxChat || hasControl(a, false) {
			return problems.Field("announcements", "validation.chatMessage", "max", strconv.Itoa(MaxChat))
		}
	}
	if len(p.Announcements) > 0 {
		if p.AnnounceEvery < MinAnnounceEvery || p.AnnounceEvery > MaxAnnounceEvery {
			return problems.Field("announceEvery", "validation.announceEvery",
				"min", strconv.Itoa(MinAnnounceEvery), "max", strconv.Itoa(MaxAnnounceEvery))
		}
		if !p.Rcon {
			return problems.Field("announcements", "validation.chatNeedsRcon")
		}
	}
	if p.AutoRestart.Daily != "" {
		if _, _, ok := ParseClock(p.AutoRestart.Daily); !ok {
			return problems.Field("autoRestartDaily", "validation.clock")
		}
	}
	if w := p.AutoRestart.WarnMessage; w != "" {
		if utf8.RuneCountInString(w) > MaxChat || hasControl(w, false) {
			return problems.Field("autoRestartWarn", "validation.chatMessage", "max", strconv.Itoa(MaxChat))
		}
		if !p.Rcon {
			return problems.Field("autoRestartWarn", "validation.chatNeedsRcon")
		}
	}
	for _, f := range []struct{ name, v string }{
		{"discordReady", p.Discord.Ready}, {"discordStopped", p.Discord.Stopped}, {"discordCrashed", p.Discord.Crashed},
		{"discordJoin", p.Discord.Join}, {"discordLeave", p.Discord.Leave},
	} {
		if utf8.RuneCountInString(f.v) > MaxDiscord || hasControl(f.v, false) {
			return problems.Field(f.name, "validation.discordText", "max", strconv.Itoa(MaxDiscord))
		}
	}
	return nil
}

// ParseClock lê "HH:MM" (00:00 a 23:59).
func ParseClock(s string) (int, int, bool) {
	h, m, ok := strings.Cut(strings.TrimSpace(s), ":")
	if !ok || len(h) != 2 || len(m) != 2 {
		return 0, 0, false
	}
	hh, err1 := strconv.Atoi(h)
	mm, err2 := strconv.Atoi(m)
	if err1 != nil || err2 != nil || hh < 0 || hh > 23 || mm < 0 || mm > 59 {
		return 0, 0, false
	}
	return hh, mm, true
}

// NextDaily é a próxima vez (depois de now) em que o relógio marca "HH:MM", no fuso de now.
func NextDaily(now time.Time, clock string) (time.Time, bool) {
	hh, mm, ok := ParseClock(clock)
	if !ok {
		return time.Time{}, false
	}
	t := time.Date(now.Year(), now.Month(), now.Day(), hh, mm, 0, 0, now.Location())
	if !t.After(now) {
		t = time.Date(now.Year(), now.Month(), now.Day()+1, hh, mm, 0, 0, now.Location())
	}
	return t, true
}

// AnnouncementCommand monta o "say" de uma mensagem programada (vazio se não há o que mandar).
func AnnouncementCommand(text string) string { return ChatCommand(text, "") }

// WarningText é o aviso de reinício para o chat ({minutes} vira o número).
func WarningText(template string, minutes int) string {
	return strings.ReplaceAll(template, MinutesTag, strconv.Itoa(minutes))
}

// SayCommand monta o "say" de um texto digitado pelo administrador na tela.
func SayCommand(text string) (string, error) {
	text = strings.TrimSpace(text)
	if text == "" || utf8.RuneCountInString(text) > MaxChat || hasControl(text, false) {
		return "", problems.Field("say", "validation.chatMessage", "max", strconv.Itoa(MaxChat))
	}
	return ChatCommand(text, ""), nil
}

// DiscordText troca os marcadores de um aviso do Discord. Nomes de jogador e de servidor não podem mencionar
// ninguém (@everyone) nem formatar o texto: os caracteres especiais do Discord são escapados.
func DiscordText(template, serverName, mapName, player string) string {
	template = strings.TrimSpace(template)
	if template == "" || hasControl(template, false) {
		return ""
	}
	r := strings.NewReplacer(ServerTag, discordEscape(serverName), MapTag, discordEscape(mapName),
		NameTag, discordEscape(cleanName(player)))
	text := r.Replace(template)
	if utf8.RuneCountInString(text) > 1900 {
		text = string([]rune(text)[:1900])
	}
	return text
}

func discordEscape(s string) string {
	var out strings.Builder
	for _, r := range s {
		switch r {
		case '\\', '*', '_', '~', '`', '|', '>', '#', '[', ']', '(', ')':
			out.WriteRune('\\')
			out.WriteRune(r)
		case '@':
			out.WriteString("@​") // sem menção
		default:
			if r >= 32 && r != 127 {
				out.WriteRune(r)
			}
		}
	}
	return out.String()
}
