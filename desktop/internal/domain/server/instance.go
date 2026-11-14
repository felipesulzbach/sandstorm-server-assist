package server

import (
	"regexp"
	"strconv"
	"strings"

	"sandstorm-server-assist/desktop/internal/domain/problems"
)

// Vários servidores numa só instalação (ADR-048). Cada servidor tem um número de instância. A instância 1 é a
// principal e usa os arquivos de sempre (Insurgency/Saved, Admins.txt, MapCycle.txt, Motd.txt). A instância N ≥ 2
// inicia com -saveddirsuffix=N (o jogo usa Insurgency/Saved_N para o Game.ini, o Engine.ini e os logs) e arquivos
// próprios em Insurgency/Config/Server (Admins_N.txt, MapCycle_N.txt, Motd_N.txt). Perfis antigos têm 0, que vale 1.

// MaxInstance limita o número de instância aceito (o app guarda no máximo 20 servidores).
const MaxInstance = 99

// InstanceNumber devolve o número de instância efetivo (0 e valores fora da faixa valem 1).
func InstanceNumber(n int) int {
	if n < 2 || n > MaxInstance {
		return 1
	}
	return n
}

// Files são os caminhos de uma instância, relativos à pasta do servidor (com /). AdminsName, MapCycleName e
// MotdName são os nomes que vão na linha de comando (sem .txt).
type Files struct {
	Saved        string
	GameIni      string
	EngineIni    string
	LogsDir      string
	Admins       string
	MapCycle     string
	Motd         string
	AdminsName   string
	MapCycleName string
	MotdName     string
}

// configServer é a pasta dos arquivos de texto do servidor.
const configServer = "Insurgency/Config/Server/"

// InstanceFiles devolve os arquivos da instância n.
func InstanceFiles(n int) Files {
	n = InstanceNumber(n)
	saved, suffix := "Insurgency/Saved", ""
	if n > 1 {
		suffix = "_" + strconv.Itoa(n)
		saved += suffix
	}
	f := Files{
		Saved:        saved,
		GameIni:      saved + "/Config/WindowsServer/Game.ini",
		EngineIni:    saved + "/Config/WindowsServer/Engine.ini",
		LogsDir:      saved + "/Logs",
		AdminsName:   "Admins" + suffix,
		MapCycleName: "MapCycle" + suffix,
		MotdName:     "Motd" + suffix,
	}
	f.Admins = configServer + f.AdminsName + ".txt"
	f.MapCycle = configServer + f.MapCycleName + ".txt"
	f.Motd = configServer + f.MotdName + ".txt"
	return f
}

// instanceArgs são os argumentos da instância na linha de comando. A principal fica como sempre: -AdminList e
// -MapCycle só com lista. As outras sempre levam o próprio Admins e Motd (senão o jogo leria os da principal) e o
// próprio MapCycle quando têm rotação (sem rotação, valem os mapas do MapCycle.txt da instalação).
func instanceArgs(p Profile) []string {
	n := InstanceNumber(p.Instance)
	f := InstanceFiles(n)
	var args []string
	if n > 1 {
		args = append(args, "-saveddirsuffix="+strconv.Itoa(n), "-AdminList="+f.AdminsName, "-Motd="+f.MotdName)
	} else if len(p.Admins) > 0 {
		args = append(args, "-AdminList="+f.AdminsName)
	}
	if len(p.MapCycle) > 0 {
		args = append(args, "-MapCycle="+f.MapCycleName)
	}
	return args
}

var reSavedSuffix = regexp.MustCompile(`(?i)(?:^|\s)-saveddirsuffix=("[^"]*"|\S*)`)

// InstanceOf diz de qual instância é um processo do servidor, pela linha de comando: sem -saveddirsuffix é a 1;
// com um número de 2 a MaxInstance é esse número; qualquer outro valor (servidor de outra ferramenta) devolve 0,
// que não é de nenhum servidor do app.
func InstanceOf(cmdline string) int {
	m := reSavedSuffix.FindStringSubmatch(cmdline)
	if m == nil {
		return 1
	}
	v := strings.Trim(m[1], `"`)
	n, err := strconv.Atoi(v)
	if err != nil || n < 2 || n > MaxInstance || strconv.Itoa(n) != v {
		return 0
	}
	return n
}

// NextInstance devolve o menor número de instância livre (used são os números dos outros servidores).
func NextInstance(used []int) int {
	taken := map[int]bool{}
	for _, n := range used {
		taken[InstanceNumber(n)] = true
	}
	for n := 1; n <= MaxInstance; n++ {
		if !taken[n] {
			return n
		}
	}
	return 0
}

// ---- portas ----

// Peer é outro servidor do app, com o nome que a tela mostra.
type Peer struct {
	ID      string
	Name    string
	Profile Profile
	Secrets Secrets
}

// PortConflict é uma porta deste servidor que outro servidor também usa.
type PortConflict struct {
	// Field é o campo deste servidor (port, queryPort ou rconPort).
	Field string `json:"field"`
	Port  int    `json:"port"`
	// Server é o nome do outro servidor.
	Server string `json:"server"`
}

// portUse é uma porta aberta pelo servidor: o jogo e a consulta em UDP, o RCON em TCP.
type portUse struct {
	field string
	port  int
	tcp   bool
}

func (p Profile) portUses() []portUse {
	out := []portUse{{"port", p.Port, false}, {"queryPort", p.QueryPort, false}}
	if p.Rcon {
		out = append(out, portUse{"rconPort", p.RconPort, true})
	}
	return out
}

// PortConflicts lista as portas de p que algum dos outros servidores também usa (mesmo protocolo).
func PortConflicts(p Profile, others []Peer) []PortConflict {
	out := []PortConflict{}
	for _, mine := range p.portUses() {
		for _, o := range others {
			hit := false
			for _, theirs := range o.Profile.portUses() {
				hit = hit || (theirs.tcp == mine.tcp && theirs.port == mine.port)
			}
			if hit {
				out = append(out, PortConflict{Field: mine.field, Port: mine.port, Server: o.Name})
				break
			}
		}
	}
	return out
}

// PortBusy diz se outro programa do PC já usa a porta (tcp=false: UDP). Quem confere é o adaptador.
type PortBusy func(port int, tcp bool) bool

// BusyPorts lista as portas de p que outro programa do PC já usa (Server vazio na resposta).
func BusyPorts(p Profile, busy PortBusy) []PortConflict {
	out := []PortConflict{}
	if busy == nil {
		return out
	}
	for _, u := range p.portUses() {
		if busy(u.port, u.tcp) {
			out = append(out, PortConflict{Field: u.field, Port: u.port})
		}
	}
	return out
}

// FreePorts devolve as primeiras portas de jogo, consulta e RCON que nenhum dos outros servidores usa e que nenhum
// outro programa do PC ocupa (busy pode ser nil), a partir das padrão (27102, 27131 e 27015) e subindo de 1 em 1.
// Ex.: com o MongoDB na 27017, o terceiro servidor fica com o RCON na 27018.
func FreePorts(others []Profile, busy PortBusy) (port, query, rcon int) {
	udp, tcp := map[int]bool{}, map[int]bool{}
	for _, o := range others {
		udp[o.Port], udp[o.QueryPort] = true, true
		if o.Rcon {
			tcp[o.RconPort] = true
		}
	}
	d := Default()
	for k := 0; k < 1000; k++ {
		port, query, rcon = d.Port+k, d.QueryPort+k, d.RconPort+k
		if udp[port] || udp[query] || tcp[rcon] || port == query || rcon == port || rcon == query {
			continue
		}
		if busy != nil && (busy(port, false) || busy(query, false) || busy(rcon, true)) {
			continue
		}
		return port, query, rcon
	}
	return d.Port, d.QueryPort, d.RconPort
}

// CheckStart diz por que p não pode iniciar junto dos servidores que já estão rodando (running): uma porta usada
// por outro deles ou o mesmo token GSLT (a Steam aceita cada token em um servidor por vez).
func CheckStart(p Profile, sec Secrets, running []Peer) error {
	if c := PortConflicts(p, running); len(c) > 0 {
		e := problems.New(problems.PortsShared, "")
		e.Field = c[0].Field
		// portField diz qual porta (jogo, consulta ou RCON): a tela mostra o nome dela.
		e.Params = map[string]string{"port": strconv.Itoa(c[0].Port), "portField": c[0].Field, "server": c[0].Server}
		return e
	}
	if sec.GSLTToken != "" {
		for _, o := range running {
			if o.Secrets.GSLTToken == sec.GSLTToken {
				e := problems.New(problems.GSLTShared, "")
				e.Params = map[string]string{"server": o.Name}
				return e
			}
		}
	}
	return nil
}
