//go:build windows

package winsys

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

const (
	runKey   = `Software\Microsoft\Windows\CurrentVersion\Run`
	runValue = "SandstormServerAssist"
)

// SetAutostart liga ou desliga a abertura do app ao entrar no Windows (registro do usuário, sem administrador).
func SetAutostart(on bool) error {
	k, _, err := registry.CreateKey(registry.CURRENT_USER, runKey, registry.SET_VALUE|registry.QUERY_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()
	if !on {
		if err := k.DeleteValue(runValue); err != nil && !errors.Is(err, registry.ErrNotExist) {
			return err
		}
		return nil
	}
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	return k.SetStringValue(runValue, `"`+exe+`" `+AutostartFlag)
}

// Autostart diz se o app está no registro para abrir com o Windows.
func Autostart() bool {
	k, err := registry.OpenKey(registry.CURRENT_USER, runKey, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer k.Close()
	v, _, err := k.GetStringValue(runValue)
	return err == nil && strings.Contains(v, AutostartFlag)
}

func netsh(args string) error {
	// Pasta do sistema pela API do Windows (não por variável de ambiente, que outro programa poderia trocar).
	sys, err := windows.GetSystemDirectory()
	if err != nil {
		return err
	}
	exe := filepath.Join(sys, "netsh.exe")
	cmd := exec.Command(exe) // #nosec G204 -- netsh do Windows; argumentos montados aqui com nomes fixos
	cmd.SysProcAttr = &syscall.SysProcAttr{
		// netsh lê a linha inteira: name="..." e program="..." com as aspas depois do =.
		CmdLine:       syscall.EscapeArg(exe) + " " + args,
		HideWindow:    true,
		CreationFlags: windows.CREATE_NO_WINDOW,
	}
	return cmd.Run()
}

// FirewallRule diz se a regra de um executável existe (não precisa de administrador).
func FirewallRule(program string) bool {
	return netsh(`advfirewall firewall show rule name="`+RuleName(program)+`"`) == nil
}

// AddFirewallRules cria (ou recria) as regras de entrada para os executáveis do servidor que existem na pasta.
// Precisa rodar como administrador (o app abre uma cópia de si com FirewallFlag).
func AddFirewallRules(serverDir string) error {
	added := 0
	for _, prog := range ServerPrograms(serverDir) {
		if strings.ContainsAny(prog, `"`) {
			return errors.New("invalid path")
		}
		if info, err := os.Stat(prog); err != nil || info.IsDir() {
			continue
		}
		name := RuleName(prog)
		_ = netsh(`advfirewall firewall delete rule name="` + name + `"`)
		if err := netsh(`advfirewall firewall add rule name="` + name + `" dir=in action=allow program="` + prog +
			`" enable=yes profile=any`); err != nil {
			return err
		}
		added++
	}
	if added == 0 {
		return errors.New("server executables not found")
	}
	return nil
}

// shellExecuteInfo é o SHELLEXECUTEINFOW.
type shellExecuteInfo struct {
	cbSize       uint32
	fMask        uint32
	hwnd         uintptr
	lpVerb       *uint16
	lpFile       *uint16
	lpParameters *uint16
	lpDirectory  *uint16
	nShow        int32
	hInstApp     uintptr
	lpIDList     uintptr
	lpClass      *uint16
	hkeyClass    uintptr
	dwHotKey     uint32
	hIcon        uintptr
	hProcess     windows.Handle
}

var procShellExecuteEx = windows.NewLazySystemDLL("shell32.dll").NewProc("ShellExecuteExW")

// ErrCanceled: o usuário não confirmou a janela de administrador do Windows.
var ErrCanceled = errors.New("canceled by the user")

// RunElevated abre o próprio app como administrador com um argumento fixo e espera terminar (até 2 minutos).
func RunElevated(arg string) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	verb, _ := windows.UTF16PtrFromString("runas")
	file, _ := windows.UTF16PtrFromString(exe)
	params, _ := windows.UTF16PtrFromString(arg)
	info := shellExecuteInfo{fMask: 0x40 | 0x100, lpVerb: verb, lpFile: file, lpParameters: params} // NOCLOSEPROCESS | NOASYNC
	info.cbSize = uint32(unsafe.Sizeof(info))                                                       // #nosec G103 G115 -- tamanho da estrutura da API do Windows
	ok, _, callErr := procShellExecuteEx.Call(uintptr(unsafe.Pointer(&info)))                       // #nosec G103 -- chamada da API do Windows
	if ok == 0 {
		if errors.Is(callErr, windows.ERROR_CANCELLED) {
			return ErrCanceled
		}
		return callErr
	}
	if info.hProcess == 0 {
		return errors.New("no process")
	}
	defer windows.CloseHandle(info.hProcess)
	if ev, err := windows.WaitForSingleObject(info.hProcess, 120000); err != nil || ev != windows.WAIT_OBJECT_0 {
		return errors.New("timeout")
	}
	var code uint32
	if err := windows.GetExitCodeProcess(info.hProcess, &code); err != nil {
		return err
	}
	if code != 0 {
		return errors.New("the rule was not created")
	}
	return nil
}

// WindowsVersion devolve a versão do Windows para o relato de problema, ex.: "Windows 11 (build 26100)".
func WindowsVersion() string {
	v := windows.RtlGetVersion()
	name := "Windows " + strconv.Itoa(int(v.MajorVersion))
	if v.MajorVersion == 10 && v.BuildNumber >= 22000 {
		name = "Windows 11"
	}
	return name + " (build " + strconv.Itoa(int(v.BuildNumber)) + ")"
}

// RevealFile abre o Explorer com o arquivo selecionado.
func RevealFile(path string) error {
	return exec.Command("explorer.exe", "/select,"+path).Start() // #nosec G204 -- explorer.exe com um arquivo criado pelo app
}
