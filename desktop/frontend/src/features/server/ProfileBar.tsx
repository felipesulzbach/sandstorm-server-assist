import { useState } from 'react';
import { Copy, Pencil, Plus, Trash2 } from 'lucide-react';
import { useI18n } from '../../i18n';
import { backend } from '../../shared/backend';
import { ErrorCard } from '../../shared/components';
import { useAction } from '../../shared/hooks';
import { useServers } from '../../shared/servers';

// Fases em que o servidor não pode ser apagado (processo rodando ou esperando a vez).
const ACTIVE = ['waiting', 'starting', 'modio', 'mods', 'traveling', 'ready', 'stopping'];

/**
 * Servidores (ADR-048, antes perfis do ADR-022): o escolhido é o que esta tela edita. Cada servidor tem as próprias
 * configurações, senhas e portas e pode rodar junto dos outros; a conta e os mods do mod.io são os mesmos. Um
 * servidor novo ganha portas livres e passa a ser o escolhido.
 */
export function ProfileBar({ dirty }: { dirty: boolean }) {
  const { t } = useI18n();
  const action = useAction();
  const { list, selected, select, refresh } = useServers();
  const [mode, setMode] = useState<'new' | 'rename' | null>(null);
  const [name, setName] = useState('');
  const [copy, setCopy] = useState(true);
  const [confirmDelete, setConfirmDelete] = useState<string | null>(null);

  const active = list.find((p) => p.id === selected);

  const submit = async () => {
    const r = await action.run(() =>
      mode === 'rename' && active
        ? backend.renameProfile(active.id, name)
        : backend.createProfile(name, copy),
    );
    if (r) {
      setMode(null);
      setName('');
      // Servidor novo: passa a ser o escolhido e as telas remontam com ele.
      await refresh();
    }
  };

  const remove = async (id: string) => {
    const r = await action.run(() => backend.deleteProfile(id));
    setConfirmDelete(null);
    if (r) await refresh();
  };

  return (
    <div className="profilebar">
      <div className="row row--wrap">
        <label className="profilebar__pick">
          <span className="field__label">{t('profiles.label')}</span>
          <select
            className="select"
            value={selected}
            disabled={dirty || action.busy}
            title={dirty ? t('profiles.saveFirst') : undefined}
            onChange={(e) => action.run(() => select(e.target.value))}
          >
            {list.map((p) => (
              <option key={p.id} value={p.id}>
                {p.name}
                {p.hostname && p.hostname !== p.name ? ` (${p.hostname})` : ''}
              </option>
            ))}
          </select>
        </label>
        <button
          type="button"
          className="btn btn--small btn--ghost"
          disabled={dirty}
          title={dirty ? t('profiles.saveFirst') : undefined}
          onClick={() => {
            setMode('new');
            setName('');
          }}
        >
          <Plus size={14} aria-hidden /> {t('profiles.new')}
        </button>
        <button
          type="button"
          className="btn btn--small btn--ghost"
          disabled={!active}
          onClick={() => {
            setMode('rename');
            setName(active?.name ?? '');
          }}
        >
          <Pencil size={14} aria-hidden /> {t('profiles.rename')}
        </button>
        <span className="hint">{t('profiles.switchHint')}</span>
      </div>

      {mode && (
        <div className="row row--wrap profilebar__form">
          <input
            className="input"
            value={name}
            maxLength={40}
            placeholder={t('profiles.namePlaceholder')}
            onChange={(e) => setName(e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && submit()}
          />
          {mode === 'new' && (
            <label className="check">
              <input type="checkbox" checked={copy} onChange={(e) => setCopy(e.target.checked)} />
              <span>
                <Copy size={14} aria-hidden /> {t('profiles.copyCurrent')}
              </span>
            </label>
          )}
          <button
            type="button"
            className="btn btn--small btn--primary"
            disabled={!name.trim() || action.busy}
            onClick={submit}
          >
            {mode === 'new' ? t('profiles.create') : t('common.save')}
          </button>
          <button type="button" className="btn btn--small" onClick={() => setMode(null)}>
            {t('common.cancel')}
          </button>
          {mode === 'new' && <span className="hint">{t('profiles.newHint')}</span>}
        </div>
      )}

      {list.length > 1 && (
        <ul className="list profilebar__others">
          {list
            .filter((p) => p.id !== selected)
            .map((p) => (
              <li key={p.id} className="list__row">
                <span>{p.name}</span>
                <span className="muted">
                  {p.hostname} · {t('servers.port', { port: p.port })}
                </span>
                {confirmDelete === p.id ? (
                  <span className="row">
                    <button
                      type="button"
                      className="btn btn--small btn--danger"
                      onClick={() => remove(p.id)}
                    >
                      {t('profiles.deleteYes')}
                    </button>
                    <button
                      type="button"
                      className="btn btn--small"
                      onClick={() => setConfirmDelete(null)}
                    >
                      {t('common.cancel')}
                    </button>
                  </span>
                ) : (
                  <button
                    type="button"
                    className="icon-btn"
                    aria-label={t('profiles.delete')}
                    title={
                      ACTIVE.includes(p.phase) ? t('profiles.deleteRunning') : t('profiles.delete')
                    }
                    disabled={ACTIVE.includes(p.phase)}
                    onClick={() => setConfirmDelete(p.id)}
                  >
                    <Trash2 size={16} />
                  </button>
                )}
              </li>
            ))}
        </ul>
      )}

      {action.error && <ErrorCard error={action.error} onClose={() => action.setError(null)} />}
    </div>
  );
}
