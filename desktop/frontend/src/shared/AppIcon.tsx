/** Ícone do app: mira sobre um escudo (desenho próprio, sem marcas do jogo). */
export function AppIcon({ className, size = 18 }: { className?: string; size?: number }) {
  return (
    <svg className={className} width={size} height={size} viewBox="0 0 24 24" aria-hidden="true">
      <path d="M12 2 4 5v6c0 5 3.4 9.3 8 11 4.6-1.7 8-6 8-11V5l-8-3z" fill="var(--accent)" />
      <circle cx="12" cy="11.5" r="4.2" fill="none" stroke="var(--accent-text)" strokeWidth="1.6" />
      <path d="M12 5.8v3M12 14.2v3M6.3 11.5h3M14.7 11.5h3" stroke="var(--accent-text)" strokeWidth="1.6" strokeLinecap="round" />
    </svg>
  );
}
