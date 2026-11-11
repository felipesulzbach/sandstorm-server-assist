import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest';
import { I18nProvider } from '../../i18n';
import type { ModRow } from '../../shared/backend';

const downloadNewMods = vi.fn(() => Promise.resolve());
vi.mock('../../shared/backend', () => ({
  backend: { downloadNewMods: () => downloadNewMods() },
}));

import { PendingMods } from './ModsPage';

const row = (id: number, name: string, installed: boolean, live = '') =>
  ({ id, name, size: 1, onDisk: 0, installed, subscribed: true, live, done: 0 }) as ModRow;

const mods = [
  row(1, 'Border', true),
  row(2, 'Operation Metro', false),
  row(3, 'Outback', false),
  row(4, 'Estates', false, 'installing'),
];

const ui = (running: boolean, list = mods) => (
  <I18nProvider initial="en-US">
    <PendingMods mods={list} running={running} />
  </I18nProvider>
);

// Mods inscritos e ainda não baixados (ADR-038): o servidor baixa ao iniciar.
describe('PendingMods', () => {
  beforeEach(() => downloadNewMods.mockClear());

  it('lists the mods not downloaded yet and starts the server', async () => {
    render(ui(false));
    expect(screen.getByText(/Not downloaded yet: Operation Metro, Outback\./)).toBeInTheDocument();
    await userEvent.click(screen.getByRole('button', { name: 'Download now' }));
    expect(downloadNewMods).toHaveBeenCalledTimes(1);
    expect(await screen.findByText(/The server is starting and downloads/)).toBeInTheDocument();
  });

  it('asks before restarting a running server', async () => {
    render(ui(true));
    await userEvent.click(screen.getByRole('button', { name: 'Restart and download now' }));
    expect(downloadNewMods).not.toHaveBeenCalled();
    expect(screen.getByText(/players leave for a moment/)).toBeInTheDocument();
    await userEvent.click(screen.getByRole('button', { name: 'Restart now' }));
    expect(downloadNewMods).toHaveBeenCalledTimes(1);
  });

  it('shows nothing when every subscribed mod is on disk', () => {
    const { container } = render(ui(false, [row(1, 'Border', true)]));
    expect(container).toBeEmptyDOMElement();
  });
});
