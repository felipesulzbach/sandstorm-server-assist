// Regras de segurança para o código TypeScript do front. Usado por `npm run audit:code`.
// Separado do lint normal para que a skill security-audit rode só as regras de segurança.
import tseslint from 'typescript-eslint';
import security from 'eslint-plugin-security';

export default tseslint.config(
  { ignores: ['**/dist/**', '**/node_modules/**', '**/*.test.*', '**/wailsjs/**'] },
  {
    files: ['**/*.{ts,tsx}'],
    languageOptions: { parser: tseslint.parser },
    plugins: { security },
    rules: {
      ...security.configs.recommended.rules,
      // Acesso a objeto por chave dinâmica gera muitos falsos positivos em código tipado.
      'security/detect-object-injection': 'off',
      'no-eval': 'error',
      'no-implied-eval': 'error',
      'no-new-func': 'error',
    },
  },
);
