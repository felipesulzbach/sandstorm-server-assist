import { useEffect, useState } from 'react';
import { backend, on, type ServerUpdate } from './backend';

/**
 * Versão nova do servidor dedicado (o Go pergunta ao SteamCMD ao abrir e a cada 6 horas; evento "server:update").
 * A comparação com a versão instalada é refeita a cada minuto: depois de atualizar, o aviso some sozinho.
 */
export function useServerUpdate(): ServerUpdate | null {
  const [u, setU] = useState<ServerUpdate | null>(null);
  useEffect(() => {
    let alive = true;
    const load = () =>
      backend
        .serverUpdate()
        .then((v) => alive && setU(v))
        .catch(() => {});
    load();
    const off = on<ServerUpdate>('server:update', setU);
    const id = window.setInterval(load, 60_000);
    return () => {
      alive = false;
      off();
      window.clearInterval(id);
    };
  }, []);
  return u;
}
