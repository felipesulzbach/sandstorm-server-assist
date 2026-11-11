import { render, screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest';
import { I18nProvider } from '../../i18n';
import type { Scenario, Status } from '../../shared/backend';

const list: Scenario[] = [
  {
    id: 'Scenario_Estates_Checkpoint_Insurgency',
    map: 'Estates',
    mode: 'Checkpoint_Insurgency',
    source: 'mod',
  },
  { id: 'Scenario_Estates_Push_Security', map: 'Estates', mode: 'Push_Security', source: 'mod' },
  {
    id: 'Scenario_Farmhouse_Checkpoint_Security',
    map: 'Farmhouse',
    mode: 'Checkpoint_Security',
    source: 'server',
  },
  {
    id: 'Scenario_RedwoodReborn_Checkpoint',
    map: 'RedwoodReborn',
    mode: 'Checkpoint',
    source: 'missing',
  },
  {
    id: 'Scenario_TORO_Checkpoint_Insurgents',
    map: 'TORO',
    mode: 'Checkpoint_Insurgents',
    source: 'mod',
  },
  {
    id: 'Scenario_TORO_Checkpoint_Security',
    map: 'TORO',
    mode: 'Checkpoint_Security',
    source: 'mod',
  },
];

vi.mock('../../shared/backend', () => ({
  backend: {
    mapChange: () => Promise.resolve({ rcon: true }),
    scenarios: () => Promise.resolve(list),
    changeMap: () => Promise.resolve(),
  },
}));

import { MapChangeCard } from './MapChangeCard';

const status = {
  phase: 'ready',
  map: 'TORO',
  scenario: 'Scenario_TORO_Checkpoint_Insurgents',
} as unknown as Status;

// Mapa e modo em listas separadas (11/11/2026): a lista única agrupada era difícil de usar.
describe('MapChangeCard', () => {
  it('picks the map first and then one of its modes', async () => {
    render(
      <I18nProvider initial="en-US">
        <MapChangeCard status={status} go={() => {}} />
      </I18nProvider>,
    );
    await userEvent.click(await screen.findByText('Show available maps'));

    const map = await screen.findByLabelText<HTMLSelectElement>('Map');
    const mode = screen.getByLabelText<HTMLSelectElement>('Mode');
    expect(map.value).toBe('TORO');
    expect(
      within(map)
        .getAllByRole('option')
        .map((o) => o.textContent),
    ).toEqual([
      'Estates (mod map)',
      'Farmhouse',
      'Redwood Reborn (mod not loaded)',
      'TORO (mod map) · now',
    ]);
    // Mapa da rotação sem o mod carregado: aparece como referência, sem poder ser escolhido (ADR-030).
    expect(within(map).getByRole('option', { name: /Redwood/ })).toBeDisabled();
    expect(
      screen.getByText(/are in your map rotation, but their mod is not loaded/),
    ).toBeInTheDocument();
    expect(mode.value).toBe('Scenario_TORO_Checkpoint_Insurgents');
    expect(within(mode).getAllByRole('option')).toHaveLength(2);

    // Outro mapa: o primeiro modo dele já vem escolhido.
    await userEvent.selectOptions(map, 'Estates');
    expect(mode.value).toBe('Scenario_Estates_Checkpoint_Insurgency');
    expect(
      within(mode)
        .getAllByRole('option')
        .map((o) => o.textContent),
    ).toEqual(['Checkpoint Insurgency', 'Push Security']);
  });
});
