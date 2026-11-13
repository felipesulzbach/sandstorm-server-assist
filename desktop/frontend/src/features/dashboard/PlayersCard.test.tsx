import { render, screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest';
import { I18nProvider } from '../../i18n';
import type { Status } from '../../shared/backend';

const p = (id: string, name: string, score: number, bot = false) => ({
  id,
  name,
  netId: bot ? '' : `SteamNWI:7656119000000000${id.slice(-1)}`,
  ip: '',
  score,
  bot,
});

vi.mock('../../shared/backend', () => ({
  backend: {
    mapChange: () => Promise.resolve({ rcon: true }),
    players: () =>
      Promise.resolve({
        players: [
          p('257', 'Alpha', 120),
          p('1', 'Marksman', 3, true),
          p('258', 'Bravo', 625),
          p('2', 'Rifleman', 0, true),
          p('259', 'Charlie', 0),
        ],
        bans: [],
      }),
  },
}));

import { PlayersCard } from './PlayersCard';

const status = { phase: 'ready' } as unknown as Status;

// 12/11/2026: com 3 jogadores e 19 bots o app mostrava só 1 (o listplayers vem numa linha só).
describe('PlayersCard', () => {
  it('lists every player by score, with the bots folded', async () => {
    render(
      <I18nProvider initial="en-US">
        <PlayersCard status={status} />
      </I18nProvider>,
    );
    await userEvent.click(await screen.findByRole('button', { name: /Show players/ }));
    const table = await screen.findByRole('table');
    const names = () =>
      within(table)
        .getAllByRole('row')
        .slice(1)
        .map((r) => within(r).getAllByRole('cell')[0]!.textContent);
    expect(names()).toEqual(['Bravo', 'Alpha', 'Charlie', 'Show the 2 bots']);

    await userEvent.click(screen.getByRole('button', { name: 'Show the 2 bots' }));
    expect(names()).toEqual([
      'Bravo',
      'Alpha',
      'Charlie',
      'Hide the 2 bots',
      'Marksman bot',
      'Rifleman bot',
    ]);
    expect(screen.getAllByRole('button', { name: /Kick/ })).toHaveLength(3);
  });
});
