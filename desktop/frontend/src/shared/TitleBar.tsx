import { useEffect, useState } from 'react';
import { useI18n } from '../i18n';
import { backend, windowControls } from './backend';
import { useUpdate } from './useUpdate';
import { AppIcon } from './AppIcon';

/**
 * Barra de título do app (mesmo modelo do midia-multi-tools e do game-in): a barra inteira arrasta a janela,
 * duplo clique maximiza, e os botões de minimizar, maximizar e fechar ficam à direita.
 */
export function TitleBar({ onHelp }: { onHelp: () => void }) {
  const { t } = useI18n();
  const [maximised, setMaximised] = useState(false);
  const update = useUpdate();
  const st = update.state;

  useEffect(() => {
    const check = () =>
      windowControls
        .isMaximised()
        .then(setMaximised)
        .catch(() => {});
    check();
    window.addEventListener('resize', check);
    return () => window.removeEventListener('resize', check);
  }, []);

  return (
    <header className="titlebar">
      <AppIcon className="titlebar__icon" />
      <span className="titlebar__title">{t('app.name')}</span>
      <span className="titlebar__beta">{t('app.beta')}</span>
      <div className="titlebar__actions">
        {update.error !== null && (
          <span className="titlebar__msg" role="alert">
            {t('update.installFailed')}
          </span>
        )}
        {st?.status === 'ready' && (
          <button
            type="button"
            className="titlebar__update"
            title={t('update.readyHint', { version: st.version })}
            onClick={update.install}
          >
            {t('update.install')}
          </button>
        )}
        {st?.status === 'manual' && (
          <button
            type="button"
            className="titlebar__update"
            title={t('update.manualHint', { version: st.version })}
            onClick={() => backend.openLink('download-page')}
          >
            {t('update.manualButton', { version: st.version })}
          </button>
        )}
        <button
          type="button"
          className="titlebar__help"
          title={t('titlebar.help')}
          aria-label={t('titlebar.help')}
          onClick={onHelp}
        >
          ?
        </button>
      </div>
      <div className="titlebar__window">
        <button
          type="button"
          className="titlebar__win"
          title={t('titlebar.minimize')}
          aria-label={t('titlebar.minimize')}
          onClick={windowControls.minimise}
        >
          <svg viewBox="0 0 10 10" aria-hidden="true">
            <path d="M0 5h10" />
          </svg>
        </button>
        <button
          type="button"
          className="titlebar__win"
          title={maximised ? t('titlebar.restore') : t('titlebar.maximize')}
          aria-label={maximised ? t('titlebar.restore') : t('titlebar.maximize')}
          onClick={windowControls.toggleMaximise}
        >
          {maximised ? (
            <svg viewBox="0 0 10 10" aria-hidden="true">
              <path d="M2.5 2.5V.5h7v7h-2M.5 2.5h7v7h-7z" />
            </svg>
          ) : (
            <svg viewBox="0 0 10 10" aria-hidden="true">
              <path d="M.5.5h9v9h-9z" />
            </svg>
          )}
        </button>
        <button
          type="button"
          className="titlebar__win titlebar__win--close"
          title={t('titlebar.close')}
          aria-label={t('titlebar.close')}
          onClick={windowControls.close}
        >
          <svg viewBox="0 0 10 10" aria-hidden="true">
            <path d="M0 0l10 10M10 0L0 10" />
          </svg>
        </button>
      </div>
    </header>
  );
}
