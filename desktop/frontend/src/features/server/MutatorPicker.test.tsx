import { render, screen, within } from '@testing-library/react';
import { vi } from 'vitest';
import { I18nProvider } from '../../i18n';

vi.mock('../../shared/backend', () => ({
  backend: {
    mutators: () =>
      Promise.resolve([
        { id: 'BulletSponge', label: 'Bullet Sponge' },
        { id: 'Ultralethal', label: 'Ultralethal' },
      ]),
    modMutators: () => Promise.resolve([]),
    status: () => Promise.resolve({ refusedMutators: ['Ultralethal'] }),
  },
}));

import { MutatorPicker } from './MapsExtras';

// 13/11/2026: o servidor recusou o Ultralethal no Checkpoint (com e sem o Bullet Sponge); só o log dizia (ADR-041).
describe('MutatorPicker', () => {
  it('marks the mutators the server refused on the current map', async () => {
    render(
      <I18nProvider initial="en-US">
        <MutatorPicker value={['BulletSponge', 'Ultralethal']} onChange={() => {}} />
      </I18nProvider>,
    );
    expect(await screen.findByText(/the server refused: Ultralethal/)).toBeInTheDocument();
    const ultra = screen.getByText('Ultralethal', { selector: 'span' }).closest('label')!;
    expect(within(ultra).getByText('refused by the server')).toBeInTheDocument();
    const sponge = screen.getByText('Bullet Sponge').closest('label')!;
    expect(within(sponge).queryByText('refused by the server')).toBeNull();
  });
});
