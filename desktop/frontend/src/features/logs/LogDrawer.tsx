import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useRef,
  useState,
  type ReactNode,
} from 'react';
import { ChevronDown, ChevronUp, FolderOpen, ScrollText, Trash2 } from 'lucide-react';
import { useI18n } from '../../i18n';
import { backend, on, type LogLine, type LogRun } from '../../shared/backend';
import { CopyButton } from '../../shared/components';
import { bytes } from '../../shared/format';

type Filter = 'all' | 'important' | 'warning' | 'error';

/** "live" é a partida atual (linhas ao vivo); os outros valores são ids de partidas guardadas pelo servidor. */
const LIVE = 'live';

interface LogsCtx {
  open: boolean;
  setOpen: (v: boolean) => void;
  /** Abre já filtrando erros (botão "Abrir logs" de um erro). */
  showErrors: () => void;
}

const Ctx = createContext<LogsCtx>({ open: false, setOpen: () => {}, showErrors: () => {} });
export const useLogs = () => useContext(Ctx);

const MAX = 6000;

/** Guarda as linhas e o estado da gaveta para o app inteiro (as linhas chegam mesmo com a gaveta fechada). */
export function LogsProvider({ children }: { children: ReactNode }) {
  const [open, setOpen] = useState(false);
  const [lines, setLines] = useState<LogLine[]>([]);
  const [filter, setFilter] = useState<Filter>('important');

  useEffect(() => {
    let alive = true;
    backend
      .recentLogs(MAX)
      .then((l) => alive && setLines(l ?? []))
      .catch(() => {});
    const off = on<LogLine[]>('log:lines', (batch) =>
      setLines((cur) => {
        const next = cur.concat(batch);
        return next.length > MAX ? next.slice(next.length - MAX) : next;
      }),
    );
    return () => {
      alive = false;
      off();
    };
  }, []);

  const value = useMemo<LogsCtx>(
    () => ({
      open,
      setOpen,
      showErrors: () => {
        setFilter('error');
        setOpen(true);
      },
    }),
    [open],
  );

  return (
    <Ctx.Provider value={value}>
      {children}
      <LogDrawer
        open={open}
        setOpen={setOpen}
        lines={lines}
        clear={() => {
          backend.clearLogs().catch(() => {});
          setLines([]);
        }}
        filter={filter}
        setFilter={setFilter}
      />
    </Ctx.Provider>
  );
}

/** Partida atual: a de número mais alto (cada início do servidor abre uma, ADR-018). */
function currentRun(lines: LogLine[]): number {
  let run = 0;
  for (const l of lines) if (l.run > run) run = l.run;
  return run;
}

function LogDrawer({
  open,
  setOpen,
  lines,
  clear,
  filter,
  setFilter,
}: {
  open: boolean;
  setOpen: (v: boolean) => void;
  lines: LogLine[];
  clear: () => void;
  filter: Filter;
  setFilter: (f: Filter) => void;
}) {
  const { t, date } = useI18n();
  const [query, setQuery] = useState('');
  const [follow, setFollow] = useState(true);
  const [view, setView] = useState<string>(LIVE);
  const [runs, setRuns] = useState<LogRun[]>([]);
  const [saved, setSaved] = useState<LogLine[]>([]);
  const [loading, setLoading] = useState(false);
  const listRef = useRef<HTMLDivElement>(null);

  const run = useMemo(() => currentRun(lines), [lines]);
  const live = useMemo(() => lines.filter((l) => l.run === run), [lines, run]);
  const liveStart = live.find((l) => l.kind === 'run');
  const shown = view === LIVE ? live : saved;

  const errors = useMemo(() => live.filter((l) => l.level === 'error').length, [live]);
  const shownErrors = useMemo(() => shown.filter((l) => l.level === 'error').length, [shown]);

  // Partidas guardadas pelo servidor. Com uma partida ao vivo, o log atual do servidor é ela (fica de fora).
  const loadRuns = useCallback(() => {
    backend
      .logRuns()
      .then((r) => setRuns((r ?? []).filter((x) => !(x.current && run > 0))))
      .catch(() => setRuns([]));
  }, [run]);

  useEffect(() => {
    if (open) loadRuns();
  }, [open, loadRuns]);

  const choose = (id: string) => {
    setView(id);
    if (id === LIVE) {
      setFollow(true);
      return;
    }
    setFollow(false);
    setLoading(true);
    backend
      .readLogRun(id)
      .then((l) => setSaved(l ?? []))
      .catch(() => setSaved([]))
      .finally(() => setLoading(false));
  };

  const visible = useMemo(() => {
    const q = query.trim().toLowerCase();
    return shown.filter((l) => {
      if (l.kind === 'run') return true;
      if (filter === 'error' && l.level !== 'error') return false;
      if (filter === 'warning' && l.level === 'info') return false;
      if (
        filter === 'important' &&
        (l.noise || (l.level === 'info' && !l.event && l.source === 'server'))
      )
        return false;
      return !q || l.text.toLowerCase().includes(q);
    });
  }, [shown, filter, query]);

  useEffect(() => {
    if (open && follow && listRef.current) listRef.current.scrollTop = listRef.current.scrollHeight;
  }, [visible, open, follow]);

  return (
    <section className={`logdrawer ${open ? 'is-open' : ''}`} aria-label={t('logs.title')}>
      <div className="logdrawer__bar">
        <button
          type="button"
          className="logdrawer__toggle"
          aria-expanded={open}
          onClick={() => setOpen(!open)}
        >
          <ScrollText size={15} aria-hidden />
          {t('logs.button')}
          {errors > 0 && (
            <span className="pill pill--danger">{t('logs.errors', { n: errors })}</span>
          )}
          {open ? <ChevronDown size={15} aria-hidden /> : <ChevronUp size={15} aria-hidden />}
        </button>
        {open && (
          <div className="logdrawer__tools">
            <select
              className="input input--small logdrawer__run"
              aria-label={t('logs.run')}
              value={view}
              onFocus={loadRuns}
              onChange={(e) => choose(e.target.value)}
            >
              <option value={LIVE}>
                {liveStart
                  ? t('logs.runLive', { time: date(liveStart.time) })
                  : t('logs.runSession')}
              </option>
              {runs.map((r) => (
                <option key={r.id} value={r.id}>
                  {t(r.current ? 'logs.runLast' : 'logs.runAt', {
                    time: date(r.started),
                    size: bytes(r.size, t),
                  })}
                </option>
              ))}
            </select>
            <div className="segmented segmented--small" role="radiogroup">
              {(['all', 'important', 'warning', 'error'] as Filter[]).map((f) => (
                <button
                  key={f}
                  type="button"
                  role="radio"
                  aria-checked={filter === f}
                  className={filter === f ? 'is-active' : ''}
                  onClick={() => setFilter(f)}
                >
                  {t(`logs.filter.${f}`)}
                </button>
              ))}
            </div>
            <input
              className="input input--small"
              type="search"
              placeholder={t('logs.search')}
              value={query}
              onChange={(e) => setQuery(e.target.value)}
            />
            <label className="check">
              <input
                type="checkbox"
                checked={follow}
                disabled={view !== LIVE}
                onChange={(e) => setFollow(e.target.checked)}
              />
              {t('logs.follow')}
            </label>
            <CopyButton
              text={visible
                .slice(-500)
                .map((l) => l.text)
                .join('\n')}
            />
            <button
              type="button"
              className="icon-btn"
              title={t('logs.openFolder')}
              aria-label={t('logs.openFolder')}
              onClick={() => backend.openLogFolder().catch(() => {})}
            >
              <FolderOpen size={16} />
            </button>
            {view === LIVE && (
              <button
                type="button"
                className="icon-btn"
                title={t('logs.clear')}
                aria-label={t('logs.clear')}
                onClick={clear}
              >
                <Trash2 size={16} />
              </button>
            )}
          </div>
        )}
      </div>
      {open && (
        <div className="logdrawer__list" ref={listRef}>
          {view !== LIVE && (
            <p className="logdrawer__saved">
              {loading ? t('logs.loading') : t('logs.savedRun', { n: shownErrors })}{' '}
              <button type="button" className="link" onClick={() => choose(LIVE)}>
                {t('logs.backToLive')}
              </button>
            </p>
          )}
          {visible.length === 0 ? (
            <p className="muted logdrawer__empty">{t('logs.empty')}</p>
          ) : (
            visible.slice(-1500).map((l) =>
              l.kind === 'run' ? (
                <div key={l.seq} className="logline logline--run" role="separator">
                  {t('logs.runStarted', { time: date(l.time) })}
                </div>
              ) : (
                <div
                  key={l.seq}
                  className={`logline logline--${l.level} ${l.event ? 'logline--event' : ''}`}
                >
                  <span className="logline__src">{t(`logs.source.${l.source}`)}</span>
                  <span className="logline__text">{l.text}</span>
                </div>
              ),
            )
          )}
          <p className="muted logdrawer__note">{t('logs.masked')}</p>
        </div>
      )}
    </section>
  );
}
