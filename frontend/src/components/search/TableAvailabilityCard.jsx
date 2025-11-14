import { Link } from 'react-router-dom';
import { CalendarClock, Users, ChefHat, CheckCircle, XCircle } from 'lucide-react';

export const TableAvailabilityCard = ({ table }) => (
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
    {table.is_available && (
      <div className="mt-auto pt-2">
        <Link
          to={`/create-reservation?table=${table.table_number}&date=${table.date}&meal_type=${table.meal_type}`}
          className="inline-flex w-full items-center justify-center rounded-2xl border border-primary-200 bg-primary-600 px-4 py-2.5 text-sm font-semibold text-white transition-all hover:bg-primary-700 dark:border-primary-800 dark:bg-primary-700 dark:hover:bg-primary-600"
        >
          Reservar mesa
        </Link>
      </div>
    )}
  </article>
);
