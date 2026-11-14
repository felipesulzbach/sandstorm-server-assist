// Idiomas da tela (ADR-009). O inglês é a referência: os outros arquivos têm o mesmo formato (o TypeScript acusa
// chave faltando). Sem biblioteca externa: t('chave.com.pontos', { param }) troca {param} no texto.
import { useCallback, useContext, useMemo, useState, type ReactNode } from 'react';
import { I18nCtx, type I18nValue, type Params } from './context';
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

export type { Params };

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

type I18n = I18nValue<Lang>;

/** Fora do provedor (não deveria acontecer): inglês, sem trocar idioma. A tela de erro nunca quebra por isso. */
const FALLBACK: I18n = {
  lang: 'en-US',
  setLang: () => {},
  t: (key, params) => translate('en-US', key, params),
  tl: (key, params) => translateList('en-US', key, params),
  date: (d) => String(d),
};

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
  return <I18nCtx.Provider value={value as I18nValue}>{children}</I18nCtx.Provider>;
}

export function useI18n(): I18n {
  return (useContext(I18nCtx) as I18n | null) ?? FALLBACK;
}
