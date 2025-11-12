import { CalendarDays, Users, BadgeCheck, Utensils } from 'lucide-react';
import { formatCurrency, formatDateTime, formatStatus } from '../../utils/formatters';

export const ReservationDetails = ({ reservation }) => {
  if (!reservation) return null;

  return (
    <section className="glass-panel p-6">
      <div className="flex flex-col gap-6">
        <div>
          <p className="text-xs uppercase tracking-wide text-slate-400">Reserva</p>
          <h2 className="font-display text-2xl font-semibold text-slate-900">#{reservation.id}</h2>
        </div>

        <div className="grid gap-4 md:grid-cols-2">
          <div className="rounded-2xl bg-slate-50 p-4">
            <p className="text-xs uppercase tracking-wide text-slate-400">Fecha & hora</p>
            <p className="mt-1 flex items-center gap-2 text-lg font-semibold text-slate-900">
              <CalendarDays className="text-primary-500" size={18} />
              {formatDateTime(reservation.date_time || reservation.dateTime)}
            </p>
          </div>
          <div className="rounded-2xl bg-slate-50 p-4">
            <p className="text-xs uppercase tracking-wide text-slate-400">Estado</p>
            <p className="mt-1 flex items-center gap-2 text-lg font-semibold capitalize text-slate-900">
              <BadgeCheck className="text-primary-500" size={18} />
              {formatStatus(reservation.status)}
            </p>
          </div>
        </div>

        <div className="grid gap-4 md:grid-cols-3">
          <DetailTile icon={<Utensils size={18} />} label="Mesa" value={`#${reservation.table_number ?? reservation.tableNumber}`} />
          <DetailTile icon={<Users size={18} />} label="Comensales" value={`${reservation.guests}`} />
          <DetailTile icon={<Utensils size={18} />} label="Tipo" value={reservation.meal_type || reservation.mealType} />
        </div>

        <div className="grid gap-4 md:grid-cols-2">
          <DetailTile label="Solicitante" value={reservation.owner_id || reservation.ownerId} />
          <DetailTile label="Monto estimado" value={formatCurrency(reservation.total_price || reservation.totalPrice)} />
        </div>

        {reservation.special_requests && (
          <div>
            <p className="text-xs uppercase tracking-wide text-slate-400">Notas</p>
            <p className="mt-2 rounded-2xl border border-slate-100 bg-slate-50 p-4 text-slate-600">
              {reservation.special_requests}
            </p>
          </div>
        )}
      </div>
    </section>
  );
};

const DetailTile = ({ icon, label, value }) => (
  <div className="rounded-2xl border border-slate-100 bg-white p-4 shadow-sm">
    <p className="text-xs uppercase tracking-wide text-slate-400">{label}</p>
    <p className="mt-2 flex items-center gap-2 text-lg font-semibold text-slate-900">
      {icon}
      <span className="capitalize">{value || '—'}</span>
    </p>
  </div>
);
