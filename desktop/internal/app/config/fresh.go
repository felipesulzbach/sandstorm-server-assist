package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"sandstorm-server-assist/desktop/internal/domain/server"
)

// Começar do zero (ADR-047). O perfil e o cofre são do app, não da pasta do servidor: mudar ou apagar a pasta não
// limpa nada, e o app não apaga sozinho (um disco desligado ou uma pasta renomeada também "somem"). Quando a pasta
// guardada não existe mais, ou não tem servidor instalado, e o perfil tem dados, a tela oferece começar do zero.

// Fresh é o que a tela precisa para oferecer começar do zero.
type Fresh struct {
	// Offer: a pasta do servidor não tem servidor instalado (ou não existe) e o perfil ou o cofre têm dados.
	Offer bool `json:"offer"`
	// Missing: a pasta guardada não existe mais (foi movida ou apagada).
	Missing bool `json:"missing"`
	// HasSecrets: o cofre tem alguma senha, token ou chave.
	HasSecrets bool `json:"hasSecrets"`
	// Others é quantos outros perfis existem (o cofre é o mesmo para todos).
	Others int `json:"others"`
	// NoServer: há uma pasta escolhida e nela não há servidor instalado.
	NoServer bool `json:"noServer"`
}

// serverExe é o executável que indica um servidor instalado na pasta.
const serverExe = "InsurgencyServer.exe"

// FreshState diz se a tela deve oferecer começar do zero.
func (s *Service) FreshState() Fresh {
	s.mu.Lock()
	p, sec, others := s.data.Profile, s.secrets, len(s.data.Others)
	s.mu.Unlock()
	f := Fresh{HasSecrets: sec != server.Secrets{}, Others: others}
	if p.ServerDir == "" {
		return f
	}
	if info, err := os.Stat(p.ServerDir); err != nil || !info.IsDir() {
		f.Missing = true
	}
	_, err := os.Stat(filepath.Join(p.ServerDir, serverExe))
	f.NoServer = err != nil
	f.Offer = f.NoServer && (f.HasSecrets || hasData(p))
	return f
}

// hasData diz se o perfil tem algo além do padrão (as pastas não contam).
func hasData(p server.Profile) bool {
	p.ServerDir, p.SteamCMDDir = "", ""
	a, errA := json.Marshal(p)
	b, errB := json.Marshal(server.Default())
	return errA != nil || errB != nil || string(a) != string(b)
}

// StartFresh volta o perfil em uso ao padrão, mantendo as pastas escolhidas, e com clearSecrets também esvazia o
// cofre (vale para todos os perfis). Os outros perfis guardados não mudam.
func (s *Service) StartFresh(clearSecrets bool) (server.Profile, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if clearSecrets {
		if s.vaultErr != nil {
			return s.data.Profile, s.vaultErr
		}
		if err := s.vault.Save(server.Secrets{}); err != nil {
			return s.data.Profile, err
		}
		s.secrets = server.Secrets{}
	}
	p := server.Default()
	p.ServerDir, p.SteamCMDDir = s.data.Profile.ServerDir, s.data.Profile.SteamCMDDir
	s.data.Profile = p
	return p, s.save()
}
