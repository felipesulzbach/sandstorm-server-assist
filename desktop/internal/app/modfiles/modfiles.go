// Package modfiles apaga um mod baixado para o servidor baixar de novo (ADR-037).
//
// O mod.io do servidor baixa os mods inscritos ao iniciar. Sem o mod no disco (pasta e registro no state.json), a
// próxima partida o baixa outra vez: é o jeito de consertar um mod corrompido ou forçar a versão mais nova. O
// servidor precisa estar parado, porque ele regrava o state.json ao fechar e desfaria a mudança.
package modfiles

import (
	"strconv"
	"sync"

	"sandstorm-server-assist/desktop/internal/domain/problems"
)

// Deps são as partes do sistema que o serviço usa (dublês nos testes).
type Deps struct {
	// Running diz se o servidor está rodando.
	Running func() bool
	// Busy diz se há outra tarefa que não pode rodar junto (instalação do servidor).
	Busy        func() bool
	StopAndWait func() error
	Start       func() error
	// Remove apaga a pasta do mod e o registro dele no state.json.
	Remove func(id int) error
	// RemoveAll apaga todos os mods baixados; Logout apaga a sessão do mod.io do servidor (ADR-047). Podem ser nil.
	RemoveAll func() error
	Logout    func() error
	Log       func(text, level string)
}

// Service faz um pedido por vez.
type Service struct {
	d    Deps
	mu   sync.Mutex
	busy bool
}

// New cria o serviço.
func New(d Deps) *Service { return &Service{d: d} }

// Delete apaga o mod. redownload=true inicia o servidor logo depois, e ele baixa o mod de novo (se o mod continua
// inscrito). Com o servidor rodando, só "apagar e baixar de novo" é aceito: o servidor para, o mod é apagado e o
// servidor volta. Se apagar falhar depois de parar, o servidor volta como estava.
func (s *Service) Delete(id int, redownload bool) error {
	if id <= 0 {
		return problems.Field("mod", "validation.modId", "value", strconv.Itoa(id))
	}
	if !s.lock() {
		return problems.New(problems.Busy, "")
	}
	defer s.unlock()
	if s.d.Busy != nil && s.d.Busy() {
		return problems.New(problems.Busy, "")
	}
	running := s.d.Running()
	if running && !redownload {
		return problems.New(problems.ServerRunning, "")
	}
	name := "mod " + strconv.Itoa(id)
	if running {
		s.d.Log("Stopping the server to delete "+name, "info")
		if err := s.d.StopAndWait(); err != nil {
			return err
		}
	}
	if err := s.d.Remove(id); err != nil {
		s.d.Log("Could not delete "+name+": "+err.Error(), "error")
		if running {
			_ = s.d.Start()
		}
		return problems.Wrap(problems.ModDelete, err)
	}
	s.d.Log("Deleted "+name+"; the server downloads it again on its next start if it is still subscribed", "info")
	if redownload {
		s.d.Log("Starting the server to download "+name+" again", "info")
		return s.d.Start()
	}
	return nil
}

// DownloadNew baixa os mods inscritos que ainda não estão no disco (ADR-038). Quem baixa é o mod.io do servidor, e
// só ao iniciar: com o servidor parado, ele inicia; rodando, ele reinicia (o mod novo não entra numa partida em
// andamento). Depois do download o runner reinicia uma vez sozinho para montar todos os mods (ADR-017).
func (s *Service) DownloadNew() error {
	if !s.lock() {
		return problems.New(problems.Busy, "")
	}
	defer s.unlock()
	if s.d.Busy != nil && s.d.Busy() {
		return problems.New(problems.Busy, "")
	}
	if s.d.Running() {
		s.d.Log("Restarting the server to download the new mods", "info")
		if err := s.d.StopAndWait(); err != nil {
			return err
		}
	} else {
		s.d.Log("Starting the server to download the new mods", "info")
	}
	return s.d.Start()
}

// Clear apaga, para começar do zero (ADR-047), os mods baixados (mods) e a autorização do servidor no mod.io
// (session). O servidor precisa estar parado: ele regrava o state.json ao fechar.
func (s *Service) Clear(mods, session bool) error {
	if !mods && !session {
		return nil
	}
	if !s.lock() {
		return problems.New(problems.Busy, "")
	}
	defer s.unlock()
	if s.d.Busy != nil && s.d.Busy() {
		return problems.New(problems.Busy, "")
	}
	if s.d.Running() {
		return problems.New(problems.ServerRunning, "")
	}
	if mods && s.d.RemoveAll != nil {
		if err := s.d.RemoveAll(); err != nil {
			s.d.Log("Could not delete every downloaded mod: "+err.Error(), "error")
			return problems.Wrap(problems.ModDelete, err)
		}
		s.d.Log("Deleted every downloaded mod (start from scratch)", "info")
	}
	if session && s.d.Logout != nil {
		if err := s.d.Logout(); err != nil {
			return problems.Wrap(problems.ModDelete, err)
		}
		s.d.Log("Removed the server's mod.io authorization (start from scratch)", "info")
	}
	return nil
}

func (s *Service) lock() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.busy {
		return false
	}
	s.busy = true
	return true
}

func (s *Service) unlock() {
	s.mu.Lock()
	s.busy = false
	s.mu.Unlock()
}
