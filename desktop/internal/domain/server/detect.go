package server

import (
	"errors"
	"net"
	"strconv"
	"strings"

	"sandstorm-server-assist/desktop/internal/domain/problems"
)

// Detected é o que foi achado num servidor já configurado: a linha de comando gravada no log ("Command Line:"),
// o Game.ini, o Admins.txt e o MapCycle.txt. Found diz quais campos de SyncKeys foram achados.
type Detected struct {
	Profile Profile
	Secrets Secrets
	Found   map[string]bool
}

// NewDetected começa vazio.
func NewDetected() Detected { return Detected{Found: map[string]bool{}} }

// Estados de um campo na comparação com o servidor existente.
const (
	SyncSame     = "same"     // igual nos dois lados
	SyncFill     = "fill"     // o app tem o valor padrão ou vazio: usa o do servidor sem perguntar
	SyncConflict = "conflict" // o usuário mudou o valor no app e o servidor tem outro: perguntar
)

// Diff é um campo comparado. Segredos nunca levam o valor (Mine e Server ficam vazios).
type Diff struct {
	Key    string `json:"key"`
	Kind   string `json:"kind"` // text, number, bool, map, list, secret
	Mine   string `json:"mine"`
	Server string `json:"server"`
	State  string `json:"state"`
}

type syncField struct {
	key, kind string
	get       func(Profile) string
	put       func(dst *Profile, src Profile)
}

func mapText(m, s string) string {
	if m == "" {
		return ""
	}
	return m + " / " + s
}

func boolText(b bool) string { return strconv.FormatBool(b) }

func adminsText(in []Admin) string {
	lines := make([]string, 0, len(in))
	for _, a := range in {
		lines = append(lines, strings.TrimSpace(a.SteamID+" "+a.Name))
	}
	return strings.Join(lines, "\n")
}

var syncFields = []syncField{
	{"hostname", "text", func(p Profile) string { return p.Hostname }, func(d *Profile, s Profile) { d.Hostname = s.Hostname }},
	{"port", "number", func(p Profile) string { return strconv.Itoa(p.Port) }, func(d *Profile, s Profile) { d.Port = s.Port }},
	{"queryPort", "number", func(p Profile) string { return strconv.Itoa(p.QueryPort) }, func(d *Profile, s Profile) { d.QueryPort = s.QueryPort }},
	{"maxPlayers", "number", func(p Profile) string { return strconv.Itoa(p.MaxPlayers) }, func(d *Profile, s Profile) { d.MaxPlayers = s.MaxPlayers }},
	{"startMap", "map", func(p Profile) string { return mapText(p.StartMap, p.StartScenario) }, func(d *Profile, s Profile) {
		d.StartMap, d.StartScenario = s.StartMap, s.StartScenario
	}},
	{"travelMap", "map", func(p Profile) string { return mapText(p.TravelMap, p.TravelScenario) }, func(d *Profile, s Profile) {
		d.TravelMap, d.TravelScenario = s.TravelMap, s.TravelScenario
	}},
	{"useMods", "bool", func(p Profile) string { return boolText(p.UseMods) }, func(d *Profile, s Profile) { d.UseMods = s.UseMods }},
	{"mutators", "list", func(p Profile) string { return strings.Join(p.Mutators, "\n") }, func(d *Profile, s Profile) {
		d.Mutators = append([]string(nil), s.Mutators...)
	}},
	{"rcon", "bool", func(p Profile) string { return boolText(p.Rcon) }, func(d *Profile, s Profile) { d.Rcon = s.Rcon }},
	{"rconPort", "number", func(p Profile) string { return strconv.Itoa(p.RconPort) }, func(d *Profile, s Profile) { d.RconPort = s.RconPort }},
	{"admins", "list", func(p Profile) string { return adminsText(p.Admins) }, func(d *Profile, s Profile) {
		d.Admins = append([]Admin(nil), s.Admins...)
	}},
	{"mapCycle", "list", func(p Profile) string { return strings.Join(p.MapCycle, "\n") }, func(d *Profile, s Profile) {
		d.MapCycle = append([]string(nil), s.MapCycle...)
	}},
	{"motd", "list", func(p Profile) string { return p.Motd }, func(d *Profile, s Profile) { d.Motd = s.Motd }},
	{"startLighting", "text", func(p Profile) string { return p.StartLighting }, func(d *Profile, s Profile) { d.StartLighting = s.StartLighting }},
	{"bindIp", "text", func(p Profile) string { return p.BindIP }, func(d *Profile, s Profile) { d.BindIP = s.BindIP }},
	{"modIds", "list", func(p Profile) string { return modIDsText(p) }, func(d *Profile, s Profile) {
		d.ModIDs, d.CmdModList = append([]int(nil), s.ModIDs...), s.CmdModList
	}},
	{"rules", "list", func(p Profile) string { return rulesText(p.Rules) }, func(d *Profile, s Profile) {
		d.Rules = map[string]string{}
		for k, v := range s.Rules {
			d.Rules[k] = v
		}
	}},
}

func modIDsText(p Profile) string {
	if !p.CmdModList {
		return ""
	}
	ids := make([]string, len(p.ModIDs))
	for i, id := range p.ModIDs {
		ids[i] = strconv.Itoa(id)
	}
	return strings.Join(ids, ", ")
}

// syncSecrets são os segredos que um servidor existente pode ter (o e-mail do mod.io não aparece no servidor).
var syncSecrets = []string{"serverPassword", "rconPassword", "gsltToken", "gameStatsToken", "modioApiKey"}

// SyncKeys são todos os campos que a sincronização conhece (para a tela mandar só estas chaves).
func SyncKeys() []string {
	out := make([]string, 0, len(syncFields)+len(syncSecrets))
	for _, f := range syncFields {
		out = append(out, f.key)
	}
	return append(out, syncSecrets...)
}

// Erros de validação apontam o campo do perfil; alguns campos da sincronização cobrem dois do perfil.
var validationToSync = map[string]string{"startScenario": "startMap", "travelScenario": "travelMap"}

// ParseCommandLine lê a linha de comando do InsurgencyServer.exe (como aparece no log, sem o executável).
// Campos com valor inválido ficam de fora.
func (d *Detected) ParseCommandLine(cmd string) {
	toks := splitArgs(cmd)
	if len(toks) == 0 {
		return
	}
	// Numa linha de comando achada, o que não aparece está desligado.
	for _, k := range []string{"useMods", "travelMap", "mutators", "rcon", "startLighting", "bindIp"} {
		d.Found[k] = true
	}
	d.Profile.UseMods, d.Profile.Rcon = false, false
	d.Profile.StartLighting, d.Profile.BindIP = "", ""
	d.Profile.TravelMap, d.Profile.TravelScenario, d.Profile.Mutators = "", "", nil

	if !strings.HasPrefix(toks[0], "-") {
		d.parseURL(toks[0])
		toks = toks[1:]
	}
	for _, t := range toks {
		if !strings.HasPrefix(t, "-") {
			continue
		}
		name, val, _ := strings.Cut(t[1:], "=")
		switch strings.ToLower(name) {
		case "port":
			d.setInt("port", val, &d.Profile.Port)
		case "queryport":
			d.setInt("queryPort", val, &d.Profile.QueryPort)
		case "maxplayers":
			d.setInt("maxPlayers", val, &d.Profile.MaxPlayers)
		case "hostname":
			if validHostname(val) == nil {
				d.Profile.Hostname, d.Found["hostname"] = val, true
			}
		case "mods":
			d.Profile.UseMods = true
		case "moddownloadtravelto":
			if m, s, ok := mapScenario(val); ok {
				d.Profile.TravelMap, d.Profile.TravelScenario = m, s
			}
		case "mutators":
			for _, m := range strings.Split(val, ",") {
				if m = strings.TrimSpace(m); mutatorRe.MatchString(m) {
					d.Profile.Mutators = append(d.Profile.Mutators, m)
				}
			}
		case "rcon":
			d.Profile.Rcon = true
		case "multihome":
			if ip := net.ParseIP(val); ip != nil && ip.To4() != nil && !ip.IsUnspecified() {
				d.Profile.BindIP = val
			}
		case "cmdmodlist":
			var ids []int
			for _, part := range strings.Split(val, ",") {
				if id, err := strconv.Atoi(strings.TrimSpace(part)); err == nil && id > 0 {
					ids = append(ids, id)
				}
			}
			if len(ids) > 0 && len(ids) <= 200 {
				d.Profile.ModIDs, d.Profile.CmdModList, d.Found["modIds"] = ids, true, true
			}
		case "rconlistenport":
			d.setInt("rconPort", val, &d.Profile.RconPort)
		case "rconpassword":
			d.setSecret("rconPassword", val)
		case "gslttoken":
			d.setSecret("gsltToken", val)
		case "gamestatstoken":
			d.setSecret("gameStatsToken", val)
		}
	}
}

func (d *Detected) parseURL(url string) {
	parts := strings.Split(url, "?")
	var m, s string
	for i, p := range parts {
		if i == 0 {
			m = p
			continue
		}
		k, v, ok := strings.Cut(p, "=")
		if !ok {
			continue
		}
		switch strings.ToLower(k) {
		case "scenario":
			s = v
		case "maxplayers":
			d.setInt("maxPlayers", v, &d.Profile.MaxPlayers)
		case "password":
			d.setSecret("serverPassword", v)
		case "lighting":
			for _, l := range []string{LightingDay, LightingNight} {
				if strings.EqualFold(v, l) {
					d.Profile.StartLighting = l
				}
			}
		}
	}
	if identRe.MatchString(m) && scenarioRe.MatchString(s) {
		d.Profile.StartMap, d.Profile.StartScenario, d.Found["startMap"] = m, s, true
	}
}

func mapScenario(v string) (string, string, bool) {
	m, rest, ok := strings.Cut(v, "?")
	if !ok {
		return "", "", false
	}
	k, s, ok := strings.Cut(rest, "=")
	if !ok || !strings.EqualFold(k, "Scenario") || !identRe.MatchString(m) || !scenarioRe.MatchString(s) {
		return "", "", false
	}
	return m, s, true
}

func (d *Detected) setInt(key, v string, dst *int) {
	n, err := strconv.Atoi(v)
	if err != nil || n < 1 || n > 65535 {
		return
	}
	*dst, d.Found[key] = n, true
}

func (d *Detected) setSecret(key, v string) {
	if ValidateSecret(key, v) != nil {
		return
	}
	if next, err := d.Secrets.With(key, v); err == nil {
		d.Secrets = next
	}
}

// splitArgs separa a linha por espaços, respeitando aspas (que saem do valor): -hostname="Meu servidor".
func splitArgs(s string) []string {
	var out []string
	var cur strings.Builder
	quoted, has := false, false
	for _, r := range s {
		switch {
		case r == '"':
			quoted, has = !quoted, true
		case (r == ' ' || r == '\t') && !quoted:
			if has {
				out = append(out, cur.String())
				cur.Reset()
				has = false
			}
		default:
			cur.WriteRune(r)
			has = true
		}
	}
	if has {
		out = append(out, cur.String())
	}
	return out
}

// ValidAdmins devolve só os administradores com SteamID64 e nome válidos.
func ValidAdmins(in []Admin) []Admin {
	var out []Admin
	for _, a := range in {
		test := Default()
		test.Admins = []Admin{a}
		if test.Validate() == nil {
			out = append(out, a)
		}
	}
	return out
}

// usable diz se o valor achado para um campo passa na validação (testado sobre o perfil padrão).
// Portas só têm a faixa conferida aqui; portas repetidas são tratadas no Merge, com o perfil inteiro.
func usable(f syncField, src Profile) bool {
	switch f.key {
	case "port":
		return validPort(src.Port)
	case "queryPort":
		return validPort(src.QueryPort)
	case "rconPort":
		return validPort(src.RconPort)
	}
	test := Default()
	f.put(&test, src)
	return test.Validate() == nil
}

func validPort(n int) bool { return n >= 1024 && n <= 65535 }

// Compare mostra, campo a campo, o que está no app e o que foi achado no servidor.
func Compare(cur Profile, curSec Secrets, d Detected) []Diff {
	def := Default()
	var out []Diff
	for _, f := range syncFields {
		if !d.Found[f.key] || !usable(f, d.Profile) {
			continue
		}
		mine, srv := f.get(cur), f.get(d.Profile)
		state := SyncConflict
		switch {
		case mine == srv:
			state = SyncSame
		case mine == "" || mine == f.get(def):
			state = SyncFill
		}
		out = append(out, Diff{Key: f.key, Kind: f.kind, Mine: mine, Server: srv, State: state})
	}
	for _, k := range syncSecrets {
		srv, _ := d.Secrets.Get(k)
		if srv == "" {
			continue
		}
		mine, _ := curSec.Get(k)
		state := SyncConflict
		switch {
		case mine == srv:
			state = SyncSame
		case mine == "":
			state = SyncFill
		}
		out = append(out, Diff{Key: k, Kind: "secret", State: state})
	}
	return out
}

// Merge aplica os valores do servidor: os campos "fill" sempre, os "conflict" só se a chave não estiver em
// keepMine. Um valor que deixa o perfil inválido (ex.: portas repetidas) volta ao do app e vai para skipped.
func Merge(cur Profile, curSec Secrets, d Detected, keepMine map[string]bool) (p Profile, sec Secrets, applied, skipped []string, err error) {
	p, sec = cur, curSec
	byKey := map[string]syncField{}
	for _, f := range syncFields {
		byKey[f.key] = f
	}
	for _, df := range Compare(cur, curSec, d) {
		if df.State == SyncSame || (df.State == SyncConflict && keepMine[df.Key]) {
			continue
		}
		if f, ok := byKey[df.Key]; ok {
			f.put(&p, d.Profile)
		} else {
			v, _ := d.Secrets.Get(df.Key)
			if next, e := sec.With(df.Key, v); e == nil {
				sec = next
			} else {
				skipped = append(skipped, df.Key)
				continue
			}
		}
		applied = append(applied, df.Key)
	}

	revert := func(key string) bool {
		for i, k := range applied {
			if k == key {
				byKey[key].put(&p, cur)
				applied = append(applied[:i], applied[i+1:]...)
				skipped = append(skipped, key)
				return true
			}
		}
		return false
	}
	for range len(syncFields) + 1 {
		verr := p.Validate()
		if verr == nil {
			return p, sec, applied, skipped, nil
		}
		field := ""
		var pe *problems.Error
		if errors.As(verr, &pe) {
			field = pe.Field
		}
		key := field
		if k, ok := validationToSync[field]; ok {
			key = k
		}
		if revert(key) {
			continue
		}
		// Portas repetidas: o erro aponta queryPort, mas a causa pode ser qualquer porta trazida do servidor.
		reverted := false
		if field == "queryPort" || field == "port" || field == "rconPort" {
			for _, k := range []string{"rconPort", "rcon", "queryPort", "port"} {
				if reverted = revert(k); reverted {
					break
				}
			}
		}
		if !reverted {
			return cur, curSec, nil, nil, verr
		}
	}
	return cur, curSec, nil, nil, p.Validate()
}
