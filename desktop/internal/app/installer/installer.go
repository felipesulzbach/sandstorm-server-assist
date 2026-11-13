// Package installer baixa o SteamCMD e instala ou atualiza o servidor, informando o progresso (ADR-007).
package installer

import (
	"archive/zip"
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"
	"time"

	"sandstorm-server-assist/desktop/internal/adapters/fsx"
	"sandstorm-server-assist/desktop/internal/adapters/netx"
	"sandstorm-server-assist/desktop/internal/adapters/proc"
	"sandstorm-server-assist/desktop/internal/adapters/trust"
	"sandstorm-server-assist/desktop/internal/app/logbuf"
	"sandstorm-server-assist/desktop/internal/domain/problems"
	"sandstorm-server-assist/desktop/internal/domain/server"
	"sandstorm-server-assist/desktop/internal/domain/steamcmd"
)

// SteamCMDURL é o endereço oficial do SteamCMD para Windows.
const SteamCMDURL = "https://steamcdn-a.akamaihd.net/client/installer/steamcmd.zip"

// Task é a tarefa em andamento.
type Task string

// Tarefas.
const (
	Idle          Task = ""
	FetchSteamCMD Task = "steamcmd"
	InstallServer Task = "server"
)

// State é o que a tela de instalação mostra.
type State struct {
	Task    Task            `json:"task"`
	Phase   steamcmd.Phase  `json:"phase"`
	Percent float64         `json:"percent"`
	Done    int64           `json:"done"`
	Total   int64           `json:"total"`
	Speed   float64         `json:"speed"` // bytes por segundo
	ETA     int64           `json:"eta"`   // segundos (0 = sem estimativa)
	Started time.Time       `json:"started"`
	Error   *problems.Error `json:"error,omitempty"`
	// Finished indica que a última tarefa terminou bem (para a tela mostrar "concluído").
	Finished Task `json:"finished"`
	// Last é a última tarefa iniciada (em andamento, concluída ou com erro): a tela mostra o andamento e o resultado
	// no cartão dela, perto do botão que a iniciou.
	Last Task `json:"last"`
	// Validate: a tarefa do servidor confere todos os arquivos ("Verificar ficheiros").
	Validate bool `json:"validate"`
	// UpToDate: o servidor já estava na versão mais recente (nada foi baixado).
	UpToDate bool `json:"upToDate"`
	// Present: o SteamCMD já estava na pasta, com a assinatura da Valve conferida (nada foi baixado).
	Present bool `json:"present"`
}

// Info é o que já está instalado.
type Info struct {
	SteamCMDFound   bool   `json:"steamcmdFound"`
	ServerInstalled bool   `json:"serverInstalled"`
	BuildID         string `json:"buildId"`
	SizeOnDisk      int64  `json:"sizeOnDisk"`
	LastUpdated     int64  `json:"lastUpdated"`
	// FoldersNested: a pasta do SteamCMD está dentro da do servidor, ou o contrário; o SteamCMD não instalaria na
	// pasta do servidor (ADR-042).
	FoldersNested bool `json:"foldersNested"`
}

// Service faz uma tarefa por vez.
type Service struct {
	mu     sync.Mutex
	st     State
	cancel context.CancelFunc
	notify func(State)
	log    *logbuf.Buffer
	// speed: amostras para calcular velocidade.
	lastDone int64
	lastAt   time.Time
}

// New cria o serviço.
func New(log *logbuf.Buffer) *Service { return &Service{log: log} }

// SetNotify define quem recebe o estado.
func (s *Service) SetNotify(fn func(State)) {
	s.mu.Lock()
	s.notify = fn
	s.mu.Unlock()
}

// State devolve o estado atual.
func (s *Service) State() State {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.st
}

// Busy diz se há tarefa em andamento.
func (s *Service) Busy() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.st.Task != Idle
}

func (s *Service) set(fn func(*State)) {
	s.mu.Lock()
	fn(&s.st)
	st, n := s.st, s.notify
	s.mu.Unlock()
	if n != nil {
		n(st)
	}
}

// progress atualiza bytes e calcula velocidade e tempo restante (média simples).
func (s *Service) progress(phase steamcmd.Phase, done, total int64, pct float64) {
	now := time.Now()
	s.set(func(st *State) {
		st.Phase, st.Done, st.Total = phase, done, total
		if pct > 0 {
			st.Percent = pct
		} else if total > 0 {
			st.Percent = float64(done) * 100 / float64(total)
		}
		if !s.lastAt.IsZero() && done >= s.lastDone {
			if dt := now.Sub(s.lastAt).Seconds(); dt >= 0.8 {
				inst := float64(done-s.lastDone) / dt
				if st.Speed == 0 {
					st.Speed = inst
				} else {
					st.Speed = st.Speed*0.7 + inst*0.3
				}
				s.lastDone, s.lastAt = done, now
			}
		} else {
			s.lastDone, s.lastAt = done, now
		}
		st.ETA = 0
		if st.Speed > 1 && total > done {
			st.ETA = int64(float64(total-done) / st.Speed)
		}
	})
}

// begin marca a tarefa e avisa a tela na hora: antes a tela só sabia da tarefa no primeiro progresso, e o SteamCMD
// passa um tempo se atualizando sem progresso (os botões ficavam ativos e nada aparecia).
func (s *Service) begin(task Task, validate bool) (context.Context, error) {
	s.mu.Lock()
	if s.st.Task != Idle {
		s.mu.Unlock()
		return nil, problems.New(problems.Busy, "")
	}
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	s.st = State{Task: task, Last: task, Validate: validate, Started: time.Now()}
	s.lastDone, s.lastAt = 0, time.Time{}
	st, n := s.st, s.notify
	s.mu.Unlock()
	if n != nil {
		n(st)
	}
	return ctx, nil
}

func (s *Service) end(task Task, err error) {
	s.mu.Lock()
	if s.cancel != nil {
		s.cancel()
		s.cancel = nil
	}
	s.mu.Unlock()
	s.set(func(st *State) {
		st.Task = Idle
		st.Speed, st.ETA = 0, 0
		if err != nil {
			var pe *problems.Error
			if !errors.As(err, &pe) {
				pe = problems.New(problems.Unknown, err.Error())
			}
			st.Error = pe
			st.Finished = Idle
		} else {
			st.Error = nil
			st.Finished = task
			st.Percent = 100
		}
	})
}

// Cancel pede para parar a tarefa atual.
func (s *Service) Cancel() {
	s.mu.Lock()
	c := s.cancel
	s.mu.Unlock()
	if c != nil {
		c()
	}
}

// SteamCMDExe é o caminho do executável numa pasta.
func SteamCMDExe(dir string) string { return filepath.Join(dir, "steamcmd.exe") }

// Inspect diz o que está instalado.
func Inspect(p server.Profile) Info {
	var in Info
	in.FoldersNested = steamcmd.NestedDirs(p.ServerDir, p.SteamCMDDir)
	if p.SteamCMDDir != "" {
		in.SteamCMDFound = fsx.Exists(SteamCMDExe(p.SteamCMDDir))
	}
	if p.ServerDir != "" {
		in.ServerInstalled = fsx.Exists(filepath.Join(p.ServerDir, "InsurgencyServer.exe"))
		if m, err := fsx.ReadLimited(manifestPath(p.ServerDir), 1<<20); err == nil {
			in.BuildID = acf(m, "buildid")
			in.SizeOnDisk, _ = strconv.ParseInt(acf(m, "SizeOnDisk"), 10, 64)
			in.LastUpdated, _ = strconv.ParseInt(acf(m, "LastUpdated"), 10, 64)
		}
	}
	return in
}

func manifestPath(serverDir string) string {
	return filepath.Join(serverDir, "steamapps", "appmanifest_"+strconv.Itoa(server.AppID)+".acf")
}

// acf lê um valor de primeiro nível do appmanifest ("chave"		"valor").
func acf(b []byte, key string) string {
	re := regexp.MustCompile(`"` + regexp.QuoteMeta(key) + `"\s+"([^"]*)"`)
	if m := re.FindSubmatch(b); m != nil {
		return string(m[1])
	}
	return ""
}

// LauncherDir lê do appmanifest a pasta do SteamCMD usado na última instalação (para importar).
func LauncherDir(serverDir string) string {
	if m, err := fsx.ReadLimited(manifestPath(serverDir), 1<<20); err == nil {
		if p := acf(m, "LauncherPath"); p != "" {
			p = filepath.Clean(regexp.MustCompile(`\\\\`).ReplaceAllString(p, `\`))
			if filepath.IsAbs(p) && fsx.Exists(p) {
				return filepath.Dir(p)
			}
		}
	}
	return ""
}

// FetchSteamCMD baixa o SteamCMD para dir (em segundo plano). A assinatura da Valve é conferida antes de usar.
func (s *Service) FetchSteamCMD(dir string) error {
	if !fsx.SafeDir(dir) {
		return problems.New(problems.ServerDirInvalid, dir)
	}
	ctx, err := s.begin(FetchSteamCMD, false)
	if err != nil {
		return err
	}
	go func() { s.end(FetchSteamCMD, s.fetchSteamCMD(ctx, dir)) }()
	return nil
}

func (s *Service) fetchSteamCMD(ctx context.Context, dir string) error {
	exe := SteamCMDExe(dir)
	if fsx.Exists(exe) {
		if err := trust.VerifySigner(exe, "Valve"); err == nil {
			s.logf("SteamCMD already present: " + exe)
			s.set(func(st *State) { st.Present = true })
			return nil
		}
	}
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return problems.Wrap(problems.ServerDirInvalid, err)
	}
	tmpZip := filepath.Join(dir, ".steamcmd-download.zip")
	_ = os.Remove(tmpZip)
	defer os.Remove(tmpZip)
	s.logf("Downloading " + SteamCMDURL)
	s.set(func(st *State) { st.Phase = steamcmd.Downloading })
	err := netx.Download(ctx, SteamCMDURL, tmpZip, 64<<20, func(done, total int64) {
		s.progress(steamcmd.Downloading, done, total, 0)
	})
	if err != nil {
		if ctx.Err() != nil {
			return problems.New(problems.Canceled, "")
		}
		return netx.NetProblem(problems.SteamCMDDownload, err)
	}
	data, err := extractOne(tmpZip, "steamcmd.exe", 32<<20)
	if err != nil {
		return problems.Wrap(problems.SteamCMDDownload, err)
	}
	// Grava com outro nome, confere a assinatura e só então troca pelo nome final.
	staged := exe + ".new"
	if err := os.WriteFile(staged, data, 0o600); err != nil {
		return problems.Wrap(problems.SteamCMDDownload, err)
	}
	if err := trust.VerifySigner(staged, "Valve"); err != nil {
		_ = os.Remove(staged)
		return problems.Wrap(problems.SteamCMDSignature, err)
	}
	if err := os.Rename(staged, exe); err != nil {
		_ = os.Remove(staged)
		return problems.Wrap(problems.SteamCMDDownload, err)
	}
	s.logf("SteamCMD ready: " + exe)
	return nil
}

func extractOne(zipPath, name string, limit int64) ([]byte, error) {
	zr, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}
	defer zr.Close()
	for _, f := range zr.File {
		if f.Name != name {
			continue
		}
		if f.UncompressedSize64 > uint64(limit) { // #nosec G115 -- limit é constante positiva (32 MiB)
			return nil, errors.New("steamcmd.exe too large")
		}
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		defer rc.Close()
		return io.ReadAll(io.LimitReader(rc, limit))
	}
	return nil, errors.New("steamcmd.exe not found in zip")
}

// Install instala ou atualiza o servidor (em segundo plano). validate=true confere todos os arquivos.
func (s *Service) Install(steamcmdDir, serverDir string, validate bool) error {
	exe := SteamCMDExe(steamcmdDir)
	if !fsx.Exists(exe) {
		return problems.New(problems.SteamCMDFailed, "steamcmd.exe not found")
	}
	if !fsx.SafeDir(serverDir) {
		return problems.New(problems.ServerDirInvalid, serverDir)
	}
	// O SteamCMD não instala numa pasta que contém a dele: instalaria o servidor em steamapps\common (ADR-042).
	if steamcmd.NestedDirs(serverDir, steamcmdDir) {
		return problems.New(problems.FoldersNested, "")
	}
	if err := trust.VerifySigner(exe, "Valve"); err != nil {
		return problems.Wrap(problems.SteamCMDSignature, err)
	}
	ctx, err := s.begin(InstallServer, validate)
	if err != nil {
		return err
	}
	go func() { s.end(InstallServer, s.install(ctx, exe, serverDir, validate)) }()
	return nil
}

func (s *Service) install(ctx context.Context, exe, serverDir string, validate bool) error {
	if err := os.MkdirAll(serverDir, 0o750); err != nil {
		return problems.Wrap(problems.ServerDirInvalid, err)
	}
	args := steamcmd.Args(serverDir, server.AppID, validate)
	// "Missing configuration" e falhas passageiras costumam passar na segunda tentativa.
	var last error
	for attempt := 1; attempt <= 2; attempt++ {
		last = s.runSteamCMD(ctx, exe, serverDir, args)
		var pe *problems.Error
		if last == nil || ctx.Err() != nil || (errors.As(last, &pe) && pe.Code == problems.FoldersNested) {
			break
		}
		s.logf("SteamCMD failed, retrying (" + strconv.Itoa(attempt) + ")")
	}
	if ctx.Err() != nil {
		return problems.New(problems.Canceled, "")
	}
	return last
}

func (s *Service) runSteamCMD(ctx context.Context, exe, serverDir string, args []string) error {
	cmd := exec.CommandContext(ctx, exe, args...) // #nosec G204 -- steamcmd.exe assinado pela Valve; argumentos do domínio
	cmd.Dir = filepath.Dir(exe)
	cmd.SysProcAttr = proc.HideConsole()
	out, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	cmd.Stderr = cmd.Stdout
	if err := cmd.Start(); err != nil {
		return problems.Wrap(problems.SteamCMDFailed, err)
	}
	// O SteamCMD nem sempre escreve o progresso quando a saída vai para um pipe: o appmanifest (atualizado
	// durante o download) serve de segunda fonte.
	pollDone := make(chan struct{})
	go s.pollManifest(ctx, serverDir, pollDone)
	defer close(pollDone)

	var failure string
	success, wrongDir := false, false
	sc := bufio.NewScanner(out)
	sc.Buffer(make([]byte, 64<<10), 1<<20)
	sc.Split(splitCRLF)
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}
		s.log.Add("steamcmd", line, "")
		if steamcmd.WrongInstallDir(line) {
			wrongDir = true
		}
		l, ok := steamcmd.Parse(line)
		if !ok {
			continue
		}
		switch l.Phase {
		case steamcmd.Done:
			success = true
			s.set(func(st *State) { st.UpToDate = l.UpToDate })
		case steamcmd.Failed:
			failure = l.Error
		case steamcmd.Downloading, steamcmd.Verifying, steamcmd.Preallocate, steamcmd.Committing:
			s.progress(l.Phase, l.Done, l.Total, l.Percent)
		default:
			s.set(func(st *State) { st.Phase = l.Phase })
		}
	}
	werr := cmd.Wait()
	if ctx.Err() != nil {
		// Cancelado: o SteamCMD pode ter aberto uma cópia atualizada de si mesmo; encerra o que sobrou.
		if list, _ := proc.FindInDir(filepath.Dir(exe), "steamcmd"); len(list) > 0 {
			for _, p := range list {
				_ = proc.Kill(p.PID)
			}
		}
		return ctx.Err()
	}
	// O SteamCMD recusou a pasta e instalou em outro lugar: o "Success!" não vale para a pasta do servidor.
	if wrongDir {
		return problems.New(problems.FoldersNested, "SteamCMD ignored the install folder")
	}
	if success {
		return nil
	}
	if failure != "" {
		code := problems.SteamCMDFailed
		if bytes.Contains([]byte(failure), []byte("Disk space")) || bytes.Contains([]byte(failure), []byte("disk space")) {
			code = problems.DiskFull
		}
		return problems.New(code, failure)
	}
	if werr != nil {
		return problems.Wrap(problems.SteamCMDFailed, werr)
	}
	// Saída sem "Success!" (versões que não escrevem no pipe): confia no executável instalado.
	if fsx.Exists(filepath.Join(serverDir, "InsurgencyServer.exe")) {
		return nil
	}
	return problems.New(problems.SteamCMDFailed, "no success message")
}

func (s *Service) pollManifest(ctx context.Context, serverDir string, done <-chan struct{}) {
	t := time.NewTicker(time.Second)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-done:
			return
		case <-t.C:
			m, err := fsx.ReadLimited(manifestPath(serverDir), 1<<20)
			if err != nil {
				continue
			}
			total, _ := strconv.ParseInt(acf(m, "BytesToDownload"), 10, 64)
			got, _ := strconv.ParseInt(acf(m, "BytesDownloaded"), 10, 64)
			if total > 0 {
				st := s.State()
				// Só usa o manifesto se a saída do SteamCMD não trouxe números mais novos.
				if st.Total == 0 || st.Total == total {
					s.progress(steamcmd.Downloading, got, total, 0)
				}
			}
		}
	}
}

// splitCRLF separa linhas por \n ou \r (o SteamCMD reescreve a mesma linha com \r).
func splitCRLF(data []byte, atEOF bool) (int, []byte, error) {
	for i, b := range data {
		if b == '\n' || b == '\r' {
			return i + 1, bytes.TrimSpace(data[:i]), nil
		}
	}
	if atEOF && len(data) > 0 {
		return len(data), bytes.TrimSpace(data), nil
	}
	return 0, nil, nil
}

func (s *Service) logf(text string) {
	if s.log != nil {
		s.log.Add("app", text, "info")
	}
}
