// Package mapchange troca o mapa com o servidor ligado, pelo RCON (ADR-019): pede a lista de cenários ao servidor,
// junta com os do perfil e manda "travelscenario" para o cenário escolhido na lista.
package mapchange

import (
	"errors"
	"strings"
	"sync"

	"sandstorm-server-assist/desktop/internal/domain/problems"
	"sandstorm-server-assist/desktop/internal/domain/server"
)

// Conn é uma conexão de RCON autenticada.
type Conn interface {
	Exec(cmd string) (string, error)
	Close() error
}

// ErrAuth é o erro que o Dial devolve quando o servidor recusa a senha.
var ErrAuth = errors.New("rcon: password refused")

// Deps são as partes do sistema que o serviço usa (dublês nos testes).
type Deps struct {
	Config func() (server.Profile, server.Secrets, error)
	// Ready diz se o servidor está ligado e com o mapa carregado.
	Ready func() bool
	// Dial conecta no RCON do servidor desta máquina. Senha errada → um erro que satisfaz errors.Is(err, ErrAuth).
	Dial func(port int, password string) (Conn, error)
	// ModScenarios devolve os cenários dos mods carregados nesta partida (ADR-029). nil = não se sabe (usa os
	// cenários do perfil). Pode ser nil.
	ModScenarios func() []string
	// Log escreve no log do app (o texto passa pela máscara de segredos).
	Log func(text, level string)
}

// Service guarda a última lista para conferir a escolha.
type Service struct {
	d    Deps
	mu   sync.Mutex
	busy bool
	list []server.Scenario
}

// New cria o serviço.
func New(d Deps) *Service { return &Service{d: d} }

// Info diz se dá para trocar o mapa pelo app.
type Info struct {
	// Rcon: RCON ligado no perfil e com senha no cofre.
	Rcon bool `json:"rcon"`
}

// Info devolve o que a tela precisa para mostrar o cartão.
func (s *Service) Info() Info {
	p, sec, err := s.d.Config()
	return Info{Rcon: err == nil && p.Rcon && sec.RconPassword != ""}
}

// connect confere as condições e abre a conexão.
func (s *Service) connect() (Conn, server.Profile, error) {
	p, sec, err := s.d.Config()
	if err != nil {
		return nil, p, err
	}
	if !p.Rcon || sec.RconPassword == "" {
		return nil, p, problems.New(problems.RconOff, "")
	}
	if !s.d.Ready() {
		return nil, p, problems.New(problems.ServerNotRunning, "")
	}
	c, err := s.d.Dial(p.RconPort, sec.RconPassword)
	if errors.Is(err, ErrAuth) {
		return nil, p, problems.New(problems.RconAuth, "")
	}
	if err != nil {
		return nil, p, problems.New(problems.RconFailed, err.Error())
	}
	return c, p, nil
}

// lock impede duas ações ao mesmo tempo (cada uma abre uma conexão no servidor).
func (s *Service) lock() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.busy {
		return problems.New(problems.Busy, "")
	}
	s.busy = true
	return nil
}

func (s *Service) unlock() {
	s.mu.Lock()
	s.busy = false
	s.mu.Unlock()
}

// Scenarios pede a lista ao servidor (uma conexão) e junta com os cenários do perfil.
func (s *Service) Scenarios() ([]server.Scenario, error) {
	if err := s.lock(); err != nil {
		return nil, err
	}
	defer s.unlock()
	c, p, err := s.connect()
	if err != nil {
		return nil, err
	}
	defer func() { _ = c.Close() }()
	resp, err := c.Exec("scenarios")
	if err != nil {
		return nil, problems.New(problems.RconFailed, err.Error())
	}
	var mods []string
	if s.d.ModScenarios != nil {
		mods = s.d.ModScenarios()
	}
	list := server.MapChoices(server.ParseScenarios(resp), p, mods)
	s.mu.Lock()
	s.list = list
	s.mu.Unlock()
	return list, nil
}

// Change troca para um cenário da última lista. A lista precisa ter sido pedida antes (a tela só mostra o botão
// com a lista aberta).
func (s *Service) Change(id, lighting string) error {
	s.mu.Lock()
	list := s.list
	s.mu.Unlock()
	p, sec, err := s.d.Config()
	if err != nil {
		return err
	}
	cmd, err := server.TravelWithLighting(id, lighting, list, p.MaxPlayers, sec.ServerPassword)
	if err != nil {
		return err
	}
	if err := s.lock(); err != nil {
		return err
	}
	defer s.unlock()
	c, _, err := s.connect()
	if err != nil {
		return err
	}
	defer func() { _ = c.Close() }()
	s.d.Log("Changing map: "+cmd, "")
	resp, err := c.Exec(cmd)
	resp = strings.TrimSpace(resp)
	if err != nil {
		return problems.New(problems.RconFailed, err.Error())
	}
	if server.TravelRefused(resp) {
		s.d.Log("The server refused the map change: "+resp, "warning")
		return problems.New(problems.MapChangeFailed, resp)
	}
	if resp != "" {
		s.d.Log("Server: "+resp, "")
	}
	return nil
}
