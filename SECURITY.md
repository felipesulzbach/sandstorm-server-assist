# Security policy

Please **do not open a public issue** for security problems (for example a way to read the password vault,
run arbitrary commands through the app, or leak tokens in logs or backups).

Report them privately through the repository's "Report a vulnerability" (GitHub Security Advisories) and
include the steps to reproduce. You will get an answer as soon as possible.

## What the app protects

- Passwords, tokens and the mod.io API key: Windows DPAPI (current Windows user, this PC only), never in config
  files or logs.
- Backups: Argon2id + AES-256-GCM with a password chosen by the user (minimum 10 characters).
- SteamCMD: only run after its Authenticode signature is verified as Valve's.
- Network: HTTPS only, to a fixed list of hosts (Valve's SteamCMD CDN and the mod.io API).

## Known limits

- The game server only accepts its passwords and tokens on the command line, which other programs running as the
  same Windows user can read.
- Anyone who has your backup file **and** its password can restore your server settings and secrets.
