import { MEAL_TYPES } from '../../utils/constants';

const normalizeDateValue = (raw) => {
  if (!raw) return '';
  const isoPattern = /^\d{4}-\d{2}-\d{2}$/;
  if (isoPattern.test(raw)) return raw;
  const match = raw.match(/^(\d{2})[\/-](\d{2})[\/-](\d{4})$/);
  if (match) {
    const [, day, month, year] = match;
    return `${year}-${month}-${day}`;
  }
  const parsed = new Date(raw);
  if (!Number.isNaN(parsed.getTime())) {
    return parsed.toISOString().slice(0, 10);
  }
  return '';
};

export const FilterPanel = ({ filters, onChange, visible }) => {
  if (!visible) return null;

  const handleChange = (event) => {
    const { name, value } = event.target;
    const normalizedValue = name === 'date' ? normalizeDateValue(value) : value;
    onChange?.({ ...filters, [name]: normalizedValue });
  };

  return (
    <div className="mt-4 rounded-3xl border border-white/40 bg-white/70 p-5 shadow-lg shadow-primary-500/10 backdrop-blur-xl dark:border-white/10 dark:bg-white/5">
      <div className="grid gap-4 sm:grid-cols-4">
        <label className="text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-300">
          Tipo de comida
          <select
            name="meal_type"
            value={filters.meal_type ?? ''}
            onChange={handleChange}
            className="mt-2 w-full rounded-2xl border border-white/50 bg-white/90 px-3 py-3 text-sm text-slate-700 outline-none transition focus:border-primary-400 focus:ring-2 focus:ring-primary-200 dark:border-white/15 dark:bg-slate-900/70 dark:text-slate-100 dark:focus:border-primary-500"
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
          Disponibilidad
          <select
            name="is_available"
            value={filters.is_available ?? 'true'}
            onChange={handleChange}
            className="mt-2 w-full rounded-2xl border border-white/50 bg-white/90 px-3 py-3 text-sm text-slate-700 outline-none transition focus:border-primary-400 focus:ring-2 focus:ring-primary-200 dark:border-white/15 dark:bg-slate-900/70 dark:text-slate-100 dark:focus:border-primary-500"
          >
            <option value="">Todas</option>
            <option value="true">Solo disponibles</option>
            <option value="false">Solo reservadas</option>
          </select>
        </label>
        <label className="text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-300">
          Capacidad mínima
          <input
            type="number"
            min="1"
            max="20"
            name="capacity"
            value={filters.capacity ?? ''}
            onChange={handleChange}
            className="mt-2 w-full rounded-2xl border border-white/50 bg-white/90 px-3 py-3 text-sm text-slate-700 outline-none transition focus:border-primary-400 focus:ring-2 focus:ring-primary-200 dark:border-white/15 dark:bg-slate-900/70 dark:text-slate-100 dark:focus:border-primary-500"
            placeholder="4"
          />
        </label>
        <label className="text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-300">
          Fecha
          <input
            type="date"
            name="date"
            value={filters.date ?? ''}
            onChange={handleChange}
            className="mt-2 w-full rounded-2xl border border-white/50 bg-white/90 px-3 py-3 text-sm text-slate-700 outline-none transition focus:border-primary-400 focus:ring-2 focus:ring-primary-200 dark:border-white/15 dark:bg-slate-900/70 dark:text-slate-100 dark:focus:border-primary-500"
          />
        </label>
      </div>
    </div>
  );
};
