# Architecture decisions (ADRs)

Short records: context, decision and consequence. New decisions go at the end. When a decision replaces an older
one, the older one says so.

## ADR-001 · Desktop app with Wails 2 (Go + React)

- **Context:** the app runs Windows processes (SteamCMD, the game server), reads and writes server files, keeps
  secrets with Windows encryption and shows screens with progress.
- **Decision:** Wails **2.17.1** + Go 1.27 + React 19 + Vite. The Wails CLI runs from the same version pinned in
  `go.mod` (`scripts/desktop.mjs`), with no global install. Electron was ruled out: it would ship a whole
  Chromium and Node for no gain here.
- **Consequence:** a small app (it uses the WebView2 that comes with Windows 11), with no server and no database.
  Rules in Go can be tested without a window.

## ADR-002 · Secrets in the Windows vault (DPAPI), not in `.ini` or `.bat` files

- **Context:** a hand-made server usually keeps the RCON password in plain text in `Game.ini`, and the GSLT and
  GameStats tokens in plain text next to the `.bat`. Any copy of the folder takes the secrets with it.
- **Decision:** join password, RCON password, GSLT, GameStats token and mod.io e-mail live in
  `%APPDATA%\sandstorm-server-assist\vault.bin`, encrypted with DPAPI (`CryptProtectData`, Windows user scope, with
  the app's own entropy). Only that user on that PC can open it. When starting, the app passes the secrets on the
  command line (`?Password=`, `-Rcon -RconPassword=`, `-GSLTToken=`, `-GameStatsToken=`). The app offers to move
  an existing RCON password from `Game.ini` into the vault.
- **On screen:** secret fields only show "set / not set". Seeing the value needs a click on "Show".
- **In logs:** every line shown goes through a mask (`api_key=`, `Password=`, `Token=`, `-SecurityCode=`).
- **Known limit:** another program of the same Windows user can read a process command line. It is the only way
  the server accepts these values without keeping them on disk.

## ADR-003 · The app starts `InsurgencyServer.exe` directly, without `.bat`

- **Decision:** Go builds the arguments from the profile (pure, tested domain code) and starts the executable with
  `exec.Command` (argument list, no shell) in the server folder. Old `.bat` files still work, but the app does not
  need them.
- The server **keeps running if the app closes**: when it opens, the app looks for `InsurgencyServer*.exe`
  processes of the configured folder and follows them again. Stop ends the processes of that folder.
- **mod.io rules the app applies** (confirmed on a real server): `-SecurityCode=none` in normal use; a new code
  only once to authorize, then back to `none`; the start map is always a base-game map and the mod map comes
  through `-ModDownloadTravelTo`.

## ADR-004 · Server state read from the log

- **Decision:** the app follows `Insurgency.log` (reading what was appended) and turns known lines into events:
  mod.io ready, port open, mod activated, mod installing/uninstalling, map travel, match state change, fatal
  error. The parsing rules live in `domain/logparse`, tested with real log lines.
- **Consequence:** if a game update changes a log text, only the phase shown becomes less precise; start and stop
  keep working.

## ADR-005 · Error catalog with cause and solution

- **Decision:** `domain/problems` holds a catalog of codes. Each code has a title, what happened, the likely cause,
  how to fix it and help links (by key). App errors (`AppError`) and log lines point to a catalog item. The screen
  always shows the three blocks: "What happened", "Possible cause", "How to fix".
- An error that is not in the catalog becomes a generic item with the original (masked) message and a hint to open
  the logs.

## ADR-006 · Portable backup encrypted with a password

- **Context:** the DPAPI vault only opens on the same PC and Windows user. A backup has to open somewhere else.
- **Decision:** a `.ssbackup` file = header + data encrypted with **AES-256-GCM**, key derived from the backup
  password with **Argon2id** (64 MiB, 3 passes, random salt). Inside, a zip with `manifest.json`, the profile, the
  secrets, the server config files (`Game.ini`, `Engine.ini`, `Admins.txt`, `MapCycle.txt`, `Bans.json`,
  `Motd.txt`) and, if the user ticks it, the mod.io session (`user.json`).
- Restoring asks for the password, shows the contents, asks for the server folder on the new PC and writes
  everything. Secrets go back into the DPAPI vault of the new PC.
- Backup password: at least 10 characters. Without it the backup cannot be opened (there is no recovery).

## ADR-007 · Downloads with visible progress

- **SteamCMD:** downloaded from `https://steamcdn-a.akamaihd.net/client/installer/steamcmd.zip` (fixed address),
  with byte progress. Only `steamcmd.exe` is extracted, and its digital signature is checked (`WinVerifyTrust`)
  before it runs.
- **Server:** SteamCMD writes `Update state (0x61) downloading, progress: 12.34 (x / y)`. The domain reads these
  lines and the screen shows phase, percentage, bytes, speed and time left.
- **Mods:** the size of each mod comes from the mod.io API (public data); progress is the size of the mod folder on
  disk compared with the final size. Log events (installing, installed) mark start and end.

## ADR-008 · Apache 2.0 license and disclaimer (public project)

- **Context:** the project is public on GitHub so that others can use it, contribute and fork it, without the
  authors answering for misuse or illegal use.
- **Decision:** **Apache License 2.0** (`LICENSE`) plus `DISCLAIMER.md`.
  - It allows use, copying, changes, forks and commercial use. Whoever redistributes keeps the license and the
    notice, and marks the files they changed.
  - Sections 7 and 8 say the software comes "AS IS", without warranty, and that the authors are not liable.
  - Section 6 grants no trademark rights: it helps keep the project apart from Insurgency, Steam and mod.io.
  - It has an explicit patent license, which MIT does not.
- **Disclaimer:** unofficial project, not affiliated with New World Interactive, Focus Entertainment, Valve or
  mod.io. Each user answers for how they run their server, for the game, Steam and mod.io terms and for the laws
  of their country. The app does not distribute the game or SteamCMD: it downloads them from the official sources.
- **Consequence:** the dependencies used (Wails MIT, React MIT, x/crypto BSD) are compatible. This is not legal
  advice.

## ADR-009 · Interface in several languages (en-US default, pt-PT, es-ES, ru, uk)

- **Decision:** screen texts live in the front end, in `desktop/frontend/src/i18n/locales/<language>.ts`. English
  (`en-US.ts`) is the reference: the other languages share its type, and TypeScript flags a missing key. No
  external library (a small `t(key, params)`, fewer dependencies to audit).
- Go **has no user-facing text**: it returns codes (`problems.*`), validation keys (`validation.*`) with
  parameters, and events. The error catalog texts (ADR-005) live in the front end, in
  `problems.<code>.{title,what,cause,solution}`.
- Start language: the Windows language if it is one of the five; otherwise English. The user changes it in
  Settings.
- Documentation language: see ADR-035.

## ADR-010 · mod.io API key out of the code

- **Context:** hand-made `.bat` files often have the mod.io API key written in them. In a public repository it
  would be exposed for anyone to use on someone else's behalf.
- **Decision:** the key is a vault secret (`modioApiKey`). The app finds it in the server's own `Insurgency.log`
  (game requests carry `api_key=`), or the user creates one at https://mod.io/me/access, with a guide on screen.
- API address fixed in code: `https://g-254.modapi.io/v1` (game 254). HTTPS only.

## ADR-011 · Color themes and layout skins

- **Decision:** `<html data-theme>` changes only the color variables and `<html data-skin>` changes the layout
  shape (`shared/appearance.ts` + `styles/themes.css`). Colors: Sand (default, amber accent), Black, Olive, Dark
  blue and Light. Skins: Sidebar (default), Compact (icons only), Cards (floating panels) and Top bar. The choice
  is saved in `localStorage` (per PC) inside try/catch: with no storage, the default stays. It is applied before
  the first paint (no flash).
- **Consequence:** a new theme = one block of variables in `themes.css` + an item in `THEMES` + translations.

## ADR-012 · Private data out of the repository

- Notes about a real server (network, e-mail, tokens) stay out of git (`.gitignore`). Code, tests and docs carry no
  keys, tokens, IP addresses or e-mail addresses.
- The author's name and GitHub handle appear only in the credits (README, About screen, `NOTICE`), as chosen by the
  author (ADR-035).

## ADR-013 · Fill the fields from an existing server

- **Context:** someone who already set up a server by hand (running or stopped) does not want to type everything
  again. If they already changed a field in the app, the app must not overwrite it without asking.
- **Decision:** the main source is the `LogInit: Command Line:` line that the server writes at the start of each log.
  It works for a running server (current log) and a stopped one (last log), wherever the `.bat` files are. Logs are
  read from the newest to the oldest by date (`Insurgency.log`, `Insurgency_2.log`, `Insurgency-backup-*.log`).
  They are completed by `Admins.txt`, `MapCycle.txt`, the `[Rcon]` section of `Game.ini` and the mod.io key found in
  the logs (ADR-010).
- **Compare rule** (domain, `server.Compare`/`server.Merge`): equal → nothing; empty or default value in the app →
  filled without asking; different from the default and from the server → **conflict**, and the screen asks field
  by field "Mine (app)" or "Server" (with "keep all mine" and "use all from the server"). "Apply" is enabled only
  when every conflict has a choice. Unsaved changes on the Server screen are saved first, so they count as "mine".
- **Security:** secrets never go to the screen (only "kept in the app" / "different on the server"); the screen
  sends only the keys to keep (`keepMine`, checked against the list); Go reads the files again when applying; an
  invalid server value (for example repeated ports) falls back to the app value and shows as "not used".
- **Parent folder:** if the chosen folder has no `InsurgencyServer.exe` but exactly one subfolder has, the app uses
  the subfolder. This avoids installing a second server by mistake.

## ADR-014 · Per-user installer and automatic updates with a signed manifest

- **Decision:**
  - **Per-user NSIS installer** (`npm run build:installer` → `wails build -nsis -installscope user`): installs into
    `%LOCALAPPDATA%\Programs\Sandstorm Server Assist` without administrator rights, in 5 languages, with LICENSE,
    NOTICE, DISCLAIMER, TERMS and PRIVACY. Output:
    `desktop/build/dist/sandstorm-server-assist-<version>-windows-setup.exe`. With `/S /relaunch` it waits for the
    app to close, installs silently and reopens the app.
  - **`latest.json` manifest:** version, installer name (with the same version), SHA-256, size and an **Ed25519**
    signature. The private key lives only on the publisher's PC
    (`~/.sandstorm-server-assist/update-signing-key.pem`, `npm run update:keygen`); the public key is in the app
    (`adapters/updater/pubkey.go`). `npm run update:manifest` signs and refuses a key that is not the app's.
  - **App** (`domain/update` pure rules, `app/updater` use case, `adapters/updater` HTTP/disk/launcher): checks 10 s
    after opening and every 4 h, only through `netx` (HTTPS and a fixed host list). It checks the signature before
    downloading and the SHA-256/size after. When ready, the title bar shows **Update** (installs and reopens); if the
    user just closes the app, it installs on close without reopening. It does not install while a server install is
    running. The portable version (no `uninstall.exe` next to it) only tells the user and opens the release page.
  - **Game server:** a separate process; it keeps running during the update.
- **Where the manifest lives:** first on a download page on the author's server (ADR-016), now on GitHub Releases
  (ADR-032).
- **Consequence:** losing the private key means installed apps cannot be updated (a version with another key would
  have to be installed by hand).

## ADR-015 · Terms of use and privacy notice (accepted on first use)

- **Decision:** both, in a single acceptance (clickwrap):
  - **Terms of use:** what the app is (unofficial), no warranty (Apache 2.0), the user's responsibility (game, Steam
    and mod.io terms, laws), what is not allowed, third-party services, automatic updates.
  - **Privacy notice:** even without collecting anything, the app keeps secrets and sends data to third parties
    (e-mail and key to mod.io, update checks, the port diagnosis). The notice says what stays on the PC, what
    leaves and to whom, the role of whoever runs the server with players' data and how to delete everything. No
    account, no ads, no telemetry.
- **How:** texts in the 5 languages (`legal.*`) and in `TERMS.md`/`PRIVACY.md`. The version is
  `domain/legal.Version`; the acceptance (version and time) is kept in the preferences (`config.json`). Until the
  current version is accepted, the acceptance screen covers the app (the title bar stays, so it can be closed). Update
  checks start only after acceptance. About shows the accepted version and opens the texts again.
- **Consequence:** when a text changes, `legal.Version` goes up and everyone accepts again.

## ADR-016 · Download page on the author's server (replaced by ADR-032)

- **Decision (2026-11-09):** a static page in `site/`, served by the author's web server, with the downloads and the
  `latest.json` of ADR-014. No inline script or style (strict CSP); English in the HTML and `site/lang.js` switches to
  the other 4 app languages without `innerHTML`. `terms.html` and `privacy.html` are generated from `TERMS.md` and
  `PRIVACY.md`. `scripts/publish-site.mjs` builds the page and sends it over SSH; the SSH address comes from
  `SSA_SSH_HOST`, never from the code.
- **Now:** GitHub Releases is the download place and the update source (ADR-032). The page and its `latest.json` stay
  online for a while only as a bridge, so apps up to 0.1.10 can update (see `docs/RELEASING.md`).

## ADR-017 · Restart once when new mods leave the installed ones out

- **Context:** on a real first run, the server came up on a base-game map, without the mod map and without a mod
  mutator. The log showed the server rule: when there are new subscriptions (`PendingChanges with state Added`),
  it downloads and mounts only the new mods and travels to `-ModDownloadTravelTo` without mounting the ones that
  were already installed. The travel falls back to the start map and the mutator becomes "invalid". With nothing
  new, the next run mounts everything.
- **Decision:** the runner keeps the "Added" mods of the run. On travel, if there was a new mod and an installed mod
  (`state.json`, `State 1`) was not activated (`OnModActivatedFromCloud`), the app restarts the server **once**
  (only a server it started; notice on the Dashboard). A server started outside the app only gets the problem
  `mods-not-mounted` ("click Restart"). It never restarts twice in a row.
- **Also:** the map cycle fixes `Scenario="Scenario=...` (pasted from a mod page; the game rejects the entry) when
  saving, opening and writing MapCycle.txt; mods uninstalled by the server leave the Mods list; errors the game
  always writes (cosmetics, replay, voice, GameStats frame, assets repeated across maps) count as warnings.

## ADR-018 · Logs per server run, optional cleanup and a single app instance

- **Context:** the log drawer mixed every run seen since the app opened (hundreds of "errors" in an afternoon) and
  nothing deleted old logs.
- **Decision:**
  - **Runs:** the server already writes one file per run (`Insurgency.log` and `Insurgency-backup-<date>.log`).
    `logbuf` numbers the runs (`StartRun` on each start, also for a server started outside the app) and the drawer
    shows only the current one, with a "Server started at …" line and its error count. A selector opens the runs
    kept by the server (`LogAPI.Runs`/`ReadRun`: the screen sends only the name, checked against the list; last
    6000 lines, masked).
  - **Cleanup** (`domain/serverlogs`): only names the server creates; never the current log; keeps N runs (minimum
    2, the current and the previous one, which the app reads for "Fetch from server" and the mod.io key; maximum
    100; suggested 10). Settings › Server logs: automatic cleanup **off by default**; when on, it runs once per run
    started by the app. A "Delete old logs now" button asks for confirmation (how many and how much space).
  - **Noise:** `Could not find outer … to create` (mod maps missing internal objects; hundreds per run) counts as a
    warning.
  - **Single instance** (`SingleInstanceLock`): two open copies saved `config.json` over each other and would follow
    the same server twice. Opening again brings the window to the front. The development build uses another id.

## ADR-019 · Change the map with the server running, through RCON

- **Decision:**
  - **Own RCON client** (Source RCON protocol): `domain/rcon` builds and reads packets (size limited to 1 MB);
    `adapters/rcon` connects only to `127.0.0.1` (the app runs on the server PC), with the password from the vault.
  - **One connection per click**, never on a timer: the Sandstorm RCON leaves a thread stuck on the server for each
    connection. One action at a time (`busy`).
  - **List:** the `scenarios` command. On a real server it lists only base-game maps. Mod maps come from the mods
    loaded in the run (ADR-029).
  - **Change:** `travelscenario <scenario>`, which keeps password, slots and mutators. The screen sends only the id;
    Go accepts only an id from the last list and rejects server error replies.
  - New problems: `rcon-off`, `rcon-failed`, `rcon-auth`, `map-change-failed`.
  - The change does not touch the profile: after a restart, the server goes back to the map of Server › Maps.

## ADR-020 · Message of the day and join/leave messages

- **Decision:**
  - **Message of the day:** the server reads `Insurgency/Config/Server/Motd.txt` when it starts and shows it on the
    loading screen. `config.Prepare` writes the file when the text is not empty (up to 1000 characters, several
    lines); "Fetch from server" reads an existing `Motd.txt`.
  - **Join and leave:** the game has no text of its own for this. The app reads the log (`logparse`):
    `Login request: …?Name=<name> userId: <id>`, then `Join succeeded: <name>` (sometimes cut; matched by the start
    of the name) and `UNetConnection::Close … GameNetDriver … PC: INSPlayerController … UniqueId: <id>` (leave, id
    only). It sends `say <text>` through RCON, with `{name}` replaced (cleaned of control characters, up to 40
    characters; the message up to 150).
  - **One open RCON connection** for messages (`app/chat`), queued, without blocking log reading; if it drops, it
    reopens once. It closes when the server stops. Lines older than 2 minutes do not create messages (the app may
    read the log of a run already in progress).
  - Needs RCON on (validation `chatNeedsRcon`); works only while the app is open (said on screen).
- **Consequence:** if a game update changes these log lines, the messages stop and nothing else breaks.

## ADR-021 · Neutral network texts and a port test from the internet (test replaced by ADR-033)

- **Decision (2026-11-10):** the network screen has two neutral guides (ports in the router, or a VPN program), the
  list of UDP ports to forward and a warning never to forward the RCON port. A small checker service (`checker/`,
  Go standard library only) ran on the author's server and tested the caller's IP from outside: an A2S query on the
  query port and a TCP connection to the RCON port, with rate limits and no logs.
- **Now:** the app diagnoses from the PC itself (ADR-033). `checker/` stays in the repository only while apps up to
  0.1.10 may still call it; it is removed after the bridge period.

## ADR-022 · More server settings (rules, players, automation, Windows, profiles)

- **Source of the options:** the official Server Admin Guide by New World Interactive and the RCON `help` of a real
  server.
- **Decision:**
  - **Game rules** (`domain/server/rules.go`): a catalog of 60 options of Game.ini (coop, checkpoint, objectives,
    friendly fire, team kills, chat and voice, match, supply, vote kick) and Engine.ini (tick rate and bandwidth).
    `Profile.Rules` nil = the app does not touch the files (a hand-made server stays the same); when on, the app
    writes the options with a value and removes the empty ones when starting (`config.Prepare`, `.bak` copy).
  - **Mutators:** the official list (25) with translated descriptions, plus mod mutators (ADR-023).
  - **Day or night:** `?Lighting=` in the start URL, `Lighting="Day|Night"` in map cycle lines and in the map change
    (with lighting it sends `travel <map>?Scenario=…?Lighting=…?MaxPlayers=…[?Password=…]`, because
    `travelscenario` does not accept lighting).
  - **Players** (`app/players`, `domain/admin`): `listplayers` and `listbans` on one connection per click; kick,
    ban (`banid` with minutes or `permban`), unban and a message to everyone (`say`). Only net IDs from the last
    list; reason without quotes or control characters, up to 100 characters.
  - **Scheduled messages:** up to 10, one at a time every N minutes (2–240), on the chat connection (ADR-020).
  - **Automatic restart** (`app/watchdog`): daily at a time, with chat warnings 5 and 1 minute before and a SteamCMD
    update; and after a crash (waits 10 s, stops after 3 times in 30 minutes). Only while the app is open.
  - **Discord** (`app/notify`, `adapters/discord`): webhook in the vault (only official Discord addresses); user
    texts for online, stopped, crashed, join and leave (`{server}`, `{map}`, `{name}`); names escaped and empty
    `allowed_mentions` (nobody is pinged). Queue of 32, one send per second.
  - **Windows** (`adapters/winsys`): start with Windows (`HKCU\…\Run`, minimized with `--autostart`; option to start
    the server too); Firewall rule: the app opens a copy of itself as administrator (UAC) with `--firewall-rule`,
    which reads the server folder from the config and creates inbound rules for `InsurgencyServer.exe` and
    `InsurgencyServer-Win64-Shipping.exe` (netsh, fixed names); network adapters list for `-MultiHome`.
  - **Profiles** (`app/config/profiles.go`): switch only with the server stopped; vault secrets are shared; one
    server at a time.
  - **Mod list** (`-CmdModList`): optional and off; with the current mod.io the server account subscriptions rule.

## ADR-023 · Mutators of installed mods

- **Decision:** a mod mutator is a file `<Plugin>/Content/Mutators/<Name>.uasset` inside the mod `.pak`, and `<Name>`
  is what goes into `-mutators=`. `domain/pak` reads the `.pak` footer (versions 10 and 11), the main index and the
  directory index, checking every size and offset (third-party file; fuzz tested). `adapters/modio` reads only those
  parts of each `.pak` of the installed mods. If no `.pak` can be read and the mod has the "Mutator" tag on mod.io,
  it uses the `PluginName` of the `metadata_blob`.
- **Screen:** one scrolling list: "From your mods" (MOD badge and mod name) and "From the game". A profile mutator
  that no installed mod provides stays ticked, with a warning.

## ADR-024 · Message formatting and mod map picker (formatting replaced by ADR-031)

- **Mod map picker (still valid):** `pak.Scan` also lists maps (`.umap`) and scenarios (`Scenario_*.uasset`); the map
  of a scenario is the shortest `.umap` name such that the scenario is `Scenario_<map>_<mode>`.
  `ConfigAPI.ModScenarios` feeds the map and scenario pickers. A profile value missing from the installed mods
  stays in the list, flagged.
- **Formatting (replaced):** the game files have rich-text style tables for the loading screen and chat, so the app
  offered `<Bold>`, headers and images with an editor. The game showed the tags as plain text (ADR-031).

## ADR-025 · Lists never null and protection against screen errors

- **Context:** in 0.1.5, "Show players" with an empty server turned the whole app black: Go returned the lists as
  `nil` (JSON `null`), the screen read `.length` and React unmounted everything.
- **Decision:**
  - Lists sent to the screen start empty in Go (`[]T{}`), never `nil`. The screen also treats `null` as empty.
  - `shared/ErrorBoundary` in two layers: per page (a warning on that page only, "Try again" and "Back to the
    Dashboard"; the menu keeps working) and for the whole app (full-screen warning with "Reload the app"). The game
    server is not affected.
  - The error goes to the app logs (`LogAPI.ReportScreenError`), one line, cut at 2000 characters, with the secret
    mask.

## ADR-026 · A reopened app recognises the server it started

- **Decision:**
  - When starting, the app writes `server-run.json` next to `config.json`: time, folder and the command line already
    masked (no secrets). It is deleted when the server stops.
  - When the app opens with the server running, the server is the app's own if the folder is the same and the oldest
    process started between 5 s before and 2 min after the recorded time. Otherwise it counts as started outside the
    app (no automatic restart to load new mods; everything else works the same).
- **Limit:** mod.io authorization with a code (ADR-003) is not resumed.

## ADR-027 · Waiting for the mod map ends at the first travel

- **Decision:** at start the server opens a base-game map and travels to the mod map (`-ModDownloadTravelTo`). The
  app waits for that map only until the server's first travel; after it, any map counts as ready (otherwise a map
  change left the Dashboard stuck on "loading the map"). A server started outside the app does not wait at all.

## ADR-028 · Travel to a mod map read from the map path

- **Decision:** for a mod map the server logs the path (`ProcessServerTravel: /Estates/Maps/Estates?Scenario=...`);
  for an official map only the name. The parser takes the last part of the path and finds the scenario in any
  option position. Tested with real lines.

## ADR-029 · Map change with the maps of the loaded mods, in map and mode lists

- **Decision:** mod maps for the map change come from the mods active in the run (read from the log), crossed with
  the scenarios of those mods' `.pak` files (ADR-024): `MapChoices(server, profile, mods)`, source `mod`. With no
  known active mod, it uses the profile map cycle as before. The screen has two lists: the map (flagged "mod map" and
  "now") and its mode.

## ADR-030 · Map cycle entries without their mod: warn, do not delete

- **Decision:** what exists is read from the files: base-game scenarios from the server `.pak` files (indexes only,
  about 20 ms) and from all installed mods, cached for 15 s. On Server › Maps each cycle line without a known map
  gets "mod not installed" and "Remove", and a notice offers "Remove these maps". The app never removes them by
  itself. In the map change, such maps show disabled and Go refuses them (`validation.scenarioNotLoaded`).

## ADR-031 · Plain-text messages (replaces the formatting of ADR-024)

- **Context:** the game showed the formatting tags as text on the loading screen (`Motd.txt`) and in chat (`say`).
- **Decision:** the message of the day and chat messages are plain text. Formatting tags are refused when saving
  (`validation.markupPlain`), and tags saved before are removed when writing `Motd.txt` and sending to chat
  (`StripMarkup`).
- **Lesson:** a game feature found only in the game files is marked "to be confirmed" and is not shipped before it is
  seen in the game.

## ADR-032 · Updates and downloads from GitHub Releases

- **Context:** the project is free and public. Depending on the author's own server for downloads and updates ties
  the app to that server; GitHub Releases is where open-source users look for downloads, and it shows the source
  code next to each version.
- **Decision:**
  - Each version is a GitHub Release `v<version>` with the installer, the portable zip, the signed `latest.json`
    (ADR-014) and `SHA256SUMS.txt`. The notes come from `CHANGELOG.md`. `npm run release:publish` builds the
    manifest, checks that the commit is on `origin/main` and publishes with the GitHub CLI.
  - The app reads `https://github.com/<repo>/releases/latest/download/latest.json` and downloads the installer from
    the release of that exact version (`releases/download/v<version>/<file>`), so a release published in the middle
    of a download does not mix versions. The signature and SHA-256 are checked as before.
  - `netx` allows `github.com` and its file hosts (`release-assets.githubusercontent.com`,
    `objects.githubusercontent.com`); the old download server left the list.
- **Bridge:** apps up to 0.1.10 only know the old download server. It keeps serving `latest.json` (pointing to the
  new version) for a while, so they update once and then use GitHub (`docs/RELEASING.md`).

## ADR-033 · Port diagnosis from the PC (replaces the port test of ADR-021)

- **Context:** the internet port test needed a service on the author's server. Without it, a PC cannot prove that
  its ports are reachable from the internet, but it can check the usual causes.
- **Decision:** "Diagnose ports" (`adapters/netdiag`, rules in `domain/netcheck`) measures, in parallel and in up to
  15 s:
  - the server answers the A2S query on `127.0.0.1` (server running and query port right);
  - the app's Windows Firewall rules exist;
  - the router, through UPnP (read only): its internet IP and the forwards of the game, query and RCON ports. Only
    `http://` addresses of the private network, from the device that answered, with no proxy and no redirects;
  - the public IPv4 through STUN (Cloudflare, Google), which also shows CGNAT when the router IP is private, in
    100.64.0.0/10 or different from the public one;
  - the A2S query through the public IP (works only on routers with NAT loopback, so a missing answer is not an
    error);
  - whether Steam lists a server at the public IP (`ISteamApps/GetServersAtAddress`, public API without a key).
    The result is a list of coded items (ok, warning, problem, info) and a verdict. The screen always explains how to
    confirm with a friend outside the network (Steam › View › Game Servers › Favorites › Add a server).
- **Privacy:** STUN servers and the Steam API see the public IP; the privacy notice says so (`legal.Version`
  2026-11-11). The old `port-check-*` problems were removed.

## ADR-034 · Feedback through GitHub issues, prepared by the app

- **Context:** users should be able to report problems and suggest improvements with a description, the log and an
  image. Sending directly from the app would need a server or a key inside the app that anyone could extract and
  abuse.
- **Decision:** the Feedback page builds a `.zip` on the PC (`%LOCALAPPDATA%\sandstorm-server-assist\reports`):
  `report.txt` (type, title, version, Windows, language, description), `app-log.txt` (optional; last 3000 lines up
  to 2 MB, already secret-masked, plus IP addresses and SteamIDs masked) and an optional screenshot (chosen in a
  Windows dialog opened by Go, up to 5 MB, PNG/JPG/GIF/WebP checked by content). Then it opens the GitHub issue form
  (`.github/ISSUE_TEMPLATE/bug_report.yml` or `feature_request.yml`) filled through the URL and shows the zip in
  Explorer. Title, description and image are optional. The app sends nothing; the user posts with their own GitHub
  account.
- **Problems:** `report-failed`, `report-image`.

## ADR-035 · Documentation in English, user docs in 5 languages

- **Decision:** every document in the repository is in English. The README and the user guide also exist in the 4
  other app languages (`README.<lang>.md`, `docs/<lang>/GUIDE.md`), each with a language bar at the top and a note
  that English wins when they differ. Developer documents (contributing, architecture, decisions, releasing, roadmap,
  dev log) are English only. `npm test` checks that every translation has the same sections and images as the
  English file. Code comments are still mostly in Portuguese; new code may use English.
- **Credits:** the author is credited by name and GitHub handle in the README, the About screen and `NOTICE`; the app
  was built with the help of Claude Code (Anthropic).

## ADR-036 · Installation feedback next to the button

- **Context:** on the Installation screen, "Check for updates" seemed to do nothing, then "Verify files" showed
  "Wait for the current task" at the top of the page, out of view. Two causes: Go marked the task as running but
  only told the screen at the first progress line, and SteamCMD spends a while updating itself before printing any
  (the buttons stayed enabled and nothing showed); and the progress, result and error lived at the top of the page.
  "Download SteamCMD" with SteamCMD already present only checked the signature and showed nothing near the button.
- **Decision:**
  - The installer tells the screen as soon as a task starts. The state also says which task ran last (`last`),
    whether it verified files (`validate`), whether the server was already up to date (`upToDate`, from SteamCMD's
    "already up to date") and whether SteamCMD was already there (`present`).
  - Progress, result and errors show inside the card of the button that started the task, and the area scrolls into
    view when the task starts. The buttons are disabled from the click on.
  - With SteamCMD present, the button is **Check SteamCMD**, with a note that SteamCMD updates itself every time the
    app runs it; the check downloads it again only if the file is damaged.

## ADR-037 · Delete a downloaded mod so the server downloads it again

- **Context:** the maintainer asked for a way to delete a downloaded mod, to force the server to download it again
  (a broken mod, or a fresh copy), either on the next start or right away. The server's mod.io plugin downloads the
  subscribed mods when it starts and keeps what is installed in `254/metadata/state.json`. Deleting only the mod
  folder is not enough: the entry would still say "installed" and the server would try to mount files that are gone.
  The server rewrites `state.json` when it closes, so it must be stopped.
- **Decision:**
  - `modio.RemoveInstalled` removes the mod from the `Mods` list of `state.json` (other mods and fields kept as they
    were, a `.bak` copy on the first change) and deletes `254/mods/<id>`. The folder comes only from the id, never
    from `PathOnDisk`, and must be inside the mod.io mods folder.
  - `app/modfiles` orders the steps: **Delete and download again** stops the server if it runs, deletes, then starts
    it (it downloads the mod; the restart for new mods of ADR-017 still applies). **Only delete** needs the server
    stopped; the mod comes back on the next start if it is still subscribed. If deleting fails after a stop, the
    server is started again. It refuses while the server is being installed or updated.
  - On the Mods screen, each downloaded mod has a delete button; the confirmation, progress, result and error show on
    the mod's row. New problem code `mod-delete`.

## ADR-038 · Download new mods now: start or restart the server

- **Context:** after subscribing to new mods and reading the subscriptions, the Mods list showed them as "will be
  downloaded" with no way to download them. Only the server's mod.io plugin downloads mods, and only when the server
  starts; a running server does not load new mods.
- **Decision:** when subscribed mods are not on disk, the Mods screen lists them above the table with **Download
  now**, which starts the server, or **Restart and download now** with a confirmation (the match ends and players
  leave for a moment) when it runs. `modfiles.DownloadNew` refuses while the server is being installed. After the
  download, the restart of ADR-017 loads every mod. The progress shows on each mod's row.
- **Not chosen:** downloading the mods with the app. It would copy the plugin's private files (`state.json`, the mod
  folders), which can change between game versions and which the server rewrites.

## ADR-039 · Find mods with problems from their files

- **Context:** the maintainer's game crashed on maps that used to work, with "Serial size mismatch" on a `/Game/...`
  object. This error comes from mod content: two mods ship the same game file in different versions, or a mod ships
  a copy made for another game version. The app only said "Mods in conflict", without naming a mod, and only after
  the server crashed.
- **Decision:**
  - `domain/pak` lists every file of a `.pak` with its uncompressed size: formats 10 and 11 from the encoded entries,
    formats 8 and 9 from the full index (with the UE 4.22 variant of the compression field). Only the footer and the
    indexes are read; every size and offset is checked (third-party files, fuzz tested).
  - `domain/modcheck.Conflicts`: the same object (path without extension) in two mods with any file of a different
    size is a conflict. Identical copies, common between map mods, do not count. The new pak entries have no hash, so
    size is the signal; a different file with the same size is not detected.
  - `modcheck.Check` reads the objects named in a crash message (`/Game/...`, `/Engine/...`, `/<Plugin>/...`) and finds
    the installed mods that ship them. The Mods screen has **Find the mod** for a pasted message; the text is not
    saved, logged or sent.
  - When the server log has "Serial size mismatch", `logparse` keeps the object and the runner adds it, with the mods
    that ship it, to the `mod-conflict` problem; the error card shows them.
- **Limit:** a mod whose copy does not match the current game version, with no other mod having the file, is found
  only through the crash message. The app never changes mod files; the fix is to unsubscribe the mod, or delete and
  download it again (ADR-037).
