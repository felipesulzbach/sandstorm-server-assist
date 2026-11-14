package config

import (
	"crypto/rand"
	"encoding/hex"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"sandstorm-server-assist/desktop/internal/domain/problems"
	"sandstorm-server-assist/desktop/internal/domain/server"
)

// Vários servidores (ADR-048, antes perfis do ADR-022). O arquivo continua no formato dos perfis, para uma versão
// antiga ainda abrir: o servidor escolhido na tela fica em file.Profile (ActiveID, ActiveName) e os outros em
// file.Others. Escolher outro servidor só troca qual deles as telas editam; todos podem rodar ao mesmo tempo, cada
// um com o seu número de instância. A lista para a tela vem em ordem de instância (não muda ao escolher).

// NamedProfile é um servidor guardado.
type NamedProfile struct {
	ID      string         `json:"id"`
	Name    string         `json:"name"`
	Profile server.Profile `json:"profile"`
}

// ProfileRef é o que a tela mostra de um servidor.
type ProfileRef struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Hostname string `json:"hostname"`
	Active   bool   `json:"active"`
	Instance int    `json:"instance"`
	Port     int    `json:"port"`
}

// MaxProfiles limita quantos servidores o app guarda.
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

// ensureActive dá id e nome ao servidor escolhido (configurações de versões antigas não têm). Chamar com o lock.
func (s *Service) ensureActive() bool {
	changed := false
	if s.data.ActiveID == "" {
		s.data.ActiveID, changed = newID(), true
	}
	if s.data.ActiveName == "" {
		s.data.ActiveName, changed = s.data.Profile.Hostname, true
	}
	return changed
}

// ensureInstances dá a cada servidor um número de instância só dele: perfis antigos têm 0 (o escolhido vira o
// principal, 1, e os outros 2, 3...). Chamar com o lock.
func (s *Service) ensureInstances() bool {
	used := map[int]bool{}
	var fix []*server.Profile
	s.eachProfile(func(p *server.Profile) {
		if p.Instance >= 1 && p.Instance <= server.MaxInstance && !used[p.Instance] {
			used[p.Instance] = true
		} else {
			fix = append(fix, p)
		}
	})
	for _, p := range fix {
		taken := make([]int, 0, len(used))
		for n := range used {
			taken = append(taken, n)
		}
		p.Instance = server.NextInstance(taken)
		used[p.Instance] = true
	}
	return len(fix) > 0
}

// eachProfile chama fn com cada servidor (o escolhido primeiro). Chamar com o lock.
func (s *Service) eachProfile(fn func(*server.Profile)) {
	fn(&s.data.Profile)
	for i := range s.data.Others {
		fn(&s.data.Others[i].Profile)
	}
}

// find devolve o perfil de um servidor pelo id (nil se não existe). Chamar com o lock.
func (s *Service) find(id string) *server.Profile {
	if id == "" {
		return nil
	}
	if id == s.data.ActiveID {
		return &s.data.Profile
	}
	for i := range s.data.Others {
		if s.data.Others[i].ID == id {
			return &s.data.Others[i].Profile
		}
	}
	return nil
}

// all devolve todos os servidores (cópias), o escolhido primeiro. Chamar com o lock.
func (s *Service) all() []NamedProfile {
	out := []NamedProfile{{ID: s.data.ActiveID, Name: s.data.ActiveName, Profile: s.data.Profile}}
	return append(out, s.data.Others...)
}

// ---- segredos por servidor ----

// secretsOf são os segredos de um servidor com os da conta do mod.io. Chamar com o lock.
func (s *Service) secretsOf(id string) server.Secrets {
	return s.vdata.Servers[id].WithAccount(s.vdata.Secrets)
}

// copyVault devolve uma cópia do cofre para mudar e gravar. Chamar com o lock.
func (s *Service) copyVault() server.Vault {
	next := server.Vault{Secrets: s.vdata.Secrets, Version: s.vdata.Version, Servers: map[string]server.Secrets{}}
	for id, sec := range s.vdata.Servers {
		next.Servers[id] = sec
	}
	return next
}

// saveVault grava o cofre e passa a usá-lo. O nível de cima leva a conta do mod.io e uma cópia dos segredos do
// servidor escolhido (uma versão antiga do app só lê esse nível). Chamar com o lock.
func (s *Service) saveVault(next server.Vault) error {
	next.Version = server.VaultVersion
	next.Secrets = next.Servers[s.data.ActiveID].WithAccount(next.Secrets)
	if err := s.vault.Save(next); err != nil {
		return err
	}
	s.vdata = next
	return nil
}

// migrateVault: no cofre de versões antigas todos os segredos ficam no nível de cima e valiam para todos os perfis.
// Os do servidor (senhas, tokens, webhook) passam para cada servidor; os da conta do mod.io ficam em cima.
// Chamar com o lock (ou antes de o serviço ser usado).
func (s *Service) migrateVault() {
	if s.vdata.Version >= server.VaultVersion {
		return
	}
	next := s.copyVault()
	// Só preenche quem ainda não tem segredos próprios: nunca troca o que já está guardado por servidor.
	for _, np := range s.all() {
		if _, ok := next.Servers[np.ID]; !ok {
			next.Servers[np.ID] = s.vdata.Secrets.ServerOnly()
		}
	}
	if err := s.saveVault(next); err != nil {
		// Sem gravar, vale na memória: os segredos de antes continuam servindo a todos os servidores.
		next.Version = server.VaultVersion
		s.vdata = next
	}
}

// ---- lista de servidores ----

// Profiles lista os servidores em ordem de instância.
func (s *Service) Profiles() []ProfileRef {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ensureActive()
	out := []ProfileRef{}
	for _, np := range s.all() {
		out = append(out, ProfileRef{ID: np.ID, Name: np.Name, Hostname: np.Profile.Hostname, Active: np.ID == s.data.ActiveID,
			Instance: server.InstanceNumber(np.Profile.Instance), Port: np.Profile.Port})
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Instance < out[j].Instance })
	return out
}

// SelectedID é o id do servidor que as telas editam.
func (s *Service) SelectedID() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.data.ActiveID
}

// Name devolve o nome de um servidor ("" se não existe).
func (s *Service) Name(id string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, np := range s.all() {
		if np.ID == id {
			return np.Name
		}
	}
	return ""
}

// Peers devolve os outros servidores (todos menos id), com nome, perfil e segredos (para conferir portas e GSLT).
func (s *Service) Peers(id string) []server.Peer {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []server.Peer
	for _, np := range s.all() {
		if np.ID != id {
			out = append(out, server.Peer{ID: np.ID, Name: np.Name, Profile: np.Profile, Secrets: s.secretsOf(np.ID)})
		}
	}
	return out
}

// CreateProfile guarda um servidor novo e devolve o id: cópia do escolhido (copy, sem os tokens GSLT e GameStats,
// que a Steam aceita em um servidor por vez) ou o padrão com as mesmas pastas. Ganha a próxima instância livre e
// portas que nenhum outro servidor usa. Não troca o servidor escolhido.
// busy confere as portas ocupadas por outros programas do PC (pode ser nil).
func (s *Service) CreateProfile(name string, copyCurrent bool, busy server.PortBusy) ([]ProfileRef, string, error) {
	name, err := cleanProfileName(name)
	if err != nil {
		return s.Profiles(), "", err
	}
	s.mu.Lock()
	s.ensureActive()
	if len(s.data.Others)+1 >= MaxProfiles {
		s.mu.Unlock()
		return s.Profiles(), "", problems.Field("profileName", "validation.profilesMax", "max", "20")
	}
	p := server.Default()
	if copyCurrent {
		p = s.data.Profile
		p.Admins = append([]server.Admin(nil), p.Admins...)
		p.MapCycle = append([]string(nil), p.MapCycle...)
		p.Mutators = append([]string(nil), p.Mutators...)
		p.Announcements = append([]string(nil), p.Announcements...)
		p.ModIDs = append([]int(nil), p.ModIDs...)
		if p.Rules != nil {
			rules := map[string]string{}
			for k, v := range p.Rules {
				rules[k] = v
			}
			p.Rules = rules
		}
	}
	p.ServerDir, p.SteamCMDDir = s.data.Profile.ServerDir, s.data.Profile.SteamCMDDir
	var others []server.Profile
	var used []int
	for _, np := range s.all() {
		others = append(others, np.Profile)
		used = append(used, np.Profile.Instance)
	}
	p.Instance = server.NextInstance(used)
	p.Port, p.QueryPort, p.RconPort = server.FreePorts(others, busy)
	id := newID()
	if copyCurrent && s.vaultErr == nil {
		sec := s.vdata.Servers[s.data.ActiveID]
		sec.GSLTToken, sec.GameStatsToken = "", ""
		if sec != (server.Secrets{}) {
			next := s.copyVault()
			next.Servers[id] = sec
			if err := s.saveVault(next); err != nil {
				s.mu.Unlock()
				return s.Profiles(), "", err
			}
		}
	}
	s.data.Others = append(s.data.Others, NamedProfile{ID: id, Name: name, Profile: p})
	err = s.save()
	s.mu.Unlock()
	return s.Profiles(), id, err
}

// RenameProfile troca o nome de um servidor.
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

// DeleteProfile apaga um servidor que não é o escolhido, com os segredos dele (quem chama confere que está parado).
// Os arquivos dele na pasta do servidor ficam.
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
	if _, ok := s.vdata.Servers[id]; ok && s.vaultErr == nil {
		next := s.copyVault()
		delete(next.Servers, id)
		if err := s.saveVault(next); err != nil {
			s.mu.Unlock()
			return s.Profiles(), err
		}
	}
	s.data.Others = append(s.data.Others[:idx], s.data.Others[idx+1:]...)
	err := s.save()
	s.mu.Unlock()
	return s.Profiles(), err
}

// SwitchProfile escolhe o servidor que as telas editam. Pode ser com servidores rodando: cada um continua com o
// próprio perfil.
func (s *Service) SwitchProfile(id string) (server.Profile, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ensureActive()
	if id == s.data.ActiveID {
		return s.data.Profile, nil
	}
	for i, o := range s.data.Others {
		if o.ID != id {
			continue
		}
		s.data.Others[i] = NamedProfile{ID: s.data.ActiveID, Name: s.data.ActiveName, Profile: s.data.Profile}
		s.data.ActiveID, s.data.ActiveName, s.data.Profile = o.ID, o.Name, o.Profile
		if s.vaultErr == nil {
			_ = s.saveVault(s.copyVault()) // a cópia para versões antigas passa a ser a do novo escolhido
		}
		return s.data.Profile, s.save()
	}
	return s.data.Profile, problems.Field("profile", "validation.unknownField")
}
