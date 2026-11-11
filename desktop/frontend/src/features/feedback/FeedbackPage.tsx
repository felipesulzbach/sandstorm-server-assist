import { useState } from 'react';
import { Bug, FolderOpen, ImagePlus, Lightbulb, Send, X } from 'lucide-react';
import { useI18n } from '../../i18n';
import {
  backend,
  type ReportImage,
  type ReportKind,
  type ReportResult,
} from '../../shared/backend';
import { ErrorCard, Field, PageHead } from '../../shared/components';
import { bytes } from '../../shared/format';
import { useAction } from '../../shared/hooks';

/**
 * Relato de problema e sugestão de melhoria (ADR-034). O Go monta um zip no PC (texto, log mascarado e imagem
 * opcional) e abre o formulário do GitHub já preenchido; a pessoa anexa o zip lá. Nada sai do PC sem ela.
 */
export function FeedbackPage() {
  const { t, tl } = useI18n();
  const action = useAction();
  const imageAction = useAction();
  const [kind, setKind] = useState<ReportKind>('bug');
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [includeLog, setIncludeLog] = useState(true);
  const [image, setImage] = useState<ReportImage | null>(null);
  const [done, setDone] = useState<ReportResult | null>(null);

  const chooseKind = (k: ReportKind) => {
    setKind(k);
    setIncludeLog(k === 'bug');
    setDone(null);
  };

  const chooseImage = async () => {
    const r = await imageAction.run(backend.chooseReportImage);
    if (r?.name) setImage(r);
  };

  const removeImage = async () => {
    await backend.clearReportImage();
    setImage(null);
  };

  const create = async () => {
    setDone(null);
    const r = await action.run(() => backend.createReport(kind, title, description, includeLog));
    if (r) {
      setDone(r);
      setTitle('');
      setDescription('');
      setImage(null);
    }
  };

  return (
    <div className="page">
      <PageHead title={t('feedback.title')} subtitle={t('feedback.subtitle')} />

      <section className="card">
        <div className="segmented" role="radiogroup" aria-label={t('feedback.kindLabel')}>
          {(['bug', 'idea'] as ReportKind[]).map((k) => (
            <button
              key={k}
              type="button"
              role="radio"
              aria-checked={kind === k}
              className={kind === k ? 'is-active' : ''}
              onClick={() => chooseKind(k)}
            >
              {k === 'bug' ? <Bug size={16} aria-hidden /> : <Lightbulb size={16} aria-hidden />}{' '}
              {t(`feedback.kind.${k}`)}
            </button>
          ))}
        </div>
        <p className="muted">{t(`feedback.intro.${kind}`)}</p>

        {action.error && <ErrorCard error={action.error} onClose={() => action.setError(null)} />}

        <div className="form form--single">
          <Field label={t('feedback.titleField')} hint={t('feedback.titleHint')}>
            <input
              className="input"
              maxLength={120}
              value={title}
              onChange={(e) => setTitle(e.target.value)}
            />
          </Field>
          <Field label={t('feedback.description')} hint={t(`feedback.descriptionHint.${kind}`)}>
            <textarea
              className="input"
              rows={6}
              maxLength={4000}
              value={description}
              onChange={(e) => setDescription(e.target.value)}
            />
          </Field>

          <label className="check check--block">
            <input
              type="checkbox"
              checked={includeLog}
              onChange={(e) => setIncludeLog(e.target.checked)}
            />
            <span>
              {t('feedback.includeLog')}
              <span className="field__hint">{t('feedback.includeLogHint')}</span>
            </span>
          </label>

          <Field label={t('feedback.image')} hint={t('feedback.imageHint')}>
            <div className="row row--wrap">
              {image ? (
                <>
                  <span>
                    {image.name} ({bytes(image.size, t)})
                  </span>
                  <button type="button" className="btn btn--small btn--ghost" onClick={removeImage}>
                    <X size={14} aria-hidden /> {t('feedback.removeImage')}
                  </button>
                </>
              ) : (
                <button
                  type="button"
                  className="btn btn--small"
                  disabled={imageAction.busy}
                  onClick={chooseImage}
                >
                  <ImagePlus size={14} aria-hidden /> {t('feedback.chooseImage')}
                </button>
              )}
            </div>
          </Field>
          {imageAction.error && (
            <ErrorCard error={imageAction.error} onClose={() => imageAction.setError(null)} />
          )}

          <button
            type="button"
            className="btn btn--primary"
            disabled={action.busy}
            onClick={create}
          >
            <Send size={16} aria-hidden /> {t('feedback.create')}
          </button>
          <p className="hint">{t('feedback.privacy')}</p>
        </div>

        {done && (
          <div className="banner banner--ok">
            <p>{t('feedback.done', { file: done.file, size: bytes(done.size, t) })}</p>
            <ol>
              {tl('feedback.next').map((s) => (
                <li key={s}>{s}</li>
              ))}
            </ol>
            <button type="button" className="btn btn--small" onClick={() => backend.showReport()}>
              <FolderOpen size={14} aria-hidden /> {t('feedback.showFile')}
            </button>
          </div>
        )}
      </section>
    </div>
  );
}
