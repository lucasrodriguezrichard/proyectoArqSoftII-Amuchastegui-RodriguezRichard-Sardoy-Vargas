import { useState } from 'react';
import { CalendarDays, Users, Utensils, Hash } from 'lucide-react';
import { useQuery } from '@tanstack/react-query';
import { MEAL_TYPES } from '../../utils/constants';
import { getAvailableTables } from '../../api/reservations';

export const CreateReservationForm = ({ onSubmit, loading, userId }) => {
  const [formData, setFormData] = useState({
    owner_id: userId || '',
    table_number: '',
    guests: '',
    meal_type: '',
    date: '',
    special_requests: '',
  });

  const [selectedTable, setSelectedTable] = useState(null);

  // Fetch available tables when date and meal_type are selected
  const { data: availableTables = [], isLoading: loadingTables } = useQuery({
    queryKey: ['available-tables', formData.date, formData.meal_type],
    queryFn: () => getAvailableTables({ date: formData.date, mealType: formData.meal_type }),
    enabled: Boolean(formData.date && formData.meal_type),
    staleTime: 1000 * 30, // 30 seconds
  });

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));

    // Reset table selection when date or meal_type changes
    if (name === 'date' || name === 'meal_type') {
      setSelectedTable(null);
      setFormData((prev) => ({ ...prev, table_number: '', guests: '' }));
    }
  };

  const handleTableSelect = (table) => {
    setSelectedTable(table);
    setFormData((prev) => ({
      ...prev,
      table_number: table.table_number,
      guests: table.capacity, // Auto-fill with table capacity
    }));
  };

  const handleSubmit = (e) => {
    e.preventDefault();

    if (!selectedTable) {
      alert('Por favor seleccioná una mesa disponible');
      return;
    }

    // Create datetime at noon (12:00) for the selected date
    const dateTime = new Date(formData.date + 'T12:00:00');

    const payload = {
      owner_id: formData.owner_id,
      table_number: formData.table_number,
      guests: formData.guests,
      meal_type: formData.meal_type,
      date_time: dateTime.toISOString(),
      special_requests: formData.special_requests || undefined,
    };

    onSubmit(payload);
  };

  const today = new Date().toISOString().slice(0, 10);

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      {/* Step 1: Select Date and Meal Type */}
      <div className="grid gap-6 sm:grid-cols-2">
        <div>
          <label htmlFor="date" className="mb-2 flex items-center gap-2 text-sm font-medium text-slate-700 dark:text-slate-300">
            <CalendarDays size={16} className="text-primary-600 dark:text-primary-400" />
            Fecha
          </label>
          <input
            type="date"
            id="date"
            name="date"
            required
            min={today}
            value={formData.date}
            onChange={handleChange}
            className="luxury-input"
          />
        </div>

        <div>
          <label htmlFor="meal_type" className="mb-2 flex items-center gap-2 text-sm font-medium text-slate-700 dark:text-slate-300">
            <Utensils size={16} className="text-primary-600 dark:text-primary-400" />
            Tipo de comida
          </label>
          <select
            id="meal_type"
            name="meal_type"
            required
            value={formData.meal_type}
            onChange={handleChange}
            className="luxury-input"
          >
            <option value="">Seleccionar...</option>
            {MEAL_TYPES.map((type) => (
              <option key={type.value} value={type.value}>
                {type.label}
              </option>
            ))}
          </select>
        </div>
      </div>

      {/* Step 2: Show Available Tables */}
      {formData.date && formData.meal_type && (
        <div>
          <label className="mb-3 block text-sm font-medium text-slate-700 dark:text-slate-300">
            Mesas disponibles
          </label>

          {loadingTables ? (
            <div className="rounded-xl border border-slate-200 bg-slate-50 p-6 text-center text-sm text-slate-500 dark:border-slate-700 dark:bg-slate-800 dark:text-slate-400">
              Cargando mesas disponibles...
            </div>
          ) : availableTables.length === 0 ? (
            <div className="rounded-xl border border-red-200 bg-red-50 p-6 text-center text-sm text-red-600 dark:border-red-900 dark:bg-red-950 dark:text-red-400">
              No hay mesas disponibles para esta fecha y tipo de comida.
            </div>
          ) : (
            <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
              {availableTables.map((table) => (
                <button
                  key={table.table_number}
                  type="button"
                  onClick={() => handleTableSelect(table)}
                  className={`rounded-xl border-2 p-4 text-left transition ${
                    selectedTable?.table_number === table.table_number
                      ? 'border-primary-500 bg-primary-50 dark:border-primary-400 dark:bg-primary-950'
                      : 'border-slate-200 bg-white hover:border-primary-300 dark:border-slate-700 dark:bg-slate-800 dark:hover:border-primary-600'
                  }`}
                >
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                      <Hash size={18} className="text-primary-600 dark:text-primary-400" />
                      <span className="font-semibold text-slate-900 dark:text-slate-100">
                        Mesa {table.table_number}
                      </span>
                    </div>
                    {selectedTable?.table_number === table.table_number && (
                      <span className="text-primary-600 dark:text-primary-400">✓</span>
                    )}
                  </div>
                  <div className="mt-2 flex items-center gap-1 text-sm text-slate-600 dark:text-slate-400">
                    <Users size={14} />
                    {table.capacity} personas
                  </div>
                </button>
              ))}
            </div>
          )}
        </div>
      )}

      {/* Step 3: Guests (auto-filled from table capacity, but editable) */}
      {selectedTable && (
        <div>
          <label htmlFor="guests" className="mb-2 flex items-center gap-2 text-sm font-medium text-slate-700 dark:text-slate-300">
            <Users size={16} className="text-primary-600 dark:text-primary-400" />
            Cantidad de comensales
          </label>
          <input
            type="number"
            id="guests"
            name="guests"
            min="1"
            max={selectedTable.capacity}
            required
            value={formData.guests}
            onChange={handleChange}
            className="luxury-input"
            placeholder={`Máximo ${selectedTable.capacity} personas`}
          />
          <p className="mt-1 text-xs text-slate-500 dark:text-slate-400">
            Esta mesa tiene capacidad para {selectedTable.capacity} personas
          </p>
        </div>
      )}

      {/* Step 4: Special Requests */}
      {selectedTable && (
        <div>
          <label htmlFor="special_requests" className="mb-2 block text-sm font-medium text-slate-700 dark:text-slate-300">
            Pedidos especiales (opcional)
          </label>
          <textarea
            id="special_requests"
            name="special_requests"
            rows="3"
            value={formData.special_requests}
            onChange={handleChange}
            className="luxury-input resize-none"
            placeholder="Ej: Mesa cerca de la ventana, cumpleaños, alergias..."
          />
        </div>
      )}

      <button
        type="submit"
        disabled={loading || !selectedTable}
        className="luxury-button w-full disabled:cursor-not-allowed disabled:opacity-50"
      >
        {loading ? 'Creando reserva...' : 'Crear reserva'}
      </button>
    </form>
  );
};
