// Package notify manda avisos para um canal do Discord por webhook (ADR-022): servidor no ar, parado, travou,
// jogador entrou ou saiu. Os textos são do usuário (Profile.Discord); o endereço do webhook fica no cofre.
// Os envios vão numa fila, um por segundo, e nunca travam o resto do app.
package notify

import (
	"context"
	"encoding/json"
	"time"

	"sandstorm-server-assist/desktop/internal/domain/problems"
	"sandstorm-server-assist/desktop/internal/domain/server"
)

// Tipos de aviso.
const (
	Ready   = "ready"
	Stopped = "stopped"
	Crashed = "crashed"
	Join    = "join"
	Leave   = "leave"
)

// Deps são as partes do sistema que o serviço usa (dublês nos testes).
type Deps struct {
	Config func() (server.Profile, server.Secrets, error)
	// Post manda o JSON para o webhook.
	Post func(ctx context.Context, url string, body []byte) error
	Log  func(text, level string)
	// Gap é o intervalo entre envios (padrão 1 s).
	Gap time.Duration
}

type job struct{ url, text string }

// Service é a fila de avisos.
type Service struct {
	d     Deps
	queue chan job
}

// New cria o serviço. async liga a fila (o app); nos testes, Message devolve o texto sem enviar.
func New(d Deps, async bool) *Service {
	if d.Gap == 0 {
		d.Gap = time.Second
	}
	s := &Service{d: d}
	if async {
		s.queue = make(chan job, 32)
		go s.worker()
	}
	return s
}

// Message monta o endereço e o texto de um aviso (texto vazio se o aviso está desligado ou não há webhook).
func (s *Service) Message(kind, mapName, player string) (string, string) {
	p, sec, err := s.d.Config()
	if err != nil || sec.DiscordWebhook == "" {
		return "", ""
	}
	tpl := map[string]string{Ready: p.Discord.Ready, Stopped: p.Discord.Stopped, Crashed: p.Discord.Crashed,
		Join: p.Discord.Join, Leave: p.Discord.Leave}[kind]
	return sec.DiscordWebhook, server.DiscordText(tpl, p.Hostname, mapName, player)
}

// Event põe um aviso na fila (sem travar).
func (s *Service) Event(kind, mapName, player string) {
	url, text := s.Message(kind, mapName, player)
	if text == "" || s.queue == nil {
		return
	}
	select {
	case s.queue <- job{url, text}:
	default: // fila cheia: perde o aviso
	}
}

// Test manda um aviso de teste agora (o texto de "servidor no ar", ou um texto fixo se estiver vazio).
func (s *Service) Test(mapName string) error {
	p, sec, err := s.d.Config()
	if err != nil {
		return err
	}
	if sec.DiscordWebhook == "" {
		return problems.Field("discordWebhook", "validation.webhook")
	}
	tpl := p.Discord.Ready
	if tpl == "" {
		tpl = server.ServerTag + ": test"
	}
	return s.post(sec.DiscordWebhook, server.DiscordText(tpl, p.Hostname, mapName, ""))
}

func (s *Service) post(url, text string) error {
	body, err := Payload(text)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := s.d.Post(ctx, url, body); err != nil {
		return problems.Wrap(problems.DiscordFailed, err)
	}
	return nil
}

func (s *Service) worker() {
	for j := range s.queue {
		if err := s.post(j.url, j.text); err != nil {
			s.d.Log("Discord message not sent: "+err.Error(), "warning")
		}
		time.Sleep(s.d.Gap)
	}
}

// Payload é o JSON do webhook: sem menções (ninguém é chamado por @) e com o nome do app como remetente.
func Payload(text string) ([]byte, error) {
	return json.Marshal(map[string]any{
		"content":          text,
		"username":         "Sandstorm Server Assist",
		"allowed_mentions": map[string]any{"parse": []string{}},
	})
}
