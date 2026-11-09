// Gera o latest.json assinado para o instalador em desktop/build/dist (ADR-014). É o arquivo que os apps
// instalados leem para se atualizar sozinhos; vai para a pasta de downloads da página, DEPOIS do instalador
// (assim ninguém vê a versão nova antes de o instalador estar no lugar).
//
// Uso: npm run update:manifest              → versão do desktop/wails.json
//      npm run update:manifest -- 0.2.0     → outra versão (o instalador dela precisa estar em build/dist)
// Saída: desktop/build/dist/latest.json
import { createHash, createPrivateKey, sign } from 'node:crypto';
import { existsSync, readFileSync, writeFileSync } from 'node:fs';
import { join, resolve } from 'node:path';
import { KEY_PATH, embeddedPublicKey, publicKeyBase64 } from './update-key.mjs';

const ROOT = resolve(import.meta.dirname, '..');
const DIST = join(ROOT, 'desktop', 'build', 'dist');

function fail(msg) {
  console.error(`✖ ${msg}`);
  process.exit(1);
}

const version =
  process.argv[2] ??
  JSON.parse(readFileSync(join(ROOT, 'desktop', 'wails.json'), 'utf8')).info.productVersion;
if (!/^\d+\.\d+\.\d+$/.test(version)) fail(`Versão inválida: ${version}`);

let signingKey;
try {
  signingKey = readFileSync(KEY_PATH, 'utf8');
} catch {
  fail(`Chave de assinatura não encontrada em ${KEY_PATH} (veja scripts/update-key.mjs).`);
}
if (publicKeyBase64(signingKey) !== embeddedPublicKey()) {
  fail(
    'A chave de assinatura não é a do app (pubkey.go): os apps instalados recusariam a atualização.',
  );
}

// Mesmo nome que o Go espera (update.InstallerName).
const file = `sandstorm-server-assist-${version}-windows-setup.exe`;
const path = join(DIST, file);
if (!existsSync(path)) fail(`Instalador não encontrado: ${path} (rode "npm run build:installer").`);
const data = readFileSync(path);

const manifest = {
  version,
  file,
  sha256: createHash('sha256').update(data).digest('hex'),
  size: data.length,
  date: new Date().toISOString().slice(0, 10),
};
// A mensagem assinada é a mesma do Go (update.Manifest.Message).
const message = [
  'sandstorm-server-assist-update-v1',
  manifest.version,
  manifest.file,
  manifest.sha256,
  manifest.size,
  '',
].join('\n');
manifest.signature = sign(null, Buffer.from(message), createPrivateKey(signingKey)).toString(
  'base64',
);
writeFileSync(join(DIST, 'latest.json'), `${JSON.stringify(manifest, null, 2)}\n`);
console.log(`✔ ${join(DIST, 'latest.json')} (versão ${version}, sha256 ${manifest.sha256})`);
