import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import '@fontsource-variable/inter';
import './styles/app.css';
import './styles/themes.css';
import { App } from './App';
import { applyAppearance, loadAppearance } from './shared/appearance';

// Aparência antes do primeiro desenho (sem piscar a cor padrão).
applyAppearance(loadAppearance());

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>,
);
