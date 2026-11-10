package config

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"unicode"
	"unicode/utf8"

	"sandstorm-server-assist/desktop/internal/domain/problems"
	"sandstorm-server-assist/desktop/internal/domain/server"
)

// Vários perfis (ADR-022): o perfil em uso fica em file.Profile (o resto do app só conhece ele); os outros ficam
// guardados em file.Others. Trocar de perfil troca as configurações do servidor (pasta, portas, mapas, regras). Os
// segredos do cofre são os mesmos para todos.

// NamedProfile é um perfil guardado.
type NamedProfile struct {
	ID      string         `json:"id"`
	Name    string         `json:"name"`
	Profile server.Profile `json:"profile"`
}

// ProfileRef é o que a tela mostra de um perfil.
type ProfileRef struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Hostname string `json:"hostname"`
	Active   bool   `json:"active"`
}

// MaxProfiles limita quantos perfis o app guarda.
const MaxProfiles = 20

func newID() string {
	b := make([]byte, 6)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// cleanProfileName confere o nome (1 a 40 caracteres, sem controle).
func cleanProfileName(name string) (string, error) {
	name = strings.Join(strings.Fields(name), " ")
	if name == "" || utf8.RuneCountInString(name) > 40 || strings.IndexFunc(name, unicode.IsControl) >= 0 {
		return "", problems.Field("profileName", "validation.profileName")
	}
	return name, nil
}

// ensureActive dá id e nome ao perfil em uso (configurações de versões antigas não têm). Chamar com o lock.
func (s *Service) ensureActive() {
	if s.data.ActiveID == "" {
		s.data.ActiveID = newID()
	}
	if s.data.ActiveName == "" {
		s.data.ActiveName = s.data.Profile.Hostname
	}
}

// Profiles lista os perfis (o em uso primeiro).
func (s *Service) Profiles() []ProfileRef {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ensureActive()
	out := []ProfileRef{{ID: s.data.ActiveID, Name: s.data.ActiveName, Hostname: s.data.Profile.Hostname, Active: true}}
	for _, o := range s.data.Others {
		out = append(out, ProfileRef{ID: o.ID, Name: o.Name, Hostname: o.Profile.Hostname})
	}
	return out
}

// CreateProfile guarda um perfil novo: cópia do em uso (copy) ou o padrão com as mesmas pastas. Não troca o perfil
// em uso.
func (s *Service) CreateProfile(name string, copyCurrent bool) ([]ProfileRef, error) {
	name, err := cleanProfileName(name)
	if err != nil {
		return s.Profiles(), err
	}
	s.mu.Lock()
	s.ensureActive()
	if len(s.data.Others)+1 >= MaxProfiles {
		s.mu.Unlock()
		return s.Profiles(), problems.Field("profileName", "validation.profilesMax", "max", "20")
	}
	p := server.Default()
	if copyCurrent {
		p = s.data.Profile
	}
	p.ServerDir, p.SteamCMDDir = s.data.Profile.ServerDir, s.data.Profile.SteamCMDDir
	s.data.Others = append(s.data.Others, NamedProfile{ID: newID(), Name: name, Profile: p})
	err = s.save()
	s.mu.Unlock()
	return s.Profiles(), err
}

// RenameProfile troca o nome de um perfil.
func (s *Service) RenameProfile(id, name string) ([]ProfileRef, error) {
	name, err := cleanProfileName(name)
	if err != nil {
		return s.Profiles(), err
	}
	s.mu.Lock()
	s.ensureActive()
	found := false
	if id == s.data.ActiveID {
		s.data.ActiveName, found = name, true
	}
	for i := range s.data.Others {
		if s.data.Others[i].ID == id {
			s.data.Others[i].Name, found = name, true
		}
	}
	if !found {
		s.mu.Unlock()
		return s.Profiles(), problems.Field("profile", "validation.unknownField")
	}
	err = s.save()
	s.mu.Unlock()
	return s.Profiles(), err
}

// DeleteProfile apaga um perfil guardado (o em uso não pode ser apagado).
func (s *Service) DeleteProfile(id string) ([]ProfileRef, error) {
	s.mu.Lock()
	idx := -1
	for i, o := range s.data.Others {
		if o.ID == id {
			idx = i
		}
	}
	if idx < 0 {
		s.mu.Unlock()
		return s.Profiles(), problems.Field("profile", "validation.unknownField")
	}
	s.data.Others = append(s.data.Others[:idx], s.data.Others[idx+1:]...)
	err := s.save()
	s.mu.Unlock()
	return s.Profiles(), err
}

// SwitchProfile passa a usar outro perfil (quem chama confere que o servidor está parado).
func (s *Service) SwitchProfile(id string) (server.Profile, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ensureActive()
	for i, o := range s.data.Others {
		if o.ID != id {
			continue
		}
		s.data.Others[i] = NamedProfile{ID: s.data.ActiveID, Name: s.data.ActiveName, Profile: s.data.Profile}
		s.data.ActiveID, s.data.ActiveName, s.data.Profile = o.ID, o.Name, o.Profile
		return s.data.Profile, s.save()
	}
	return s.data.Profile, problems.Field("profile", "validation.unknownField")
}
