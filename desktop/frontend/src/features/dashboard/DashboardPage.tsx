import { useEffect, useState } from 'react';
import { Check, Circle, Loader2, Play, RotateCw, Square } from 'lucide-react';
import type { PageId } from '../../App';
import { useI18n } from '../../i18n';
import { backend, type Connect, type Status, type Step } from '../../shared/backend';
import { CopyButton, ErrorCard, PageHead } from '../../shared/components';
import { useAction } from '../../shared/hooks';
import { useLogs } from '../logs/LogDrawer';
import { MapChangeCard } from './MapChangeCard';
import { PlayersCard } from './PlayersCard';

// Ordem das fases para a linha do tempo.
const ORDER = ['starting', 'modio', 'mods', 'traveling', 'ready'];
const STEPS = [
  { id: 'start', at: 'starting' },
  { id: 'modio', at: 'modio' },
  { id: 'mods', at: 'mods' },
  { id: 'map', at: 'traveling' },
  { id: 'ready', at: 'ready' },
];

export function DashboardPage({ go, status }: { go: (p: PageId) => void; status: Status | null }) {
  const { t, date } = useI18n();
  const logs = useLogs();
  const action = useAction();
  const [connect, setConnect] = useState<Connect | null>(null);
  const [steps, setSteps] = useState<Step[]>([]);

  const phase = status?.phase ?? 'stopped';
  const running = (status?.pids?.length ?? 0) > 0;
  const busy = ['starting', 'modio', 'mods', 'traveling', 'stopping'].includes(phase);

  useEffect(() => {
    backend
      .connect()
      .then(setConnect)
      .catch(() => {});
    backend
      .checklist()
      .then(setSteps)
      .catch(() => {});
  }, [phase]);

  const idx = ORDER.indexOf(phase);
  const required = steps.filter((s) => !s.optional);
  const pending = steps.filter((s) => !s.done);
  const activeMods = status?.mods?.filter((m) => m.state === 'active').length ?? 0;
  const problem = status?.problem ?? null;

  return (
    <div className="page">
      <PageHead title={t('dashboard.title')} subtitle={t('dashboard.subtitle')} />

      {action.error && (
        <ErrorCard
          error={action.error}
          onClose={() => action.setError(null)}
          onOpenLogs={logs.showErrors}
        />
      )}
      {problem && (phase === 'crashed' || !running) && (
        <ErrorCard error={problem} onOpenLogs={logs.showErrors} />
      )}

      <div className="grid grid--2">
        <section className="card status">
          <div className={`status__phase status__phase--${phase}`}>
            {busy ? (
              <Loader2 className="spin" size={22} aria-hidden />
            ) : (
              <span className={`dot dot--lg dot--${phase}`} aria-hidden />
            )}
            <div>
              <h2>{t(`phase.${phase}`)}</h2>
              <p className="muted">{t(`phaseHint.${phase}`)}</p>
            </div>
          </div>

          <div className="status__actions">
            {!running ? (
              <button
                type="button"
                className="btn btn--primary btn--large"
                disabled={action.busy || phase === 'not-installed'}
                onClick={() => action.run(backend.start)}
              >
                <Play size={18} aria-hidden />
                {t('dashboard.start')}
              </button>
            ) : (
              <>
                <button
                  type="button"
                  className="btn btn--danger btn--large"
                  disabled={action.busy || phase === 'stopping'}
                  onClick={() => action.run(backend.stop)}
                >
                  <Square size={16} aria-hidden />
                  {t('dashboard.stop')}
                </button>
                <button
                  type="button"
                  className="btn btn--large"
                  disabled={action.busy || phase === 'stopping'}
                  onClick={() => action.run(backend.restart)}
                >
                  <RotateCw size={16} aria-hidden />
                  {t('dashboard.restart')}
                </button>
              </>
            )}
          </div>

          {running && (
            <>
              <h3 className="card__title">{t('dashboard.timeline')}</h3>
              <ol className="timeline">
                {STEPS.map((s) => {
                  const at = ORDER.indexOf(s.at);
                  const state =
                    idx > at || phase === 'ready' ? 'done' : idx === at ? 'now' : 'todo';
                  return (
                    <li key={s.id} className={`timeline__step is-${state}`}>
                      {state === 'done' ? (
                        <Check size={14} />
                      ) : state === 'now' ? (
                        <Loader2 size={14} className="spin" />
                      ) : (
                        <Circle size={10} />
                      )}
                      {t(`dashboard.step.${s.id}`)}
                    </li>
                  );
                })}
              </ol>
              <dl className="facts">
                {status?.map && (
                  <>
                    <dt>{t('dashboard.map')}</dt>
                    <dd>{status.map}</dd>
                  </>
                )}
                {status?.scenario && (
                  <>
                    <dt>{t('dashboard.scenario')}</dt>
                    <dd>{status.scenario}</dd>
                  </>
                )}
                {!!status?.port && (
                  <>
                    <dt>{t('dashboard.port')}</dt>
                    <dd>{status.port}</dd>
                  </>
                )}
                {status?.since && (
                  <>
                    <dt>{t('dashboard.since')}</dt>
                    <dd>{date(status.since as unknown as string)}</dd>
                  </>
                )}
                {activeMods > 0 && (
                  <>
                    <dt>{t('nav.mods')}</dt>
                    <dd>
                      <button type="button" className="link" onClick={() => go('mods')}>
                        {t('dashboard.modsActive', { n: activeMods })}
                      </button>
                    </dd>
                  </>
                )}
              </dl>
              {status?.external && <p className="banner banner--info">{t('dashboard.external')}</p>}
              {status?.authorizing && (
                <p className="banner banner--info">{t('dashboard.authorizing')}</p>
              )}
              {status?.remount && <p className="banner banner--info">{t('dashboard.remount')}</p>}
              {status?.gameStatsMissing && (
                <p className="hint">{t('dashboard.gameStatsMissing')}</p>
              )}
              {problem && running && <ErrorCard error={problem} onOpenLogs={logs.showErrors} />}
              {status?.cmdLine && (
                <details className="cmdline">
                  <summary>{t('dashboard.cmdLine')}</summary>
                  <code>{status.cmdLine}</code>
                </details>
              )}
            </>
          )}
        </section>

        <div className="stack">
          {running && <MapChangeCard status={status} go={go} />}
          {running && <PlayersCard status={status} />}
          <section className="card">
            <h2 className="card__title">{t('dashboard.connect.title')}</h2>
            <p className="muted">{t('dashboard.connect.text')}</p>
            {connect?.command ? (
              <>
                <p>{t('dashboard.connect.console')}</p>
                <div className="copyline">
                  <code>{connect.command}</code>
                  <CopyButton text={connect.command} />
                </div>
              </>
            ) : (
              <p className="hint">
                <button type="button" className="link" onClick={() => go('server')}>
                  {t('dashboard.connect.noIp')}
                </button>
              </p>
            )}
          </section>

          <section className="card">
            <h2 className="card__title">{t('dashboard.checklist.title')}</h2>
            {required.length > 0 && pending.filter((s) => !s.optional).length === 0 && (
              <p className="banner banner--ok">{t('dashboard.checklist.allDone')}</p>
            )}
            <ul className="checklist">
              {steps.map((s) => (
                <li key={s.id} className={s.done ? 'is-done' : ''}>
                  {s.done ? (
                    <Check size={16} className="ok" aria-hidden />
                  ) : (
                    <Circle size={14} aria-hidden />
                  )}
                  <span>
                    {t(`dashboard.checklist.${s.id}`)}
                    {s.optional && <span className="badge">{t('common.optional')}</span>}
                  </span>
                  {!s.done && (
                    <button
                      type="button"
                      className="btn btn--small btn--ghost"
                      onClick={() => go(s.page as PageId)}
                    >
                      {t('common.go')}
                    </button>
                  )}
                </li>
              ))}
            </ul>
          </section>
        </div>
      </div>
    </div>
  );
}
