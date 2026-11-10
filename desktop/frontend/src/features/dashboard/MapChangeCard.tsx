import { useEffect, useMemo, useState } from 'react';
import { Map as MapIcon, RefreshCw } from 'lucide-react';
import type { PageId } from '../../App';
import { useI18n } from '../../i18n';
import { backend, type Lighting, type Scenario, type Status } from '../../shared/backend';
import { ErrorCard } from '../../shared/components';
import { useAction } from '../../shared/hooks';
import { useLogs } from '../logs/LogDrawer';

/** "LastLight" → "Last Light", "Old_World" → "Old World" (nomes do jogo e dos mods, sem tradução). */
const human = (s: string) => s.replace(/_/g, ' ').replace(/([a-z])([A-Z])/g, '$1 $2');

/**
 * Troca de mapa com o servidor ligado (ADR-019). O Go pede a lista ao servidor pelo RCON (uma conexão por clique,
 * nunca em intervalos) e troca com "travelscenario". A tela manda só o id do cenário escolhido na lista.
 */
export function MapChangeCard({ status, go }: { status: Status | null; go: (p: PageId) => void }) {
  const { t } = useI18n();
  const logs = useLogs();
  const action = useAction();
  const [rcon, setRcon] = useState<boolean | null>(null);
  const [list, setList] = useState<Scenario[] | null>(null);
  const [choice, setChoice] = useState('');
  const [lighting, setLighting] = useState<Lighting>('');
  const [confirming, setConfirming] = useState(false);
  const [sentTo, setSentTo] = useState<Scenario | null>(null);

  const phase = status?.phase ?? 'stopped';
  const ready = phase === 'ready';

  useEffect(() => {
    backend
      .mapChange()
      .then((i) => setRcon(i.rcon))
      .catch(() => setRcon(false));
  }, [phase]);

  // Troca concluída: o servidor voltou a "pronto" no cenário pedido.
  useEffect(() => {
    if (sentTo && ready && status?.scenario === sentTo.id) setSentTo(null);
  }, [sentTo, ready, status?.scenario]);

  const groups = useMemo(() => {
    const g = new Map<string, Scenario[]>();
    for (const s of list ?? []) g.set(s.map, [...(g.get(s.map) ?? []), s]);
    return [...g.entries()];
  }, [list]);

  const load = async () => {
    const got = await action.run(backend.scenarios);
    if (got === undefined) return;
    const l = got ?? [];
    setList(l);
    const cur = l.find((s) => s.id === status?.scenario);
    setChoice((c) => (l.some((s) => s.id === c) ? c : ((cur ?? l[0])?.id ?? '')));
  };

  const chosen = list?.find((s) => s.id === choice) ?? null;
  const label = (s: Scenario) => (s.mode ? `${human(s.map)} · ${human(s.mode)}` : human(s.map));

  const change = async () => {
    if (!chosen) return;
    const target = chosen;
    setConfirming(false);
    const ok = await action.run(async () => {
      await backend.changeMap(target.id, lighting);
      return true;
    });
    if (ok) setSentTo(target);
  };

  if (rcon === null) return null;

  return (
    <section className="card">
      <h2 className="card__title">{t('dashboard.mapChange.title')}</h2>

      {!rcon ? (
        <>
          <p className="muted">{t('dashboard.mapChange.needsRcon')}</p>
          <button type="button" className="btn btn--small" onClick={() => go('server')}>
            {t('dashboard.mapChange.openServer')}
          </button>
        </>
      ) : (
        <>
          <p className="muted">{t('dashboard.mapChange.text')}</p>

          {action.error && (
            <ErrorCard
              error={action.error}
              onClose={() => action.setError(null)}
              onOpenLogs={logs.showErrors}
            />
          )}

          {sentTo && (
            <p className="banner banner--info">
              {t('dashboard.mapChange.sent', { map: label(sentTo) })}
            </p>
          )}

          {!list ? (
            <button type="button" className="btn" disabled={!ready || action.busy} onClick={load}>
              <MapIcon size={16} aria-hidden /> {t('dashboard.mapChange.load')}
            </button>
          ) : list.length === 0 ? (
            <p className="hint">{t('dashboard.mapChange.empty')}</p>
          ) : (
            <>
              <div className="row row--wrap mapchange__pick">
                <select
                  className="select"
                  aria-label={t('dashboard.mapChange.title')}
                  value={choice}
                  disabled={action.busy || confirming}
                  onChange={(e) => setChoice(e.target.value)}
                >
                  {groups.map(([map, items]) => (
                    <optgroup key={map} label={human(map)}>
                      {items.map((s) => (
                        <option key={s.id} value={s.id}>
                          {(s.mode ? human(s.mode) : human(s.map)) +
                            (s.source === 'config'
                              ? ` (${t('dashboard.mapChange.fromConfig')})`
                              : '') +
                            (s.id === status?.scenario
                              ? ` · ${t('dashboard.mapChange.current')}`
                              : '')}
                        </option>
                      ))}
                    </optgroup>
                  ))}
                </select>
                <select
                  className="select select--short"
                  aria-label={t('lighting.label')}
                  value={lighting}
                  disabled={action.busy || confirming}
                  onChange={(e) => setLighting(e.target.value as Lighting)}
                >
                  <option value="">{t('lighting.default')}</option>
                  <option value="Day">{t('lighting.Day')}</option>
                  <option value="Night">{t('lighting.Night')}</option>
                </select>
                <button
                  type="button"
                  className="icon-btn"
                  title={t('dashboard.mapChange.reload')}
                  aria-label={t('dashboard.mapChange.reload')}
                  disabled={!ready || action.busy}
                  onClick={load}
                >
                  <RefreshCw size={16} />
                </button>
              </div>

              {/* O servidor só lista o jogo base: este aviso explica de onde vêm os mapas de mod. */}
              <p className="hint">{t('dashboard.mapChange.fromConfigHint')}</p>

              <div className="row row--wrap mapchange__actions">
                {confirming && chosen ? (
                  <>
                    <span>{t('dashboard.mapChange.confirm', { map: label(chosen) })}</span>
                    <button
                      type="button"
                      className="btn btn--primary btn--small"
                      disabled={action.busy || !ready}
                      onClick={change}
                    >
                      {t('dashboard.mapChange.confirmYes')}
                    </button>
                    <button
                      type="button"
                      className="btn btn--small"
                      onClick={() => setConfirming(false)}
                    >
                      {t('common.cancel')}
                    </button>
                  </>
                ) : (
                  <button
                    type="button"
                    className="btn btn--primary"
                    disabled={
                      !ready ||
                      action.busy ||
                      !chosen ||
                      (chosen.id === status?.scenario && lighting === '')
                    }
                    onClick={() => setConfirming(true)}
                  >
                    <MapIcon size={16} aria-hidden /> {t('dashboard.mapChange.change')}
                  </button>
                )}
              </div>
            </>
          )}

          {!ready && <p className="hint">{t('dashboard.mapChange.waitReady')}</p>}
          <p className="hint">{t('dashboard.mapChange.restartNote')}</p>
        </>
      )}
    </section>
  );
}
