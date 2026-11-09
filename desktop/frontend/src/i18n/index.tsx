// Idiomas da tela (ADR-009). O inglês é a referência: os outros arquivos têm o mesmo formato (o TypeScript acusa
// chave faltando). Sem biblioteca externa: t('chave.com.pontos', { param }) troca {param} no texto.
import { createContext, useCallback, useContext, useMemo, useState, type ReactNode } from 'react';
import { enUS, type Messages } from './locales/en-US';
import { ptPT } from './locales/pt-PT';
import { esES } from './locales/es-ES';
import { ru } from './locales/ru';
import { uk } from './locales/uk';

export const LANGUAGES = [
  { id: 'en-US', label: 'English (US)' },
  { id: 'pt-PT', label: 'Português (Portugal)' },
  { id: 'es-ES', label: 'Español (España)' },
  { id: 'ru', label: 'Русский' },
  { id: 'uk', label: 'Українська' },
] as const;

export type Lang = (typeof LANGUAGES)[number]['id'];

export const MESSAGES: Record<Lang, Messages> = {
  'en-US': enUS,
  'pt-PT': ptPT,
  'es-ES': esES,
  ru,
  uk,
};

export type Params = Record<string, string | number>;

/** Busca um valor pelo caminho com pontos ("server.tabs.general"). */
export function lookup(messages: unknown, key: string): unknown {
  let cur: unknown = messages;
  for (const part of key.split('.')) {
    if (cur === null || typeof cur !== 'object') return undefined;
    cur = (cur as Record<string, unknown>)[part];
  }
  return cur;
}

export function format(text: string, params?: Params): string {
  if (!params) return text;
  return text.replace(/\{(\w+)\}/g, (m, name: string) =>
    params[name] !== undefined ? String(params[name]) : m,
  );
}

/** Tradução com volta para o inglês e, por fim, para a própria chave. */
export function translate(lang: Lang, key: string, params?: Params): string {
  const v = lookup(MESSAGES[lang], key) ?? lookup(enUS, key);
  return typeof v === 'string' ? format(v, params) : key;
}

/** Lista de textos (passos de um guia, soluções de um erro). */
export function translateList(lang: Lang, key: string, params?: Params): string[] {
  const v = lookup(MESSAGES[lang], key) ?? lookup(enUS, key);
  return Array.isArray(v) ? v.map((s) => format(String(s), params)) : [];
}

/** Idioma inicial: o do Windows, se for um dos cinco (pt-BR cai em pt-PT, es-* em es-ES). */
export function detectLang(navLang: string | undefined): Lang {
  const l = (navLang ?? '').toLowerCase();
  if (l.startsWith('pt')) return 'pt-PT';
  if (l.startsWith('es')) return 'es-ES';
  if (l.startsWith('ru')) return 'ru';
  if (l.startsWith('uk')) return 'uk';
  return 'en-US';
}

export function isLang(v: string | null | undefined): v is Lang {
  return LANGUAGES.some((l) => l.id === v);
}

interface I18n {
  lang: Lang;
  setLang: (l: Lang) => void;
  t: (key: string, params?: Params) => string;
  tl: (key: string, params?: Params) => string[];
  /** Data e hora no formato do idioma. */
  date: (d: Date | string) => string;
}

const Ctx = createContext<I18n | null>(null);

export function I18nProvider({ initial, children }: { initial: Lang; children: ReactNode }) {
  const [lang, setLangState] = useState<Lang>(initial);
  const setLang = useCallback((l: Lang) => {
    setLangState(l);
    document.documentElement.lang = l;
  }, []);
  const value = useMemo<I18n>(
    () => ({
      lang,
      setLang,
      t: (key, params) => translate(lang, key, params),
      tl: (key, params) => translateList(lang, key, params),
      date: (d) => {
        const date = typeof d === 'string' ? new Date(d) : d;
        if (Number.isNaN(date.getTime()) || date.getFullYear() < 2000) return '—';
        return date.toLocaleString(lang, { dateStyle: 'medium', timeStyle: 'short' });
      },
    }),
    [lang, setLang],
  );
  return <Ctx.Provider value={value}>{children}</Ctx.Provider>;
}

export function useI18n(): I18n {
  const v = useContext(Ctx);
  if (!v) throw new Error('useI18n fora do I18nProvider');
  return v;
}
