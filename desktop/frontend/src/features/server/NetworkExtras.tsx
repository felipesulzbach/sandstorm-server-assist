import { useEffect, useState } from 'react';
import { CheckCircle2, ShieldCheck } from 'lucide-react';
import { useI18n } from '../../i18n';
import { backend, type FirewallInfo, type LocalAddress } from '../../shared/backend';
import { ErrorCard, Field } from '../../shared/components';
import { useAction } from '../../shared/hooks';

/** Placa de rede do servidor (-MultiHome): todas ou uma das placas desta máquina (inclui as de programas de VPN). */
export function BindIpField({
  value,
  onChange,
  error,
}: {
  value: string;
  onChange: (ip: string) => void;
  error: string | null;
}) {
  const { t } = useI18n();
  const [list, setList] = useState<LocalAddress[]>([]);

  useEffect(() => {
    backend
      .localAddresses()
      .then((l) => setList(l ?? []))
      .catch(() => {});
  }, []);

  const known = value === '' || list.some((a) => a.ip === value);
  return (
    <Field
      label={t('server.network.bindIp')}
      hint={known ? t('server.network.bindIpHint') : t('server.network.bindIpMissing')}
      error={error}
      wide
    >
      <select className="select" value={value} onChange={(e) => onChange(e.target.value)}>
        <option value="">{t('server.network.bindIpAll')}</option>
        {!known && <option value={value}>{value}</option>}
        {list.map((a) => (
          <option key={a.ip} value={a.ip}>
            {a.ip} · {a.name}
          </option>
        ))}
      </select>
    </Field>
  );
}

/** Regra do Firewall do Windows para o servidor (o Windows pede a confirmação de administrador). */
export function FirewallCard() {
  const { t } = useI18n();
  const action = useAction();
  const [info, setInfo] = useState<FirewallInfo | null>(null);

  useEffect(() => {
    backend
      .firewall()
      .then(setInfo)
      .catch(() => {});
  }, []);
  if (!info) return null;

  return (
    <div className="firewall field--wide">
      <span className="field__label">{t('firewall.title')}</span>
      {!info.ready ? (
        <p className="hint">{t('firewall.notInstalled')}</p>
      ) : info.rules ? (
        <p className="row">
          <CheckCircle2 size={16} aria-hidden className="text-ok" /> {t('firewall.ok')}
        </p>
      ) : (
        <>
          <p className="muted">{t('firewall.text')}</p>
          <div className="row">
            <button
              type="button"
              className="btn"
              disabled={action.busy}
              onClick={async () => {
                const r = await action.run(backend.createFirewallRule);
                if (r) setInfo(r);
              }}
            >
              <ShieldCheck size={16} aria-hidden /> {t('firewall.create')}
            </button>
            <span className="hint">{t('firewall.uac')}</span>
          </div>
        </>
      )}
      {action.error && <ErrorCard error={action.error} onClose={() => action.setError(null)} />}
    </div>
  );
}
