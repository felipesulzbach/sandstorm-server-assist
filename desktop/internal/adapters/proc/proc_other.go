//go:build !windows

package proc

import (
	"errors"
	"os/exec"
	"syscall"
	"time"
)

// Process é um processo encontrado.
type Process struct {
	PID     uint32    `json:"pid"`
	Path    string    `json:"path"`
	Started time.Time `json:"started"`
}

var errWindowsOnly = errors.New("only supported on Windows")

func FindInDir(string, string) ([]Process, error)           { return nil, nil }
func Kill(uint32) error                                     { return errWindowsOnly }
func Start(string, string, string, bool) (*exec.Cmd, error) { return nil, errWindowsOnly }
func HideConsole() *syscall.SysProcAttr                     { return nil }
func Exists(uint32) bool                                    { return false }
func IsSelf(uint32) bool                                    { return false }
