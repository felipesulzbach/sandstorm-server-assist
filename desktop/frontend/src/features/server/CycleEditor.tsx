import { useEffect, useMemo, useState } from 'react';
import { AlertTriangle, ArrowDown, ArrowUp, Moon, Plus, Sun, Trash2, X } from 'lucide-react';
import { useI18n } from '../../i18n';
import { backend, type CycleMap, type CycleView, type Lighting } from '../../shared/backend';

/** "Checkpoint_Security" → "Checkpoint Security" (nomes do jogo e dos mods, sem tradução). */
const human = (s: string) => s.replace(/_/g, ' ');

/** Chave de um mapa do catálogo: o mesmo mapa pode vir do jogo e de um mod. */
const mapKey = (m: CycleMap) => `${m.modId}:${m.map}`;

/** Iluminações que a rotação usa: toda linha é de dia ou de noite (ADR-046). */
const LIGHTS: Lighting[] = ['Day', 'Night'];

/** Chave de um cenário com a iluminação. */
const litKey = (scenario: string, lighting: Lighting) => `${scenario.toLowerCase()}|${lighting}`;

/** Quantos cenários do mapa já estão na rotação. */
const addedCount = (m: CycleMap, inCycle: Set<string>) =>
  m.scenarios.filter((s) => inCycle.has(s.id.toLowerCase())).length;

/**
 * Rotação de mapas por mapa (ADR-046): a tela mostra as linhas do MapCycle em blocos por mapa, com dia ou noite em
 * cada cenário, e muda a lista pelo Go (adicionar, tirar, subir e descer um mapa). O arquivo continua uma lista
 * simples, na ordem dos blocos.
 */
export function CycleEditor({
  lines,
  onChange,
}: {
  lines: string[];
  onChange: (lines: string[]) => void;
}) {
  const { t } = useI18n();
  const [view, setView] = useState<CycleView | null>(null);
  const [adding, setAdding] = useState(false);
  const key = lines.join('\n');

  useEffect(() => {
    backend
      .mapCycleView(key ? key.split('\n') : [])
      .then(setView)
      .catch(() => {});
  }, [key]);

  const groups = useMemo(() => view?.groups ?? [], [view]);
  const maps = view?.maps ?? [];
  const entries = useMemo(
    () => groups.flatMap((g) => g.entries).filter((e) => e.scenario),
    [groups],
  );
  // Cenários na rotação (com qualquer iluminação) e pares cenário + iluminação: o mesmo cenário pode estar de dia e
  // de noite (ADR-046).
  const inCycle = useMemo(() => new Set(entries.map((e) => e.scenario.toLowerCase())), [entries]);
  const inCycleLit = useMemo(
    () => new Set(entries.map((e) => litKey(e.scenario, e.lighting as Lighting))),
    [entries],
  );
  const missing = groups.flatMap((g) => g.entries).filter((e) => e.missing);
  // Todos os mapas disponíveis já estão na rotação: "Adicionar mapa" não tem o que oferecer.
  const allIn = maps.length > 0 && maps.every((m) => addedCount(m, inCycle) > 0);

  // Toda mudança passa pelo Go; uma recusa (lista longa demais, linha estranha) deixa a rotação como está.
  const run = async (p: Promise<string[]>) => {
    try {
      onChange(await p);
    } catch {
      /* a rotação fica como está */
    }
  };
  const setLighting = (index: number, lighting: Lighting) =>
    run(backend.setMapCycleLighting(lines, index, lighting));
  const remove = (indexes: number[]) => run(backend.removeFromMapCycle(lines, indexes));
  const move = (index: number, up: boolean) => run(backend.moveMapCycleGroup(lines, index, up));
  const add = (after: number, ids: string[], lightings: Lighting[]) =>
    run(backend.addToMapCycle(lines, after, ids, lightings));

  /** Cenários do mesmo mapa (jogo ou mods), de dia e de noite, que ainda não estão na rotação. */
  const moreFor = (map: string) =>
    maps
      .filter((m) => m.map.toLowerCase() === map.toLowerCase())
      .flatMap((m) => m.scenarios)
      .flatMap((s) => LIGHTS.map((light) => ({ ...s, light })))
      .filter((s) => !inCycleLit.has(litKey(s.id, s.light)));

  return (
    <div className="cycle">
      <p className="hint">{t('cycle.hint')}</p>
      {missing.length > 0 && (
        <div className="banner banner--warn row row--wrap" role="status">
          <AlertTriangle size={16} aria-hidden />
          <span>{t('lighting.cycleMissing')}</span>
          <button
            type="button"
            className="btn btn--small"
            onClick={() => remove(missing.map((e) => e.index))}
          >
            {t('lighting.cycleRemoveMissing')}
          </button>
        </div>
      )}
      {lines.length > 0 && (
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
        </div>
      )}

      {lines.length === 0 && <p className="muted">{t('cycle.empty')}</p>}
      <ol className="cyclemaps">
        {groups.map((g, gi) => {
          const first = g.entries[0]!.index;
          const last = g.entries[g.entries.length - 1]!.index;
          const name = g.map ? g.label || human(g.map) : t('cycle.unknownGroup');
          const more = g.map ? moreFor(g.map) : [];
          return (
            <li key={`${first}-${g.map}`} className="cyclemap">
              <div className="cyclemap__head">
                <div className="cyclemap__name">
                  <strong>{name}</strong>
                  {g.modName && (
                    <span className="muted">
                      <span className="mutators__badge">{t('server.mutatorsBadge')}</span>{' '}
                      {t('cycle.fromMod', { name: g.modName })}
                    </span>
                  )}
                </div>
                <div className="cyclemap__tools">
                  <button
                    type="button"
                    className="icon-btn"
                    title={t('cycle.moveUp', { map: name })}
                    aria-label={t('cycle.moveUp', { map: name })}
                    disabled={gi === 0}
                    onClick={() => move(first, true)}
                  >
                    <ArrowUp size={16} aria-hidden />
                  </button>
                  <button
                    type="button"
                    className="icon-btn"
                    title={t('cycle.moveDown', { map: name })}
                    aria-label={t('cycle.moveDown', { map: name })}
                    disabled={gi === groups.length - 1}
                    onClick={() => move(first, false)}
                  >
                    <ArrowDown size={16} aria-hidden />
                  </button>
                  <button
                    type="button"
                    className="icon-btn"
                    title={t('cycle.removeMap', { map: name })}
                    aria-label={t('cycle.removeMap', { map: name })}
                    onClick={() => remove(g.entries.map((e) => e.index))}
                  >
                    <Trash2 size={16} aria-hidden />
                  </button>
                </div>
              </div>
              <ul className="cycle__list">
                {g.entries.map((e) => (
                  <li key={e.index} className="cycle__row">
                    <code title={e.scenario || e.raw}>
                      {e.scenario ? human(e.mode || e.scenario) : e.raw}
                    </code>
                    {e.missing && (
                      <span className="pill pill--warn">{t('lighting.notInstalled')}</span>
                    )}
                    {e.scenario ? (
                      <select
                        className="select select--short"
                        aria-label={t('lighting.label')}
                        value={e.lighting}
                        onChange={(ev) => setLighting(e.index, ev.target.value as Lighting)}
                      >
                        <option value="Day">{t('lighting.Day')}</option>
                        <option value="Night">{t('lighting.Night')}</option>
                      </select>
                    ) : (
                      <span className="muted">{t('lighting.unknownLine')}</span>
                    )}
                    <button
                      type="button"
                      className="icon-btn"
                      title={t('cycle.removeLine', { scenario: e.scenario || e.raw })}
                      aria-label={t('cycle.removeLine', { scenario: e.scenario || e.raw })}
                      onClick={() => remove([e.index])}
                    >
                      <X size={16} aria-hidden />
                    </button>
                  </li>
                ))}
              </ul>
              {more.length > 0 && (
                <select
                  className="select cyclemap__more"
                  aria-label={t('cycle.addScenarioTo', { map: name })}
                  value=""
                  onChange={(ev) => {
                    const [id, light] = ev.target.value.split('|');
                    if (id && light) void add(last, [id], [light as Lighting]);
                  }}
                >
                  <option value="">{t('cycle.addScenario')}</option>
                  {more.map((s) => (
                    <option key={`${s.id}|${s.light}`} value={`${s.id}|${s.light}`}>
                      {`${human(s.mode || s.id)} · ${t(`lighting.${s.light}`)}`}
                    </option>
                  ))}
                </select>
              )}
            </li>
          );
        })}
      </ol>

      {adding ? (
        <AddMap
          maps={maps}
          inCycle={inCycle}
          onAdd={(ids, lightings) => {
            void add(-1, ids, lightings);
            setAdding(false);
          }}
          onCancel={() => setAdding(false)}
        />
      ) : (
        <div className="row row--wrap">
          <button type="button" className="btn" disabled={allIn} onClick={() => setAdding(true)}>
            <Plus size={16} aria-hidden /> {t('cycle.addMap')}
          </button>
          {allIn && <span className="hint">{t('cycle.allMapsIn')}</span>}
        </div>
      )}
    </div>
  );
}

/** Escolha de um mapa (do jogo ou de um mod instalado) e dos cenários dele para pôr no fim da rotação. */
function AddMap({
  maps,
  inCycle,
  onAdd,
  onCancel,
}: {
  maps: CycleMap[];
  inCycle: Set<string>;
  onAdd: (ids: string[], lightings: Lighting[]) => void;
  onCancel: () => void;
}) {
  const { t } = useI18n();
  const [chosen, setChosen] = useState('');
  const [picked, setPicked] = useState<string[]>([]);
  // Cada cenário marcado entra uma vez em cada iluminação marcada (Dia, Noite ou as duas).
  const [lights, setLights] = useState<Lighting[]>(['Day', 'Night']);
  const map = maps.find((m) => mapKey(m) === chosen);
  const fromMods = maps.filter((m) => m.modId !== 0);
  const fromGame = maps.filter((m) => m.modId === 0);

  const choose = (k: string) => {
    setChosen(k);
    const m = maps.find((x) => mapKey(x) === k);
    const free = (m?.scenarios ?? []).filter((s) => !inCycle.has(s.id.toLowerCase()));
    // Sugestão: os dois lados do Checkpoint, se o mapa tiver; senão o primeiro cenário.
    const checkpoint = free.filter((s) => /^Checkpoint_(Security|Insurgents)$/.test(s.mode));
    setPicked((checkpoint.length ? checkpoint : free.slice(0, 1)).map((s) => s.id));
  };
  // Mapa que já está na rotação (com um ou mais cenários): ✓, cor e sem poder escolher. Mais cenários dele entram
  // pelo bloco do mapa ("Adicionar um cenário deste mapa").
  const label = (m: CycleMap) => {
    const name =
      m.modName && m.modName.toLowerCase() !== m.label.toLowerCase()
        ? `${m.label} · ${m.modName}`
        : m.label;
    const n = addedCount(m, inCycle);
    return n === 0 ? name : `✓ ${name} — ${t('cycle.mapIn', { n, total: m.scenarios.length })}`;
  };
  const inRotation = (m: CycleMap) => addedCount(m, inCycle) > 0;

  return (
    <div className="cycleadd">
      <strong>{t('cycle.addMapTitle')}</strong>
      {maps.length === 0 ? (
        <p className="hint">{t('cycle.noMaps')}</p>
      ) : (
        <label className="field">
          <span className="field__label">{t('cycle.pickMap')}</span>
          <select className="select" value={chosen} onChange={(e) => choose(e.target.value)}>
            <option value="">{t('common.choose')}</option>
            {fromMods.length > 0 && (
              <optgroup label={t('cycle.groupMods')}>
                {fromMods.map((m) => (
                  <option
                    key={mapKey(m)}
                    value={mapKey(m)}
                    className={inRotation(m) ? 'cycleadd__in' : undefined}
                    disabled={inRotation(m)}
                  >
                    {label(m)}
                  </option>
                ))}
              </optgroup>
            )}
            <optgroup label={t('cycle.groupGame')}>
              {fromGame.map((m) => (
                <option
                  key={mapKey(m)}
                  value={mapKey(m)}
                  className={inRotation(m) ? 'cycleadd__in' : undefined}
                  disabled={inRotation(m)}
                >
                  {label(m)}
                </option>
              ))}
            </optgroup>
          </select>
        </label>
      )}
      {maps.length > 0 && <span className="field__hint">{t('cycle.addMapHint')}</span>}
      {map && (
        <fieldset className="cycleadd__scenarios">
          <legend className="field__label">{t('cycle.pickScenarios')}</legend>
          {map.scenarios.map((s) => {
            const already = inCycle.has(s.id.toLowerCase());
            return (
              <label key={s.id} className="check">
                <input
                  type="checkbox"
                  checked={already || picked.includes(s.id)}
                  disabled={already}
                  onChange={(e) =>
                    setPicked((p) =>
                      e.target.checked ? [...p, s.id] : p.filter((x) => x !== s.id),
                    )
                  }
                />
                <span>
                  {human(s.mode || s.id)}
                  {already && <span className="field__hint">{t('cycle.inCycle')}</span>}
                </span>
              </label>
            );
          })}
        </fieldset>
      )}
      {map && (
        <fieldset className="cycleadd__lights">
          <legend className="field__label">{t('lighting.label')}</legend>
          {LIGHTS.map((light) => (
            <label key={light} className="check">
              <input
                type="checkbox"
                checked={lights.includes(light)}
                onChange={(e) =>
                  setLights((l) =>
                    e.target.checked
                      ? LIGHTS.filter((x) => x === light || l.includes(x))
                      : l.filter((x) => x !== light),
                  )
                }
              />
              <span>{t(`lighting.${light}`)}</span>
            </label>
          ))}
          <span className="field__hint">{t('cycle.lightsHint')}</span>
        </fieldset>
      )}
      <div className="row row--wrap">
        <button
          type="button"
          className="btn btn--primary"
          disabled={picked.length === 0 || lights.length === 0}
          onClick={() => onAdd(picked, lights)}
        >
          <Plus size={16} aria-hidden /> {t('cycle.addSelected')}
        </button>
        <button type="button" className="btn btn--ghost" onClick={onCancel}>
          {t('common.cancel')}
        </button>
      </div>
    </div>
  );
}
