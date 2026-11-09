//go:build windows

// Package trust confere a assinatura digital (Authenticode) de um executável antes de rodar (ADR-007).
package trust

import (
	"errors"
	"fmt"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	wintrust                = windows.NewLazySystemDLL("wintrust.dll")
	procProvDataFromState   = wintrust.NewProc("WTHelperProvDataFromStateData")
	procProvSignerFromChain = wintrust.NewProc("WTHelperGetProvSignerFromChain")
)

// Layout de CRYPT_PROVIDER_SGNR e CRYPT_PROVIDER_CERT (wintrust.h), só os campos usados.
type providerSigner struct {
	CbStruct     uint32
	VerifyAsOf   windows.Filetime
	CsCertChain  uint32
	PasCertChain *providerCert
}

type providerCert struct {
	CbStruct uint32
	PCert    *windows.CertContext
}

// ErrUntrusted indica assinatura ausente, inválida ou de outro signatário.
var ErrUntrusted = errors.New("untrusted signature")

// VerifySigner confere que o arquivo tem assinatura válida (cadeia confiável no Windows) e que o nome do
// signatário contém expected (ex.: "Valve").
func VerifySigner(path, expected string) error {
	p, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return err
	}
	file := &windows.WinTrustFileInfo{Size: uint32(unsafe.Sizeof(windows.WinTrustFileInfo{})), FilePath: p}
	data := &windows.WinTrustData{
		Size:                            uint32(unsafe.Sizeof(windows.WinTrustData{})),
		UIChoice:                        windows.WTD_UI_NONE,
		RevocationChecks:                windows.WTD_REVOKE_NONE,
		UnionChoice:                     windows.WTD_CHOICE_FILE,
		StateAction:                     windows.WTD_STATEACTION_VERIFY,
		FileOrCatalogOrBlobOrSgnrOrCert: unsafe.Pointer(file), // #nosec G103 -- estrutura exigida pela API
	}
	verr := windows.WinVerifyTrustEx(windows.InvalidHWND, &windows.WINTRUST_ACTION_GENERIC_VERIFY_V2, data)
	defer func() {
		data.StateAction = windows.WTD_STATEACTION_CLOSE
		_ = windows.WinVerifyTrustEx(windows.InvalidHWND, &windows.WINTRUST_ACTION_GENERIC_VERIFY_V2, data)
	}()
	if verr != nil {
		return fmt.Errorf("%w: %v", ErrUntrusted, verr)
	}
	name, err := signerName(data.StateData)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrUntrusted, err)
	}
	if !strings.Contains(strings.ToLower(name), strings.ToLower(expected)) {
		return fmt.Errorf("%w: signer %q", ErrUntrusted, name)
	}
	return nil
}

func signerName(state windows.Handle) (string, error) {
	prov, _, _ := procProvDataFromState.Call(uintptr(state))
	if prov == 0 {
		return "", errors.New("no provider data")
	}
	sgnr, _, _ := procProvSignerFromChain.Call(prov, 0, 0, 0)
	if sgnr == 0 {
		return "", errors.New("no signer")
	}
	s := *(**providerSigner)(unsafe.Pointer(&sgnr)) // #nosec G103 -- ponteiro devolvido pela própria API
	if s.CsCertChain == 0 || s.PasCertChain == nil {
		return "", errors.New("empty chain")
	}
	c := s.PasCertChain // primeiro certificado = o do signatário
	if c.PCert == nil {
		return "", errors.New("no certificate")
	}
	buf := make([]uint16, 256)
	n := windows.CertGetNameString(c.PCert, windows.CERT_NAME_SIMPLE_DISPLAY_TYPE, 0, nil, &buf[0], uint32(len(buf))) // #nosec G115 -- buffer fixo de 256
	if n <= 1 {
		return "", errors.New("no signer name")
	}
	return windows.UTF16ToString(buf[:n]), nil
}
