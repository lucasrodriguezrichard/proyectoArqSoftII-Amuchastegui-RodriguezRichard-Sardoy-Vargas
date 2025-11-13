import { Link, NavLink, useNavigate } from 'react-router-dom';
import { UtensilsCrossed, LogOut, UserRound } from 'lucide-react';
import { APP_NAME } from '../../utils/constants';
import { useAuth } from '../../hooks/useAuth';

const navItems = [
  { to: '/', label: 'Buscar', private: false },
  { to: '/create-reservation', label: 'Crear reserva', private: true },
  { to: '/my-reservations', label: 'Mis reservas', private: true },
  { to: '/admin', label: 'Admin', private: true, adminOnly: true },
];

const linkClasses = ({ isActive }) =>
  `px-3 py-2 rounded-full text-sm font-medium ${
    isActive ? 'bg-primary-100 text-primary-700' : 'text-slate-500 hover:text-slate-900'
  }`;

export const Navbar = () => {
  const navigate = useNavigate();
  const { isAuthenticated, isAdmin, user, logout } = useAuth();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <header className="sticky top-0 z-30 border-b border-slate-100 bg-white/90 backdrop-blur">
      <nav className="mx-auto flex max-w-6xl items-center justify-between px-4 py-3">
        <Link to="/" className="flex items-center gap-2 font-display text-lg font-semibold text-primary-700">
          <UtensilsCrossed size={22} />
          {APP_NAME}
        </Link>
        <div className="flex items-center gap-4">
          <div className="hidden items-center gap-1 rounded-full bg-slate-50 px-2 py-1 sm:flex">
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
              <span className="hidden text-sm font-medium text-slate-700 sm:flex sm:flex-col sm:items-end">
                <span className="flex items-center gap-1 text-xs uppercase tracking-wide text-slate-400">
                  <UserRound size={14} />
                  {user?.role ?? 'user'}
                </span>
                {user?.first_name || user?.username}
              </span>
              <button
                type="button"
                onClick={handleLogout}
                className="inline-flex items-center gap-1 rounded-full border border-slate-200 px-3 py-1 text-sm font-medium text-slate-600 hover:border-slate-300 hover:text-slate-900"
              >
                <LogOut size={16} />
                Salir
              </button>
            </div>
          ) : (
            <div className="flex items-center gap-2">
              <Link
                to="/login"
                className="rounded-full border border-primary-100 px-4 py-2 text-sm font-semibold text-primary-600 transition hover:bg-primary-50"
              >
                Ingresar
              </Link>
              <Link
                to="/register"
                className="rounded-full bg-primary-600 px-4 py-2 text-sm font-semibold text-white shadow-md shadow-primary-200 transition hover:bg-primary-700"
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
