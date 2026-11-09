//go:build !windows

package vault

import "errors"

var errNoDPAPI = errors.New("DPAPI is only available on Windows")

func protect([]byte) ([]byte, error)   { return nil, errNoDPAPI }
func unprotect([]byte) ([]byte, error) { return nil, errNoDPAPI }
