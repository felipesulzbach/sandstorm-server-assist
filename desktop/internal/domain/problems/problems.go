// Package problems é o catálogo de erros do app (ADR-005). O Go só conhece os códigos e os links de ajuda; os
// textos (o que aconteceu, causa provável, como resolver) ficam nas traduções do front (ADR-009), na chave
// "problems.<código>".
package problems

import (
	"encoding/json"
	"errors"

	"sandstorm-server-assist/desktop/internal/domain/redact"
)

// Códigos do catálogo.
const (
	Unknown            = "unknown"
	ServerNotInstalled = "server-not-installed"
	ServerDirInvalid   = "server-dir-invalid"
	ServerRunning      = "server-already-running"
	ServerNotRunning   = "server-not-running"
	ServerStartFailed  = "server-start-failed"
	ServerCrashed      = "server-crashed"
	ServerStopFailed   = "server-stop-failed"
	PortInUse          = "port-in-use"
	ModConflict        = "mod-conflict"
	ModioCodeInvalid   = "modio-code-invalid"
	ModioNotAuthorized = "modio-not-authorized"
	ModioRequestFailed = "modio-request-failed"
	ModioEmailInvalid  = "modio-email-invalid"
	ModioNoAPIKey      = "modio-no-api-key" // #nosec G101 -- código de erro, não credencial
	SteamCMDDownload   = "steamcmd-download"
	SteamCMDSignature  = "steamcmd-signature"
	SteamCMDFailed     = "steamcmd-failed"
	DiskFull           = "disk-full"
	NoInternet         = "no-internet"
	InvalidField       = "invalid-field"
	VaultUnavailable   = "vault-unavailable"
	BackupPassword     = "backup-password"
	BackupCorrupt      = "backup-corrupt"
	BackupWrite        = "backup-write"
	ConfigWrite        = "config-write"
	Busy               = "busy"
	Canceled           = "canceled"
	UpdateFailed       = "update-failed"
	ModsNotMounted     = "mods-not-mounted"
	LogDelete          = "log-delete"
	RconOff            = "rcon-off"
	RconFailed         = "rcon-failed"
	RconAuth           = "rcon-auth"
	MapChangeFailed    = "map-change-failed"
	DiscordFailed      = "discord-failed"
	FirewallFailed     = "firewall-failed"
	AutostartFailed    = "autostart-failed"
	ServerUpdateFailed = "server-update-failed"
	ReportFailed       = "report-failed"
	ReportImage        = "report-image"
	ModDelete          = "mod-delete"
)

// links de ajuda de cada problema (chaves de links conhecidos, abertos pelo Go).
var links = map[string][]string{
	ModConflict:        {"modio-subscriptions"},
	ModioNotAuthorized: {"modio"},
	ModioRequestFailed: {"modio-status"},
	ModioNoAPIKey:      {"modio-api-access"},
	SteamCMDDownload:   {"steamcmd-docs"},
	SteamCMDFailed:     {"steamcmd-docs"},
	UpdateFailed:       {"download-page"},
	ModsNotMounted:     {"modio-subscriptions"},
}

// Codes lista todos os códigos (o teste do front confere que cada um tem tradução).
var Codes = []string{
	Unknown, ServerNotInstalled, ServerDirInvalid, ServerRunning, ServerNotRunning, ServerStartFailed,
	ServerCrashed, ServerStopFailed, PortInUse, ModConflict, ModioCodeInvalid, ModioNotAuthorized,
	ModioRequestFailed, ModioEmailInvalid, ModioNoAPIKey, SteamCMDDownload, SteamCMDSignature, SteamCMDFailed,
	DiskFull, NoInternet, InvalidField, VaultUnavailable, BackupPassword, BackupCorrupt, BackupWrite,
	ConfigWrite, Busy, Canceled, UpdateFailed, ModsNotMounted, LogDelete,
	RconOff, RconFailed, RconAuth, MapChangeFailed,
	DiscordFailed, FirewallFailed, AutostartFailed, ServerUpdateFailed, ReportFailed, ReportImage, ModDelete,
}

// Links devolve os links de ajuda de um código.
func Links(code string) []string { return links[code] }

// Error é um erro do app ligado a um item do catálogo.
type Error struct {
	Code string `json:"code"`
	// Field é o campo da tela com problema (validação).
	Field string `json:"field,omitempty"`
	// Msg é a chave de tradução da mensagem do campo ("validation.port"), com Params para preencher.
	Msg    string            `json:"msg,omitempty"`
	Params map[string]string `json:"params,omitempty"`
	// Detail é texto técnico extra (mensagem original do sistema), já sem segredos. Não é traduzido.
	Detail string   `json:"detail,omitempty"`
	Links  []string `json:"links,omitempty"`
}

// Error vira JSON: é assim que chega à tela pela ponte do Wails (a promessa rejeitada traz a mensagem).
func (e *Error) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

// New cria um erro do catálogo com detalhe técnico (mascarado).
func New(code, detail string) *Error {
	return &Error{Code: code, Detail: redact.Text(detail), Links: links[code]}
}

// Field cria um erro de campo inválido. msg é a chave de tradução; params são pares chave, valor.
func Field(field, msg string, params ...string) *Error {
	e := &Error{Code: InvalidField, Field: field, Msg: msg}
	if len(params) > 1 {
		e.Params = map[string]string{}
		for i := 0; i+1 < len(params); i += 2 {
			e.Params[params[i]] = redact.Text(params[i+1])
		}
	}
	return e
}

// Wrap liga um erro qualquer a um código. Se err já é um *Error, ele é mantido.
func Wrap(code string, err error) error {
	if err == nil {
		return nil
	}
	var pe *Error
	if errors.As(err, &pe) {
		return pe
	}
	return New(code, err.Error())
}

// CodeOf devolve o código de um erro (Unknown se não for do catálogo).
func CodeOf(err error) string {
	var pe *Error
	if errors.As(err, &pe) {
		return pe.Code
	}
	return Unknown
}
