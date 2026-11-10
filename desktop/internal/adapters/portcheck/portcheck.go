// Package portcheck pede ao verificador de portas do projeto (ADR-021) o teste das portas desta máquina pela
// internet. Só HTTPS para o host permitido (netx); vai só o número das portas.
package portcheck

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"sandstorm-server-assist/desktop/internal/adapters/netx"
	"sandstorm-server-assist/desktop/internal/domain/netcheck"
	"sandstorm-server-assist/desktop/internal/domain/problems"
)

// URL é o endereço do verificador (servidor da página de download).
const URL = "https://game-in.duckdns.org/sandstorm/api/check"

// timeout cobre os testes do verificador (alguns segundos por porta).
const timeout = 30 * time.Second

// Check testa as portas e devolve o resultado.
func Check(ctx context.Context, req netcheck.Request, gamePort int) (netcheck.Result, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	body, err := json.Marshal(req)
	if err != nil {
		return netcheck.Result{}, problems.Wrap(problems.PortCheckFailed, err)
	}
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, URL, bytes.NewReader(body))
	if err != nil {
		return netcheck.Result{}, problems.Wrap(problems.PortCheckFailed, err)
	}
	r.Header.Set("Content-Type", "application/json")
	resp, err := netx.Do(r)
	if err != nil {
		return netcheck.Result{}, netx.NetProblem(problems.PortCheckFailed, err)
	}
	defer func() { _ = resp.Body.Close() }()
	b, err := io.ReadAll(io.LimitReader(resp.Body, netcheck.MaxReply+1))
	if err != nil || len(b) > netcheck.MaxReply {
		return netcheck.Result{}, problems.New(problems.PortCheckFailed, "invalid reply")
	}
	return netcheck.Parse(resp.StatusCode, b, req, gamePort)
}
