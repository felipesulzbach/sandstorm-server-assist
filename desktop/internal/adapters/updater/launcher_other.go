//go:build !windows

package updater

import "errors"

// Launcher: a atualização automática só existe no Windows (instalador NSIS).
type Launcher struct{}

func (Launcher) Start(string, ...string) error { return errors.New("automatic update only on Windows") }
