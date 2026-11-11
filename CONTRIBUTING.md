# Contributing

Thanks for helping! Bug reports, ideas, translations and pull requests are welcome.

## Reporting problems and ideas

- In the app, open **Feedback**: it prepares a `.zip` with the details (the log is masked) and opens the GitHub form
  already filled in. Or use the [issue forms](https://github.com/felipesulzbach/sandstorm-server-assist/issues/new/choose).
- Security problems: please follow [SECURITY.md](SECURITY.md) instead of opening a public issue.

## Ground rules

- **Rules live in Go.** `desktop/internal/domain` is pure (no I/O) and tested; `desktop/internal/app` holds use
  cases; `desktop/internal/adapters` talks to Windows, processes, network and the UI bridge.
- **The UI never sends commands or free paths.** It sends ids and options; paths only come from Windows dialogs
  opened by Go. Process arguments are built in Go as lists (no shell).
- **Secrets** (passwords, tokens, mod.io session) only go to the DPAPI vault or the encrypted backup. Never log
  them; every log line shown in the UI goes through `domain/redact`.
- **Errors** are codes from `domain/problems`; their texts (what happened / cause / how to fix) are translated
  in the front-end locales.
- **External links** are opened by key (`links` in `adapters/desktop/system_api.go`). Network access is limited
  to the hosts in `adapters/netx` (HTTPS only).
- **All user-facing text** must exist in the 5 locales (`npm test` checks the keys).
- **Docs are in English** ([ADR-035](docs/DECISIONS.md)). `README.md` and `docs/GUIDE.md` also exist in pt-PT, es-ES, ru
  and uk: when you change the English file, update the translations with the same sections and images, or open an
  issue asking for help (`npm test` checks the structure). Screenshots live in `docs/images/` and must not show real
  IP addresses, server names, players or secrets.
- **User-visible changes** go into `CHANGELOG.md`. Architecture decisions go into `docs/DECISIONS.md`.

See [architecture](docs/ARCHITECTURE.md), [decisions](docs/DECISIONS.md), [roadmap](docs/dev/ROADMAP.md) and
[releasing](docs/RELEASING.md).

## Before opening a pull request

```bash
npm test
npm run lint
npm run typecheck
npm run audit:code
```

After changing methods exposed in `desktop/internal/adapters/desktop`, run `npm run gen:bindings`.

By contributing you agree that your contribution is licensed under the Apache License 2.0 (section 5 of the
[LICENSE](LICENSE)).
