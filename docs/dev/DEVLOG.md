# Dev log

Short notes of each development session, newest last. They help a new session pick up the work. Decisions are in
[DECISIONS.md](../DECISIONS.md); the plan is in [ROADMAP.md](ROADMAP.md).

## What the app was asked to do (2026-11-09)

1. Help set up and start an Insurgency: Sandstorm dedicated server, following a server that was first configured by
   hand (those private notes stay out of git).
2. Guides with links for what the user has to get elsewhere (SteamID, GSLT, GameStats, mod.io).
3. Visible download progress (SteamCMD, server, mods).
4. Errors always with what happened, the likely cause and a solution.
5. Server logs hidden by default, opened with one click.
6. Sensitive data kept safe.
7. Backup and restore of everything needed to run the server on another PC.
8. A public project with a license that protects the authors from misuse.
9. Interface in en-US (default), pt-PT, es-ES, ru and uk.
10. Color themes and layout skins.

## 2026-11-09 · 0.1.0 to 0.1.2

- **0.1.0:** all phases done and verified (Go back end, React screens, 5 languages, themes and skins, app icon,
  portable zip, per-user installer, signed automatic updates, terms and privacy acceptance).
- **Existing server (ADR-013):** the fields are filled from the server's own `LogInit: Command Line:` plus
  `Game.ini`, `Admins.txt` and `MapCycle.txt`, asking field by field when the app and the server differ. Found in a real
  test: the config pointed to the parent folder of the server; the app now picks the subfolder by itself. "Running
  since" now uses the real process start time.
- **Download page (ADR-016)** prepared and later published on the author's server.
- **0.1.1 (ADR-017):** on the first real run only 2 of 7 mods were mounted. Cause: with new subscriptions the server
  mounts only the new mods. The app now restarts once by itself; it also fixes repeated `Scenario=` in the map
  cycle, hides removed mods and treats known game noise as warnings.
- **0.1.2 (ADR-018):** logs per server run, a selector for old runs, optional cleanup of old logs (off by default) and
  a single app instance (two open copies had overwritten `config.json`).
- The public repository was created on GitHub (first commit 0.1.2) after removing server paths from docs and tests
  and a secret scan.

## 2026-11-10 · 0.1.3 to 0.1.7

- **Map change by RCON (ADR-019)**, message of the day and join/leave messages (ADR-020), neutral network texts and
  an internet port test service (ADR-021), and 14 more server settings (ADR-022): game rules, mutators, day/night,
  players panel, scheduled messages, automatic restart, Discord, start with Windows, Firewall rule, network adapter,
  tick rate, profiles and mod list. Published as **0.1.3**.
- Lesson from publishing: after changing the web server config, the container must be recreated, or it keeps the old
  file.
- **0.1.4 (ADR-023):** mutators found inside the installed mods' `.pak` files.
- **0.1.5 (ADR-024):** mod map picker; message formatting (later removed).
- **0.1.6 (ADR-025):** "Show players" with an empty server turned the app black (Go sent `null`). Lists are never null
  now, and screen errors show a message instead of a black window.
- **0.1.7 (ADR-026):** a reopened app recognises the server it started.

## 2026-11-11 · 0.1.8 to 0.1.10

- **0.1.8 (ADR-027):** after a map change the Dashboard stayed on "loading the map". The wait for the mod map now ends
  at the first travel.
- **0.1.9 (ADR-028):** the Dashboard showed the old map after changing to a mod map; the travel line has the map path.
- **0.1.10 (ADR-029 to 031):** map change uses the mods loaded in the run, with map and mode lists; map cycle entries
  without their mod are flagged; messages are plain text (the game showed the tags). Built with Go 1.27.4 after
  `govulncheck` found GO-2026-6859 (HTTP/1 server crash); the port checker service was rebuilt with it too.

## 2026-11-11 · 0.2.0 (community release)

- The maintainer decided to make the app a free community project and to stop depending on their own server:
  - **Updates from GitHub Releases (ADR-032)**, with the old server as a bridge for apps up to 0.1.10.
  - **Port diagnosis from the PC (ADR-033)** instead of the internet test service, plus step-by-step instructions to
    confirm with a friend outside the network.
  - **Feedback page (ADR-034):** bug reports and ideas through GitHub issues; the app prepares a zip (masked log,
    optional screenshot) and opens the form filled in. Description and image are optional.
  - **Docs in English (ADR-035)**, README and user guide in the 5 app languages with a language bar; the development
    notes became this log and the roadmap; credits added.
- Screenshots taken from the dev app with empty temporary data and demo values.
- **Published** the same day: GitHub Release `v0.2.0` and, as the bridge, the old download server. Installed apps
  up to 0.1.10 update from the old server to 0.2.0 and then use GitHub.

## Cautions for future sessions

- The maintainer's real server folder has secrets in plain text (`Game.ini`, notes next to the `.bat`) and the mod.io
  session token (`user.json`): never display their values.
- Never start, stop or restart the maintainer's real game server unless asked.
- The update signing key (`~/.sandstorm-server-assist/update-signing-key.pem`) never goes into git or logs; without
  it installed apps reject updates.
- Windows Bash: text with backticks inside `node -e "..."` or quoted heredocs breaks; write a file instead.
