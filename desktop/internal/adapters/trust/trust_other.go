//go:build !windows

// Package trust confere a assinatura digital (Authenticode) de um executável antes de rodar (ADR-007).
package trust

import "errors"

// ErrUntrusted indica assinatura ausente, inválida ou de outro signatário.
var ErrUntrusted = errors.New("untrusted signature")

// VerifySigner só existe no Windows.
func VerifySigner(string, string) error { return ErrUntrusted }
