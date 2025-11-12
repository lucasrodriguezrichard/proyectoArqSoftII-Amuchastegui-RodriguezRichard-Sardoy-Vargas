import { useEffect, useState } from 'react';
import { MEAL_TYPES, RESERVATION_STATUSES } from '../../utils/constants';

const fieldDefinitions = [
  { name: 'table_number', label: 'Mesa', type: 'number' },
  { name: 'guests', label: 'Comensales', type: 'number' },
  { name: 'meal_type', label: 'Tipo de comida', type: 'select', options: MEAL_TYPES },
  { name: 'status', label: 'Estado', type: 'select', options: RESERVATION_STATUSES },
];

export const EditModal = ({ open, reservation, onClose, onSave, loading }) => {
  const [formValues, setFormValues] = useState({});

  useEffect(() => {
    if (reservation) {
      setFormValues({
        table_number: reservation.table_number,
        guests: reservation.guests,
        meal_type: reservation.meal_type,
        status: reservation.status,
      });
    }
  }, [reservation]);

  if (!open || !reservation) return null;

  const handleChange = (event) => {
    const { name, value } = event.target;
    setFormValues((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = (event) => {
    event.preventDefault();
    onSave?.({
      reservationId: reservation.id,
      payload: {
        table_number: Number(formValues.table_number),
        guests: Number(formValues.guests),
        meal_type: formValues.meal_type,
        status: formValues.status,
      },
    });
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/40 px-4">
      <form onSubmit={handleSubmit} className="w-full max-w-xl rounded-3xl bg-white p-6 shadow-2xl">
        <h3 className="text-lg font-semibold text-slate-900">Editar reserva #{reservation.id}</h3>
        <div className="mt-4 grid gap-4 sm:grid-cols-2">
          {fieldDefinitions.map((field) => (
            <label key={field.name} className="text-sm font-medium text-slate-600">
              {field.label}
              {field.type === 'select' ? (
                <select
                  name={field.name}
                  value={formValues[field.name] ?? ''}
                  onChange={handleChange}
                  className="mt-1 w-full rounded-xl border border-slate-200 px-3 py-2 text-slate-700 focus:border-primary-300 focus:outline-none focus:ring-2 focus:ring-primary-100"
                >
                  {field.options.map((option) => (
                    <option key={option.value} value={option.value}>
                      {option.label}
                    </option>
                  ))}
                </select>
              ) : (
                <input
                  type={field.type}
                  name={field.name}
                  min="1"
                  value={formValues[field.name] ?? ''}
                  onChange={handleChange}
                  className="mt-1 w-full rounded-xl border border-slate-200 px-3 py-2 text-slate-700 focus:border-primary-300 focus:outline-none focus:ring-2 focus:ring-primary-100"
                />
              )}
            </label>
          ))}
        </div>
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
            {loading ? 'Guardando...' : 'Guardar'}
          </button>
        </div>
      </form>
    </div>
  );
};
