import { Link } from 'react-router-dom';
import { CalendarClock, Users, ChefHat } from 'lucide-react';
import { formatCurrency, formatDateTime, formatStatus } from '../../utils/formatters';

export const ReservationCard = ({ reservation }) => (
  <article className="flex flex-col gap-3 rounded-2xl border border-slate-100 bg-white/80 p-4 shadow-soft transition hover:-translate-y-1 hover:shadow-lg">
    <header className="flex items-center justify-between gap-2">
      <span className="text-xs font-semibold uppercase tracking-wide text-slate-400">
        Mesa #{reservation.table_number ?? reservation.tableNumber}
      </span>
      <span className="rounded-full bg-slate-100 px-3 py-1 text-xs font-semibold capitalize text-slate-600">
        {formatStatus(reservation.status)}
      </span>
    </header>
    <div className="space-y-2 text-sm text-slate-600">
      <p className="flex items-center gap-2">
        <CalendarClock size={16} className="text-primary-500" />
        {formatDateTime(reservation.date_time || reservation.dateTime)}
      </p>
      <p className="flex items-center gap-2">
        <Users size={16} className="text-primary-500" />
        {reservation.guests} comensales
      </p>
      <p className="flex items-center gap-2 capitalize">
        <ChefHat size={16} className="text-primary-500" />
        {reservation.meal_type || reservation.mealType}
      </p>
      <p className="text-base font-semibold text-slate-900">{formatCurrency(reservation.total_price || reservation.totalPrice)}</p>
    </div>
    <div className="mt-auto">
      <Link
        to={`/reservations/${reservation.id}`}
        className="inline-flex w-full items-center justify-center rounded-xl border border-primary-200 px-4 py-2 text-sm font-semibold text-primary-700 hover:bg-primary-50"
      >
        Ver detalles
      </Link>
    </div>
  </article>
);
