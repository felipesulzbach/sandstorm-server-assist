import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from 'react';
import { backend, on, type ServerSummary } from './backend';

/**
 * Vários servidores (ADR-048): a lista com o estado de cada um (evento "fleet:status") e o servidor escolhido, que
 * as telas Servidor, Mods, Logs e o Painel mostram. Trocar o escolhido remonta as telas (o App usa o id como key).
 */
type ServersCtx = {
  list: ServerSummary[];
  selected: string;
  loaded: boolean;
  /** Servidores marcados no Painel para o botão principal iniciar juntos (começa com o mostrado). */
  checked: string[];
  toggle: (id: string, on: boolean) => void;
  /** Marca exatamente estes (Selecionar todos / nenhum). */
  setChecked: (ids: string[]) => void;
  select: (id: string) => Promise<void>;
  refresh: () => Promise<void>;
};

const Ctx = createContext<ServersCtx>({
  list: [],
  selected: '',
  loaded: false,
  checked: [],
  toggle: () => {},
  setChecked: () => {},
  select: async () => {},
  refresh: async () => {},
});

export function ServersProvider({ children }: { children: ReactNode }) {
  const [list, setList] = useState<ServerSummary[]>([]);
  const [loaded, setLoaded] = useState(false);
  const [marks, setMarks] = useState<string[] | null>(null);

  const refresh = useCallback(async () => {
    try {
      setList((await backend.servers()) ?? []);
    } catch {
      // sem a lista, o app continua com o servidor escolhido no Go
    } finally {
      setLoaded(true);
    }
  }, []);

  useEffect(() => {
    // Se a lista não chegar (ponte com o Go ainda a abrir), as telas aparecem na mesma: nunca fica tudo preto.
    const fallback = setTimeout(() => setLoaded(true), 3000);
    void refresh().finally(() => clearTimeout(fallback));
    // O "selected" que chega no evento é de quando ele saiu; vale o da lista (muda só por select/refresh).
    const off = on<ServerSummary>('fleet:status', (s) =>
      setList((cur) =>
        cur.some((x) => x.id === s.id)
          ? cur.map((x) => (x.id === s.id ? ({ ...s, selected: x.selected } as ServerSummary) : x))
          : cur,
      ),
    );
    return () => {
      clearTimeout(fallback);
      off();
    };
  }, [refresh]);

  const select = useCallback(
    async (id: string) => {
      await backend.switchProfile(id);
      await refresh();
    },
    [refresh],
  );

  const selected = list.find((s) => s.selected)?.id ?? '';
  // Sem escolha ainda, vale o servidor mostrado; servidores apagados saem da escolha.
  const checked = useMemo(
    () => (marks ?? (selected ? [selected] : [])).filter((id) => list.some((s) => s.id === id)),
    [marks, selected, list],
  );
  const toggle = useCallback(
    (id: string, on: boolean) =>
      setMarks((cur) => {
        const base = cur ?? (selected ? [selected] : []);
        return on ? [...base.filter((x) => x !== id), id] : base.filter((x) => x !== id);
      }),
    [selected],
  );
  const value = useMemo(
    () => ({ list, selected, loaded, checked, toggle, setChecked: setMarks, select, refresh }),
    [list, selected, loaded, checked, toggle, select, refresh],
  );
  return <Ctx.Provider value={value}>{children}</Ctx.Provider>;
}

export function useServers() {
  return useContext(Ctx);
}
