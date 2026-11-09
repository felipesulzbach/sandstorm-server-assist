# Privacy Notice

Version **2026-11-09**. It is accepted together with the [Terms of Use](TERMS.md) the first time you open the
app.

## 1. In short

The app has no account, no ads and no telemetry. It does not send your data to the authors of the project.

Everything you type stays on your PC, except what goes to the services you choose to use (section 3).

## 2. What stays on your PC

- **Server settings and app preferences:** in `%APPDATA%\sandstorm-server-assist\config.json`.
- **Secrets:** passwords, tokens, the mod.io e-mail and API key. They are kept in
  `%APPDATA%\sandstorm-server-assist\vault.bin`, protected by Windows (DPAPI). Only your Windows user on that PC
  can open it.
- **Backups:** only where you save them, encrypted with the password you choose.
- **Server logs:** read from the server folder to show them in the app. Passwords and tokens are masked before
  they are shown.

## 3. What leaves your PC, and to whom

| Destination                  | What and why                                                                                    |
| ---------------------------- | ----------------------------------------------------------------------------------------------- |
| Valve (Steam)                | Download of SteamCMD and of the game server.                                                    |
| mod.io                       | Server account e-mail (to get the login code), the code and the server API key, for mods.       |
| Project download page        | Update check: asks only for the latest version. Like any website, it sees your IP and the time. |
| Sites of the links you click | They open in your browser, under each site's rules.                                             |

The update check sends no identifier of yours. The app accepts only HTTPS connections, and only to these hosts.

## 4. Your game server

The game server you run connects to Steam, GameStats, mod.io and the players, under the rules of those
services. Admin SteamIDs you add are written to the server files (`Admins.txt`).

If players join your server, you are responsible for their data on it (names, SteamIDs, chat, logs).

## 5. Your control

You can see and change everything in the app.

To delete everything:

1. Uninstall the app.
2. Delete the folder `%APPDATA%\sandstorm-server-assist`.

Tokens and keys can be revoked at any time on the Steam and mod.io sites.

## 6. Changes

If this notice changes, the app shows it again and asks for your acceptance.
