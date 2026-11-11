import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest';
import { I18nProvider } from '../../i18n';

const lines = [
  '(Scenario="Scenario_Farmhouse_Checkpoint_Security",Lighting="Day")',
  '(Scenario="Scenario_Outback_Checkpoint_Security",Lighting="Day")',
  '(Scenario="Scenario_RedwoodReborn_Checkpoint",Lighting="Winter")',
];
const scenarioOf = (l: string) => /Scenario_\w+/.exec(l)?.[0] ?? '';

vi.mock('../../shared/backend', () => ({
  backend: {
    // O Go marca as linhas sem mapa no jogo base nem nos mods instalados.
    mapCycleLines: (ls: string[]) =>
      Promise.resolve(
        ls.map((l) => ({
          scenario: scenarioOf(l),
          lighting: '',
          missing: !l.includes('Farmhouse'),
        })),
      ),
    setMapCycleLighting: (ls: string[]) => Promise.resolve(ls),
  },
}));

import { CycleLighting } from './MapsExtras';

// Rotação com mapas de mods não instalados (11/11/2026): aviso, "Remover" por linha e "Remover estes mapas".
describe('CycleLighting', () => {
  it('warns about maps without their mod and removes them', async () => {
    const onChange = vi.fn();
    render(
      <I18nProvider initial="en-US">
        <CycleLighting lines={lines} onChange={onChange} />
      </I18nProvider>,
    );
    expect(
      await screen.findByText(/not in the base game or in any installed mod/),
    ).toBeInTheDocument();
    expect(screen.getAllByText('mod not installed')).toHaveLength(2);

    const [first] = screen.getAllByRole('button', { name: 'Remove' });
    await userEvent.click(first!);
    expect(onChange).toHaveBeenLastCalledWith([lines[0], lines[2]]);

    await userEvent.click(screen.getByRole('button', { name: 'Remove these maps' }));
    expect(onChange).toHaveBeenLastCalledWith([lines[0]]);
  });
});
