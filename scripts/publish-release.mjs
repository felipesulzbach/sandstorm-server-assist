// Publishes a version as a GitHub Release (ADR-032): the installer, the portable zip, the signed latest.json that
// installed apps read to update themselves, and SHA256SUMS.txt. The release notes come from CHANGELOG.md.
//
// Usage: npm run release:publish                 → version from desktop/wails.json
//        npm run release:publish -- 0.2.0         → another version (its files must be in desktop/build/dist)
//        npm run release:publish -- --dry-run     → checks everything and prints the notes, publishes nothing
// Before: npm run build:installer && npm run package, commit and push (the release tag points to the pushed commit).
// Needs: the GitHub CLI (gh) signed in with access to the repository, and the update signing key.
import { spawnSync } from 'node:child_process';
import { createHash } from 'node:crypto';
import { existsSync, readFileSync, writeFileSync } from 'node:fs';
import { join, resolve } from 'node:path';

const ROOT = resolve(import.meta.dirname, '..');
const DIST = join(ROOT, 'desktop', 'build', 'dist');
const APP = 'sandstorm-server-assist';
// Same repository as the Go updater (adapters/updater.DefaultBaseURL).
const REPO = 'felipesulzbach/sandstorm-server-assist';

const args = process.argv.slice(2);
const dryRun = args.includes('--dry-run');
const version =
  args.find((a) => !a.startsWith('--')) ??
  JSON.parse(readFileSync(join(ROOT, 'desktop', 'wails.json'), 'utf8')).info.productVersion;
const tag = `v${version}`;

function fail(msg) {
  console.error(`✖ ${msg}`);
  process.exit(1);
}

function run(cmd, cmdArgs, { allowFail = false } = {}) {
  const r = spawnSync(cmd, cmdArgs, {
    cwd: ROOT,
    encoding: 'utf8',
    stdio: ['ignore', 'pipe', 'pipe'],
  });
  if (r.error) fail(`${cmd} not found (${r.error.message})`);
  if (r.status !== 0 && !allowFail) fail(`${cmd} ${cmdArgs.join(' ')} failed:\n${r.stderr}`);
  return { ok: r.status === 0, out: (r.stdout ?? '').trim() };
}

if (!/^\d+\.\d+\.\d+$/.test(version)) fail(`Invalid version: ${version}`);

// ---- 1. Files ----
const exe = `${APP}-${version}-windows-setup.exe`;
const zip = `${APP}-${version}-windows-portable.zip`;
for (const f of [exe, zip]) {
  if (!existsSync(join(DIST, f)))
    fail(`Missing ${join(DIST, f)} (run "npm run build:installer" and "npm run package").`);
}
const sha = (f) =>
  createHash('sha256')
    .update(readFileSync(join(DIST, f)))
    .digest('hex');
const sums = `${sha(exe)}  ${exe}\n${sha(zip)}  ${zip}\n`;
writeFileSync(join(DIST, 'SHA256SUMS.txt'), sums);

// Signed manifest (same script and checks as "npm run update:manifest").
const m = spawnSync(process.execPath, [join(ROOT, 'scripts', 'update-manifest.mjs'), version], {
  cwd: ROOT,
  stdio: 'inherit',
});
if (m.status !== 0) fail('Could not create latest.json');
const manifest = JSON.parse(readFileSync(join(DIST, 'latest.json'), 'utf8'));

// ---- 2. Notes from CHANGELOG.md ----
const changelog = readFileSync(join(ROOT, 'CHANGELOG.md'), 'utf8').replace(/\r/g, '');
const start = changelog.search(new RegExp(`^## \\[${version.replace(/\./g, '\\.')}\\]`, 'm'));
if (start < 0) fail(`CHANGELOG.md has no "## [${version}]" section.`);
const rest = changelog.slice(start);
const next = rest.slice(3).search(/^## \[/m);
const section = (next < 0 ? rest : rest.slice(0, next + 3)).trim();
const notes = `${section.split('\n').slice(1).join('\n').trim()}

### Download

- **Installer** (recommended, updates itself): \`${exe}\`
- **Portable** (no install, no automatic updates): \`${zip}\`

Windows may show a SmartScreen warning because the app is not code-signed: click "More info" › "Run anyway".

### SHA-256

\`\`\`
${sums.trim()}
\`\`\`
`;
const notesFile = join(DIST, 'release-notes.md');
writeFileSync(notesFile, notes);

// ---- 3. Repository state ----
const origin = run('git', ['remote', 'get-url', 'origin']).out;
if (!origin.includes(REPO)) fail(`The origin remote is not ${REPO}: ${origin}`);
if (run('git', ['status', '--porcelain']).out)
  fail('There are uncommitted changes. Commit and push first.');
run('git', ['fetch', '-q', 'origin']);
const head = run('git', ['rev-parse', 'HEAD']).out;
if (!run('git', ['merge-base', '--is-ancestor', head, 'origin/main'], { allowFail: true }).ok)
  fail('The current commit is not on origin/main yet. Push it first.');
if (run('gh', ['release', 'view', tag, '--repo', REPO], { allowFail: true }).ok)
  fail(`Release ${tag} already exists.`);

if (dryRun) {
  console.log(`\n--- notes for ${tag} ---\n${notes}`);
  console.log(`✔ Dry run: ${tag} ready to publish from ${head.slice(0, 7)} (nothing sent).`);
  process.exit(0);
}

// ---- 4. Publish ----
// gh uploads the assets to a draft and publishes it at the end, so "latest" never points to a release without
// its files.
console.log(`Publishing ${tag} on GitHub...`);
run('gh', [
  'release',
  'create',
  tag,
  join(DIST, exe),
  join(DIST, zip),
  join(DIST, 'latest.json'),
  join(DIST, 'SHA256SUMS.txt'),
  '--repo',
  REPO,
  '--target',
  head,
  '--title',
  `Sandstorm Server Assist ${version}`,
  '--notes-file',
  notesFile,
  '--latest',
]);

// ---- 5. Check from the internet, at the same address the app uses ----
const latestUrl = `https://github.com/${REPO}/releases/latest/download/latest.json`;
const published = JSON.parse(run('curl', ['-fsSL', latestUrl]).out);
if (published.version !== version || published.signature !== manifest.signature)
  fail(`${latestUrl} is not the manifest of ${version}`);
console.log(`✔ Published: https://github.com/${REPO}/releases/tag/${tag}`);
console.log(`  installed apps will see ${version} at ${latestUrl}`);
