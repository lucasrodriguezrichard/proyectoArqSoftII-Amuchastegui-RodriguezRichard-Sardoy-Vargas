import { Link } from 'react-router-dom';
import { CalendarClock, Users, ChefHat } from 'lucide-react';
import { formatCurrency, formatDateTime, formatStatus } from '../../utils/formatters';

export const ReservationCard = ({ reservation }) => (
  <article className="relative flex flex-col gap-4 rounded-3xl border border-white/50 bg-white/80 p-5 shadow-xl shadow-primary-500/10 backdrop-blur-xl transition hover:-translate-y-1 hover:border-primary-200 dark:border-white/10 dark:bg-slate-900/50">
    <header className="flex items-center justify-between gap-2 border-b border-white/70 pb-3 dark:border-white/10">
      <span className="text-[0.65rem] font-semibold uppercase tracking-[0.3em] text-slate-500 dark:text-slate-300">
        Mesa #{reservation.table_number ?? reservation.tableNumber}
      </span>
      <span className="rounded-full bg-gradient-to-r from-primary-500/15 to-sky-500/20 px-3 py-1 text-xs font-semibold capitalize text-primary-600 dark:text-primary-200">
        {formatStatus(reservation.status)}
      </span>
    </header>
    <div className="space-y-3 text-sm text-slate-600 dark:text-slate-300">
      <p className="flex items-center gap-3 font-semibold">
        <CalendarClock size={18} className="text-primary-500 dark:text-primary-300" />
        {formatDateTime(reservation.date_time || reservation.dateTime)}
      </p>
      <p className="flex items-center gap-3 font-semibold">
        <Users size={18} className="text-primary-500 dark:text-primary-300" />
        {reservation.guests} comensales
      </p>
      <p className="flex items-center gap-3 font-semibold capitalize">
        <ChefHat size={18} className="text-primary-500 dark:text-primary-300" />
        {reservation.meal_type || reservation.mealType}
      </p>
      <p className="text-2xl font-semibold text-slate-900 dark:text-white">{formatCurrency(reservation.total_price || reservation.totalPrice)}</p>
    </div>
    <div className="mt-auto pt-2">
      <Link
        to={`/reservations/${reservation.id}`}
        className="inline-flex w-full items-center justify-center rounded-2xl border border-primary-200 px-4 py-2.5 text-sm font-semibold text-primary-700 transition-all hover:bg-primary-50 dark:border-primary-800 dark:text-primary-200 dark:hover:bg-primary-900/30"
      >
        Ver detalles
      </Link>
    </div>
  </article>
);
