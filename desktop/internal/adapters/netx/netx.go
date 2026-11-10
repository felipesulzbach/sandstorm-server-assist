// Package netx tem o cliente HTTP do app: só HTTPS, só para hosts permitidos, com limite de tamanho e progresso.
package netx

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"sandstorm-server-assist/desktop/internal/domain/problems"
)

// AllowedHosts são os únicos endereços que o app acessa (ADR-007, ADR-010).
var AllowedHosts = map[string]bool{
	"steamcdn-a.akamaihd.net": true, // SteamCMD (Valve)
	"g-254.modapi.io":         true, // API do mod.io do Insurgency: Sandstorm
	"game-in.duckdns.org":     true, // página de download e atualização automática do app (ADR-014)
	// Webhook do Discord, só se o usuário configurar (ADR-022).
	"discord.com":        true,
	"discordapp.com":     true,
	"ptb.discord.com":    true,
	"canary.discord.com": true,
}

// Client faz pedidos só para AllowedHosts, sem seguir redirecionamento para outro host.
var Client = &http.Client{
	Timeout: 0, // downloads longos: o limite vem do contexto
	Transport: &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           (&net.Dialer{Timeout: 15 * time.Second}).DialContext,
		TLSHandshakeTimeout:   15 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
		ForceAttemptHTTP2:     true,
	},
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) > 5 {
			return errors.New("too many redirects")
		}
		return CheckURL(req.URL)
	},
}

// CheckURL recusa o que não for HTTPS para um host permitido.
func CheckURL(u *url.URL) error {
	if u.Scheme != "https" || !AllowedHosts[strings.ToLower(u.Hostname())] {
		return fmt.Errorf("blocked address: %s", u.Hostname())
	}
	return nil
}

// Do faz o pedido conferindo o endereço.
func Do(req *http.Request) (*http.Response, error) {
	if err := CheckURL(req.URL); err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "SandstormServerAssist")
	return Client.Do(req) // #nosec G704 -- endereço conferido por CheckURL (lista fixa de hosts)
}

// NetProblem traduz erros de rede para o catálogo.
func NetProblem(code string, err error) error {
	var ne net.Error
	var de *net.DNSError
	if errors.As(err, &de) || (errors.As(err, &ne) && ne.Timeout()) {
		return problems.Wrap(problems.NoInternet, err)
	}
	return problems.Wrap(code, err)
}

// Progress é chamado durante um download.
type Progress func(done, total int64)

// Download baixa rawURL para o arquivo dest (criado do zero), até maxBytes, chamando onProgress.
func Download(ctx context.Context, rawURL, dest string, maxBytes int64, onProgress Progress) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return err
	}
	resp, err := Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	if resp.ContentLength > maxBytes {
		return fmt.Errorf("file too large (%d bytes)", resp.ContentLength)
	}
	f, err := os.OpenFile(dest, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600) // #nosec G304 -- arquivo temporário do app
	if err != nil {
		return err
	}
	var done int64
	buf := make([]byte, 64<<10)
	last := time.Time{}
	for {
		n, rerr := resp.Body.Read(buf)
		if n > 0 {
			done += int64(n)
			if done > maxBytes {
				_ = f.Close()
				return errors.New("file too large")
			}
			if _, werr := f.Write(buf[:n]); werr != nil {
				_ = f.Close()
				return werr
			}
			if onProgress != nil && time.Since(last) > 150*time.Millisecond {
				onProgress(done, resp.ContentLength)
				last = time.Now()
			}
		}
		if rerr == io.EOF {
			break
		}
		if rerr != nil {
			_ = f.Close()
			return rerr
		}
	}
	if onProgress != nil {
		onProgress(done, resp.ContentLength)
	}
	return f.Close()
}
