# Changelog

All notable changes to Sandstorm Server Assist. Versions follow `MAJOR.MINOR.PATCH`. The GitHub release notes of
each version come from its section here (`npm run release:publish`).

## [0.2.2] - 2026-11-12

### Added

- **Game rules from the server:** when "Set the game rules here" is off and the server's `Game.ini` already sets
  rules, the Game rules tab says how many and offers **Use the server rules**.
- **Refused mutators:** when the server refuses a mutator on the current map (seen with Ultralethal in Checkpoint),
  the Dashboard and the mutator list in Server › Maps say which one. The game gives no reason.

### Fixed

- The Players card showed only one player: the server sends every player on one line, and the app read one per line.
  Now every player shows, with the highest score on top, and the bots fold into one row you can open.

## [0.2.1] - 2026-11-11

### Added

- **Mods with problems** on the Mods screen: the app compares the files of the installed mods and lists pairs that
  ship the same game file in different versions, which can crash the game or the server with "Serial size mismatch"
  when a map loads. Identical copies are not listed.
- **Find the mod of a crash:** paste the crash message a player got (or the server's) and the app names the
  installed mod that ships the file in it. When the server itself closes with "Serial size mismatch", the error card
  now names the file and the mod.
- Mods made for older versions of the game (older `.pak` format) are now read too: their mutators and maps show up,
  and they take part in the checks above.
- **Delete a downloaded mod** on the Mods screen, so the server downloads it again: **Delete and download again**
  (stops and starts the server if needed) or **Only delete** (the server downloads it on its next start, if you are
  still subscribed). Useful for a broken mod or to force a fresh copy.
- **Download new mods now:** subscribed mods that are not downloaded yet show above the Mods list with **Download
  now** (starts the server) or **Restart and download now** (asks first, since the match ends).

### Changed

- **Installation screen:** progress, result and errors now show in the card of the button you clicked, and the
  screen reacts as soon as you click (before, nothing showed while SteamCMD updated itself). It also says when the
  server already had the latest version.
- With SteamCMD already installed, the button is **Check SteamCMD**: SteamCMD updates itself every time the app uses
  it, so there is nothing to download again unless the file is damaged.

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
