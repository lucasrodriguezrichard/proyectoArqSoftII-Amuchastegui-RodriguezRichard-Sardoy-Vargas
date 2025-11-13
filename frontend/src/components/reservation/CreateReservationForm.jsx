import { useState } from 'react';
import { CalendarDays, Users, Utensils, Hash } from 'lucide-react';
import { MEAL_TYPES } from '../../utils/constants';

export const CreateReservationForm = ({ onSubmit, loading, userId }) => {
  const [formData, setFormData] = useState({
    owner_id: userId || '',
    table_number: '',
    guests: '',
    meal_type: '',
    date_time: '',
    special_requests: '',
  });

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = (e) => {
    e.preventDefault();

    // Convert to correct types
    const payload = {
      owner_id: formData.owner_id,
      table_number: parseInt(formData.table_number, 10),
      guests: parseInt(formData.guests, 10),
      meal_type: formData.meal_type,
      date_time: new Date(formData.date_time).toISOString(),
      special_requests: formData.special_requests || undefined,
    };

    onSubmit(payload);
  };

  const today = new Date().toISOString().slice(0, 16);

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="grid gap-4 sm:grid-cols-2">
        <div>
          <label htmlFor="table_number" className="mb-2 flex items-center gap-2 text-sm font-medium text-slate-700">
            <Hash size={16} />
            Número de mesa
          </label>
          <input
            type="number"
            id="table_number"
            name="table_number"
            min="1"
            required
            value={formData.table_number}
            onChange={handleChange}
            className="w-full rounded-xl border border-slate-200 px-4 py-2.5 text-sm focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-100"
            placeholder="Ej: 5"
          />
        </div>

        <div>
          <label htmlFor="guests" className="mb-2 flex items-center gap-2 text-sm font-medium text-slate-700">
            <Users size={16} />
            Cantidad de comensales
          </label>
          <input
            type="number"
            id="guests"
            name="guests"
            min="1"
            max="20"
            required
            value={formData.guests}
            onChange={handleChange}
            className="w-full rounded-xl border border-slate-200 px-4 py-2.5 text-sm focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-100"
            placeholder="Ej: 4"
          />
        </div>

        <div>
          <label htmlFor="meal_type" className="mb-2 flex items-center gap-2 text-sm font-medium text-slate-700">
            <Utensils size={16} />
            Tipo de comida
          </label>
          <select
            id="meal_type"
            name="meal_type"
            required
            value={formData.meal_type}
            onChange={handleChange}
            className="w-full rounded-xl border border-slate-200 px-4 py-2.5 text-sm focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-100"
          >
            <option value="">Seleccionar...</option>
            {MEAL_TYPES.map((type) => (
              <option key={type.value} value={type.value}>
                {type.label}
              </option>
            ))}
          </select>
        </div>

        <div>
          <label htmlFor="date_time" className="mb-2 flex items-center gap-2 text-sm font-medium text-slate-700">
            <CalendarDays size={16} />
            Fecha y hora
          </label>
          <input
            type="datetime-local"
            id="date_time"
            name="date_time"
            required
            min={today}
            value={formData.date_time}
            onChange={handleChange}
            className="w-full rounded-xl border border-slate-200 px-4 py-2.5 text-sm focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-100"
          />
        </div>
      </div>

      <div>
        <label htmlFor="special_requests" className="mb-2 block text-sm font-medium text-slate-700">
          Pedidos especiales (opcional)
        </label>
        <textarea
          id="special_requests"
          name="special_requests"
          rows="3"
          value={formData.special_requests}
          onChange={handleChange}
          className="w-full rounded-xl border border-slate-200 px-4 py-2.5 text-sm focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-100"
          placeholder="Ej: Mesa cerca de la ventana, cumpleaños, alergias..."
        />
      </div>

      <button
        type="submit"
        disabled={loading}
        className="w-full rounded-xl bg-primary-600 px-6 py-3 font-semibold text-white shadow-lg shadow-primary-200 transition hover:bg-primary-700 disabled:cursor-not-allowed disabled:opacity-50"
      >
        {loading ? 'Creando reserva...' : 'Crear reserva'}
      </button>
    </form>
  );
};
