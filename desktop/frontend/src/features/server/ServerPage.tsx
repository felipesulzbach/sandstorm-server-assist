import { useEffect, useMemo, useState, type ReactNode } from 'react';
import { Plus, RotateCcw, Save, ServerCog, Trash2 } from 'lucide-react';
import { useI18n } from '../../i18n';
import { backend, type MapsInfo, type Profile, type SecretKey } from '../../shared/backend';
import { ErrorCard, Field, Guide, PageHead, SecretField } from '../../shared/components';
import { useAction } from '../../shared/hooks';
import { fieldError, toAppError } from '../../shared/problem';
import { useExistingSync } from '../sync/ExistingSync';
import { AutomationTab } from './AutomationTab';
import { CycleLighting, ModMapPicker, MutatorPicker } from './MapsExtras';
import { MarkupEditor } from './MarkupEditor';
import { BindIpField, FirewallCard } from './NetworkExtras';
import { PortCheck } from './PortCheck';
import { ProfileBar } from './ProfileBar';
import { RulesTab } from './RulesTab';

type Tab =
  'general' | 'maps' | 'rules' | 'admins' | 'secrets' | 'network' | 'messages' | 'automation';
const TABS: Tab[] = [
  'general',
  'maps',
  'rules',
  'admins',
  'secrets',
  'network',
  'messages',
  'automation',
];

export function ServerPage() {
  const { t } = useI18n();
  const [tab, setTab] = useState<Tab>('general');
  const [saved, setSaved] = useState<Profile | null>(null);
  const [form, setForm] = useState<Profile | null>(null);
  const [maps, setMaps] = useState<MapsInfo | null>(null);
  const [okAt, setOkAt] = useState(0);
  const [syncedAt, setSyncedAt] = useState(0);
  const action = useAction();
  const sync = useExistingSync((r) => {
    setSaved(r.profile);
    setForm(r.profile);
    setSyncedAt(Date.now()); // recarrega a aba de senhas (o cofre pode ter mudado)
  });

  useEffect(() => {
    backend.profile().then((p) => {
      setSaved(p);
      setForm(p);
    });
    backend
      .maps()
      .then(setMaps)
      .catch(() => {});
  }, []);

  const dirty = useMemo(() => JSON.stringify(saved) !== JSON.stringify(form), [saved, form]);
  if (!form) return null;

  const patch = (p: Partial<Profile>) => setForm((f) => ({ ...f!, ...p }) as Profile);
  const set = <K extends keyof Profile>(k: K, v: Profile[K]) =>
    patch({ [k]: v } as Partial<Profile>);
  // O Go manda o id da regra ({rule}); a tela mostra o nome traduzido.
  const withLabels = (p?: Record<string, string>) =>
    p?.rule ? { ...p, rule: t(`rules.items.${p.rule}.label`) } : p;
  const err = (f: string) => {
    const e = fieldError(action.error, f);
    return e?.msg ? t(e.msg, withLabels(e.params)) : null;
  };

  const save = async () => {
    const p = await action.run(() => backend.saveProfile(form));
    if (p) {
      setSaved(p);
      setForm(p);
      setOkAt(Date.now());
    }
    return p;
  };

  // O que o usuário digitou e ainda não salvou também conta como "meu" na comparação: salva antes.
  const syncExisting = async () => {
    if (dirty && !(await save())) return;
    await sync.check('manual');
  };

  const generalError = action.error && !action.error.field ? action.error : null;
  const fieldTab: Record<string, Tab> = {
    hostname: 'general',
    port: 'general',
    queryPort: 'general',
    maxPlayers: 'general',
    rconPort: 'secrets',
    startMap: 'maps',
    startScenario: 'maps',
    travelMap: 'maps',
    travelScenario: 'maps',
    mutators: 'maps',
    mapCycle: 'maps',
    admins: 'admins',
    connectIp: 'network',
    motd: 'messages',
    joinMessage: 'messages',
    leaveMessage: 'messages',
    announcements: 'messages',
    announceEvery: 'messages',
    startLighting: 'maps',
    bindIp: 'network',
    modIds: 'general',
    rules: 'rules',
    autoRestartDaily: 'automation',
    autoRestartWarn: 'automation',
    discordReady: 'automation',
    discordStopped: 'automation',
    discordCrashed: 'automation',
    discordJoin: 'automation',
    discordLeave: 'automation',
  };
  const errTab = action.error?.field ? fieldTab[action.error.field] : undefined;

  return (
    <div className="page">
      <PageHead title={t('server.title')} subtitle={t('server.subtitle')}>
        {dirty && <span className="pill pill--warn">{t('common.unsaved')}</span>}
        {!dirty && Date.now() - okAt < 4000 && (
          <span className="pill pill--ok">{t('common.saved')}</span>
        )}
        <button
          type="button"
          className="btn btn--ghost"
          disabled={!form.serverDir || sync.busy || action.busy}
          title={form.serverDir ? t('sync.buttonHint') : t('sync.needDir')}
          onClick={syncExisting}
        >
          <ServerCog size={16} aria-hidden />
          {t('sync.button')}
        </button>
        <button
          type="button"
          className="btn btn--ghost"
          onClick={async () => {
            const d = await backend.defaultProfile();
            setForm({
              ...d,
              serverDir: form.serverDir,
              steamcmdDir: form.steamcmdDir,
              admins: form.admins,
              mapCycle: form.mapCycle,
              motd: form.motd,
              joinMessage: form.joinMessage,
              leaveMessage: form.leaveMessage,
              announcements: form.announcements,
              announceEvery: form.announceEvery,
              rules: form.rules,
              autoRestart: form.autoRestart,
              discord: form.discord,
              bindIp: form.bindIp,
              modIds: form.modIds,
              cmdModList: form.cmdModList,
            } as Profile);
          }}
        >
          <RotateCcw size={16} aria-hidden />
          {t('server.reset')}
        </button>
        <button
          type="button"
          className="btn btn--primary"
          disabled={!dirty || action.busy}
          onClick={save}
        >
          <Save size={16} aria-hidden />
          {t('common.save')}
        </button>
      </PageHead>

      <ProfileBar
        dirty={dirty}
        onSwitched={(p) => {
          setSaved(p);
          setForm(p);
          setSyncedAt(Date.now());
        }}
      />
      {generalError && <ErrorCard error={generalError} onClose={() => action.setError(null)} />}
      {sync.node}
      {action.error?.field && errTab && errTab !== tab && (
        <p className="banner banner--error">
          {t(action.error.msg ?? 'problems.invalid-field.what', withLabels(action.error.params))}{' '}
          <button type="button" className="link" onClick={() => setTab(errTab)}>
            {t(`server.tabs.${errTab}`)}
          </button>
        </p>
      )}

      <div className="tabs" role="tablist">
        {TABS.map((id) => (
          <button
            key={id}
            type="button"
            role="tab"
            aria-selected={tab === id}
            className={`tabs__tab ${tab === id ? 'is-active' : ''}`}
            onClick={() => setTab(id)}
          >
            {t(`server.tabs.${id}`)}
          </button>
        ))}
      </div>

      <section className="card">
        {tab === 'general' && (
          <div className="form">
            <Field
              label={t('server.hostname')}
              hint={t('server.hostnameHint')}
              error={err('hostname')}
              wide
            >
              <input
                className="input"
                value={form.hostname}
                maxLength={64}
                onChange={(e) => set('hostname', e.target.value)}
              />
            </Field>
            <Field label={t('server.maxPlayers')} error={err('maxPlayers')}>
              <input
                className="input"
                type="number"
                min={1}
                max={64}
                value={form.maxPlayers}
                onChange={(e) => set('maxPlayers', Number(e.target.value))}
              />
            </Field>
            <Field label={t('server.port')} error={err('port')}>
              <input
                className="input"
                type="number"
                min={1024}
                max={65535}
                value={form.port}
                onChange={(e) => set('port', Number(e.target.value))}
              />
            </Field>
            <Field
              label={t('server.queryPort')}
              hint={t('server.portsHint')}
              error={err('queryPort')}
            >
              <input
                className="input"
                type="number"
                min={1024}
                max={65535}
                value={form.queryPort}
                onChange={(e) => set('queryPort', Number(e.target.value))}
              />
            </Field>
            <label className="check check--block field--wide">
              <input
                type="checkbox"
                checked={form.useMods}
                onChange={(e) => set('useMods', e.target.checked)}
              />
              <span>
                {t('server.useMods')}
                <span className="field__hint">{t('server.useModsHint')}</span>
              </span>
            </label>
            {form.useMods && (
              <>
                <label className="check check--block field--wide">
                  <input
                    type="checkbox"
                    checked={form.cmdModList}
                    onChange={(e) => set('cmdModList', e.target.checked)}
                  />
                  <span>
                    {t('server.modList')}
                    <span className="field__hint">{t('server.modListHint')}</span>
                  </span>
                </label>
                {form.cmdModList && (
                  <Field
                    label={t('server.modIds')}
                    hint={t('server.modIdsHint')}
                    error={err('modIds')}
                    wide
                  >
                    <ModIdsInput value={form.modIds ?? []} onChange={(v) => set('modIds', v)} />
                  </Field>
                )}
              </>
            )}
          </div>
        )}

        {tab === 'maps' && <MapsTab form={form} set={set} patch={patch} maps={maps} err={err} />}
        {tab === 'rules' && (
          <RulesTab
            form={form}
            set={(r) => set('rules', r as Profile['rules'])}
            error={err('rules')}
          />
        )}
        {tab === 'automation' && <AutomationTab form={form} set={set} err={err} />}
        {tab === 'admins' && <AdminsTab form={form} set={set} err={err} />}
        {tab === 'secrets' && <SecretsTab key={syncedAt} form={form} set={set} err={err} />}

        {tab === 'network' && (
          <div className="form">
            <Field
              label={t('server.network.connectIp')}
              hint={t('server.network.connectIpHint')}
              error={err('connectIp')}
              wide
            >
              <input
                className="input"
                value={form.connectIp}
                placeholder="203.0.113.10"
                onChange={(e) => set('connectIp', e.target.value.trim())}
              />
            </Field>
            <BindIpField
              value={form.bindIp ?? ''}
              onChange={(ip) => set('bindIp', ip)}
              error={err('bindIp')}
            />
            <div className="field--wide">
              <p>
                {t('server.network.ports', {
                  port: form.port,
                  queryPort: form.queryPort,
                  rconPort: form.rconPort,
                })}
              </p>
              <Guide
                id="ports"
                label={t('guides.ports.title')}
                params={{ port: form.port, queryPort: form.queryPort, rconPort: form.rconPort }}
              />
              <Guide id="vpn" label={t('guides.vpn.title')} />
              <p className="hint">{t('server.network.firewall')}</p>
              <FirewallCard />
              <PortCheck connectIp={form.connectIp} onUseIp={(ip) => set('connectIp', ip)} />
            </div>
          </div>
        )}
        {tab === 'messages' && (
          <div className="form">
            <Block
              label={t('server.messages.motd')}
              hint={t('server.messages.motdHint')}
              error={err('motd')}
            >
              <MarkupEditor
                value={form.motd ?? ''}
                maxLength={1000}
                onChange={(v) => set('motd', v)}
              />
            </Block>
            <Block
              label={t('server.messages.join')}
              hint={t('server.messages.chatHint')}
              error={err('joinMessage')}
            >
              <MarkupEditor
                chat
                multiline={false}
                preview={false}
                maxLength={150}
                value={form.joinMessage ?? ''}
                placeholder={t('server.messages.joinExample')}
                onChange={(v) => set('joinMessage', v)}
              />
            </Block>
            <Block
              label={t('server.messages.leave')}
              hint={t('server.messages.chatHint')}
              error={err('leaveMessage')}
            >
              <MarkupEditor
                chat
                multiline={false}
                preview={false}
                maxLength={150}
                value={form.leaveMessage ?? ''}
                placeholder={t('server.messages.leaveExample')}
                onChange={(v) => set('leaveMessage', v)}
              />
            </Block>
            <Block
              label={t('server.messages.announcements')}
              hint={t('server.messages.announcementsHint')}
              error={err('announcements')}
            >
              <MarkupEditor
                chat
                rows={4}
                value={(form.announcements ?? []).join('\n')}
                placeholder={t('server.messages.announcementsExample')}
                onChange={(v) =>
                  set(
                    'announcements',
                    v
                      .split('\n')
                      .map((x) => x.trimStart())
                      .filter((x, i, all) => x !== '' || i === all.length - 1),
                  )
                }
              />
            </Block>
            <Field label={t('server.messages.every')} error={err('announceEvery')}>
              <input
                className="input input--short"
                type="number"
                min={2}
                max={240}
                value={form.announceEvery}
                onChange={(e) => set('announceEvery', Number(e.target.value))}
              />
            </Field>
            <p className={`field--wide ${form.rcon ? 'hint' : 'banner banner--info'}`}>
              {t(form.rcon ? 'server.messages.howChat' : 'server.messages.needsRcon')}
            </p>
          </div>
        )}
      </section>
      <p className="hint">{t('common.appliesNextStart')}</p>
    </div>
  );
}

/** Campo com controles próprios dentro (botões, caixas): bloco em vez de <label>. */
function Block({
  label,
  hint,
  error,
  children,
}: {
  label: string;
  hint?: string;
  error?: string | null;
  children: ReactNode;
}) {
  return (
    <div className={`field field--wide ${error ? 'has-error' : ''}`}>
      <span className="field__label">{label}</span>
      {children}
      {error ? (
        <span className="field__error">{error}</span>
      ) : hint ? (
        <span className="field__hint">{hint}</span>
      ) : null}
    </div>
  );
}

type SetFn = <K extends keyof Profile>(k: K, v: Profile[K]) => void;
type ErrFn = (field: string) => string | null;

function MapsTab({
  form,
  set,
  patch,
  maps,
  err,
}: {
  form: Profile;
  set: SetFn;
  patch: (p: Partial<Profile>) => void;
  maps: MapsInfo | null;
  err: ErrFn;
}) {
  const { t } = useI18n();
  const known = maps?.maps.find((m) => m.map === form.startMap);
  const modeId = known ? form.startScenario.replace(`Scenario_${known.scenario}_`, '') : '';
  const knownMode = !!known && (maps?.modes ?? []).some((m) => m.id === modeId);

  return (
    <div className="form">
      <Field label={t('server.startMap')} hint={t('server.startMapHint')} error={err('startMap')}>
        <select
          className="select"
          value={known ? known.map : ''}
          onChange={(e) => {
            const m = maps?.maps.find((x) => x.map === e.target.value);
            if (m) {
              const mode = knownMode ? modeId : 'Checkpoint_Security';
              patch({ startMap: m.map, startScenario: `Scenario_${m.scenario}_${mode}` });
            }
          }}
        >
          {!known && <option value="">{form.startMap}</option>}
          {(maps?.maps ?? []).map((m) => (
            <option key={m.map} value={m.map}>
              {m.label}
            </option>
          ))}
        </select>
      </Field>
      <Field label={t('lighting.start')} error={err('startLighting')}>
        <select
          className="select"
          value={form.startLighting ?? ''}
          onChange={(e) => set('startLighting', e.target.value)}
        >
          <option value="">{t('lighting.default')}</option>
          <option value="Day">{t('lighting.Day')}</option>
          <option value="Night">{t('lighting.Night')}</option>
        </select>
      </Field>
      <Field label={t('server.mode')} error={err('startScenario')}>
        <select
          className="select"
          value={knownMode ? modeId : 'custom'}
          onChange={(e) => {
            if (known && e.target.value !== 'custom')
              set('startScenario', `Scenario_${known.scenario}_${e.target.value}`);
          }}
        >
          {(maps?.modes ?? []).map((m) => (
            <option key={m.id} value={m.id}>
              {t(`maps.mode.${m.id}`)}
            </option>
          ))}
          <option value="custom">{t('maps.mode.custom')}</option>
        </select>
      </Field>
      {!knownMode && (
        <Field label={t('server.customScenario')} error={err('startScenario')} wide>
          <input
            className="input"
            value={form.startScenario}
            onChange={(e) => set('startScenario', e.target.value.trim())}
          />
        </Field>
      )}

      <ModMapPicker
        map={form.travelMap}
        scenario={form.travelScenario}
        disabled={!form.useMods}
        onChange={(travelMap, travelScenario) => patch({ travelMap, travelScenario })}
        mapError={err('travelMap')}
        scenarioError={err('travelScenario')}
      />

      {/* Bloco, não <label>: a lista tem as próprias caixas de marcar. */}
      <div className={`field field--wide ${err('mutators') ? 'has-error' : ''}`}>
        <span className="field__label">{t('server.mutators')}</span>
        <MutatorPicker value={form.mutators ?? []} onChange={(v) => set('mutators', v)} />
        <span className={err('mutators') ? 'field__error' : 'field__hint'}>
          {err('mutators') ?? t('server.mutatorsHint')}
        </span>
      </div>
      <Field
        label={t('server.mapCycle')}
        hint={t('server.mapCycleHint')}
        error={err('mapCycle')}
        wide
      >
        <textarea
          className="input textarea"
          rows={8}
          spellCheck={false}
          value={(form.mapCycle ?? []).join('\n')}
          onChange={(e) =>
            set(
              'mapCycle',
              e.target.value
                .split('\n')
                .map((s) => s.trim())
                .filter(Boolean),
            )
          }
        />
      </Field>
      <div className="field--wide">
        <CycleLighting lines={form.mapCycle ?? []} onChange={(l) => set('mapCycle', l)} />
      </div>
    </div>
  );
}

/** Ids dos mods separados por vírgula (a tela só junta os números; o Go valida). */
function ModIdsInput({ value, onChange }: { value: number[]; onChange: (v: number[]) => void }) {
  const [text, setText] = useState(value.join(', '));
  return (
    <input
      className="input"
      value={text}
      spellCheck={false}
      placeholder="1234567, 2345678"
      onChange={(e) => {
        setText(e.target.value);
        onChange(
          e.target.value
            .split(/[\s,;]+/)
            .map((x) => Number(x))
            .filter((n) => Number.isInteger(n) && n > 0),
        );
      }}
    />
  );
}

function AdminsTab({ form, set, err }: { form: Profile; set: SetFn; err: ErrFn }) {
  const { t } = useI18n();
  const [text, setText] = useState('');
  const [name, setName] = useState('');
  const [msg, setMsg] = useState<string | null>(null);
  const [bad, setBad] = useState<string | null>(null);
  const admins = form.admins ?? [];

  const add = async () => {
    setBad(null);
    setMsg(null);
    try {
      const id = await backend.extractSteamId(text);
      if (!admins.some((a) => a.steamId === id))
        set('admins', [...admins, { steamId: id, name: name.trim() }]);
      setMsg(t('server.admins.found', { id }));
      setText('');
      setName('');
    } catch (e) {
      const pe = toAppError(e);
      setBad(pe.msg ? t(pe.msg, pe.params) : t(`problems.${pe.code}.title`));
    }
  };

  return (
    <div className="form">
      <p className="muted field--wide">{t('server.admins.text')}</p>
      <Field
        label={t('server.admins.steamId')}
        error={bad ?? err('admins')}
        hint={msg ?? undefined}
      >
        <input
          className="input"
          value={text}
          spellCheck={false}
          onChange={(e) => setText(e.target.value)}
          onKeyDown={(e) => e.key === 'Enter' && add()}
        />
      </Field>
      <Field label={t('server.admins.name')}>
        <div className="row">
          <input
            className="input"
            value={name}
            maxLength={40}
            onChange={(e) => setName(e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && add()}
          />
          <button type="button" className="btn btn--primary" disabled={!text.trim()} onClick={add}>
            <Plus size={16} aria-hidden />
            {t('common.add')}
          </button>
        </div>
      </Field>
      <div className="field--wide">
        <Guide id="steamId" links={['steam-profile', 'steamid-io']} />
      </div>
      <ul className="list field--wide">
        {admins.length === 0 && <li className="muted">{t('server.admins.empty')}</li>}
        {admins.map((a) => (
          <li key={a.steamId} className="list__row">
            <code>{a.steamId}</code>
            <span>{a.name}</span>
            <button
              type="button"
              className="icon-btn"
              aria-label={t('common.remove')}
              title={t('common.remove')}
              onClick={() =>
                set(
                  'admins',
                  admins.filter((x) => x.steamId !== a.steamId),
                )
              }
            >
              <Trash2 size={16} />
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
}

function SecretsTab({ form, set, err }: { form: Profile; set: SetFn; err: ErrFn }) {
  const { t } = useI18n();
  const [state, setState] = useState<Record<SecretKey, boolean> | null>(null);
  const [rconInIni, setRconInIni] = useState(false);
  const [moved, setMoved] = useState(false);
  const action = useAction();

  useEffect(() => {
    backend
      .secretsState()
      .then(setState)
      .catch(() => {});
    backend
      .rconInIni()
      .then(setRconInIni)
      .catch(() => {});
  }, []);
  if (!state) return null;

  return (
    <div className="form form--single">
      <p className="banner banner--info">{t('server.secrets.vaultNote')}</p>
      {action.error && <ErrorCard error={action.error} onClose={() => action.setError(null)} />}
      {rconInIni && (
        <div className="banner banner--warn">
          <span>{t('server.secrets.rconInIni')}</span>
          <button
            type="button"
            className="btn btn--small"
            onClick={async () => {
              const ok = await action.run(async () => {
                await backend.moveRconToVault();
                return true;
              });
              if (ok) {
                setRconInIni(false);
                setMoved(true);
                backend.secretsState().then(setState);
              }
            }}
          >
            {t('server.secrets.moveToVault')}
          </button>
        </div>
      )}
      {moved && <p className="banner banner--ok">{t('server.secrets.moved')}</p>}

      <SecretField
        name="serverPassword"
        label={t('server.secrets.serverPassword')}
        hint={t('server.secrets.serverPasswordHint')}
        isSet={state.serverPassword}
        onChanged={setState}
      />

      <label className="check check--block">
        <input
          type="checkbox"
          checked={form.rcon}
          onChange={(e) => set('rcon', e.target.checked)}
        />
        <span>{t('server.secrets.rcon')}</span>
      </label>
      {form.rcon && (
        <>
          <SecretField
            name="rconPassword"
            label={t('server.secrets.rconPassword')}
            isSet={state.rconPassword}
            onChanged={setState}
          />
          <Field label={t('server.secrets.rconPort')} error={err('rconPort')}>
            <input
              className="input input--short"
              type="number"
              min={1024}
              max={65535}
              value={form.rconPort}
              onChange={(e) => set('rconPort', Number(e.target.value))}
            />
          </Field>
        </>
      )}

      <SecretField
        name="gsltToken"
        label={t('server.secrets.gsltToken')}
        hint={t('server.secrets.gsltTokenHint')}
        isSet={state.gsltToken}
        onChanged={setState}
      >
        <Guide id="gslt" links={['steam-gslt']} />
      </SecretField>
      <SecretField
        name="gameStatsToken"
        label={t('server.secrets.gameStatsToken')}
        hint={t('server.secrets.gameStatsTokenHint')}
        isSet={state.gameStatsToken}
        onChanged={setState}
      >
        <Guide id="gameStats" links={['gamestats']} />
      </SecretField>
    </div>
  );
}
