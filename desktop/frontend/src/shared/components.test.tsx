import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest';
import { I18nProvider } from '../i18n';

vi.mock('./backend', () => ({
  backend: { openLink: vi.fn(), copyText: vi.fn(), setSecret: vi.fn(), revealSecret: vi.fn() },
}));

import { ErrorCard, Guide } from './components';

const wrap = (ui: React.ReactNode, lang: 'en-US' | 'ru' = 'en-US') => render(<I18nProvider initial={lang}>{ui}</I18nProvider>);

describe('ErrorCard', () => {
  it('shows what happened, cause and solution', () => {
    wrap(<ErrorCard error={{ code: 'mod-conflict', detail: 'Serial size mismatch', links: ['modio-subscriptions'] }} />);
    expect(screen.getByText('Mods in conflict')).toBeInTheDocument();
    expect(screen.getByText('What happened')).toBeInTheDocument();
    expect(screen.getByText('Possible cause')).toBeInTheDocument();
    expect(screen.getByText('How to fix')).toBeInTheDocument();
    expect(screen.getByText('Serial size mismatch')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /mod.io: subscriptions/ })).toBeInTheDocument();
  });

  it('falls back to the generic problem', () => {
    wrap(<ErrorCard error={{ code: 'never-seen' }} />);
    expect(screen.getByText('Something went wrong')).toBeInTheDocument();
  });

  it('shows the field message for validation errors', () => {
    wrap(<ErrorCard error={{ code: 'invalid-field', field: 'admins', msg: 'validation.steamId', params: { value: '123' } }} />);
    expect(screen.getByText(/Invalid SteamID64: 123/)).toBeInTheDocument();
  });
});

describe('Guide', () => {
  it('is closed by default and opens on click', async () => {
    wrap(<Guide id="gslt" links={['steam-gslt']} />);
    expect(screen.queryByText(/581320/)).not.toBeInTheDocument();
    await userEvent.click(screen.getByRole('button', { name: /How do I get this/ }));
    expect(screen.getByText(/581320/)).toBeInTheDocument();
  });
});
