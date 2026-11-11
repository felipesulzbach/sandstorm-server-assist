import { useState, type ReactNode } from 'react';
import { AlertTriangle, ChevronDown, ExternalLink, HelpCircle, X } from 'lucide-react';
import { useI18n } from '../i18n';
import { backend, type SecretKey } from './backend';
import { toAppError, type AppError } from './problem';
import { percent } from './format';

/** Erro do catálogo: o que aconteceu, possível causa e como resolver (ADR-005). */
export function ErrorCard({
  error,
  onClose,
  onOpenLogs,
}: {
  error: AppError;
  onClose?: () => void;
  onOpenLogs?: () => void;
}) {
  const { t, tl } = useI18n();
  const base = `problems.${error.code}`;
  const title = t(`${base}.title`);
  const known = title !== `${base}.title`;
  const p = known ? base : 'problems.unknown';
  const fieldMsg = error.msg ? t(error.msg, error.params) : '';
  return (
    <div className="errorcard" role="alert">
      <div className="errorcard__head">
        <AlertTriangle size={18} aria-hidden />
        <strong>{t(`${p}.title`)}</strong>
        {onClose && (
          <button
            type="button"
            className="icon-btn"
            aria-label={t('common.close')}
            onClick={onClose}
          >
            <X size={16} />
          </button>
        )}
      </div>
      {fieldMsg && <p className="errorcard__field">{fieldMsg}</p>}
      {!fieldMsg && error.params?.object && (
        <p className="errorcard__field">
          {t(error.params.mods ? 'error.modOwners' : 'error.modFile', error.params)}
        </p>
      )}
      {!fieldMsg && (
        <dl className="errorcard__body">
          <dt>{t('error.what')}</dt>
          <dd>{t(`${p}.what`)}</dd>
          <dt>{t('error.cause')}</dt>
          <dd>{t(`${p}.cause`)}</dd>
          <dt>{t('error.solution')}</dt>
          <dd>
            <ol>
              {tl(`${p}.solution`).map((s) => (
                <li key={s}>{s}</li>
              ))}
            </ol>
          </dd>
        </dl>
      )}
      <div className="errorcard__actions">
        {(error.links ?? []).map((k) => (
          <LinkButton key={k} linkKey={k} />
        ))}
        {onOpenLogs && (
          <button type="button" className="btn btn--small" onClick={onOpenLogs}>
            {t('error.openLogs')}
          </button>
        )}
      </div>
      {error.detail && (
        <details className="errorcard__details">
          <summary>{t('error.details')}</summary>
          <code>{error.detail}</code>
        </details>
      )}
    </div>
  );
}

/** Botão que abre um link conhecido (o Go tem a lista; a tela manda só a chave). */
export function LinkButton({ linkKey, label }: { linkKey: string; label?: string }) {
  const { t } = useI18n();
  return (
    <button
      type="button"
      className="btn btn--small btn--ghost"
      onClick={() => backend.openLink(linkKey)}
    >
      <ExternalLink size={14} aria-hidden />
      {label ?? t(`links.${linkKey}`)}
    </button>
  );
}

/** Guia "Como conseguir?": passos e links, fechado por padrão. */
/** Guia passo a passo (guides.<id>). label troca o "Como fazer?" quando há mais de um guia lado a lado. */
export function Guide({
  id,
  links,
  label,
  params,
}: {
  id: string;
  links?: string[];
  label?: string;
  params?: Record<string, string | number>;
}) {
  const { t, tl } = useI18n();
  const [open, setOpen] = useState(false);
  return (
    <div className={`guide ${open ? 'is-open' : ''}`}>
      <button
        type="button"
        className="guide__toggle"
        aria-expanded={open}
        onClick={() => setOpen(!open)}
      >
        <HelpCircle size={15} aria-hidden />
        {label ?? t('common.howTo')}
        <ChevronDown size={15} className="guide__chev" aria-hidden />
      </button>
      {open && (
        <div className="guide__body">
          <strong>{t(`guides.${id}.title`)}</strong>
          <ol>
            {tl(`guides.${id}.steps`, params).map((s) => (
              <li key={s}>{s}</li>
            ))}
          </ol>
          {links && links.length > 0 && (
            <div className="guide__links">
              {links.map((k) => (
                <LinkButton key={k} linkKey={k} />
              ))}
            </div>
          )}
        </div>
      )}
    </div>
  );
}

/** Campo de formulário com rótulo, dica e erro. */
export function Field({
  label,
  hint,
  error,
  children,
  wide,
}: {
  label: string;
  hint?: ReactNode;
  error?: string | null;
  children: ReactNode;
  wide?: boolean;
}) {
  return (
    <label className={`field ${wide ? 'field--wide' : ''} ${error ? 'has-error' : ''}`}>
      <span className="field__label">{label}</span>
      {children}
      {error ? (
        <span className="field__error">{error}</span>
      ) : hint ? (
        <span className="field__hint">{hint}</span>
      ) : null}
    </label>
  );
}

/**
 * Segredo guardado no cofre: mostra só "definido/não definido". "Mostrar" pede o valor ao Go; "Trocar" abre um
 * campo; nada fica no estado da tela depois de salvar.
 */
export function SecretField({
  name,
  label,
  hint,
  isSet,
  onChanged,
  plain,
  children,
}: {
  name: SecretKey;
  label: string;
  hint?: ReactNode;
  isSet: boolean;
  onChanged: (state: Record<SecretKey, boolean>) => void;
  /** Campo que não é senha (e-mail): mostra o valor em texto e sem botão "Mostrar". */
  plain?: boolean;
  children?: ReactNode;
}) {
  const { t } = useI18n();
  const [editing, setEditing] = useState(false);
  const [value, setValue] = useState('');
  const [shown, setShown] = useState<string | null>(null);
  const [err, setErr] = useState<string | null>(null);

  const save = async (v: string) => {
    setErr(null);
    try {
      const st = await backend.setSecret(name, v);
      setEditing(false);
      setValue('');
      setShown(null);
      onChanged(st);
    } catch (e) {
      const pe = toAppError(e);
      setErr(pe.msg ? t(pe.msg, pe.params) : t(`problems.${pe.code}.title`));
    }
  };

  const reveal = async () => {
    if (shown !== null) {
      setShown(null);
      return;
    }
    try {
      setShown(await backend.revealSecret(name));
    } catch {
      setShown(null);
    }
  };

  return (
    <div className="secret">
      <div className="secret__head">
        <span className="field__label">{label}</span>
        <span className={`pill ${isSet ? 'pill--ok' : ''}`}>
          {isSet ? t('server.secrets.set') : t('server.secrets.notSet')}
        </span>
      </div>
      {editing ? (
        <div className="secret__edit">
          <input
            className="input"
            type={plain ? 'text' : 'password'}
            autoComplete="off"
            spellCheck={false}
            value={value}
            aria-label={label}
            onChange={(e) => setValue(e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && save(value)}
            autoFocus
          />
          <button type="button" className="btn btn--primary btn--small" onClick={() => save(value)}>
            {t('common.save')}
          </button>
          <button
            type="button"
            className="btn btn--small"
            onClick={() => {
              setEditing(false);
              setValue('');
              setErr(null);
            }}
          >
            {t('common.cancel')}
          </button>
        </div>
      ) : (
        <div className="secret__row">
          <code className="secret__value">{isSet ? (shown ?? '••••••••') : '—'}</code>
          {isSet && (
            <button type="button" className="btn btn--small btn--ghost" onClick={reveal}>
              {shown !== null ? t('common.hide') : t('common.show')}
            </button>
          )}
          <button type="button" className="btn btn--small" onClick={() => setEditing(true)}>
            {isSet ? t('common.change') : t('common.set')}
          </button>
          {isSet && (
            <button type="button" className="btn btn--small btn--ghost" onClick={() => save('')}>
              {t('common.remove')}
            </button>
          )}
        </div>
      )}
      {err ? (
        <span className="field__error">{err}</span>
      ) : hint ? (
        <span className="field__hint">{hint}</span>
      ) : null}
      {children}
    </div>
  );
}

/** Barra de progresso; sem total, fica indeterminada (animada). */
export function ProgressBar({
  done,
  total,
  value,
  label,
}: {
  done?: number;
  total?: number;
  value?: number;
  label?: string;
}) {
  const pct = value ?? (total ? percent(done ?? 0, total) : undefined);
  return (
    <div
      className={`progress ${pct === undefined ? 'progress--indeterminate' : ''}`}
      role="progressbar"
      aria-label={label}
      aria-valuemin={0}
      aria-valuemax={100}
      aria-valuenow={pct === undefined ? undefined : Math.round(pct)}
    >
      <div className="progress__bar" style={pct === undefined ? undefined : { width: `${pct}%` }} />
    </div>
  );
}

/** Cabeçalho de página: título e subtítulo, como no midia-multi-tools. */
export function PageHead({
  title,
  subtitle,
  children,
}: {
  title: string;
  subtitle?: string;
  children?: ReactNode;
}) {
  return (
    <header className="page__head">
      <div>
        <h1>{title}</h1>
        {subtitle && <p className="muted">{subtitle}</p>}
      </div>
      {children && <div className="page__actions">{children}</div>}
    </header>
  );
}

/** Botão de copiar com "Copiado!" por 2 s. */
export function CopyButton({ text }: { text: string }) {
  const { t } = useI18n();
  const [done, setDone] = useState(false);
  return (
    <button
      type="button"
      className="btn btn--small"
      onClick={async () => {
        await backend.copyText(text);
        setDone(true);
        setTimeout(() => setDone(false), 2000);
      }}
    >
      {done ? t('common.copied') : t('common.copy')}
    </button>
  );
}
