// Package chat manda as mensagens de entrada e saída de jogadores para o chat do servidor, pelo RCON (ADR-020).
//
// O jogo não tem texto próprio para isso: o app lê no log quem entrou (pedido de entrada com nome e id, depois a
// entrada confirmada) e quem saiu (conexão fechada, só com o id), e manda "say <mensagem>". Uma conexão de RCON
// fica aberta enquanto houver mensagens: o RCON do Sandstorm deixa uma thread presa a cada conexão nova.
package chat

import (
	"strings"
	"sync"
	"time"

	"sandstorm-server-assist/desktop/internal/domain/logparse"
	"sandstorm-server-assist/desktop/internal/domain/server"
)

// Conn é uma conexão de RCON autenticada.
type Conn interface {
	Exec(cmd string) (string, error)
	Close() error
}

// Deps são as partes do sistema que o serviço usa (dublês nos testes).
type Deps struct {
	Config func() (server.Profile, server.Secrets, error)
	Dial   func(port int, password string) (Conn, error)
	// Log escreve no log do app (o texto passa pela máscara de segredos).
	Log func(text, level string)
	Now func() time.Time
}

// MaxAge: linhas mais velhas que isso não geram mensagem (o app pode ler o log de uma partida já em andamento).
const MaxAge = 2 * time.Minute

// maxTracked limita quantos jogadores o serviço lembra (pedidos de entrada sem confirmação não crescem sem fim).
const maxTracked = 256

// Service guarda quem pediu para entrar e quem entrou, e a conexão aberta.
type Service struct {
	d      Deps
	mu     sync.Mutex
	logins map[string]string // id → nome do pedido de entrada
	joined map[string]string // id → nome de quem entrou

	queue  chan string
	sendMu sync.Mutex
	conn   Conn
}

// New cria o serviço. Async liga a fila de envio (o app); nos testes, Handle devolve o comando sem enviar.
func New(d Deps, async bool) *Service {
	if d.Now == nil {
		d.Now = time.Now
	}
	s := &Service{d: d, logins: map[string]string{}, joined: map[string]string{}}
	if async {
		s.queue = make(chan string, 32)
		go s.worker()
	}
	return s
}

// Observe recebe um evento do log e, se for o caso, põe a mensagem na fila (sem travar a leitura do log).
func (s *Service) Observe(ev logparse.Event) {
	cmd := s.Handle(ev)
	if cmd == "" || s.queue == nil {
		return
	}
	select {
	case s.queue <- cmd:
	default: // fila cheia: perde a mensagem, nunca trava o log
	}
}

// Handle atualiza quem está no servidor e devolve o comando "say" a mandar (vazio se nada).
func (s *Service) Handle(ev logparse.Event) string {
	s.mu.Lock()
	var template, name string
	switch ev.Kind {
	case logparse.PlayerLogin:
		if len(s.logins) >= maxTracked {
			clear(s.logins)
		}
		s.logins[ev.PlayerID] = ev.Player
	case logparse.PlayerJoined:
		id, full := s.match(ev.Player)
		name = full
		if id != "" {
			delete(s.logins, id)
			if len(s.joined) >= maxTracked {
				clear(s.joined)
			}
			s.joined[id] = full
		}
		template = "join"
	case logparse.PlayerLeft:
		n, ok := s.joined[ev.PlayerID]
		if !ok {
			s.mu.Unlock()
			return ""
		}
		delete(s.joined, ev.PlayerID)
		name, template = n, "leave"
	}
	s.mu.Unlock()

	if template == "" || (!ev.Time.IsZero() && s.d.Now().Sub(ev.Time) > MaxAge) {
		return ""
	}
	p, sec, err := s.d.Config()
	if err != nil || !p.Rcon || sec.RconPassword == "" {
		return ""
	}
	if template == "join" {
		return server.ChatCommand(p.JoinMessage, name)
	}
	return server.ChatCommand(p.LeaveMessage, name)
}

// match acha o pedido de entrada de quem entrou. A entrada confirmada pode trazer o nome cortado, então vale o
// nome igual ou o que começa com ele. Devolve o id e o nome completo (ou o nome recebido, se não achar).
func (s *Service) match(joined string) (string, string) {
	joined = strings.TrimSpace(joined)
	for id, n := range s.logins {
		if n == joined {
			return id, n
		}
	}
	for id, n := range s.logins {
		if joined != "" && strings.HasPrefix(n, joined) {
			return id, n
		}
	}
	return "", joined
}

// Reset esquece os jogadores e fecha a conexão (servidor parou ou reiniciou).
func (s *Service) Reset() {
	s.mu.Lock()
	clear(s.logins)
	clear(s.joined)
	s.mu.Unlock()
	s.sendMu.Lock()
	s.closeConn()
	s.sendMu.Unlock()
}

func (s *Service) worker() {
	for cmd := range s.queue {
		s.Send(cmd)
	}
}

func (s *Service) closeConn() {
	if s.conn != nil {
		_ = s.conn.Close()
		s.conn = nil
	}
}

// Send manda um comando pela conexão aberta (abre se precisar; tenta uma vez de novo com conexão nova).
func (s *Service) Send(cmd string) {
	s.sendMu.Lock()
	defer s.sendMu.Unlock()
	p, sec, err := s.d.Config()
	if err != nil || !p.Rcon || sec.RconPassword == "" {
		return
	}
	for attempt := 0; attempt < 2; attempt++ {
		if s.conn == nil {
			c, err := s.d.Dial(p.RconPort, sec.RconPassword)
			if err != nil {
				s.d.Log("Chat message not sent (RCON): "+err.Error(), "warning")
				return
			}
			s.conn = c
		}
		if _, err := s.conn.Exec(cmd); err == nil {
			return
		}
		s.closeConn()
	}
	s.d.Log("Chat message not sent: the RCON connection failed twice", "warning")
}
