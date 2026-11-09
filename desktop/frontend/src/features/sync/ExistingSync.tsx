import { useEffect, useState, type ReactNode } from 'react';
import { ServerCog, X } from 'lucide-react';
import { useI18n } from '../../i18n';
import { backend, type ApplyResult, type ExistingScan, type SyncDiff } from '../../shared/backend';
import { ErrorCard } from '../../shared/components';
import { useAction } from '../../shared/hooks';

type Choice = 'mine' | 'server';
type Notice = { kind: 'ok' | 'info'; lines: string[] };

/**
 * Busca os dados de um servidor já configurado (rodando ou parado) e completa os campos do app. O Go compara campo
 * a campo; aqui só se pergunta o que o usuário mudou e o servidor tem diferente (manter o meu ou usar o do servidor).
 * mode "auto" (logo depois de escolher a pasta) não avisa quando não acha nada.
 */
export function useExistingSync(onApplied: (r: ApplyResult) => void) {
  const { t } = useI18n();
  const action = useAction();
  const [scan, setScan] = useState<ExistingScan | null>(null);
  const [notice, setNotice] = useState<Notice | null>(null);

  const check = async (mode: 'auto' | 'manual') => {
    setNotice(null);
    const s = await action.run(backend.scanExisting);
    if (!s) return;
    const pending = (s.diffs ?? []).filter((d) => d.state !== 'same');
    if (pending.length > 0) {
      setScan(s);
    } else if (mode === 'manual') {
      setNotice({ kind: 'info', lines: [s.found ? t('sync.upToDate') : t('sync.notFound')] });
    }
  };

  const apply = async (keepMine: string[]) => {
    const r = await action.run(() => backend.applyExisting(keepMine));
    if (!r) return;
    setScan(null);
    const lines = [t('sync.done', { applied: r.applied?.length ?? 0, kept: r.kept?.length ?? 0 })];
    if (r.skipped?.length) {
      lines.push(
        t('sync.skipped', { fields: r.skipped.map((k) => t(`sync.fields.${k}`)).join(', ') }),
      );
    }
    if (r.rconInIni) lines.push(t('sync.rconInIni'));
    setNotice({ kind: 'ok', lines });
    onApplied(r);
  };

  const node: ReactNode = (
    <>
      {action.error && <ErrorCard error={action.error} onClose={() => action.setError(null)} />}
      {notice && (
        <div className={`banner banner--${notice.kind}`}>
          {notice.lines.map((l) => (
            <div key={l}>{l}</div>
          ))}
        </div>
      )}
      {scan && (
        <SyncDialog scan={scan} busy={action.busy} onCancel={() => setScan(null)} onApply={apply} />
      )}
    </>
  );
  return { check, busy: action.busy, node };
}

function SyncDialog({
  scan,
  busy,
  onCancel,
  onApply,
}: {
  scan: ExistingScan;
  busy: boolean;
  onCancel: () => void;
  onApply: (keepMine: string[]) => void;
}) {
  const { t } = useI18n();
  const diffs = scan.diffs ?? [];
  const conflicts = diffs.filter((d) => d.state === 'conflict');
  const fills = diffs.filter((d) => d.state === 'fill');
  const same = diffs.filter((d) => d.state === 'same');
  const [choice, setChoice] = useState<Record<string, Choice>>({});
  const decided = conflicts.every((d) => choice[d.key]);

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => e.key === 'Escape' && !busy && onCancel();
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
  }, [busy, onCancel]);

  const setAll = (c: Choice) => setChoice(Object.fromEntries(conflicts.map((d) => [d.key, c])));
  const source = [
    scan.logFile ? t('sync.source', { file: scan.logFile }) : t('sync.sourceFiles'),
    scan.running ? t('sync.running') : t('sync.stopped'),
  ].join(' · ');

  return (
    <div className="modal-backdrop">
      <div className="modal" role="dialog" aria-modal="true" aria-labelledby="sync-title">
        <header className="modal__head">
          <div>
            <h2 id="sync-title" className="card__title">
              <ServerCog size={18} aria-hidden /> {t('sync.title')}
            </h2>
            <p className="muted">{source}</p>
          </div>
          <button
            type="button"
            className="btn btn--ghost btn--small"
            aria-label={t('common.close')}
            disabled={busy}
            onClick={onCancel}
          >
            <X size={16} aria-hidden />
          </button>
        </header>

        <div className="modal__body">
          {conflicts.length > 0 && (
            <section className="sync__group">
              <h3 className="sync__title">{t('sync.conflictsTitle', { n: conflicts.length })}</h3>
              <p className="hint">{t('sync.conflictsText')}</p>
              <div className="row row--wrap">
                <button type="button" className="btn btn--small" onClick={() => setAll('mine')}>
                  {t('sync.keepAllMine')}
                </button>
                <button type="button" className="btn btn--small" onClick={() => setAll('server')}>
                  {t('sync.useAllServer')}
                </button>
              </div>
              <ul className="sync__list">
                {conflicts.map((d) => (
                  <li key={d.key} className="sync__row sync__row--conflict">
                    <span className="sync__label">{t(`sync.fields.${d.key}`)}</span>
                    <div
                      className="sync__choices"
                      role="radiogroup"
                      aria-label={t(`sync.fields.${d.key}`)}
                    >
                      {(['mine', 'server'] as const).map((c) => (
                        <button
                          key={c}
                          type="button"
                          role="radio"
                          aria-checked={choice[d.key] === c}
                          className={`sync__choice ${choice[d.key] === c ? 'is-active' : ''}`}
                          onClick={() => setChoice((x) => ({ ...x, [d.key]: c }))}
                        >
                          <span className="sync__side">
                            {t(c === 'mine' ? 'sync.mine' : 'sync.server')}
                          </span>
                          <span className="sync__value">
                            <Value diff={d} side={c} />
                          </span>
                        </button>
                      ))}
                    </div>
                  </li>
                ))}
              </ul>
            </section>
          )}

          {fills.length > 0 && (
            <section className="sync__group">
              <h3 className="sync__title">{t('sync.fillTitle', { n: fills.length })}</h3>
              <p className="hint">{t('sync.fillText')}</p>
              <ul className="sync__list">
                {fills.map((d) => (
                  <li key={d.key} className="sync__row">
                    <span className="sync__label">{t(`sync.fields.${d.key}`)}</span>
                    <span className="sync__value">
                      <Value diff={d} side="server" />
                    </span>
                  </li>
                ))}
              </ul>
            </section>
          )}

          {same.length > 0 && <p className="hint">{t('sync.sameCount', { n: same.length })}</p>}
        </div>

        <footer className="modal__foot">
          {!decided && <span className="hint">{t('sync.chooseEach')}</span>}
          <button type="button" className="btn" disabled={busy} onClick={onCancel}>
            {t('common.cancel')}
          </button>
          <button
            type="button"
            className="btn btn--primary"
            disabled={busy || !decided}
            onClick={() =>
              onApply(conflicts.filter((d) => choice[d.key] === 'mine').map((d) => d.key))
            }
          >
            {t('sync.apply')}
          </button>
        </footer>
      </div>
    </div>
  );
}

function Value({ diff, side }: { diff: SyncDiff; side: Choice }) {
  const { t } = useI18n();
  const v = side === 'mine' ? diff.mine : diff.server;
  switch (diff.kind) {
    case 'secret':
      if (side === 'mine') return <em>{t('sync.secretMine')}</em>;
      return <em>{t(diff.state === 'fill' ? 'sync.secretFound' : 'sync.secretServer')}</em>;
    case 'bool':
      return <>{v === 'true' ? t('common.yes') : t('common.no')}</>;
    case 'list': {
      const lines = v ? v.split('\n') : [];
      if (lines.length === 0) return <em>{t('sync.empty')}</em>;
      if (lines.length <= 3) return <span title={v}>{lines.join(', ')}</span>;
      return (
        <span title={lines.slice(0, 20).join('\n')}>{t('sync.items', { n: lines.length })}</span>
      );
    }
    default:
      return v ? <>{v}</> : <em>{t('sync.empty')}</em>;
  }
}
