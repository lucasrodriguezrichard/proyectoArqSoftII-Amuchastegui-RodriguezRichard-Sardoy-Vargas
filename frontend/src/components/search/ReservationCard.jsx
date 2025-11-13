import { Link } from 'react-router-dom';
import { CalendarClock, Users, ChefHat } from 'lucide-react';
import { formatCurrency, formatDateTime, formatStatus } from '../../utils/formatters';

export const ReservationCard = ({ reservation }) => (
  <article className="elegant-card flex flex-col gap-4 p-5 transition-all hover:-translate-y-1">
    <header className="flex items-center justify-between gap-2 border-b border-slate-100 pb-3 dark:border-slate-700">
      <span className="text-xs font-light uppercase tracking-widest text-slate-500 dark:text-slate-400">
        Mesa #{reservation.table_number ?? reservation.tableNumber}
      </span>
      <span className="rounded-full bg-slate-100 px-3 py-1 text-xs font-medium capitalize text-slate-700 dark:bg-slate-700 dark:text-slate-300">
        {formatStatus(reservation.status)}
      </span>
    </header>
    <div className="space-y-3 text-sm text-slate-600 dark:text-slate-400">
      <p className="flex items-center gap-3 font-light">
        <CalendarClock size={18} className="text-primary-600 dark:text-primary-400" />
        {formatDateTime(reservation.date_time || reservation.dateTime)}
      </p>
      <p className="flex items-center gap-3 font-light">
        <Users size={18} className="text-primary-600 dark:text-primary-400" />
        {reservation.guests} comensales
      </p>
      <p className="flex items-center gap-3 font-light capitalize">
        <ChefHat size={18} className="text-primary-600 dark:text-primary-400" />
        {reservation.meal_type || reservation.mealType}
      </p>
      <p className="text-lg font-medium text-slate-900 dark:text-slate-100">{formatCurrency(reservation.total_price || reservation.totalPrice)}</p>
    </div>
    <div className="mt-auto pt-2">
      <Link
        to={`/reservations/${reservation.id}`}
        className="inline-flex w-full items-center justify-center rounded-lg border border-primary-200 px-4 py-2.5 text-sm font-medium text-primary-700 transition-all hover:bg-primary-50 hover:border-primary-300 dark:border-primary-800 dark:text-primary-300 dark:hover:bg-primary-900/30 dark:hover:border-primary-700"
      >
        Ver detalles
      </Link>
    </div>
  </article>
);
