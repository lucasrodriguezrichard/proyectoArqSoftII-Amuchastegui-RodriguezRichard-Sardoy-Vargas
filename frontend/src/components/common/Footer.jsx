import { APP_NAME } from '../../utils/constants';

export const Footer = () => (
  <footer className="border-t border-slate-200/50 bg-white/95 backdrop-blur-sm dark:border-slate-700/50 dark:bg-slate-900/90">
    <div className="mx-auto flex max-w-6xl flex-col gap-2 px-4 py-8 text-center text-sm text-slate-600 dark:text-slate-400 sm:flex-row sm:items-center sm:justify-between">
      <span className="font-light tracking-wide">&copy; {new Date().getFullYear()} {APP_NAME}. Todos los derechos reservados.</span>
      <span className="text-xs text-slate-400 dark:text-slate-500 font-light">
        Sistema de reservas de alta calidad
      </span>
    </div>
  </footer>
);
