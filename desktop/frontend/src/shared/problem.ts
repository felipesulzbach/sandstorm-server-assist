// Erros do Go chegam como JSON na mensagem da promessa rejeitada (domain/problems.Error, ADR-005).
export interface AppError {
  code: string;
  field?: string;
  msg?: string;
  params?: Record<string, string>;
  detail?: string;
  links?: string[];
}

export function toAppError(e: unknown): AppError {
  const raw = typeof e === 'string' ? e : e instanceof Error ? e.message : String(e ?? '');
  try {
    const v = JSON.parse(raw) as Partial<AppError>;
    if (v && typeof v.code === 'string') return v as AppError;
  } catch {
    // não é JSON: erro solto do Wails ou do navegador
  }
  return { code: 'unknown', detail: raw };
}

/** Erro de um campo específico (para mostrar ao lado do campo). */
export function fieldError(err: AppError | null, field: string): AppError | null {
  return err && err.field === field ? err : null;
}
