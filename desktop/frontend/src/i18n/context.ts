// Contexto do idioma num módulo só dele: quando um arquivo de tradução muda no modo de desenvolvimento, o Vite
// recarrega o index.tsx, mas o contexto continua o mesmo objeto. Antes ele era recriado, as telas abertas deixavam
// de achar o provedor e a janela ficava preta.
import { createContext } from 'react';

export type Params = Record<string, string | number>;

export interface I18nValue<L extends string = string> {
  lang: L;
  setLang: (l: L) => void;
  t: (key: string, params?: Params) => string;
  tl: (key: string, params?: Params) => string[];
  /** Data e hora no formato do idioma. */
  date: (d: Date | string) => string;
}

export const I18nCtx = createContext<I18nValue | null>(null);
