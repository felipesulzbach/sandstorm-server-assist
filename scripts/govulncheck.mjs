// Verifica vulnerabilidades conhecidas nas dependências Go do app (e na biblioteca padrão).
// O govulncheck só acusa o que o código realmente chama (análise de alcance).
// Uso: npm run audit:go
// Exceções conscientes: govulncheck-allowlist.json (cada uma com motivo e prazo; vencidas não valem).
import { spawnSync } from 'node:child_process';
import { readFileSync } from 'node:fs';

const TOOL = 'golang.org/x/vuln/cmd/govulncheck@v1.9.0';

const result = spawnSync('go', ['run', TOOL, './...'], { cwd: 'desktop', encoding: 'utf8' });
const out = result.stdout ?? '';
process.stdout.write(out);

// "go run" devolve 1 tanto para "achou vulnerabilidade" quanto para erro: decide pela saída.
if (out.includes('No vulnerabilities found')) {
  console.log('✔ govulncheck: nenhuma vulnerabilidade alcançável');
  process.exit(0);
}
const found = new Set([...out.matchAll(/Vulnerability #\d+: (GO-\d{4}-\d+)/g)].map((m) => m[1]));
if (found.size === 0) {
  process.stderr.write(result.stderr ?? '');
  console.error('✖ govulncheck não rodou (veja o erro acima)');
  process.exit(1);
}

const today = new Date().toISOString().slice(0, 10);
const { allowlist } = JSON.parse(readFileSync('govulncheck-allowlist.json', 'utf8'));
const accepted = new Set(
  allowlist.filter((e) => e.id && e.notes && e.expiry && e.expiry >= today).map((e) => e.id),
);

const failing = [...found].filter((id) => !accepted.has(id));
if (failing.length === 0) {
  console.log(`\n✔ Só vulnerabilidades aceitas na allowlist: ${[...found].join(', ')}`);
  process.exit(0);
}
console.error(`\n✖ Vulnerabilidades alcançáveis sem exceção válida: ${failing.join(', ')}`);
process.exit(1);
