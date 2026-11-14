import { Server } from 'lucide-react';
import type { PageId } from '../App';
import { useI18n } from '../i18n';
import { ErrorCard } from './components';
import { useAction } from './hooks';
import { useServers } from './servers';

/**
 * Seletor do servidor no topo (ADR-048): só aparece com mais de um servidor. As telas Painel, Servidor, Mods e Logs
 * mostram o servidor escolhido; trocar pode ser com servidores rodando.
 */
export function ServerBar({ go }: { go: (p: PageId) => void }) {
  const { t } = useI18n();
  const { list, selected, select } = useServers();
  const action = useAction();
  if (list.length < 2) return null;
  return (
    <div className="serverbar">
      <label className="serverbar__pick">
        <Server size={16} aria-hidden />
        <span className="serverbar__label">{t('servers.label')}</span>
        <select
          className="select"
          value={selected}
          disabled={action.busy}
          onChange={(e) => action.run(() => select(e.target.value))}
        >
          {list.map((s) => (
            <option key={s.id} value={s.id}>
              {s.name} · {t(`phase.${s.phase}`)}
            </option>
          ))}
        </select>
      </label>
      <span
        className={`dot dot--${list.find((s) => s.id === selected)?.phase ?? 'stopped'}`}
        aria-hidden
      />
      <button type="button" className="btn btn--small btn--ghost" onClick={() => go('dashboard')}>
        {t('servers.all')}
      </button>
      {action.error && <ErrorCard error={action.error} onClose={() => action.setError(null)} />}
    </div>
  );
}
