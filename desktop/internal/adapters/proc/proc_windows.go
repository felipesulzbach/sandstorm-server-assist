//go:build windows

// Package proc inicia, encontra e encerra os processos do servidor no Windows (ADR-003).
package proc

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Process é um processo encontrado. Started é a hora em que o processo começou (zero se o Windows não informar).
// CmdLine é a linha de comando do processo (vazia se o Windows não deixar ler): diz de qual instância do servidor
// ele é (ADR-048). Nunca vai para a tela nem para o log (pode ter senhas).
type Process struct {
	PID     uint32    `json:"pid"`
	Path    string    `json:"path"`
	Started time.Time `json:"started"`
	CmdLine string    `json:"-"`
}

// FindInDir devolve os processos cujo executável começa com prefix (sem diferenciar maiúsculas) e fica dentro de
// dir. Processos que o app não consegue abrir (de outro usuário, protegidos) são ignorados.
func FindInDir(dir, prefix string) ([]Process, error) {
	snap, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, err
	}
	defer windows.CloseHandle(snap)
	var e windows.ProcessEntry32
	e.Size = uint32(unsafe.Sizeof(e))
	root := strings.ToLower(filepath.Clean(dir)) + string(filepath.Separator)
	prefix = strings.ToLower(prefix)
	var out []Process
	for err = windows.Process32First(snap, &e); err == nil; err = windows.Process32Next(snap, &e) {
		name := strings.ToLower(windows.UTF16ToString(e.ExeFile[:]))
		if !strings.HasPrefix(name, prefix) {
			continue
		}
		path, started, cmdline, perr := imageInfo(e.ProcessID)
		if perr != nil {
			continue
		}
		if strings.HasPrefix(strings.ToLower(path), root) {
			out = append(out, Process{PID: e.ProcessID, Path: path, Started: started, CmdLine: cmdline})
		}
	}
	return out, nil
}

// imageInfo devolve o caminho do executável, a hora de início e a linha de comando do processo.
func imageInfo(pid uint32) (string, time.Time, string, error) {
	h, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if err != nil {
		return "", time.Time{}, "", err
	}
	defer windows.CloseHandle(h)
	buf := make([]uint16, windows.MAX_LONG_PATH)
	n := uint32(len(buf)) // #nosec G115 -- buffer fixo (MAX_LONG_PATH)
	if err := windows.QueryFullProcessImageName(h, 0, &buf[0], &n); err != nil {
		return "", time.Time{}, "", err
	}
	var started time.Time
	var created, exited, kernel, user windows.Filetime
	if windows.GetProcessTimes(h, &created, &exited, &kernel, &user) == nil && created.Nanoseconds() > 0 {
		started = time.Unix(0, created.Nanoseconds())
	}
	return windows.UTF16ToString(buf[:n]), started, commandLine(h), nil
}

// maxCmdLine limita o tamanho da linha de comando lida (o Windows aceita até 32 767 caracteres).
const maxCmdLine = 64 << 10

// commandLine lê a linha de comando do processo (ProcessCommandLineInformation, Windows 8.1 ou mais novo; basta
// PROCESS_QUERY_LIMITED_INFORMATION). Vazia se não der.
func commandLine(h windows.Handle) string {
	var size uint32
	_ = windows.NtQueryInformationProcess(h, windows.ProcessCommandLineInformation, nil, 0, &size)
	if size < uint32(unsafe.Sizeof(windows.NTUnicodeString{})) || size > maxCmdLine {
		return ""
	}
	// Buffer em uint64 para o UNICODE_STRING do começo ficar alinhado.
	buf := make([]uint64, (size+7)/8)
	if err := windows.NtQueryInformationProcess(h, windows.ProcessCommandLineInformation,
		unsafe.Pointer(&buf[0]), size, &size); err != nil { // #nosec G103 -- API do Windows, buffer do tamanho pedido
		return ""
	}
	us := (*windows.NTUnicodeString)(unsafe.Pointer(&buf[0])) // #nosec G103 -- o Windows grava um UNICODE_STRING
	return us.String()
}

// Kill encerra um processo.
func Kill(pid uint32) error {
	h, err := windows.OpenProcess(windows.PROCESS_TERMINATE, false, pid)
	if err != nil {
		return err
	}
	defer windows.CloseHandle(h)
	return windows.TerminateProcess(h, 1)
}

// Start inicia exe na pasta dir com a linha de comando já montada pelo domínio (sem shell). O processo continua
// vivo se o app fechar. hidden=true não abre janela de console (os logs aparecem no app).
func Start(exe, cmdline, dir string, hidden bool) (*exec.Cmd, error) {
	cmd := exec.Command(exe) // #nosec G204 -- executável do servidor na pasta configurada; argumentos do domínio
	cmd.Dir = dir
	var flags uint32 = windows.CREATE_NEW_PROCESS_GROUP
	if hidden {
		flags |= windows.CREATE_NO_WINDOW
	} else {
		flags |= windows.CREATE_NEW_CONSOLE
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{
		// O Unreal lê a linha de comando inteira; montada pelo domínio (server.BuildLaunch), com aspas onde precisa.
		CmdLine:       syscall.EscapeArg(exe) + " " + cmdline,
		CreationFlags: flags,
		HideWindow:    hidden,
	}
	cmd.Stdin, cmd.Stdout, cmd.Stderr = nil, nil, nil
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return cmd, nil
}

// HideConsole devolve atributos para rodar um programa de console (SteamCMD) sem abrir janela.
func HideConsole() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{HideWindow: true, CreationFlags: windows.CREATE_NO_WINDOW}
}

// Exists diz se o processo ainda está vivo.
func Exists(pid uint32) bool {
	h, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if err != nil {
		return false
	}
	defer windows.CloseHandle(h)
	var code uint32
	if err := windows.GetExitCodeProcess(h, &code); err != nil {
		return false
	}
	return code == 259 // STILL_ACTIVE
}

// self é o PID do app (nunca encerrado por engano).
var self = uint32(os.Getpid()) // #nosec G115 -- PID cabe em uint32 no Windows

// IsSelf diz se o PID é o do próprio app.
func IsSelf(pid uint32) bool { return pid == self }
