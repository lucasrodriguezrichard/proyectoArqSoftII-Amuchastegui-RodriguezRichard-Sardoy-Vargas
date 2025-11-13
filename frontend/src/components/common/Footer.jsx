import { APP_NAME } from '../../utils/constants';

export const Footer = () => (
  <footer className="border-t border-white/40 bg-white/80 backdrop-blur-xl dark:border-white/10 dark:bg-white/5">
    <div className="mx-auto flex max-w-6xl flex-col gap-3 px-4 py-10 text-center text-sm text-slate-600 dark:text-slate-300 sm:flex-row sm:items-center sm:justify-between">
      <span className="font-light tracking-wide">&copy; {new Date().getFullYear()} {APP_NAME}. Todos los derechos reservados.</span>
      <div className="flex items-center justify-center gap-3 text-xs uppercase tracking-[0.3em] text-slate-400 dark:text-slate-500">
        <span>Experiencia</span>
        <span className="h-px w-10 bg-slate-300 dark:bg-slate-600" />
        <span>Premium</span>
      </div>
    </div>
  </footer>
);
