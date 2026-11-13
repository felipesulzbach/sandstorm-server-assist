import { useEffect, useState } from 'react';
import { useI18n } from '../../i18n';
import { backend, type Profile, type RuleDef, type RulesInfo } from '../../shared/backend';

/**
 * Regras do jogo (ADR-022): opções do Game.ini e do Engine.ini que o Go sabe gravar. A lista, os limites e os
 * padrões vêm do Go; a tela só mostra e coleta. Vazio = padrão do jogo (o app tira a opção do arquivo).
 */
export function RulesTab({
  form,
  set,
  error,
}: {
  form: Profile;
  set: (rules: Record<string, string> | null) => void;
  error: string | null;
}) {
  const { t } = useI18n();
  const [info, setInfo] = useState<RulesInfo | null>(null);
  const [fromServer, setFromServer] = useState<Record<string, string>>({});
  const rules = (form.rules ?? null) as Record<string, string> | null;
  const managed = rules !== null;

  useEffect(() => {
    backend
      .rules()
      .then(setInfo)
      .catch(() => {});
  }, []);
  // Com a opção desligada, o Go lê as regras que o servidor já tem (ADR-040) para oferecer usá-las.
  useEffect(() => {
    if (managed) return;
    backend
      .serverRules()
      .then((r) => setFromServer(r ?? {}))
      .catch(() => setFromServer({}));
  }, [managed, form.serverDir]);
  if (!info) return null;
  const serverCount = Object.keys(fromServer).length;
  const value = (id: string) => rules?.[id] ?? '';
  const put = (id: string, v: string) => {
    const next = { ...(rules ?? {}) };
    if (v === '') delete next[id];
    else next[id] = v;
    set(next);
  };

  const input = (r: RuleDef) => {
    if (r.kind === 'bool') {
      return (
        <select
          className="select"
          value={value(r.id)}
          disabled={!managed}
          onChange={(e) => put(r.id, e.target.value)}
        >
          <option value="">
            {t('rules.default', { value: t(r.default === 'True' ? 'rules.on' : 'rules.off') })}
          </option>
          <option value="True">{t('rules.on')}</option>
          <option value="False">{t('rules.off')}</option>
        </select>
      );
    }
    return (
      <input
        className="input"
        type="number"
        min={r.min}
        max={r.max}
        step={r.kind === 'float' ? 0.05 : 1}
        value={value(r.id)}
        disabled={!managed}
        placeholder={r.default ? t('rules.default', { value: r.default }) : t('rules.gameDefault')}
        onChange={(e) => put(r.id, e.target.value)}
      />
    );
  };

  return (
    <div className="form">
      <label className="check check--block field--wide">
        <input
          type="checkbox"
          checked={managed}
          onChange={(e) => set(e.target.checked ? {} : null)}
        />
        <span>
          {t('rules.manage')}
          <span className="field__hint">{t('rules.manageHint')}</span>
        </span>
      </label>
      {!managed && serverCount > 0 && (
        <div className="banner banner--info field--wide rules__server" role="status">
          <p>{t('rules.fromServer', { n: serverCount })}</p>
          <button type="button" className="btn btn--small" onClick={() => set({ ...fromServer })}>
            {t('rules.useServer')}
          </button>
        </div>
      )}
      {error && <p className="banner banner--error field--wide">{error}</p>}
      {managed &&
        info.groups.map((g) => (
          <fieldset key={g} className="rules__group field--wide">
            <legend>{t(`rules.groups.${g}`)}</legend>
            <div className="form">
              {info.rules
                .filter((r) => r.group === g)
                .map((r) => (
                  <label key={r.id} className="field">
                    <span className="field__label">{t(`rules.items.${r.id}.label`)}</span>
                    {input(r)}
                    <span className="field__hint">
                      {t(`rules.items.${r.id}.hint`)}
                      {r.kind !== 'bool' && ` (${r.min}–${r.max})`}
                    </span>
                  </label>
                ))}
            </div>
          </fieldset>
        ))}
    </div>
  );
}
