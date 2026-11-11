import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest';
import { I18nProvider } from '../../i18n';

const checkMods = vi.fn();
const checkCrash = vi.fn();
vi.mock('../../shared/backend', () => ({
  backend: {
    checkMods: () => checkMods(),
    checkCrash: (text: string) => checkCrash(text),
    openLink: () => Promise.resolve(),
  },
}));

import { ModHealth } from './ModHealth';

const ui = () => (
  <I18nProvider initial="en-US">
    <ModHealth refreshKey={1} />
  </I18nProvider>
);

// Mods com problema (ADR-039): conflitos entre mods e o mod de um crash.
describe('ModHealth', () => {
  beforeEach(() => {
    checkMods.mockReset();
    checkCrash.mockReset();
  });

  it('lists mod pairs with the same game files in different versions', async () => {
    checkMods.mockResolvedValue({
      mods: 3,
      unreadable: [],
      conflicts: [
        {
          a: { id: 2123731, name: 'Containment' },
          b: { id: 2853082, name: 'Containment Storm' },
          count: 12,
          example: '/Game/Subway/M_Wall',
        },
      ],
    });
    render(ui());
    expect(await screen.findByText('Containment (#2123731)')).toBeInTheDocument();
    expect(screen.getByText('Containment Storm (#2853082)')).toBeInTheDocument();
    expect(screen.getByText(/12 file\(s\) in different versions/)).toBeInTheDocument();
    expect(screen.getByText('e.g. /Game/Subway/M_Wall')).toBeInTheDocument();
  });

  it('says when there is no conflict', async () => {
    checkMods.mockResolvedValue({ mods: 7, unreadable: null, conflicts: null });
    render(ui());
    expect(
      await screen.findByText('No conflicts found between your 7 installed mods.'),
    ).toBeInTheDocument();
  });

  it('finds the mod of a pasted crash message', async () => {
    checkMods.mockResolvedValue({ mods: 1, unreadable: [], conflicts: [] });
    checkCrash.mockResolvedValue({
      mods: 1,
      sizeMismatch: true,
      findings: [{ object: '/Game/Props/SM_Railing', mods: [{ id: 5748299, name: 'Daybreak' }] }],
    });
    render(ui());
    await userEvent.click(screen.getByText(/Find the mod$/, { selector: 'summary' }));
    await userEvent.type(
      screen.getByRole('textbox'),
      'StaticMesh /Game/Props/SM_Railing: Serial size mismatch',
    );
    await userEvent.click(screen.getByRole('button', { name: 'Find the mod' }));
    expect(checkCrash).toHaveBeenCalledWith(
      'StaticMesh /Game/Props/SM_Railing: Serial size mismatch',
    );
    expect(await screen.findByText(/shipped by Daybreak \(#5748299\)/)).toBeInTheDocument();
    expect(screen.getByText(/made for another version of the game/)).toBeInTheDocument();
  });
});
