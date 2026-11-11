import { useEffect, useRef, useState, type ReactNode } from 'react';
import { Download, FolderOpen, Import, RefreshCw, ShieldCheck, X } from 'lucide-react';
import { useI18n } from '../../i18n';
import {
  backend,
  type InstallInfo,
  type InstallState,
  type Profile,
  type Status,
} from '../../shared/backend';
import { ErrorCard, PageHead, ProgressBar } from '../../shared/components';
import { bytes, duration } from '../../shared/format';
import { useAction, useInstallState } from '../../shared/hooks';
import { useLogs } from '../logs/LogDrawer';
import { useExistingSync } from '../sync/ExistingSync';

/** Cartão de onde veio a última ação: o andamento, o resultado e o erro aparecem nele, perto do botão. */
type Origin = 'folders' | 'steamcmd' | 'server';

export function InstallPage({ status }: { status: Status | null }) {
  const { t, date } = useI18n();
  const logs = useLogs();
  const st = useInstallState();
  const action = useAction();
  const [profile, setProfile] = useState<Profile | null>(null);
  const [info, setInfo] = useState<InstallInfo | null>(null);
  const [origin, setOrigin] = useState<Origin | null>(null);
  const sync = useExistingSync((r) => {
    setProfile(r.profile);
    backend
      .installInfo()
      .then(setInfo)
      .catch(() => {});
  });

  const running = (status?.pids?.length ?? 0) > 0;
  const task = st?.task ?? '';
  const working = task !== '' || action.busy;

  const reload = () => {
    backend
      .profile()
      .then(setProfile)
      .catch(() => {});
    backend
      .installInfo()
      .then(setInfo)
      .catch(() => {});
  };
  useEffect(reload, [st?.finished, task]);

  const choose = async (which: 'server' | 'steamcmd') => {
    setOrigin('folders');
    const before = profile?.serverDir;
    const p = await action.run(
      which === 'server' ? backend.chooseServerDir : backend.chooseSteamCmdDir,
    );
    if (p) setProfile(p);
    const next = await backend.installInfo().catch(() => null);
    if (next) setInfo(next);
    // Pasta nova com um servidor já instalado: busca os dados dele para completar os campos.
    if (which === 'server' && p?.serverDir && p.serverDir !== before && next?.serverInstalled) {
      await sync.check('auto');
    }
  };

  const start = (which: 'steamcmd' | 'server', fn: () => Promise<unknown>) => {
    setOrigin(which);
    void action.run(fn);
  };

  /** Erro da ação (o Go recusou o pedido) no cartão de onde ela veio. */
  const actionError = (where: Origin) =>
    action.error && origin === where ? (
      <ErrorCard
        error={action.error}
        onClose={() => action.setError(null)}
        onOpenLogs={logs.showErrors}
      />
    ) : null;

  return (
    <div className="page">
      <PageHead title={t('install.title')} subtitle={t('install.subtitle')} />

      <section className="card">
        <h2 className="card__title">{t('install.folders')}</h2>
        {actionError('folders')}
        <div className="dirrow">
          <div>
            <span className="field__label">{t('install.serverDir')}</span>
            <code className="dirrow__path">{profile?.serverDir || t('common.notChosen')}</code>
            <span className="field__hint">{t('install.serverDirHint')}</span>
          </div>
          <button type="button" className="btn" disabled={working} onClick={() => choose('server')}>
            <FolderOpen size={16} aria-hidden />
            {t('common.choose')}
          </button>
        </div>
        <div className="dirrow">
          <div>
            <span className="field__label">{t('install.steamcmdDir')}</span>
            <code className="dirrow__path">{profile?.steamcmdDir || t('common.notChosen')}</code>
            <span className="field__hint">{t('install.steamcmdDirHint')}</span>
          </div>
          <button
            type="button"
            className="btn"
            disabled={working}
            onClick={() => choose('steamcmd')}
          >
            <FolderOpen size={16} aria-hidden />
            {t('common.choose')}
          </button>
        </div>
      </section>

      <div className="grid grid--2">
        <section className="card">
          <h2 className="card__title">
            <ShieldCheck size={18} aria-hidden /> {t('install.steamcmd.title')}
          </h2>
          <p className="muted">{t('install.steamcmd.text')}</p>
          <p>
            <span className={`pill ${info?.steamcmdFound ? 'pill--ok' : ''}`}>
              {info?.steamcmdFound ? t('install.steamcmd.ready') : t('install.steamcmd.missing')}
            </span>
          </p>
          <button
            type="button"
            className={`btn ${info?.steamcmdFound ? '' : 'btn--primary'}`}
            disabled={working || !profile?.serverDir}
            onClick={() => start('steamcmd', backend.fetchSteamCmd)}
          >
            {info?.steamcmdFound ? (
              <ShieldCheck size={16} aria-hidden />
            ) : (
              <Download size={16} aria-hidden />
            )}
            {info?.steamcmdFound ? t('install.steamcmd.check') : t('install.steamcmd.fetch')}
          </button>
          {info?.steamcmdFound && <p className="hint">{t('install.steamcmd.selfUpdates')}</p>}
          {actionError('steamcmd')}
          <TaskPanel which="steamcmd" st={st} onOpenLogs={logs.showErrors} />
        </section>

        <section className="card">
          <h2 className="card__title">{t('install.server.title')}</h2>
          <p>
            <span className={`pill ${info?.serverInstalled ? 'pill--ok' : ''}`}>
              {info?.serverInstalled
                ? t('install.server.installed')
                : t('install.server.notInstalled')}
            </span>
          </p>
          {info?.serverInstalled && (
            <p className="muted">
              {info.buildId && t('install.server.build', { id: info.buildId })}
              {info.sizeOnDisk > 0 &&
                ` · ${t('install.server.size', { size: bytes(info.sizeOnDisk, t) })}`}
              {info.lastUpdated > 0 &&
                ` · ${t('install.server.updated', { date: date(new Date(info.lastUpdated * 1000)) })}`}
            </p>
          )}
          {running && <p className="banner banner--warn">{t('install.server.runningWarn')}</p>}
          <div className="row row--wrap">
            <button
              type="button"
              className="btn btn--primary"
              disabled={working || running || !info?.steamcmdFound || !profile?.serverDir}
              onClick={() => start('server', () => backend.installServer(false))}
            >
              {info?.serverInstalled ? (
                <RefreshCw size={16} aria-hidden />
              ) : (
                <Download size={16} aria-hidden />
              )}
              {info?.serverInstalled ? t('install.server.update') : t('install.server.install')}
            </button>
            {info?.serverInstalled && (
              <button
                type="button"
                className="btn"
                disabled={working || running || !info?.steamcmdFound}
                onClick={() => start('server', () => backend.installServer(true))}
              >
                {t('install.server.validate')}
              </button>
            )}
          </div>
          <p className="hint">{t('install.server.validateHint')}</p>
          {actionError('server')}
          <TaskPanel which="server" st={st} onOpenLogs={logs.showErrors} />
        </section>
      </div>

      {info?.serverInstalled && (
        <section className="card">
          <h2 className="card__title">{t('install.import.title')}</h2>
          <p className="muted">{t('install.import.text')}</p>
          <button
            type="button"
            className="btn"
            disabled={working || sync.busy}
            onClick={() => sync.check('manual')}
          >
            <Import size={16} aria-hidden />
            {t('install.import.button')}
          </button>
        </section>
      )}
      {sync.node}
    </div>
  );
}

/**
 * Andamento e resultado da última tarefa no cartão que a iniciou (o Go avisa assim que ela começa). Ao começar,
 * a área entra na tela se a página estava rolada.
 */
export function TaskPanel({
  which,
  st,
  onOpenLogs,
}: {
  which: 'steamcmd' | 'server';
  st: InstallState | null;
  onOpenLogs: () => void;
}) {
  const { t } = useI18n();
  const ref = useRef<HTMLDivElement>(null);
  const mine = st?.last === which;
  const working = mine && st?.task === which;

  useEffect(() => {
    if (working) ref.current?.scrollIntoView?.({ block: 'nearest', behavior: 'smooth' });
  }, [working]);

  if (!st || !mine) return null;
  const phase = st.phase || '';

  let body: ReactNode = null;
  if (working) {
    body = (
      <>
        <div className="progress-head">
          <strong>{phase ? t(`install.progress.${phase}`) : t('install.progress.waiting')}</strong>
          <span className="muted">
            {st.total > 0 && `${bytes(st.done, t)} / ${bytes(st.total, t)}`}
            {st.speed > 0 && ` · ${t('units.perSecond', { value: bytes(st.speed, t) })}`}
            {st.eta > 0 && ` · ${t('units.left', { time: duration(st.eta, t) })}`}
          </span>
        </div>
        <ProgressBar value={st.percent > 0 ? st.percent : undefined} label={phase} />
        <div className="row">
          <span className="muted">{st.percent > 0 ? `${st.percent.toFixed(1)}%` : ''}</span>
          <button type="button" className="btn btn--small" onClick={() => backend.cancelInstall()}>
            <X size={14} aria-hidden />
            {t('install.progress.cancel')}
          </button>
        </div>
      </>
    );
  } else if (st.error) {
    body = <ErrorCard error={st.error} onOpenLogs={onOpenLogs} />;
  } else if (st.finished === which) {
    const key =
      which === 'steamcmd'
        ? st.present
          ? 'steamcmdPresent'
          : 'steamcmd'
        : st.validate
          ? 'validated'
          : st.upToDate
            ? 'upToDate'
            : 'server';
    body = <p className="banner banner--ok">{t(`install.progress.finished.${key}`)}</p>;
  }
  if (!body) return null;
  return (
    <div ref={ref} className="taskpanel" aria-live="polite">
      {body}
    </div>
  );
}
