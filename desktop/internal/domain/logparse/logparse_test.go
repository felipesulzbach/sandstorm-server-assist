package logparse

import (
	"testing"
	"time"

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
		`[2026.11.11-02.23.16:423][874]LogWorldSettings: Error: No "Max Raycast Distance" defined, falling back to 1km (100000 units).`,
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

// Entradas e saídas de jogadores, no formato do log real de 10/11/2026 (nomes e ids inventados).
func TestParsePlayers(t *testing.T) {
	at := time.Date(2026, 11, 10, 2, 20, 41, 435*int(time.Millisecond), time.UTC)
	login := `[2026.11.10-02.20.41:435][337]LogNet: Login request: ?InitialConnectTimeout=30?ConnectionType=2?Name=[TAG] Alpha Bravo userId: SteamNWI:76561190000000001 platform: SteamNWI`
	got, ok := Parse(login)
	if !ok || got.Kind != PlayerLogin || got.Player != "[TAG] Alpha Bravo" || got.PlayerID != "SteamNWI:76561190000000001" || !got.Time.Equal(at) {
		t.Fatalf("login = %+v %v", got, ok)
	}
	joined := `[2026.11.10-02.20.43:924][487]LogNet: Join succeeded: [TAG] Alpha `
	if got, ok := Parse(joined); !ok || got.Kind != PlayerJoined || got.Player != "[TAG] Alpha" {
		t.Fatalf("joined = %+v %v", got, ok)
	}
	left := `[2026.11.10-03.27.05:054][915]LogNet: UNetConnection::Close: [UNetConnection] RemoteAddr: 192.0.2.10:7777, Name: IpConnection_2147481362, Driver: GameNetDriver IpNetDriver_2147482533, IsServer: YES, PC: INSPlayerController_2147481358, Owner: INSPlayerController_2147481358, UniqueId: SteamNWI:76561190000000001, Channels: 229, Time: 2026.11.10-03.27.05`
	if got, ok := Parse(left); !ok || got.Kind != PlayerLeft || got.PlayerID != "SteamNWI:76561190000000001" {
		t.Fatalf("left = %+v %v", got, ok)
	}
	// Replay (DemoNetDriver) e conexão que caiu antes de entrar (PC: NULL) não são saída de jogador.
	for _, l := range []string{
		`[2026.11.10-02.22.55:100][459]LogNet: UNetConnection::Close: [UNetConnection] RemoteAddr: 192.0.2.10:7777, Name: DemoNetConnection_1, Driver: DemoNetDriver INSDemoNetDriver_2, IsServer: YES, PC: INSDemoSpectator_3, Owner: INSDemoSpectator_3, UniqueId: INVALID, Channels: 205, Time: 2026.11.10-02.22.55`,
		`[2026.11.10-01.45.59:396][292]LogNet: UNetConnection::Close: [UNetConnection] RemoteAddr: 192.0.2.10:7777, Name: IpConnection_1, Driver: GameNetDriver IpNetDriver_2, IsServer: YES, PC: NULL, Owner: NULL, UniqueId: SteamNWI:76561190000000002, Channels: 2, Time: 2026.11.10-01.45.59`,
	} {
		if got, ok := Parse(l); ok && got.Kind == PlayerLeft {
			t.Fatalf("not a player leaving: %q", l)
		}
	}
	if !Stamp("no time here").IsZero() {
		t.Fatal("line without time must give zero")
	}
}
