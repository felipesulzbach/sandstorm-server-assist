//go:build windows

package updater

import (
	"os/exec"
	"syscall"

	"golang.org/x/sys/windows"
)

// Launcher roda o instalador num processo separado, que continua depois que o app fecha.
type Launcher struct{}

// Start inicia o instalador (baixado na pasta de atualizações e conferido pela assinatura + SHA-256).
func (Launcher) Start(path string, args ...string) error {
	cmd := exec.Command(path, args...) // #nosec G204 -- instalador conferido (assinatura + SHA-256), argumentos fixos
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: windows.DETACHED_PROCESS | windows.CREATE_NEW_PROCESS_GROUP}
	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Process.Release()
}
