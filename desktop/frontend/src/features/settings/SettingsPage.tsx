import { useEffect, useState } from 'react';
import { LANGUAGES, useI18n, type Lang } from '../../i18n';
import { backend, type Settings } from '../../shared/backend';
import { SKINS, THEMES, updateAppearance, useAppearance, type Skin } from '../../shared/appearance';
import { PageHead } from '../../shared/components';
import { LogCleanupCard } from './LogCleanupCard';

/** Miniatura do layout de cada skin (como no game-in). */
function SkinSketch({ skin }: { skin: Skin }) {
  return (
    <span className={`skin-sketch skin-sketch--${skin}`} aria-hidden>
      <span className="skin-sketch__nav" />
      <span className="skin-sketch__main" />
    </span>
  );
}

export function SettingsPage() {
  const { t, lang, setLang } = useI18n();
  const look = useAppearance();
  const [settings, setSettings] = useState<Settings | null>(null);

  useEffect(() => {
    backend
      .settings()
      .then(setSettings)
      .catch(() => {});
  }, []);

  const changeLang = async (l: Lang) => {
    setLang(l);
    try {
      setSettings(await backend.setLanguage(l));
    } catch {
      // fica só nesta sessão
    }
  };

  return (
    <div className="page page--narrow">
      <PageHead title={t('settings.title')} subtitle={t('settings.subtitle')} />

      <section className="card">
        <h2 className="card__title">{t('settings.language')}</h2>
        <div className="chips" role="radiogroup" aria-label={t('settings.language')}>
          {LANGUAGES.map((l) => (
            <button
              key={l.id}
              type="button"
              role="radio"
              aria-checked={lang === l.id}
              className={`chip ${lang === l.id ? 'is-active' : ''}`}
              onClick={() => changeLang(l.id)}
            >
              {l.label}
            </button>
          ))}
        </div>
      </section>

      <section className="card">
        <h2 className="card__title">{t('settings.appearance')}</h2>
        <fieldset className="appearance-group">
          <legend>{t('settings.color')}</legend>
          <div className="appearance-options" role="radiogroup" aria-label={t('settings.color')}>
            {THEMES.map((th) => (
              <button
                key={th.id}
                type="button"
                role="radio"
                aria-checked={look.theme === th.id}
                className={`appearance-option ${look.theme === th.id ? 'is-selected' : ''}`}
                onClick={() => updateAppearance({ theme: th.id })}
              >
                <span className="theme-swatch" aria-hidden>
                  {th.swatch.map((c) => (
                    <span key={c} style={{ background: c }} />
                  ))}
                </span>
                {t(`settings.themes.${th.id}`)}
              </button>
            ))}
          </div>
        </fieldset>
        <fieldset className="appearance-group">
          <legend>{t('settings.skin')}</legend>
          <div className="appearance-options" role="radiogroup" aria-label={t('settings.skin')}>
            {SKINS.map((s) => (
              <button
                key={s}
                type="button"
                role="radio"
                aria-checked={look.skin === s}
                title={t(`settings.skins.${s}Desc`)}
                className={`appearance-option ${look.skin === s ? 'is-selected' : ''}`}
                onClick={() => updateAppearance({ skin: s })}
              >
                <SkinSketch skin={s} />
                {t(`settings.skins.${s}`)}
              </button>
            ))}
          </div>
          <p className="hint">
            {t(`settings.skins.${look.skin}Desc`)} {t('settings.savedHere')}
          </p>
        </fieldset>
      </section>

      <section className="card">
        <label className="check check--block">
          <input
            type="checkbox"
            checked={settings?.showConsole ?? false}
            onChange={async (e) => setSettings(await backend.setShowConsole(e.target.checked))}
          />
          <span>
            {t('settings.console')}
            <span className="field__hint">{t('settings.consoleHint')}</span>
          </span>
        </label>
      </section>

      <section className="card">
        <label className="check check--block">
          <input
            type="checkbox"
            checked={settings?.startWithWindows ?? false}
            onChange={async (e) =>
              setSettings(
                await backend.setAutostart(e.target.checked, settings?.autoStartServer ?? false),
              )
            }
          />
          <span>
            {t('settings.startWithWindows')}
            <span className="field__hint">{t('settings.startWithWindowsHint')}</span>
          </span>
        </label>
        <label className="check check--block">
          <input
            type="checkbox"
            checked={settings?.autoStartServer ?? false}
            disabled={!settings?.startWithWindows}
            onChange={async (e) => setSettings(await backend.setAutostart(true, e.target.checked))}
          />
          <span>
            {t('settings.autoStartServer')}
            <span className="field__hint">{t('settings.autoStartServerHint')}</span>
          </span>
        </label>
      </section>

      <LogCleanupCard />
    </div>
  );
}
