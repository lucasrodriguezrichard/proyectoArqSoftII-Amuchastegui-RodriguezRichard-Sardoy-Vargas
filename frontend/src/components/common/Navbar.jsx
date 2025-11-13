import { Link, NavLink, useNavigate } from 'react-router-dom';
import { LogOut, UserRound, Sun, Moon } from 'lucide-react';
import { APP_NAME } from '../../utils/constants';
import { useAuth } from '../../hooks/useAuth';
import { useTheme } from '../../hooks/useTheme';

const navItems = [
  { to: '/', label: 'Buscar', private: false },
  { to: '/create-reservation', label: 'Crear reserva', private: true },
  { to: '/my-reservations', label: 'Mis reservas', private: true },
  { to: '/admin', label: 'Admin', private: true, adminOnly: true },
];

const linkClasses = ({ isActive }) =>
  `px-3 py-1.5 rounded-full text-sm font-medium transition ${
    isActive
      ? 'bg-white/80 text-slate-900 shadow-sm shadow-white/60 dark:bg-white/20 dark:text-white'
      : 'text-slate-500 hover:text-slate-900 dark:text-slate-400 dark:hover:text-white'
  }`;

export const Navbar = () => {
  const navigate = useNavigate();
  const { isAuthenticated, isAdmin, user, logout } = useAuth();
  const { theme, toggleTheme } = useTheme();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <header className="sticky top-0 z-30 border-b border-white/40 bg-white/70 backdrop-blur-3xl shadow-lg shadow-slate-900/5 dark:border-white/5 dark:bg-slate-900/60">
      <nav className="mx-auto flex max-w-6xl items-center justify-between px-4 py-3">
        <Link to="/" className="group flex items-center gap-2 font-display text-lg font-semibold text-slate-900 dark:text-white">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100" className="h-6 w-6">
            <rect x="20" y="30" width="60" height="40" fill="none" stroke="currentColor" strokeWidth="2.5" rx="3"/>
            <circle cx="50" cy="50" r="12" fill="none" stroke="currentColor" strokeWidth="2"/>
            <line x1="32" y1="45" x2="32" y2="55" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round"/>
            <line x1="68" y1="45" x2="68" y2="55" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round"/>
          </svg>
          {APP_NAME}
        </Link>
        <div className="flex items-center gap-4">
          <button
            onClick={toggleTheme}
            className="rounded-full border border-white/40 bg-white/30 p-2 text-slate-700 shadow-sm shadow-white/30 transition hover:scale-105 hover:bg-white/80 dark:border-white/10 dark:bg-white/5 dark:text-slate-100 dark:hover:bg-white/20"
            aria-label="Toggle theme"
          >
            {theme === 'light' ? <Moon size={18} /> : <Sun size={18} />}
          </button>
          <div className="hidden items-center gap-1 rounded-full border border-white/40 bg-white/50 px-2 py-1 shadow-inner shadow-white/60 dark:border-white/10 dark:bg-white/5 sm:flex">
            {navItems
              .filter((item) => {
                if (item.private && !isAuthenticated) return false;
                if (item.adminOnly && !isAdmin) return false;
                return true;
              })
              .map((item) => (
                <NavLink key={item.to} to={item.to} className={linkClasses}>
                  {item.label}
                </NavLink>
              ))}
          </div>
          {isAuthenticated ? (
            <div className="flex items-center gap-3">
              <span className="hidden text-sm font-medium text-slate-700 dark:text-slate-200 sm:flex sm:flex-col sm:items-end">
                <span className="flex items-center gap-1 text-xs uppercase tracking-wide text-slate-500 dark:text-slate-400">
                  <UserRound size={14} />
                  {user?.role ?? 'user'}
                </span>
                {user?.first_name || user?.username}
              </span>
              <button
                type="button"
                onClick={handleLogout}
                className="inline-flex items-center gap-2 rounded-full border border-white/50 bg-white/70 px-4 py-2 text-sm font-semibold text-slate-700 shadow-sm transition hover:-translate-y-0.5 hover:shadow-lg dark:border-white/10 dark:bg-white/10 dark:text-white dark:hover:bg-white/20"
              >
                <LogOut size={16} />
                Salir
              </button>
            </div>
          ) : (
            <div className="flex items-center gap-2">
              <Link
                to="/login"
                className="rounded-full border border-white/60 px-4 py-2 text-sm font-semibold text-slate-700 transition hover:-translate-y-0.5 hover:bg-white/80 dark:border-white/15 dark:text-white dark:hover:bg-white/10"
              >
                Ingresar
              </Link>
              <Link
                to="/register"
                className="rounded-full bg-gradient-to-r from-primary-500 via-sky-500 to-emerald-400 px-5 py-2 text-sm font-semibold text-white shadow-lg shadow-primary-500/40 transition hover:-translate-y-0.5"
              >
                Crear cuenta
              </Link>
            </div>
          )}
        </div>
      </nav>
    </header>
  );
};
