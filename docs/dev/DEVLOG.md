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

## 2026-11-11 · 0.2.1 (installation feedback)

- The maintainer clicked "Check for updates" and saw nothing, then "Verify files" showed "Wait for the current task"
  at the top of the page, out of view; "Download SteamCMD" also seemed to do nothing. The installer now tells the
  screen when a task starts, and progress, result and errors show in the card of the button (ADR-036). With SteamCMD
  present the button is "Check SteamCMD" and explains that SteamCMD updates itself.
- Also in 0.2.1: the maintainer asked to delete a downloaded mod so the server downloads it again, right away or on
  the next start (ADR-037). The app removes the mod from mod.io's `state.json` and deletes its folder; with the
  server running, "Delete and download again" stops it, deletes and starts it.
- Also in 0.2.1: new subscriptions showed as "will be downloaded" with no way to download them. The Mods screen now
  offers to start the server (or restart it, with a confirmation) so mod.io downloads them (ADR-038).

## 2026-11-11 · 0.2.1 (mods with problems)

- The maintainer's game crashed on maps that used to work: "Serial size mismatch" on a `/Game/...` static mesh.
  Reading the index of every installed mod showed that only the Daybreak mod ships that file, with a copy made for
  another game version; the app was not involved. The maintainer asked whether the app can find this kind of
  problem, ideally before starting the server (it had happened before when loading two map mods).
- Done (ADR-039): the Mods screen compares the mods' files (same game file in different sizes = conflict), finds the
  mod of a pasted crash message, the server's "Serial size mismatch" error names the file and the mod, and mods
  with the older `.pak` format are read too.
- Limit: a mod whose file does not match the game version (like Daybreak here) does not show as a conflict, because
  no other mod has the file; it shows when the crash message is checked.
- Published the same day as `v0.2.1` (GitHub and the bridge server), together with the 0.2.1 work from the other PC
  (installation feedback, delete a mod, download new mods now).

## 2026-11-12 · 0.2.2 (every player listed)

- The maintainer saw only one player on the Players card with three on the server. The real `listplayers` answer
  (read by the app, names and addresses masked) has every player and bot on one line, 5 cells each; the parser read
  one player per line. Fixed and tested with that shape; the Players card now lists players by score and folds the
  bots (19 in that match).
- Also in 0.2.2: "Set the game rules here" was off on a server with rules in `Game.ini`. The import was not at
  fault (run on a copy of the maintainer's config: it turns the option on with the 15 rules); the rules arrived
  after the first import. The Game rules tab now offers the server's rules with one click (ADR-040).
- Also in 0.2.2: asked whether Bullet Sponge and Ultralethal override each other. The logs show the server refusing
  Ultralethal every time the two were together, without a reason; the app now warns about refused mutators, read
  from the log (ADR-041). The maintainer then started the server without Bullet Sponge: Ultralethal was refused
  again, so Bullet Sponge was not the cause.
- Published on 2026-11-13 as `v0.2.2` (GitHub and the bridge server).

## 2026-11-13 · 0.2.3 (SteamCMD folder)

- On a new PC the maintainer chose the server folder and, inside it, the SteamCMD folder. "Install server" ended
  with "SteamCMD could not install". The SteamCMD log said "Please set the game install path to something other
  than the Steam install folder", then "Missing configuration" (first try) and, on the automatic second try,
  "Success!", but into `steamcmd\steamapps\common\sandstorm_server`. The app did not find the server in its
  folder and showed the generic error.
- Fixed (ADR-042): nested folders are refused when chosen and before installing, the Installation screen warns about
  them, and the SteamCMD warning becomes the new `folders-nested` problem.
- Published the same day as `v0.2.3` (GitHub and the bridge server).

## 2026-11-13 · 0.2.4 (error card layout)

- The maintainer saw the new `folders-nested` card overflow its column (horizontal scroll). Cause: long paths
  without spaces in the text and a grid column that would not shrink. Fixed in CSS and checked in the dev app at
  two window widths.

- Also in 0.2.4: a fresh install showed NoRestrictedArea among the mutators before mod.io was set up. The mod.io
  folder is shared with the game on that PC; `state.json` lists the subscribed accounts of each mod, so the app
  now keeps only the server account's mods (ADR-043). The maintainer also asked for a default map rotation:
  a new server now starts with the game maps in the mode of the start map.
- Published the same day as `v0.2.4` (GitHub and the bridge server).

## 2026-11-13 · 0.2.5 (mod.io code)

- On the new server, "Authorize" failed with mod.io 401 "Invalid security code" (`error_ref` 11014). mod.io's
  documentation says the exchange must use the API key of the original request; the server exchanges with the
  game's own key, but the app had requested the code with the key typed in Mods (the new mod.io plugin no longer
  writes the key in the server log, so it could not be found there). The game key is in
  `Insurgency/Config/DefaultGame.ini` inside `pakchunk0`; the app now reads it and requests the code with it
  (ADR-044).
- Published the same day as `v0.2.5` (GitHub and the bridge server).

## 2026-11-13 · 0.2.6 (layout)

- The maintainer saw the Dashboard start-up steps and the scenario name go past the status card. Cause: grid
  columns sized `1fr` (minimum = content) and a fixed 5-column step grid. Fixed in every screen, not only the
  Dashboard, and checked with a scan script that walks all screens and tabs through the DevTools protocol and
  reports any element outside its card (kept outside the repository, with demo data only).

## Cautions for future sessions

- The maintainer's real server folder has secrets in plain text (`Game.ini`, notes next to the `.bat`) and the mod.io
  session token (`user.json`): never display their values.
- Never start, stop or restart the maintainer's real game server unless asked.
- The update signing key (`~/.sandstorm-server-assist/update-signing-key.pem`) never goes into git or logs; without
  it installed apps reject updates.
- Windows Bash: text with backticks inside `node -e "..."` or quoted heredocs breaks; write a file instead.
