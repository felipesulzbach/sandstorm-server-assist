// Package discord manda o aviso para o webhook do Discord (ADR-022), pelo cliente HTTP do app (só HTTPS, só hosts
// da lista do netx).
package discord

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"sandstorm-server-assist/desktop/internal/adapters/netx"
	"sandstorm-server-assist/desktop/internal/domain/problems"
)

// Post manda o JSON. O endereço já foi validado no domínio (server.ValidateSecret) e passa de novo pela lista de
// hosts do netx.
func Post(ctx context.Context, endpoint string, body []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("discord: invalid request")
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := netx.Do(req)
	if err != nil {
		// O erro do Go traz o endereço inteiro, e o endereço do webhook é segredo: fica só a causa.
		var ue *url.Error
		if errors.As(err, &ue) {
			err = ue.Err
		}
		return netx.NetProblem(problems.DiscordFailed, err)
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 64<<10))
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("discord: HTTP %d", resp.StatusCode)
	}
	return nil
}
