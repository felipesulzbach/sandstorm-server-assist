import { useState } from 'react';
import { FileUp, FolderOpen, Lock, Save, Unlock } from 'lucide-react';
import { useI18n } from '../../i18n';
import { backend, type BackupSummary, type Status } from '../../shared/backend';
import { ErrorCard, Field, PageHead } from '../../shared/components';
import { useAction } from '../../shared/hooks';

const MIN = 10;

export function BackupPage({ status }: { status: Status | null }) {
  const { t } = useI18n();
  return (
    <div className="page">
      <PageHead title={t('backup.title')} subtitle={t('backup.subtitle')} />
      <div className="grid grid--2">
        <CreateCard />
        <RestoreCard running={(status?.pids?.length ?? 0) > 0} />
      </div>
    </div>
  );
}

function CreateCard() {
  const { t } = useI18n();
  const action = useAction();
  const [pw, setPw] = useState('');
  const [pw2, setPw2] = useState('');
  const [modio, setModio] = useState(false);
  const [done, setDone] = useState<string | null>(null);
  const mismatch = pw2.length > 0 && pw !== pw2;

  return (
    <section className="card">
      <h2 className="card__title">
        <Lock size={18} aria-hidden /> {t('backup.create.title')}
      </h2>
      <p className="muted">{t('backup.create.contents')}</p>
      {action.error && <ErrorCard error={action.error} onClose={() => action.setError(null)} />}
      <div className="form form--single">
        <Field label={t('backup.create.password')} hint={t('backup.create.passwordHint')}>
          <input className="input" type="password" autoComplete="new-password" value={pw} onChange={(e) => setPw(e.target.value)} />
        </Field>
        <Field label={t('backup.create.confirm')} error={mismatch ? t('backup.create.mismatch') : null}>
          <input className="input" type="password" autoComplete="new-password" value={pw2} onChange={(e) => setPw2(e.target.value)} />
        </Field>
        <label className="check check--block">
          <input type="checkbox" checked={modio} onChange={(e) => setModio(e.target.checked)} />
          <span>
            {t('backup.create.includeModio')}
            {modio && <span className="field__hint">{t('backup.create.includeModioWarn')}</span>}
          </span>
        </label>
        <button
          type="button"
          className="btn btn--primary"
          disabled={action.busy || pw.length < MIN || pw !== pw2}
          onClick={async () => {
            setDone(null);
            const r = await action.run(() => backend.createBackup(pw, modio));
            if (r?.saved) {
              setDone(r.file);
              setPw('');
              setPw2('');
            }
          }}
        >
          <Save size={16} aria-hidden />
          {t('backup.create.button')}
        </button>
        {done && <p className="banner banner--ok">{t('backup.create.done', { file: done })}</p>}
      </div>
    </section>
  );
}

function RestoreCard({ running }: { running: boolean }) {
  const { t, date } = useI18n();
  const action = useAction();
  const [file, setFile] = useState<string | null>(null);
  const [pw, setPw] = useState('');
  const [summary, setSummary] = useState<BackupSummary | null>(null);
  const [dir, setDir] = useState<string | null>(null);
  const [modio, setModio] = useState(true);
  const [done, setDone] = useState(false);

  return (
    <section className="card">
      <h2 className="card__title">
        <Unlock size={18} aria-hidden /> {t('backup.restore.title')}
      </h2>
      {action.error && <ErrorCard error={action.error} onClose={() => action.setError(null)} />}
      <div className="form form--single">
        <div className="row row--wrap">
          <button
            type="button"
            className="btn"
            onClick={async () => {
              const f = await action.run(backend.chooseBackupFile);
              if (f) {
                setFile(f);
                setSummary(null);
                setDone(false);
              }
            }}
          >
            <FileUp size={16} aria-hidden />
            {t('backup.restore.choose')}
          </button>
          {file && <span className="muted">{t('backup.restore.file', { file })}</span>}
        </div>

        {file && !summary && (
          <>
            <Field label={t('backup.restore.password')}>
              <input className="input" type="password" autoComplete="off" value={pw} onChange={(e) => setPw(e.target.value)} />
            </Field>
            <button
              type="button"
              className="btn btn--primary"
              disabled={action.busy || pw.length < MIN}
              onClick={async () => {
                const s = await action.run(() => backend.inspectBackup(pw));
                if (s) setSummary(s);
              }}
            >
              {t('backup.restore.open')}
            </button>
          </>
        )}

        {summary && (
          <>
            <dl className="facts">
              <dt>{t('backup.restore.createdAt')}</dt>
              <dd>{date(summary.createdAt as unknown as string)}</dd>
              <dt>{t('backup.restore.machine')}</dt>
              <dd>{summary.machine || '—'}</dd>
              <dt>{t('backup.restore.hostname')}</dt>
              <dd>{summary.hostname}</dd>
              <dt>{t('backup.restore.files')}</dt>
              <dd>{(summary.files ?? []).join(', ') || '—'}</dd>
              <dt>{t('backup.restore.secrets')}</dt>
              <dd>{(summary.secrets ?? []).length}</dd>
              <dt>{t('backup.restore.modio')}</dt>
              <dd>{summary.hasModio ? t('backup.restore.included') : t('backup.restore.notIncluded')}</dd>
            </dl>
            <div className="dirrow">
              <div>
                <span className="field__label">{t('backup.restore.serverDir')}</span>
                <code className="dirrow__path">{dir || summary.oldDir || t('common.notChosen')}</code>
              </div>
              <button
                type="button"
                className="btn"
                onClick={async () => {
                  const d = await action.run(backend.chooseRestoreDir);
                  if (d) setDir(d);
                }}
              >
                <FolderOpen size={16} aria-hidden />
                {t('common.choose')}
              </button>
            </div>
            {summary.hasModio && (
              <label className="check check--block">
                <input type="checkbox" checked={modio} onChange={(e) => setModio(e.target.checked)} />
                <span>{t('backup.restore.restoreModio')}</span>
              </label>
            )}
            {running && <p className="banner banner--warn">{t('backup.restore.runningWarn')}</p>}
            <button
              type="button"
              className="btn btn--primary"
              disabled={action.busy || !dir || running}
              onClick={async () => {
                const s = await action.run(() => backend.restoreBackup(pw, modio && summary.hasModio));
                if (s) {
                  setDone(true);
                  setPw('');
                }
              }}
            >
              {t('backup.restore.button')}
            </button>
            {done && <p className="banner banner--ok">{t('backup.restore.done')}</p>}
          </>
        )}
      </div>
    </section>
  );
}
