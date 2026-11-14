package desktop

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"sandstorm-server-assist/desktop/internal/adapters/fsx"
	"sandstorm-server-assist/desktop/internal/adapters/modio"
	"sandstorm-server-assist/desktop/internal/adapters/netdiag"
	"sandstorm-server-assist/desktop/internal/adapters/winsys"
	"sandstorm-server-assist/desktop/internal/app/config"
	"sandstorm-server-assist/desktop/internal/app/mapchange"
	"sandstorm-server-assist/desktop/internal/app/runner"
	"sandstorm-server-assist/desktop/internal/domain/netcheck"
	"sandstorm-server-assist/desktop/internal/domain/problems"
	"sandstorm-server-assist/desktop/internal/domain/server"
)

// ServerAPI: iniciar, parar, autorizar no mod.io e lista de mods (evento "server:status").
type ServerAPI struct {
	a  *app
	mu sync.Mutex
	// subscribed é a última lista de inscrições lida do mod.io.
	subscribed   []modio.Mod
	subscribedAt time.Time
}

// Status devolve o estado do servidor.
func (s *ServerAPI) Status() runner.Status { return s.a.runner.Status() }

// Start inicia o servidor (modo normal, -SecurityCode=none).
func (s *ServerAPI) Start() error {
	if s.a.install.Busy() {
		return problems.New(problems.Busy, "")
	}
	return s.a.runner.Start()
}

// Stop encerra o servidor.
func (s *ServerAPI) Stop() error { return s.a.runner.Stop() }

// Restart para e inicia de novo.
func (s *ServerAPI) Restart() error {
	if s.a.install.Busy() {
		return problems.New(problems.Busy, "")
	}
	return s.a.runner.Restart()
}

// Connect é o que os jogadores precisam para entrar.
type Connect struct {
	IP      string `json:"ip"`
	Port    int    `json:"port"`
	Command string `json:"command"`
}

// Connect devolve o IP configurado e o comando do console do jogo.
func (s *ServerAPI) Connect() Connect {
	p := s.a.cfg.Profile()
	c := Connect{IP: p.ConnectIP, Port: p.Port}
	if p.ConnectIP != "" {
		c.Command = fmt.Sprintf("open %s:%d", p.ConnectIP, p.Port)
	}
	return c
}

// ---- troca de mapa (RCON, ADR-019) ----

// MapChange diz se dá para trocar o mapa pelo app (RCON ligado e com senha).
func (s *ServerAPI) MapChange() mapchange.Info { return s.a.maps.Info() }

// Scenarios pede ao servidor ligado a lista de mapas (inclui os mapas dos mods montados) e junta com os do perfil.
func (s *ServerAPI) Scenarios() ([]server.Scenario, error) { return s.a.maps.Scenarios() }

// ChangeMap troca para um cenário da lista (a tela manda só o id e, se quiser, Day ou Night).
func (s *ServerAPI) ChangeMap(id, lighting string) error { return s.a.maps.Change(id, lighting) }

// ---- diagnóstico de portas feito no PC (ADR-033) ----

// CheckPorts faz o diagnóstico das portas pelo próprio PC: servidor respondendo, Firewall, roteador (UPnP), IP
// público e CGNAT (STUN), retorno pelo IP público e lista da Steam. Funciona com o servidor parado (mede menos).
func (s *ServerAPI) CheckPorts() netcheck.Result {
	p := s.a.cfg.Profile()
	fw := netcheck.Unknown
	if info := (&SystemAPI{a: s.a}).Firewall(); info.Ready {
		fw = netcheck.No
		if info.Rules {
			fw = netcheck.Yes
		}
	}
	var local []string
	for _, a := range winsys.LocalAddresses() {
		local = append(local, a.IP)
	}
	return netdiag.Run(s.a.ctx, netdiag.Input{
		GamePort: p.Port, QueryPort: p.QueryPort, RconPort: p.RconPort, Rcon: p.Rcon,
		ServerRunning: s.a.runner.Status().Phase == runner.Ready,
		LocalIPs:      local, Firewall: fw,
	})
}

// ---- mod.io ----

// ModsView é a tela de mods.
type ModsView struct {
	Session   modio.Session `json:"session"`
	HasAPIKey bool          `json:"hasApiKey"`
	// GameKey: a chave do mod.io do próprio jogo foi achada nos arquivos do servidor; o pedido do código usa ela e o
	// usuário não precisa de chave própria (ADR-044).
	GameKey      bool      `json:"gameKey"`
	HasEmail     bool      `json:"hasEmail"`
	Root         string    `json:"root"`
	Mods         []ModRow  `json:"mods"`
	SubscribedAt time.Time `json:"subscribedAt"`
}

// ModRow junta o que está instalado, o que está inscrito e o que acontece agora no servidor.
type ModRow struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	ProfileURL string `json:"-"`
	Size       int64  `json:"size"`
	OnDisk     int64  `json:"onDisk"`
	Installed  bool   `json:"installed"`
	Subscribed bool   `json:"subscribed"`
	// Live é o estado nesta partida (active, installing...), vazio se o servidor está parado.
	Live string `json:"live"`
	Done int64  `json:"done"`
}

// Mods monta a lista de mods.
func (s *ServerAPI) Mods() ModsView {
	root := modio.RootDir()
	v := ModsView{Session: modioSession(), Root: root}
	sec, _ := s.a.cfg.SecretsState()
	v.HasAPIKey, v.HasEmail = sec["modioApiKey"], sec["modioEmail"]
	v.GameKey = modio.GameAPIKey(s.a.cfg.Profile().ServerDir) != ""
	rows := map[int]*ModRow{}
	get := func(id int) *ModRow {
		if r, ok := rows[id]; ok {
			return r
		}
		r := &ModRow{ID: id}
		rows[id] = r
		return r
	}
	if inst, err := modio.ServerInstalled(root); err == nil {
		for _, m := range inst {
			r := get(m.ID)
			r.Name, r.Size, r.Installed, r.OnDisk = m.Name, m.Size, true, m.SizeOnDisk
		}
	}
	s.mu.Lock()
	v.SubscribedAt = s.subscribedAt
	for _, m := range s.subscribed {
		r := get(m.ID)
		r.Subscribed = true
		if r.Name == "" {
			r.Name = m.Name
		}
		if r.Size == 0 {
			r.Size = m.Size
		}
	}
	s.mu.Unlock()
	for _, m := range s.a.runner.Status().Mods {
		r := get(m.ID)
		r.Live, r.Done = m.State, m.Done
		if r.Name == "" {
			r.Name = m.Name
		}
		if m.Total > 0 {
			r.Size = m.Total
		}
	}
	for _, r := range rows {
		// Desinstalado pelo servidor nesta partida e sem arquivos: some da lista.
		if r.Live == "removed" && !r.Installed {
			continue
		}
		if r.Name == "" {
			r.Name = fmt.Sprintf("#%d", r.ID)
		}
		v.Mods = append(v.Mods, *r)
	}
	sort.Slice(v.Mods, func(i, j int) bool { return strings.ToLower(v.Mods[i].Name) < strings.ToLower(v.Mods[j].Name) })
	return v
}

// DeleteMod apaga um mod baixado (pasta e registro do mod.io) para o servidor baixar de novo (ADR-037).
// redownload=true inicia o servidor logo depois; com o servidor rodando, ele para, apaga e volta.
func (s *ServerAPI) DeleteMod(id int, redownload bool) error {
	return s.a.modFiles.Delete(id, redownload)
}

// DownloadNewMods faz o servidor baixar os mods inscritos que ainda não estão no disco: inicia o servidor, ou
// reinicia se ele está rodando (ADR-038).
func (s *ServerAPI) DownloadNewMods() error { return s.a.modFiles.DownloadNew() }

// ClearModio apaga, para começar do zero, os mods baixados (mods) e a autorização do servidor no mod.io (session).
// Servidor parado (ADR-047).
func (s *ServerAPI) ClearModio(mods, session bool) error { return s.a.modFiles.Clear(mods, session) }

// RefreshSubscriptions lê do mod.io as inscrições da conta do servidor (usa a sessão do user.json).
func (s *ServerAPI) RefreshSubscriptions() (ModsView, error) {
	uf, err := modio.UserFile()
	if err != nil {
		return s.Mods(), problems.Wrap(problems.ModioNotAuthorized, err)
	}
	sess, token, _ := modio.ReadSession(uf)
	if !sess.Valid || token == "" {
		return s.Mods(), problems.New(problems.ModioNotAuthorized, "")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	list, err := modio.New().Subscribed(ctx, token)
	if err != nil {
		return s.Mods(), err
	}
	s.mu.Lock()
	s.subscribed, s.subscribedAt = list, time.Now()
	s.mu.Unlock()
	return s.Mods(), nil
}

// RequestCode pede o código de autorização para o e-mail guardado no cofre.
func (s *ServerAPI) RequestCode() error {
	email, err := s.a.cfg.Secret("modioEmail")
	if err != nil {
		return err
	}
	if email == "" {
		return problems.Field("modioEmail", "validation.email")
	}
	// O servidor troca o código com a chave do próprio jogo, e o mod.io só aceita a troca com a mesma chave do pedido
	// (ADR-044). A chave guardada no cofre só serve se a do jogo não for achada.
	key := modio.GameAPIKey(s.a.cfg.Profile().ServerDir)
	if key == "" {
		key, _ = s.a.cfg.Secret("modioApiKey")
	}
	if key == "" {
		return problems.New(problems.ModioNoAPIKey, "")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := modio.New().RequestCode(ctx, key, email); err != nil {
		return err
	}
	s.a.log.Add("app", "mod.io: security code requested by e-mail", "info")
	return nil
}

// Authorize inicia o servidor uma vez com o código recebido; depois o app volta sozinho para "none".
func (s *ServerAPI) Authorize(code string) error {
	if s.a.install.Busy() {
		return problems.New(problems.Busy, "")
	}
	return s.a.runner.Authorize(code)
}

// DetectAPIKey procura a chave de API do mod.io nos logs do próprio servidor e guarda no cofre (ADR-010).
func (s *ServerAPI) DetectAPIKey() (bool, error) {
	p := s.a.cfg.Profile()
	if p.ServerDir == "" {
		return false, problems.New(problems.ServerNotInstalled, "")
	}
	key := findAPIKey(p.ServerDir)
	if key == "" {
		return false, nil
	}
	if err := s.a.cfg.SetSecret("modioApiKey", key); err != nil {
		return false, err
	}
	return true, nil
}

// findAPIKey procura a chave de API do mod.io nos logs do servidor, do mais novo para o mais velho.
func findAPIKey(serverDir string) string {
	if serverDir == "" {
		return ""
	}
	for _, path := range config.LogFiles(serverDir) {
		b, err := fsx.ReadLimited(path, 64<<20)
		if err != nil {
			continue
		}
		if key := modio.FindAPIKey(string(b)); key != "" {
			return key
		}
	}
	return ""
}
