import { useState } from 'react';
import { AlertTriangle, CheckCircle2, Globe, Info, Stethoscope, XCircle } from 'lucide-react';
import { useI18n } from '../../i18n';
import { backend, type PortCheck as Result } from '../../shared/backend';
import { ErrorCard } from '../../shared/components';
import { useAction } from '../../shared/hooks';

/**
 * Diagnóstico das portas feito no próprio PC (ADR-033). O Go mede (servidor, Firewall, roteador por UPnP, IP
 * público e CGNAT por STUN, retorno pelo IP público, lista da Steam) e devolve itens com código; a tela só traduz.
 * Só um teste de fora da rede prova que a porta abre: a tela ensina o teste com um amigo.
 */
const ICONS = {
  ok: CheckCircle2,
  bad: XCircle,
  warn: AlertTriangle,
  info: Info,
} as const;

export function PortCheck({
  connectIp,
  onUseIp,
}: {
  connectIp: string;
  onUseIp: (ip: string) => void;
}) {
  const { t, tl } = useI18n();
  const action = useAction();
  const [res, setRes] = useState<Result | null>(null);

  const run = async () => {
    setRes(null);
    const r = await action.run(backend.checkPorts);
    if (r) setRes(r);
  };

  return (
    <div className="portcheck field--wide">
      <div className="row row--wrap">
        <button type="button" className="btn" disabled={action.busy} onClick={run}>
          <Stethoscope size={16} aria-hidden />{' '}
          {action.busy ? t('server.portCheck.running') : t('server.portCheck.button')}
        </button>
        <span className="hint">{t('server.portCheck.hint')}</span>
      </div>

      {action.error && <ErrorCard error={action.error} onClose={() => action.setError(null)} />}

      {res && (
        <>
          <p className={`portcheck__verdict portcheck__verdict--${res.verdict}`}>
            {t(`server.portCheck.verdict.${res.verdict}`)}
          </p>
          <ul className="portcheck__list">
            {(res.items ?? []).map((it) => {
              const Icon = ICONS[it.status as keyof typeof ICONS] ?? Info;
              return (
                <li key={it.key + JSON.stringify(it.params ?? {})} className={`is-${it.status}`}>
                  <Icon size={16} aria-hidden />
                  <span>{t(`server.portCheck.items.${it.key}`, it.params ?? {})}</span>
                </li>
              );
            })}
            <li className="is-info">
              <Info size={16} aria-hidden />
              <span>{t('server.portCheck.gamePort', { port: res.gamePort })}</span>
            </li>
            {res.publicIp && (
              <li className="is-info">
                <Globe size={16} aria-hidden />
                <span>{t('server.portCheck.publicIp', { ip: res.publicIp })}</span>
                {res.publicIp !== connectIp && (
                  <button
                    type="button"
                    className="btn btn--small btn--ghost"
                    onClick={() => onUseIp(res.publicIp)}
                  >
                    {t('server.portCheck.useIp')}
                  </button>
                )}
              </li>
            )}
          </ul>
        </>
      )}

      <details className="portcheck__manual" open={res !== null && res.verdict !== 'open'}>
        <summary>{t('server.portCheck.manual.title')}</summary>
        <p className="hint">{t('server.portCheck.manual.intro')}</p>
        <ol>
          {tl('server.portCheck.manual.steps', {
            ip: res?.publicIp || connectIp || t('server.portCheck.manual.yourIp'),
            port: res?.queryPort ?? '',
          }).map((step) => (
            <li key={step}>{step}</li>
          ))}
        </ol>
        <p className="hint">{t('server.portCheck.manual.tipsTitle')}</p>
        <ul>
          {tl('server.portCheck.manual.tips').map((tip) => (
            <li key={tip}>{tip}</li>
          ))}
        </ul>
      </details>
    </div>
  );
}
