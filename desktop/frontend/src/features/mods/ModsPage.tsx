import { useCallback, useEffect, useState } from 'react';
import { Download, KeyRound, Mail, RefreshCw, Search, ShieldCheck, Trash2 } from 'lucide-react';
import { useI18n } from '../../i18n';
import {
  backend,
  type ModRow,
  type ModsView,
  type SecretKey,
  type Status,
} from '../../shared/backend';
import {
  ErrorCard,
  Field,
  Guide,
  LinkButton,
  PageHead,
  ProgressBar,
  SecretField,
} from '../../shared/components';
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
    backend
      .mods()
      .then(setView)
      .catch(() => {});
    backend
      .secretsState()
      .then(setSecrets)
      .catch(() => {});
  }, []);
  useEffect(load, [load]);
  // Enquanto o servidor roda, a lista acompanha o progresso dos downloads.
  useEffect(() => {
    if (!status?.pids?.length) return;
    backend
      .mods()
      .then(setView)
      .catch(() => {});
  }, [status]);

  const session = view?.session;
  const expires = session?.expires ? date(session.expires as unknown as string) : '';
  const running = (status?.pids?.length ?? 0) > 0;

  return (
    <div className="page">
      <PageHead title={t('mods.title')} subtitle={t('mods.subtitle')}>
        <LinkButton linkKey="modio" label={t('mods.list.openSite')} />
      </PageHead>
      {action.error && (
        <ErrorCard
          error={action.error}
          onClose={() => action.setError(null)}
          onOpenLogs={logs.showErrors}
        />
      )}

      <div className="grid grid--2">
        <section className="card">
          <h2 className="card__title">
            <ShieldCheck size={18} aria-hidden /> {t('mods.session.title')}
          </h2>
          {session?.valid ? (
            <p className="banner banner--ok">
              {t('mods.session.valid', { user: session.username || '—', date: expires })}
            </p>
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
              <SecretField
                name="modioEmail"
                plain
                label={t('mods.auth.email')}
                isSet={secrets.modioEmail}
                onChanged={setSecrets}
              />
              <SecretField
                name="modioApiKey"
                label={t('mods.auth.apiKey')}
                isSet={secrets.modioApiKey}
                onChanged={setSecrets}
              >
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
          {status?.authorizing && (
            <p className="banner banner--info">{t('dashboard.authorizing')}</p>
          )}
        </section>
      </div>

      <section className="card">
        <div className="card__head">
          <h2 className="card__title">{t('mods.list.title')}</h2>
          <div className="row">
            {view?.subscribedAt && date(view.subscribedAt as unknown as string) !== '—' && (
              <span className="muted">
                {t('mods.list.refreshed', { date: date(view.subscribedAt as unknown as string) })}
              </span>
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
        <PendingMods mods={view?.mods ?? []} running={running} />
        {!view?.mods?.length ? (
          <p className="muted">{t('mods.list.empty')}</p>
        ) : (
          <table className="table">
            <thead>
              <tr>
                <th>{t('mods.list.name')}</th>
                <th>{t('mods.list.size')}</th>
                <th>{t('mods.list.status')}</th>
                <th>
                  <span className="sr-only">{t('mods.delete.column')}</span>
                </th>
              </tr>
            </thead>
            <tbody>
              {view.mods.map((m) => (
                <ModLine
                  key={m.id}
                  m={m}
                  hasSubs={
                    !!view.subscribedAt && date(view.subscribedAt as unknown as string) !== '—'
                  }
                  running={running}
                  onDeleted={load}
                />
              ))}
            </tbody>
          </table>
        )}
        <p className="hint">{t('mods.list.conflict')}</p>
      </section>
    </div>
  );
}

/**
 * Mods inscritos que ainda não estão no disco (ADR-038). Quem baixa é o mod.io do servidor, e só ao iniciar:
 * "Transferir agora" inicia o servidor, ou reinicia se ele está ligado (com confirmação: a partida acaba e os
 * jogadores saem por um momento). O andamento aparece na lista, em cada mod.
 */
export function PendingMods({ mods, running }: { mods: ModRow[]; running: boolean }) {
  const { t } = useI18n();
  const logs = useLogs();
  const action = useAction();
  const [confirming, setConfirming] = useState(false);
  const [started, setStarted] = useState(false);
  const pending = mods.filter((m) => m.subscribed && !m.installed && !m.live);

  const download = async () => {
    const ok = await action.run(async () => {
      await backend.downloadNewMods();
      return true;
    });
    setConfirming(false);
    if (ok) setStarted(true);
  };

  if (pending.length === 0 && !started && !action.error) return null;
  return (
    <div className="pendingmods" role="status">
      {action.error && (
        <ErrorCard
          error={action.error}
          onClose={() => action.setError(null)}
          onOpenLogs={logs.showErrors}
        />
      )}
      {started ? (
        <p className="banner banner--ok">{t('mods.pending.started')}</p>
      ) : (
        pending.length > 0 && (
          <div className="banner banner--info">
            <p>
              {t('mods.pending.text', { names: pending.map((m) => m.name).join(', ') })}{' '}
              {running ? t('mods.pending.runningHint') : t('mods.pending.stoppedHint')}
            </p>
            {confirming ? (
              <div className="row row--wrap">
                <span>{t('mods.pending.confirm')}</span>
                <button
                  type="button"
                  className="btn btn--primary btn--small"
                  disabled={action.busy}
                  onClick={download}
                >
                  {t('mods.pending.confirmYes')}
                </button>
                <button
                  type="button"
                  className="btn btn--small btn--ghost"
                  disabled={action.busy}
                  onClick={() => setConfirming(false)}
                >
                  {t('common.cancel')}
                </button>
              </div>
            ) : (
              <button
                type="button"
                className="btn btn--primary btn--small"
                disabled={action.busy}
                onClick={() => (running ? setConfirming(true) : download())}
              >
                <Download size={14} aria-hidden />
                {running ? t('mods.pending.restart') : t('mods.pending.start')}
              </button>
            )}
          </div>
        )
      )}
    </div>
  );
}

/**
 * Linha de um mod. Mod baixado pode ser apagado (ADR-037): "Apagar e baixar de novo" (para e inicia o servidor se
 * preciso) ou "Só apagar" (servidor parado; ele baixa de novo na próxima vez que iniciar, se o mod continua
 * inscrito). A confirmação, o andamento, o resultado e o erro ficam na própria linha.
 */
export function ModLine({
  m,
  hasSubs,
  running,
  onDeleted,
}: {
  m: ModRow;
  hasSubs: boolean;
  running: boolean;
  onDeleted: () => void;
}) {
  const { t } = useI18n();
  const logs = useLogs();
  const action = useAction();
  const [confirming, setConfirming] = useState(false);
  const [done, setDone] = useState<string | null>(null);
  const downloading = m.live === 'installing' || m.live === 'updating';
  let state = m.live;
  if (!state)
    state = m.installed
      ? hasSubs && !m.subscribed
        ? 'orphan'
        : 'installed'
      : m.subscribed
        ? 'subscribed'
        : 'installed';
  const canDelete = m.installed && !downloading;

  const remove = async (redownload: boolean) => {
    setDone(null);
    const ok = await action.run(async () => {
      await backend.deleteMod(m.id, redownload);
      return true;
    });
    if (!ok) return;
    setConfirming(false);
    setDone(t(redownload ? 'mods.delete.doneRedownload' : 'mods.delete.done', { name: m.name }));
    onDeleted();
  };

  return (
    <>
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
                {m.size > 0
                  ? `${percent(m.done, m.size).toFixed(0)}% · ${bytes(m.done, t)} / ${bytes(m.size, t)}`
                  : bytes(m.done, t)}
              </span>
            </div>
          )}
        </td>
        <td className="table__actions">
          {canDelete && (
            <button
              type="button"
              className="icon-btn"
              title={t('mods.delete.button', { name: m.name })}
              aria-label={t('mods.delete.button', { name: m.name })}
              disabled={action.busy || confirming}
              onClick={() => {
                setDone(null);
                setConfirming(true);
              }}
            >
              <Trash2 size={16} />
            </button>
          )}
        </td>
      </tr>
      {(confirming || done || action.error) && (
        <tr className="table__detail">
          <td colSpan={4}>
            {action.error && (
              <ErrorCard
                error={action.error}
                onClose={() => action.setError(null)}
                onOpenLogs={logs.showErrors}
              />
            )}
            {done && <p className="banner banner--ok">{done}</p>}
            {confirming && (
              <div
                className="moddelete"
                role="group"
                aria-label={t('mods.delete.button', { name: m.name })}
              >
                <p>{t('mods.delete.confirm', { name: m.name })}</p>
                <div className="row row--wrap">
                  <button
                    type="button"
                    className="btn btn--primary btn--small"
                    disabled={action.busy}
                    onClick={() => remove(true)}
                  >
                    {action.busy ? t('mods.delete.working') : t('mods.delete.redownload')}
                  </button>
                  <button
                    type="button"
                    className="btn btn--small"
                    disabled={action.busy || running}
                    onClick={() => remove(false)}
                  >
                    {t('mods.delete.only')}
                  </button>
                  <button
                    type="button"
                    className="btn btn--small btn--ghost"
                    disabled={action.busy}
                    onClick={() => setConfirming(false)}
                  >
                    {t('common.cancel')}
                  </button>
                </div>
                <p className="hint">
                  {running ? t('mods.delete.runningHint') : t('mods.delete.stoppedHint')}
                </p>
              </div>
            )}
          </td>
        </tr>
      )}
    </>
  );
}
