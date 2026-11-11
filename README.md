# Sandstorm Server Assist

**English** · [Português](README.pt-PT.md) · [Español](README.es-ES.md) · [Русский](README.ru.md) · [Українська](README.uk.md)

[![Latest release](https://img.shields.io/github/v/release/felipesulzbach/sandstorm-server-assist)](https://github.com/felipesulzbach/sandstorm-server-assist/releases/latest)
[![License: Apache 2.0](https://img.shields.io/badge/license-Apache%202.0-blue)](LICENSE)

A free Windows app that helps you **install, configure and run an Insurgency: Sandstorm dedicated server with
mod.io mods**, without editing `.bat` and `.ini` files by hand.

> Unofficial community tool. Not affiliated with New World Interactive, Focus Entertainment, Valve or mod.io.
> See [DISCLAIMER.md](DISCLAIMER.md).

![Dashboard with the server running](docs/images/dashboard.png)

## Download

Get the latest version from the **[Releases page](https://github.com/felipesulzbach/sandstorm-server-assist/releases/latest)**:

- **Installer** (`…-windows-setup.exe`, recommended): installs for your Windows user, no administrator needed, and
  updates itself.
- **Portable** (`…-windows-portable.zip`): unzip and run; it only tells you when a new version exists.

Windows 10 or 11 is required. Windows may show a SmartScreen warning because the app is not code-signed: click
**More info** › **Run anyway**. Each release lists the SHA-256 of its files.

## Features

- **Step-by-step setup:** a checklist on the Dashboard shows what is ready and what is missing.
- **Install with progress:** downloads SteamCMD (Valve's signature is checked) and installs or updates the server,
  with speed and time left.
- **Use the server you already have:** pick its folder and the app reads its settings, asking before it changes
  anything you typed.
- **All the settings in one place:** name, ports, players, maps and map cycle, mod maps, day or night, mutators
  (including those of your mods), 60 game rules, admins, passwords and tokens, with "How to get it?" guides.
- **mod.io mods:** authorize the server once, see each mod download, and the app restarts once when new mods need
  it.
- **Run it:** start and stop with live phases, change the map while it runs, players panel (kick, ban, message),
  message of the day and join/leave messages.
- **Automation:** scheduled messages, daily restart with warning and update, restart after a crash, Discord
  notices, start with Windows.
- **Port diagnosis:** checks the Firewall, your router (UPnP), your public IP, CGNAT and Steam, and explains how to
  confirm with a friend.
- **Explained errors and clean logs:** every known error says what happened, why and how to fix it; passwords and
  tokens are always masked.
- **Safe secrets and portable backup:** secrets are encrypted by Windows (DPAPI); one password-protected file moves
  everything to another PC.
- **Feedback:** report a problem or suggest an idea from the app, with the log and an optional screenshot.
- **5 languages and 5 themes:** English, Português, Español, Русский, Українська.

## Quick start

1. Install the app and accept the terms of use.
2. **Installation:** choose a server folder (about 25 GB free), download SteamCMD and install the server. If you
   already have a server, choose its folder instead.
3. **Server:** set the name, the maps and the passwords. Get a Steam game server token (GSLT) with the guide.
4. **Mods:** authorize the server on mod.io and subscribe to mods on the mod.io website.
5. **Server › Network:** forward the game and query ports (UDP) in your router, create the Firewall rule and run
   **Diagnose ports**.
6. **Dashboard:** click **Start server** and share the join command with your players.

The full tutorial, with pictures, is in the **[user guide](docs/GUIDE.md)**.

## Screenshots

| Installation                                  | Server settings                                     |
| --------------------------------------------- | --------------------------------------------------- |
| ![Installation](docs/images/install.png)      | ![Server settings](docs/images/server-settings.png) |
| **Port diagnosis**                            | **Feedback**                                        |
| ![Port diagnosis](docs/images/port-check.png) | ![Feedback](docs/images/feedback.png)               |

## Feedback and support

- **Found a problem or have an idea?** In the app, open **Feedback**: it prepares a file with the details and opens
  the GitHub form already filled in. Or [open an issue](https://github.com/felipesulzbach/sandstorm-server-assist/issues/new/choose).
- **Security problem?** Please do not open a public issue: see [SECURITY.md](SECURITY.md).

## Build from source

Requirements: Windows 10/11 (WebView2), Go 1.27, Node.js 24; NSIS for the installer.

```bash
npm install
npm run dev               # app in development mode
npm test                  # Go, front-end and docs tests
npm run build             # desktop/build/bin/sandstorm-server-assist.exe
npm run build:installer   # installer in desktop/build/dist
```

Wails 2 (Go, Clean/Hexagonal) + React 19 + Vite. Rules live in Go (`desktop/internal`); the screens only send ids and
options. See [architecture](docs/ARCHITECTURE.md), [decisions](docs/DECISIONS.md), [releasing](docs/RELEASING.md) and
[contributing](CONTRIBUTING.md).

## Translations

The app and these docs are available in 5 languages. Screen texts are in `desktop/frontend/src/i18n/locales/`
(`en-US.ts` is the reference). Translations of the README and the guide follow the English files; if they differ,
the English version is correct. Help with translations is welcome.

## Credits

Created by **Felipe Sulzbach** ([@felipesulzbach](https://github.com/felipesulzbach)), with the help of
[Claude Code](https://claude.com/claude-code) by Anthropic. Made for the Insurgency: Sandstorm community.

## License

[Apache License 2.0](LICENSE). See also [NOTICE](NOTICE), [DISCLAIMER.md](DISCLAIMER.md), [TERMS.md](TERMS.md) and
[PRIVACY.md](PRIVACY.md).
