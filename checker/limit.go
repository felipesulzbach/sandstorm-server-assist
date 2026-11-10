package main

import (
	"sync"
	"time"
)

// Limiter conta pedidos por minuto, por IP e no total (janela fixa de um minuto).
type Limiter struct {
	perIP, all int
	mu         sync.Mutex
	window     time.Time
	total      int
	counts     map[string]int
}

// NewLimiter cria o limite.
func NewLimiter(perIP, all int) *Limiter {
	return &Limiter{perIP: perIP, all: all, counts: map[string]int{}}
}

// Allow diz se o pedido pode seguir e o conta.
func (l *Limiter) Allow(key string, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	w := now.Truncate(time.Minute)
	if !w.Equal(l.window) {
		l.window, l.total = w, 0
		clear(l.counts)
	}
	if l.total >= l.all || l.counts[key] >= l.perIP {
		return false
	}
	l.total++
	l.counts[key]++
	return true
}
