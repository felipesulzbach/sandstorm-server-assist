import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest';
import { I18nProvider } from '../../i18n';

const state = {
  offer: true,
  missing: true,
  hasSecrets: true,
  others: 1,
  mods: 6,
  modioSession: true,
};
const startFresh = vi.fn((_clear: boolean) => Promise.resolve({ hostname: 'Default' }));
const clearModio = vi.fn((_mods: boolean, _session: boolean) => Promise.resolve());
vi.mock('../../shared/backend', () => ({
  backend: {
    freshState: () => Promise.resolve(state),
    startFresh: (c: boolean) => startFresh(c),
    clearModio: (m: boolean, s: boolean) => clearModio(m, s),
  },
}));

import { FreshStart } from './FreshStart';

const ui = (onDone = vi.fn()) => (
  <I18nProvider initial="en-US">
    <FreshStart checkKey="x" onDone={onDone} />
  </I18nProvider>
);

// 14/11/2026: a pasta do servidor foi movida para começar do zero e o app guardava tudo do servidor antigo.
describe('FreshStart', () => {
  beforeEach(() => {
    startFresh.mockClear();
    clearModio.mockClear();
  });

  it('offers to start from scratch, asking before clearing the passwords', async () => {
    const onDone = vi.fn();
    render(ui(onDone));
    expect(await screen.findByText('The server folder no longer exists.')).toBeInTheDocument();
    await userEvent.click(screen.getByRole('button', { name: /Start from scratch/ }));
    expect(screen.getByText(/other 1 profile\(s\) use the same passwords/)).toBeInTheDocument();
    // Senhas marcadas por padrão; mods e mod.io desmarcados (ficam no Windows, fora da pasta do servidor).
    const [secrets, mods, session] = screen.getAllByRole('checkbox');
    expect(secrets).toBeChecked();
    expect(mods).not.toBeChecked();
    expect(session).not.toBeChecked();
    await userEvent.click(secrets!);
    await userEvent.click(mods!);
    await userEvent.click(screen.getByRole('button', { name: 'Clear and start from scratch' }));
    expect(startFresh).toHaveBeenCalledWith(false);
    expect(clearModio).toHaveBeenCalledWith(true, false);
    expect(onDone).toHaveBeenCalled();
    expect(await screen.findByText(/the settings are back to the defaults/)).toBeInTheDocument();
  });

  it('can be dismissed without clearing anything', async () => {
    const { container } = render(ui());
    await userEvent.click(await screen.findByRole('button', { name: 'Keep the settings' }));
    expect(container).toBeEmptyDOMElement();
    expect(startFresh).not.toHaveBeenCalled();
  });
});
