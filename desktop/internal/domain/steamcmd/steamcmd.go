// Package steamcmd lê a saída do SteamCMD: progresso da instalação, sucesso e erros (ADR-007).
package steamcmd

import (
	"regexp"
	"strconv"
	"strings"
)

// Phase é a fase mostrada na tela.
type Phase string

// Fases conhecidas.
const (
	SelfUpdate    Phase = "self-update"
	LoggingIn     Phase = "login"
	Preallocate   Phase = "preallocating"
	Downloading   Phase = "downloading"
	Verifying     Phase = "verifying"
	Committing    Phase = "committing"
	Reconfiguring Phase = "reconfiguring"
	Done          Phase = "done"
	Failed        Phase = "failed"
)

// Line é o significado de uma linha.
type Line struct {
	Phase   Phase   `json:"phase"`
	Percent float64 `json:"percent"`
	Done    int64   `json:"done"`
	Total   int64   `json:"total"`
	// Error é o motivo do erro (ex.: "Missing configuration").
	Error string `json:"error,omitempty"`
	// UpToDate: o SteamCMD terminou dizendo que o app já estava na versão mais recente (nada baixado).
	UpToDate bool `json:"upToDate,omitempty"`
}

var (
	reState   = regexp.MustCompile(`Update state \(0x[0-9a-fA-F]+\) ([a-z ]+), progress: ([\d.]+) \((\d+) / (\d+)\)`)
	reSelfUpd = regexp.MustCompile(`^\[\s*(\d+)%\]|Downloading update|Checking for available updates|Extracting package|Installing update`)
	reSuccess = regexp.MustCompile(`Success! App '\d+' (fully installed|already up to date)`)
	reError   = regexp.MustCompile(`ERROR! (?:Failed to install app|Download item|Timed out)[^(]*\(([^)]+)\)|ERROR! (.+)`)
	reLogin   = regexp.MustCompile(`Logging in user|Connecting anonymously|Waiting for user info`)
)

// Parse lê uma linha do SteamCMD.
func Parse(line string) (Line, bool) {
	line = strings.TrimSpace(line)
	if line == "" {
		return Line{}, false
	}
	if m := reState.FindStringSubmatch(line); m != nil {
		pct, _ := strconv.ParseFloat(m[2], 64)
		done, _ := strconv.ParseInt(m[3], 10, 64)
		total, _ := strconv.ParseInt(m[4], 10, 64)
		return Line{Phase: statePhase(m[1]), Percent: pct, Done: done, Total: total}, true
	}
	if m := reSuccess.FindStringSubmatch(line); m != nil {
		return Line{Phase: Done, Percent: 100, UpToDate: m[1] == "already up to date"}, true
	}
	if m := reError.FindStringSubmatch(line); m != nil {
		msg := m[1]
		if msg == "" {
			msg = m[2]
		}
		return Line{Phase: Failed, Error: strings.TrimSpace(msg)}, true
	}
	if reLogin.MatchString(line) {
		return Line{Phase: LoggingIn}, true
	}
	if m := reSelfUpd.FindStringSubmatch(line); m != nil {
		l := Line{Phase: SelfUpdate}
		if len(m) > 1 && m[1] != "" {
			l.Percent, _ = strconv.ParseFloat(m[1], 64)
		}
		return l, true
	}
	return Line{}, false
}

func statePhase(s string) Phase {
	switch {
	case strings.Contains(s, "downloading"):
		return Downloading
	case strings.Contains(s, "verifying"):
		return Verifying
	case strings.Contains(s, "preallocating"):
		return Preallocate
	case strings.Contains(s, "committing"):
		return Committing
	case strings.Contains(s, "reconfiguring"):
		return Reconfiguring
	}
	return Downloading
}

// Args monta os argumentos do SteamCMD para instalar ou atualizar o servidor numa pasta. validate=true confere
// todos os arquivos (mais lento, conserta instalação corrompida).
func Args(installDir string, appID int, validate bool) []string {
	args := []string{"+force_install_dir", installDir, "+login", "anonymous", "+app_update", strconv.Itoa(appID)}
	if validate {
		args = append(args, "validate")
	}
	return append(args, "+quit")
}
