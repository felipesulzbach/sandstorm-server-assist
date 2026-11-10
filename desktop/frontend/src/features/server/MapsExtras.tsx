import { useEffect, useState } from 'react';
import { Moon, Sun } from 'lucide-react';
import { useI18n } from '../../i18n';
import { backend, type CycleLine, type Lighting, type Mutator } from '../../shared/backend';

/**
 * Mutators (ADR-022): os oficiais do jogo com a descrição, mais os de mods digitados pelo nome. A tela só monta a
 * lista de nomes; o Go valida.
 */
export function MutatorPicker({
  value,
  onChange,
}: {
  value: string[];
  onChange: (v: string[]) => void;
}) {
  const { t } = useI18n();
  const [list, setList] = useState<Mutator[]>([]);

  useEffect(() => {
    backend
      .mutators()
      .then(setList)
      .catch(() => {});
  }, []);

  const official = new Set(list.map((m) => m.id));
  const custom = value.filter((m) => !official.has(m));
  const toggle = (id: string, on: boolean) =>
    onChange(on ? [...value.filter((m) => m !== id), id] : value.filter((m) => m !== id));

  return (
    <div className="mutators">
      <div className="mutators__grid">
        {list.map((m) => (
          <label key={m.id} className="check mutators__item">
            <input
              type="checkbox"
              checked={value.includes(m.id)}
              onChange={(e) => toggle(m.id, e.target.checked)}
            />
            <span>
              {m.label}
              <span className="field__hint">{t(`mutators.${m.id}`)}</span>
            </span>
          </label>
        ))}
      </div>
      <label className="field">
        <span className="field__label">{t('server.mutatorsCustom')}</span>
        <input
          className="input"
          value={custom.join(', ')}
          placeholder={t('server.mutatorsCustomHint')}
          onChange={(e) =>
            onChange([
              ...value.filter((m) => official.has(m)),
              ...e.target.value
                .split(',')
                .map((s) => s.trim())
                .filter(Boolean),
            ])
          }
        />
      </label>
    </div>
  );
}

/** Dia ou noite em cada mapa da rotação (Lighting="Day|Night" na linha). */
export function CycleLighting({
  lines,
  onChange,
}: {
  lines: string[];
  onChange: (lines: string[]) => void;
}) {
  const { t } = useI18n();
  const [info, setInfo] = useState<CycleLine[]>([]);
  const key = lines.join('\n');

  useEffect(() => {
    backend
      .mapCycleLines(key ? key.split('\n') : [])
      .then(setInfo)
      .catch(() => {});
  }, [key]);

  const setLighting = async (index: number, lighting: Lighting) => {
    try {
      onChange(await backend.setMapCycleLighting(lines, index, lighting));
    } catch {
      /* linha que o app não entende: fica como está */
    }
  };

  if (lines.length === 0) return null;
  return (
    <div className="cycle">
      <div className="row row--wrap">
        <span className="field__label">{t('lighting.cycleTitle')}</span>
        <button
          type="button"
          className="btn btn--small btn--ghost"
          onClick={() => setLighting(-1, 'Day')}
        >
          <Sun size={14} aria-hidden /> {t('lighting.allDay')}
        </button>
        <button
          type="button"
          className="btn btn--small btn--ghost"
          onClick={() => setLighting(-1, 'Night')}
        >
          <Moon size={14} aria-hidden /> {t('lighting.allNight')}
        </button>
        <button
          type="button"
          className="btn btn--small btn--ghost"
          onClick={() => setLighting(-1, '')}
        >
          {t('lighting.allDefault')}
        </button>
      </div>
      <ul className="cycle__list">
        {info.map((l, i) => (
          <li key={`${i}-${l.scenario}`} className="cycle__row">
            <code>{l.scenario || lines[i]}</code>
            {l.scenario ? (
              <select
                className="select select--short"
                aria-label={t('lighting.label')}
                value={l.lighting}
                onChange={(e) => setLighting(i, e.target.value as Lighting)}
              >
                <option value="">{t('lighting.default')}</option>
                <option value="Day">{t('lighting.Day')}</option>
                <option value="Night">{t('lighting.Night')}</option>
              </select>
            ) : (
              <span className="muted">{t('lighting.unknownLine')}</span>
            )}
          </li>
        ))}
      </ul>
    </div>
  );
}
