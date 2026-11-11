# Roadmap

> Living document, updated at each step. Last update: 2026-11-11 (0.2.0 prepared).
> See also: [Architecture](../ARCHITECTURE.md) · [Decisions](../DECISIONS.md) · [Dev log](DEVLOG.md) ·
> [Releasing](../RELEASING.md)

## Status

| Phase                                | What it delivers                                                                                           | Status              |
| ------------------------------------ | ---------------------------------------------------------------------------------------------------------- | ------------------- |
| 0. Foundation                        | Docs, Wails 2 + React, scripts, lint, tests, `security-audit` skill                                        | ✅ Done             |
| 1. Profile, vault and settings       | Profile, DPAPI vault, Server screen with "How to get it?" guides, import an existing server                | ✅ Done             |
| 2. Install with progress             | SteamCMD download (Valve signature checked), install/update/verify with bar, speed and time left           | ✅ Done             |
| 3. Start and follow                  | Command line in Go, start/stop/restart, phases from the log, server started outside the app                | ✅ Done             |
| 4. mod.io                            | Request a code, authorize (back to `none` by itself), API key from the log, mods with progress             | ✅ Done             |
| 5. Errors and logs                   | Error catalog (what/cause/solution + links), log drawer closed by default, filters, mask                   | ✅ Done             |
| 6. Backup and restore                | `.ssbackup` (Argon2id + AES-256-GCM), create, check and restore on another folder or PC                    | ✅ Done             |
| 7. Languages                         | en-US (default), pt-PT, es-ES, ru, uk; a test checks every key                                             | ✅ Done             |
| 8. Public project                    | Apache 2.0, NOTICE, DISCLAIMER, CONTRIBUTING, SECURITY, no personal data in code                           | ✅ Done             |
| 9. Appearance                        | 5 color themes and 4 layout skins (ADR-011)                                                                | ✅ Done             |
| 10. Distribution                     | App icon, portable zip with the legal files                                                                | ✅ Done             |
| 11. Existing server (ADR-013)        | Fill the fields from a configured server (log + ini/txt), ask to keep or overwrite                         | ✅ Done             |
| 12. Installer and updates            | Per-user installer, signed automatic updates, terms and privacy acceptance (ADR-014/015)                   | ✅ Done             |
| 13. Download page                    | Static page on the author's server (ADR-016), now a bridge only                                            | ✅ Done             |
| 14–15. First real use (0.1.1–0.1.2)  | Restart once when new mods leave others out (ADR-017), logs per run and cleanup (ADR-018)                  | ✅ Done             |
| 16–19. Server control (0.1.3)        | Map change by RCON, messages, neutral network texts, game rules, players, automation, Discord, profiles    | ✅ Done             |
| 20–26. Mods and fixes (0.1.4–0.1.10) | Mod mutators and maps, screen error protection, own-server recognition, map change fixes (ADR-023 to 031)  | ✅ Done             |
| 27. Community release (0.2.0)        | GitHub Releases updates, local port diagnosis, Feedback page, English docs in 5 languages (ADR-032 to 035) | 🔄 Ready to publish |

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

## Bridge from the old download server

Apps up to 0.1.10 only know the old download server. Plan:

1. Publish 0.2.0 on GitHub (`npm run release:publish`) **and** on the old server (`npm run site:publish`), so old apps
   update to 0.2.0. From then on they update from GitHub.
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
