package logparse

import (
	"testing"

	"sandstorm-server-assist/desktop/internal/domain/problems"
)

// Linhas reais do Insurgency.log (09/11/2026).
func TestParseRealLines(t *testing.T) {
	cases := []struct {
		line string
		want Event
	}{
		{`[2026.11.09-04.39.24:437][  0]LogModioGame: ModSubsystem: mod.io initialization complete with result: A opera??o foi conclu?da com ?xito.`, Event{Kind: ModioReady}},
		{`[2026.11.09-04.39.25:323][  0]LogNet: GameNetDriver IpNetDriver_2147482532 IpNetDriver listening on port 27102`, Event{Kind: Listening, Port: 27102}},
		{`[2026.11.09-04.39.26:320][  0]LogModioGame: ModSubsystem: ModManager enabled`, Event{Kind: ModsEnabled}},
		{`[2026.11.09-04.39.26:440][  1]LogModioGame: ModSubsystem: OnModActivatedFromCloud: No Restricted Area Mutator [1.15][98145] bFromGameServer = true`, Event{Kind: ModActivated, ModID: 98145, ModName: "No Restricted Area Mutator [1.15]"}},
		{`[2026.11.09-04.39.26:465][  1]LogModioGame: ModSubsystem: OnModActivatedFromCloud: Containment[2123731] bFromGameServer = true`, Event{Kind: ModActivated, ModID: 2123731, ModName: "Containment"}},
		{`[2026.11.09-04.39.27:622][ 44]LogModioGame: ModSubsystem: ModManager Event received. ModId: 1134624, ModEvent: BeginUninstall, Message: ok`, Event{Kind: ModEvent, ModID: 1134624, Action: "BeginUninstall"}},
		{`[2026.11.09-04.36.32:282][  0]LogModioGame: ModSubsystem: Mod added, id: 2853082 to PendingChanges with state Removed`, Event{Kind: ModPending, ModID: 2853082, Action: "Removed"}},
		{`[2026.11.09-04.39.26:483][  1]LogGameMode: ProcessServerTravel: Containment?Scenario=Scenario_Containment_Checkpoint_Security`, Event{Kind: Travel, Map: "Containment", Scenario: "Scenario_Containment_Checkpoint_Security"}},
		{`[2026.11.09-04.39.24:437][  0]LogLoad: LoadMap: /Game/Maps/Farmhouse/Farmhouse?Name=Player?Scenario=Scenario_Farmhouse_Checkpoint_Security`, Event{Kind: LoadMap, Map: "Farmhouse", Scenario: "Scenario_Farmhouse_Checkpoint_Security"}},
		{`[2026.11.09-04.39.26:411][  0]LogGameState: Match State Changed from EnteringMap to LoadingAssets`, Event{Kind: MatchState, From: "EnteringMap", To: "LoadingAssets"}},
		{`[2026.11.09-04.36.46:759][ 21]LogWindows: Error: StaticMesh /Game/Environment/SM_MetalWindow_01d: Serial size mismatch. Got 1, Expected 2`, Event{Kind: Problem, Problem: problems.ModConflict}},
		{`[2026.11.09-04.36.46:759][ 21]LogWindows: Error: Fatal error: [File:Unknown] [Line: 3237]`, Event{Kind: Problem, Problem: problems.ServerCrashed}},
		{`LogModio: [Error][Http] 401 Invalid security code`, Event{Kind: Problem, Problem: problems.ModioCodeInvalid}},
		{`[2026.11.09-04.39.21:792][  0]LogGameStatsServer: Display: GameStatsServer INIT: No authentication provided, -GameStatsToken= required for statistics collection.`, Event{Kind: GameStatsNoKey}},
	}
	for _, c := range cases {
		got, ok := Parse(c.line)
		if !ok || got != c.want {
			t.Errorf("Parse(%q)\n got %+v, %v\nwant %+v", c.line, got, ok, c.want)
		}
	}
	if _, ok := Parse(`LogPakFile: Display: Mounting pak file ../../../Insurgency/Content/Paks/pakchunk28-WindowsServer.pak.`); ok {
		t.Error("linha de montagem de pak não deveria ser evento")
	}
}

func TestLevel(t *testing.T) {
	if Level(`[x][ 0]LogMapVoteManager: Error: Not adding 'Scenario_X' to map vote, no asset found.`) != "error" {
		t.Error("erro não reconhecido")
	}
	if Level(`[x][ 0]LogLoad: Warning: UINSGameInstance::OnPreLoadMap: Invalid NextURL`) != "warning" {
		t.Error("aviso não reconhecido")
	}
	if Level(`LogNet: GameNetDriver listening on port 27102`) != "info" {
		t.Error("info não reconhecido")
	}
}

func TestHarmlessErrorsCountAsWarnings(t *testing.T) {
	harmlessLines := []string{
		`[x][ 9]LogCosmetics: Error: GetSelectedCosmeticFromItemId called with id: 123 was invalid!`,
		`[x][ 9]LogStreaming: Error: Couldn't find file for package /Game/Game/Actors/Weapons/Upgrades/Underbarrel/BP_UG_Foregrip_Flashlight requested by async loading code.`,
		`[x][ 9]LogDemo: Error: UDemoNetDriver::TickFlush: ReplayStreamer ERROR: ServiceUnavailable`,
		`[x][ 9]LogRedstoneSocket: Error: [network/stun] Integrity check failed`,
	}
	for _, l := range harmlessLines {
		if Level(l) != "warning" {
			t.Errorf("should be warning: %s", l)
		}
	}
	// Rotação com cenário que não existe continua erro: é configuração a corrigir.
	if Level(`[x][ 9]LogMapVoteManager: Error: Not adding 'Scenario_West_Checkpoint_Insurgents' to map vote, no asset found.`) != "error" {
		t.Error("map vote error should stay an error")
	}
}

func TestModMapPackagingErrorsAreWarnings(t *testing.T) {
	l := `[2026.11.09-22.56.20:123][ 5]LogStreaming: Error: Could not find outer LandscapeComponent_12 to create LandscapeHeightfieldCollisionComponent_3`
	if Level(l) != "warning" {
		t.Errorf("should be warning: %s", l)
	}
}
