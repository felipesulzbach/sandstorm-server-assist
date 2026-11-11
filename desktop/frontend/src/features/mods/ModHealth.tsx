import { useCallback, useEffect, useRef, useState } from 'react';
import { AlertTriangle, CheckCircle2, RefreshCw, Search } from 'lucide-react';
import { useI18n } from '../../i18n';
import { backend, type CrashCheck, type ModCheck, type ModRef } from '../../shared/backend';
import { ErrorCard, LinkButton } from '../../shared/components';
import { useAction } from '../../shared/hooks';

/** "Nome (#id)" de um mod. */
const label = (m: ModRef) => `${m.name} (#${m.id})`;

/**
 * Mods com problema (ADR-039). O Go lê os índices dos .pak dos mods instalados e compara os mods dois a dois: o mesmo
 * arquivo do jogo em tamanhos diferentes pode fazer o jogo fechar com "Serial size mismatch". A segunda parte
 * recebe a mensagem de um crash (de um jogador ou do servidor) e diz qual mod traz o arquivo citado.
 */
export function ModHealth({ refreshKey }: { refreshKey: number }) {
  const { t } = useI18n();
  const check = useAction();
  const [result, setResult] = useState<ModCheck | null>(null);

  const run = useCallback(async () => {
    const r = await check.run(backend.checkMods);
    if (r) setResult(r);
  }, [check]);

  // Verifica de novo quando a lista de mods muda (refreshKey); a referência evita depender de "run".
  const runRef = useRef(run);
  runRef.current = run;
  useEffect(() => {
    void runRef.current();
  }, [refreshKey]);

  const conflicts = result?.conflicts ?? [];
  const unreadable = result?.unreadable ?? [];

  return (
    <section className="card modhealth">
      <div className="card__head">
        <h2 className="card__title">{t('mods.health.title')}</h2>
        <button
          type="button"
          className="btn btn--small btn--ghost"
          disabled={check.busy}
          onClick={run}
        >
          <RefreshCw size={14} aria-hidden /> {t('mods.health.recheck')}
        </button>
      </div>
      <p className="muted">{t('mods.health.text')}</p>
      {check.error && <ErrorCard error={check.error} onClose={() => check.setError(null)} />}

      {result && conflicts.length === 0 && (
        <p className="modhealth__ok">
          <CheckCircle2 size={16} aria-hidden /> {t('mods.health.none', { n: result.mods })}
        </p>
      )}
      {conflicts.length > 0 && (
        <div className="banner banner--warn">
          <p>
            <AlertTriangle size={16} aria-hidden />{' '}
            {t('mods.health.found', { n: conflicts.length })}
          </p>
          <ul className="modhealth__list">
            {conflicts.map((c) => (
              <li key={`${c.a.id}-${c.b.id}`}>
                <strong>{label(c.a)}</strong> {t('mods.health.and')} <strong>{label(c.b)}</strong>
                {': '}
                {t('mods.health.pair', { n: c.count })}
                <span className="modhealth__example">
                  {t('mods.health.example', { file: c.example })}
                </span>
              </li>
            ))}
          </ul>
          <p className="hint">{t('mods.health.advice')}</p>
          <LinkButton linkKey="modio-subscriptions" />
        </div>
      )}
      {unreadable.length > 0 && (
        <p className="hint">
          {t('mods.health.unreadable', { names: unreadable.map(label).join(', ') })}
        </p>
      )}

      <CrashLookup />
    </section>
  );
}

/** Colar a mensagem de um crash e achar o mod que traz o arquivo citado. */
function CrashLookup() {
  const { t } = useI18n();
  const action = useAction();
  const [text, setText] = useState('');
  const [result, setResult] = useState<CrashCheck | null>(null);

  const find = async () => {
    setResult(null);
    const r = await action.run(() => backend.checkCrash(text));
    if (r) setResult(r);
  };

  const findings = result?.findings ?? [];
  const found = findings.filter((f) => (f.mods ?? []).length > 0);

  return (
    <details className="modhealth__crash">
      <summary>{t('mods.crash.title')}</summary>
      <p className="hint">{t('mods.crash.text')}</p>
      <textarea
        className="input"
        rows={5}
        maxLength={65536}
        spellCheck={false}
        placeholder={t('mods.crash.placeholder')}
        value={text}
        onChange={(e) => setText(e.target.value)}
      />
      <div className="row row--wrap">
        <button
          type="button"
          className="btn btn--small"
          disabled={action.busy || text.trim() === ''}
          onClick={find}
        >
          <Search size={14} aria-hidden /> {t('mods.crash.find')}
        </button>
      </div>
      {action.error && <ErrorCard error={action.error} onClose={() => action.setError(null)} />}
      {result && findings.length === 0 && <p className="hint">{t('mods.crash.noObject')}</p>}
      {findings.length > 0 && (
        <ul className="modhealth__list">
          {findings.map((f) => (
            <li key={f.object}>
              <code>{f.object}</code>
              {': '}
              {(f.mods ?? []).length > 0
                ? t('mods.crash.inMods', { names: (f.mods ?? []).map(label).join(', ') })
                : t('mods.crash.inNoMod')}
            </li>
          ))}
        </ul>
      )}
      {found.length > 0 && (
        <p className="hint">
          {t(result?.sizeMismatch ? 'mods.crash.adviceMismatch' : 'mods.crash.advice')}
        </p>
      )}
      {result && findings.length > 0 && found.length === 0 && (
        <p className="hint">{t('mods.crash.adviceNone')}</p>
      )}
    </details>
  );
}
