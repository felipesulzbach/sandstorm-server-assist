import { useEffect, useState } from 'react';
import { Send } from 'lucide-react';
import { useI18n } from '../../i18n';
import { backend, type AutomationInfo, type Profile, type SecretKey } from '../../shared/backend';
import { ErrorCard, Field, SecretField } from '../../shared/components';
import { useAction } from '../../shared/hooks';

type SetFn = <K extends keyof Profile>(k: K, v: Profile[K]) => void;
type ErrFn = (field: string) => string | null;
type DiscordKey = 'ready' | 'stopped' | 'crashed' | 'join' | 'leave';
const DISCORD: DiscordKey[] = ['ready', 'stopped', 'crashed', 'join', 'leave'];
const fieldOf = (k: DiscordKey) => `discord${k[0]!.toUpperCase()}${k.slice(1)}`;

/**
 * Automação (ADR-022): reinício diário (com aviso no chat e atualização pelo SteamCMD), reinício quando o
 * servidor fecha sozinho e avisos num canal do Discord. As regras ficam no Go; a tela só coleta.
 */
export function AutomationTab({ form, set, err }: { form: Profile; set: SetFn; err: ErrFn }) {
  const { t, date } = useI18n();
  const action = useAction();
  const [info, setInfo] = useState<AutomationInfo | null>(null);
  const [secrets, setSecrets] = useState<Record<SecretKey, boolean> | null>(null);
  const [tested, setTested] = useState(false);
  const ar = form.autoRestart;
  const dc = form.discord;

  useEffect(() => {
    backend
      .automation()
      .then(setInfo)
      .catch(() => {});
    backend
      .secretsState()
      .then(setSecrets)
      .catch(() => {});
  }, [form.autoRestart.daily]);

  const setAr = (p: Partial<Profile['autoRestart']>) =>
    set('autoRestart', { ...ar, ...p } as Profile['autoRestart']);
  const setDc = (p: Partial<Profile['discord']>) =>
    set('discord', { ...dc, ...p } as Profile['discord']);

  const next = info?.nextRestart ? new Date(info.nextRestart as unknown as string) : null;
  const nextOk = next && next.getFullYear() > 2000;

  return (
    <div className="form">
      <h3 className="form__title field--wide">{t('automation.restartTitle')}</h3>
      <Field
        label={t('automation.daily')}
        hint={nextOk ? t('automation.next', { when: date(next) }) : t('automation.dailyHint')}
        error={err('autoRestartDaily')}
      >
        <input
          className="input input--short"
          type="time"
          value={ar.daily}
          onChange={(e) => setAr({ daily: e.target.value })}
        />
      </Field>
      <label className="check check--block">
        <input
          type="checkbox"
          checked={ar.update}
          disabled={!ar.daily}
          onChange={(e) => setAr({ update: e.target.checked })}
        />
        <span>
          {t('automation.update')}
          <span className="field__hint">
            {info && !info.hasSteamCmd
              ? t('automation.updateNoSteamCmd')
              : t('automation.updateHint')}
          </span>
        </span>
      </label>
      <Field
        label={t('automation.warn')}
        hint={form.rcon ? t('automation.warnHint') : t('server.messages.needsRcon')}
        error={err('autoRestartWarn')}
        wide
      >
        <input
          className="input"
          maxLength={150}
          value={ar.warnMessage}
          disabled={!ar.daily}
          placeholder={t('automation.warnExample')}
          onChange={(e) => setAr({ warnMessage: e.target.value })}
        />
      </Field>
      <label className="check check--block field--wide">
        <input
          type="checkbox"
          checked={ar.onCrash}
          onChange={(e) => setAr({ onCrash: e.target.checked })}
        />
        <span>
          {t('automation.onCrash')}
          <span className="field__hint">{t('automation.onCrashHint')}</span>
        </span>
      </label>
      <p className="hint field--wide">{t('automation.appOpen')}</p>

      <h3 className="form__title field--wide">{t('automation.discordTitle')}</h3>
      <p className="muted field--wide">{t('automation.discordText')}</p>
      {secrets && (
        <div className="field--wide">
          <SecretField
            name="discordWebhook"
            label={t('automation.webhook')}
            hint={t('automation.webhookHint')}
            isSet={secrets.discordWebhook}
            onChanged={(st) => {
              setSecrets(st);
              backend
                .automation()
                .then(setInfo)
                .catch(() => {});
            }}
          />
        </div>
      )}
      {DISCORD.map((k) => (
        <Field key={k} label={t(`automation.discord.${k}`)} error={err(fieldOf(k))} wide>
          <input
            className="input"
            maxLength={300}
            value={dc[k]}
            placeholder={t(`automation.discordExample.${k}`)}
            onChange={(e) => setDc({ [k]: e.target.value })}
          />
        </Field>
      ))}
      <div className="row row--wrap field--wide">
        <button
          type="button"
          className="btn btn--small btn--ghost"
          onClick={() =>
            setDc(
              Object.fromEntries(
                DISCORD.map((k) => [k, t(`automation.discordExample.${k}`)]),
              ) as Partial<Profile['discord']>,
            )
          }
        >
          {t('automation.fillExamples')}
        </button>
        <button
          type="button"
          className="btn btn--small"
          disabled={!secrets?.discordWebhook || action.busy}
          onClick={async () => {
            setTested(false);
            const ok = await action.run(async () => {
              await backend.testDiscord();
              return true;
            });
            if (ok) setTested(true);
          }}
        >
          <Send size={14} aria-hidden /> {t('automation.test')}
        </button>
        {tested && <span className="pill pill--ok">{t('automation.testSent')}</span>}
      </div>
      <p className="hint field--wide">{t('automation.tags')}</p>
      {action.error && (
        <div className="field--wide">
          <ErrorCard error={action.error} onClose={() => action.setError(null)} />
        </div>
      )}
    </div>
  );
}
