import { render, screen, within } from '@testing-library/react';
import { vi } from 'vitest';
import { I18nProvider } from '../../i18n';

vi.mock('../../shared/backend', () => ({
  backend: {
    mutators: () =>
      Promise.resolve([
        { id: 'BulletSponge', label: 'Bullet Sponge', changes: [], conflicts: [] },
        { id: 'Ultralethal', label: 'Ultralethal', changes: [], conflicts: [] },
        {
          id: 'AllYouCanEat',
          label: 'All You Can Eat',
          changes: ['supply'],
          conflicts: ['Strapped'],
        },
        { id: 'Strapped', label: 'Strapped', changes: ['supply'], conflicts: ['AllYouCanEat'] },
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

  // Mutators que trocam os mesmos pontos de abastecimento iniciais: um sobrescreve o outro (ADR-045).
  it('turns off the game mutators that change the same setting as a chosen one', async () => {
    render(
      <I18nProvider initial="en-US">
        <MutatorPicker value={['AllYouCanEat']} onChange={() => {}} />
      </I18nProvider>,
    );
    const strapped = (await screen.findByText('Strapped', { selector: 'span' })).closest('label')!;
    expect(within(strapped).getByRole('checkbox')).toBeDisabled();
    expect(within(strapped).getByText(/same setting as All You Can Eat/)).toBeInTheDocument();
    const sponge = screen.getByText('Bullet Sponge').closest('label')!;
    expect(within(sponge).getByRole('checkbox')).toBeEnabled();
    const eat = screen.getByText('All You Can Eat', { selector: 'span' }).closest('label')!;
    expect(within(eat).getByRole('checkbox')).toBeEnabled();
  });

  it('flags two conflicting mutators kept from an older profile without locking them', async () => {
    render(
      <I18nProvider initial="en-US">
        <MutatorPicker value={['AllYouCanEat', 'Strapped']} onChange={() => {}} />
      </I18nProvider>,
    );
    const strapped = (await screen.findByText('Strapped', { selector: 'span' })).closest('label')!;
    expect(within(strapped).getByRole('checkbox')).toBeEnabled();
    expect(within(strapped).getByText('conflict')).toBeInTheDocument();
  });
});
