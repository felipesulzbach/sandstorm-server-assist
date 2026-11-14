import { render, screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest';
import { I18nProvider } from '../../i18n';

const lines = [
  'Scenario_Farmhouse_Checkpoint_Security',
  'Scenario_Farmhouse_Checkpoint_Insurgents',
  'Scenario_Outback_Checkpoint_Security',
];

const backend = vi.hoisted(() => ({
  // O Go agrupa a rotação por mapa e marca o Outback como sem mod instalado.
  mapCycleView: vi.fn(() =>
    Promise.resolve({
      groups: [
        {
          map: 'Farmhouse',
          label: 'Farmhouse',
          modName: '',
          entries: [
            {
              index: 0,
              scenario: lines[0],
              mode: 'Checkpoint_Security',
              lighting: 'Day',
              missing: false,
              raw: '',
            },
            {
              index: 1,
              scenario: lines[1],
              mode: 'Checkpoint_Insurgents',
              lighting: 'Night',
              missing: false,
              raw: '',
            },
          ],
        },
        {
          map: 'Outback',
          label: 'Outback',
          modName: 'Outback',
          entries: [
            {
              index: 2,
              scenario: lines[2],
              mode: 'Checkpoint_Security',
              lighting: 'Day',
              missing: true,
              raw: '',
            },
          ],
        },
      ],
      maps: [
        {
          map: 'Farmhouse',
          label: 'Farmhouse',
          modId: 0,
          modName: '',
          scenarios: [
            { id: lines[0], mode: 'Checkpoint_Security' },
            { id: lines[1], mode: 'Checkpoint_Insurgents' },
            { id: 'Scenario_Farmhouse_Push_Security', mode: 'Push_Security' },
          ],
        },
        {
          map: 'Outback',
          label: 'Outback',
          modId: 7,
          modName: 'Outback',
          scenarios: [{ id: lines[2], mode: 'Checkpoint_Security' }],
        },
        {
          map: 'Daybreak',
          label: 'Daybreak',
          modId: 5,
          modName: 'Daybreak',
          scenarios: [
            { id: 'Scenario_Daybreak_Checkpoint_Security', mode: 'Checkpoint_Security' },
            { id: 'Scenario_Daybreak_Checkpoint_Insurgents', mode: 'Checkpoint_Insurgents' },
            { id: 'Scenario_Daybreak_Push_Security', mode: 'Push_Security' },
          ],
        },
      ],
    }),
  ),
  removeFromMapCycle: vi.fn((ls: string[], idx: number[]) =>
    Promise.resolve(ls.filter((_, i) => !idx.includes(i))),
  ),
  addToMapCycle: vi.fn((ls: string[], _after: number, ids: string[]) =>
    Promise.resolve([...ls, ...ids]),
  ),
  moveMapCycleGroup: vi.fn((ls: string[]) => Promise.resolve(ls)),
  setMapCycleLighting: vi.fn((ls: string[]) => Promise.resolve(ls)),
}));
vi.mock('../../shared/backend', () => ({ backend }));

import { CycleEditor } from './CycleEditor';

const renderEditor = (onChange = vi.fn()) => {
  render(
    <I18nProvider initial="en-US">
      <CycleEditor lines={lines} onChange={onChange} />
    </I18nProvider>,
  );
  return onChange;
};

describe('CycleEditor', () => {
  it('shows the rotation in one block per map', async () => {
    renderEditor();
    const farm = (await screen.findByText('Farmhouse', { selector: 'strong' })).closest('li')!;
    expect(within(farm).getByText('Checkpoint Security')).toBeInTheDocument();
    expect(within(farm).getByText('Checkpoint Insurgents')).toBeInTheDocument();
    expect(within(farm).getByRole('button', { name: 'Move Farmhouse up' })).toBeDisabled();
    const outback = screen.getByText('Outback', { selector: 'strong' }).closest('li')!;
    expect(within(outback).getByText('Mod: Outback')).toBeInTheDocument();
    expect(within(outback).getByText('mod not installed')).toBeInTheDocument();
  });

  // Rotação com mapas de mods não instalados (11/11/2026): aviso e "Remover estes mapas".
  it('removes the maps without their mod, a scenario or a whole map', async () => {
    const onChange = renderEditor();
    await userEvent.click(await screen.findByRole('button', { name: 'Remove these maps' }));
    expect(backend.removeFromMapCycle).toHaveBeenLastCalledWith(lines, [2]);
    expect(onChange).toHaveBeenLastCalledWith(lines.slice(0, 2));

    await userEvent.click(
      screen.getByRole('button', { name: 'Remove Farmhouse from the rotation' }),
    );
    expect(backend.removeFromMapCycle).toHaveBeenLastCalledWith(lines, [0, 1]);

    await userEvent.click(screen.getByRole('button', { name: `Remove ${lines[1]}` }));
    expect(backend.removeFromMapCycle).toHaveBeenLastCalledWith(lines, [1]);
  });

  it('adds another scenario of a map right after its block', async () => {
    renderEditor();
    const pick = await screen.findByRole('combobox', { name: 'Add a scenario of Farmhouse' });
    // Checkpoint Security já está de dia e Insurgents de noite: cada um só aparece na outra iluminação.
    expect(within(pick).queryByText('Checkpoint Security · Day')).toBeNull();
    expect(within(pick).getByText('Checkpoint Security · Night')).toBeInTheDocument();
    expect(within(pick).queryByText('Checkpoint Insurgents · Night')).toBeNull();
    expect(within(pick).getByText('Checkpoint Insurgents · Day')).toBeInTheDocument();
    await userEvent.selectOptions(pick, 'Scenario_Farmhouse_Push_Security|Night');
    expect(backend.addToMapCycle).toHaveBeenLastCalledWith(
      lines,
      1,
      ['Scenario_Farmhouse_Push_Security'],
      ['Night'],
    );
  });

  it('adds a mod map with both Checkpoint sides suggested', async () => {
    const onChange = renderEditor();
    await userEvent.click(await screen.findByRole('button', { name: 'Add map' }));
    // Mapas que já estão na rotação (todos ou parte dos cenários): ✓, quantos cenários e sem poder escolher; mais
    // cenários deles entram pelo bloco do mapa. O nome do mod igual ao do mapa não se repete.
    expect(
      screen.getByRole('option', { name: '✓ Farmhouse — in the rotation (2 of 3 scenarios)' }),
    ).toBeDisabled();
    expect(
      screen.getByRole('option', { name: '✓ Outback — in the rotation (1 of 1 scenarios)' }),
    ).toBeDisabled();
    expect(screen.getByRole('option', { name: 'Daybreak' })).toBeEnabled();
    await userEvent.selectOptions(screen.getByRole('combobox', { name: 'Map' }), '5:Daybreak');
    expect(screen.getByRole('checkbox', { name: 'Checkpoint Security' })).toBeChecked();
    expect(screen.getByRole('checkbox', { name: 'Checkpoint Insurgents' })).toBeChecked();
    expect(screen.getByRole('checkbox', { name: 'Push Security' })).not.toBeChecked();
    // Dia e Noite vêm marcados: cada cenário entra de dia e de noite.
    expect(screen.getByRole('checkbox', { name: 'Day' })).toBeChecked();
    expect(screen.getByRole('checkbox', { name: 'Night' })).toBeChecked();
    await userEvent.click(screen.getByRole('button', { name: 'Add to the rotation' }));
    expect(backend.addToMapCycle).toHaveBeenLastCalledWith(
      lines,
      -1,
      ['Scenario_Daybreak_Checkpoint_Security', 'Scenario_Daybreak_Checkpoint_Insurgents'],
      ['Day', 'Night'],
    );
    expect(onChange).toHaveBeenCalled();
  });

  it('says every map is already in the rotation and turns off Add map', async () => {
    const full = await backend.mapCycleView();
    backend.mapCycleView.mockResolvedValueOnce({
      ...full,
      maps: full.maps.filter((m) => m.map !== 'Daybreak'),
    });
    renderEditor();
    expect(
      await screen.findByText(/Every available map is already in the rotation/),
    ).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Add map' })).toBeDisabled();
  });
});
