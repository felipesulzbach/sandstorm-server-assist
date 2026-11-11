package desktop

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"sandstorm-server-assist/desktop/internal/adapters/winsys"
	"sandstorm-server-assist/desktop/internal/domain/problems"
	"sandstorm-server-assist/desktop/internal/domain/redact"
	"sandstorm-server-assist/desktop/internal/domain/report"
)

// ReportAPI: relato de problema e sugestão de melhoria (ADR-034). O app monta um zip no PC e abre o formulário do
// GitHub já preenchido; a pessoa anexa o zip. Nada sai do PC sem ela.
type ReportAPI struct {
	a     *app
	mu    sync.Mutex
	image string // caminho da imagem escolhida no diálogo do Windows ("" = nenhuma)
	last  string // último zip criado
}

// ImageInfo é o que a tela sabe da imagem escolhida (nunca o caminho).
type ImageInfo struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

// ReportRequest é o que a tela manda: tipo, textos e se o log vai junto.
type ReportRequest struct {
	Kind        string `json:"kind"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IncludeLog  bool   `json:"includeLog"`
}

// ReportResult é o zip criado.
type ReportResult struct {
	File string `json:"file"`
	Size int64  `json:"size"`
}

// ChooseImage abre o diálogo do Windows para escolher a imagem (opcional). Confere tamanho e tipo pelo conteúdo.
func (r *ReportAPI) ChooseImage() (ImageInfo, error) {
	path, err := runtime.OpenFileDialog(r.a.ctx, runtime.OpenDialogOptions{
		Filters: []runtime.FileFilter{{DisplayName: "Images (*.png, *.jpg, *.gif, *.webp)", Pattern: "*.png;*.jpg;*.jpeg;*.gif;*.webp"}},
	})
	if err != nil || path == "" {
		return ImageInfo{}, err
	}
	st, err := os.Stat(path)
	if err != nil || !st.Mode().IsRegular() {
		return ImageInfo{}, problems.New(problems.ReportImage, "not a file")
	}
	if st.Size() > report.MaxImageBytes {
		return ImageInfo{}, problems.New(problems.ReportImage, "too large")
	}
	if _, err := imageExt(path); err != nil {
		return ImageInfo{}, err
	}
	r.mu.Lock()
	r.image = path
	r.mu.Unlock()
	return ImageInfo{Name: filepath.Base(path), Size: st.Size()}, nil
}

// ClearImage tira a imagem escolhida.
func (r *ReportAPI) ClearImage() {
	r.mu.Lock()
	r.image = ""
	r.mu.Unlock()
}

// imageExt lê o começo do arquivo e devolve a extensão do tipo real.
func imageExt(path string) (string, error) {
	f, err := os.Open(path) // #nosec G304 -- arquivo escolhido pela pessoa no diálogo do Windows
	if err != nil {
		return "", problems.Wrap(problems.ReportImage, err)
	}
	defer func() { _ = f.Close() }()
	head := make([]byte, 16)
	n, _ := io.ReadFull(f, head)
	ext := report.ImageExt(head[:n])
	if ext == "" {
		return "", problems.New(problems.ReportImage, "not an image")
	}
	return ext, nil
}

// Create monta o zip do relato, abre o formulário do GitHub e mostra o zip no Explorer.
func (r *ReportAPI) Create(req ReportRequest) (ReportResult, error) {
	if !report.ValidKind(req.Kind) {
		return ReportResult{}, problems.New(problems.Unknown, "invalid report kind")
	}
	secrets := r.a.cfg.SecretValues()
	scrub := func(s string, max int, lines bool) string {
		return report.Scrub(redact.Values(redact.Text(report.Clean(s, max, lines)), secrets...))
	}
	meta := report.Meta{
		Kind:       req.Kind,
		Title:      scrub(req.Title, report.MaxTitle, false),
		AppVersion: AppVersion,
		Windows:    winsys.WindowsVersion(),
		Language:   r.a.cfg.Settings().Language,
		Created:    time.Now().Format("2006-01-02 15:04"),
	}
	if meta.Title == "" {
		meta.Title = report.DefaultTitle(req.Kind, AppVersion)
	}
	desc := scrub(req.Description, report.MaxDescription, true)

	r.mu.Lock()
	image := r.image
	r.mu.Unlock()

	zipBytes, err := r.build(meta, desc, req.IncludeLog, image)
	if err != nil {
		return ReportResult{}, err
	}
	dir, err := reportsDir()
	if err != nil {
		return ReportResult{}, problems.Wrap(problems.ReportFailed, err)
	}
	name := "sandstorm-server-assist-" + req.Kind + "-" + time.Now().Format("20060102-150405") + ".zip"
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, zipBytes, 0o600); err != nil {
		return ReportResult{}, problems.Wrap(problems.ReportFailed, err)
	}
	r.mu.Lock()
	r.last = path
	r.image = ""
	r.mu.Unlock()

	runtime.BrowserOpenURL(r.a.ctx, report.IssueLink(meta, desc))
	_ = winsys.RevealFile(path)
	r.a.log.Add("app", "Report saved: "+name, "info")
	return ReportResult{File: name, Size: int64(len(zipBytes))}, nil
}

// ShowLast mostra no Explorer o último zip criado.
func (r *ReportAPI) ShowLast() error {
	r.mu.Lock()
	path := r.last
	r.mu.Unlock()
	if path == "" {
		return problems.New(problems.ReportFailed, "no report yet")
	}
	return winsys.RevealFile(path)
}

// build monta o zip em memória: report.txt, app-log.txt (opcional) e a imagem (opcional).
func (r *ReportAPI) build(meta report.Meta, desc string, withLog bool, image string) ([]byte, error) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	add := func(name string, data []byte) error {
		w, err := zw.Create(name)
		if err != nil {
			return err
		}
		_, err = w.Write(data)
		return err
	}
	if err := add("report.txt", report.Text(meta, desc)); err != nil {
		return nil, problems.Wrap(problems.ReportFailed, err)
	}
	if withLog {
		if err := add("app-log.txt", r.logText()); err != nil {
			return nil, problems.Wrap(problems.ReportFailed, err)
		}
	}
	if image != "" {
		ext, err := imageExt(image)
		if err != nil {
			return nil, err
		}
		data, err := readImage(image)
		if err != nil {
			return nil, err
		}
		if err := add("screenshot"+ext, data); err != nil {
			return nil, problems.Wrap(problems.ReportFailed, err)
		}
	}
	if err := zw.Close(); err != nil {
		return nil, problems.Wrap(problems.ReportFailed, err)
	}
	return buf.Bytes(), nil
}

// logText são as últimas linhas do log do app (já com os segredos mascarados), sem IPs nem SteamIDs, com limite
// de tamanho: guarda as mais recentes.
func (r *ReportAPI) logText() []byte {
	lines := r.a.log.Recent(report.MaxLogLines)
	rows := make([]string, 0, len(lines))
	size := 0
	for i := len(lines) - 1; i >= 0; i-- {
		l := lines[i]
		row := l.Time.Format("2006-01-02 15:04:05") + " [" + l.Source + "] [" + l.Level + "] " + report.Scrub(l.Text) + "\n"
		if size+len(row) > report.MaxLogBytes {
			break
		}
		size += len(row)
		rows = append(rows, row)
	}
	var out bytes.Buffer
	for i := len(rows) - 1; i >= 0; i-- {
		out.WriteString(rows[i])
	}
	return out.Bytes()
}

// readImage lê a imagem até o limite (o arquivo pode ter mudado depois da escolha).
func readImage(path string) ([]byte, error) {
	f, err := os.Open(path) // #nosec G304 -- arquivo escolhido pela pessoa no diálogo do Windows
	if err != nil {
		return nil, problems.Wrap(problems.ReportImage, err)
	}
	defer func() { _ = f.Close() }()
	data, err := io.ReadAll(io.LimitReader(f, report.MaxImageBytes+1))
	if err != nil {
		return nil, problems.Wrap(problems.ReportImage, err)
	}
	if len(data) > report.MaxImageBytes {
		return nil, problems.New(problems.ReportImage, "too large")
	}
	return data, nil
}

// reportsDir é a pasta dos relatos (%LOCALAPPDATA%\sandstorm-server-assist\reports).
func reportsDir() (string, error) {
	base, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	d := filepath.Join(base, "sandstorm-server-assist", "reports")
	if err := os.MkdirAll(d, 0o700); err != nil {
		return "", err
	}
	return d, nil
}
