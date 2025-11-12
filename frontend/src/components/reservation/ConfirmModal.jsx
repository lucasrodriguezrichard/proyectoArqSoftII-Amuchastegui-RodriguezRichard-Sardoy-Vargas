import { useEffect, useState } from 'react';

export const ConfirmModal = ({ open, onClose, onConfirm, loading }) => {
  const [notes, setNotes] = useState('');

  useEffect(() => {
    if (!open) {
      setNotes('');
    }
  }, [open]);

  if (!open) return null;

  const handleSubmit = (event) => {
    event.preventDefault();
    onConfirm?.({ confirmation_notes: notes });
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/40 px-4">
      <form onSubmit={handleSubmit} className="w-full max-w-lg rounded-3xl bg-white p-6 shadow-2xl">
        <h3 className="text-lg font-semibold text-slate-900">Confirmar reserva</h3>
        <p className="mt-1 text-sm text-slate-500">El cálculo concurrente se ejecutará en el backend.</p>
        <label className="mt-4 block text-sm font-medium text-slate-600">
          Notas
          <textarea
            value={notes}
            onChange={(event) => setNotes(event.target.value)}
            rows={4}
            className="mt-1 w-full rounded-2xl border border-slate-200 px-3 py-2 text-slate-700 focus:border-primary-300 focus:outline-none focus:ring-2 focus:ring-primary-100"
          />
        </label>
        <div className="mt-6 flex justify-end gap-2">
          <button
            type="button"
            onClick={onClose}
            className="rounded-full border border-slate-200 px-4 py-2 text-sm font-semibold text-slate-600 hover:bg-slate-50"
          >
            Cancelar
          </button>
          <button
            type="submit"
            disabled={loading}
            className="rounded-full bg-primary-600 px-5 py-2 text-sm font-semibold text-white shadow-soft hover:bg-primary-700 disabled:cursor-not-allowed disabled:bg-slate-300"
          >
            {loading ? 'Confirmando...' : 'Confirmar'}
          </button>
        </div>
      </form>
    </div>
  );
};
