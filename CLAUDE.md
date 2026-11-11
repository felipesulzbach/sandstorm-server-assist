# Sandstorm Server Assist

Free, open-source Windows desktop app to install, configure and run an Insurgency: Sandstorm dedicated server with
mod.io mods. `desktop/` is the Wails 2 project (Go, Clean/Hexagonal) and `desktop/frontend/` the React + Vite front
end (workspace `@ssa/frontend`). Overview in the README; architecture in `docs/ARCHITECTURE.md`; decisions in
`docs/DECISIONS.md`; plan in `docs/dev/ROADMAP.md`; session notes in `docs/dev/DEVLOG.md`; releases in
`docs/RELEASING.md`. The private notes of how the maintainer's server was first set up by hand (`resumo_sessao.md`
at the root) stay out of git: they are the source of the mod.io rules the app applies.

## Project rules

- **Keep the plan current:** when a step ends or something is decided, update `docs/dev/ROADMAP.md` (and
  `docs/DECISIONS.md` for architecture decisions). At the end of a session, add a short entry to
  `docs/dev/DEVLOG.md`. User-visible changes go into `CHANGELOG.md`.
- **Security on every code change:** run the `security-audit` skill (at least `npm run audit:code` and a review of the
  diff) and fix before finishing. Never `npm audit fix --force`.
- **Git:** commit and push only with the maintainer's approval.
- **Rules only in Go** (`internal/domain` pure + `internal/app` use cases). The front end only shows and collects
  choices.
- **The front end never sends a command or a free path:** it sends ids and options. Paths come only from Windows
  dialogs opened by Go. Go builds the server and SteamCMD arguments (list, no shell).
- **Secrets** (passwords, tokens, the mod.io session) only in the DPAPI vault (`adapters/vault`) or the encrypted
  backup. Never in logs, never in the front end unless the user clicks "Show". Every log line shown goes through
  `domain/logparse.Redact`.
- **Errors** always as an item of the `domain/problems` catalog (what happened, cause, solution).
- **External links** only by key (`links` in `adapters/desktop/system_api.go`). Internet only to the fixed hosts of
  `adapters/netx`.
- **Wails:** after changing exposed methods in `internal/adapters/desktop`, run `npm run gen:bindings`.
- **Screen texts in 5 languages** (ADR-009): en-US is the reference (`desktop/frontend/src/i18n/locales/en-US.ts`);
  pt-PT, es-ES, ru and uk with the same keys (`npm test` checks). Go has no user-facing text: it returns codes
  (`problems.*`) and keys (`validation.*`). Plain language.
- **Docs in English** (ADR-035). The README and `docs/GUIDE.md` also exist in pt-PT, es-ES, ru and uk
  (`README.<lang>.md`, `docs/<lang>/GUIDE.md`): when the English one changes, update the translations with the same
  sections and images (`npm test` checks the structure).
- **Public project** (ADR-008, ADR-012): Apache 2.0; no keys, tokens, IP addresses or e-mail addresses in code, tests
  or docs. The author's name appears only in the credits (README, About, `NOTICE`).
- **Terms and data** (ADR-015): if `TERMS.md`/`PRIVACY.md` or the `legal.*` texts change, bump `legal.Version`.
- **Appearance** (ADR-011): colors only through CSS variables; a new theme in `styles/themes.css` +
  `shared/appearance.ts`.

## Cautions

- The maintainer's real server folder has secrets in plain text (`Game.ini`, notes next to the `.bat`) and the mod.io
  session token (`user.json`): never display their values.
- Never start, stop or restart the maintainer's real game server unless asked.
- The update signing key (`~/.sandstorm-server-assist/update-signing-key.pem`) never goes into git or logs.

## Commands

- `npm install` → dependencies.
- `npm run dev` → app in development mode.
- `npm run build` → `desktop/build/bin/sandstorm-server-assist.exe`.
- `npm run build:installer` → NSIS installer in `desktop/build/dist` (ADR-014); `npm run package` → portable zip.
- `npm run release:publish` → GitHub Release with the signed `latest.json` (ADR-032, `docs/RELEASING.md`; publishing
  only with the maintainer's approval). `-- --dry-run` checks without publishing.
- `npm run site:preview` / `npm run site:publish` → old download page, bridge for apps up to 0.1.10 (ADR-016;
  publishing only with the maintainer's approval).
- `npm run update:manifest` → signed `latest.json` (key in `~/.sandstorm-server-assist/`, never in git).
- `npm test` (Go + front end + docs translations), `npm run lint`, `npm run typecheck`.
- `npm run audit`, `npm run audit:go`, `npm run audit:code`.
