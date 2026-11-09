// Análise estática de segurança do código Go do app (processos, caminhos, segredos no código,
// permissões de arquivo, etc.). Uso: npm run audit:code
// Falso positivo: `// #nosec G<regra> -- motivo` na linha, sempre com o motivo.
import { spawnSync } from 'node:child_process';

const TOOL = 'github.com/securego/gosec/v2/cmd/gosec@v2.30.0';

const result = spawnSync(
  'go',
  ['run', TOOL, '-quiet', '-exclude-generated', '-exclude-dir=frontend', './...'],
  { cwd: 'desktop', stdio: 'inherit' },
);

if (result.status === 0) {
  console.log('✔ gosec: nenhum problema encontrado');
  process.exit(0);
}
console.error('✖ gosec encontrou problemas (veja acima)');
process.exit(1);
