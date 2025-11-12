import { MEAL_TYPES, RESERVATION_STATUSES } from '../../utils/constants';

export const FilterPanel = ({ filters, onChange, visible }) => {
  if (!visible) return null;

  const handleChange = (event) => {
    const { name, value } = event.target;
    onChange?.({ ...filters, [name]: value });
  };

  return (
    <div className="glass-panel mt-4 p-4">
      <div className="grid gap-4 sm:grid-cols-3">
        <label className="text-sm font-medium text-slate-600">
          Tipo de comida
          <select
            name="meal_type"
            value={filters.meal_type ?? ''}
            onChange={handleChange}
            className="mt-1 w-full rounded-xl border border-slate-200 px-3 py-2 text-slate-700 focus:border-primary-300 focus:outline-none focus:ring-2 focus:ring-primary-100"
          >
            <option value="">Todos</option>
            {MEAL_TYPES.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
        </label>
        <label className="text-sm font-medium text-slate-600">
          Estado
          <select
            name="status"
            value={filters.status ?? ''}
            onChange={handleChange}
            className="mt-1 w-full rounded-xl border border-slate-200 px-3 py-2 text-slate-700 focus:border-primary-300 focus:outline-none focus:ring-2 focus:ring-primary-100"
          >
            <option value="">Todos</option>
            {RESERVATION_STATUSES.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
        </label>
        <label className="text-sm font-medium text-slate-600">
          Cantidad de comensales
          <input
            type="number"
            min="1"
            max="20"
            name="guests"
            value={filters.guests ?? ''}
            onChange={handleChange}
            className="mt-1 w-full rounded-xl border border-slate-200 px-3 py-2 text-slate-700 focus:border-primary-300 focus:outline-none focus:ring-2 focus:ring-primary-100"
            placeholder="4"
          />
        </label>
      </div>
    </div>
  );
};
