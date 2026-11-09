import { useCallback, useEffect, useState } from 'react';
import { backend, on, type UpdateState } from './backend';

/** Situação da atualização automática (o Go procura sozinho; a tela só mostra e manda instalar). */
export function useUpdate() {
  const [state, setState] = useState<UpdateState | null>(null);
  const [error, setError] = useState<unknown>(null);

  useEffect(() => {
    let alive = true;
    backend
      .updateState()
      // Só preenche se nenhum evento chegou antes (o evento é sempre mais novo).
      .then((s) => alive && setState((prev) => prev ?? s))
      .catch(() => {});
    const off = on<UpdateState>('update:state', setState);
    return () => {
      alive = false;
      off();
    };
  }, []);

  const check = useCallback(async () => {
    setError(null);
    try {
      setState(await backend.checkUpdates());
    } catch (e) {
      setError(e);
    }
  }, []);

  const install = useCallback(async () => {
    setError(null);
    try {
      await backend.installUpdate();
    } catch (e) {
      setError(e);
    }
  }, []);

  return { state, error, check, install };
}
