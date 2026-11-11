// Package report são as regras do relato de problema e da sugestão de melhoria (ADR-034). O app monta um zip no
// PC (texto, log mascarado e imagem opcional) e abre o formulário de issue do GitHub já preenchido; a pessoa
// anexa o zip lá. Nada é enviado pelo app. Este pacote não faz I/O.
package report

import (
	"bytes"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Tipos de relato.
const (
	Bug  = "bug"
	Idea = "idea"
)

// Limites.
const (
	MaxTitle       = 120     // caracteres do título
	MaxDescription = 4000    // caracteres da descrição
	MaxImageBytes  = 5 << 20 // imagem anexada
	MaxLogLines    = 3000    // linhas do log no zip
	MaxLogBytes    = 2 << 20 // tamanho do log no zip
	maxURLText     = 1500    // caracteres da descrição que vão no endereço do GitHub (o resto fica no zip)
)

// IssuesURL é o formulário de issue do projeto no GitHub.
const IssuesURL = "https://github.com/felipesulzbach/sandstorm-server-assist/issues/new"

// ValidKind diz se o tipo existe.
func ValidKind(kind string) bool { return kind == Bug || kind == Idea }

// Clean tira caracteres de controle (menos quebra de linha e tab), corta no limite e apara espaços.
func Clean(s string, max int, lines bool) string {
	s = strings.ToValidUTF8(s, "")
	s = strings.ReplaceAll(s, "\r\n", "\n")
	var b strings.Builder
	n := 0
	for _, r := range s {
		if n >= max {
			break
		}
		switch {
		case r == '\n' || r == '\t':
			if !lines {
				r = ' '
			}
		case unicode.IsControl(r):
			continue
		}
		b.WriteRune(r)
		n++
	}
	return strings.TrimSpace(b.String())
}

var (
	reIPv4    = regexp.MustCompile(`\b(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})\b`)
	reSteamID = regexp.MustCompile(`\b7656119\d{10}\b`)
)

// Scrub esconde dados de terceiros e da rede no texto que vai ao relato: IPs (menos 127.0.0.1 e 0.0.0.0, úteis
// para entender o problema) e SteamIDs de jogadores. Os segredos já saem mascarados do log do app.
func Scrub(s string) string {
	s = reIPv4.ReplaceAllStringFunc(s, func(ip string) string {
		if ip == "127.0.0.1" || ip == "0.0.0.0" {
			return ip
		}
		return "x.x.x.x"
	})
	return reSteamID.ReplaceAllString(s, "7656119xxxxxxxxxx")
}

// Meta são os dados do app e do PC que vão no relato.
type Meta struct {
	Kind       string
	Title      string
	AppVersion string
	Windows    string
	Language   string
	Created    string // data e hora, já formatada
}

// DefaultTitle é o título quando a pessoa não escreve um.
func DefaultTitle(kind, version string) string {
	if kind == Idea {
		return "Idea for version " + version
	}
	return "Problem in version " + version
}

// Text é o report.txt do zip.
func Text(m Meta, description string) []byte {
	var b bytes.Buffer
	fmt.Fprintf(&b, "Sandstorm Server Assist report\n\nType: %s\nTitle: %s\nApp version: %s\nWindows: %s\nLanguage: %s\nCreated: %s\n\n",
		m.Kind, m.Title, m.AppVersion, m.Windows, m.Language, m.Created)
	b.WriteString("Description:\n")
	if description == "" {
		b.WriteString("(none)\n")
	} else {
		b.WriteString(description + "\n")
	}
	return b.Bytes()
}

// IssueLink monta o endereço do formulário do GitHub já preenchido (modelos em .github/ISSUE_TEMPLATE). Os campos
// vão pelo id do formulário; a descrição longa é cortada (fica inteira no zip).
func IssueLink(m Meta, description string) string {
	template, label := "bug_report.yml", "bug"
	if m.Kind == Idea {
		template, label = "feature_request.yml", "enhancement"
	}
	if utf8.RuneCountInString(description) > maxURLText {
		description = string([]rune(description)[:maxURLText]) + " […]"
	}
	q := url.Values{
		"template":    {template},
		"labels":      {label},
		"title":       {m.Title},
		"version":     {m.AppVersion},
		"windows":     {m.Windows},
		"description": {description},
	}
	return IssuesURL + "?" + q.Encode()
}

// ImageExt diz o tipo da imagem pelo começo do arquivo (não pelo nome). "" se não for uma imagem aceita.
func ImageExt(head []byte) string {
	switch {
	case bytes.HasPrefix(head, []byte("\x89PNG\r\n\x1a\n")):
		return ".png"
	case bytes.HasPrefix(head, []byte{0xFF, 0xD8, 0xFF}):
		return ".jpg"
	case bytes.HasPrefix(head, []byte("GIF87a")), bytes.HasPrefix(head, []byte("GIF89a")):
		return ".gif"
	case len(head) >= 12 && bytes.Equal(head[0:4], []byte("RIFF")) && bytes.Equal(head[8:12], []byte("WEBP")):
		return ".webp"
	}
	return ""
}
