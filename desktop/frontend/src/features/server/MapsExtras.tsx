import { useEffect, useRef, useState } from 'react';
import { ListRestart } from 'lucide-react';
import { useI18n } from '../../i18n';
import {
  backend,
  type CycleDefault,
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
  // Mutators que o servidor recusou no mapa atual (ADR-041), lidos do estado do servidor.
  const [refused, setRefused] = useState<string[]>([]);

  useEffect(() => {
    backend
      .mutators()
      .then((l) => setList(l ?? []))
      .catch(() => {});
    backend
      .status()
      .then((s) => setRefused(s?.refusedMutators ?? []))
      .catch(() => {});
    backend
      .modMutators()
      .then((l) => setMods(l ?? []))
      .catch(() => setMods([]));
  }, []);

  const official = new Set(list.map((m) => m.id));
  const fromMods = new Set((mods ?? []).map((m) => m.id));
  // Marcados no perfil, mas sem mod instalado que traga (só depois de as listas chegarem).
  const missing =
    list.length > 0 && mods ? value.filter((m) => !official.has(m) && !fromMods.has(m)) : [];
  const toggle = (id: string, on: boolean) =>
    onChange(on ? [...value.filter((m) => m !== id), id] : value.filter((m) => m !== id));
  // Mutators do jogo marcados que trocam o mesmo ajuste que este (a lista de conflitos vem do Go, ADR-045).
  const labelOf = new Map(list.map((m) => [m.id, m.label]));
  const blockers = (m: Mutator) => (m.conflicts ?? []).filter((c) => value.includes(c));

  const item = (
    id: string,
    label: string,
    hint: string,
    kind: 'mod' | 'game' | 'missing',
    blockedBy: string[] = [],
  ) => {
    const checked = value.includes(id);
    const names = blockedBy.map((c) => labelOf.get(c) ?? c).join(', ');
    return (
      <label
        key={`${kind}-${id}`}
        className={`check mutators__item mutators__item--${kind}${!checked && blockedBy.length ? ' is-blocked' : ''}`}
      >
        <input
          type="checkbox"
          checked={checked}
          disabled={!checked && blockedBy.length > 0}
          onChange={(e) => toggle(id, e.target.checked)}
        />
        <span>
          {label}
          {kind !== 'game' && <span className="mutators__badge">{t('server.mutatorsBadge')}</span>}
          {refused.includes(id) && (
            <span className="pill pill--warn">{t('server.mutatorsRefusedBadge')}</span>
          )}
          {checked && blockedBy.length > 0 && (
            <span className="pill pill--warn">{t('server.mutatorsConflictBadge')}</span>
          )}
          <span className="field__hint">{hint}</span>
          {blockedBy.length > 0 && (
            <span className="field__hint mutators__conflict">
              {t(checked ? 'server.mutatorsConflictBoth' : 'server.mutatorsConflict', { names })}
            </span>
          )}
        </span>
      </label>
    );
  };

  return (
    <div className="mutators">
      {refused.length > 0 && (
        <p className="banner banner--warn">
          {t('server.mutatorsRefused', { names: refused.join(', ') })}
        </p>
      )}
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
        <p className="hint">{t('server.mutatorsConflictHint')}</p>
        <div className="mutators__grid">
          {list.map((m) => item(m.id, m.label, t(`mutators.${m.id}`), 'game', blockers(m)))}
        </div>
      </div>
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

/**
 * Rotação padrão com os mapas do jogo, no modo do mapa inicial (ADR-043). Numa rotação vazia de um servidor sem
 * MapCycle.txt (instalação nova), preenche sozinha uma vez; o botão troca a rotação a qualquer momento, com
 * confirmação se ela já tem linhas. Só muda o formulário: vale ao Guardar.
 */
export function DefaultCycle({
  lines,
  startScenario,
  onFill,
}: {
  lines: string[];
  startScenario: string;
  onFill: (lines: string[]) => void;
}) {
  const { t } = useI18n();
  const [def, setDef] = useState<CycleDefault | null>(null);
  const [auto, setAuto] = useState(false);
  const [ask, setAsk] = useState(false);
  const filled = useRef(false);
  const current = useRef({ lines, onFill });
  current.current = { lines, onFill };

  useEffect(() => {
    let alive = true;
    backend
      .defaultMapCycle()
      .then((d) => {
        if (!alive) return;
        setDef(d);
        const next = d.lines ?? [];
        if (
          !filled.current &&
          current.current.lines.length === 0 &&
          !d.serverHasFile &&
          next.length > 0
        ) {
          filled.current = true;
          setAuto(true);
          current.current.onFill(next);
        }
      })
      .catch(() => {});
    return () => {
      alive = false;
    };
  }, [startScenario]);

  const next = def?.lines ?? [];
  const maps = def?.maps || next.length;
  if (next.length === 0) return null;
  const fill = () => {
    onFill(next);
    setAsk(false);
    setAuto(false);
  };
  return (
    <div className="defaultcycle">
      {auto && <p className="hint">{t('server.cycleDefault.filled', { n: maps })}</p>}
      {ask ? (
        <div className="row row--wrap">
          <span>{t('server.cycleDefault.replace')}</span>
          <button type="button" className="btn btn--small btn--primary" onClick={fill}>
            {t('common.yes')}
          </button>
          <button type="button" className="btn btn--small btn--ghost" onClick={() => setAsk(false)}>
            {t('common.cancel')}
          </button>
        </div>
      ) : (
        <button
          type="button"
          className="btn btn--small"
          onClick={() => (lines.length > 0 ? setAsk(true) : fill())}
        >
          <ListRestart size={14} aria-hidden /> {t('server.cycleDefault.button', { n: maps })}
        </button>
      )}
    </div>
  );
}
