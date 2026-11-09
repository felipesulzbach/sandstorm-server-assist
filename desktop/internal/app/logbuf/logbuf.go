// Package logbuf guarda as últimas linhas de log (servidor, SteamCMD e app) para a gaveta de logs da tela.
// Toda linha entra já mascarada (ADR-002).
package logbuf

import (
	"sync"
	"time"

	"sandstorm-server-assist/desktop/internal/domain/logparse"
	"sandstorm-server-assist/desktop/internal/domain/redact"
)

// Line é uma linha mostrada na tela.
type Line struct {
	Seq    int64     `json:"seq"`
	Time   time.Time `json:"time"`
	Source string    `json:"source"` // server | steamcmd | app
	Level  string    `json:"level"`  // error | warning | info
	Noise  bool      `json:"noise"`
	Event  bool      `json:"event"` // a linha virou um evento (fase, mod, erro)
	Text   string    `json:"text"`
	// Run é a partida do servidor (cresce a cada início); Kind "run" marca a linha que abre uma partida.
	Run  int64  `json:"run"`
	Kind string `json:"kind,omitempty"`
}

// Buffer é um anel com capacidade fixa.
type Buffer struct {
	mu      sync.Mutex
	lines   []Line
	next    int64
	cap     int
	secrets func() []string
	notify  func([]Line)
	pending []Line
	run     int64
}

// New cria o buffer. secrets devolve os valores do cofre a mascarar (pode ser nil).
func New(capacity int, secrets func() []string) *Buffer {
	return &Buffer{cap: capacity, secrets: secrets}
}

// SetNotify define quem recebe as linhas novas (chamado por Flush, em lote).
func (b *Buffer) SetNotify(fn func([]Line)) {
	b.mu.Lock()
	b.notify = fn
	b.mu.Unlock()
}

// Add guarda uma linha (mascarada).
func (b *Buffer) Add(source, text string, level string) {
	l := b.Format(source, text, level)
	b.mu.Lock()
	defer b.mu.Unlock()
	b.next++
	l.Seq, l.Run = b.next, b.run
	b.push(l)
}

// Format mascara e classifica uma linha sem guardar (também usada para ler logs de partidas antigas).
func (b *Buffer) Format(source, text, level string) Line {
	text = redact.Text(text)
	if b.secrets != nil {
		text = redact.Values(text, b.secrets()...)
	}
	if level == "" {
		level = logparse.Level(text)
	}
	_, isEvent := logparse.Parse(text)
	return Line{Time: time.Now(), Source: source, Level: level, Noise: logparse.Noise(text), Event: isEvent, Text: text}
}

// StartRun abre uma partida nova: as linhas seguintes ganham o número dela, e entra uma linha divisória com a
// hora de início (a tela escreve o texto).
func (b *Buffer) StartRun(at time.Time) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.run++
	b.next++
	b.push(Line{Seq: b.next, Time: at, Source: "app", Level: "info", Run: b.run, Kind: "run"})
}

// push guarda a linha no anel (com b.mu travado).
func (b *Buffer) push(l Line) {
	b.lines = append(b.lines, l)
	if len(b.lines) > b.cap {
		b.lines = append(b.lines[:0], b.lines[len(b.lines)-b.cap:]...)
	}
	b.pending = append(b.pending, l)
}

// Flush envia as linhas pendentes ao notify (a ponte chama a cada ~300 ms).
func (b *Buffer) Flush() {
	b.mu.Lock()
	p, fn := b.pending, b.notify
	b.pending = nil
	b.mu.Unlock()
	if fn != nil && len(p) > 0 {
		fn(p)
	}
}

// Recent devolve as últimas n linhas.
func (b *Buffer) Recent(n int) []Line {
	b.mu.Lock()
	defer b.mu.Unlock()
	if n <= 0 || n > len(b.lines) {
		n = len(b.lines)
	}
	out := make([]Line, n)
	copy(out, b.lines[len(b.lines)-n:])
	return out
}

// Clear esvazia o buffer.
func (b *Buffer) Clear() {
	b.mu.Lock()
	b.lines = nil
	b.pending = nil
	b.mu.Unlock()
}
