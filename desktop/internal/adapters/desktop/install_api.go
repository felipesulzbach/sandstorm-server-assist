package desktop

import (
	"path/filepath"

	"sandstorm-server-assist/desktop/internal/app/installer"
	"sandstorm-server-assist/desktop/internal/domain/problems"
)

// InstallAPI: SteamCMD e instalação/atualização do servidor com progresso (evento "install:state").
type InstallAPI struct{ a *app }

// State devolve a tarefa em andamento.
func (i *InstallAPI) State() installer.State { return i.a.install.State() }

// Info diz o que está instalado.
func (i *InstallAPI) Info() installer.Info { return installer.Inspect(i.a.cfg.Profile()) }

// FetchSteamCMD baixa o SteamCMD. Sem pasta escolhida, usa "steamcmd" ao lado da pasta do servidor.
func (i *InstallAPI) FetchSteamCMD() error {
	p := i.a.cfg.Profile()
	dir := p.SteamCMDDir
	if dir == "" {
		if p.ServerDir == "" {
			return problems.New(problems.ServerDirInvalid, "")
		}
		dir = filepath.Join(filepath.Dir(p.ServerDir), "steamcmd")
		if _, err := i.a.cfg.SetSteamCMDDir(dir); err != nil {
			return err
		}
	}
	return i.a.install.FetchSteamCMD(dir)
}

// InstallServer instala ou atualiza o servidor. validate=true confere todos os arquivos.
func (i *InstallAPI) InstallServer(validate bool) error {
	p := i.a.cfg.Profile()
	if p.ServerDir == "" {
		return problems.New(problems.ServerDirInvalid, "")
	}
	if len(i.a.runner.Status().PIDs) > 0 {
		return problems.New(problems.ServerRunning, "")
	}
	return i.a.install.Install(p.SteamCMDDir, p.ServerDir, validate)
}

// Cancel cancela a tarefa atual.
func (i *InstallAPI) Cancel() { i.a.install.Cancel() }
