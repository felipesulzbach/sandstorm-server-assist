import { useCallback, useEffect, useState } from 'react';
import { backend, on, type InstallState, type Status } from './backend';
import { useServers } from './servers';
import { toAppError, type AppError } from './problem';

/** Estado do servidor (evento "server:status" do Go). */
export function useServerStatus(): Status | null {
  const [st, setSt] = useState<Status | null>(null);
  // Só o estado do servidor escolhido (um aviso atrasado de outro servidor não troca a tela, ADR-048).
  const { selected } = useServers();
  useEffect(() => {
    let alive = true;
    const mine = (s: Status) => !selected || !s.server || s.server === selected;
    const load = () =>
      backend
        .status()
        .then((s) => alive && mine(s) && setSt(s))
        .catch(() => {});
    void load();
    const off = on<Status>('server:status', (s) => {
      if (mine(s)) setSt(s);
    });
    // Reserva: se um aviso se perder (a tela recarregada no meio), o estado se acerta sozinho em poucos segundos.
    const poll = setInterval(load, 5000);
    return () => {
      alive = false;
      clearInterval(poll);
      off();
    };
  }, [selected]);
  return st;
}

/** Estado da instalação (evento "install:state"). */
export function useInstallState(): InstallState | null {
  const [st, setSt] = useState<InstallState | null>(null);
  useEffect(() => {
    let alive = true;
    backend
      .installState()
      .then((s) => alive && setSt(s))
      .catch(() => {});
    const off = on<InstallState>('install:state', (s) => setSt(s));
    return () => {
      alive = false;
      off();
    };
  }, []);
  return st;
}

/**
 * Roda uma ação do Go guardando "ocupado" e o erro do catálogo. run devolve o resultado ou undefined se falhou.
 */
export function useAction() {
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<AppError | null>(null);
  const run = useCallback(async <R>(fn: () => Promise<R>): Promise<R | undefined> => {
    setBusy(true);
    setError(null);
    try {
      return await fn();
    } catch (e) {
      setError(toAppError(e));
      return undefined;
    } finally {
      setBusy(false);
    }
  }, []);
  return { busy, error, setError, run };
}
