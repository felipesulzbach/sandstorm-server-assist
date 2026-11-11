---
name: security-audit
description: Checks and fixes vulnerabilities in the code and dependencies of Sandstorm Server Assist (Wails desktop app with Go, React front end on npm, DPAPI vault, SteamCMD and game server processes, network diagnosis). Use after ANY code change, when adding or updating dependencies, when a security audit is requested or when the security CI fails.
---

# Vulnerability audit and fixes

The risks of this app: server secrets, external processes (SteamCMD, InsurgencyServer.exe), network (mod.io, Valve,
GitHub, STUN, UPnP on the local network), backups and the feedback reports.

| Command              | What it does                                                                                                  |
| -------------------- | ------------------------------------------------------------------------------------------------------------- |
| `npm run audit`      | `audit-ci` (fails from _moderate_, config in `audit-ci.jsonc`) + `npm audit signatures`                       |
| `npm run audit:go`   | `govulncheck ./...` in `desktop/` and `checker/` (reachable only; exceptions in `govulncheck-allowlist.json`) |
| `npm run audit:code` | `gosec ./...` in `desktop/` + ESLint security rules in `desktop/frontend/src`                                 |

## What to run when

| Change                          | Minimum check                              |
| ------------------------------- | ------------------------------------------ |
| Any code                        | `npm run audit:code` + diff review (below) |
| npm dependency added or updated | + `npm run audit`                          |
| `go.mod`, Go or Wails version   | + `npm run audit:go`                       |

## Diff review (what the tools miss)

- **Secrets:** join password, RCON, GSLT, GameStats, mod.io API key and e-mail only in the vault (`adapters/vault`,
  DPAPI with entropy) or the encrypted backup. Never in logs, `.ini`, `.bat`, error messages or JSON to the screen.
  Every log text goes through `domain/redact` (`logbuf.Add` does it). Errors go through `problems.New/Wrap`. The screen
  gets a value only with "Show" (`RevealSecret`).
- **Processes:** `InsurgencyServer.exe` only with the line built by `server.BuildLaunch` (fields validated by regex: no
  quotes, spaces, `?`, `&`). SteamCMD only with `steamcmd.Args` and only after `trust.VerifySigner(..., "Valve")`.
  Nothing from the screen becomes a free argument. Stop only processes of the server folder (`proc.FindInDir`).
- **Exposed methods (Wails):** every exported method of the `*API` types in `adapters/desktop` can be called from JS.
  They take ids, keys and options; paths only from dialogs opened in Go; links only by key from the `links` list.
- **Files:** atomic writes (`fsx.WriteAtomic`) with a `.bak` before the first change of a user file. Chosen folders go
  through `fsx.SafeDir`. Reads with a limit (`fsx.ReadLimited`).
- **Backup:** only names from `backup.ConfigFiles` + `modio/user.json`; size limits when opening; bounded Argon2id
  parameters; authenticated header.
- **Network:** HTTPS only, to `netx.AllowedHosts`; redirects checked. UDP only to the fixed STUN servers, to the local
  server (`127.0.0.1`) and to the user's own public IP (A2S). UPnP only to private IPv4 addresses of the device that
  answered SSDP, with no proxy and no redirects; XML and SOAP bodies size-limited.
- **Feedback reports:** the log in the zip is already secret-masked and also has IP addresses and SteamIDs masked
  (`report.Scrub`); the image is chosen in a Go dialog, size-limited and checked by content; the app opens only the
  fixed GitHub issue URL.
- **Front end:** no `dangerouslySetInnerHTML`, `eval` or HTML from logs or mod names; build CSP without `unsafe-eval`
  and without external origins.
- **Public project:** no key, token or personal data in the repository (the mod.io key is not in code, ADR-010). The
  maintainer's private server notes (`resumo_sessao.md`) stay out of git.

## Fixing

1. npm: `npm audit fix` (never `--force`); Go: `go get <module>@<fixed version>` + `go mod tidy`; for the Go standard
   library, raise the `go` line in `go.mod` (and the `checker/Dockerfile` image).
2. A gosec suppression only with `// #nosec G<rule> -- reason` on the same line.
3. With no fix available: an exception with a reason and an expiry date (`audit-ci.jsonc` or
   `govulncheck-allowlist.json`).
4. Validate: `npm test`, `npm run lint`, `npm run typecheck`, `npm run build`.
5. Record below the decisions that are not obvious.

## Decision history

- 2026-11-09: project created; gosec suppressions with a reason: G204 (server, SteamCMD and explorer.exe processes
  with profile paths), G304 (reads of paths built by the app), G103 (DPAPI and WinVerifyTrust need pointers), G115
  (bounded sizes).
- 2026-11-11: Go 1.27.4 for GO-2026-6859 (HTTP/1 server crash); the app has no HTTP server, the legacy `checker/` has.
- 2026-11-11: port diagnosis (`adapters/netdiag`): G107/G704 suppressed for UPnP requests, with the address checked by
  `lanURL` (private IPv4 of the device that answered). Feedback zip: G304 for the image chosen in the Windows dialog.
