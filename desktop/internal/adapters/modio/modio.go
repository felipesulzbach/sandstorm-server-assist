// Package modio fala com a API do mod.io (só HTTPS, endereço fixo, ADR-010) e lê os arquivos locais do mod.io
// (sessão do servidor em user.json e mods instalados em state.json).
package modio

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"sandstorm-server-assist/desktop/internal/adapters/fsx"
	"sandstorm-server-assist/desktop/internal/adapters/netx"
	"sandstorm-server-assist/desktop/internal/domain/problems"
)

// GameID é o Insurgency: Sandstorm no mod.io.
const GameID = 254

// BaseURL é a API do mod.io do jogo.
const BaseURL = "https://g-254.modapi.io/v1"

// Mod é o que a tela mostra de um mod.
type Mod struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	ProfileURL string `json:"profileUrl"`
	// Size é o tamanho final instalado (filesize_uncompressed).
	Size int64 `json:"size"`
}

type apiMod struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	ProfileURL string `json:"profile_url"`
	Modfile    struct {
		Filesize             int64  `json:"filesize"`
		FilesizeUncompressed int64  `json:"filesize_uncompressed"`
		MetadataBlob         string `json:"metadata_blob"`
	} `json:"modfile"`
	Tags []struct {
		Name string `json:"name"`
	} `json:"tags"`
}

func (m apiMod) toMod() Mod {
	size := m.Modfile.FilesizeUncompressed
	if size == 0 {
		size = m.Modfile.Filesize
	}
	return Mod{ID: m.ID, Name: m.Name, ProfileURL: m.ProfileURL, Size: size}
}

// Client chama a API.
type Client struct{ BaseURL string }

// New usa o endereço oficial.
func New() Client { return Client{BaseURL: BaseURL} }

// RequestCode pede ao mod.io que envie o código de 5 caracteres para o e-mail (mesmo pedido do .bat antigo).
func (c Client) RequestCode(ctx context.Context, apiKey, email string) error {
	form := url.Values{"email": {email}}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/oauth/emailrequest?api_key="+url.QueryEscape(apiKey), strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	_, err = c.send(req)
	return err
}

// Subscribed lista os mods em que a conta do servidor está inscrita (token da sessão do user.json).
func (c Client) Subscribed(ctx context.Context, token string) ([]Mod, error) {
	var out []Mod
	for offset := 0; offset < 1000; offset += 100 {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet,
			fmt.Sprintf("%s/me/subscribed?game_id=%d&_limit=100&_offset=%d", c.BaseURL, GameID, offset), nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/json")
		body, err := c.send(req)
		if err != nil {
			return nil, err
		}
		var page struct {
			Data        []apiMod `json:"data"`
			ResultTotal int      `json:"result_total"`
		}
		if err := json.Unmarshal(body, &page); err != nil {
			return nil, problems.New(problems.ModioRequestFailed, "invalid response")
		}
		for _, m := range page.Data {
			out = append(out, m.toMod())
		}
		if len(page.Data) < 100 || len(out) >= page.ResultTotal {
			break
		}
	}
	return out, nil
}

func (c Client) send(req *http.Request) ([]byte, error) {
	resp, err := netx.Do(req)
	if err != nil {
		return nil, netx.NetProblem(problems.ModioRequestFailed, err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, problems.New(problems.ModioNotAuthorized, fmt.Sprintf("HTTP %d", resp.StatusCode))
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		var e struct {
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		_ = json.Unmarshal(body, &e)
		return nil, problems.New(problems.ModioRequestFailed, fmt.Sprintf("HTTP %d %s", resp.StatusCode, e.Error.Message))
	}
	return body, nil
}

// ---- arquivos locais ----

// Session é a sessão do servidor no mod.io (sem o token, que nunca sai do Go).
type Session struct {
	Found    bool      `json:"found"`
	Username string    `json:"username"`
	Expires  time.Time `json:"expires"`
	Valid    bool      `json:"valid"`
}

type userFile struct {
	OAuth struct {
		Expiry int64  `json:"expiry"`
		Status int    `json:"status"`
		Token  string `json:"token"`
	} `json:"OAuth"`
	Profile struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
	} `json:"Profile"`
}

// UserFile é o user.json do servidor dedicado (%LOCALAPPDATA%\mod.io\254\ModServer\user.json).
func UserFile() (string, error) {
	base, err := os.UserCacheDir() // %LOCALAPPDATA% no Windows
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "mod.io", "254", "ModServer", "user.json"), nil
}

// ReadSession lê a sessão e devolve também o token (só para chamar a API, nunca para a tela).
func ReadSession(path string) (Session, string, error) {
	b, err := fsx.ReadLimited(path, 1<<20)
	if err != nil {
		if os.IsNotExist(err) {
			return Session{}, "", nil
		}
		return Session{}, "", err
	}
	var u userFile
	if err := json.Unmarshal(b, &u); err != nil {
		return Session{Found: true}, "", nil
	}
	s := Session{Found: true, Username: u.Profile.Username}
	if u.OAuth.Expiry > 0 {
		s.Expires = time.Unix(u.OAuth.Expiry, 0)
	}
	s.Valid = u.OAuth.Token != "" && u.OAuth.Status == 0 && s.Expires.After(time.Now())
	return s, u.OAuth.Token, nil
}

// RootDir é a pasta dos mods (globalsettings.json → RootLocalStoragePath, padrão C:\Users\Public\mod.io).
func RootDir() string {
	def := filepath.Join(os.Getenv("PUBLIC"), "mod.io")
	if os.Getenv("PUBLIC") == "" {
		def = `C:\Users\Public\mod.io`
	}
	base, err := os.UserCacheDir()
	if err != nil {
		return def
	}
	b, err := fsx.ReadLimited(filepath.Join(base, "mod.io", "globalsettings.json"), 64<<10)
	if err != nil {
		return def
	}
	var g struct {
		RootLocalStoragePath string `json:"RootLocalStoragePath"`
	}
	if json.Unmarshal(b, &g) != nil || g.RootLocalStoragePath == "" || !filepath.IsAbs(g.RootLocalStoragePath) {
		return def
	}
	return filepath.Clean(g.RootLocalStoragePath)
}

// Installed é um mod no disco (state.json).
type Installed struct {
	Mod
	Path       string `json:"path"`
	SizeOnDisk int64  `json:"sizeOnDisk"`
	State      int    `json:"state"`
	// Tags do mod no mod.io ("Mutator", "Checkpoint"...) e o metadata_blob do arquivo (PluginName:...).
	Tags []string `json:"-"`
	Blob string   `json:"-"`
	// Subscribers são as contas do mod.io inscritas no mod (state.json, "SubscriptionCount").
	Subscribers []int `json:"-"`
}

// ReadInstalled lê o state.json da pasta de mods.
// StateInstalled é o estado "instalado" de um mod no state.json do mod.io (os outros são download, pendente...).
const StateInstalled = 1

func ReadInstalled(root string) ([]Installed, error) {
	b, err := fsx.ReadLimited(filepath.Join(root, "254", "metadata", "state.json"), 64<<20)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var st struct {
		Mods []struct {
			ID         int    `json:"ID"`
			PathOnDisk string `json:"PathOnDisk"`
			SizeOnDisk int64  `json:"SizeOnDisk"`
			State      int    `json:"State"`
			Profile    apiMod `json:"Profile"`
			// Contas do mod.io inscritas no mod (a pasta é dividida entre o jogo e o servidor desta máquina).
			SubscriptionCount []int `json:"SubscriptionCount"`
		} `json:"Mods"`
	}
	if err := json.Unmarshal(b, &st); err != nil {
		return nil, err
	}
	out := make([]Installed, 0, len(st.Mods))
	for _, m := range st.Mods {
		mod := m.Profile.toMod()
		mod.ID = m.ID
		in := Installed{Mod: mod, Path: m.PathOnDisk, SizeOnDisk: m.SizeOnDisk, State: m.State,
			Blob: m.Profile.Modfile.MetadataBlob, Subscribers: m.SubscriptionCount}
		for _, t := range m.Profile.Tags {
			in.Tags = append(in.Tags, t.Name)
		}
		out = append(out, in)
	}
	return out, nil
}

// ModDir é a pasta de um mod.
func ModDir(root string, id int) string {
	return filepath.Join(root, "254", "mods", fmt.Sprint(id))
}

var apiKeyRe = regexp.MustCompile(`api_key=([0-9a-fA-F]{32})`)

// FindAPIKey procura a chave de API nas requisições que o próprio servidor registrou no log (ADR-010).
func FindAPIKey(logText string) string {
	if m := apiKeyRe.FindStringSubmatch(logText); m != nil {
		return m[1]
	}
	return ""
}

// ServerUserID é o id da conta do mod.io do servidor (user.json do ModServer); 0 se o servidor não foi autorizado.
// Só lê o id: o token continua só com ReadSession.
func ServerUserID() int {
	path, err := UserFile()
	if err != nil {
		return 0
	}
	b, err := fsx.ReadLimited(path, 1<<20)
	if err != nil {
		return 0
	}
	var u userFile
	if json.Unmarshal(b, &u) != nil {
		return 0
	}
	return u.Profile.ID
}

// ForServer deixa só os mods do servidor (ADR-043). A pasta de mods é dividida entre o jogo e o servidor da mesma
// máquina, e o state.json diz quais contas estão inscritas em cada mod. Fica o mod em que a conta do servidor está
// inscrita; sem a conta (servidor não autorizado), nenhum. Mod sem a lista de contas (versões antigas do mod.io)
// fica, como antes.
func ForServer(list []Installed, serverID int) []Installed {
	out := make([]Installed, 0, len(list))
	for _, m := range list {
		if len(m.Subscribers) == 0 {
			out = append(out, m)
			continue
		}
		for _, id := range m.Subscribers {
			if serverID != 0 && id == serverID {
				out = append(out, m)
				break
			}
		}
	}
	return out
}

// ServerInstalled são os mods instalados da conta do servidor (ForServer sobre ReadInstalled).
func ServerInstalled(root string) ([]Installed, error) {
	list, err := ReadInstalled(root)
	if err != nil {
		return nil, err
	}
	return ForServer(list, ServerUserID()), nil
}
