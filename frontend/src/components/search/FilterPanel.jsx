import { MEAL_TYPES, RESERVATION_STATUSES } from '../../utils/constants';

export const FilterPanel = ({ filters, onChange, visible }) => {
  if (!visible) return null;

  const handleChange = (event) => {
    const { name, value } = event.target;
    onChange?.({ ...filters, [name]: value });
  };

  return (
    <div className="mt-4 rounded-3xl border border-white/40 bg-white/70 p-5 shadow-lg shadow-primary-500/10 backdrop-blur-xl dark:border-white/10 dark:bg-white/5">
      <div className="grid gap-4 sm:grid-cols-3">
        <label className="text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-300">
          Tipo de comida
          <select
            name="meal_type"
            value={filters.meal_type ?? ''}
            onChange={handleChange}
            className="mt-2 w-full rounded-2xl border border-white/50 bg-white/80 px-3 py-3 text-sm text-slate-700 outline-none transition focus:border-primary-400 focus:ring-2 focus:ring-primary-200 dark:border-white/10 dark:bg-white/10 dark:text-slate-100 dark:focus:border-primary-500"
          >
            <option value="">Todos</option>
            {MEAL_TYPES.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
        </label>
        <label className="text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-300">
          Estado
          <select
            name="status"
            value={filters.status ?? ''}
            onChange={handleChange}
            className="mt-2 w-full rounded-2xl border border-white/50 bg-white/80 px-3 py-3 text-sm text-slate-700 outline-none transition focus:border-primary-400 focus:ring-2 focus:ring-primary-200 dark:border-white/10 dark:bg-white/10 dark:text-slate-100 dark:focus:border-primary-500"
          >
            <option value="">Todos</option>
            {RESERVATION_STATUSES.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
        </label>
        <label className="text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-300">
          Cantidad de comensales
          <input
            type="number"
            min="1"
            max="20"
            name="guests"
            value={filters.guests ?? ''}
            onChange={handleChange}
            className="mt-2 w-full rounded-2xl border border-white/50 bg-white/80 px-3 py-3 text-sm text-slate-700 outline-none transition focus:border-primary-400 focus:ring-2 focus:ring-primary-200 dark:border-white/10 dark:bg-white/10 dark:text-slate-100 dark:focus:border-primary-500"
            placeholder="4"
          />
        </label>
      </div>
    </div>
  );
};
