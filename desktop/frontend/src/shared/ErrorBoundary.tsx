import { Component, type ErrorInfo, type ReactNode } from 'react';
import { AlertTriangle } from 'lucide-react';
import { useI18n } from '../i18n';
import { backend } from './backend';

/**
 * Proteção contra erro de tela (ADR-025). Em vez de o app ficar preto:
 *  - scope "page": só a página mostra o aviso, com "Tentar de novo" e "Voltar ao Painel"; o menu continua;
 *  - scope "app": a moldura do app falhou; aviso em tela cheia com "Recarregar o app".
 * O erro vai para os Registos do app (o Go mascara segredos). resetKey (a página) limpa o erro ao trocar de página.
 */
type Props = {
  scope: 'page' | 'app';
  resetKey?: string;
  onHome?: () => void;
  children: ReactNode;
};
type State = { error: Error | null; key?: string };

export class ErrorBoundary extends Component<Props, State> {
  state: State = { error: null, key: this.props.resetKey };

  static getDerivedStateFromError(error: Error): Partial<State> {
    return { error };
  }

  static getDerivedStateFromProps(props: Props, state: State): Partial<State> | null {
    return props.resetKey !== state.key ? { error: null, key: props.resetKey } : null;
  }

  componentDidCatch(error: Error, info: ErrorInfo) {
    const where = `${this.props.scope}${this.props.resetKey ? ` ${this.props.resetKey}` : ''}`;
    // Primeiras linhas da pilha dos componentes: dizem qual parte da tela quebrou.
    const stack = (info.componentStack ?? '').trim().split('\n').slice(0, 4).join(' <');
    try {
      void backend
        .reportScreenError(where, `${error.name}: ${error.message} <${stack}`)
        .catch(() => {});
    } catch {
      /* sem ponte com o Go (testes): só o aviso na tela */
    }
  }

  render() {
    if (!this.state.error) return this.props.children;
    return (
      <ScreenError
        scope={this.props.scope}
        detail={`${this.state.error.name}: ${this.state.error.message}`}
        onRetry={() => this.setState({ error: null })}
        onHome={this.props.onHome}
      />
    );
  }
}

function ScreenError({
  scope,
  detail,
  onRetry,
  onHome,
}: {
  scope: 'page' | 'app';
  detail: string;
  onRetry: () => void;
  onHome?: () => void;
}) {
  const { t } = useI18n();
  return (
    <div className={scope === 'app' ? 'screen-error__full' : 'page'}>
      <section className="card screen-error" role="alert">
        <h2 className="card__title">
          <AlertTriangle size={18} aria-hidden />{' '}
          {t(scope === 'app' ? 'error.appTitle' : 'error.screenTitle')}
        </h2>
        <p>{t(scope === 'app' ? 'error.appText' : 'error.screenText')}</p>
        <code className="screen-error__detail">{detail}</code>
        <div className="row row--wrap">
          {scope === 'app' ? (
            <button
              type="button"
              className="btn btn--primary"
              onClick={() => window.location.reload()}
            >
              {t('error.reload')}
            </button>
          ) : (
            <>
              <button type="button" className="btn btn--primary" onClick={onRetry}>
                {t('error.retry')}
              </button>
              {onHome && (
                <button type="button" className="btn" onClick={onHome}>
                  {t('error.home')}
                </button>
              )}
            </>
          )}
        </div>
      </section>
    </div>
  );
}
