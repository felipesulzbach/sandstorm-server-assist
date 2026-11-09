package updater

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"sync"
	"testing"

	"sandstorm-server-assist/desktop/internal/domain/problems"
	"sandstorm-server-assist/desktop/internal/domain/update"
)

type fakeFetcher struct {
	manifest []byte
	err      error
	sha      string
	size     int64
	dlErr    error
	calls    int
}

func (f *fakeFetcher) Manifest(context.Context) ([]byte, error) { return f.manifest, f.err }
func (f *fakeFetcher) Download(_ context.Context, _, _ string, onProgress func(int64)) (string, int64, error) {
	f.calls++
	onProgress(f.size / 2)
	return f.sha, f.size, f.dlErr
}

type fakeFiles struct{ renamed, removed []string }

func (f *fakeFiles) Rename(from, _ string) error { f.renamed = append(f.renamed, from); return nil }
func (f *fakeFiles) Remove(p string) error       { f.removed = append(f.removed, p); return nil }

type fakeLauncher struct {
	mu   sync.Mutex
	runs [][]string
}

func (l *fakeLauncher) Start(path string, args ...string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.runs = append(l.runs, append([]string{path}, args...))
	return nil
}

type env struct {
	svc      *Service
	fetch    *fakeFetcher
	files    *fakeFiles
	launcher *fakeLauncher
}

const sha = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

func setup(t *testing.T, current string, installed bool, mutate func(*update.Manifest), signWith ed25519.PrivateKey) *env {
	t.Helper()
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)
	if signWith == nil {
		signWith = priv
	}
	m := update.Manifest{Version: "0.2.0", File: update.InstallerName("0.2.0"), SHA256: sha, Size: 1000, Date: "2026-11-09"}
	if mutate != nil {
		mutate(&m)
	}
	m.Signature = base64.StdEncoding.EncodeToString(ed25519.Sign(signWith, m.Message()))
	raw, _ := json.Marshal(m)
	e := &env{fetch: &fakeFetcher{manifest: raw, sha: sha, size: 1000}, files: &fakeFiles{}, launcher: &fakeLauncher{}}
	e.svc = New(Config{
		Current: current, Installed: installed, Dir: `C:\cache\updates`, PublicKey: pub,
		Fetcher: e.fetch, Files: e.files, Launcher: e.launcher,
	})
	return e
}

func TestDownloadsAndInstalls(t *testing.T) {
	e := setup(t, "0.1.0", true, nil, nil)
	st := e.svc.Check(context.Background())
	if st.Status != StatusReady || st.Version != "0.2.0" || st.Current != "0.1.0" {
		t.Fatalf("state %+v", st)
	}
	if len(e.files.renamed) != 1 {
		t.Fatal("installer not moved into place")
	}
	// Pronta: não procura de novo.
	e.svc.Check(context.Background())
	if e.fetch.calls != 1 {
		t.Fatal("downloaded twice")
	}
	if err := e.svc.Install(); err != nil {
		t.Fatal(err)
	}
	e.svc.InstallOnExit() // já lançado: não roda de novo
	if len(e.launcher.runs) != 1 {
		t.Fatalf("runs %v", e.launcher.runs)
	}
	run := e.launcher.runs[0]
	if run[0] != `C:\cache\updates\sandstorm-server-assist-0.2.0-windows-setup.exe` || run[1] != "/S" || run[2] != "/relaunch" {
		t.Fatalf("run %v", run)
	}
}

func TestOnExitDoesNotRelaunch(t *testing.T) {
	e := setup(t, "0.1.0", true, nil, nil)
	e.svc.Check(context.Background())
	e.svc.InstallOnExit()
	if len(e.launcher.runs) != 1 || len(e.launcher.runs[0]) != 2 {
		t.Fatalf("runs %v", e.launcher.runs)
	}
}

func TestRefusesBadSignatureAndHash(t *testing.T) {
	_, other, _ := ed25519.GenerateKey(rand.Reader)
	e := setup(t, "0.1.0", true, nil, other)
	if st := e.svc.Check(context.Background()); st.Status != StatusError || st.Stage != StageCheck {
		t.Fatalf("bad signature: %+v", st)
	}
	if e.fetch.calls != 0 {
		t.Fatal("downloaded with a bad signature")
	}

	e = setup(t, "0.1.0", true, nil, nil)
	e.fetch.sha = "ff" + sha[2:]
	if st := e.svc.Check(context.Background()); st.Status != StatusError || st.Stage != StageDownload {
		t.Fatalf("bad hash: %+v", st)
	}
	if len(e.files.renamed) != 0 || len(e.files.removed) == 0 {
		t.Fatal("bad download kept")
	}
	if err := e.svc.Install(); problems.CodeOf(err) != problems.UpdateFailed {
		t.Fatalf("install without ready update: %v", err)
	}
}

func TestStatesWithoutDownload(t *testing.T) {
	if st := setup(t, "0.2.0", true, nil, nil).svc.Check(context.Background()); st.Status != StatusUpToDate {
		t.Fatalf("same version: %+v", st)
	}
	if st := setup(t, "0.1.0", false, nil, nil).svc.Check(context.Background()); st.Status != StatusManual || st.Version != "0.2.0" {
		t.Fatalf("portable: %+v", st)
	}
	e := setup(t, "0.1.0-dev", true, nil, nil)
	if st := e.svc.Check(context.Background()); st.Status != StatusDev || e.fetch.calls != 0 {
		t.Fatalf("dev: %+v", st)
	}
	e = setup(t, "0.1.0", true, nil, nil)
	e.fetch.err = errors.New("offline")
	if st := e.svc.Check(context.Background()); st.Status != StatusError || st.Stage != StageCheck {
		t.Fatalf("offline: %+v", st)
	}
	if New(Config{Current: "0.1.0"}).State().Status != StatusDev {
		t.Fatal("without fetcher should be off")
	}
}
