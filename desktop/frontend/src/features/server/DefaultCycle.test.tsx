import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest';
import { I18nProvider } from '../../i18n';

const defaultMapCycle = vi.fn();
vi.mock('../../shared/backend', () => ({
  backend: { defaultMapCycle: () => defaultMapCycle() },
}));

import { DefaultCycle } from './MapsExtras';

const lines = ['Scenario_Bab_Checkpoint_Security', 'Scenario_Bab_Checkpoint_Insurgents'];
const ui = (current: string[], onFill: (l: string[]) => void) => (
  <I18nProvider initial="en-US">
    <DefaultCycle
      lines={current}
      startScenario="Scenario_Bab_Checkpoint_Security"
      onFill={onFill}
    />
  </I18nProvider>
);

// Rotação padrão com os mapas do jogo (ADR-043).
describe('DefaultCycle', () => {
  beforeEach(() => defaultMapCycle.mockReset());

  it('fills an empty rotation of a new server once', async () => {
    defaultMapCycle.mockResolvedValue({ lines, serverHasFile: false });
    const onFill = vi.fn();
    render(ui([], onFill));
    expect(
      await screen.findByText(/Filled with the default rotation: 2 game maps/),
    ).toBeInTheDocument();
    expect(onFill).toHaveBeenCalledWith(lines);
  });

  it('keeps an empty rotation when the server has its own MapCycle.txt', async () => {
    defaultMapCycle.mockResolvedValue({ lines, serverHasFile: true });
    const onFill = vi.fn();
    render(ui([], onFill));
    await userEvent.click(
      await screen.findByRole('button', { name: /Use the default rotation \(2 game maps\)/ }),
    );
    expect(onFill).toHaveBeenCalledTimes(1);
  });

  it('asks before replacing a rotation that has lines', async () => {
    defaultMapCycle.mockResolvedValue({ lines, serverHasFile: true });
    const onFill = vi.fn();
    render(ui(['Scenario_Tell_Push_Security'], onFill));
    await userEvent.click(await screen.findByRole('button', { name: /Use the default rotation/ }));
    expect(onFill).not.toHaveBeenCalled();
    expect(
      screen.getByText('Replace the current rotation with the default one?'),
    ).toBeInTheDocument();
    await userEvent.click(screen.getByRole('button', { name: 'Yes' }));
    expect(onFill).toHaveBeenCalledWith(lines);
  });
});
