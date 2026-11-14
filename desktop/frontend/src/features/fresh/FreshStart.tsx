import { useEffect, useState } from 'react';
import { Eraser } from 'lucide-react';
import { useI18n } from '../../i18n';
import { backend, type Fresh, type Profile } from '../../shared/backend';
import { ErrorCard } from '../../shared/components';
import { useAction } from '../../shared/hooks';

/**
 * Começar do zero (ADR-047). Quando a pasta do servidor sumiu ou não tem servidor instalado e o perfil ainda tem os
 * dados do servidor antigo, oferece voltar o perfil ao padrão (pastas mantidas) e, se o usuário quiser, esvaziar o
 * cofre. O app nunca apaga sozinho: a pasta pode ter sumido por um disco desligado. Quem decide é o Go (FreshState).
 */
export function FreshStart({
  checkKey,
  onDone,
}: {
  /** Muda quando a pasta ou a instalação mudam: a tela pergunta de novo ao Go. */
  checkKey: string;
  onDone: (p: Profile) => void;
}) {
  const { t } = useI18n();
  const action = useAction();
  const [state, setState] = useState<Fresh | null>(null);
  const [hidden, setHidden] = useState(false);
  const [confirming, setConfirming] = useState(false);
  const [clearSecrets, setClearSecrets] = useState(true);
  // Mods e autorização do mod.io ficam no Windows, fora da pasta do servidor: desmarcados por padrão.
  const [clearMods, setClearMods] = useState(false);
  const [logout, setLogout] = useState(false);
  const [done, setDone] = useState(false);

  useEffect(() => {
    backend
      .freshState()
      .then(setState)
      .catch(() => setState(null));
  }, [checkKey]);

  if (done) {
    return (
      <p className="banner banner--ok" role="status">
        {t('fresh.done')}
      </p>
    );
  }
  if (!state?.offer || hidden) return null;

  const start = async () => {
    const p = await action.run(async () => {
      const profile = await backend.startFresh(state.hasSecrets && clearSecrets);
      const mods = state.mods > 0 && clearMods;
      const session = state.modioSession && logout;
      if (mods || session) await backend.clearModio(mods, session);
      return profile;
    });
    if (!p) return;
    setDone(true);
    setConfirming(false);
    onDone(p);
  };

  return (
    <div className="banner banner--warn fresh" role="status">
      <p>
        <strong>{state.missing ? t('fresh.missing') : t('fresh.notInstalled')}</strong>{' '}
        {t('fresh.text')}
      </p>
      {action.error && <ErrorCard error={action.error} onClose={() => action.setError(null)} />}
      {confirming ? (
        <>
          <p>{t('fresh.confirm')}</p>
          {state.hasSecrets && (
            <label className="check">
              <input
                type="checkbox"
                checked={clearSecrets}
                onChange={(e) => setClearSecrets(e.target.checked)}
              />
              <span>
                {t('fresh.clearSecrets')}
                {state.others > 0 && (
                  <span className="field__hint">
                    {t('fresh.sharedSecrets', { n: state.others })}
                  </span>
                )}
              </span>
            </label>
          )}
          {state.mods > 0 && (
            <label className="check">
              <input
                type="checkbox"
                checked={clearMods}
                onChange={(e) => setClearMods(e.target.checked)}
              />
              <span>
                {t('fresh.clearMods', { n: state.mods })}
                <span className="field__hint">{t('fresh.clearModsHint')}</span>
              </span>
            </label>
          )}
          {state.modioSession && (
            <label className="check">
              <input
                type="checkbox"
                checked={logout}
                onChange={(e) => setLogout(e.target.checked)}
              />
              <span>
                {t('fresh.logout')}
                <span className="field__hint">{t('fresh.logoutHint')}</span>
              </span>
            </label>
          )}
          <div className="row row--wrap">
            <button
              type="button"
              className="btn btn--danger btn--small"
              disabled={action.busy}
              onClick={start}
            >
              {t('fresh.confirmYes')}
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
        </>
      ) : (
        <div className="row row--wrap">
          <button type="button" className="btn btn--small" onClick={() => setConfirming(true)}>
            <Eraser size={14} aria-hidden /> {t('fresh.start')}
          </button>
          <button
            type="button"
            className="btn btn--small btn--ghost"
            onClick={() => setHidden(true)}
          >
            {t('fresh.keep')}
          </button>
        </div>
      )}
    </div>
  );
}
