import { useEffect, useRef, useState } from 'react';
import { Bold, Heading1, Heading2, ImageIcon } from 'lucide-react';
import { useI18n } from '../../i18n';
import { backend, type Segment } from '../../shared/backend';
import { toAppError } from '../../shared/problem';

type Tag = 'Bold' | 'Header1' | 'Header2';

/**
 * Editor das mensagens com formatação do jogo (ADR-024): botões que envolvem o texto selecionado com as tags do
 * texto rico do Unreal (<Bold>…</>, <Header1>…</>, <Header2>…</>, <extimg src="…"/>) e uma prévia. Quem confere as
 * tags e monta a prévia é o Go; a tela só edita o texto e desenha os pedaços. A imagem não é baixada pelo app.
 */
export function MarkupEditor({
  value,
  onChange,
  chat = false,
  multiline = true,
  rows = 5,
  maxLength,
  placeholder,
  preview = true,
}: {
  value: string;
  onChange: (v: string) => void;
  chat?: boolean;
  multiline?: boolean;
  rows?: number;
  maxLength?: number;
  placeholder?: string;
  preview?: boolean;
}) {
  const { t } = useI18n();
  const ref = useRef<HTMLTextAreaElement & HTMLInputElement>(null);
  const [segments, setSegments] = useState<Segment[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [imageUrl, setImageUrl] = useState<string | null>(null);

  // Prévia conferida pelo Go, um pouco depois de parar de digitar.
  useEffect(() => {
    if (!preview) return;
    const id = setTimeout(() => {
      backend
        .markupPreview(value, chat)
        .then((s) => {
          setSegments(s ?? []);
          setError(null);
        })
        .catch((e) => {
          const pe = toAppError(e);
          setError(pe.msg ? t(pe.msg, pe.params) : t(`problems.${pe.code}.title`));
        });
    }, 300);
    return () => clearTimeout(id);
  }, [value, chat, preview, t]);

  /** Troca a seleção por before + seleção + after e devolve o foco com a seleção no texto envolvido. */
  const insert = (before: string, after = '') => {
    const el = ref.current;
    const start = el?.selectionStart ?? value.length;
    const end = el?.selectionEnd ?? value.length;
    const next =
      value.slice(0, start) + before + value.slice(start, end) + after + value.slice(end);
    onChange(next);
    requestAnimationFrame(() => {
      if (!el) return;
      el.focus();
      el.setSelectionRange(start + before.length, end + before.length);
    });
  };
  const wrap = (tag: Tag) => insert(`<${tag}>`, '</>');

  const tools: { tag: Tag; icon: typeof Bold; label: string }[] = [
    { tag: 'Bold', icon: Bold, label: t('markup.bold') },
    ...(chat
      ? []
      : [
          { tag: 'Header1' as Tag, icon: Heading1, label: t('markup.header1') },
          { tag: 'Header2' as Tag, icon: Heading2, label: t('markup.header2') },
        ]),
  ];

  const common = {
    ref,
    className: 'input markup__input',
    value,
    maxLength,
    placeholder,
    spellCheck: true,
    onChange: (e: { target: { value: string } }) => onChange(e.target.value),
  };

  return (
    <div className="markup">
      <div className="markup__toolbar" role="toolbar" aria-label={t('markup.toolbar')}>
        {tools.map(({ tag, icon: Icon, label }) => (
          <button
            key={tag}
            type="button"
            className="icon-btn"
            title={label}
            aria-label={label}
            onMouseDown={(e) => e.preventDefault()}
            onClick={() => wrap(tag)}
          >
            <Icon size={16} />
          </button>
        ))}
        {!chat && (
          <button
            type="button"
            className="icon-btn"
            title={t('markup.image')}
            aria-label={t('markup.image')}
            onMouseDown={(e) => e.preventDefault()}
            onClick={() => setImageUrl((u) => (u === null ? '' : null))}
          >
            <ImageIcon size={16} />
          </button>
        )}
        {chat && <span className="hint">{t('markup.chatNote')}</span>}
      </div>

      {imageUrl !== null && (
        <div className="row row--wrap markup__image">
          <input
            className="input"
            value={imageUrl}
            placeholder="https://…/image.png"
            spellCheck={false}
            onChange={(e) => setImageUrl(e.target.value.trim())}
          />
          <button
            type="button"
            className="btn btn--small btn--primary"
            disabled={!imageUrl.startsWith('https://')}
            onClick={() => {
              insert(`<extimg src="${imageUrl.replace(/"/g, '')}"/>`);
              setImageUrl(null);
            }}
          >
            {t('markup.insert')}
          </button>
          <span className="hint">{t('markup.imageNote')}</span>
        </div>
      )}

      {multiline ? <textarea {...common} rows={rows} /> : <input {...common} />}

      {preview && (value.trim() !== '' || error) && (
        <div className={`markup__preview ${error ? 'has-error' : ''}`} aria-live="polite">
          <span className="markup__previewLabel">{t('markup.preview')}</span>
          {error ? (
            <span className="field__error">{error}</span>
          ) : (
            <div className="markup__text">
              {segments.map((s, i) =>
                s.src ? (
                  <span key={i} className="markup__img" title={s.src}>
                    <ImageIcon size={14} aria-hidden /> {t('markup.imageAt', { url: s.src })}
                  </span>
                ) : (
                  <span key={i} className={s.style ? `mk mk--${s.style}` : 'mk'}>
                    {s.text}
                  </span>
                ),
              )}
            </div>
          )}
        </div>
      )}
    </div>
  );
}
