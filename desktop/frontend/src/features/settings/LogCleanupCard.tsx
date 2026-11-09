import { useEffect, useState } from 'react';
import { Trash2 } from 'lucide-react';
import { useI18n } from '../../i18n';
import { backend, type LogCleanup } from '../../shared/backend';
import { ErrorCard } from '../../shared/components';
import { bytes } from '../../shared/format';
import { toAppError, type AppError } from '../../shared/problem';

/**
 * Logs do servidor (ADR-018): o servidor guarda um arquivo por partida. A limpeza automática vem desligada (não
 * apaga nada); ligada, apaga as partidas mais antigas a cada início, mantendo a quantidade escolhida. O botão
 * apaga agora, sempre com confirmação. O Go decide o que apagar e nunca apaga o log da partida atual.
 */
export function LogCleanupCard() {
  const { t } = useI18n();
  const [info, setInfo] = useState<LogCleanup | null>(null);
  const [keep, setKeep] = useState('');
  const [error, setError] = useState<AppError | null>(null);
  const [confirming, setConfirming] = useState(false);
  const [busy, setBusy] = useState(false);
  const [result, setResult] = useState<LogCleanup | null>(null);

  useEffect(() => {
    backend
      .logCleanup()
      .then((c) => {
        setInfo(c);
        setKeep(String(c.keep));
      })
      .catch(() => {});
  }, []);

  if (!info) return null;

  const save = async (enabled: boolean, keepText: string) => {
    setError(null);
    setResult(null);
    const n = Number.parseInt(keepText, 10);
    try {
      const c = await backend.setLogCleanup(enabled, Number.isFinite(n) ? n : 0);
      setInfo(c);
      setKeep(String(c.keep));
    } catch (e) {
      setError(toAppError(e));
    }
  };

  const clean = async () => {
    setBusy(true);
    setError(null);
    try {
      const c = await backend.cleanOldLogs();
      setInfo(c);
      setResult(c);
    } catch (e) {
      setError(toAppError(e));
    } finally {
      setBusy(false);
      setConfirming(false);
    }
  };

  const fieldError = error?.field === 'logKeep' ? t(error.msg ?? '', error.params) : '';

  return (
    <section className="card">
      <h2 className="card__title">{t('settings.logs.title')}</h2>
      <p className="muted">{t('settings.logs.text', { n: info.runs })}</p>

      <label className="check check--block">
        <input
          type="checkbox"
          checked={info.enabled}
          onChange={(e) => save(e.target.checked, keep)}
        />
        <span>
          {t('settings.logs.auto')}
          <span className="field__hint">
            {t(info.enabled ? 'settings.logs.autoOn' : 'settings.logs.autoOff')}
          </span>
        </span>
      </label>

      <label className="field logclean__keep">
        <span className="field__label">{t('settings.logs.keep')}</span>
        <input
          className={`input input--short ${fieldError ? 'is-invalid' : ''}`}
          type="number"
          inputMode="numeric"
          min={info.min}
          max={info.max}
          value={keep}
          onChange={(e) => setKeep(e.target.value)}
          onBlur={() => keep !== String(info.keep) && save(info.enabled, keep)}
        />
        <span className={fieldError ? 'field__error' : 'field__hint'}>
          {fieldError || t('settings.logs.keepHint', { min: info.min, max: info.max })}
        </span>
      </label>

      {error && !fieldError && <ErrorCard error={error} onClose={() => setError(null)} />}

      <div className="row row--wrap logclean__actions">
        {confirming ? (
          <>
            <span>
              {t('settings.logs.confirm', {
                n: info.old,
                size: bytes(info.oldSize, t),
                keep: info.keep,
              })}
            </span>
            <button
              type="button"
              className="btn btn--danger btn--small"
              disabled={busy}
              onClick={clean}
            >
              <Trash2 size={14} aria-hidden /> {t('settings.logs.confirmYes', { n: info.old })}
            </button>
            <button
              type="button"
              className="btn btn--small"
              disabled={busy}
              onClick={() => setConfirming(false)}
            >
              {t('common.cancel')}
            </button>
          </>
        ) : (
          <button
            type="button"
            className="btn btn--small"
            disabled={info.old === 0}
            onClick={() => {
              setResult(null);
              setConfirming(true);
            }}
          >
            <Trash2 size={14} aria-hidden /> {t('settings.logs.cleanNow')}
          </button>
        )}
      </div>
      <p className="hint">
        {result
          ? t('settings.logs.done', { n: result.deleted, size: bytes(result.freed, t) })
          : info.old > 0
            ? t('settings.logs.wouldDelete', {
                n: info.old,
                size: bytes(info.oldSize, t),
                keep: info.keep,
              })
            : t('settings.logs.nothing', { keep: info.keep })}
      </p>
    </section>
  );
}
