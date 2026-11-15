import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest';
import { I18nProvider } from '../../i18n';
import type { Profile } from '../../shared/backend';

vi.mock('../../shared/backend', () => ({
  backend: {
    rules: () =>
      Promise.resolve({
        groups: ['coop'],
        rules: [
          {
            id: 'aiDifficulty',
            group: 'coop',
            key: 'AIDifficulty',
            kind: 'float',
            min: 0,
            max: 1,
            default: '0.5',
          },
          {
            id: 'minimumEnemies',
            group: 'coop',
            key: 'MinimumEnemies',
            kind: 'int',
            min: 0,
            max: 100,
            default: '6',
          },
        ],
      }),
    serverRules: () => Promise.resolve({ aiDifficulty: '0.1', minimumEnemies: '8' }),
  },
}));

import { RulesTab } from './RulesTab';

const ui = (
  rules: Record<string, string> | null,
  set: (r: Record<string, string> | null) => void,
) => (
  <I18nProvider initial="en-US">
    <RulesTab form={{ serverDir: 'x', rules } as unknown as Profile} set={set} error={null} />
  </I18nProvider>
);

// Regras que o servidor já tem (ADR-040): com a opção desligada, a aba oferece usá-las com um clique.
describe('RulesTab', () => {
  it('offers the rules the server already has when the app does not manage them', async () => {
    const set = vi.fn();
    render(ui(null, set));
    expect(await screen.findByText(/already sets 2 game rules/)).toBeInTheDocument();
    await userEvent.click(screen.getByRole('button', { name: 'Use the server rules' }));
    expect(set).toHaveBeenCalledWith({ aiDifficulty: '0.1', minimumEnemies: '8' });
  });

  // Pedido de 14/11/2026: nome traduzido com a chave original entre parênteses, e busca pelos dois.
  it('shows the Game.ini key and filters by name or key', async () => {
    render(ui({}, () => {}));
    expect(await screen.findByText('(AIDifficulty)')).toBeInTheDocument();
    expect(screen.getByText('(MinimumEnemies)')).toBeInTheDocument();
    const search = screen.getByRole('searchbox', { name: 'Search the game rules' });
    await userEvent.type(search, 'aidiff');
    expect(screen.getByText('(AIDifficulty)')).toBeInTheDocument();
    expect(screen.queryByText('(MinimumEnemies)')).toBeNull();
    await userEvent.clear(search);
    await userEvent.type(search, 'minimum enem');
    expect(screen.getByText('(MinimumEnemies)')).toBeInTheDocument();
    expect(screen.queryByText('(AIDifficulty)')).toBeNull();
  });

  it('does not offer them when the app already manages the rules', async () => {
    render(ui({ aiDifficulty: '0.3' }, () => {}));
    expect(await screen.findByRole('checkbox')).toBeChecked();
    expect(screen.queryByRole('button', { name: 'Use the server rules' })).toBeNull();
  });
});
