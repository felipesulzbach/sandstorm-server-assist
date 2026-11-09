import { useCallback, useEffect, useState } from 'react';
import { KeyRound, Mail, RefreshCw, Search, ShieldCheck } from 'lucide-react';
import { useI18n } from '../../i18n';
import { backend, type ModRow, type ModsView, type SecretKey, type Status } from '../../shared/backend';
import { ErrorCard, Field, Guide, LinkButton, PageHead, ProgressBar, SecretField } from '../../shared/components';
import { bytes, percent } from '../../shared/format';
import { useAction } from '../../shared/hooks';
import { useLogs } from '../logs/LogDrawer';

export function ModsPage({ status }: { status: Status | null }) {
  const { t, date } = useI18n();
  const logs = useLogs();
  const action = useAction();
  const [view, setView] = useState<ModsView | null>(null);
  const [secrets, setSecrets] = useState<Record<SecretKey, boolean> | null>(null);
  const [code, setCode] = useState('');
  const [note, setNote] = useState<string | null>(null);

  const load = useCallback(() => {
    backend.mods().then(setView).catch(() => {});
    backend.secretsState().then(setSecrets).catch(() => {});
  }, []);
  useEffect(load, [load]);
  // Enquanto o servidor roda, a lista acompanha o progresso dos downloads.
  useEffect(() => {
    if (!status?.pids?.length) return;
    backend.mods().then(setView).catch(() => {});
  }, [status]);

  const session = view?.session;
  const expires = session?.expires ? date(session.expires as unknown as string) : '';
  const running = (status?.pids?.length ?? 0) > 0;

  return (
    <div className="page">
      <PageHead title={t('mods.title')} subtitle={t('mods.subtitle')}>
        <LinkButton linkKey="modio" label={t('mods.list.openSite')} />
      </PageHead>
      {action.error && <ErrorCard error={action.error} onClose={() => action.setError(null)} onOpenLogs={logs.showErrors} />}

      <div className="grid grid--2">
        <section className="card">
          <h2 className="card__title">
            <ShieldCheck size={18} aria-hidden /> {t('mods.session.title')}
          </h2>
          {session?.valid ? (
            <p className="banner banner--ok">{t('mods.session.valid', { user: session.username || '—', date: expires })}</p>
          ) : session?.found && expires !== '—' ? (
            <p className="banner banner--warn">{t('mods.session.expired', { date: expires })}</p>
          ) : (
            <p className="banner banner--warn">{t('mods.session.none')}</p>
          )}
          <Guide id="modioAccount" links={['modio']} />
        </section>

        <section className="card">
          <h2 className="card__title">{t('mods.auth.title')}</h2>
          <p className="muted">{t('mods.auth.text')}</p>
          {secrets && (
            <>
              <SecretField name="modioEmail" plain label={t('mods.auth.email')} isSet={secrets.modioEmail} onChanged={setSecrets} />
              <SecretField name="modioApiKey" label={t('mods.auth.apiKey')} isSet={secrets.modioApiKey} onChanged={setSecrets}>
                <div className="row row--wrap">
                  <button
                    type="button"
                    className="btn btn--small"
                    onClick={async () => {
                      const found = await action.run(backend.detectApiKey);
                      if (found === undefined) return;
                      setNote(found ? t('mods.auth.detected') : t('mods.auth.notFound'));
                      load();
                    }}
                  >
                    <Search size={14} aria-hidden />
                    {t('mods.auth.detect')}
                  </button>
                </div>
                <Guide id="modioApiKey" links={['modio-api-access']} />
              </SecretField>
            </>
          )}
          {note && <p className="hint">{note}</p>}
          <div className="row row--wrap">
            <button
              type="button"
              className="btn"
              disabled={action.busy || !secrets?.modioEmail || !secrets?.modioApiKey}
              onClick={async () => {
                const ok = await action.run(async () => {
                  await backend.requestCode();
                  return true;
                });
                if (ok) setNote(t('mods.auth.codeSent'));
              }}
            >
              <Mail size={16} aria-hidden />
              {t('mods.auth.requestCode')}
            </button>
          </div>
          <Field label={t('mods.auth.code')} hint={t('mods.auth.authorizeHint')}>
            <div className="row">
              <input
                className="input input--code"
                value={code}
                maxLength={5}
                autoComplete="off"
                spellCheck={false}
                onChange={(e) => setCode(e.target.value.trim().toUpperCase())}
              />
              <button
                type="button"
                className="btn btn--primary"
                disabled={action.busy || code.length !== 5 || running}
                onClick={async () => {
                  const ok = await action.run(async () => {
                    await backend.authorize(code);
                    return true;
                  });
                  if (ok) setCode('');
                }}
              >
                <KeyRound size={16} aria-hidden />
                {t('mods.auth.authorize')}
              </button>
            </div>
          </Field>
          {status?.authorizing && <p className="banner banner--info">{t('dashboard.authorizing')}</p>}
        </section>
      </div>

      <section className="card">
        <div className="card__head">
          <h2 className="card__title">{t('mods.list.title')}</h2>
          <div className="row">
            {view?.subscribedAt && date(view.subscribedAt as unknown as string) !== '—' && (
              <span className="muted">{t('mods.list.refreshed', { date: date(view.subscribedAt as unknown as string) })}</span>
            )}
            <button
              type="button"
              className="btn btn--small"
              disabled={action.busy || !session?.valid}
              onClick={async () => {
                const v = await action.run(backend.refreshSubscriptions);
                if (v) setView(v);
              }}
            >
              <RefreshCw size={14} aria-hidden />
              {t('mods.list.refresh')}
            </button>
          </div>
        </div>
        {!view?.mods?.length ? (
          <p className="muted">{t('mods.list.empty')}</p>
        ) : (
          <table className="table">
            <thead>
              <tr>
                <th>{t('mods.list.name')}</th>
                <th>{t('mods.list.size')}</th>
                <th>{t('mods.list.status')}</th>
              </tr>
            </thead>
            <tbody>
              {view.mods.map((m) => (
                <ModLine key={m.id} m={m} hasSubs={!!view.subscribedAt && date(view.subscribedAt as unknown as string) !== '—'} />
              ))}
            </tbody>
          </table>
        )}
        <p className="hint">{t('mods.list.conflict')}</p>
      </section>
    </div>
  );
}

function ModLine({ m, hasSubs }: { m: ModRow; hasSubs: boolean }) {
  const { t } = useI18n();
  const downloading = m.live === 'installing' || m.live === 'updating';
  let state = m.live;
  if (!state) state = m.installed ? (hasSubs && !m.subscribed ? 'orphan' : 'installed') : m.subscribed ? 'subscribed' : 'installed';
  return (
    <tr>
      <td>
        <strong>{m.name}</strong> <span className="muted">#{m.id}</span>
      </td>
      <td>{m.size > 0 ? bytes(m.size, t) : '—'}</td>
      <td className="table__status">
        <span className={`pill pill--${state}`}>{t(`mods.state.${state}`)}</span>
        {downloading && (
          <div className="table__progress">
            <ProgressBar done={m.done} total={m.size} label={m.name} />
            <span className="muted">
              {m.size > 0 ? `${percent(m.done, m.size).toFixed(0)}% · ${bytes(m.done, t)} / ${bytes(m.size, t)}` : bytes(m.done, t)}
            </span>
          </div>
        )}
      </td>
    </tr>
  );
}
