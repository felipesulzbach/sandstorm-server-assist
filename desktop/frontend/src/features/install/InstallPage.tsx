import { useEffect, useState } from 'react';
import { Download, FolderOpen, Import, RefreshCw, ShieldCheck, X } from 'lucide-react';
import { useI18n } from '../../i18n';
import { backend, type InstallInfo, type Profile, type Status } from '../../shared/backend';
import { ErrorCard, PageHead, ProgressBar } from '../../shared/components';
import { bytes, duration } from '../../shared/format';
import { useAction, useInstallState } from '../../shared/hooks';
import { useLogs } from '../logs/LogDrawer';
import { useExistingSync } from '../sync/ExistingSync';

export function InstallPage({ status }: { status: Status | null }) {
  const { t, date } = useI18n();
  const logs = useLogs();
  const st = useInstallState();
  const action = useAction();
  const [profile, setProfile] = useState<Profile | null>(null);
  const [info, setInfo] = useState<InstallInfo | null>(null);
  const sync = useExistingSync((r) => {
    setProfile(r.profile);
    backend.installInfo().then(setInfo).catch(() => {});
  });

  const running = (status?.pids?.length ?? 0) > 0;
  const task = st?.task ?? '';
  const working = task !== '';

  const reload = () => {
    backend.profile().then(setProfile).catch(() => {});
    backend.installInfo().then(setInfo).catch(() => {});
  };
  useEffect(reload, [st?.finished, task]);

  const choose = async (which: 'server' | 'steamcmd') => {
    const before = profile?.serverDir;
    const p = await action.run(which === 'server' ? backend.chooseServerDir : backend.chooseSteamCmdDir);
    if (p) setProfile(p);
    const next = await backend.installInfo().catch(() => null);
    if (next) setInfo(next);
    // Pasta nova com um servidor já instalado: busca os dados dele para completar os campos.
    if (which === 'server' && p?.serverDir && p.serverDir !== before && next?.serverInstalled) {
      await sync.check('auto');
    }
  };

  const phase = st?.phase || '';
  const shownError = action.error ?? (st?.error && !working ? st.error : null);

  return (
    <div className="page">
      <PageHead title={t('install.title')} subtitle={t('install.subtitle')} />
      {shownError && <ErrorCard error={shownError} onClose={() => action.setError(null)} onOpenLogs={logs.showErrors} />}

      <section className="card">
        <h2 className="card__title">{t('install.folders')}</h2>
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
          <button type="button" className="btn" disabled={working} onClick={() => choose('steamcmd')}>
            <FolderOpen size={16} aria-hidden />
            {t('common.choose')}
          </button>
        </div>
      </section>

      {working || st?.finished ? (
        <section className="card">
          <h2 className="card__title">{t('install.progress.title')}</h2>
          {working ? (
            <>
              <div className="progress-head">
                <strong>{phase ? t(`install.progress.${phase}`) : t('install.progress.waiting')}</strong>
                <span className="muted">
                  {st!.total > 0 && `${bytes(st!.done, t)} / ${bytes(st!.total, t)}`}
                  {st!.speed > 0 && ` · ${t('units.perSecond', { value: bytes(st!.speed, t) })}`}
                  {st!.eta > 0 && ` · ${t('units.left', { time: duration(st!.eta, t) })}`}
                </span>
              </div>
              <ProgressBar value={st!.percent > 0 ? st!.percent : undefined} label={phase} />
              <div className="row">
                <span className="muted">{st!.percent > 0 ? `${st!.percent.toFixed(1)}%` : ''}</span>
                <button type="button" className="btn btn--small" onClick={() => backend.cancelInstall()}>
                  <X size={14} aria-hidden />
                  {t('install.progress.cancel')}
                </button>
              </div>
            </>
          ) : (
            <p className="banner banner--ok">{t(`install.progress.finished.${st!.finished}`)}</p>
          )}
        </section>
      ) : null}

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
            onClick={() => action.run(backend.fetchSteamCmd)}
          >
            <Download size={16} aria-hidden />
            {t('install.steamcmd.fetch')}
          </button>
        </section>

        <section className="card">
          <h2 className="card__title">{t('install.server.title')}</h2>
          <p>
            <span className={`pill ${info?.serverInstalled ? 'pill--ok' : ''}`}>
              {info?.serverInstalled ? t('install.server.installed') : t('install.server.notInstalled')}
            </span>
          </p>
          {info?.serverInstalled && (
            <p className="muted">
              {info.buildId && t('install.server.build', { id: info.buildId })}
              {info.sizeOnDisk > 0 && ` · ${t('install.server.size', { size: bytes(info.sizeOnDisk, t) })}`}
              {info.lastUpdated > 0 && ` · ${t('install.server.updated', { date: date(new Date(info.lastUpdated * 1000)) })}`}
            </p>
          )}
          {running && <p className="banner banner--warn">{t('install.server.runningWarn')}</p>}
          <div className="row row--wrap">
            <button
              type="button"
              className="btn btn--primary"
              disabled={working || running || !info?.steamcmdFound || !profile?.serverDir}
              onClick={() => action.run(() => backend.installServer(false))}
            >
              {info?.serverInstalled ? <RefreshCw size={16} aria-hidden /> : <Download size={16} aria-hidden />}
              {info?.serverInstalled ? t('install.server.update') : t('install.server.install')}
            </button>
            {info?.serverInstalled && (
              <button
                type="button"
                className="btn"
                disabled={working || running || !info?.steamcmdFound}
                onClick={() => action.run(() => backend.installServer(true))}
              >
                {t('install.server.validate')}
              </button>
            )}
          </div>
          <p className="hint">{t('install.server.validateHint')}</p>
        </section>
      </div>

      {info?.serverInstalled && (
        <section className="card">
          <h2 className="card__title">{t('install.import.title')}</h2>
          <p className="muted">{t('install.import.text')}</p>
          <button type="button" className="btn" disabled={working || sync.busy} onClick={() => sync.check('manual')}>
            <Import size={16} aria-hidden />
            {t('install.import.button')}
          </button>
        </section>
      )}
      {sync.node}
    </div>
  );
}
