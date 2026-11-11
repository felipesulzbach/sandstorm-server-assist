import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest';
import { I18nProvider } from '../../i18n';
import type { ModRow } from '../../shared/backend';

const deleteMod = vi.fn((_id: number, _redownload: boolean) => Promise.resolve());
vi.mock('../../shared/backend', () => ({
  backend: { deleteMod: (id: number, r: boolean) => deleteMod(id, r) },
}));

import { ModLine } from './ModsPage';

const mod = {
  id: 6096282,
  name: 'Estates',
  size: 1000,
  onDisk: 1000,
  installed: true,
  subscribed: true,
  live: '',
  done: 0,
} as ModRow;

const ui = (running: boolean, onDeleted = () => {}, m: ModRow = mod) => (
  <I18nProvider initial="en-US">
    <table>
      <tbody>
        <ModLine m={m} hasSubs running={running} onDeleted={onDeleted} />
      </tbody>
    </table>
  </I18nProvider>
);

// Apagar o mod baixado para o servidor baixar de novo (ADR-037).
describe('ModLine', () => {
  beforeEach(() => deleteMod.mockClear());

  it('deletes and downloads again, also with the server running', async () => {
    const onDeleted = vi.fn();
    render(ui(true, onDeleted));
    await userEvent.click(
      screen.getByRole('button', { name: 'Delete the downloaded files of Estates' }),
    );
    expect(screen.getByRole('button', { name: 'Only delete' })).toBeDisabled();
    expect(
      screen.getByText(/stops the server, deletes the mod and starts the server/),
    ).toBeInTheDocument();

    await userEvent.click(screen.getByRole('button', { name: 'Delete and download again' }));
    expect(deleteMod).toHaveBeenCalledWith(6096282, true);
    expect(onDeleted).toHaveBeenCalled();
    expect(await screen.findByText(/Estates deleted. The server is starting/)).toBeInTheDocument();
  });

  it('only deletes with the server stopped', async () => {
    render(ui(false));
    await userEvent.click(
      screen.getByRole('button', { name: 'Delete the downloaded files of Estates' }),
    );
    await userEvent.click(screen.getByRole('button', { name: 'Only delete' }));
    expect(deleteMod).toHaveBeenCalledWith(6096282, false);
    expect(await screen.findByText(/on its next start/)).toBeInTheDocument();
  });

  it('has no delete button while the mod downloads or when it is not on disk', () => {
    const { rerender } = render(ui(false, () => {}, { ...mod, live: 'installing' } as ModRow));
    expect(screen.queryByRole('button', { name: /Delete the downloaded/ })).toBeNull();
    rerender(ui(false, () => {}, { ...mod, installed: false } as ModRow));
    expect(screen.queryByRole('button', { name: /Delete the downloaded/ })).toBeNull();
  });
});
