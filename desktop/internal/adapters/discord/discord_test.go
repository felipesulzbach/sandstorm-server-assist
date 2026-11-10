package discord

import (
	"context"
	"strings"
	"testing"
)

// Endereço fora da lista do netx: o erro não pode trazer o endereço (o token do webhook é segredo).
func TestErrorHidesURL(t *testing.T) {
	err := Post(context.Background(), "https://example.invalid/api/webhooks/1/secret-token-value", []byte("{}"))
	if err == nil || strings.Contains(err.Error(), "secret-token-value") {
		t.Fatalf("error leaks the address: %v", err)
	}
}
