// Package ini lê e troca chaves de arquivos .ini do Unreal sem mexer no resto (comentários, ordem, outras
// seções e quebras de linha continuam iguais).
package ini

import "strings"

// Get devolve o valor de uma chave numa seção ([Rcon], [/Script/...]). Seção e chave sem diferenciar maiúsculas.
func Get(content, section, key string) (string, bool) {
	cur := ""
	for _, line := range splitLines(content) {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "[") && strings.HasSuffix(t, "]") {
			cur = t[1 : len(t)-1]
			continue
		}
		if strings.EqualFold(cur, section) {
			if k, v, ok := strings.Cut(t, "="); ok && strings.EqualFold(strings.TrimSpace(k), key) {
				return strings.TrimSpace(v), true
			}
		}
	}
	return "", false
}

// Remove apaga uma chave de uma seção. Devolve o conteúdo novo e se algo mudou.
func Remove(content, section, key string) (string, bool) {
	lines := splitLines(content)
	out := make([]string, 0, len(lines))
	cur, changed := "", false
	for _, line := range lines {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "[") && strings.HasSuffix(t, "]") {
			cur = t[1 : len(t)-1]
		} else if strings.EqualFold(cur, section) {
			if k, _, ok := strings.Cut(t, "="); ok && strings.EqualFold(strings.TrimSpace(k), key) {
				changed = true
				continue
			}
		}
		out = append(out, line)
	}
	return strings.Join(out, newline(content)), changed
}

// Set grava uma chave (troca se existe; senão cria na seção, ou cria a seção no fim).
func Set(content, section, key, value string) string {
	nl := newline(content)
	lines := splitLines(content)
	cur, sectionEnd, found := "", -1, false
	for i, line := range lines {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "[") && strings.HasSuffix(t, "]") {
			cur = t[1 : len(t)-1]
			if strings.EqualFold(cur, section) {
				sectionEnd, found = i+1, true
			}
			continue
		}
		if strings.EqualFold(cur, section) {
			if k, _, ok := strings.Cut(t, "="); ok && strings.EqualFold(strings.TrimSpace(k), key) {
				lines[i] = key + "=" + value
				return strings.Join(lines, nl)
			}
			if t != "" {
				sectionEnd = i + 1
			}
		}
	}
	entry := key + "=" + value
	if !found {
		if len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) != "" {
			lines = append(lines, "")
		}
		lines = append(lines, "["+section+"]", entry, "")
		return strings.Join(lines, nl)
	}
	lines = append(lines[:sectionEnd], append([]string{entry}, lines[sectionEnd:]...)...)
	return strings.Join(lines, nl)
}

func splitLines(s string) []string {
	return strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n")
}

func newline(s string) string {
	if strings.Contains(s, "\r\n") || s == "" {
		return "\r\n"
	}
	return "\n"
}
