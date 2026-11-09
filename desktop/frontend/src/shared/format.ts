import type { Params } from '../i18n';

type T = (key: string, params?: Params) => string;

/** 1536 → "1.5 KB" (unidades traduzidas). */
export function bytes(n: number, t: T): string {
  if (!Number.isFinite(n) || n <= 0) return `0 ${t('units.b')}`;
  const units = ['units.b', 'units.kb', 'units.mb', 'units.gb'];
  let i = 0;
  let v = n;
  while (v >= 1024 && i < units.length - 1) {
    v /= 1024;
    i++;
  }
  const digits = i === 0 ? 0 : v < 10 ? 2 : v < 100 ? 1 : 0;
  return `${v.toFixed(digits)} ${t(units[i]!)}`;
}

/** Segundos → "45 s", "12 min", "1 h 05 min". */
export function duration(sec: number, t: T): string {
  if (sec < 60) return t('units.seconds', { n: Math.max(1, Math.round(sec)) });
  if (sec < 3600) return t('units.minutes', { n: Math.round(sec / 60) });
  const h = Math.floor(sec / 3600);
  const m = Math.round((sec % 3600) / 60);
  return t('units.hours', { h, m: String(m).padStart(2, '0') });
}

export function percent(done: number, total: number): number {
  if (!total || total <= 0) return 0;
  return Math.min(100, Math.max(0, (done * 100) / total));
}
