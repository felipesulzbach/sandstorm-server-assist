package installer

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"sandstorm-server-assist/desktop/internal/domain/server"
)

const manifest = "\"AppState\"\n{\n\t\"appid\"\t\t\"581330\"\n\t\"LauncherPath\"\t\t\"%s\"\n\t\"SizeOnDisk\"\t\t\"4994664175\"\n\t\"buildid\"\t\t\"25534773\"\n\t\"LastUpdated\"\t\t\"1794188468\"\n\t\"BytesToDownload\"\t\t\"100\"\n\t\"BytesDownloaded\"\t\t\"40\"\n}\n"

func TestInspectAndLauncherDir(t *testing.T) {
	root := t.TempDir()
	srv := filepath.Join(root, "sandstorm_server")
	sc := filepath.Join(root, "steamcmd")
	for _, d := range []string{filepath.Join(srv, "steamapps"), sc} {
		if err := os.MkdirAll(d, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	_ = os.WriteFile(filepath.Join(srv, "InsurgencyServer.exe"), []byte("x"), 0o600)
	_ = os.WriteFile(SteamCMDExe(sc), []byte("x"), 0o600)
	// O appmanifest guarda o caminho com barras duplicadas.
	launcher := strings.ReplaceAll(SteamCMDExe(sc), `\`, `\\`)
	_ = os.WriteFile(manifestPath(srv), []byte(strings.Replace(manifest, "%s", launcher, 1)), 0o600)

	in := Inspect(server.Profile{ServerDir: srv, SteamCMDDir: sc})
	if !in.ServerInstalled || !in.SteamCMDFound || in.BuildID != "25534773" || in.SizeOnDisk != 4994664175 || in.LastUpdated != 1794188468 {
		t.Errorf("%+v", in)
	}
	if got := LauncherDir(srv); !strings.EqualFold(got, sc) {
		t.Errorf("LauncherDir = %q, want %q", got, sc)
	}
	if LauncherDir(t.TempDir()) != "" {
		t.Error("sem manifesto deveria ser vazio")
	}
}

func TestSplitCRLF(t *testing.T) {
	in := "Loading...\r Update state (0x61) downloading, progress: 1.00 (1 / 100)\r Update state (0x61) downloading, progress: 2.00 (2 / 100)\nSuccess! App '581330' fully installed.\n"
	sc := bufio.NewScanner(strings.NewReader(in))
	sc.Split(splitCRLF)
	var lines []string
	for sc.Scan() {
		if sc.Text() != "" {
			lines = append(lines, sc.Text())
		}
	}
	if len(lines) != 4 || !strings.HasPrefix(lines[3], "Success!") {
		t.Errorf("%q", lines)
	}
}

func TestProgressSpeedAndBusy(t *testing.T) {
	s := New(nil)
	if _, err := s.begin(InstallServer); err != nil {
		t.Fatal(err)
	}
	if _, err := s.begin(FetchSteamCMD); err == nil {
		t.Error("duas tarefas ao mesmo tempo")
	}
	s.progress("downloading", 10, 100, 0)
	st := s.State()
	if st.Percent != 10 || st.Total != 100 {
		t.Errorf("%+v", st)
	}
	s.end(InstallServer, nil)
	if st := s.State(); st.Task != Idle || st.Finished != InstallServer || st.Percent != 100 {
		t.Errorf("%+v", st)
	}
}

func TestRefusesUnsafeDirs(t *testing.T) {
	s := New(nil)
	if err := s.FetchSteamCMD(`C:\`); err == nil {
		t.Error("raiz de disco aceita")
	}
	if err := s.Install(t.TempDir(), "relativo", false); err == nil {
		t.Error("sem steamcmd/pasta relativa aceita")
	}
}
