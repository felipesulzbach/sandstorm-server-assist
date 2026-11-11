// Ponte com o Go. Única parte do front que importa o código gerado pelo Wails (wailsjs/).
// O front manda ids e opções; caminhos só vêm de diálogos abertos pelo Go.
import * as Sys from '../../wailsjs/go/desktop/SystemAPI';
import * as Cfg from '../../wailsjs/go/desktop/ConfigAPI';
import * as Inst from '../../wailsjs/go/desktop/InstallAPI';
import * as Srv from '../../wailsjs/go/desktop/ServerAPI';
import * as Logs from '../../wailsjs/go/desktop/LogAPI';
import * as Bak from '../../wailsjs/go/desktop/BackupAPI';
import * as Rep from '../../wailsjs/go/desktop/ReportAPI';
import {
  EventsOn,
  Quit,
  WindowIsMaximised,
  WindowMinimise,
  WindowToggleMaximise,
} from '../../wailsjs/runtime/runtime';
import type {
  admin,
  backupsvc,
  config,
  desktop,
  installer,
  logbuf,
  mapchange,
  modio,
  netcheck,
  players,
  runner,
  server,
  updater,
  winsys,
} from '../../wailsjs/go/models';

export type Profile = server.Profile;
export type Admin = server.Admin;
export type Status = runner.Status;
export type InstallState = installer.State;
export type InstallInfo = installer.Info;
export type LogLine = logbuf.Line;
export type LogRun = desktop.LogRun;
export type LogCleanup = desktop.LogCleanup;
export type ModsView = desktop.ModsView;
export type ModRow = desktop.ModRow;
export type Step = desktop.Step;
export type Connect = desktop.Connect;
export type MapsInfo = desktop.MapsInfo;
export type MapChange = mapchange.Info;
export type Scenario = server.Scenario;
export type PortCheck = netcheck.Result;
export type Mutator = server.Mutator;
export type ModMutator = modio.ModMutator;
export type ModScenario = modio.ModScenario;
export type RuleDef = server.RuleDef;
export type RulesInfo = desktop.RulesInfo;
export type CycleLine = desktop.CycleLine;
export type LocalAddress = winsys.Address;
export type ProfileRef = config.ProfileRef;
export type PlayersView = players.View;
export type Player = admin.Player;
export type Ban = admin.Ban;
export type AutomationInfo = desktop.AutomationInfo;
export type FirewallInfo = desktop.FirewallInfo;
export type Lighting = '' | 'Day' | 'Night';
export type AppInfo = desktop.Info;
export type Settings = config.Settings;
export type ExistingScan = desktop.ExistingScan;
export type ApplyResult = config.ApplyResult;
export type SyncDiff = server.Diff;
export type UpdateState = updater.State;
export type Terms = desktop.Terms;
export type BackupSummary = backupsvc.Summary;
export type CreateResult = desktop.CreateResult;
export type ReportImage = desktop.ImageInfo;
export type ReportResult = desktop.ReportResult;
export type ReportKind = 'bug' | 'idea';
export type SecretKey =
  | 'serverPassword'
  | 'rconPassword'
  | 'gsltToken'
  | 'gameStatsToken'
  | 'modioEmail'
  | 'modioApiKey'
  | 'discordWebhook';

export const backend = {
  info: () => Sys.Info(),
  openLink: (key: string) => Sys.OpenLink(key),
  problemLinks: (code: string) => Sys.ProblemLinks(code),
  settings: () => Sys.Settings(),
  setLanguage: (lang: string) => Sys.SetLanguage(lang),
  setShowConsole: (on: boolean) => Sys.SetShowConsole(on),
  copyText: (text: string) => Sys.CopyText(text),
  checklist: () => Sys.Checklist(),
  terms: () => Sys.Terms(),
  acceptTerms: (version: string) => Sys.AcceptTerms(version),
  updateState: () => Sys.UpdateState(),
  checkUpdates: () => Sys.CheckUpdates(),
  installUpdate: () => Sys.InstallUpdate(),
  setAutostart: (withWindows: boolean, startServer: boolean) =>
    Sys.SetAutostart(withWindows, startServer),
  firewall: () => Sys.Firewall(),
  createFirewallRule: () => Sys.CreateFirewallRule(),

  profile: () => Cfg.Profile(),
  saveProfile: (p: Profile) => Cfg.SaveProfile(p),
  defaultProfile: () => Cfg.DefaultProfile(),
  maps: () => Cfg.Maps(),
  secretsState: () => Cfg.SecretsState() as Promise<Record<SecretKey, boolean>>,
  setSecret: (key: SecretKey, value: string) =>
    Cfg.SetSecret(key, value) as Promise<Record<SecretKey, boolean>>,
  revealSecret: (key: SecretKey) => Cfg.RevealSecret(key),
  extractSteamId: (text: string) => Cfg.ExtractSteamID(text),
  chooseServerDir: () => Cfg.ChooseServerDir(),
  chooseSteamCmdDir: () => Cfg.ChooseSteamCMDDir(),
  scanExisting: () => Cfg.ScanExisting(),
  applyExisting: (keepMine: string[]) => Cfg.ApplyExisting(keepMine),
  rconInIni: () => Cfg.RconInIni(),
  moveRconToVault: () => Cfg.MoveRconToVault(),
  mutators: () => Cfg.Mutators(),
  modMutators: () => Cfg.ModMutators(),
  modScenarios: () => Cfg.ModScenarios(),
  rules: () => Cfg.Rules(),
  mapCycleLines: (lines: string[]) => Cfg.MapCycleLines(lines),
  setMapCycleLighting: (lines: string[], index: number, lighting: Lighting) =>
    Cfg.SetMapCycleLighting(lines, index, lighting),
  localAddresses: () => Cfg.LocalAddresses(),
  profiles: () => Cfg.Profiles(),
  createProfile: (name: string, copyCurrent: boolean) => Cfg.CreateProfile(name, copyCurrent),
  renameProfile: (id: string, name: string) => Cfg.RenameProfile(id, name),
  deleteProfile: (id: string) => Cfg.DeleteProfile(id),
  switchProfile: (id: string) => Cfg.SwitchProfile(id),

  installState: () => Inst.State(),
  installInfo: () => Inst.Info(),
  fetchSteamCmd: () => Inst.FetchSteamCMD(),
  installServer: (validate: boolean) => Inst.InstallServer(validate),
  cancelInstall: () => Inst.Cancel(),

  status: () => Srv.Status(),
  start: () => Srv.Start(),
  stop: () => Srv.Stop(),
  restart: () => Srv.Restart(),
  connect: () => Srv.Connect(),
  mapChange: () => Srv.MapChange(),
  scenarios: () => Srv.Scenarios(),
  changeMap: (id: string, lighting: Lighting) => Srv.ChangeMap(id, lighting),
  players: () => Srv.Players(),
  kick: (netId: string, reason: string) => Srv.Kick(netId, reason),
  ban: (netId: string, minutes: number, reason: string) => Srv.Ban(netId, minutes, reason),
  unban: (netId: string) => Srv.Unban(netId),
  say: (text: string) => Srv.Say(text),
  automation: () => Srv.Automation(),
  testDiscord: () => Srv.TestDiscord(),
  checkPorts: () => Srv.CheckPorts(),
  mods: () => Srv.Mods(),
  refreshSubscriptions: () => Srv.RefreshSubscriptions(),
  requestCode: () => Srv.RequestCode(),
  authorize: (code: string) => Srv.Authorize(code),
  detectApiKey: () => Srv.DetectAPIKey(),

  recentLogs: (limit: number) => Logs.Recent(limit),
  clearLogs: () => Logs.Clear(),
  reportScreenError: (where: string, detail: string) => Logs.ReportScreenError(where, detail),
  openLogFolder: () => Logs.OpenFolder(),
  logRuns: () => Logs.Runs(),
  readLogRun: (id: string) => Logs.ReadRun(id),
  logCleanup: () => Logs.Cleanup(),
  setLogCleanup: (enabled: boolean, keep: number) => Logs.SetCleanup(enabled, keep),
  cleanOldLogs: () => Logs.CleanOld(),

  createBackup: (password: string, includeModio: boolean) => Bak.Create(password, includeModio),
  chooseBackupFile: () => Bak.ChooseFile(),

  chooseReportImage: () => Rep.ChooseImage(),
  clearReportImage: () => Rep.ClearImage(),
  createReport: (kind: ReportKind, title: string, description: string, includeLog: boolean) =>
    Rep.Create({ kind, title, description, includeLog }),
  showReport: () => Rep.ShowLast(),
  inspectBackup: (password: string) => Bak.Inspect(password),
  chooseRestoreDir: () => Bak.ChooseRestoreDir(),
  restoreBackup: (password: string, restoreModio: boolean) => Bak.Restore(password, restoreModio),
};

/** Assina um evento do Go. Devolve a função que cancela. */
export function on<T>(
  event: 'server:status' | 'install:state' | 'log:lines' | 'update:state',
  fn: (v: T) => void,
) {
  return EventsOn(event, fn);
}

export const windowControls = {
  minimise: () => WindowMinimise(),
  toggleMaximise: () => WindowToggleMaximise(),
  isMaximised: () => WindowIsMaximised(),
  close: () => Quit(),
};
