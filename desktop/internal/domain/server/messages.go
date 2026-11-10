package server

import (
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"sandstorm-server-assist/desktop/internal/domain/problems"
)

// Mensagens do servidor (ADR-020).
//
//   - Mensagem do dia: o servidor lê Insurgency/Config/Server/Motd.txt ao iniciar e mostra na tela de carregamento.
//   - Entrada e saída: o jogo não tem texto próprio para isso. O app lê no log quem entrou ou saiu e manda a
//     mensagem para o chat com o comando "say" do RCON. {name} vira o nome do jogador.

// Limites das mensagens.
const (
	MaxMotd     = 1000
	MaxChat     = 150
	maxName     = 40
	maxChatSent = 200
	// NameTag é o marcador do nome do jogador nas mensagens.
	NameTag = "{name}"
)

// validateMessages confere a mensagem do dia e as de entrada e saída.
func (p Profile) validateMessages() error {
	if utf8.RuneCountInString(p.Motd) > MaxMotd || hasControl(p.Motd, true) {
		return problems.Field("motd", "validation.motd", "max", strconv.Itoa(MaxMotd))
	}
	for _, f := range []struct{ name, v string }{{"joinMessage", p.JoinMessage}, {"leaveMessage", p.LeaveMessage}} {
		if f.v == "" {
			continue
		}
		if utf8.RuneCountInString(f.v) > MaxChat || hasControl(f.v, false) {
			return problems.Field(f.name, "validation.chatMessage", "max", strconv.Itoa(MaxChat))
		}
		if !p.Rcon {
			return problems.Field(f.name, "validation.chatNeedsRcon")
		}
	}
	return nil
}

// hasControl diz se o texto tem caracteres de controle (quebra de linha só é aceita quando lines é true).
func hasControl(s string, lines bool) bool {
	for _, r := range s {
		if lines && (r == '\n' || r == '\r') {
			continue
		}
		if unicode.IsControl(r) {
			return true
		}
	}
	return false
}

// MotdFile é o conteúdo do Motd.txt (linhas com CRLF, sem espaços sobrando no fim).
func MotdFile(text string) string {
	text = strings.TrimSpace(strings.ReplaceAll(text, "\r\n", "\n"))
	if text == "" {
		return ""
	}
	return strings.ReplaceAll(text, "\n", "\r\n") + "\r\n"
}

// ParseMotd lê um Motd.txt existente. Texto inválido (longo demais ou com controle) fica de fora.
func ParseMotd(content string) (string, bool) {
	text := strings.TrimSpace(strings.ReplaceAll(content, "\r\n", "\n"))
	if text == "" || utf8.RuneCountInString(text) > MaxMotd || hasControl(text, true) {
		return "", false
	}
	return text, true
}

// cleanName tira do nome do jogador o que não pode ir para o chat (controle, espaços repetidos) e corta o tamanho.
func cleanName(name string) string {
	name = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return ' '
		}
		return r
	}, name)
	name = strings.Join(strings.Fields(name), " ")
	if utf8.RuneCountInString(name) > maxName {
		name = string([]rune(name)[:maxName])
	}
	return name
}

// ChatCommand monta o comando "say" com a mensagem e o nome do jogador. Vazio se não há o que mandar.
func ChatCommand(template, player string) string {
	template = strings.TrimSpace(template)
	if template == "" || hasControl(template, false) {
		return ""
	}
	text := strings.ReplaceAll(template, NameTag, cleanName(player))
	text = strings.Join(strings.Fields(text), " ")
	if utf8.RuneCountInString(text) > maxChatSent {
		text = string([]rune(text)[:maxChatSent])
	}
	if text == "" {
		return ""
	}
	return "say " + text
}
