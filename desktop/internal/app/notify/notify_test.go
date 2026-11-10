package notify

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"sandstorm-server-assist/desktop/internal/domain/server"
)

const hook = "https://discord.com/api/webhooks/123456789012/abcdefghijklmnopqrstuvwxyz_-ABC"

func TestMessages(t *testing.T) {
	var posted []string
	s := New(Deps{
		Config: func() (server.Profile, server.Secrets, error) {
			p := server.Default()
			p.Hostname = "Coop BR"
			p.Discord = server.Discord{Ready: "{server} online ({map})", Join: "{name} joined"}
			return p, server.Secrets{DiscordWebhook: hook}, nil
		},
		Post: func(_ context.Context, url string, body []byte) error {
			var m map[string]any
			_ = json.Unmarshal(body, &m)
			posted = append(posted, url+" "+m["content"].(string))
			return nil
		},
		Log: func(string, string) {},
	}, false)
	if _, text := s.Message(Ready, "Farmhouse", ""); text != "Coop BR online (Farmhouse)" {
		t.Fatal(text)
	}
	if _, text := s.Message(Stopped, "", ""); text != "" {
		t.Fatal("empty template must not send")
	}
	if _, text := s.Message(Join, "", "Fulano"); text != "Fulano joined" {
		t.Fatal(text)
	}
	if err := s.Test("Farmhouse"); err != nil || len(posted) != 1 || !strings.HasPrefix(posted[0], hook) {
		t.Fatal(posted, err)
	}
	body, _ := Payload("x")
	if !strings.Contains(string(body), `"parse":[]`) {
		t.Fatal(string(body))
	}
}
