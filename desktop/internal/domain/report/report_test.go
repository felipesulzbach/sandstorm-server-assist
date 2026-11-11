package report

import (
	"net/url"
	"strings"
	"testing"
)

func TestClean(t *testing.T) {
	if got := Clean("  a\x00b\r\nc\td  ", 100, true); got != "ab\nc\td" {
		t.Fatalf("lines: %q", got)
	}
	if got := Clean("a\nb", 100, false); got != "a b" {
		t.Fatalf("single line: %q", got)
	}
	if got := Clean("ação", 2, false); got != "aç" {
		t.Fatalf("limit: %q", got)
	}
}

func TestScrub(t *testing.T) {
	in := "join from 203.0.113.9:7777 id 76561198000000001, rcon on 127.0.0.1 and bind 0.0.0.0"
	got := Scrub(in)
	if strings.Contains(got, "203.0.113.9") || strings.Contains(got, "76561198000000001") {
		t.Fatalf("not scrubbed: %s", got)
	}
	if !strings.Contains(got, "127.0.0.1") || !strings.Contains(got, "0.0.0.0") {
		t.Fatalf("scrubbed too much: %s", got)
	}
}

func TestIssueLink(t *testing.T) {
	m := Meta{Kind: Bug, Title: "Crash & burn", AppVersion: "0.2.0", Windows: "Windows 11 (26100)"}
	u, err := url.Parse(IssueLink(m, strings.Repeat("é", 3000)))
	if err != nil {
		t.Fatal(err)
	}
	if u.Host != "github.com" || u.Path != "/felipesulzbach/sandstorm-server-assist/issues/new" {
		t.Fatalf("url %s", u)
	}
	q := u.Query()
	if q.Get("template") != "bug_report.yml" || q.Get("labels") != "bug" || q.Get("title") != "Crash & burn" {
		t.Fatalf("query %v", q)
	}
	if n := len([]rune(q.Get("description"))); n > maxURLText+4 {
		t.Fatalf("description not cut: %d", n)
	}
	if q2, _ := url.Parse(IssueLink(Meta{Kind: Idea}, "")); q2.Query().Get("template") != "feature_request.yml" {
		t.Fatal("idea template")
	}
}

func TestImageExt(t *testing.T) {
	cases := map[string]string{
		"\x89PNG\r\n\x1a\nxxxx":     ".png",
		"\xFF\xD8\xFF\xE0":          ".jpg",
		"GIF89a..":                  ".gif",
		"RIFF\x00\x00\x00\x00WEBPV": ".webp",
		"MZ\x90\x00":                "",
		"<svg":                      "",
	}
	for head, want := range cases {
		if got := ImageExt([]byte(head)); got != want {
			t.Errorf("%q: got %q want %q", head, got, want)
		}
	}
}

func TestTextAndDefaults(t *testing.T) {
	b := string(Text(Meta{Kind: Idea, Title: DefaultTitle(Idea, "0.2.0")}, ""))
	if !strings.Contains(b, "Idea for version 0.2.0") || !strings.Contains(b, "(none)") {
		t.Fatal(b)
	}
	if !ValidKind(Bug) || ValidKind("other") {
		t.Fatal("kind")
	}
}
