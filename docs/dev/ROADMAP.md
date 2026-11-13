# Roadmap

> Living document, updated at each step. Last update: 2026-11-13 (0.2.4 prepared).
> See also: [Architecture](../ARCHITECTURE.md) · [Decisions](../DECISIONS.md) · [Dev log](DEVLOG.md) ·
> [Releasing](../RELEASING.md)

## Status

| Phase                                | What it delivers                                                                                           | Status  |
| ------------------------------------ | ---------------------------------------------------------------------------------------------------------- | ------- |
| 0. Foundation                        | Docs, Wails 2 + React, scripts, lint, tests, `security-audit` skill                                        | ✅ Done |
| 1. Profile, vault and settings       | Profile, DPAPI vault, Server screen with "How to get it?" guides, import an existing server                | ✅ Done |
| 2. Install with progress             | SteamCMD download (Valve signature checked), install/update/verify with bar, speed and time left           | ✅ Done |
| 3. Start and follow                  | Command line in Go, start/stop/restart, phases from the log, server started outside the app                | ✅ Done |
| 4. mod.io                            | Request a code, authorize (back to `none` by itself), API key from the log, mods with progress             | ✅ Done |
| 5. Errors and logs                   | Error catalog (what/cause/solution + links), log drawer closed by default, filters, mask                   | ✅ Done |
| 6. Backup and restore                | `.ssbackup` (Argon2id + AES-256-GCM), create, check and restore on another folder or PC                    | ✅ Done |
| 7. Languages                         | en-US (default), pt-PT, es-ES, ru, uk; a test checks every key                                             | ✅ Done |
| 8. Public project                    | Apache 2.0, NOTICE, DISCLAIMER, CONTRIBUTING, SECURITY, no personal data in code                           | ✅ Done |
| 9. Appearance                        | 5 color themes and 4 layout skins (ADR-011)                                                                | ✅ Done |
| 10. Distribution                     | App icon, portable zip with the legal files                                                                | ✅ Done |
| 11. Existing server (ADR-013)        | Fill the fields from a configured server (log + ini/txt), ask to keep or overwrite                         | ✅ Done |
| 12. Installer and updates            | Per-user installer, signed automatic updates, terms and privacy acceptance (ADR-014/015)                   | ✅ Done |
| 13. Download page                    | Static page on the author's server (ADR-016), now a bridge only                                            | ✅ Done |
| 14–15. First real use (0.1.1–0.1.2)  | Restart once when new mods leave others out (ADR-017), logs per run and cleanup (ADR-018)                  | ✅ Done |
| 16–19. Server control (0.1.3)        | Map change by RCON, messages, neutral network texts, game rules, players, automation, Discord, profiles    | ✅ Done |
| 20–26. Mods and fixes (0.1.4–0.1.10) | Mod mutators and maps, screen error protection, own-server recognition, map change fixes (ADR-023 to 031)  | ✅ Done |
| 27. Community release (0.2.0)        | GitHub Releases updates, local port diagnosis, Feedback page, English docs in 5 languages (ADR-032 to 035) | ✅ Done |
| 28. Installation feedback (0.2.1)    | Progress, result and errors next to the button; task shown as soon as it starts; Check SteamCMD (ADR-036)  | ✅ Done |
| 29. Delete a mod (0.2.1)             | Delete a downloaded mod so the server downloads it again, now or on the next start (ADR-037)               | ✅ Done |
| 30. Download new mods now (0.2.1)    | Subscribed mods not on disk: start or restart the server to download them (ADR-038)                        | ✅ Done |
| 31. Mods with problems (0.2.1)       | Conflicts between mods, mod of a crash message, server error names the mod (ADR-039)                       | ✅ Done |
| 32. Every player listed (0.2.2)      | `listplayers` read as one row of records; players by score, bots folded                                    | ✅ Done |
| 33. Server rules offered (0.2.2)     | Game rules tab offers the rules the server's Game.ini already sets (ADR-040)                               | ✅ Done |
| 34. Refused mutators (0.2.2)         | Warn when the server refuses a mutator on the current map, read from its log (ADR-041)                     | ✅ Done |
| 35. SteamCMD folder (0.2.3)          | Refuse a SteamCMD folder inside the server folder, and the other way around (ADR-042)                      | ✅ Done |

Details of each version are in the [dev log](DEVLOG.md) and in the [changelog](../../CHANGELOG.md).

## Phase 27 · Community release (0.2.0)

- **Updates from GitHub Releases** (ADR-032): `adapters/updater` reads `releases/latest/download/latest.json` and
  downloads `releases/download/v<version>/<installer>`; `netx` allows GitHub hosts; `npm run release:publish`
  publishes the release with notes from `CHANGELOG.md`.
- **Port diagnosis from the PC** (ADR-033): `domain/netcheck` (rules), `adapters/netdiag` (STUN, UPnP read only,
  A2S, Steam `GetServersAtAddress`); new screen with verdict, items and the "test with a friend" steps.
- **Feedback** (ADR-034): `domain/report`, `ReportAPI`, Feedback page, GitHub issue forms in `.github/ISSUE_TEMPLATE`.
- **Docs** (ADR-035): all docs in English; README and user guide in 5 languages with a language bar; `npm test`
  checks the translations' structure; credits in README, About and `NOTICE`.
- Terms and privacy version 2026-11-11 (new hosts: GitHub, STUN, Steam Web API).
- Verified: Go and front-end tests (new: netcheck rules, STUN parser, UPnP against a fake router, A2S against a fake
  server, report zip masking, i18n codes), lint, typecheck, audits.
- Real network check (`SSA_NETDIAG_E2E=1 go test ./internal/adapters/netdiag/ -run TestRealNetwork -v`): STUN
  found the public IP; the maintainer's router does not answer UPnP, which shows as information. Screenshots taken
  from the dev app with empty temporary data and demo values (no real IP, server or player).

- **Published** on 2026-11-11: GitHub Release `v0.2.0` (installer, portable zip, `latest.json`, `SHA256SUMS.txt`)
  and, as the bridge, on the old download server. Both `latest.json` point to 0.2.0 and serve the same installer
  (SHA-256 checked after download).

## Phase 28 · Installation feedback (0.2.1)

- `installer.begin` notifies the screen at once; `State` gains `last`, `validate`, `upToDate` and `present`;
  `steamcmd.Parse` reads "already up to date".
- `InstallPage`: `TaskPanel` in the SteamCMD and server cards (progress with Cancel, result, error), scrolled into
  view on start; the old "Progress" card at the top is gone. **Check SteamCMD** when it is present. ADR-036.
- Verified: Go tests (notify on start, up-to-date parsing) and front-end tests of the panel; not run against the
  real server folder.

## Phase 29 · Delete a downloaded mod (0.2.1)

- `modio.RemoveInstalled` (state.json entry and folder), `app/modfiles` (stop, delete, start), `ServerAPI.DeleteMod`;
  delete button with the confirmation on each downloaded mod's row. ADR-037.
- Verified: adapter test with the real `state.json` shape, use case tests (stopped, running, failure, busy) and
  front-end tests of the row; not run on the real mods.

## Phase 30 · Download new mods now (0.2.1)

- `modfiles.DownloadNew` (start, or stop and start), `ServerAPI.DownloadNewMods`; `PendingMods` above the Mods table
  with the confirmation for a running server. ADR-038.
- Verified: use case and front-end tests; not run on the real server.

## Phase 31 · Mods with problems (0.2.1)

- `domain/pak`: `ModernFiles` (all files with their size, from the encoded entries), `LegacyFiles` (formats 8 and 9,
  with the UE 4.22 variant), `FooterAny`, `ContentsOf`; fuzz tests. `adapters/modio`: `openPak`, `pakFiles`,
  `InstalledFiles`; old mods also get mutators and maps.
- `domain/modcheck`: `Conflicts` (same asset, different sizes), `Assets` (objects named in a message), `Owners`,
  `Check`. `logparse` keeps the object of a "Serial size mismatch" line; the runner adds the file and the mod to the
  `mod-conflict` problem (`ModOwners`). `ServerAPI.CheckMods` and `CheckCrash`.
- Mods screen: **Mods with problems** card (pairs in conflict, unreadable mods) and **Find the mod** for a pasted crash
  message. ADR-039.
- Verified: unit and fuzz tests; real check on the maintainer's PC (10 installed mods read in about 10 ms, 5 pairs in
  conflict, and the crash of the maintainer's game traced to the Daybreak mod).

- **0.2.1 published** on 2026-11-11: GitHub Release `v0.2.1` and the bridge server; both `latest.json` serve the
  same installer (SHA-256 checked after download).

## Phase 32 · Every player listed (0.2.2)

- The server sends `listplayers` with every player on one line, one 5-cell record after another; the app read one
  player per line and showed only the first. `admin.ParsePlayers` now reads the cells as one sequence and starts a
  record at each numeric ID (one player per line still works).
- Players card: players by score, the bots folded in one row ("Show the 19 bots").
- Verified: parser test with the real one-line shape, screen test, and the real server's answer read by the app
  (3 players and 19 bots).

## Phase 33 · Server rules offered (0.2.2)

- `config.ServerRules` and `ConfigAPI.ServerRules` (read the rules of `Game.ini`/`Engine.ini`); the Game rules tab,
  with the option off, shows how many the server sets and **Use the server rules**. ADR-040.
- Verified: config test (no folder, no file, two rules, profile unchanged) and screen tests.

## Phase 34 · Refused mutators (0.2.2)

- `logparse` `MutatorList`/`MutatorRefused` and `MutatorName`; runner `Status.RefusedMutators`; warning on the
  Dashboard and in the mutator picker, with the refused ones marked. ADR-041.
- Verified: parser and runner tests with the real lines, a real log replayed (Ultralethal refused), screen test.

- **0.2.2 published** on 2026-11-13: GitHub Release `v0.2.2` and the bridge server; both `latest.json` serve the
  same installer (SHA-256 checked after download).

## Phase 35 · SteamCMD folder inside the server folder (0.2.3)

- `steamcmd.NestedDirs` and `steamcmd.WrongInstallDir`; `config.SetServerDir`/`SetSteamCMDDir` refuse nested
  folders, `installer.Install` checks again before running, `runSteamCMD` turns the SteamCMD warning into the
  problem even when SteamCMD then says "Success!", and there is no second attempt for it. `installer.Info` has
  `foldersNested` for a warning on the Installation screen. New problem `folders-nested`. ADR-042.
- Verified: domain, config and installer tests; the maintainer's real SteamCMD log shows the warning followed by a
  successful install into `steamcmd\steamapps\common`.

- **0.2.3 published** on 2026-11-13: GitHub Release `v0.2.3` and the bridge server; both `latest.json` serve the
  same installer (SHA-256 checked after download).

## 0.2.4 · Server mods only, default rotation (ADR-043)

- `modio.ReadInstalled` keeps `SubscriptionCount` (`Subscribers`); `modio.ServerUserID` reads only the id of the
  server account; `ForServer`/`ServerInstalled` keep the mods of that account (none before mod.io is authorized;
  mods without the list stay). Every server-side reader uses `ServerInstalled`.
- `server.DefaultMapCycle` (mode of the start map, both sides in Checkpoint, from the server `.pak` scenarios or
  the `Maps` list); `ConfigAPI.DefaultMapCycle`; `DefaultCycle` on the Maps tab fills a new server once and offers
  the default rotation with a confirmation.
- Verified: unit tests; on the maintainer's PC the game's 42 mods are hidden (all subscribed by the player account)
  and the default rotation from the real server files has 20 maps (40 lines, Hold and Trainyard included).

## 0.2.4 · Error card layout

- The `folders-nested` card overflowed its column with a horizontal scroll: the long paths in its text have no
  spaces and the `1fr` column would not shrink below them. `.errorcard` wraps anywhere and is a size container;
  its grid uses `minmax(0, 1fr)` and stacks labels over the text below 520 px; grid items and cards get
  `min-width: 0`; `.grid--2` columns are at most 100% wide. Checked in the dev app at 1060 and 800 px.

## Bridge from the old download server

Apps up to 0.1.10 only know the old download server. Plan:

1. ✅ 0.2.0 published on GitHub (`npm run release:publish`) **and** on the old server (`npm run site:publish`), so
   old apps update to 0.2.0. From then on they update from GitHub.
2. Keep the old server's `latest.json` and port checker online for about 3 months (until 2027-02).
3. Then remove the `/sandstorm/` route and the `sandstorm-check` service from the author's server, delete
   `checker/` and `scripts/publish-site.mjs` and `site/` from this repository, and update ADR-016/021 and
   `docs/RELEASING.md`.

## Pending (needs the maintainer)

- Human checks in the native window: choose folders, import the real server, move the RCON password to the vault,
  start/stop, backup/restore with the Windows dialogs.
- Live checks with players: kick/ban, map change with day/night, game rules, Discord webhook, Firewall rule (asks
  for administrator), daily restart, join/leave messages.
- 0.2.0: run "Diagnose ports" with the server running, and send a test report from the Feedback page.
- Screenshots for the README and the user guide (`docs/images/`).

## Ideas

- Detect mod conflicts before starting (repeated paths inside the `.pak` files).
- Open the router ports through UPnP from the app (today the diagnosis only reads).
- CI on GitHub (tests, lint and audits on every push).
- Translate the code comments to English.
