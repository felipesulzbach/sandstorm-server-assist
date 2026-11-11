package server

import (
	"regexp"

	"sandstorm-server-assist/desktop/internal/domain/problems"
)

// Mensagens sem formatação (ADR-031). Os arquivos do jogo têm tabelas de estilo de texto rico (MotdTextStyleSet,
// ChatTextStyleSet) e o ADR-024 pôs negrito, títulos e imagem nas mensagens. Visto no jogo em 11/11/2026: a tela de
// carregamento mostrou "<Bold>" escrito vindo do Motd.txt, e o chat mostrou o mesmo vindo do "say". Nenhuma das
// duas interpreta as tags; as mensagens são texto simples.

// reStyleTag são as tags de formatação que o app punha nas mensagens (<Bold>, <Header1>, <Header2>, </> e a
// imagem). Só essas são recusadas e tiradas: outro texto entre < e > (como "<3" ou "a < b") fica como está.
var reStyleTag = regexp.MustCompile(`(?i)</>|<(?:Bold|Header1|Header2)>|<extimg(?:\s+[A-Za-z0-9_.\-]+="[^"]*")*\s*/>`)

// StripMarkup tira as tags de formatação de uma mensagem salva antes da 0.1.10, quando o editor deixava formatar.
// Vale para o Motd.txt e para o chat.
func StripMarkup(text string) string { return reStyleTag.ReplaceAllString(text, "") }

// checkPlain recusa uma mensagem com tags: o jogo mostraria a tag escrita. field é o campo apontado no erro.
func checkPlain(text, field string) error {
	if tag := reStyleTag.FindString(text); tag != "" {
		return problems.Field(field, "validation.markupPlain", "tag", tag)
	}
	return nil
}

// validateMarkup confere que a mensagem do dia e as mensagens de chat não têm tags.
func (p Profile) validateMarkup() error {
	for _, f := range []struct{ name, v string }{
		{"motd", p.Motd}, {"joinMessage", p.JoinMessage}, {"leaveMessage", p.LeaveMessage},
	} {
		if err := checkPlain(f.v, f.name); err != nil {
			return err
		}
	}
	for _, a := range p.Announcements {
		if err := checkPlain(a, "announcements"); err != nil {
			return err
		}
	}
	return nil
}
