import { useState } from 'react';
import { RefreshCw } from 'lucide-react';
import { useI18n } from '../../i18n';
import type { AppInfo, Terms, UpdateState } from '../../shared/backend';
import { AppIcon } from '../../shared/AppIcon';
import { ErrorCard, LinkButton, PageHead, ProgressBar } from '../../shared/components';
import { toAppError } from '../../shared/problem';
import { useUpdate } from '../../shared/useUpdate';
import { LegalDialog } from '../legal/Legal';

/** Texto da situação da atualização automática. */
function updateText(
  t: (k: string, p?: Record<string, string | number>) => string,
  s: UpdateState,
): string {
  if (s.status === 'error')
    return t(s.stage === 'download' ? 'update.errorDownload' : 'update.errorCheck');
  return t(`update.status.${s.status}`, {
    version: s.version,
    percent: Math.floor(s.progress * 100),
  });
}

function UpdatesCard() {
  const { t } = useI18n();
  const { state, error, check, install } = useUpdate();
  const busy = state?.status === 'checking' || state?.status === 'downloading';
  return (
    <section className="card">
      <h2 className="card__title">{t('update.title')}</h2>
      {state && <p>{updateText(t, state)}</p>}
      {state?.status === 'downloading' && (
        <ProgressBar value={state.progress * 100} label={t('update.title')} />
      )}
      {error !== null && <ErrorCard error={toAppError(error)} />}
      <div className="row row--wrap">
        {state?.status === 'ready' && (
          <button type="button" className="btn btn--primary btn--small" onClick={install}>
            {t('update.install')}
          </button>
        )}
        {state?.status === 'manual' && <LinkButton linkKey="download-page" />}
        {state && state.status !== 'dev' && state.status !== 'ready' && (
          <button type="button" className="btn btn--small" disabled={busy} onClick={check}>
            <RefreshCw size={14} aria-hidden /> {t('update.check')}
          </button>
        )}
      </div>
      <p className="hint">{t('update.hint')}</p>
    </section>
  );
}

export function AboutPage({ info, terms }: { info: AppInfo | null; terms: Terms | null }) {
  const { t, date } = useI18n();
  const [reading, setReading] = useState<'terms' | 'privacy' | null>(null);
  return (
    <div className="page page--narrow">
      <PageHead title={t('about.title')} />
      <section className="card about">
        <AppIcon size={48} />
        <div>
          <h2>{t('app.name')}</h2>
          {info && <p className="muted">{t('app.version', { version: info.appVersion })}</p>}
          <p>{t('about.text')}</p>
        </div>
      </section>
      <UpdatesCard />
      <section className="card">
        <h2 className="card__title">{t('legal.title')}</h2>
        {terms?.accepted && (
          <p className="muted">
            {t('legal.acceptedOn', { version: terms.version, date: date(terms.acceptedAt) })}
          </p>
        )}
        <div className="row row--wrap">
          <button
            type="button"
            className="btn btn--small btn--ghost"
            onClick={() => setReading('terms')}
          >
            {t('legal.tabs.terms')}
          </button>
          <button
            type="button"
            className="btn btn--small btn--ghost"
            onClick={() => setReading('privacy')}
          >
            {t('legal.tabs.privacy')}
          </button>
        </div>
      </section>
      <section className="card">
        <h2 className="card__title">{t('about.credits')}</h2>
        <p>{t('about.creditsText')}</p>
        <div className="row row--wrap">
          <LinkButton linkKey="project" />
          <LinkButton linkKey="user-guide" />
          <LinkButton linkKey="author" />
        </div>
      </section>
      <section className="card">
        <h2 className="card__title">{t('about.license')}</h2>
        <p>{t('about.licenseText')}</p>
        <LinkButton linkKey="apache-license" />
      </section>
      <section className="card">
        <h2 className="card__title">{t('about.disclaimer')}</h2>
        <p className="muted">{t('about.disclaimerText')}</p>
      </section>
      <section className="card">
        <h2 className="card__title">{t('about.links')}</h2>
        <div className="row row--wrap">
          <LinkButton linkKey="server-guide" label={t('about.serverGuide')} />
          <LinkButton linkKey="modio-admin-guide" label={t('about.modioGuide')} />
          <LinkButton linkKey="steamcmd-docs" />
          <LinkButton linkKey="wails" />
        </div>
      </section>
      {reading && <LegalDialog initial={reading} onClose={() => setReading(null)} />}
    </div>
  );
}
