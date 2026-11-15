import { Search, X } from 'lucide-react';
import { useI18n } from '../i18n';

/** Texto para comparar na busca: sem maiúsculas, acentos nem sublinhados ("Dificuldade" acha "dificuldade"). */
export const fold = (s: string) =>
  s.normalize('NFD').replace(/\p{M}/gu, '').replace(/_/g, ' ').toLowerCase();

/** Diz se algum dos textos contém todas as palavras da busca (vazia acha tudo). */
export function matches(query: string, ...texts: string[]): boolean {
  const words = fold(query).split(/\s+/).filter(Boolean);
  if (words.length === 0) return true;
  const all = texts.map(fold).join(' \u0000 ');
  // "AIDifficulty" também se acha por "ai difficulty": compara junto e separado.
  const joined = all.replace(/\s+/g, '');
  return words.every((w) => all.includes(w) || joined.includes(w));
}

/** Campo de busca de uma lista, com botão de limpar. Só filtra o que a tela mostra. */
export function SearchBox({
  value,
  onChange,
  placeholder,
  label,
}: {
  value: string;
  onChange: (v: string) => void;
  placeholder: string;
  label: string;
}) {
  const { t } = useI18n();
  return (
    <div className="searchbox" role="search">
      <Search size={15} aria-hidden className="searchbox__icon" />
      <input
        className="input searchbox__input"
        type="search"
        value={value}
        aria-label={label}
        placeholder={placeholder}
        onChange={(e) => onChange(e.target.value)}
        onKeyDown={(e) => {
          if (e.key === 'Escape') onChange('');
        }}
      />
      {value && (
        <button
          type="button"
          className="icon-btn searchbox__clear"
          title={t('common.clearSearch')}
          aria-label={t('common.clearSearch')}
          onClick={() => onChange('')}
        >
          <X size={14} aria-hidden />
        </button>
      )}
    </div>
  );
}
