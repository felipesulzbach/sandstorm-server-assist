package server

import (
	"sort"
	"strconv"
	"strings"

	"sandstorm-server-assist/desktop/internal/domain/ini"
	"sandstorm-server-assist/desktop/internal/domain/problems"
)

// Regras do jogo (ADR-022): opções do Game.ini e do Engine.ini que o app sabe gravar. Os nomes, seções e padrões
// vêm do guia oficial de administração do servidor (Server Admin Guide da New World Interactive). O texto da tela
// vem das traduções (rules.<id>).
//
// Profile.Rules nil = o app não mexe nessas opções (servidor configurado à mão continua igual). Depois que o
// usuário salva a aba Regras, o app passa a cuidar das opções do catálogo: as que têm valor são gravadas, as vazias
// são apagadas do arquivo (o jogo volta ao padrão dele).

// Arquivos de configuração das regras.
const (
	FileGame   = "game"
	FileEngine = "engine"
)

// Tipos de valor.
const (
	KindBool  = "bool"
	KindInt   = "int"
	KindFloat = "float"
)

// RuleDef é uma opção do catálogo.
type RuleDef struct {
	ID      string  `json:"id"`
	Group   string  `json:"group"`
	File    string  `json:"file"`
	Section string  `json:"section"`
	Key     string  `json:"key"`
	Kind    string  `json:"kind"`
	Min     float64 `json:"min"`
	Max     float64 `json:"max"`
	// Default é o padrão do jogo (vazio quando varia por modo).
	Default string `json:"default"`
	// also são chaves gravadas junto quando a opção liga (votação: a lista de assuntos).
	also [][3]string
}

// Grupos da aba Regras, na ordem da tela.
var RuleGroups = []string{"coop", "checkpoint", "objectives", "friendlyFire", "teamKills", "chat", "match", "supply", "votes", "network"}

const (
	secGame   = "/Script/Insurgency.INSGameMode"
	secMulti  = "/Script/Insurgency.INSMultiplayerMode"
	secCoop   = "/Script/Insurgency.INSCoopMode"
	secCheck  = "/Script/Insurgency.INSCheckpointGameMode"
	secTeam   = "/Script/Insurgency.TeamInfo"
	secVote   = "/Script/Insurgency.VoteIssueKick"
	secNet    = "/Script/OnlineSubsystemUtils.IpNetDriver"
	voteIssue = "/Script/Insurgency.VoteIssueKick"
)

func b(id, group, sec, key, def string) RuleDef {
	return RuleDef{ID: id, Group: group, File: FileGame, Section: sec, Key: key, Kind: KindBool, Max: 1, Default: def}
}

func n(id, group, sec, key string, minV, maxV float64, def string) RuleDef {
	return RuleDef{ID: id, Group: group, File: FileGame, Section: sec, Key: key, Kind: KindInt, Min: minV, Max: maxV, Default: def}
}

func f(id, group, sec, key string, minV, maxV float64, def string) RuleDef {
	return RuleDef{ID: id, Group: group, File: FileGame, Section: sec, Key: key, Kind: KindFloat, Min: minV, Max: maxV, Default: def}
}

func eng(r RuleDef) RuleDef {
	r.File = FileEngine
	return r
}

// Rules é o catálogo (Server Admin Guide; tick rate e banda são opções do Unreal no Engine.ini).
var Rules = []RuleDef{
	// Coop (bots)
	f("aiDifficulty", "coop", secCoop, "AIDifficulty", 0, 1, "0.5"),
	n("minimumEnemies", "coop", secCoop, "MinimumEnemies", 0, 100, "6"),
	n("maximumEnemies", "coop", secCoop, "MaximumEnemies", 0, 100, "12"),
	n("friendlyBotQuota", "coop", secCoop, "FriendlyBotQuota", 0, 32, "4"),
	b("vehicleInsertion", "coop", secCoop, "bUseVehicleInsertion", "True"),
	// Checkpoint
	n("defendTimer", "checkpoint", secCheck, "DefendTimer", 0, 3600, "90"),
	n("defendTimerFinal", "checkpoint", secCheck, "DefendTimerFinal", 0, 3600, "180"),
	n("retreatTimer", "checkpoint", secCheck, "RetreatTimer", 0, 600, "10"),
	f("respawnDpr", "checkpoint", secCheck, "RespawnDPR", 0, 1, "0.1"),
	n("respawnDelay", "checkpoint", secCheck, "RespawnDelay", 0, 600, "20"),
	n("postCaptureRushTimer", "checkpoint", secCheck, "PostCaptureRushTimer", 0, 600, "30"),
	f("counterAttackRespawnDpr", "checkpoint", secCheck, "CounterAttackRespawnDPR", 0, 1, "0.2"),
	n("counterAttackRespawnDelay", "checkpoint", secCheck, "CounterAttackRespawnDelay", 0, 600, "20"),
	f("enemyRespawnMin", "checkpoint", secCheck, "ObjectiveTotalEnemyRespawnMultiplierMin", 0, 10, "1"),
	f("enemyRespawnMax", "checkpoint", secCheck, "ObjectiveTotalEnemyRespawnMultiplierMax", 0, 10, "1"),
	f("finalCacheBots", "checkpoint", secCheck, "FinalCacheBotQuotaMultiplier", 0, 10, "1.5"),
	// Objetivos
	n("captureTime", "objectives", secGame, "ObjectiveCaptureTime", 1, 600, ""),
	n("resetTime", "objectives", secGame, "ObjectiveResetTime", -1, 600, "-1"),
	f("captureSpeedup", "objectives", secGame, "ObjectiveSpeedup", 0, 10, "0.25"),
	n("captureSpeedupPlayers", "objectives", secGame, "ObjectiveMaxSpeedupPlayers", 1, 64, "4"),
	// Fogo amigo e equipes
	b("friendlyFire", "friendlyFire", secMulti, "bAllowFriendlyFire", "True"),
	f("friendlyFireModifier", "friendlyFire", secMulti, "FriendlyFireModifier", 0, 1, "0.2"),
	f("friendlyFireReflect", "friendlyFire", secMulti, "FriendlyFireReflect", 0, 1, "0"),
	b("autoBalance", "friendlyFire", secMulti, "bAutoBalanceTeams", "True"),
	n("autoBalanceDelay", "friendlyFire", secMulti, "AutoBalanceDelay", 0, 600, "10"),
	b("teamSelect", "friendlyFire", secMulti, "bAllowPlayerTeamSelect", "True"),
	// Team kills
	n("teamKillLimit", "teamKills", secGame, "TeamKillLimit", 0, 100, "3"),
	f("teamKillGrace", "teamKills", secGame, "TeamKillGrace", 0, 60, "0.2"),
	n("teamKillReduceTime", "teamKills", secGame, "TeamKillReduceTime", 0, 3600, "90"),
	// Chat, voz e kill feed
	b("killFeed", "chat", secGame, "bKillFeed", "False"),
	b("killFeedSpectator", "chat", secGame, "bKillFeedSpectator", "True"),
	b("killerInfo", "chat", secGame, "bKillerInfo", "True"),
	b("killDistance", "chat", secGame, "bKillerInfoRevealDistance", "False"),
	b("deadSay", "chat", secGame, "bDeadSay", "False"),
	b("deadSayTeam", "chat", secGame, "bDeadSayTeam", "True"),
	b("deadVoice", "chat", secGame, "bVoiceAllowDeadChat", "False"),
	b("enemyHearsLocal", "chat", secGame, "bVoiceEnemyHearsLocal", "True"),
	b("intermissionVoice", "chat", secMulti, "bVoiceIntermissionAllowAll", "True"),
	// Partida
	b("mapVoting", "match", secMulti, "bMapVoting", "True"),
	b("useMapCycle", "match", secMulti, "bUseMapCycle", "True"),
	n("minimumPlayers", "match", secMulti, "MinimumPlayers", 1, 32, "1"),
	n("roundLimit", "match", secMulti, "RoundLimit", 1, 100, ""),
	n("winLimit", "match", secMulti, "WinLimit", 1, 100, ""),
	n("roundTime", "match", secMulti, "RoundTime", 60, 7200, ""),
	n("preRoundTime", "match", secMulti, "PreRoundTime", 0, 600, "10"),
	n("postRoundTime", "match", secMulti, "PostRoundTime", 0, 600, "15"),
	n("idleLimit", "match", secMulti, "IdleLimit", 0, 3600, "150"),
	n("idleLimitLow", "match", secMulti, "IdleLimitLowReinforcements", 0, 3600, "90"),
	b("bots", "match", secMulti, "bBots", "False"),
	n("botQuota", "match", secMulti, "BotQuota", 0, 32, ""),
	// Suprimentos
	n("initialSupply", "supply", secMulti, "InitialSupply", 0, 100, "15"),
	n("maximumSupply", "supply", secMulti, "MaximumSupply", 0, 100, "15"),
	b("supplyGain", "supply", secMulti, "bSupplyGainEnabled", "False"),
	// Votação para expulsar (Player Voting)
	func() RuleDef {
		r := b("voteKick", "votes", secTeam, "bVotingEnabled", "False")
		r.also = [][3]string{{secTeam, "TeamVoteIssues", voteIssue}}
		return r
	}(),
	n("voteMinPlayers", "votes", secVote, "MinimumPlayersRequired", 1, 64, "3"),
	f("votePassRatio", "votes", secVote, "VotePassRatio", 0, 1, "0.75"),
	n("voteBanDuration", "votes", secVote, "DefaultVoteBanDuration", -1, 604800, "120"),
	// Rede (Engine.ini)
	eng(n("tickRate", "network", secNet, "NetServerMaxTickRate", 10, 128, "")),
	eng(n("maxClientRate", "network", secNet, "MaxClientRate", 2500, 1000000, "")),
	eng(n("maxInternetClientRate", "network", secNet, "MaxInternetClientRate", 2500, 1000000, "")),
}

func ruleByID(id string) (RuleDef, bool) {
	for _, r := range Rules {
		if r.ID == id {
			return r, true
		}
	}
	return RuleDef{}, false
}

// NormalizeRule confere e padroniza um valor ("true"/"1" → "True"; números sem espaços). Vazio é aceito (padrão do
// jogo) e devolvido vazio.
func NormalizeRule(id, v string) (string, error) {
	r, ok := ruleByID(id)
	if !ok {
		return "", problems.Field("rules", "validation.unknownField")
	}
	v = strings.TrimSpace(v)
	if v == "" {
		return "", nil
	}
	bad := problems.Field("rules", "validation.rule", "rule", id,
		"min", strconv.FormatFloat(r.Min, 'f', -1, 64), "max", strconv.FormatFloat(r.Max, 'f', -1, 64))
	switch r.Kind {
	case KindBool:
		switch strings.ToLower(v) {
		case "true", "1":
			return "True", nil
		case "false", "0":
			return "False", nil
		}
		return "", problems.Field("rules", "validation.ruleBool", "rule", id)
	case KindInt:
		x, err := strconv.Atoi(v)
		if err != nil || float64(x) < r.Min || float64(x) > r.Max {
			return "", bad
		}
		return strconv.Itoa(x), nil
	default:
		x, err := strconv.ParseFloat(strings.Replace(v, ",", ".", 1), 64)
		if err != nil || x < r.Min || x > r.Max {
			return "", bad
		}
		return strconv.FormatFloat(x, 'f', -1, 64), nil
	}
}

// validateRules confere as regras do perfil.
func (p Profile) validateRules() error {
	for id, v := range p.Rules {
		if _, err := NormalizeRule(id, v); err != nil {
			return err
		}
	}
	minE, maxE := p.Rules["minimumEnemies"], p.Rules["maximumEnemies"]
	if minE != "" && maxE != "" {
		a, _ := strconv.Atoi(minE)
		z, _ := strconv.Atoi(maxE)
		if a > z {
			return problems.Field("rules", "validation.ruleMinMax", "rule", "minimumEnemies")
		}
	}
	return nil
}

// CleanRules devolve as regras com valores padronizados e sem as vazias (nil continua nil: o app não cuida delas).
func CleanRules(in map[string]string) map[string]string {
	if in == nil {
		return nil
	}
	out := map[string]string{}
	for id, v := range in {
		if nv, err := NormalizeRule(id, v); err == nil && nv != "" {
			out[id] = nv
		}
	}
	return out
}

// ApplyRules grava as regras no Game.ini e no Engine.ini (conteúdo atual → novo). rules nil não muda nada.
// Opções do catálogo sem valor são apagadas; o resto dos arquivos (outras seções, comentários) fica igual.
func ApplyRules(game, engine string, rules map[string]string) (string, string) {
	if rules == nil {
		return game, engine
	}
	for _, r := range Rules {
		text := &game
		if r.File == FileEngine {
			text = &engine
		}
		v, err := NormalizeRule(r.ID, rules[r.ID])
		if err != nil {
			v = ""
		}
		if v == "" {
			*text, _ = ini.Remove(*text, r.Section, r.Key)
		} else {
			*text = ini.Set(*text, r.Section, r.Key, v)
		}
		for _, a := range r.also {
			if v == "True" {
				*text = ini.Set(*text, a[0], a[1], a[2])
			} else if cur, ok := ini.Get(*text, a[0], a[1]); ok && cur == a[2] {
				*text, _ = ini.Remove(*text, a[0], a[1])
			}
		}
	}
	return game, engine
}

// ReadRules lê do Game.ini e do Engine.ini as opções do catálogo que estão nos arquivos (valores inválidos ficam de
// fora). Devolve nil se nenhuma foi achada.
func ReadRules(game, engine string) map[string]string {
	out := map[string]string{}
	for _, r := range Rules {
		text := game
		if r.File == FileEngine {
			text = engine
		}
		if v, ok := ini.Get(text, r.Section, r.Key); ok {
			if nv, err := NormalizeRule(r.ID, v); err == nil && nv != "" {
				out[r.ID] = nv
			}
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// rulesText mostra as regras em linhas "id=valor" ordenadas (comparação com o servidor existente).
func rulesText(m map[string]string) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	lines := make([]string, 0, len(keys))
	for _, k := range keys {
		lines = append(lines, k+"="+m[k])
	}
	return strings.Join(lines, "\n")
}
