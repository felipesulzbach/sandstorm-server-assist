// Checks the translated user docs (ADR-035): each translation of README.md and docs/GUIDE.md must have the same
// sections and images as the English file, a language bar with all 5 languages, and only local links that exist.
// Runs in "npm test".
import { existsSync, readFileSync } from 'node:fs';
import { dirname, join, resolve } from 'node:path';

const ROOT = resolve(import.meta.dirname, '..');
const LANGS = ['pt-PT', 'es-ES', 'ru', 'uk'];
const SETS = [
  { en: 'README.md', tr: (l) => `README.${l}.md` },
  { en: 'docs/GUIDE.md', tr: (l) => `docs/${l}/GUIDE.md` },
];

const errors = [];
const read = (f) => readFileSync(join(ROOT, f), 'utf8').replace(/\r/g, '');
const withoutCode = (md) => md.replace(/```[\s\S]*?```/g, '');

function structure(file) {
  const md = withoutCode(read(file));
  return {
    h2: (md.match(/^## /gm) ?? []).length,
    h3: (md.match(/^### /gm) ?? []).length,
    images: [...md.matchAll(/!\[[^\]]*\]\(([^)\s]+)\)/g)]
      .map((m) => m[1])
      .filter((u) => !u.startsWith('http'))
      .map((u) => u.split('/').pop()),
    numbered: (md.match(/^\d+\. /gm) ?? []).length,
  };
}

function checkLinks(file) {
  const md = withoutCode(read(file));
  for (const m of md.matchAll(/\]\(([^)\s]+)\)/g)) {
    const target = m[1].split('#')[0];
    if (!target || /^[a-z]+:/.test(target)) continue;
    if (!existsSync(join(ROOT, dirname(file), target)))
      errors.push(`${file}: link to missing file ${m[1]}`);
  }
}

function checkLanguageBar(file) {
  const bar = read(file)
    .split('\n')
    .find((l) => l.includes('English') && l.includes('Українська'));
  if (!bar) {
    errors.push(`${file}: no language bar`);
    return;
  }
  const links = [...bar.matchAll(/\]\(([^)]+)\)/g)].length;
  if (links !== 4)
    errors.push(`${file}: the language bar should link the 4 other languages (found ${links})`);
}

for (const set of SETS) {
  if (!existsSync(join(ROOT, set.en))) {
    errors.push(`missing ${set.en}`);
    continue;
  }
  const en = structure(set.en);
  checkLinks(set.en);
  checkLanguageBar(set.en);
  for (const l of LANGS) {
    const f = set.tr(l);
    if (!existsSync(join(ROOT, f))) {
      errors.push(`missing translation ${f}`);
      continue;
    }
    const tr = structure(f);
    for (const k of ['h2', 'h3', 'numbered']) {
      if (tr[k] !== en[k]) errors.push(`${f}: ${k} count ${tr[k]}, English has ${en[k]}`);
    }
    if (tr.images.join(',') !== en.images.join(','))
      errors.push(`${f}: images [${tr.images}] differ from English [${en.images}]`);
    checkLinks(f);
    checkLanguageBar(f);
  }
}

if (errors.length) {
  console.error(`✖ Docs translations:\n  ${errors.join('\n  ')}`);
  process.exit(1);
}
console.log(`✔ Docs translations: ${SETS.length} documents × ${LANGS.length + 1} languages match`);
