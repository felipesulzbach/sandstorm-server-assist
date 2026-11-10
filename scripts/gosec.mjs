// Análise estática de segurança do código Go do app e do verificador de portas (processos, caminhos, segredos no código,
// permissões de arquivo, etc.). Uso: npm run audit:code
// Falso positivo: `// #nosec G<regra> -- motivo` na linha, sempre com o motivo.
import { spawnSync } from 'node:child_process';

const TOOL = 'github.com/securego/gosec/v2/cmd/gosec@v2.30.0';

// O app (desktop) e o verificador de portas (checker, ADR-021) são módulos Go separados.
const results = ['desktop', 'checker'].map((cwd) =>
  spawnSync('go', ['run', TOOL, '-quiet', '-exclude-generated', '-exclude-dir=frontend', './...'], {
    cwd,
    stdio: 'inherit',
  }),
);

if (results.every((r) => r.status === 0)) {
  console.log('✔ gosec: nenhum problema encontrado');
  process.exit(0);
}
console.error('✖ gosec encontrou problemas (veja acima)');
process.exit(1);
