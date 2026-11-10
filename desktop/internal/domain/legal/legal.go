// Package legal guarda a versão em vigor dos termos de uso e do aviso de dados (ADR-015). Os textos ficam nas
// traduções do front (chave "legal.*") e em TERMS.md / PRIVACY.md; mudou o texto, muda a versão, e o app pede o
// aceite de novo.
package legal

// Version é a versão dos termos de uso e do aviso de dados (data em que passaram a valer).
const Version = "2026-11-10"

// Accepted diz se a versão aceita é a que está em vigor.
func Accepted(version string) bool { return version == Version }
