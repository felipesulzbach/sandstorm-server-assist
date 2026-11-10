import { useEffect, useState } from 'react';
import { Moon, Sun } from 'lucide-react';
import { useI18n } from '../../i18n';
import {
  backend,
  type CycleLine,
  type Lighting,
  type ModMutator,
  type ModScenario,
  type Mutator,
} from '../../shared/backend';

/**
 * Mutators (ADR-022, ADR-023): os dos mods instalados (o Go lê os .pak) e os oficiais do jogo, numa lista com
 * rolagem. Mutator do perfil que não está em nenhum mod instalado continua marcado, com aviso. A tela só monta a
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
  const [mods, setMods] = useState<ModMutator[] | null>(null);

  useEffect(() => {
    backend
      .mutators()
      .then(setList)
      .catch(() => {});
    backend
      .modMutators()
      .then(setMods)
      .catch(() => setMods([]));
  }, []);

  const official = new Set(list.map((m) => m.id));
  const fromMods = new Set((mods ?? []).map((m) => m.id));
  // Marcados no perfil, mas sem mod instalado que traga (só depois de as listas chegarem).
  const missing =
    list.length > 0 && mods ? value.filter((m) => !official.has(m) && !fromMods.has(m)) : [];
  const toggle = (id: string, on: boolean) =>
    onChange(on ? [...value.filter((m) => m !== id), id] : value.filter((m) => m !== id));

  const item = (id: string, label: string, hint: string, kind: 'mod' | 'game' | 'missing') => (
    <label key={`${kind}-${id}`} className={`check mutators__item mutators__item--${kind}`}>
      <input
        type="checkbox"
        checked={value.includes(id)}
        onChange={(e) => toggle(id, e.target.checked)}
      />
      <span>
        {label}
        {kind !== 'game' && <span className="mutators__badge">{t('server.mutatorsBadge')}</span>}
        <span className="field__hint">{hint}</span>
      </span>
    </label>
  );

  return (
    <div className="mutators">
      <div className="mutators__scroll">
        <p className="mutators__group">{t('server.mutatorsMods')}</p>
        <div className="mutators__grid">
          {(mods ?? []).map((m) =>
            item(m.id, m.id, t('server.mutatorsFromMod', { name: m.modName }), 'mod'),
          )}
          {missing.map((id) => item(id, id, t('server.mutatorsMissing'), 'missing'))}
        </div>
        {mods && mods.length === 0 && missing.length === 0 && (
          <p className="hint">{t('server.mutatorsNoMods')}</p>
        )}
        <p className="mutators__group">{t('server.mutatorsGame')}</p>
        <div className="mutators__grid">
          {list.map((m) => item(m.id, m.label, t(`mutators.${m.id}`), 'game'))}
        </div>
      </div>
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

/** "Checkpoint_Security_Hardcore" → "Checkpoint Security Hardcore" (nomes dos mods, sem tradução). */
const humanMode = (s: string) => s.replace(/_/g, ' ');

/** Cenário preferido de um mapa: Checkpoint Security, se houver; senão o primeiro. */
function preferred(list: ModScenario[]): ModScenario | undefined {
  return list.find((s) => /_Checkpoint_Security$/.test(s.scenario)) ?? list[0];
}

/**
 * Mapa de mod depois de carregar os mods (-ModDownloadTravelTo, ADR-024): mapa e cenário escolhidos de listas
 * montadas pelo Go a partir dos .pak dos mods instalados, sem digitar. Um valor do perfil que não está nos mods
 * instalados continua na lista, marcado.
 */
export function ModMapPicker({
  map,
  scenario,
  disabled,
  onChange,
  mapError,
  scenarioError,
}: {
  map: string;
  scenario: string;
  disabled: boolean;
  onChange: (map: string, scenario: string) => void;
  mapError: string | null;
  scenarioError: string | null;
}) {
  const { t } = useI18n();
  const [list, setList] = useState<ModScenario[] | null>(null);

  useEffect(() => {
    backend
      .modScenarios()
      .then((l) => setList(l ?? []))
      .catch(() => setList([]));
  }, []);

  const on = map !== '' || scenario !== '';
  const maps = [...new Map((list ?? []).map((s) => [s.map, s.modName])).entries()];
  const ofMap = (list ?? []).filter((s) => s.map === map);
  const mapKnown = maps.some(([m]) => m === map);
  const scenarioKnown = ofMap.some((s) => s.scenario === scenario);
  const modeOf = (sc: string) => humanMode(sc.replace(`Scenario_${map}_`, ''));

  return (
    <>
      <label className="check check--block field--wide">
        <input
          type="checkbox"
          checked={on}
          disabled={disabled || (!on && (list ?? []).length === 0)}
          onChange={(e) => {
            const first = preferred(list ?? []);
            if (e.target.checked && first) {
              const pick = preferred((list ?? []).filter((s) => s.map === first.map)) ?? first;
              onChange(pick.map, pick.scenario);
            } else onChange('', '');
          }}
        />
        <span>
          {t('server.travel')}
          <span className="field__hint">
            {list && list.length === 0 && !on ? t('server.travelNoMods') : t('server.travelHint')}
          </span>
        </span>
      </label>
      {on && (
        <>
          <label className={`field ${mapError ? 'has-error' : ''}`}>
            <span className="field__label">{t('server.travelMap')}</span>
            <select
              className="select"
              value={map}
              onChange={(e) => {
                const pick = preferred((list ?? []).filter((s) => s.map === e.target.value));
                if (pick) onChange(pick.map, pick.scenario);
              }}
            >
              {!mapKnown && (
                <option value={map}>{t('server.travelNotInstalled', { name: map })}</option>
              )}
              {maps.map(([m, mod]) => (
                <option key={m} value={m}>
                  {m === mod ? m : `${m} · ${mod}`}
                </option>
              ))}
            </select>
            {mapError && <span className="field__error">{mapError}</span>}
          </label>
          <label className={`field ${scenarioError ? 'has-error' : ''}`}>
            <span className="field__label">{t('server.travelScenario')}</span>
            <select
              className="select"
              value={scenario}
              title={scenario}
              onChange={(e) => onChange(map, e.target.value)}
            >
              {!scenarioKnown && (
                <option value={scenario}>
                  {t('server.travelNotInstalled', { name: scenario })}
                </option>
              )}
              {ofMap.map((s) => (
                <option key={s.scenario} value={s.scenario}>
                  {modeOf(s.scenario)}
                </option>
              ))}
            </select>
            {scenarioError ? (
              <span className="field__error">{scenarioError}</span>
            ) : (
              <span className="field__hint">{scenario}</span>
            )}
          </label>
        </>
      )}
    </>
  );
}
