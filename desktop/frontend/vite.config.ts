import { defineConfig, type Plugin } from 'vitest/config';
import react from '@vitejs/plugin-react';

// CSP só no build: em dev o Vite injeta scripts inline (recarga rápida) que ela bloquearia.
// O front só fala com o Go pela ponte do Wails (não usa rede): nada de origem externa.
const CSP = [
  "default-src 'self'",
  "script-src 'self'",
  "style-src 'self' 'unsafe-inline'",
  "img-src 'self' data:",
  "font-src 'self' data:",
  "connect-src 'self'",
  "object-src 'none'",
  "base-uri 'none'",
  "form-action 'none'",
].join('; ');

function cspOnBuild(): Plugin {
  return {
    name: 'csp-on-build',
    apply: 'build',
    transformIndexHtml: (html) =>
      html.replace(
        '<head>',
        `<head>\n    <meta http-equiv="Content-Security-Policy" content="${CSP}" />`,
      ),
  };
}

export default defineConfig({
  plugins: [react(), cspOnBuild()],
  server: { port: 5173 },
  build: { sourcemap: false, chunkSizeWarningLimit: 900 },
  test: { environment: 'jsdom', globals: true, setupFiles: ['./src/test-setup.ts'] },
});
