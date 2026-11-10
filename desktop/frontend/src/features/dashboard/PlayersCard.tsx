import { useEffect, useState } from 'react';
import { Ban as BanIcon, MessageSquare, RefreshCw, UserX, Users } from 'lucide-react';
import { useI18n } from '../../i18n';
import { backend, type Player, type PlayersView, type Status } from '../../shared/backend';
import { ErrorCard } from '../../shared/components';
import { useAction } from '../../shared/hooks';

/** Durações de banimento oferecidas (minutos; 0 = permanente). */
const BAN_MINUTES = [30, 60, 1440, 10080, 0];

/**
 * Jogadores (ADR-022): lista, expulsar, banir, tirar banimento e mensagem a todos, pelo RCON. Uma conexão por
 * clique (nunca em intervalos). A tela manda só o net ID da lista e o texto digitado; o Go confere.
 */
export function PlayersCard({ status }: { status: Status | null }) {
  const { t } = useI18n();
  const action = useAction();
  const [view, setView] = useState<PlayersView | null>(null);
  const [target, setTarget] = useState<{ player: Player; kind: 'kick' | 'ban' } | null>(null);
  const [reason, setReason] = useState('');
  const [minutes, setMinutes] = useState(60);
  const [say, setSay] = useState('');
  const [done, setDone] = useState<string | null>(null);
  const [rcon, setRcon] = useState(false);
  const ready = status?.phase === 'ready';

  useEffect(() => {
    backend
      .mapChange()
      .then((i) => setRcon(i.rcon))
      .catch(() => setRcon(false));
  }, [status?.phase]);

  if (!rcon) return null;

  const load = async () => {
    const v = await action.run(backend.players);
    if (v) setView(v);
  };

  const confirm = async () => {
    if (!target) return;
    const { player, kind } = target;
    const ok = await action.run(async () => {
      if (kind === 'kick') await backend.kick(player.netId, reason);
      else await backend.ban(player.netId, minutes, reason);
      return true;
    });
    if (ok) {
      setDone(t(kind === 'kick' ? 'players.kicked' : 'players.banned', { name: player.name }));
      setTarget(null);
      setReason('');
      await load();
    }
  };

  const minutesLabel = (m: number) =>
    m === 0
      ? t('players.permanent')
      : m >= 1440
        ? t('players.days', { n: m / 1440 })
        : t('players.minutes', { n: m });

  return (
    <section className="card">
      <h2 className="card__title">{t('players.title')}</h2>
      {action.error && <ErrorCard error={action.error} onClose={() => action.setError(null)} />}
      {done && <p className="banner banner--ok">{done}</p>}

      <div className="row row--wrap">
        <button type="button" className="btn" disabled={!ready || action.busy} onClick={load}>
          {view ? <RefreshCw size={16} aria-hidden /> : <Users size={16} aria-hidden />}{' '}
          {view ? t('players.reload') : t('players.load')}
        </button>
        {!ready && <span className="hint">{t('dashboard.mapChange.waitReady')}</span>}
      </div>

      {view && (
        <>
          {view.players.length === 0 ? (
            <p className="hint">{t('players.empty')}</p>
          ) : (
            <table className="players">
              <thead>
                <tr>
                  <th>{t('players.name')}</th>
                  <th>{t('players.score')}</th>
                  <th aria-label={t('players.actions')} />
                </tr>
              </thead>
              <tbody>
                {view.players.map((p) => (
                  <tr key={`${p.id}-${p.netId}`}>
                    <td>
                      {p.name || '—'}
                      {p.bot && <span className="pill">{t('players.bot')}</span>}
                    </td>
                    <td>{p.score}</td>
                    <td className="players__actions">
                      {!p.bot && (
                        <>
                          <button
                            type="button"
                            className="btn btn--small"
                            disabled={action.busy}
                            onClick={() => setTarget({ player: p, kind: 'kick' })}
                          >
                            <UserX size={14} aria-hidden /> {t('players.kick')}
                          </button>
                          <button
                            type="button"
                            className="btn btn--small btn--danger"
                            disabled={action.busy}
                            onClick={() => setTarget({ player: p, kind: 'ban' })}
                          >
                            <BanIcon size={14} aria-hidden /> {t('players.ban')}
                          </button>
                        </>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}

          {target && (
            <div className="players__confirm">
              <p>
                {t(target.kind === 'kick' ? 'players.confirmKick' : 'players.confirmBan', {
                  name: target.player.name,
                })}
              </p>
              <div className="row row--wrap">
                {target.kind === 'ban' && (
                  <select
                    className="select select--short"
                    aria-label={t('players.duration')}
                    value={minutes}
                    onChange={(e) => setMinutes(Number(e.target.value))}
                  >
                    {BAN_MINUTES.map((m) => (
                      <option key={m} value={m}>
                        {minutesLabel(m)}
                      </option>
                    ))}
                  </select>
                )}
                <input
                  className="input"
                  maxLength={100}
                  value={reason}
                  placeholder={t('players.reason')}
                  onChange={(e) => setReason(e.target.value)}
                />
                <button
                  type="button"
                  className="btn btn--small btn--primary"
                  disabled={action.busy}
                  onClick={confirm}
                >
                  {t('players.confirmYes')}
                </button>
                <button type="button" className="btn btn--small" onClick={() => setTarget(null)}>
                  {t('common.cancel')}
                </button>
              </div>
            </div>
          )}

          <h3 className="form__title">{t('players.bans')}</h3>
          {view.bans.length === 0 ? (
            <p className="hint">{t('players.noBans')}</p>
          ) : (
            <ul className="list">
              {view.bans.map((b) => (
                <li key={b.netId} className="list__row">
                  <code>{b.netId}</code>
                  <span className="muted">{b.text}</span>
                  <button
                    type="button"
                    className="btn btn--small"
                    disabled={action.busy}
                    onClick={async () => {
                      const ok = await action.run(async () => {
                        await backend.unban(b.netId);
                        return true;
                      });
                      if (ok) {
                        setDone(t('players.unbanned'));
                        await load();
                      }
                    }}
                  >
                    {t('players.unban')}
                  </button>
                </li>
              ))}
            </ul>
          )}
        </>
      )}

      <div className="row row--wrap players__say">
        <input
          className="input"
          maxLength={150}
          value={say}
          placeholder={t('players.sayPlaceholder')}
          onChange={(e) => setSay(e.target.value)}
        />
        <button
          type="button"
          className="btn"
          disabled={!ready || !say.trim() || action.busy}
          onClick={async () => {
            const ok = await action.run(async () => {
              await backend.say(say);
              return true;
            });
            if (ok) {
              setSay('');
              setDone(t('players.said'));
            }
          }}
        >
          <MessageSquare size={16} aria-hidden /> {t('players.say')}
        </button>
      </div>
    </section>
  );
}
