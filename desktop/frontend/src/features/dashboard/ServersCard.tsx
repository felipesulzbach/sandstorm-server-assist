import { Eye, Loader2, Play, Plus, Square } from 'lucide-react';
import type { PageId } from '../../App';
import { useI18n } from '../../i18n';
import { backend, type ServerSummary } from '../../shared/backend';
import { ErrorCard } from '../../shared/components';
import { useAction } from '../../shared/hooks';
import { useServers } from '../../shared/servers';

// Fases em que o servidor está a caminho (botão com roda a girar).
const BUSY = ['waiting', 'starting', 'modio', 'mods', 'traveling', 'stopping'];
// Fases em que há o que parar (processo rodando ou esperando a vez).
const STOPPABLE = ['waiting', 'starting', 'modio', 'mods', 'traveling', 'ready'];

/**
 * Todos os servidores do PC (ADR-048), com o estado de cada um e Iniciar/Parar. "Mostrar" escolhe o servidor que as
 * outras telas mostram. Só aparece com mais de um servidor; o primeiro outro se cria em Servidor › Novo servidor.
 */
export function ServersCard({ go }: { go: (p: PageId) => void }) {
  const { t } = useI18n();
  const { list, selected, select, checked, toggle, setChecked } = useServers();
  const action = useAction();
  if (list.length < 2) return null;
  // Os que dá para iniciar ou parar pelo botão principal (instalados).
  const pickable = list.filter((s) => s.phase !== 'not-installed').map((s) => s.id);
  const allPicked = pickable.length > 0 && pickable.every((id) => checked.includes(id));

  const row = (s: ServerSummary) => {
    const stoppable = STOPPABLE.includes(s.phase);
    return (
      <li key={s.id} className={`list__row servers__row ${s.id === selected ? 'is-selected' : ''}`}>
        <input
          type="checkbox"
          className="servers__check"
          aria-label={t('servers.pick', { name: s.name })}
          title={t('servers.pick', { name: s.name })}
          checked={checked.includes(s.id)}
          disabled={s.phase === 'not-installed'}
          onChange={(e) => toggle(s.id, e.target.checked)}
        />
        <span className="servers__state">
          {BUSY.includes(s.phase) ? (
            <Loader2 className="spin" size={16} aria-hidden />
          ) : (
            <span className={`dot dot--${s.phase}`} aria-hidden />
          )}
          <span>{t(`phase.${s.phase}`)}</span>
        </span>
        <span className="servers__name">
          <strong>{s.name}</strong>
          {s.hostname && s.hostname !== s.name && <span className="muted">{s.hostname}</span>}
          <span className="muted servers__meta">
            {t('servers.port', { port: s.port })}
            {s.map && ` · ${s.map}`}
            {s.external && ` · ${t('servers.external')}`}
          </span>
          {(s.conflicts ?? []).slice(0, 1).map((c) => (
            <span key={c.field} className="hint hint--warn">
              {t('servers.conflict', { port: c.port, server: c.server })}
            </span>
          ))}
        </span>
        <span className="row">
          {stoppable ? (
            <button
              type="button"
              className="btn btn--small"
              disabled={action.busy || s.phase === 'stopping'}
              onClick={() => action.run(() => backend.stopServer(s.id))}
            >
              <Square size={14} aria-hidden /> {t('servers.stop')}
            </button>
          ) : (
            <button
              type="button"
              className="btn btn--small btn--primary"
              disabled={action.busy || s.phase === 'not-installed' || s.phase === 'stopping'}
              onClick={() => action.run(() => backend.startServer(s.id))}
            >
              <Play size={14} aria-hidden /> {t('servers.start')}
            </button>
          )}
          <button
            type="button"
            className="btn btn--small btn--ghost"
            disabled={action.busy || s.id === selected}
            onClick={() => action.run(() => select(s.id))}
          >
            <Eye size={14} aria-hidden />{' '}
            {s.id === selected ? t('servers.shown') : t('servers.show')}
          </button>
        </span>
      </li>
    );
  };

  return (
    <section className="card servers">
      <div className="card__head">
        <h2>{t('servers.title')}</h2>
        <button type="button" className="btn btn--small btn--ghost" onClick={() => go('server')}>
          <Plus size={14} aria-hidden /> {t('servers.add')}
        </button>
      </div>
      <p className="muted">{t('servers.subtitle')}</p>
      <div className="row row--wrap servers__pickrow">
        <label className="check">
          <input
            type="checkbox"
            className="servers__check"
            checked={allPicked}
            ref={(el) => {
              if (el) el.indeterminate = checked.length > 0 && !allPicked;
            }}
            onChange={(e) => setChecked(e.target.checked ? pickable : [])}
          />
          <span>{t('servers.selectAll')}</span>
        </label>
        <span className="hint">{t('servers.pickHint')}</span>
      </div>
      {action.error && <ErrorCard error={action.error} onClose={() => action.setError(null)} />}
      <ul className="list">{list.map(row)}</ul>
    </section>
  );
}
