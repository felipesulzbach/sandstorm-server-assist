// Package config guarda o perfil do servidor e as preferências do app (JSON) e os segredos (cofre DPAPI), importa
// uma instalação já existente e grava os arquivos de configuração do servidor antes de iniciar.
package config

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"sandstorm-server-assist/desktop/internal/adapters/fsx"
	"sandstorm-server-assist/desktop/internal/domain/ini"
	"sandstorm-server-assist/desktop/internal/domain/legal"
	"sandstorm-server-assist/desktop/internal/domain/problems"
	"sandstorm-server-assist/desktop/internal/domain/server"
	"sandstorm-server-assist/desktop/internal/domain/serverlogs"
)

// Languages são os idiomas da tela (ADR-009). O primeiro é o padrão.
var Languages = []string{"en-US", "pt-PT", "es-ES", "ru", "uk"}

// Settings são as preferências do app.
type Settings struct {
	Language string `json:"language"`
	// ShowConsole abre a janela de console do servidor (padrão: escondida, os logs aparecem no app).
	ShowConsole bool `json:"showConsole"`
	// TermsAccepted é a versão dos termos de uso e do aviso de dados que o usuário aceitou (ADR-015).
	TermsAccepted string `json:"termsAccepted,omitempty"`
	// TermsAcceptedAt é quando aceitou (RFC 3339, hora local).
	TermsAcceptedAt string `json:"termsAcceptedAt,omitempty"`
	// LogCleanup apaga os logs antigos do servidor a cada início (ADR-018). Padrão: desligado, não apaga nada.
	LogCleanup bool `json:"logCleanup"`
	// LogKeep é quantas partidas manter (serverlogs.MinKeep a MaxKeep; 0 = padrão).
	LogKeep int `json:"logKeep,omitempty"`
	// StartWithWindows abre o app ao entrar no Windows (minimizado). AutoStartServer inicia o servidor quando o app
	// abre assim (ADR-022).
	StartWithWindows bool `json:"startWithWindows"`
	AutoStartServer  bool `json:"autoStartServer"`
}

type file struct {
	Version int            `json:"version"`
	Profile server.Profile `json:"profile"`
	// ActiveID e ActiveName são do perfil em uso (Profile). Others são os outros perfis guardados (ADR-022).
	ActiveID   string         `json:"activeId,omitempty"`
	ActiveName string         `json:"activeName,omitempty"`
	Others     []NamedProfile `json:"others,omitempty"`
	Settings   Settings       `json:"settings"`
}

// SecretStore é o cofre (DPAPI no app, memória nos testes).
type SecretStore interface {
	Load() (server.Secrets, error)
	Save(server.Secrets) error
}

// Service guarda o estado em memória e grava a cada mudança.
type Service struct {
	mu       sync.Mutex
	path     string
	data     file
	vault    SecretStore
	secrets  server.Secrets
	vaultErr error
}

// Paths dos arquivos do servidor dentro da pasta.
const (
	GameIni     = `Insurgency\Saved\Config\WindowsServer\Game.ini`
	EngineIni   = `Insurgency\Saved\Config\WindowsServer\Engine.ini`
	AdminsTxt   = `Insurgency\Config\Server\Admins.txt`
	MapCycleTxt = `Insurgency\Config\Server\MapCycle.txt`
	MotdTxt     = `Insurgency\Config\Server\Motd.txt`
)

// Open carrega o arquivo de configuração (ausente = padrões) e o cofre.
func Open(path string, vault SecretStore) *Service {
	s := &Service{path: path, vault: vault}
	s.data = file{Version: 1, Profile: server.Default()}
	if b, err := fsx.ReadLimited(path, 1<<20); err == nil {
		var f file
		if json.Unmarshal(b, &f) == nil && f.Version > 0 {
			s.data = f
		}
	}
	if s.data.Settings.Language == "" {
		s.data.Settings.Language = ""
	}
	// Rotação com "Scenario=" repetido (colada da página do mod): corrige uma vez.
	if fixed := server.FixMapCycle(s.data.Profile.MapCycle); !slices.Equal(fixed, s.data.Profile.MapCycle) {
		s.data.Profile.MapCycle = fixed
		_ = s.save()
	}
	// Configuração antiga apontando para a pasta de cima do servidor: corrige uma vez.
	if dir := s.data.Profile.ServerDir; dir != "" {
		if fixed := ResolveServerDir(dir); fixed != dir {
			s.data.Profile.ServerDir = fixed
			_ = s.save()
		}
	}
	s.secrets, s.vaultErr = vault.Load()
	return s
}

// DefaultPath é %APPDATA%\sandstorm-server-assist\config.json.
func DefaultPath() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "sandstorm-server-assist", "config.json"), nil
}

func (s *Service) save() error {
	b, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	return fsx.WriteAtomic(s.path, append(b, '\n'), false)
}

// Profile devolve o perfil.
func (s *Service) Profile() server.Profile {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.data.Profile
}

// Current devolve perfil e segredos (para iniciar o servidor).
func (s *Service) Current() (server.Profile, server.Secrets, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.vaultErr != nil {
		return s.data.Profile, server.Secrets{}, s.vaultErr
	}
	return s.data.Profile, s.secrets, nil
}

// SaveProfile valida e grava. As pastas não mudam por aqui (só pelos diálogos: SetServerDir/SetSteamCMDDir).
func (s *Service) SaveProfile(p server.Profile) (server.Profile, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p.ServerDir, p.SteamCMDDir = s.data.Profile.ServerDir, s.data.Profile.SteamCMDDir
	p.MapCycle = server.FixMapCycle(p.MapCycle)
	if err := p.Validate(); err != nil {
		return s.data.Profile, err
	}
	p.Rules = server.CleanRules(p.Rules)
	s.data.Profile = p
	return p, s.save()
}

// SetServerDir troca a pasta do servidor (escolhida num diálogo do Windows).
func (s *Service) SetServerDir(dir string) (server.Profile, error) {
	if !fsx.SafeDir(dir) {
		return s.Profile(), problems.New(problems.ServerDirInvalid, dir)
	}
	if fixed := ResolveServerDir(filepath.Clean(dir)); fsx.SafeDir(fixed) {
		dir = fixed
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data.Profile.ServerDir = filepath.Clean(dir)
	return s.data.Profile, s.save()
}

// SetSteamCMDDir troca a pasta do SteamCMD.
func (s *Service) SetSteamCMDDir(dir string) (server.Profile, error) {
	if !fsx.SafeDir(dir) {
		return s.Profile(), problems.New(problems.ServerDirInvalid, dir)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data.Profile.SteamCMDDir = filepath.Clean(dir)
	return s.data.Profile, s.save()
}

// Settings devolve as preferências.
func (s *Service) Settings() Settings {
	s.mu.Lock()
	defer s.mu.Unlock()
	st := s.data.Settings
	if !serverlogs.ValidKeep(st.LogKeep) {
		st.LogKeep = serverlogs.DefaultKeep
	}
	return st
}

// SetLanguage grava o idioma (só os da lista).
func (s *Service) SetLanguage(lang string) (Settings, error) {
	ok := false
	for _, l := range Languages {
		ok = ok || l == lang
	}
	if !ok {
		return s.Settings(), problems.Field("language", "validation.unknownField")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data.Settings.Language = lang
	return s.data.Settings, s.save()
}

// SetShowConsole liga ou desliga a janela de console do servidor.
func (s *Service) SetShowConsole(on bool) (Settings, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data.Settings.ShowConsole = on
	return s.data.Settings, s.save()
}

// SetAutostart grava as preferências de abrir com o Windows (quem mexe no registro é a ponte, ADR-022).
func (s *Service) SetAutostart(withWindows, startServer bool) (Settings, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data.Settings.StartWithWindows = withWindows
	s.data.Settings.AutoStartServer = withWindows && startServer
	return s.data.Settings, s.save()
}

// AcceptTerms grava o aceite dos termos de uso e do aviso de dados. Só aceita a versão em vigor.
func (s *Service) AcceptTerms(version string, now time.Time) (Settings, error) {
	if !legal.Accepted(version) {
		return s.Settings(), problems.Field("terms", "validation.unknownField")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data.Settings.TermsAccepted = version
	s.data.Settings.TermsAcceptedAt = now.Format(time.RFC3339)
	return s.data.Settings, s.save()
}

// SecretsState diz quais segredos estão definidos (a tela nunca recebe o valor sem pedir).
func (s *Service) SecretsState() (map[string]bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := map[string]bool{}
	for _, k := range server.SecretKeys {
		v, _ := s.secrets.Get(k)
		out[k] = v != ""
	}
	return out, s.vaultErr
}

// Secret devolve um segredo (para "Mostrar" na tela ou uso interno).
func (s *Service) Secret(key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.vaultErr != nil {
		return "", s.vaultErr
	}
	v, ok := s.secrets.Get(key)
	if !ok {
		return "", problems.Field(key, "validation.unknownField")
	}
	return v, nil
}

// SetSecret valida e grava um segredo no cofre (vazio apaga).
func (s *Service) SetSecret(key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	next, err := s.secrets.With(key, value)
	if err != nil {
		return err
	}
	if err := s.vault.Save(next); err != nil {
		return err
	}
	s.secrets, s.vaultErr = next, nil
	return nil
}

// ReplaceAll troca perfil e segredos de uma vez (restaurar backup). A pasta do servidor é a da máquina nova.
func (s *Service) ReplaceAll(p server.Profile, sec server.Secrets, serverDir string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	p.ServerDir = serverDir
	if p.SteamCMDDir != "" && !fsx.Exists(p.SteamCMDDir) {
		p.SteamCMDDir = ""
	}
	if err := p.Validate(); err != nil {
		return err
	}
	if err := s.vault.Save(sec); err != nil {
		return err
	}
	s.data.Profile, s.secrets, s.vaultErr = p, sec, nil
	return s.save()
}

// SecretValues devolve os valores para mascarar em log.
func (s *Service) SecretValues() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.secrets.Values()
}

// ---- servidor existente ----

// LogsDir é a pasta de logs do servidor (Insurgency.log e os backups Insurgency-backup-*.log).
const LogsDir = `Insurgency\Saved\Logs`

// LogFiles devolve os logs do servidor (Insurgency.log, Insurgency_2.log quando o primeiro estava aberto, e os
// Insurgency-backup-*.log), do mais recente para o mais velho pela data de gravação.
func LogFiles(serverDir string) []string {
	dir := filepath.Join(serverDir, LogsDir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	type logFile struct {
		path string
		mod  int64
	}
	var logs []logFile
	for _, e := range entries {
		n := e.Name()
		if e.IsDir() || !strings.HasPrefix(n, "Insurgency") || !strings.HasSuffix(n, ".log") {
			continue
		}
		if info, err := e.Info(); err == nil {
			logs = append(logs, logFile{filepath.Join(dir, n), info.ModTime().UnixNano()})
		}
	}
	// Mesmo horário: o Insurgency.log (log atual) vem antes; depois, nome maior (backup mais novo) antes.
	sort.SliceStable(logs, func(i, j int) bool {
		if logs[i].mod != logs[j].mod {
			return logs[i].mod > logs[j].mod
		}
		if a, b := filepath.Base(logs[i].path) == "Insurgency.log", filepath.Base(logs[j].path) == "Insurgency.log"; a != b {
			return a
		}
		return logs[i].path > logs[j].path
	})
	out := make([]string, 0, len(logs))
	for _, l := range logs {
		out = append(out, l.path)
	}
	return out
}

// ServerExe marca a pasta de um servidor instalado.
const ServerExe = "InsurgencyServer.exe"

// ResolveServerDir corrige a escolha da pasta de cima: se a pasta não tem o servidor, mas exatamente uma
// subpasta direta tem, devolve a subpasta. Senão devolve a pasta como veio.
func ResolveServerDir(dir string) string {
	if dir == "" || fsx.Exists(filepath.Join(dir, ServerExe)) {
		return dir
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return dir
	}
	found := ""
	for _, e := range entries {
		if e.IsDir() && fsx.Exists(filepath.Join(dir, e.Name(), ServerExe)) {
			if found != "" {
				return dir // mais de um servidor: o usuário escolhe
			}
			found = filepath.Join(dir, e.Name())
		}
	}
	if found == "" {
		return dir
	}
	return found
}

// lastCommandLine acha a linha de comando da última partida ("LogInit: Command Line:"), que fica no começo do
// log. Serve com o servidor rodando (log atual) ou parado (último log). Devolve também o nome do log.
func lastCommandLine(serverDir string) (string, string) {
	const marker = "Command Line:"
	for _, path := range LogFiles(serverDir) {
		b, err := fsx.ReadHead(path, 512<<10)
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(b), "\n") {
			if _, after, ok := strings.Cut(line, "LogInit: "+marker); ok {
				return strings.TrimSpace(after), filepath.Base(path)
			}
		}
	}
	return "", ""
}

// ServerRules lê as regras do jogo que o servidor já tem no Game.ini e no Engine.ini (ADR-040). nil quando não há
// nenhuma, ou sem pasta do servidor. Só lê: a aba Regras oferece usá-las quando o app ainda não cuida das regras.
func (s *Service) ServerRules() map[string]string {
	dir := s.Profile().ServerDir
	if dir == "" {
		return nil
	}
	game, _ := fsx.ReadLimited(filepath.Join(dir, GameIni), 1<<20)
	engine, _ := fsx.ReadLimited(filepath.Join(dir, EngineIni), 1<<20)
	return server.ReadRules(string(game), string(engine))
}

// Detect lê o servidor existente na pasta escolhida: a linha de comando da última partida (log), o Admins.txt,
// o MapCycle.txt e o [Rcon] do Game.ini. apiKey é a chave do mod.io achada nos logs (ou vazia).
func (s *Service) Detect(apiKey string) (server.Detected, string, error) {
	d := server.NewDetected()
	dir := s.Profile().ServerDir
	if dir == "" {
		return d, "", problems.New(problems.ServerNotInstalled, "")
	}
	cmd, logFile := lastCommandLine(dir)
	d.ParseCommandLine(cmd)

	if b, err := fsx.ReadLimited(filepath.Join(dir, AdminsTxt), 1<<20); err == nil {
		d.Profile.Admins, d.Found["admins"] = server.ValidAdmins(server.ParseAdmins(string(b))), true
	}
	if b, err := fsx.ReadLimited(filepath.Join(dir, MapCycleTxt), 1<<20); err == nil {
		var lines []string
		for _, l := range server.ParseMapCycle(string(b)) {
			if len(l) <= 200 {
				lines = append(lines, l)
			}
		}
		d.Profile.MapCycle, d.Found["mapCycle"] = lines, true
	}
	if b, err := fsx.ReadLimited(filepath.Join(dir, MotdTxt), 64<<10); err == nil {
		if text, ok := server.ParseMotd(string(b)); ok {
			d.Profile.Motd, d.Found["motd"] = text, true
		}
	}
	game, gerr := fsx.ReadLimited(filepath.Join(dir, GameIni), 1<<20)
	engine, _ := fsx.ReadLimited(filepath.Join(dir, EngineIni), 1<<20)
	if rules := server.ReadRules(string(game), string(engine)); rules != nil {
		d.Profile.Rules, d.Found["rules"] = rules, true
	}
	if b, err := game, gerr; err == nil {
		text := string(b)
		if on, _ := ini.Get(text, "Rcon", "bEnabled"); strings.EqualFold(on, "true") {
			d.Profile.Rcon, d.Found["rcon"] = true, true
		}
		if port, ok := ini.Get(text, "Rcon", "ListenPort"); ok && !d.Found["rconPort"] {
			if n, err := strconv.Atoi(port); err == nil && n >= 1024 && n <= 65535 {
				d.Profile.RconPort, d.Found["rconPort"] = n, true
			}
		}
		if pw, ok := ini.Get(text, "Rcon", "Password"); ok && d.Secrets.RconPassword == "" {
			if next, err := d.Secrets.With("rconPassword", pw); err == nil {
				d.Secrets = next
			}
		}
	}
	if apiKey != "" {
		if next, err := d.Secrets.With("modioApiKey", apiKey); err == nil {
			d.Secrets = next
		}
	}
	return d, logFile, nil
}

// Scan é a comparação entre o app e o servidor existente (nenhum segredo vai para a tela).
type Scan struct {
	Found     bool          `json:"found"`
	LogFile   string        `json:"logFile"`
	Diffs     []server.Diff `json:"diffs"`
	RconInIni bool          `json:"rconInIni"`
}

// ScanExisting compara o que está no app com o servidor existente, campo a campo.
func (s *Service) ScanExisting(apiKey string) (Scan, error) {
	d, logFile, err := s.Detect(apiKey)
	if err != nil {
		return Scan{}, err
	}
	p, sec, err := s.Current()
	if err != nil {
		return Scan{}, err
	}
	diffs := server.Compare(p, sec, d)
	return Scan{Found: len(diffs) > 0, LogFile: logFile, Diffs: diffs, RconInIni: s.RconInIni()}, nil
}

// ApplyResult diz o que mudou ao trazer os dados do servidor.
type ApplyResult struct {
	Profile   server.Profile `json:"profile"`
	Applied   []string       `json:"applied"`
	Kept      []string       `json:"kept"`
	Skipped   []string       `json:"skipped"`
	RconInIni bool           `json:"rconInIni"`
}

// ApplyExisting traz os dados do servidor existente. keepMine são os campos em conflito em que o usuário quer
// manter o valor do app (a tela manda só as chaves). steamcmdDir é o SteamCMD achado pelo appmanifest.
func (s *Service) ApplyExisting(steamcmdDir, apiKey string, keepMine []string) (ApplyResult, error) {
	var r ApplyResult
	known := map[string]bool{}
	for _, k := range server.SyncKeys() {
		known[k] = true
	}
	keep := map[string]bool{}
	for _, k := range keepMine {
		if !known[k] {
			return r, problems.Field(k, "validation.unknownField")
		}
		keep[k] = true
	}
	d, _, err := s.Detect(apiKey)
	if err != nil {
		return r, err
	}
	r.RconInIni = s.RconInIni() // antes do lock: RconInIni também trava

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.vaultErr != nil {
		return r, s.vaultErr
	}
	for _, df := range server.Compare(s.data.Profile, s.secrets, d) {
		if df.State == server.SyncConflict && keep[df.Key] {
			r.Kept = append(r.Kept, df.Key)
		}
	}
	p, sec, applied, skipped, err := server.Merge(s.data.Profile, s.secrets, d, keep)
	if err != nil {
		return r, err
	}
	if p.SteamCMDDir == "" && steamcmdDir != "" && fsx.SafeDir(steamcmdDir) {
		p.SteamCMDDir = filepath.Clean(steamcmdDir)
	}
	if sec != s.secrets {
		if err := s.vault.Save(sec); err != nil {
			return r, err
		}
		s.secrets = sec
	}
	s.data.Profile = p
	r.Profile, r.Applied, r.Skipped = p, applied, skipped
	return r, s.save()
}

// MoveRconToVault tira a senha de RCON do Game.ini (com cópia .bak antes). O app passa a senha pela linha de
// comando ao iniciar (ADR-002).
func (s *Service) MoveRconToVault() error {
	p := s.Profile()
	path := filepath.Join(p.ServerDir, GameIni)
	b, err := fsx.ReadLimited(path, 1<<20)
	if err != nil {
		return problems.Wrap(problems.ConfigWrite, err)
	}
	if pw, ok := ini.Get(string(b), "Rcon", "Password"); ok && pw != "" {
		if cur, _ := s.Secret("rconPassword"); cur == "" {
			if err := s.SetSecret("rconPassword", pw); err != nil {
				return err
			}
		}
	}
	out, changed := ini.Remove(string(b), "Rcon", "Password")
	if !changed {
		return nil
	}
	if err := fsx.WriteAtomic(path, []byte(out), true); err != nil {
		return problems.Wrap(problems.ConfigWrite, err)
	}
	return nil
}

// RconInIni diz se o Game.ini ainda tem senha de RCON em texto puro.
func (s *Service) RconInIni() bool {
	p := s.Profile()
	if p.ServerDir == "" {
		return false
	}
	b, err := fsx.ReadLimited(filepath.Join(p.ServerDir, GameIni), 1<<20)
	if err != nil {
		return false
	}
	pw, ok := ini.Get(string(b), "Rcon", "Password")
	return ok && pw != ""
}

// Prepare grava Admins.txt, MapCycle.txt, Motd.txt e as regras (Game.ini e Engine.ini) a partir do perfil antes de
// iniciar (cópia .bak na primeira vez). Listas e textos vazios não mexem nos arquivos.
func Prepare(p server.Profile) error {
	if err := writeRules(p); err != nil {
		return err
	}
	if len(p.Admins) > 0 {
		if err := fsx.WriteAtomic(filepath.Join(p.ServerDir, AdminsTxt), []byte(server.AdminsFile(p.Admins)), true); err != nil {
			return err
		}
	}
	if len(p.MapCycle) > 0 {
		if err := fsx.WriteAtomic(filepath.Join(p.ServerDir, MapCycleTxt), []byte(server.MapCycleFile(p.MapCycle)), true); err != nil {
			return err
		}
	}
	if motd := server.MotdFile(p.Motd); motd != "" {
		if err := fsx.WriteAtomic(filepath.Join(p.ServerDir, MotdTxt), []byte(motd), true); err != nil {
			return err
		}
	}
	return nil
}

// writeRules grava as regras do perfil no Game.ini e no Engine.ini (só os arquivos que mudam).
func writeRules(p server.Profile) error {
	if p.Rules == nil {
		return nil
	}
	gamePath, enginePath := filepath.Join(p.ServerDir, GameIni), filepath.Join(p.ServerDir, EngineIni)
	read := func(path string) (string, error) {
		b, err := fsx.ReadLimited(path, 1<<20)
		if err != nil && !IsNotExist(err) {
			return "", err
		}
		return string(b), nil
	}
	game, err := read(gamePath)
	if err != nil {
		return err
	}
	engine, err := read(enginePath)
	if err != nil {
		return err
	}
	g, e := server.ApplyRules(game, engine, p.Rules)
	if g != game {
		if err := fsx.WriteAtomic(gamePath, []byte(g), true); err != nil {
			return err
		}
	}
	if e != engine {
		if err := fsx.WriteAtomic(enginePath, []byte(e), true); err != nil {
			return err
		}
	}
	return nil
}

// IsNotExist ajuda a tela a separar "não existe" de outros erros.
func IsNotExist(err error) bool { return errors.Is(err, fs.ErrNotExist) }
