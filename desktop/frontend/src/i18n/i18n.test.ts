import { readFileSync } from 'node:fs';
import { resolve } from 'node:path';
import { detectLang, lookup, MESSAGES, translate, translateList, LANGUAGES } from './index';
import { enUS } from './locales/en-US';
import { bytes, duration } from '../shared/format';
import { toAppError } from '../shared/problem';

const GO = resolve(__dirname, '../../../internal');

/** Lê constantes de string de um arquivo Go (ex.: os códigos do catálogo). */
function goCodes(): string[] {
  const src = readFileSync(resolve(GO, 'domain/problems/problems.go'), 'utf8');
  const block = src.slice(src.indexOf('const ('), src.indexOf(')', src.indexOf('const (')));
  return [...block.matchAll(/=\s*"([a-z-]+)"/g)].map((m) => m[1]!);
}

/** Todas as chaves "validation.x" usadas no Go. */
function goValidationKeys(): string[] {
  const files = [
    'domain/server/profile.go',
    'domain/server/cmdline.go',
    'domain/server/scenarios.go',
    'domain/server/messages.go',
    'domain/server/rules.go',
    'domain/server/mutators.go',
    'domain/server/automation.go',
    'domain/server/markup.go',
    'domain/admin/admin.go',
    'app/config/profiles.go',
    'app/notify/notify.go',
    'adapters/desktop/automation.go',
    'domain/backup/backup.go',
    'adapters/desktop/config_api.go',
    'adapters/desktop/server_api.go',
    'app/runner/runner.go',
    'app/config/config.go',
    'app/config/logs.go',
  ];
  const keys = new Set<string>();
  for (const f of files) {
    const src = readFileSync(resolve(GO, f), 'utf8');
    for (const m of src.matchAll(/"(validation\.[A-Za-z]+)"/g)) keys.add(m[1]!);
  }
  return [...keys];
}

/** Todas as folhas (caminho com pontos) de um objeto de mensagens. */
function leaves(o: unknown, prefix = ''): string[] {
  if (typeof o === 'string' || Array.isArray(o)) return [prefix];
  return Object.entries(o as Record<string, unknown>).flatMap(([k, v]) => leaves(v, prefix ? `${prefix}.${k}` : k));
}

describe('i18n', () => {
  it('has text for every problem code of the Go catalog', () => {
    const codes = goCodes();
    expect(codes.length).toBeGreaterThan(20);
    for (const code of codes) {
      for (const part of ['title', 'what', 'cause']) {
        expect(typeof lookup(enUS, `problems.${code}.${part}`), `${code}.${part}`).toBe('string');
      }
      expect(Array.isArray(lookup(enUS, `problems.${code}.solution`)), `${code}.solution`).toBe(true);
    }
  });

  it('has text for every validation key used in Go', () => {
    for (const key of goValidationKeys()) {
      expect(typeof lookup(enUS, key), key).toBe('string');
    }
  });

  it('every locale has the same keys as en-US', () => {
    const ref = leaves(enUS).sort();
    for (const { id } of LANGUAGES) {
      expect(leaves(MESSAGES[id]).sort(), id).toEqual(ref);
    }
  });

  it('formats params and falls back', () => {
    expect(translate('en-US', 'app.version', { version: '1.2.3' })).toBe('Version 1.2.3');
    expect(translate('ru', 'no.such.key')).toBe('no.such.key');
    expect(translateList('en-US', 'problems.busy.solution').length).toBeGreaterThan(0);
  });

  it('detects the Windows language', () => {
    expect(detectLang('pt-BR')).toBe('pt-PT');
    expect(detectLang('es-MX')).toBe('es-ES');
    expect(detectLang('uk-UA')).toBe('uk');
    expect(detectLang('ru')).toBe('ru');
    expect(detectLang('de-DE')).toBe('en-US');
    expect(detectLang(undefined)).toBe('en-US');
  });
});

describe('format', () => {
  const t = (k: string, p?: Record<string, string | number>) => translate('en-US', k, p);
  it('bytes and duration', () => {
    expect(bytes(0, t)).toBe('0 B');
    expect(bytes(1536, t)).toBe('1.50 KB');
    expect(bytes(5.2 * 1024 ** 3, t)).toBe('5.20 GB');
    expect(duration(30, t)).toBe('30 s');
    expect(duration(600, t)).toBe('10 min');
    expect(duration(3900, t)).toBe('1 h 05 min');
  });
});

describe('toAppError', () => {
  it('parses the Go JSON error', () => {
    const e = toAppError('{"code":"invalid-field","field":"port","msg":"validation.port"}');
    expect(e).toEqual({ code: 'invalid-field', field: 'port', msg: 'validation.port' });
  });
  it('wraps unknown errors', () => {
    expect(toAppError(new Error('boom'))).toEqual({ code: 'unknown', detail: 'boom' });
  });
});
