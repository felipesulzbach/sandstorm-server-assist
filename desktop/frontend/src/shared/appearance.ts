// Aparência (ADR-011, mesmo modelo do game-in): <html data-theme> troca só as cores; <html data-skin> muda a forma
// do layout. Salva neste computador (localStorage); sem armazenamento, fica o padrão.
import { useEffect, useState } from 'react';

export type Theme = 'sand' | 'black' | 'olive' | 'blue' | 'light';
export type Skin = 'sidebar' | 'compact' | 'cards' | 'top';

export interface Appearance {
  theme: Theme;
  skin: Skin;
}

export const THEMES: { id: Theme; swatch: [string, string, string] }[] = [
  { id: 'sand', swatch: ['#0f0e0c', '#1d1c19', '#e0a43a'] },
  { id: 'black', swatch: ['#000000', '#111111', '#e0a43a'] },
  { id: 'olive', swatch: ['#11130c', '#1d2116', '#9cb35a'] },
  { id: 'blue', swatch: ['#0a1120', '#142037', '#3b82f6'] },
  { id: 'light', swatch: ['#f3f0ea', '#ffffff', '#b7791f'] },
];

export const SKINS: Skin[] = ['sidebar', 'compact', 'cards', 'top'];

const KEY = 'ssa.appearance';
const DEFAULT: Appearance = { theme: 'sand', skin: 'sidebar' };

export function loadAppearance(): Appearance {
  try {
    const saved = JSON.parse(localStorage.getItem(KEY) ?? '{}') as Partial<Appearance>;
    return {
      theme: THEMES.some((t) => t.id === saved.theme) ? saved.theme! : DEFAULT.theme,
      skin: SKINS.includes(saved.skin as Skin) ? saved.skin! : DEFAULT.skin,
    };
  } catch {
    return DEFAULT;
  }
}

export function applyAppearance(a: Appearance, doc: Document = document): void {
  doc.documentElement.dataset['theme'] = a.theme;
  doc.documentElement.dataset['skin'] = a.skin;
}

const listeners = new Set<(a: Appearance) => void>();

export function updateAppearance(patch: Partial<Appearance>): void {
  const a = { ...loadAppearance(), ...patch };
  try {
    localStorage.setItem(KEY, JSON.stringify(a));
  } catch {
    // sem armazenamento: vale só até fechar
  }
  applyAppearance(a);
  listeners.forEach((cb) => cb(a));
}

export function useAppearance(): Appearance {
  const [a, setA] = useState(loadAppearance);
  useEffect(() => {
    listeners.add(setA);
    return () => {
      listeners.delete(setA);
    };
  }, []);
  return a;
}
