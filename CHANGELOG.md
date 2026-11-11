# Changelog

All notable changes to Sandstorm Server Assist. Versions follow `MAJOR.MINOR.PATCH`. The GitHub release notes of
each version come from its section here (`npm run release:publish`).

## [0.2.0] - 2026-11-11

### Added

- **Feedback page:** report a problem or suggest an improvement. The app prepares a `.zip` (your text, the app
  log with passwords, tokens, IP addresses and SteamIDs masked, and an optional screenshot up to 5 MB) and opens
  the GitHub form already filled in. Title, description and screenshot are optional. Nothing is sent by the app.
- **Port diagnosis from your PC:** checks that the server answers, the Windows Firewall rules, the router port
  forwards (UPnP), your public IP and CGNAT (STUN), the route through your public IP and whether Steam lists your
  server. It also explains, step by step, how to confirm with a friend outside your network.
- **Credits and project links** in About.

### Changed

- **Updates now come from GitHub Releases.** Installed apps read the signed `latest.json` of the latest release
  and download the installer of that exact version. Signatures are checked as before.
- The port test no longer uses the project server.
- Terms of use and privacy notice updated (version 2026-11-11): the app asks you to accept them again.
- Built with Go 1.27.4 (security fix GO-2026-6859).

## [0.1.10] - 2026-11-11

- Map change uses the mods loaded on the server, with separate map and mode lists.
- Map cycle entries whose mod is not installed are flagged, with a "Remove" button.
- Messages are plain text again (the game showed formatting tags as text).

## [0.1.9] - 2026-11-11

- The dashboard shows the new map after changing to a mod map.

## [0.1.8] - 2026-11-11

- The dashboard no longer stays on "loading the map" after a map change.
- A harmless warning of the Border map no longer counts as an error.

## [0.1.7] - 2026-11-10

- A reopened app recognises the server it started.

## [0.1.6] - 2026-11-10

- Fixed a black screen when the player list was empty; screen errors now show a message instead.

## [0.1.5] - 2026-11-10

- Mod map picker for the map change.

## [0.1.4] - 2026-11-10

- Mutators read from installed mods.

## [0.1.3] and earlier

- First public versions: installation with SteamCMD, server settings, mod.io mods, dashboard, logs, backups,
  automatic updates, terms of use and privacy notice in 5 languages.
