//go:build windows

package vault

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

func blob(b []byte) *windows.DataBlob {
	if len(b) == 0 {
		return &windows.DataBlob{}
	}
	return &windows.DataBlob{Size: uint32(len(b)), Data: &b[0]} // #nosec G115 -- cofre tem no máximo 1 MiB
}

func protect(plain []byte) ([]byte, error) {
	var out windows.DataBlob
	if err := windows.CryptProtectData(blob(plain), nil, blob(entropy), 0, nil, windows.CRYPTPROTECT_UI_FORBIDDEN, &out); err != nil {
		return nil, err
	}
	return takeBlob(&out), nil
}

func unprotect(data []byte) ([]byte, error) {
	var out windows.DataBlob
	if err := windows.CryptUnprotectData(blob(data), nil, blob(entropy), 0, nil, windows.CRYPTPROTECT_UI_FORBIDDEN, &out); err != nil {
		return nil, err
	}
	return takeBlob(&out), nil
}

// takeBlob copia o resultado para a memória do Go e libera (zerando) o buffer do Windows.
func takeBlob(b *windows.DataBlob) []byte {
	if b.Data == nil {
		return nil
	}
	src := unsafe.Slice(b.Data, b.Size) // #nosec G103 -- buffer devolvido pelo DPAPI, tamanho informado por ele
	out := make([]byte, len(src))
	copy(out, src)
	for i := range src {
		src[i] = 0
	}
	_, _ = windows.LocalFree(windows.Handle(unsafe.Pointer(b.Data))) // #nosec G103 -- liberar memória do DPAPI
	return out
}
