import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest';
import { I18nProvider } from '../i18n';
import type { AppInfo } from './backend';

let serverUpdate = { installed: '100', latest: '101', available: true, checkedAt: '' };

vi.mock('./backend', () => ({
  backend: {
    serverUpdate: () => Promise.resolve(serverUpdate),
    updateState: () => Promise.resolve({ status: 'idle' }),
    openLink: () => Promise.resolve(),
  },
  windowControls: { isMaximised: () => Promise.resolve(false) },
  on: () => () => {},
}));

import { TitleBar } from './TitleBar';

const show = (info: Partial<AppInfo>, onHelp = vi.fn(), onInstall = vi.fn()) =>
  render(
    <I18nProvider initial="en-US">
      <TitleBar info={info as AppInfo} onHelp={onHelp} onInstall={onInstall} />
    </I18nProvider>,
  );

// Pedidos de 14/11/2026: versão e data sempre à vista; aviso de versão nova do servidor dedicado.
describe('TitleBar', () => {
  it('shows the version with its release date and the server update notice', async () => {
    const onHelp = vi.fn();
    const onInstall = vi.fn();
    show({ appVersion: '0.3.0', released: '2026-11-14' }, onHelp, onInstall);
    const version = screen.getByRole('button', { name: /v0\.3\.0/ });
    expect(version).toHaveTextContent('Nov 14, 2026');
    expect(version).toHaveAttribute('title', 'Version 0.3.0, released on Nov 14, 2026');
    await userEvent.click(version);
    expect(onHelp).toHaveBeenCalled();

    const notice = await screen.findByRole('button', { name: /Server update available/ });
    expect(notice).toHaveAttribute('title', expect.stringContaining('installed build 100'));
    await userEvent.click(notice);
    expect(onInstall).toHaveBeenCalled();
  });

  it('shows no date in development and no notice without a newer server build', async () => {
    serverUpdate = { installed: '101', latest: '101', available: false, checkedAt: '' };
    show({ appVersion: '0.1.0-dev', released: '' });
    expect(screen.getByRole('button', { name: 'v0.1.0-dev' })).toBeInTheDocument();
    await new Promise((r) => setTimeout(r, 0));
    expect(screen.queryByRole('button', { name: /Server update available/ })).toBeNull();
  });
});
