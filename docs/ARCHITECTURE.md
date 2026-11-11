# Architecture

> What the app is, the stack, how the code is organised and the project rules.
> See also: [Decisions](DECISIONS.md) · [Roadmap](dev/ROADMAP.md) · [Releasing](RELEASING.md) ·
> [User guide](GUIDE.md)

## 1. What it is

A Windows desktop app that helps you **install, configure and run** an Insurgency: Sandstorm dedicated server with
mod.io mods, without editing `.bat` and `.ini` files by hand.

- **First steps:** a checklist shows what is ready and what is missing.
- **Install:** downloads SteamCMD and installs or updates the server (app 581330), with progress, speed and time
  left.
- **Settings:** server name, ports, players, start map, mod map (`-ModDownloadTravelTo`), admins (SteamID64), map
  cycle, game rules, mutators, join password, RCON, GSLT and GameStats. Anything the user must get elsewhere has a
  "How to get it?" guide with steps and links.
- **mod.io:** requests the authorization code by e-mail, authorizes the server, lists the subscribed mods and shows
  each download.
- **Start and stop:** the app builds the command line and follows the phases from the log: starting, mod.io ready,
  mods activated, traveling to the map, ready to play.
- **While running:** change the map, players panel, messages, scheduled restarts, Discord notices.
- **Explained errors:** each known error shows what happened, the likely cause and how to fix it.
- **Logs:** hidden by default in a drawer, with filters and search. Secrets are masked.
- **Security:** passwords and tokens live in a vault encrypted by Windows (DPAPI), not in `.ini` or `.bat` files.
- **Backup:** one password-protected file with everything needed to run the server on another PC.
- **Port diagnosis** from the PC, **Feedback** through GitHub issues, **automatic updates** from GitHub Releases.
- **Languages:** English (default), European Portuguese, Spanish (Spain), Russian and Ukrainian.
- **Appearance:** 5 color themes and 4 layout skins.

Plain language in every text.

## 2. Stack

| Layer         | Technology                     | Role                                                       |
| ------------- | ------------------------------ | ---------------------------------------------------------- |
| Window/bridge | Wails 2.17 (WebView2)          | Native window, Go ↔ JavaScript bridge, file dialogs        |
| Back end      | Go 1.27 (Clean/Hexagonal)      | Rules, processes (SteamCMD, server), vault, backup, mod.io |
| Front end     | React 19 + TypeScript + Vite   | Screens                                                    |
| External      | SteamCMD, InsurgencyServer.exe | Install/update and run the server (separate processes)     |

## 3. Code layout

```
[React (WebView2)] ──Wails bindings──> [adapters/desktop] ──> [app/* (use cases)] ──> [domain/* (pure rules)]
                                               │
                                               ├──> adapters/vault    (Windows DPAPI)
                                               ├──> adapters/proc     (start, find and stop processes)
                                               ├──> adapters/modio    (mod.io API, HTTPS and fixed host only)
                                               ├──> adapters/netx     (HTTPS client with a fixed host list)
                                               ├──> adapters/rcon     (Source RCON on 127.0.0.1)
                                               ├──> adapters/netdiag  (port diagnosis: STUN, UPnP, A2S, Steam)
                                               ├──> adapters/updater  (GitHub Releases, signed manifest)
                                               └──> adapters/winsys   (autostart, Firewall, Windows info)
```

```
desktop/internal/
  domain/   pure rules, no I/O: profile and command line, log events, error catalog, SteamCMD progress,
            SteamID, INI, backup format, secret mask, RCON packets, .pak index, port diagnosis, reports
  app/      use cases: install, server runner, mods, backup, map change, chat, players, watchdog, updater
  adapters/ Windows, processes, network, files and the bridge with the screens (desktop)
desktop/frontend/src/
  features/ one folder per screen (dashboard, install, server, mods, backup, settings, feedback, about...)
  i18n/     the 5 languages (en-US is the reference)
  shared/   bridge with Go (backend.ts), components, hooks, formatting, appearance
```

## 4. Project rules

- **Rules only in Go** (`internal/domain` pure + `internal/app` use cases). The front end only shows and collects
  choices.
- The front end sends **ids and options**, never commands or free paths. Paths come only from Windows dialogs opened
  by Go. Go builds the server and SteamCMD arguments (list, no shell).
- Secrets never reach the front end in plain text, except when the user clicks "Show". They never go to logs. Every
  log line shown goes through `domain/logparse.Redact`.
- Errors are always an item of the `domain/problems` catalog (what happened, cause, solution).
- External links only by key (`links` in `adapters/desktop/system_api.go`).
- Internet only to fixed hosts (`adapters/netx`): Valve (SteamCMD), mod.io, GitHub (updates), Steam Web API (port
  diagnosis), Discord (only the user's webhook). STUN servers are UDP and also fixed.
- Never delete a user file without confirmation. Config files are written atomically, with a `.bak` copy before the
  first change.
- Screen texts in 5 languages (ADR-009); Go returns codes and keys, never text.
