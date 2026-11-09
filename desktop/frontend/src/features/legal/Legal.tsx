import { useEffect, useState } from 'react';
import { FileText, ShieldCheck } from 'lucide-react';
import { LANGUAGES, isLang, useI18n } from '../../i18n';
import { backend, windowControls, type Terms } from '../../shared/backend';
import { ErrorCard } from '../../shared/components';
import { toAppError, type AppError } from '../../shared/problem';

// Seções dos textos (ADR-015). Os textos ficam em legal.terms.* e legal.privacy.* nos 5 idiomas.
const TERMS = [
  'about',
  'asIs',
  'responsibility',
  'prohibited',
  'thirdParty',
  'updates',
  'changes',
] as const;
const PRIVACY = ['summary', 'stored', 'sent', 'server', 'control', 'changes'] as const;

type Tab = 'terms' | 'privacy';

/** Texto dos termos de uso ou do aviso de dados. */
function LegalText({ tab }: { tab: Tab }) {
  const { t, tl } = useI18n();
  const ids: readonly string[] = tab === 'terms' ? TERMS : PRIVACY;
  return (
    <div className="legal__text">
      {ids.map((id) => (
        <section key={id}>
          <h3>{t(`legal.${tab}.${id}.title`)}</h3>
          {tl(`legal.${tab}.${id}.text`).map((p) => (
            <p key={p}>{p}</p>
          ))}
        </section>
      ))}
    </div>
  );
}

function Tabs({ tab, setTab }: { tab: Tab; setTab: (t: Tab) => void }) {
  const { t } = useI18n();
  return (
    <div className="tabs" role="tablist">
      {(['terms', 'privacy'] as const).map((id) => (
        <button
          key={id}
          type="button"
          role="tab"
          aria-selected={tab === id}
          className={`tabs__tab ${tab === id ? 'is-active' : ''}`}
          onClick={() => setTab(id)}
        >
          {id === 'terms' ? (
            <FileText size={14} aria-hidden />
          ) : (
            <ShieldCheck size={14} aria-hidden />
          )}{' '}
          {t(`legal.tabs.${id}`)}
        </button>
      ))}
    </div>
  );
}

/**
 * Aceite dos termos de uso e do aviso de dados (clickwrap, como no game-in): sem aceitar a versão em vigor, o
 * app fica coberto por esta tela. A barra de título continua livre (fechar o app funciona).
 */
export function LegalGate({ terms, onAccepted }: { terms: Terms; onAccepted: (t: Terms) => void }) {
  const { t, lang, setLang } = useI18n();
  const [tab, setTab] = useState<Tab>('terms');
  const [agree, setAgree] = useState(false);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<AppError | null>(null);

  const changeLang = (v: string) => {
    if (!isLang(v)) return;
    setLang(v);
    backend.setLanguage(v).catch(() => {});
  };

  const accept = async () => {
    setBusy(true);
    setError(null);
    try {
      onAccepted(await backend.acceptTerms(terms.version));
    } catch (e) {
      setError(toAppError(e));
    } finally {
      setBusy(false);
    }
  };

  return (
    <div className="legal-gate" role="dialog" aria-modal="true" aria-labelledby="legal-title">
      <div className="modal legal">
        <div className="modal__head">
          <div>
            <h2 id="legal-title">{t('legal.title')}</h2>
            <p className="muted">{t('legal.intro')}</p>
          </div>
          <select
            className="input legal__lang"
            aria-label={t('settings.language')}
            value={lang}
            onChange={(e) => changeLang(e.target.value)}
          >
            {LANGUAGES.map((l) => (
              <option key={l.id} value={l.id}>
                {l.label}
              </option>
            ))}
          </select>
        </div>
        <Tabs tab={tab} setTab={setTab} />
        <div className="modal__body">
          <LegalText tab={tab} />
        </div>
        <div className="modal__foot legal__foot">
          {error && <ErrorCard error={error} onClose={() => setError(null)} />}
          <label className="check check--block">
            <input type="checkbox" checked={agree} onChange={(e) => setAgree(e.target.checked)} />
            <span>{t('legal.agree')}</span>
          </label>
          <div className="row">
            <button type="button" className="btn" onClick={windowControls.close}>
              {t('legal.exit')}
            </button>
            <button
              type="button"
              className="btn btn--primary"
              disabled={!agree || busy}
              onClick={accept}
            >
              {t('legal.accept')}
            </button>
          </div>
          <p className="hint">{t('legal.version', { version: terms.version })}</p>
        </div>
      </div>
    </div>
  );
}

/** Os mesmos textos só para ler (Sobre), já aceitos. */
export function LegalDialog({ initial, onClose }: { initial: Tab; onClose: () => void }) {
  const { t } = useI18n();
  const [tab, setTab] = useState<Tab>(initial);

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => e.key === 'Escape' && onClose();
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
  }, [onClose]);

  return (
    <div
      className="modal-backdrop"
      role="dialog"
      aria-modal="true"
      aria-labelledby="legal-read-title"
    >
      <div className="modal legal">
        <div className="modal__head">
          <h2 id="legal-read-title">{t('legal.title')}</h2>
        </div>
        <Tabs tab={tab} setTab={setTab} />
        <div className="modal__body">
          <LegalText tab={tab} />
        </div>
        <div className="modal__foot">
          <button type="button" className="btn" onClick={onClose}>
            {t('common.close')}
          </button>
        </div>
      </div>
    </div>
  );
}

/** Carrega a situação do aceite; null enquanto carrega. */
export function useTerms() {
  const [terms, setTerms] = useState<Terms | null>(null);
  useEffect(() => {
    backend
      .terms()
      .then(setTerms)
      // Sem resposta do Go não dá para gravar o aceite: mostra a tela de aceite mesmo assim.
      .catch(() => setTerms({ version: '', accepted: false, acceptedAt: '' }));
  }, []);
  return [terms, setTerms] as const;
}
