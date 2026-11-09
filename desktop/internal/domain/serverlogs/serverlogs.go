// Package serverlogs tem as regras dos arquivos de log do servidor (ADR-018): o servidor grava uma partida por
// arquivo (Insurgency.log é a atual; a anterior vira Insurgency-backup-<data>.log a cada início). Aqui ficam quais
// arquivos são logs do servidor, a ordem das partidas e quais apagar mantendo as N mais novas. Sem I/O.
package serverlogs

import (
	"regexp"
	"sort"
	"time"
)

// Current é o log da partida atual (nunca é apagado).
const Current = "Insurgency.log"

// Quantas partidas manter: no mínimo a atual e a anterior (o app lê os logs recentes para "Trazer do servidor" e
// para achar a chave do mod.io).
const (
	MinKeep     = 2
	MaxKeep     = 100
	DefaultKeep = 10
)

// File é um arquivo de log da pasta do servidor.
type File struct {
	Name    string    `json:"name"`
	ModTime time.Time `json:"modTime"`
	Size    int64     `json:"size"`
}

// reName aceita só os nomes que o servidor cria: Insurgency.log, Insurgency_2.log (segunda instância) e
// Insurgency-backup-2026.11.09-22.56.07.log.
var reName = regexp.MustCompile(`^Insurgency(?:_\d{1,3}|-backup-\d{4}\.\d{2}\.\d{2}-\d{2}\.\d{2}\.\d{2})?\.log$`)

// IsServerLog diz se o nome é de um log do servidor.
func IsServerLog(name string) bool { return reName.MatchString(name) }

// ValidKeep diz se a quantidade de partidas a manter está no intervalo permitido.
func ValidKeep(keep int) bool { return keep >= MinKeep && keep <= MaxKeep }

// Sort ordena do mais novo para o mais velho; no mesmo horário, o log atual vem antes e depois o nome maior.
func Sort(files []File) {
	sort.SliceStable(files, func(i, j int) bool {
		a, b := files[i], files[j]
		if !a.ModTime.Equal(b.ModTime) {
			return a.ModTime.After(b.ModTime)
		}
		if (a.Name == Current) != (b.Name == Current) {
			return a.Name == Current
		}
		return a.Name > b.Name
	})
}

// ToDelete devolve os logs a apagar para ficar só com as keep partidas mais novas. O log atual sempre conta como
// mantido e nunca é devolvido; nomes que não são logs do servidor são ignorados. keep fora do intervalo vira o
// limite mais próximo.
func ToDelete(files []File, keep int) []File {
	if keep < MinKeep {
		keep = MinKeep
	}
	if keep > MaxKeep {
		keep = MaxKeep
	}
	var logs []File
	hasCurrent := false
	for _, f := range files {
		if !IsServerLog(f.Name) {
			continue
		}
		if f.Name == Current {
			hasCurrent = true
			continue
		}
		logs = append(logs, f)
	}
	Sort(logs)
	if hasCurrent {
		keep--
	}
	if len(logs) <= keep {
		return nil
	}
	return logs[keep:]
}
