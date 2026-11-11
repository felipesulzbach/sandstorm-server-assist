package desktop

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"sandstorm-server-assist/desktop/internal/app/logbuf"
	"sandstorm-server-assist/desktop/internal/domain/report"
)

func readZip(t *testing.T, b []byte) map[string]string {
	t.Helper()
	zr, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]string{}
	for _, f := range zr.File {
		rc, err := f.Open()
		if err != nil {
			t.Fatal(err)
		}
		data, _ := io.ReadAll(rc)
		_ = rc.Close()
		out[f.Name] = string(data)
	}
	return out
}

func TestReportZipScrubsLogAndChecksImage(t *testing.T) {
	log := logbuf.New(100, func() []string { return []string{"supersecret"} })
	log.Add("server", "Join from 203.0.113.9 id 76561198000000001 pass supersecret", "info")
	r := &ReportAPI{a: &app{log: log}}

	png := filepath.Join(t.TempDir(), "shot.png")
	if err := os.WriteFile(png, []byte("\x89PNG\r\n\x1a\nrest-of-image"), 0o600); err != nil {
		t.Fatal(err)
	}
	b, err := r.build(report.Meta{Kind: report.Bug, Title: "t"}, "desc", true, png)
	if err != nil {
		t.Fatal(err)
	}
	files := readZip(t, b)
	logText := files["app-log.txt"]
	for _, leak := range []string{"203.0.113.9", "76561198000000001", "supersecret"} {
		if strings.Contains(logText, leak) {
			t.Fatalf("log leaks %q: %s", leak, logText)
		}
	}
	if _, ok := files["screenshot.png"]; !ok || !strings.Contains(files["report.txt"], "desc") {
		t.Fatalf("files %v", files)
	}

	// Sem log e sem imagem: só o report.txt.
	b, err = r.build(report.Meta{Kind: report.Idea}, "", false, "")
	if err != nil {
		t.Fatal(err)
	}
	if files := readZip(t, b); len(files) != 1 {
		t.Fatalf("files %v", files)
	}

	// Arquivo que não é imagem (mesmo com nome .png) é recusado.
	fake := filepath.Join(t.TempDir(), "fake.png")
	if err := os.WriteFile(fake, []byte("MZ\x90\x00not an image"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := r.build(report.Meta{Kind: report.Bug}, "", false, fake); err == nil {
		t.Fatal("accepted a file that is not an image")
	}
}
