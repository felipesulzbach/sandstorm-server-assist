import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest';
import { I18nProvider } from '../i18n';

const report = vi.fn(() => Promise.resolve());
vi.mock('./backend', () => ({
  backend: { reportScreenError: (...a: unknown[]) => report(...(a as [])) },
}));

import { ErrorBoundary } from './ErrorBoundary';

function Broken(): React.ReactNode {
  // Igual ao bug da 0.1.5: lista nula lida como lista.
  const players = null as unknown as unknown[];
  return <p>{players.length}</p>;
}

describe('ErrorBoundary', () => {
  beforeEach(() => {
    vi.spyOn(console, 'error').mockImplementation(() => {});
    report.mockClear();
  });

  it('shows the page warning instead of a blank app and goes back home', async () => {
    const home = vi.fn();
    render(
      <I18nProvider initial="en-US">
        <ErrorBoundary scope="page" resetKey="dashboard" onHome={home}>
          <Broken />
        </ErrorBoundary>
      </I18nProvider>,
    );
    expect(screen.getByRole('alert')).toBeInTheDocument();
    expect(screen.getByText('This screen had an error')).toBeInTheDocument();
    expect(report).toHaveBeenCalledTimes(1);
    await userEvent.click(screen.getByText('Back to the Dashboard'));
    expect(home).toHaveBeenCalled();
  });

  it('clears the error when the page changes', () => {
    const ui = (page: string, child: React.ReactNode) => (
      <I18nProvider initial="en-US">
        <ErrorBoundary scope="page" resetKey={page} onHome={() => {}}>
          {child}
        </ErrorBoundary>
      </I18nProvider>
    );
    const { rerender } = render(ui('dashboard', <Broken />));
    expect(screen.getByRole('alert')).toBeInTheDocument();
    rerender(ui('server', <p>server page</p>));
    expect(screen.getByText('server page')).toBeInTheDocument();
  });

  it('offers reloading when the whole app fails', () => {
    render(
      <I18nProvider initial="pt-PT">
        <ErrorBoundary scope="app">
          <Broken />
        </ErrorBoundary>
      </I18nProvider>,
    );
    expect(screen.getByText('Recarregar a aplicação')).toBeInTheDocument();
  });
});
