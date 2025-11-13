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
  `px-3 py-2 rounded-full text-sm font-medium ${
    isActive ? 'bg-primary-100 text-primary-700 dark:bg-primary-900/30 dark:text-primary-300' : 'text-slate-500 hover:text-slate-900 dark:text-slate-400 dark:hover:text-slate-100'
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
    <header className="sticky top-0 z-30 border-b border-slate-100 bg-white/90 backdrop-blur dark:border-slate-800 dark:bg-slate-900/90">
      <nav className="mx-auto flex max-w-6xl items-center justify-between px-4 py-3">
        <Link to="/" className="flex items-center gap-2 font-display text-lg font-semibold text-primary-700 dark:text-primary-400">
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
            className="rounded-full p-2 text-slate-600 hover:bg-slate-100 dark:text-slate-400 dark:hover:bg-slate-800"
            aria-label="Toggle theme"
          >
            {theme === 'light' ? <Moon size={18} /> : <Sun size={18} />}
          </button>
          <div className="hidden items-center gap-1 rounded-full bg-slate-50 px-2 py-1 dark:bg-slate-800 sm:flex">
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
              <span className="hidden text-sm font-medium text-slate-700 dark:text-slate-300 sm:flex sm:flex-col sm:items-end">
                <span className="flex items-center gap-1 text-xs uppercase tracking-wide text-slate-400 dark:text-slate-500">
                  <UserRound size={14} />
                  {user?.role ?? 'user'}
                </span>
                {user?.first_name || user?.username}
              </span>
              <button
                type="button"
                onClick={handleLogout}
                className="inline-flex items-center gap-1 rounded-full border border-slate-200 px-3 py-1 text-sm font-medium text-slate-600 hover:border-slate-300 hover:text-slate-900 dark:border-slate-700 dark:text-slate-400 dark:hover:border-slate-600 dark:hover:text-slate-100"
              >
                <LogOut size={16} />
                Salir
              </button>
            </div>
          ) : (
            <div className="flex items-center gap-2">
              <Link
                to="/login"
                className="rounded-full border border-primary-100 px-4 py-2 text-sm font-semibold text-primary-600 transition hover:bg-primary-50 dark:border-primary-900 dark:text-primary-400 dark:hover:bg-primary-950"
              >
                Ingresar
              </Link>
              <Link
                to="/register"
                className="rounded-full bg-primary-600 px-4 py-2 text-sm font-semibold text-white shadow-md shadow-primary-200 transition hover:bg-primary-700 dark:bg-primary-700 dark:shadow-primary-900/50 dark:hover:bg-primary-600"
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
