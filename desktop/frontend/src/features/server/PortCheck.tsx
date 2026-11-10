import { useState } from 'react';
import { AlertTriangle, CheckCircle2, Globe, Info, XCircle } from 'lucide-react';
import { useI18n } from '../../i18n';
import { backend, type PortCheck as Result } from '../../shared/backend';
import { ErrorCard } from '../../shared/components';
import { useAction } from '../../shared/hooks';

/**
 * Teste das portas pela internet (ADR-021). O Go pede ao verificador do projeto, que testa de fora o IP desta
 * máquina: a porta de consulta (o servidor responde?) e a do RCON (não deveria estar aberta). A tela só mostra.
 */
export function PortCheck({
  connectIp,
  onUseIp,
}: {
  connectIp: string;
  onUseIp: (ip: string) => void;
}) {
  const { t } = useI18n();
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
          <Globe size={16} aria-hidden />{' '}
          {action.busy ? t('server.portCheck.running') : t('server.portCheck.button')}
        </button>
        <span className="hint">{t('server.portCheck.hint')}</span>
      </div>

      {action.error && <ErrorCard error={action.error} onClose={() => action.setError(null)} />}

      {res && (
        <ul className="portcheck__list">
          <li className={res.queryOpen ? 'is-ok' : 'is-bad'}>
            {res.queryOpen ? (
              <CheckCircle2 size={16} aria-hidden />
            ) : (
              <XCircle size={16} aria-hidden />
            )}
            <span>
              {res.queryOpen
                ? t('server.portCheck.queryOpen', {
                    port: res.queryPort,
                    name: res.serverName,
                    map: res.map,
                  })
                : t('server.portCheck.queryClosed', { port: res.queryPort })}
            </span>
          </li>
          <li className="is-info">
            <Info size={16} aria-hidden />
            <span>{t('server.portCheck.gamePort', { port: res.gamePort })}</span>
          </li>
          {res.rconTested && (
            <li className={res.rconExposed ? 'is-bad' : 'is-ok'}>
              {res.rconExposed ? (
                <AlertTriangle size={16} aria-hidden />
              ) : (
                <CheckCircle2 size={16} aria-hidden />
              )}
              <span>
                {t(
                  res.rconExposed ? 'server.portCheck.rconExposed' : 'server.portCheck.rconClosed',
                  {
                    port: res.rconPort,
                  },
                )}
              </span>
            </li>
          )}
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
        </ul>
      )}
      {res && !res.queryOpen && <p className="hint">{t('server.portCheck.closedHelp')}</p>}
    </div>
  );
}
