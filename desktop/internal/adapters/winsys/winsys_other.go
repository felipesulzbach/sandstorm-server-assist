//go:build !windows

package winsys

import "errors"

var errOnlyWindows = errors.New("only on Windows")

// ErrCanceled: o usuário não confirmou a janela de administrador.
var ErrCanceled = errors.New("canceled by the user")

// SetAutostart só existe no Windows.
func SetAutostart(bool) error { return errOnlyWindows }

// Autostart só existe no Windows.
func Autostart() bool { return false }

// FirewallRule só existe no Windows.
func FirewallRule(string) bool { return false }

// AddFirewallRules só existe no Windows.
func AddFirewallRules(string) error { return errOnlyWindows }

// RunElevated só existe no Windows.
func RunElevated(string) error { return errOnlyWindows }
