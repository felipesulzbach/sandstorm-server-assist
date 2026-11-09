// Publica a página de download (site/) e os arquivos da versão no servidor do game-in, em /sandstorm/ (ADR-016),
// e o latest.json assinado que os apps instalados leem para se atualizar sozinhos (ADR-014). Mesmo modelo do
// midia-multi-tools, mas a partir de desktop/build/dist (ainda não há repositório nem CI).
//
// Uso: npm run site:preview                → só monta a página em desktop/build/dist/site (nada é enviado)
//      npm run site:publish                → monta e publica a versão do desktop/wails.json
//      npm run site:publish -- 0.2.0       → outra versão (os arquivos dela precisam estar em build/dist)
// Antes: npm run build:installer && npm run package.
// Variáveis: SSA_SSH_HOST (obrigatória para publicar, ex.: usuario@endereço; não fica no código porque o projeto
// é público), SSA_SSH_KEY (padrão ~/.ssh/game-in.pem), SSA_SITE_URL (padrão https://game-in.duckdns.org/sandstorm/).
import { spawnSync } from 'node:child_process';
import { createHash, createPrivateKey, sign } from 'node:crypto';
import {
  cpSync,
  existsSync,
  mkdirSync,
  readFileSync,
  readdirSync,
  rmSync,
  writeFileSync,
} from 'node:fs';
import { homedir } from 'node:os';
import { join, resolve } from 'node:path';
import { KEY_PATH, embeddedPublicKey, publicKeyBase64 } from './update-key.mjs';

const ROOT = resolve(import.meta.dirname, '..');
const DIST = join(ROOT, 'desktop', 'build', 'dist');
const OUT = join(DIST, 'site');
const APP = 'sandstorm-server-assist';
const REMOTE = 'game-in/deploy/sandstorm'; // montado no Caddy do game-in como /srv/sandstorm
const SITE_URL = process.env.SSA_SITE_URL ?? 'https://game-in.duckdns.org/sandstorm/';
const KEEP = 3; // versões guardadas no servidor

const args = process.argv.slice(2);
const preview = args.includes('--preview');
const version =
  args.find((a) => !a.startsWith('--')) ??
  JSON.parse(readFileSync(join(ROOT, 'desktop', 'wails.json'), 'utf8')).info.productVersion;

function fail(msg) {
  console.error(`✖ ${msg}`);
  process.exit(1);
}

if (!/^\d+\.\d+\.\d+$/.test(version)) fail(`Versão inválida: ${version}`);

// Mesmos nomes do build (scripts/desktop.mjs) e do Go (update.InstallerName).
const exe = `${APP}-${version}-windows-setup.exe`;
const zip = `${APP}-${version}-windows-portable.zip`;
for (const f of [exe, zip]) {
  if (!existsSync(join(DIST, f)))
    fail(`Falta ${join(DIST, f)} (rode "npm run build:installer" e "npm run package").`);
}

const info = (f) => {
  const data = readFileSync(join(DIST, f));
  return {
    sha: createHash('sha256').update(data).digest('hex'),
    bytes: data.length,
    size: `${(data.length / 1048576).toFixed(1)} MB`,
  };
};
const iExe = info(exe);
const iZip = info(zip);

// ---- 1. Página ----

const escapeHtml = (s) =>
  s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');

/** Converte o Markdown simples de TERMS.md / PRIVACY.md (títulos, parágrafos, listas, tabelas, links). */
function markdown(md) {
  const inline = (s) =>
    escapeHtml(s)
      .replace(/`([^`]+)`/g, '<code>$1</code>')
      .replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>')
      .replace(/\[([^\]]+)\]\(([^)]+)\)/g, (m, text, href) => {
        const map = {
          'TERMS.md': 'terms.html',
          'PRIVACY.md': 'privacy.html',
          LICENSE: 'https://www.apache.org/licenses/LICENSE-2.0',
          'DISCLAIMER.md': 'terms.html',
        };
        const url = map[href] ?? href;
        if (!/^(https:\/\/|[a-z]+\.html$)/.test(url)) return text;
        const ext = url.startsWith('https://') ? ' rel="noopener noreferrer"' : '';
        return `<a href="${url}"${ext}>${text}</a>`;
      });
  const out = [];
  const lines = md.replace(/\r/g, '').split('\n');
  for (let i = 0; i < lines.length;) {
    const line = lines[i];
    if (!line.trim()) {
      i++;
    } else if (/^#{1,3} /.test(line)) {
      const level = line.match(/^#+/)[0].length;
      out.push(`<h${level}>${inline(line.slice(level + 1))}</h${level}>`);
      i++;
    } else if (/^\|/.test(line)) {
      const rows = [];
      while (i < lines.length && /^\|/.test(lines[i])) rows.push(lines[i++]);
      const cells = (r) =>
        r
          .split('|')
          .slice(1, -1)
          .map((c) => inline(c.trim()));
      const body = rows.slice(2).map(
        (r) =>
          `<tr>${cells(r)
            .map((c) => `<td>${c}</td>`)
            .join('')}</tr>`,
      );
      out.push(
        `<table><thead><tr>${cells(rows[0])
          .map((c) => `<th>${c}</th>`)
          .join('')}</tr></thead><tbody>${body.join('')}</tbody></table>`,
      );
    } else if (/^(- |\d+\. )/.test(line)) {
      const ordered = /^\d+\. /.test(line);
      const items = [];
      while (i < lines.length && (/^(- |\d+\. )/.test(lines[i]) || /^ {2,}\S/.test(lines[i]))) {
        if (/^ {2,}\S/.test(lines[i])) items[items.length - 1] += ` ${lines[i].trim()}`;
        else items.push(lines[i].replace(/^(- |\d+\. )/, ''));
        i++;
      }
      const tag = ordered ? 'ol' : 'ul';
      out.push(`<${tag}>${items.map((it) => `<li>${inline(it)}</li>`).join('')}</${tag}>`);
    } else {
      const para = [];
      while (i < lines.length && lines[i].trim() && !/^(#|\||- |\d+\. )/.test(lines[i]))
        para.push(lines[i++].trim());
      out.push(`<p>${inline(para.join(' '))}</p>`);
    }
  }
  return out.join('\n');
}

function docPage(title, mdFile) {
  const body = markdown(readFileSync(join(ROOT, mdFile), 'utf8'));
  return `<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>${title} · Sandstorm Server Assist</title>
    <link rel="icon" href="favicon.svg" type="image/svg+xml" />
    <link rel="stylesheet" href="style.css" />
  </head>
  <body>
    <main class="doc">
      <a class="back" href="./">← Sandstorm Server Assist</a>
${body}
    </main>
  </body>
</html>
`;
}

rmSync(OUT, { recursive: true, force: true });
mkdirSync(OUT, { recursive: true });
cpSync(join(ROOT, 'site'), OUT, { recursive: true });
const values = {
  VERSION: version,
  SIZE_EXE: iExe.size,
  SIZE_ZIP: iZip.size,
  SHA_EXE: iExe.sha,
  SHA_ZIP: iZip.sha,
  DATE: new Date().toISOString().slice(0, 10),
};
let html = readFileSync(join(OUT, 'index.html'), 'utf8').replace(/<!--[\s\S]*?-->\n/, '');
html = html.replace(/\{\{(\w+)\}\}/g, (m, k) => values[k] ?? m);
if (/\{\{\w+\}\}/.test(html)) fail('Sobrou marcador {{...}} na página');
writeFileSync(join(OUT, 'index.html'), html);
writeFileSync(join(OUT, 'terms.html'), docPage('Terms of Use', 'TERMS.md'));
writeFileSync(join(OUT, 'privacy.html'), docPage('Privacy Notice', 'PRIVACY.md'));

// Manifesto da atualização automática, assinado (mesma mensagem do Go: update.Manifest.Message).
let signingKey;
try {
  signingKey = readFileSync(KEY_PATH, 'utf8');
} catch {
  fail(`Chave de assinatura não encontrada em ${KEY_PATH} (veja scripts/update-key.mjs).`);
}
if (publicKeyBase64(signingKey) !== embeddedPublicKey())
  fail(
    'A chave de assinatura não é a do app (pubkey.go): os apps instalados recusariam a atualização.',
  );
const manifest = { version, file: exe, sha256: iExe.sha, size: iExe.bytes, date: values.DATE };
const message = ['sandstorm-server-assist-update-v1', version, exe, iExe.sha, iExe.bytes, ''].join(
  '\n',
);
manifest.signature = sign(null, Buffer.from(message), createPrivateKey(signingKey)).toString(
  'base64',
);
writeFileSync(join(DIST, 'latest.json'), `${JSON.stringify(manifest, null, 2)}\n`);

if (preview) {
  console.log(`✔ Página montada em ${OUT} (nada enviado).`);
  process.exit(0);
}

// ---- 2. Envio ----

const HOST = process.env.SSA_SSH_HOST;
if (!HOST) fail('Defina SSA_SSH_HOST (usuario@endereço do servidor do game-in) para publicar.');
const KEY = process.env.SSA_SSH_KEY ?? join(homedir(), '.ssh', 'game-in.pem');

function run(cmd, cmdArgs) {
  const r = spawnSync(cmd, cmdArgs, { stdio: ['ignore', 'pipe', 'inherit'], encoding: 'utf8' });
  if (r.status !== 0) fail(`${cmd} ${cmdArgs[0] ?? ''} falhou`);
  return r.stdout ?? '';
}
const sshOpts = ['-i', KEY, '-o', 'BatchMode=yes', '-o', 'ConnectTimeout=20'];
const ssh = (command) => run('ssh', [...sshOpts, HOST, command]);
const scp = (files, dest) => run('scp', [...sshOpts, '-q', ...files, `${HOST}:${dest}`]);

// Os arquivos chegam com o nome da versão e só depois viram os nomes fixos da página (quem baixa no meio da
// publicação nunca recebe um arquivo pela metade).
console.log(`Enviando para ${REMOTE}...`);
ssh(`mkdir -p ${REMOTE}/downloads && rm -rf ${REMOTE}/site.new && mkdir -p ${REMOTE}/site.new`);
scp([join(DIST, exe), join(DIST, zip)], `${REMOTE}/downloads/`);
ssh(
  `cd ${REMOTE}/downloads && cp -f ${exe} ${APP}-setup.exe.tmp && mv -f ${APP}-setup.exe.tmp ${APP}-setup.exe` +
    ` && cp -f ${zip} ${APP}-portable.zip.tmp && mv -f ${APP}-portable.zip.tmp ${APP}-portable.zip`,
);
scp(
  readdirSync(OUT).map((f) => join(OUT, f)),
  `${REMOTE}/site.new/`,
);
ssh(
  `cd ${REMOTE} && rm -rf site.old && { [ ! -d site ] || mv site site.old; } && mv site.new site && rm -rf site.old`,
);
ssh(
  `cd ${REMOTE}/downloads && ls -t ${APP}-*-windows-setup.exe | tail -n +${KEEP + 1} | xargs -r rm -f` +
    ` && ls -t ${APP}-*-windows-portable.zip | tail -n +${KEEP + 1} | xargs -r rm -f`,
);
// Por último o manifesto: só agora os apps instalados veem a versão nova.
scp([join(DIST, 'latest.json')], `${REMOTE}/downloads/latest.json.tmp`);
ssh(`cd ${REMOTE}/downloads && mv -f latest.json.tmp latest.json`);

// ---- 3. Conferência pela internet ----
const page = run('curl', ['-fsS', SITE_URL]);
if (!page.includes(`Sandstorm Server Assist ${version}`))
  fail(`A página em ${SITE_URL} não mostra a versão ${version}`);
const head = run('curl', ['-fsSI', `${SITE_URL}downloads/${APP}-setup.exe`]);
if (!/content-disposition:\s*attachment/i.test(head))
  fail('O instalador não está sendo servido como download');
const latest = JSON.parse(run('curl', ['-fsS', `${SITE_URL}downloads/latest.json`]));
if (latest.version !== version || latest.signature !== manifest.signature)
  fail('O latest.json publicado não é o desta versão');
console.log(`✔ Publicado: ${SITE_URL} (versão ${version}, atualização automática liberada)`);
console.log(`  instalador ${iExe.size}  sha256 ${iExe.sha}`);
console.log(`  portátil   ${iZip.size}  sha256 ${iZip.sha}`);
