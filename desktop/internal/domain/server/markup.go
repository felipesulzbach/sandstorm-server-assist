package server

import (
	"regexp"
	"strings"

	"sandstorm-server-assist/desktop/internal/domain/problems"
)

// Formatação das mensagens (ADR-024). O jogo mostra a mensagem do dia com texto rico do Unreal: a tabela de estilos
// MotdTextStyleSet tem as linhas Default, Bold, Header1 e Header2, e o decorador WebImageDecorator aceita
// <extimg src="https://..."/> (lidos dos arquivos do servidor, versão 1.22). O chat tem a tabela ChatTextStyleSet
// com Default e Bold; se o "say" passa pelo texto rico não foi confirmado com o jogo aberto.
//
// Sintaxe do Unreal: <Estilo>texto</> e <decorador atributo="valor"/>. Sem cores livres: a cor vem do estilo.

// Estilos aceitos.
var (
	MotdStyles = []string{"Bold", "Header1", "Header2"}
	ChatStyles = []string{"Bold"}
)

// ImageTag é o decorador de imagem da mensagem do dia.
const ImageTag = "extimg"

// Segment é um pedaço do texto para a prévia: Style vazio = texto normal; Src = imagem.
type Segment struct {
	Style string `json:"style"`
	Text  string `json:"text"`
	Src   string `json:"src,omitempty"`
}

// reTag acha o que o Unreal trata como tag: <Nome>, <Nome a="b"/>, <Nome a="b">, </>.
var (
	reTag   = regexp.MustCompile(`<(/?)([A-Za-z0-9_.\-]*)((?:\s+[A-Za-z0-9_.\-]+="[^"]*")*)\s*(/?)>`)
	reAttr  = regexp.MustCompile(`([A-Za-z0-9_.\-]+)="([^"]*)"`)
	reImage = regexp.MustCompile(`^https://[^\s"<>]{3,300}$`)
)

// ParseMarkup confere a formatação e devolve os pedaços. chat=true aceita só os estilos do chat e nenhuma imagem.
// field é o campo apontado no erro.
func ParseMarkup(text string, chat bool, field string) ([]Segment, error) {
	styles := MotdStyles
	if chat {
		styles = ChatStyles
	}
	known := map[string]bool{}
	for _, s := range styles {
		known[s] = true
	}
	var out []Segment
	open := ""
	add := func(s Segment) {
		if s.Text == "" && s.Src == "" {
			return
		}
		if n := len(out); n > 0 && s.Src == "" && out[n-1].Src == "" && out[n-1].Style == s.Style {
			out[n-1].Text += s.Text
			return
		}
		out = append(out, s)
	}
	last := 0
	for _, m := range reTag.FindAllStringSubmatchIndex(text, -1) {
		add(Segment{Style: open, Text: text[last:m[0]]})
		last = m[1]
		closing, name := text[m[2]:m[3]] == "/", text[m[4]:m[5]]
		attrs, selfClose := text[m[6]:m[7]], text[m[8]:m[9]] == "/"
		switch {
		case closing && name == "":
			if open == "" {
				return nil, problems.Field(field, "validation.markupClose")
			}
			open = ""
		case !chat && strings.EqualFold(name, ImageTag) && selfClose:
			src := ""
			for _, a := range reAttr.FindAllStringSubmatch(attrs, -1) {
				if a[1] == "src" {
					src = a[2]
				}
			}
			if !reImage.MatchString(src) {
				return nil, problems.Field(field, "validation.markupImage")
			}
			add(Segment{Src: src})
		case known[name] && !selfClose && attrs == "":
			if open != "" {
				return nil, problems.Field(field, "validation.markupNested")
			}
			open = name
		default:
			return nil, problems.Field(field, "validation.markupTag", "tag", text[m[0]:m[1]])
		}
	}
	if open != "" {
		return nil, problems.Field(field, "validation.markupUnclosed", "tag", "<"+open+">")
	}
	add(Segment{Style: open, Text: text[last:]})
	return out, nil
}

// validateMarkup confere a formatação da mensagem do dia e das mensagens de chat.
func (p Profile) validateMarkup() error {
	if _, err := ParseMarkup(p.Motd, false, "motd"); err != nil {
		return err
	}
	for _, f := range []struct{ name, v string }{{"joinMessage", p.JoinMessage}, {"leaveMessage", p.LeaveMessage}} {
		if _, err := ParseMarkup(f.v, true, f.name); err != nil {
			return err
		}
	}
	for _, a := range p.Announcements {
		if _, err := ParseMarkup(a, true, "announcements"); err != nil {
			return err
		}
	}
	return nil
}
