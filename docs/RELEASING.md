# Releasing a version

Versions are published as **GitHub Releases** (ADR-032). Installed apps read the signed `latest.json` of the latest
release and update themselves.

## What you need

- Windows with Go, Node.js and NSIS (`winget install NSIS.NSIS`) for the installer.
- The GitHub CLI (`gh`) signed in with write access to the repository.
- The update signing key at `~/.sandstorm-server-assist/update-signing-key.pem`. It never goes into git. The scripts
  refuse a key that does not match the public key inside the app (`desktop/internal/adapters/updater/pubkey.go`).
  Without it, installed apps reject updates, so keep a safe copy.

## Steps

1. Bump the version in `desktop/wails.json` (`productVersion`) and `desktop/frontend/package.json`, then run
   `npm install` to update the lockfile.
2. Add a `## [x.y.z] - YYYY-MM-DD` section to `CHANGELOG.md`. It becomes the release notes.
3. Run the checks:

   ```bash
   npm test
   npm run lint
   npm run typecheck
   npm run audit
   npm run audit:go
   npm run audit:code
   ```

4. Build:

   ```bash
   npm run build:installer
   npm run package
   ```

5. Commit and push to `main`.
6. Check, then publish:

   ```bash
   npm run release:publish -- --dry-run
   npm run release:publish
   ```

   The script signs `latest.json`, writes `SHA256SUMS.txt`, checks that the commit is on `origin/main`, creates the
   release `v<version>` with the installer, the portable zip, `latest.json` and the checksums, and finally downloads
   `latest.json` from the same address the app uses to confirm it.

## Bridge for apps up to 0.1.10

Versions up to 0.1.10 check for updates on the author's old download server, not on GitHub. Until the bridge ends
(see the [roadmap](dev/ROADMAP.md)), publish each version there too, so those apps update once to a version that
uses GitHub:

```bash
npm run site:preview
SSA_SSH_HOST=user@host npm run site:publish
```

- `SSA_SSH_HOST` is the SSH address of that server; it is never written in the repository. It must be a name or
  address that the PC's `known_hosts` already knows.
- `SSA_SSH_KEY` changes the SSH key (the default is in `scripts/publish-site.mjs`); `SSA_SITE_URL` changes the address checked at the
  end.
- The script uploads the installer and the zip, swaps the page in one step, keeps the last 3 versions and uploads
  `latest.json` last.

The old server also runs the legacy port checker (`checker/`) used by apps up to 0.1.10. If `checker/` changes
during the bridge, rebuild that service on the server (`docker compose build sandstorm-check` and
`up -d sandstorm-check`). If the web server config changes, recreate its container with `--force-recreate`.
