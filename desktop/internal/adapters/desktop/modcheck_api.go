package desktop

import (
	"strconv"

	"sandstorm-server-assist/desktop/internal/adapters/modio"
	"sandstorm-server-assist/desktop/internal/domain/modcheck"
)

// ---- mods com problema (ADR-039) ----

// ModCheck é a verificação dos mods instalados: pares de mods com o mesmo arquivo do jogo em versões diferentes.
type ModCheck struct {
	Mods       int                 `json:"mods"`
	Conflicts  []modcheck.Conflict `json:"conflicts"`
	Unreadable []modcheck.ModRef   `json:"unreadable"`
}

// CheckMods lê os índices dos .pak dos mods instalados (só leitura) e compara os mods dois a dois.
func (s *ServerAPI) CheckMods() ModCheck {
	mf := modio.InstalledFiles(modio.RootDir())
	return ModCheck{Mods: len(mf.Mods), Conflicts: modcheck.Conflicts(mf.Mods), Unreadable: mf.Unreadable}
}

// CrashCheck é o resultado de uma mensagem de crash colada na tela.
type CrashCheck struct {
	Mods         int                `json:"mods"`
	SizeMismatch bool               `json:"sizeMismatch"`
	Findings     []modcheck.Finding `json:"findings"`
}

// CheckCrash procura, nos mods instalados, os objetos citados numa mensagem de crash (de um jogador ou do
// servidor). O texto não é guardado nem vai para o log.
func (s *ServerAPI) CheckCrash(text string) CrashCheck {
	if len(text) > modcheck.MaxText {
		text = text[:modcheck.MaxText]
	}
	mf := modio.InstalledFiles(modio.RootDir())
	findings, mismatch := modcheck.Check(text, mf.Mods)
	return CrashCheck{Mods: len(mf.Mods), SizeMismatch: mismatch, Findings: findings}
}

// modOwnerNames devolve "Nome (#id)" dos mods instalados que trazem o objeto.
func modOwnerNames(object string) []string {
	mf := modio.InstalledFiles(modio.RootDir())
	var out []string
	for _, m := range modcheck.Owners(object, mf.Mods) {
		out = append(out, m.Name+" (#"+strconv.Itoa(m.ID)+")")
	}
	return out
}
