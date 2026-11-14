import { render, screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest';
import { I18nProvider } from '../../i18n';

const row = (id: string, name: string, phase: string, selected: boolean, extra = {}) => ({
  id,
  name,
  hostname: name,
  instance: id === 'a' ? 1 : 2,
  port: id === 'a' ? 27102 : 27103,
  selected,
  phase,
  map: phase === 'ready' ? 'Hideout' : '',
  since: '',
  external: false,
  conflicts: [],
  ...extra,
});

const startServer = vi.fn((_id: string) => Promise.resolve());
const stopServer = vi.fn((_id: string) => Promise.resolve());
const switchProfile = vi.fn((_id: string) => Promise.resolve({}));
let list = [
  row('a', 'Coop', 'ready', true),
  row('b', 'Versus', 'stopped', false, {
    conflicts: [{ field: 'queryPort', port: 27131, server: 'Coop' }],
  }),
];

vi.mock('../../shared/backend', () => ({
  backend: {
    servers: () => Promise.resolve(list),
    startServer: (id: string) => startServer(id),
    stopServer: (id: string) => stopServer(id),
    switchProfile: (id: string) => switchProfile(id),
  },
  on: () => () => {},
}));

import { ServersProvider, useServers } from '../../shared/servers';
import { ServersCard } from './ServersCard';

// Mostra os servidores marcados (o botão principal do Painel inicia estes).
function Checked() {
  return <output>{useServers().checked.join(',')}</output>;
}

const show = () =>
  render(
    <I18nProvider initial="en-US">
      <ServersProvider>
        <ServersCard go={() => {}} />
        <Checked />
      </ServersProvider>
    </I18nProvider>,
  );

// ADR-048: o Painel lista todos os servidores, com Iniciar/Parar de cada um e Mostrar para escolher o das telas.
describe('ServersCard', () => {
  it('starts, stops and shows each server by id', async () => {
    show();
    const coop = (await screen.findByText('Coop')).closest('li')!;
    const versus = screen.getByText('Versus').closest('li')!;
    expect(within(coop).getByText('Ready to play')).toBeInTheDocument();
    expect(within(versus).getByText(/Port 27131 is also used by Coop/)).toBeInTheDocument();

    await userEvent.click(within(coop).getByRole('button', { name: /Stop/ }));
    expect(stopServer).toHaveBeenCalledWith('a');
    await userEvent.click(within(versus).getByRole('button', { name: /Start/ }));
    expect(startServer).toHaveBeenCalledWith('b');

    expect(within(coop).getByRole('button', { name: /Shown/ })).toBeDisabled();
    list = [row('a', 'Coop', 'ready', false), row('b', 'Versus', 'waiting', true)];
    await userEvent.click(within(versus).getByRole('button', { name: /Show/ }));
    expect(switchProfile).toHaveBeenCalledWith('b');
    expect(await within(versus).findByText('Waiting for the other server')).toBeInTheDocument();
  });

  it('ticks any server, and all of them at once', async () => {
    list = [
      row('a', 'Coop', 'stopped', true),
      row('b', 'Versus', 'stopped', false),
      row('c', 'Hard', 'ready', false),
    ];
    show();
    const coop = await screen.findByRole('checkbox', { name: /Select Coop/ });
    const versus = screen.getByRole('checkbox', { name: /Select Versus/ });
    // O servidor mostrado começa marcado; um que já roda também se marca (para parar junto).
    expect(coop).toBeChecked();
    expect(screen.getByRole('checkbox', { name: /Select Hard/ })).toBeEnabled();
    await userEvent.click(versus);
    expect(screen.getByRole('status')).toHaveTextContent('a,b');
    await userEvent.click(coop);
    expect(screen.getByRole('status')).toHaveTextContent(/^b$/);
    const all = screen.getByRole('checkbox', { name: /Select all/ });
    await userEvent.click(all);
    expect(screen.getByRole('status')).toHaveTextContent('a,b,c');
    await userEvent.click(all);
    expect(screen.getByRole('status')).toHaveTextContent(/^$/);
  });

  it('stays hidden with a single server', async () => {
    list = [row('a', 'Coop', 'stopped', true)];
    show();
    await new Promise((r) => setTimeout(r, 0));
    expect(screen.queryByText('Servers')).not.toBeInTheDocument();
  });
});
