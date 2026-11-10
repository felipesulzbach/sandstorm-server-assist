import { useEffect, useState } from 'react';
import { Copy, Pencil, Plus, Trash2 } from 'lucide-react';
import { useI18n } from '../../i18n';
import { backend, type Profile, type ProfileRef } from '../../shared/backend';
import { ErrorCard } from '../../shared/components';
import { useAction } from '../../shared/hooks';

/**
 * Perfis do servidor (ADR-022): um perfil em uso e outros guardados (ex.: Coop e Versus). Trocar de perfil troca
 * todas as configurações do servidor; os segredos do cofre são os mesmos. Só com o servidor parado.
 */
export function ProfileBar({
  dirty,
  onSwitched,
}: {
  dirty: boolean;
  onSwitched: (p: Profile) => void;
}) {
  const { t } = useI18n();
  const action = useAction();
  const [list, setList] = useState<ProfileRef[]>([]);
  const [mode, setMode] = useState<'new' | 'rename' | null>(null);
  const [name, setName] = useState('');
  const [copy, setCopy] = useState(true);
  const [confirmDelete, setConfirmDelete] = useState<string | null>(null);

  useEffect(() => {
    backend
      .profiles()
      .then(setList)
      .catch(() => {});
  }, []);

  const active = list.find((p) => p.active);

  const submit = async () => {
    const r = await action.run(() =>
      mode === 'rename' && active
        ? backend.renameProfile(active.id, name)
        : backend.createProfile(name, copy),
    );
    if (r) {
      setList(r);
      setMode(null);
      setName('');
    }
  };

  const switchTo = async (id: string) => {
    const p = await action.run(() => backend.switchProfile(id));
    if (p) {
      onSwitched(p);
      setList(await backend.profiles());
    }
  };

  const remove = async (id: string) => {
    const r = await action.run(() => backend.deleteProfile(id));
    if (r) setList(r);
    setConfirmDelete(null);
  };

  return (
    <div className="profilebar">
      <div className="row row--wrap">
        <label className="profilebar__pick">
          <span className="field__label">{t('profiles.label')}</span>
          <select
            className="select"
            value={active?.id ?? ''}
            disabled={dirty || action.busy}
            title={dirty ? t('profiles.saveFirst') : undefined}
            onChange={(e) => switchTo(e.target.value)}
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
        {list.length > 1 && <span className="hint">{t('profiles.switchHint')}</span>}
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
        </div>
      )}

      {list.length > 1 && (
        <ul className="list profilebar__others">
          {list
            .filter((p) => !p.active)
            .map((p) => (
              <li key={p.id} className="list__row">
                <span>{p.name}</span>
                <span className="muted">{p.hostname}</span>
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
                    title={t('profiles.delete')}
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
