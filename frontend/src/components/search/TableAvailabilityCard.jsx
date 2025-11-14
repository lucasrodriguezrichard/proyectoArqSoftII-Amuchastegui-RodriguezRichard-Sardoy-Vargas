import { useState } from 'react';
import { Link } from 'react-router-dom';
import { CalendarClock, Users, ChefHat, CheckCircle, XCircle } from 'lucide-react';

export const TableAvailabilityCard = ({ table, onReserve }) => {
  const [reserving, setReserving] = useState(false);
  const reservationId = table.reservation_id || table.reservationId || table.id;

  const handleReserveClick = async () => {
    if (!onReserve || reserving) return;
    setReserving(true);
    try {
      await onReserve(table);
    } finally {
      setReserving(false);
    }
  };

  return (
    <article className="relative flex flex-col gap-4 rounded-3xl border border-white/50 bg-white/80 p-5 shadow-xl shadow-primary-500/10 backdrop-blur-xl transition hover:-translate-y-1 hover:border-primary-200 dark:border-white/10 dark:bg-slate-900/50">
      <header className="flex items-center justify-between gap-2 border-b border-white/70 pb-3 dark:border-white/10">
        <span className="text-[0.65rem] font-semibold uppercase tracking-[0.3em] text-slate-500 dark:text-slate-300">
          Mesa #{table.table_number}
        </span>
        <span
          className={`flex items-center gap-1.5 rounded-full px-3 py-1 text-xs font-semibold ${
            table.is_available
              ? 'bg-gradient-to-r from-green-500/15 to-emerald-500/20 text-green-600 dark:text-green-200'
              : 'bg-gradient-to-r from-red-500/15 to-rose-500/20 text-red-600 dark:text-red-200'
          }`}
        >
          {table.is_available ? (
            <>
              <CheckCircle size={14} />
              Disponible
            </>
          ) : (
            <>
              <XCircle size={14} />
              Reservada
            </>
          )}
        </span>
      </header>
      <div className="space-y-3 text-sm text-slate-600 dark:text-slate-300">
        <p className="flex items-center gap-3 font-semibold">
          <CalendarClock size={18} className="text-primary-500 dark:text-primary-300" />
          {table.date}
        </p>
        <p className="flex items-center gap-3 font-semibold">
          <Users size={18} className="text-primary-500 dark:text-primary-300" />
          Capacidad: {table.capacity} personas
        </p>
        <p className="flex items-center gap-3 font-semibold capitalize">
          <ChefHat size={18} className="text-primary-500 dark:text-primary-300" />
          {table.meal_type}
        </p>
      </div>
      <div className="mt-auto pt-2">
        {table.is_available ? (
          <button
            type="button"
            onClick={handleReserveClick}
            disabled={reserving}
            className="inline-flex w-full items-center justify-center rounded-2xl border border-primary-200 bg-primary-600 px-4 py-2.5 text-sm font-semibold text-white transition-all hover:bg-primary-700 disabled:cursor-not-allowed disabled:opacity-70 dark:border-primary-800 dark:bg-primary-700 dark:hover:bg-primary-600"
          >
            {reserving ? 'Reservando...' : 'Reservar mesa'}
          </button>
        ) : reservationId ? (
          <Link
            to={`/reservations/${reservationId}`}
            className="inline-flex w-full items-center justify-center rounded-2xl border border-slate-200 px-4 py-2.5 text-sm font-semibold text-slate-700 transition-all hover:bg-slate-100 dark:border-white/30 dark:text-white dark:hover:bg-white/10"
          >
            Ver detalle
          </Link>
        ) : (
          <p className="text-center text-xs text-slate-400 dark:text-slate-500">Sin información de reserva</p>
        )}
      </div>
    </article>
  );
};
