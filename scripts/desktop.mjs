// Roda o Wails com a mesma versão do desktop/go.mod (sem instalar a CLI global) — ADR-001.
// Uso: node scripts/desktop.mjs dev | build [--nsis] | package | generate | vet | test
import { spawnSync } from 'node:child_process';
import { copyFileSync, existsSync, mkdirSync, readFileSync, rmSync, writeFileSync } from 'node:fs';
import { join, resolve } from 'node:path';

const ROOT = resolve(import.meta.dirname, '..');
const DESKTOP = join(ROOT, 'desktop');
const APP = 'sandstorm-server-assist';
const exe = (n) => (process.platform === 'win32' ? `${n}.exe` : n);

const gomod = readFileSync(join(DESKTOP, 'go.mod'), 'utf8');
const wailsVersion = gomod.match(/github\.com\/wailsapp\/wails\/v2 (v[\w.-]+)/)?.[1];
if (!wailsVersion) {
  console.error('✖ Versão do Wails não encontrada em desktop/go.mod');
  process.exit(1);
}
const WAILS = `github.com/wailsapp/wails/v2/cmd/wails@${wailsVersion}`;
const env = { ...process.env };

// O Wails chama o makensis pelo PATH; no Windows ele costuma ficar fora do PATH, na pasta padrão do NSIS.
for (const dir of ['C:/Program Files (x86)/NSIS', 'C:/Program Files/NSIS']) {
  if (existsSync(join(dir, 'makensis.exe'))) {
    env.PATH = `${dir};${env.PATH ?? env.Path ?? ''}`;
    delete env.Path;
    break;
  }
}

function run(cmd, args, opts = {}) {
  const r = spawnSync(cmd, args, { cwd: DESKTOP, stdio: 'inherit', env, ...opts });
  if (r.status !== 0) process.exit(r.status ?? 1);
}

const wails = (...args) => run('go', ['run', WAILS, ...args]);

// O embed do main.go (frontend/dist) precisa de pelo menos um arquivo para compilar antes do primeiro build do front.
function ensureDist() {
  const dist = join(DESKTOP, 'frontend', 'dist');
  mkdirSync(dist, { recursive: true });
  if (!existsSync(join(dist, 'index.html'))) writeFileSync(join(dist, '.placeholder'), '');
}

function syncWailsVersion(version) {
  const path = join(DESKTOP, 'wails.json');
  const cfg = JSON.parse(readFileSync(path, 'utf8'));
  if (cfg.info.productVersion !== version) {
    cfg.info.productVersion = version;
    writeFileSync(path, `${JSON.stringify(cfg, null, 2)}\n`);
  }
}

switch (process.argv[2]) {
  case 'dev':
    ensureDist();
    // O Wails 2.17 recusa no servidor de dev o host "wails.localhost" que a própria janela usa no Windows.
    env.WAILS_DEV_ALLOWED_HOSTS ??= 'wails.localhost';
    wails('dev', ...process.argv.slice(3));
    break;
  case 'generate':
    ensureDist();
    wails('generate', 'module');
    break;
  case 'build': {
    const pkgVersion = JSON.parse(
      readFileSync(join(DESKTOP, 'frontend', 'package.json'), 'utf8'),
    ).version;
    const appVersion = (process.env.APP_VERSION ?? pkgVersion).replace(/^v/, '');
    if (!/^\d+\.\d+\.\d+$/.test(appVersion)) {
      console.error(`✖ Versão inválida: ${appVersion}`);
      process.exit(1);
    }
    syncWailsVersion(appVersion);
    const nsis = process.argv.includes('--nsis');
    wails(
      'build',
      '-clean',
      '-trimpath',
      '-ldflags',
      `-s -w -X ${APP}/desktop/internal/adapters/desktop.AppVersion=${appVersion}`,
      // Instalador por usuário (sem administrador): permite a atualização automática (ADR-014).
      ...(nsis ? ['-nsis', '-installscope', 'user'] : []),
    );
    if (nsis) {
      // Mesmo nome que o Go espera no manifesto (update.InstallerName).
      const dist = join(DESKTOP, 'build', 'dist');
      mkdirSync(dist, { recursive: true });
      const setup = join(dist, `${APP}-${appVersion}-windows-setup.exe`);
      copyFileSync(join(DESKTOP, 'build', 'bin', `${APP}-amd64-installer.exe`), setup);
      console.log(`✔ Instalador: ${setup}`);
    }
    console.log(`✔ App ${appVersion} em ${join(DESKTOP, 'build', 'bin')}`);
    break;
  }
  case 'package': {
    const bin = join(DESKTOP, 'build', 'bin');
    const appExe = join(bin, exe(APP));
    if (!existsSync(appExe)) {
      console.error('✖ Rode "npm run build" antes.');
      process.exit(1);
    }
    const version = JSON.parse(readFileSync(join(DESKTOP, 'wails.json'), 'utf8')).info
      .productVersion;
    const dist = join(DESKTOP, 'build', 'dist');
    mkdirSync(dist, { recursive: true });
    const zip = join(dist, `${APP}-${version}-windows-portable.zip`);
    rmSync(zip, { force: true });
    // A licença Apache 2.0 pede LICENSE e NOTICE junto da distribuição; o aviso de isenção vai também.
    const docs = ['LICENSE', 'NOTICE', 'DISCLAIMER.md', 'TERMS.md', 'PRIVACY.md'];
    for (const f of docs) copyFileSync(join(ROOT, f), join(bin, f));
    const tar = join(process.env.SystemRoot ?? 'C:/Windows', 'System32', 'tar.exe');
    run(tar, ['-a', '-cf', zip, '-C', bin, exe(APP), ...docs]);
    console.log(`✔ ${zip}`);
    break;
  }
  case 'vet':
    ensureDist();
    run('go', ['vet', './...']);
    break;
  case 'test':
    ensureDist();
    run('go', ['test', '-count=1', './...']);
    break;
  default:
    console.error(
      'Uso: node scripts/desktop.mjs dev | build [--nsis] | package | generate | vet | test',
    );
    process.exit(1);
}
