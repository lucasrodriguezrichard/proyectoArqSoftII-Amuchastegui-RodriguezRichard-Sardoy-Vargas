import { APP_NAME } from '../../utils/constants';

export const Footer = () => (
  <footer className="border-t border-slate-100 bg-white/80">
    <div className="mx-auto flex max-w-6xl flex-col gap-1 px-4 py-6 text-center text-sm text-slate-500 sm:flex-row sm:items-center sm:justify-between">
      <span>&copy; {new Date().getFullYear()} {APP_NAME}. Todos los derechos reservados.</span>
      <span className="text-xs text-slate-400">
        Construido con React, Tailwind, React Query y los microservicios de ArqSoft II.
      </span>
    </div>
  </footer>
);
