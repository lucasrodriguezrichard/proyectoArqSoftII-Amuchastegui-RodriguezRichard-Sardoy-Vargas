import { Pencil, Trash2 } from 'lucide-react';
import { formatCurrency, formatDateTime, formatStatus } from '../../utils/formatters';

export const ReservationTable = ({ reservations = [], onEdit, onDelete }) => {
  if (!reservations.length) {
    return <p className="rounded-2xl border border-slate-100 bg-slate-50 px-4 py-6 text-center text-sm text-slate-500">No hay reservas registradas.</p>;
  }

  return (
    <div className="overflow-x-auto rounded-2xl border border-slate-100 bg-white shadow-soft">
      <table className="min-w-full text-left text-sm text-slate-600">
        <thead className="bg-slate-50 text-xs uppercase tracking-wide text-slate-500">
          <tr>
            <th className="px-4 py-3">ID</th>
            <th className="px-4 py-3">Cliente</th>
            <th className="px-4 py-3">Fecha</th>
            <th className="px-4 py-3">Comensales</th>
            <th className="px-4 py-3">Estado</th>
            <th className="px-4 py-3">Total</th>
            <th className="px-4 py-3 text-right">Acciones</th>
          </tr>
        </thead>
        <tbody>
          {reservations.map((reservation) => (
            <tr key={reservation.id} className="border-t border-slate-100">
              <td className="px-4 py-3 font-mono text-xs text-slate-400">{reservation.id}</td>
              <td className="px-4 py-3">{reservation.owner_id}</td>
              <td className="px-4 py-3">{formatDateTime(reservation.date_time)}</td>
              <td className="px-4 py-3">{reservation.guests}</td>
              <td className="px-4 py-3 font-semibold capitalize">{formatStatus(reservation.status)}</td>
              <td className="px-4 py-3 font-semibold text-slate-900">{formatCurrency(reservation.total_price)}</td>
              <td className="px-4 py-3 text-right">
                <div className="flex justify-end gap-2">
                  <button
                    type="button"
                    onClick={() => onEdit?.(reservation)}
                    className="rounded-full border border-slate-200 p-2 text-slate-500 hover:text-primary-600"
                  >
                    <Pencil size={16} />
                  </button>
                  <button
                    type="button"
                    onClick={() => onDelete?.(reservation)}
                    className="rounded-full border border-slate-200 p-2 text-slate-500 hover:text-rose-600"
                  >
                    <Trash2 size={16} />
                  </button>
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};
