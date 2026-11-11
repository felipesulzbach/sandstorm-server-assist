import { useEffect, useState } from 'react';
import {
  Archive,
  Download,
  Info,
  MessageSquareWarning,
  LayoutDashboard,
  Package,
  Server,
  Settings as SettingsIcon,
  type LucideIcon,
} from 'lucide-react';
import { detectLang, I18nProvider, isLang, useI18n, type Lang } from './i18n';
import { backend, type AppInfo } from './shared/backend';
import { ErrorBoundary } from './shared/ErrorBoundary';
import { TitleBar } from './shared/TitleBar';
import { useServerStatus } from './shared/hooks';
import { LogsProvider } from './features/logs/LogDrawer';
import { DashboardPage } from './features/dashboard/DashboardPage';
import { InstallPage } from './features/install/InstallPage';
import { ServerPage } from './features/server/ServerPage';
import { ModsPage } from './features/mods/ModsPage';
import { BackupPage } from './features/backup/BackupPage';
import { SettingsPage } from './features/settings/SettingsPage';
import { AboutPage } from './features/about/AboutPage';
import { FeedbackPage } from './features/feedback/FeedbackPage';
import { LegalGate, useTerms } from './features/legal/Legal';

export type PageId =
  'dashboard' | 'install' | 'server' | 'mods' | 'backup' | 'settings' | 'feedback' | 'about';

const MAIN: { id: PageId; icon: LucideIcon }[] = [
  { id: 'dashboard', icon: LayoutDashboard },
  { id: 'install', icon: Download },
  { id: 'server', icon: Server },
  { id: 'mods', icon: Package },
  { id: 'backup', icon: Archive },
];
const BOTTOM: { id: PageId; icon: LucideIcon }[] = [
  { id: 'settings', icon: SettingsIcon },
  { id: 'feedback', icon: MessageSquareWarning },
  { id: 'about', icon: Info },
];

export function App() {
  const [lang, setLang] = useState<Lang | null>(null);
  useEffect(() => {
    backend
      .settings()
      .then((s) => setLang(isLang(s.language) ? s.language : detectLang(navigator.language)))
      .catch(() => setLang(detectLang(navigator.language)));
  }, []);
  if (!lang) return null;
  return (
    <I18nProvider initial={lang}>
      <LogsProvider>
        {/* Erro na moldura do app (menu, barra de título): aviso em tela cheia com "Recarregar". */}
        <ErrorBoundary scope="app">
          <Shell />
        </ErrorBoundary>
      </LogsProvider>
    </I18nProvider>
  );
}

function Shell() {
  const { t } = useI18n();
  const [page, setPage] = useState<PageId>('dashboard');
  const [info, setInfo] = useState<AppInfo | null>(null);
  const status = useServerStatus();
  const [terms, setTerms] = useTerms();

  useEffect(() => {
    backend
      .info()
      .then(setInfo)
      .catch(() => {});
  }, []);

  const navItem = ({ id, icon: Icon }: { id: PageId; icon: LucideIcon }) => (
    <li key={id}>
      <button
        type="button"
        className={`nav__item ${page === id ? 'is-active' : ''}`}
        aria-current={page === id ? 'page' : undefined}
        title={t(`nav.${id}`)}
        onClick={() => setPage(id)}
      >
        <Icon size={18} aria-hidden />
        <span className="nav__label">{t(`nav.${id}`)}</span>
        {id === 'dashboard' && status && (
          <span className={`dot dot--${status.phase}`} aria-hidden />
        )}
      </button>
    </li>
  );

  return (
    <div className="shell">
      <TitleBar onHelp={() => setPage('about')} />
      <div className="app">
        <aside className="sidebar">
          <nav aria-label={t('nav.label')}>
            <ul className="nav">{MAIN.map(navItem)}</ul>
          </nav>
          <ul className="nav nav--bottom">{BOTTOM.map(navItem)}</ul>
          <footer className="sidebar__foot">
            {status && <div>{t(`phase.${status.phase}`)}</div>}
            {info && <div>{t('app.version', { version: info.appVersion })}</div>}
          </footer>
        </aside>
        <main className="content">
          {/* Erro numa página mostra um aviso nela; o resto do app continua (antes ficava tudo preto). */}
          <ErrorBoundary scope="page" resetKey={page} onHome={() => setPage('dashboard')}>
            {page === 'dashboard' && <DashboardPage go={setPage} status={status} />}
            {page === 'install' && <InstallPage status={status} />}
            {page === 'server' && <ServerPage />}
            {page === 'mods' && <ModsPage status={status} />}
            {page === 'backup' && <BackupPage status={status} />}
            {page === 'settings' && <SettingsPage />}
            {page === 'feedback' && <FeedbackPage />}
            {page === 'about' && <AboutPage info={info} terms={terms} />}
          </ErrorBoundary>
        </main>
      </div>
      {terms && !terms.accepted && <LegalGate terms={terms} onAccepted={setTerms} />}
    </div>
  );
}
