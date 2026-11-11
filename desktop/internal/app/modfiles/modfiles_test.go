package modfiles

import (
	"errors"
	"reflect"
	"testing"

	"sandstorm-server-assist/desktop/internal/domain/problems"
)

type rig struct {
	running, busy bool
	removeErr     error
	calls         []string
}

func (r *rig) service() *Service {
	return New(Deps{
		Running: func() bool { return r.running },
		Busy:    func() bool { return r.busy },
		StopAndWait: func() error {
			r.calls = append(r.calls, "stop")
			r.running = false
			return nil
		},
		Start: func() error {
			r.calls = append(r.calls, "start")
			r.running = true
			return nil
		},
		Remove: func(id int) error {
			r.calls = append(r.calls, "remove")
			return r.removeErr
		},
		Log: func(string, string) {},
	})
}

func code(err error) string {
	var pe *problems.Error
	if errors.As(err, &pe) {
		return pe.Code
	}
	return ""
}

func TestDeleteWithServerStopped(t *testing.T) {
	r := &rig{}
	if err := r.service().Delete(42, false); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(r.calls, []string{"remove"}) {
		t.Fatalf("só apagar: %v", r.calls)
	}
	r.calls = nil
	if err := r.service().Delete(42, true); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(r.calls, []string{"remove", "start"}) {
		t.Fatalf("apagar e baixar: %v", r.calls)
	}
}

// Com o servidor rodando: para, apaga e inicia (o servidor regravaria o state.json ao fechar).
func TestDeleteWithServerRunning(t *testing.T) {
	r := &rig{running: true}
	if err := r.service().Delete(42, false); code(err) != problems.ServerRunning || r.calls != nil {
		t.Fatalf("só apagar com o servidor rodando: %v %v", err, r.calls)
	}
	if err := r.service().Delete(42, true); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(r.calls, []string{"stop", "remove", "start"}) {
		t.Fatalf("%v", r.calls)
	}
}

// Falha ao apagar depois de parar: o servidor volta como estava.
func TestDeleteFailureRestoresServer(t *testing.T) {
	r := &rig{running: true, removeErr: errors.New("in use")}
	if err := r.service().Delete(42, true); code(err) != problems.ModDelete {
		t.Fatalf("%v", err)
	}
	if !reflect.DeepEqual(r.calls, []string{"stop", "remove", "start"}) {
		t.Fatalf("%v", r.calls)
	}
}

func TestDeleteRefusals(t *testing.T) {
	r := &rig{busy: true}
	if err := r.service().Delete(42, true); code(err) != problems.Busy || r.calls != nil {
		t.Fatalf("instalação em andamento: %v %v", err, r.calls)
	}
	var pe *problems.Error
	if err := (&rig{}).service().Delete(0, true); !errors.As(err, &pe) || pe.Msg != "validation.modId" {
		t.Fatalf("id 0: %v", err)
	}
}

// Mods inscritos e ainda não baixados: o servidor baixa ao iniciar (parado: inicia; rodando: reinicia).
func TestDownloadNew(t *testing.T) {
	r := &rig{}
	if err := r.service().DownloadNew(); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(r.calls, []string{"start"}) {
		t.Fatalf("parado: %v", r.calls)
	}
	r = &rig{running: true}
	if err := r.service().DownloadNew(); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(r.calls, []string{"stop", "start"}) {
		t.Fatalf("rodando: %v", r.calls)
	}
	r = &rig{busy: true}
	if err := r.service().DownloadNew(); code(err) != problems.Busy || r.calls != nil {
		t.Fatalf("instalação em andamento: %v %v", err, r.calls)
	}
}
