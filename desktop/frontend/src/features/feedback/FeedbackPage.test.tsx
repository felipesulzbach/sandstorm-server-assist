import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest';
import { I18nProvider } from '../../i18n';

const createReport = vi.fn((..._args: unknown[]) =>
  Promise.reject(
    new Error(
      JSON.stringify({ code: 'invalid-field', field: 'title', msg: 'validation.reportTitle' }),
    ),
  ),
);

vi.mock('../../shared/backend', () => ({
  backend: {
    createReport: (...args: unknown[]) => createReport(...args),
    chooseReportImage: () => Promise.resolve(null),
    clearReportImage: () => Promise.resolve(),
    showReport: () => Promise.resolve(),
  },
}));

import { FeedbackPage } from './FeedbackPage';

// Pedido de 14/11/2026: o título do relato é obrigatório; descrição, log e imagem continuam opcionais.
describe('FeedbackPage', () => {
  it('needs a title before creating the report', async () => {
    render(
      <I18nProvider initial="en-US">
        <FeedbackPage />
      </I18nProvider>,
    );
    const create = screen.getByRole('button', { name: /Prepare and open GitHub/ });
    const title = screen.getByRole('textbox', { name: /Title/ });
    expect(title).toBeRequired();
    expect(create).toBeDisabled();
    await userEvent.type(title, '   ');
    expect(create).toBeDisabled();
    await userEvent.type(title, 'Crash on start');
    expect(create).toBeEnabled();
    // Se o Go recusar (título só com caracteres que a limpeza tira), o erro aparece no próprio campo.
    await userEvent.click(create);
    expect(await screen.findByText('Write a title.')).toBeInTheDocument();
    expect(createReport).toHaveBeenCalledWith('bug', '   Crash on start', '', true);
  });
});
