package config

import (
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"testing"
	"time"

	"sandstorm-server-assist/desktop/internal/domain/server"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

// Servidor configurado à mão (valores de exemplo): o app lê a linha de comando do log mais novo, os arquivos de
// configuração e pergunta só o que o usuário mudou.
func TestScanAndApplyExisting(t *testing.T) {
	dir := t.TempDir()
	srv := filepath.Join(dir, "sandstorm_server")
	logs := filepath.Join(srv, LogsDir)
	writeFile(t, filepath.Join(logs, "Insurgency-backup-2026.01.01-10.00.00.log"),
		"LogInit: Command Line:  Hideout?Scenario=Scenario_Hideout_Checkpoint_Security -Port=27900 -QueryPort=27931\r\n")
	writeFile(t, filepath.Join(logs, "Insurgency.log"), "Log file open\r\n"+
		"LogInit: Command Line:  Farmhouse?Scenario=Scenario_Farmhouse_Checkpoint_Security?MaxPlayers=10 -SecurityCode=none -mods "+
		"-ModDownloadTravelTo=Containment?Scenario=Scenario_Containment_Checkpoint_Security -Port=27200 -QueryPort=27231 "+
		"-hostname=\"[XX] Example Server\"\r\nLogInit: other line\r\n")
	writeFile(t, filepath.Join(srv, AdminsTxt), ";Alice\r\n76561198000000001\r\nnot-a-steamid\r\n")
	writeFile(t, filepath.Join(srv, GameIni), "[Rcon]\r\nbEnabled=True\r\nPassword=rcon#Senha1\r\nListenPort=27015\r\n")

	v := &memVault{s: server.Secrets{RconPassword: "minhaSenha9"}}
	s := Open(filepath.Join(dir, "config.json"), v)
	if _, err := s.SetServerDir(srv); err != nil {
		t.Fatal(err)
	}
	p := s.Profile()
	p.Hostname = "Digitado por mim"
	if _, err := s.SaveProfile(p); err != nil {
		t.Fatal(err)
	}

	scan, err := s.ScanExisting("")
	if err != nil || !scan.Found || scan.LogFile != "Insurgency.log" || !scan.RconInIni {
		t.Fatalf("%+v %v", scan, err)
	}
	states := map[string]string{}
	for _, d := range scan.Diffs {
		states[d.Key] = d.State
		if d.Kind == "secret" && (d.Mine != "" || d.Server != "") {
			t.Errorf("segredo na tela: %+v", d)
		}
	}
	if states["hostname"] != server.SyncConflict || states["rconPassword"] != server.SyncConflict ||
		states["port"] != server.SyncFill || states["admins"] != server.SyncFill || states["maxPlayers"] != server.SyncFill {
		t.Errorf("estados: %v", states)
	}

	if _, err := s.ApplyExisting("", "", []string{"comando-qualquer"}); err == nil {
		t.Error("chave desconhecida aceita")
	}
	r, err := s.ApplyExisting("", "", []string{"hostname"})
	if err != nil {
		t.Fatal(err)
	}
	got := s.Profile()
	if got.Hostname != "Digitado por mim" || got.Port != 27200 || got.MaxPlayers != 10 || got.TravelMap != "Containment" ||
		len(got.Admins) != 1 || !got.Rcon || got.ServerDir != srv {
		t.Errorf("perfil: %+v", got)
	}
	if v.s.RconPassword != "rcon#Senha1" || !slices.Contains(r.Kept, "hostname") || !slices.Contains(r.Applied, "rconPassword") {
		t.Errorf("resultado %+v cofre %+v", r, v.s)
	}

	// Depois de aplicar, nada mais é diferente (só o nome que o usuário decidiu manter).
	again, _ := s.ScanExisting("")
	for _, d := range again.Diffs {
		if d.State != server.SyncSame && d.Key != "hostname" {
			t.Errorf("ainda diferente: %+v", d)
		}
	}
}

// O usuário escolheu a pasta de cima do servidor: o app usa a subpasta com o InsurgencyServer.exe.
func TestResolveServerDirAndLogOrder(t *testing.T) {
	top := t.TempDir()
	srv := filepath.Join(top, "sandstorm_server")
	writeFile(t, filepath.Join(srv, ServerExe), "x")
	writeFile(t, filepath.Join(top, "support", "start.bat"), "x")
	if got := ResolveServerDir(top); got != srv {
		t.Errorf("subpasta: %q", got)
	}
	if got := ResolveServerDir(srv); got != srv {
		t.Errorf("pasta certa mudou: %q", got)
	}
	writeFile(t, filepath.Join(top, "outro", ServerExe), "x")
	if got := ResolveServerDir(top); got != top {
		t.Errorf("dois servidores: deveria deixar a escolha do usuário, veio %q", got)
	}

	// Config antiga com a pasta de cima é corrigida ao abrir.
	cfgPath := filepath.Join(t.TempDir(), "config.json")
	only := t.TempDir()
	writeFile(t, filepath.Join(only, "sandstorm_server", ServerExe), "x")
	writeFile(t, cfgPath, `{"version":1,"profile":{"serverDir":`+strconvQuote(only)+`,"hostname":"x","port":27102,"queryPort":27131,"maxPlayers":8}}`)
	if got := Open(cfgPath, &memVault{}).Profile().ServerDir; got != filepath.Join(only, "sandstorm_server") {
		t.Errorf("config antiga: %q", got)
	}

	// Logs: o mais recente pela data, inclusive Insurgency_2.log (criado quando o primeiro estava aberto).
	logs := filepath.Join(srv, LogsDir)
	old, mid, newest := filepath.Join(logs, "Insurgency-backup-2026.01.01-10.00.00.log"), filepath.Join(logs, "Insurgency.log"), filepath.Join(logs, "Insurgency_2.log")
	for i, p := range []string{old, mid, newest} {
		writeFile(t, p, "LogInit: Command Line:  -Port=2720"+string(rune('0'+i))+"\r\n")
		at := time.Date(2026, 1, 1, 10, i, 0, 0, time.UTC)
		if err := os.Chtimes(p, at, at); err != nil {
			t.Fatal(err)
		}
	}
	writeFile(t, filepath.Join(logs, "Modio.log"), "x")
	if got := LogFiles(srv); !slices.Equal(got, []string{newest, mid, old}) {
		t.Errorf("ordem: %v", got)
	}
	if cmd, file := lastCommandLine(srv); cmd != "-Port=27202" || file != "Insurgency_2.log" {
		t.Errorf("linha: %q %q", cmd, file)
	}
}

func strconvQuote(s string) string { return strconv.Quote(s) }

func TestScanWithoutServer(t *testing.T) {
	s := Open(filepath.Join(t.TempDir(), "config.json"), &memVault{})
	if _, err := s.ScanExisting(""); err == nil {
		t.Error("sem pasta do servidor deveria dar erro")
	}
	if _, err := s.SetServerDir(filepath.Join(t.TempDir(), "vazio")); err != nil {
		t.Fatal(err)
	}
	scan, err := s.ScanExisting("")
	if err != nil || scan.Found {
		t.Errorf("pasta vazia: %+v %v", scan, err)
	}
}

func TestLogRetention(t *testing.T) {
	root := t.TempDir()
	logs := filepath.Join(root, LogsDir)
	if err := os.MkdirAll(logs, 0o700); err != nil {
		t.Fatal(err)
	}
	base := time.Now().Add(-time.Hour)
	names := []string{
		"Insurgency.log",
		"Insurgency-backup-2026.11.09-01.00.00.log",
		"Insurgency-backup-2026.11.09-02.00.00.log",
		"Insurgency-backup-2026.11.09-03.00.00.log",
		"Insurgency_2.log",
		"Game.log", // não é log do servidor: nunca entra
	}
	for i, n := range names {
		p := filepath.Join(logs, n)
		if err := os.WriteFile(p, []byte("x"), 0o600); err != nil {
			t.Fatal(err)
		}
		mt := base.Add(time.Duration(i) * time.Minute)
		if n == "Insurgency.log" {
			mt = base.Add(time.Hour)
		}
		_ = os.Chtimes(p, mt, mt)
	}
	if got := len(ServerLogs(root)); got != 5 {
		t.Fatalf("ServerLogs = %d", got)
	}
	n, freed, err := DeleteOldLogs(root, 3)
	if err != nil || n != 2 || freed != 2 {
		t.Fatalf("deleted %d (%d bytes) %v", n, freed, err)
	}
	left := ServerLogs(root)
	if len(left) != 3 || left[0].Name != "Insurgency.log" || left[1].Name != "Insurgency_2.log" {
		t.Fatalf("left %+v", left)
	}
	if _, err := os.Stat(filepath.Join(logs, "Game.log")); err != nil {
		t.Fatal("deleted a file that is not a server log")
	}

	s := Open(filepath.Join(root, "config.json"), &memVault{})
	if st := s.Settings(); st.LogCleanup || st.LogKeep != 10 {
		t.Fatalf("defaults %+v", st)
	}
	if _, err := s.SetLogRetention(true, 1); err == nil {
		t.Fatal("keep below the minimum accepted")
	}
	if st, err := s.SetLogRetention(true, 5); err != nil || !st.LogCleanup || st.LogKeep != 5 {
		t.Fatalf("%+v %v", st, err)
	}
}
