import { render, screen } from '@testing-library/react';
import { vi } from 'vitest';
import { I18nProvider } from '../../i18n';
import type { InstallState } from '../../shared/backend';

vi.mock('../../shared/backend', () => ({ backend: { cancelInstall: () => Promise.resolve() } }));

import { TaskPanel } from './InstallPage';

const base = {
  task: '',
  phase: '',
  percent: 0,
  done: 0,
  total: 0,
  speed: 0,
  eta: 0,
  finished: '',
  last: '',
  validate: false,
  upToDate: false,
  present: false,
} as unknown as InstallState;

const ui = (st: InstallState, which: 'steamcmd' | 'server' = 'server') => (
  <I18nProvider initial="en-US">
    <TaskPanel which={which} st={st} onOpenLogs={() => {}} />
  </I18nProvider>
);

// Andamento e resultado no cartão do botão (11/11/2026: nada aparecia perto dos botões e o erro ficava no topo).
describe('TaskPanel', () => {
  it('shows the task as soon as it starts and scrolls it into view', () => {
    const scroll = vi.fn();
    Element.prototype.scrollIntoView = scroll;
    render(ui({ ...base, task: 'server', last: 'server', phase: 'self-update' } as InstallState));
    expect(screen.getByText('Updating SteamCMD')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();
    expect(scroll).toHaveBeenCalled();
  });

  it('says when the server was already up to date or the files were checked', () => {
    const { rerender } = render(
      ui({ ...base, finished: 'server', last: 'server', upToDate: true } as InstallState),
    );
    expect(screen.getByText(/already has the latest version/)).toBeInTheDocument();
    rerender(ui({ ...base, finished: 'server', last: 'server', validate: true } as InstallState));
    expect(screen.getByText(/All server files were checked/)).toBeInTheDocument();
  });

  it('shows the error in the card of the task and nothing in the other card', () => {
    const st = {
      ...base,
      last: 'server',
      error: { code: 'busy', msg: '', field: '' },
    } as unknown as InstallState;
    const { container, rerender } = render(ui(st));
    expect(screen.getByText('Wait for the current task')).toBeInTheDocument();
    rerender(ui(st, 'steamcmd'));
    expect(container).toBeEmptyDOMElement();
  });

  it('explains when SteamCMD was already there', () => {
    render(
      ui(
        { ...base, finished: 'steamcmd', last: 'steamcmd', present: true } as InstallState,
        'steamcmd',
      ),
    );
    expect(screen.getByText(/SteamCMD is already here/)).toBeInTheDocument();
  });
});
