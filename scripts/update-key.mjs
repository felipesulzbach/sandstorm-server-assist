// Chave de assinatura das atualizações (ADR-014, mesmo modelo do midia-multi-tools).
// Uso: npm run update:keygen   → cria a chave (só uma vez) e mostra a chave pública para o app
//      npm run update:pubkey   → mostra a chave pública da chave existente
// A chave PRIVADA fica fora do repositório, em ~/.sandstorm-server-assist/update-signing-key.pem (SSA_UPDATE_KEY
// muda o caminho). Guarde uma cópia em lugar seguro: sem ela, os apps instalados não aceitam atualizações novas
// (seria preciso instalar à mão uma versão com outra chave pública).
import { createPrivateKey, createPublicKey, generateKeyPairSync } from 'node:crypto';
import { existsSync, mkdirSync, readFileSync, writeFileSync } from 'node:fs';
import { homedir } from 'node:os';
import { dirname, join, resolve } from 'node:path';
import { pathToFileURL } from 'node:url';

export const KEY_PATH =
  process.env.SSA_UPDATE_KEY ??
  join(homedir(), '.sandstorm-server-assist', 'update-signing-key.pem');

export const PUBKEY_GO = resolve(
  import.meta.dirname,
  '..',
  'desktop',
  'internal',
  'adapters',
  'updater',
  'pubkey.go',
);

/** Chave pública Ed25519 crua (32 bytes) em base64, como o Go espera. */
export function publicKeyBase64(privatePem) {
  const spki = createPublicKey(createPrivateKey(privatePem)).export({
    format: 'der',
    type: 'spki',
  });
  return spki.subarray(spki.length - 32).toString('base64');
}

/** Chave pública embutida no app (pubkey.go). */
export function embeddedPublicKey() {
  return readFileSync(PUBKEY_GO, 'utf8').match(/publicKeyB64 = "([^"]*)"/)?.[1] ?? '';
}

if (import.meta.url === pathToFileURL(process.argv[1]).href) {
  const cmd = process.argv[2];
  if (cmd === 'keygen') {
    if (existsSync(KEY_PATH)) {
      console.error(
        `✖ Já existe uma chave em ${KEY_PATH}. Não sobrescrevo (os apps instalados dependem dela).`,
      );
      process.exit(1);
    }
    const { privateKey } = generateKeyPairSync('ed25519');
    mkdirSync(dirname(KEY_PATH), { recursive: true });
    writeFileSync(KEY_PATH, privateKey.export({ format: 'pem', type: 'pkcs8' }), {
      mode: 0o600,
      flag: 'wx',
    });
    console.log(`✔ Chave privada criada em ${KEY_PATH} (guarde uma cópia em lugar seguro).`);
  } else if (cmd !== 'pubkey') {
    console.error('Uso: node scripts/update-key.mjs keygen | pubkey');
    process.exit(1);
  }
  console.log(
    `Chave pública (desktop/internal/adapters/updater/pubkey.go): ${publicKeyBase64(readFileSync(KEY_PATH, 'utf8'))}`,
  );
}
