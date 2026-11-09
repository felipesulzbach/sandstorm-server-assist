# Contributing

Thanks for helping! Bug reports, translations and pull requests are welcome.

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
