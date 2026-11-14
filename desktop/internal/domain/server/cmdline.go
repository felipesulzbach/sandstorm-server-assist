package server

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"sandstorm-server-assist/desktop/internal/domain/problems"
	"sandstorm-server-assist/desktop/internal/domain/redact"
)

// SecurityNone é o modo normal do mod.io depois de autorizado (comprovado em 09/11/2026, ADR-003).
const SecurityNone = "none"

// Launch é o que o app precisa para iniciar o servidor.
type Launch struct {
	// Args são os argumentos, um por item (para mostrar e testar).
	Args []string
	// CmdLine é a linha de comando que vai para o Windows (sem o executável). O Unreal lê a linha inteira e só
	// entende -hostname="nome com espaço" com as aspas depois do =, então ela é montada aqui e não pelo Go.
	CmdLine string
}

// Shown devolve a linha de comando com os segredos mascarados (para a tela e o log).
func (l Launch) Shown(s Secrets) string {
	return redact.Values(redact.Text(l.CmdLine), s.Values()...)
}

// BuildLaunch monta a linha de comando do InsurgencyServer.exe. securityCode é "none" no uso normal ou um código
// novo do mod.io (5 caracteres) só na partida que autoriza.
func BuildLaunch(p Profile, s Secrets, securityCode string) (Launch, error) {
	if err := p.Validate(); err != nil {
		return Launch{}, err
	}
	if securityCode == "" {
		securityCode = SecurityNone
	}
	if securityCode != SecurityNone && !ValidCode(securityCode) {
		return Launch{}, problems.Field("securityCode", "validation.code")
	}
	for _, k := range SecretKeys {
		if v, _ := s.Get(k); v != "" {
			if err := ValidateSecret(k, v); err != nil {
				return Launch{}, err
			}
		}
	}

	url := fmt.Sprintf("%s?Scenario=%s?MaxPlayers=%d", p.StartMap, p.StartScenario, p.MaxPlayers)
	if p.StartLighting != "" {
		url += "?Lighting=" + p.StartLighting
	}
	if s.ServerPassword != "" {
		url += "?Password=" + s.ServerPassword
	}
	args := []string{
		url,
		"-Port=" + strconv.Itoa(p.Port),
		"-QueryPort=" + strconv.Itoa(p.QueryPort),
		`-hostname="` + p.Hostname + `"`,
	}
	if p.UseMods {
		args = append(args, "-SecurityCode="+securityCode, "-mods")
		if p.TravelMap != "" {
			args = append(args, fmt.Sprintf("-ModDownloadTravelTo=%s?Scenario=%s", p.TravelMap, p.TravelScenario))
		}
		if p.CmdModList && len(p.ModIDs) > 0 {
			ids := make([]string, len(p.ModIDs))
			for i, id := range p.ModIDs {
				ids[i] = strconv.Itoa(id)
			}
			args = append(args, `-CmdModList="`+strings.Join(ids, ",")+`"`)
		}
	}
	if p.BindIP != "" {
		args = append(args, "-MultiHome="+p.BindIP)
	}
	if len(p.Mutators) > 0 {
		args = append(args, "-mutators="+strings.Join(p.Mutators, ","))
	}
	if len(p.Admins) > 0 {
		args = append(args, "-AdminList=Admins")
	}
	if len(p.MapCycle) > 0 {
		args = append(args, "-MapCycle=MapCycle")
	}
	if p.Rcon && s.RconPassword != "" {
		args = append(args, "-Rcon", "-RconPassword="+s.RconPassword, "-RconListenPort="+strconv.Itoa(p.RconPort))
	}
	if s.GSLTToken != "" {
		args = append(args, "-GSLTToken="+s.GSLTToken)
	}
	if s.GameStatsToken != "" {
		args = append(args, "-GameStats", "-GameStatsToken="+s.GameStatsToken)
	}
	return Launch{Args: args, CmdLine: strings.Join(args, " ")}, nil
}

// AdminsFile é o conteúdo do Admins.txt: um SteamID64 por linha, com o nome como comentário antes.
func AdminsFile(admins []Admin) string {
	var b strings.Builder
	for _, a := range admins {
		if a.Name != "" {
			b.WriteString(";" + a.Name + "\r\n")
		}
		b.WriteString(a.SteamID + "\r\n")
	}
	return b.String()
}

// ParseAdmins lê um Admins.txt (linhas com ; ou // são comentários; o comentário logo antes vira o nome).
func ParseAdmins(content string) []Admin {
	var out []Admin
	name := ""
	for _, line := range strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n") {
		line = strings.TrimSpace(line)
		switch {
		case line == "":
			name = ""
		case strings.HasPrefix(line, ";"):
			name = strings.TrimSpace(strings.TrimPrefix(line, ";"))
		case strings.HasPrefix(line, "//"):
			name = strings.TrimSpace(strings.TrimPrefix(line, "//"))
		default:
			out = append(out, Admin{SteamID: line, Name: name})
			name = ""
		}
	}
	return out
}

// reScenarioTwice acha o "Scenario=" repetido que vem de colar a linha da página do mod dentro do
// (Scenario="..."), como em (Scenario="Scenario=Scenario_Outback_Checkpoint_Security",Lighting="Day"). O jogo
// recusa essa entrada ("Not adding ... to map vote, no asset found").
var reScenarioTwice = regexp.MustCompile(`(?i)(Scenario\s*=\s*"?)\s*Scenario\s*=\s*`)

// FixMapCycleLine corrige o "Scenario=" repetido e a linha solta "Scenario=Scenario_X", e escreve toda linha de
// cenário no formato (Scenario="Scenario_X",Lighting="Day|Night") (ADR-046). Sem iluminação, vale Day, a iluminação
// padrão dos cenários do jogo. Linhas que o app não entende voltam só sem os espaços das pontas.
func FixMapCycleLine(line string) string {
	line = strings.TrimSpace(line)
	line = reScenarioTwice.ReplaceAllString(line, "$1")
	if rest, ok := strings.CutPrefix(line, "Scenario="); ok && !strings.HasPrefix(rest, `"`) {
		line = rest
	}
	if scenarioRe.MatchString(line) {
		return CycleLine(line, LightingDay)
	}
	if strings.HasPrefix(line, "(") && strings.HasSuffix(line, ")") && LineScenario(line) != "" &&
		!reLightingAttr.MatchString(line) {
		return line[:len(line)-1] + `,Lighting="` + LightingDay + `")`
	}
	return line
}

// CycleLine é a linha da rotação de um cenário com a iluminação.
func CycleLine(scenario, lighting string) string {
	return `(Scenario="` + scenario + `",Lighting="` + lighting + `")`
}

// FixMapCycle corrige todas as linhas (FixMapCycleLine).
func FixMapCycle(lines []string) []string {
	if lines == nil {
		return nil
	}
	out := make([]string, len(lines))
	for i, l := range lines {
		out[i] = FixMapCycleLine(l)
	}
	return out
}

// MapCycleFile é o conteúdo do MapCycle.txt (uma entrada por linha, já corrigida).
func MapCycleFile(lines []string) string {
	if len(lines) == 0 {
		return ""
	}
	return strings.Join(FixMapCycle(lines), "\r\n") + "\r\n"
}

// ParseMapCycle lê um MapCycle.txt e devolve só as linhas ativas (sem comentários e linhas vazias).
func ParseMapCycle(content string) []string {
	var out []string
	for _, line := range strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n") {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "//") && !strings.HasPrefix(line, ";") {
			out = append(out, FixMapCycleLine(line))
		}
	}
	return out
}
