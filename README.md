# Sandstorm Server Assist

A Windows desktop app that helps you **install, configure and run an Insurgency: Sandstorm dedicated server
with mod.io mods** — without editing `.bat` and `.ini` files by hand.

> Unofficial community tool. Not affiliated with New World Interactive, Focus Entertainment, Valve or mod.io.
> See [DISCLAIMER.md](DISCLAIMER.md).

## Features

- **First steps checklist** on the dashboard: what is ready and what is missing.
- **Installation with progress**: downloads SteamCMD (Valve's digital signature is checked before running it)
  and installs or updates the server, with progress bar, speed and time left.
- **Server settings**: name, ports, players, start map, mod map (`-ModDownloadTravelTo`), mutators, map
  rotation, administrators (SteamID64), join password, RCON, GSLT and GameStats tokens. Everything you must
  get somewhere else has a **"How do I get this?"** guide with links.
- **Use a server you already have**: pick its folder and the app reads the last launch command line from the
  server log plus `Game.ini`, `Admins.txt` and `MapCycle.txt` (server running or stopped). Empty or default
  fields are filled in; when you already typed something different, the app asks whether to keep your value or
  use the server's.
- **mod.io**: requests the authorization code by e-mail, authorizes the server once and switches back to
  normal mode by itself, lists subscribed/installed mods and shows each mod's download.
- **Start / stop** with live phases read from the server log: starting, mod.io connected, mods mounted,
  loading map, ready. Shows the console command players use to join (over the internet with the router ports open, or over a VPN program).
- **Change the map while the server runs**, including mod maps from your map rotation (RCON).
- **Server messages**: message of the day on the loading screen, and chat messages when a player joins or
  leaves (sent by the app through RCON; `{name}` becomes the player name).
- **Port test from the internet**: one click checks, from outside, whether the query port answers and whether the
  RCON port is wrongly exposed, and shows your public IP. Guides for router port forwarding or a VPN program.
- **Game rules**: bots and difficulty, checkpoint timers, friendly fire, team kills, kill feed, chat and voice,
  map vote, vote kick, supply points, tick rate and more, written to Game.ini and Engine.ini (60 options). Mutators
  with descriptions, day or night per map.
- **Players**: list, kick, ban (minutes or permanent), remove bans and message everyone (RCON).
- **Automation**: scheduled chat messages, daily restart with chat warning and SteamCMD update, restart when the
  server closes by itself, Discord notices (webhook), open with Windows and start the server.
- **Windows and network**: firewall rule for the server (with administrator confirmation), network card choice
  (-MultiHome). Several profiles (for example Coop and Versus).
- **Errors explained**: every known error shows what happened, the likely cause and how to fix it.
- **Logs** hidden by default, in a drawer that opens with one click, with filters and search. Passwords and
  tokens are masked.
- **Security**: secrets are encrypted with Windows DPAPI, never written to `.ini`/`.bat` files.
- **Portable backup**: one password-protected file (Argon2id + AES-256-GCM) with everything needed to set up
  the server on another PC.
- **5 languages**: English (US), Português (Portugal), Español (España), Русский, Українська.
- **Appearance**: 5 color themes and 4 layouts.
- **Installer and automatic updates**: a per-user Windows installer (no administrator). The installed app
  checks for new versions, downloads them in the background and installs only updates signed with the project
  key (Ed25519 + SHA-256). The game server keeps running while the app updates.
- **Terms of use and privacy notice**: accepted on first use ([TERMS.md](TERMS.md), [PRIVACY.md](PRIVACY.md)). No
  account, no ads, no telemetry.

## Development

Requirements: Windows 10/11 (WebView2), Go 1.27, Node 24.

```bash
npm install
npm run dev        # app in development mode
npm test           # Go + front-end tests
npm run lint
npm run typecheck
npm run build      # desktop/build/bin/sandstorm-server-assist.exe
npm run build:installer   # + desktop/build/dist/sandstorm-server-assist-<version>-windows-setup.exe (needs NSIS)
npm run package           # portable zip in desktop/build/dist
```

Releases: `npm run update:keygen` creates the update signing key once (outside the repository, in
`~/.sandstorm-server-assist/`); its public key goes in `desktop/internal/adapters/updater/pubkey.go`.
`npm run update:manifest` writes the signed `latest.json` for the installer in `desktop/build/dist`. Upload the
installer first and `latest.json` last to the download folder.

Download page: `site/` (5 languages). `npm run site:preview` builds it in `desktop/build/dist/site`;
`npm run site:publish` uploads it with the files and the signed `latest.json` (needs `SSA_SSH_HOST`). See
[docs/PAGINA_DOWNLOAD.md](docs/PAGINA_DOWNLOAD.md).

Security checks: `npm run audit` (npm deps), `npm run audit:go` (govulncheck), `npm run audit:code` (gosec +
ESLint security rules).

Architecture: Wails 2 (Go, Clean/Hexagonal) + React 19 + Vite. Business rules live in Go
(`desktop/internal/domain` and `desktop/internal/app`); the UI only sends ids and options. Details and
decisions (in Portuguese) are in [docs/](docs/).

## Translations

UI texts live in `desktop/frontend/src/i18n/locales/`. `en-US.ts` is the reference; the other files must have
exactly the same keys (a test checks it). To add a language, copy `en-US.ts`, translate it and register it in
`desktop/frontend/src/i18n/index.tsx`.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Security issues: [SECURITY.md](SECURITY.md).

## License

[Apache License 2.0](LICENSE). See also [NOTICE](NOTICE), [DISCLAIMER.md](DISCLAIMER.md), [TERMS.md](TERMS.md) and [PRIVACY.md](PRIVACY.md).
