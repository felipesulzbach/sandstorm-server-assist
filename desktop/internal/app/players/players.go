// Package players é o painel de jogadores pelo RCON (ADR-022): lista quem está no servidor, expulsa, bane, mostra
// e tira banimentos e manda mensagens a todos. Uma conexão por clique (o RCON do Sandstorm deixa uma thread presa
// no servidor a cada conexão nova; a tela nunca consulta em intervalos).
package players

import (
	"errors"
	"strings"
	"sync"

	"sandstorm-server-assist/desktop/internal/domain/admin"
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
	Dial  func(port int, password string) (Conn, error)
	Log   func(text, level string)
}

// Service guarda a última lista para conferir as escolhas da tela.
type Service struct {
	d       Deps
	mu      sync.Mutex
	busy    bool
	players []admin.Player
	bans    []admin.Ban
}

// New cria o serviço.
func New(d Deps) *Service { return &Service{d: d} }

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

// exec manda um comando (uma conexão).
func (s *Service) exec(cmd string) (string, error) {
	out, err := s.execAll(cmd)
	if err != nil {
		return "", err
	}
	return out[0], nil
}

// execAll abre a conexão, manda os comandos em ordem e fecha. Devolve as respostas.
func (s *Service) execAll(cmds ...string) ([]string, error) {
	if err := s.lock(); err != nil {
		return nil, err
	}
	defer s.unlock()
	p, sec, err := s.d.Config()
	if err != nil {
		return nil, err
	}
	if !p.Rcon || sec.RconPassword == "" {
		return nil, problems.New(problems.RconOff, "")
	}
	if !s.d.Ready() {
		return nil, problems.New(problems.ServerNotRunning, "")
	}
	c, err := s.d.Dial(p.RconPort, sec.RconPassword)
	if errors.Is(err, ErrAuth) {
		return nil, problems.New(problems.RconAuth, "")
	}
	if err != nil {
		return nil, problems.New(problems.RconFailed, err.Error())
	}
	defer func() { _ = c.Close() }()
	out := make([]string, 0, len(cmds))
	for _, cmd := range cmds {
		resp, err := c.Exec(cmd)
		if err != nil {
			return nil, problems.New(problems.RconFailed, err.Error())
		}
		out = append(out, strings.TrimSpace(resp))
	}
	return out, nil
}

// View é o que o painel mostra.
type View struct {
	Players []admin.Player `json:"players"`
	Bans    []admin.Ban    `json:"bans"`
}

// Refresh pede a lista de jogadores e de banidos (uma conexão).
func (s *Service) Refresh() (View, error) {
	out, err := s.execAll("listplayers", "listbans")
	if err != nil {
		return View{}, err
	}
	// ParsePlayers e ParseBans devolvem listas vazias (não nil): servidor sem ninguém vira [] na tela.
	v := View{Players: admin.ParsePlayers(out[0]), Bans: admin.ParseBans(out[1])}
	s.mu.Lock()
	s.players, s.bans = v.Players, v.Bans
	s.mu.Unlock()
	return v, nil
}

// known diz se o net ID está na última lista (jogadores ou banidos).
func (s *Service) known(netID string, banned bool) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if banned {
		for _, b := range s.bans {
			if b.NetID == netID {
				return netID, true
			}
		}
		return "", false
	}
	for _, p := range s.players {
		if p.NetID == netID && !p.Bot {
			return p.Name, true
		}
	}
	return "", false
}

// Kick expulsa um jogador da última lista.
func (s *Service) Kick(netID, reason string) error {
	name, ok := s.known(netID, false)
	if !ok {
		return problems.Field("netId", "validation.netId")
	}
	cmd, err := admin.KickCommand(netID, reason)
	if err != nil {
		return err
	}
	if _, err := s.exec(cmd); err != nil {
		return err
	}
	s.d.Log("Player kicked: "+name, "info")
	return nil
}

// Ban bane um jogador da última lista (minutes 0 = permanente).
func (s *Service) Ban(netID string, minutes int, reason string) error {
	name, ok := s.known(netID, false)
	if !ok {
		return problems.Field("netId", "validation.netId")
	}
	cmd, err := admin.BanCommand(netID, minutes, reason)
	if err != nil {
		return err
	}
	if _, err := s.exec(cmd); err != nil {
		return err
	}
	s.d.Log("Player banned: "+name, "info")
	return nil
}

// Unban tira o banimento de um net ID da última lista de banidos.
func (s *Service) Unban(netID string) error {
	if _, ok := s.known(netID, true); !ok {
		return problems.Field("netId", "validation.netId")
	}
	cmd, err := admin.UnbanCommand(netID)
	if err != nil {
		return err
	}
	if _, err := s.exec(cmd); err != nil {
		return err
	}
	s.d.Log("Ban removed", "info")
	return nil
}

// Say manda uma mensagem a todos no chat.
func (s *Service) Say(text string) error {
	cmd, err := server.SayCommand(text)
	if err != nil {
		return err
	}
	_, err = s.exec(cmd)
	return err
}
