import { useEffect, useState } from 'react';
import { Search, SlidersHorizontal } from 'lucide-react';

export const SearchBar = ({ initialQuery = '', onSearch, onToggleFilters }) => {
  const [query, setQuery] = useState(initialQuery);

  useEffect(() => {
    setQuery(initialQuery);
  }, [initialQuery]);

  const handleSubmit = (event) => {
    event.preventDefault();
    onSearch?.(query);
  };

  return (
    <form
      onSubmit={handleSubmit}
      className="elegant-card flex flex-col gap-3 p-5 sm:flex-row sm:items-center"
    >
      <div className="flex flex-1 items-center gap-3 rounded-xl bg-slate-50 px-4 py-3 dark:bg-slate-900/50">
        <Search className="text-primary-600 dark:text-primary-400" size={20} />
        <input
          value={query}
          onChange={(event) => setQuery(event.target.value)}
          placeholder="Buscar por nombre, cliente, mesa..."
          className="w-full border-none bg-transparent text-base text-slate-700 outline-none placeholder:text-slate-400 dark:text-slate-100 dark:placeholder:text-slate-500"
        />
      </div>
      <div className="flex gap-2">
        <button
          type="button"
          onClick={onToggleFilters}
          className="inline-flex items-center gap-2 rounded-xl border border-slate-200 px-4 py-2.5 text-sm font-medium text-slate-600 transition-all hover:border-primary-300 hover:text-primary-700 dark:border-slate-700 dark:text-slate-300 dark:hover:border-primary-600 dark:hover:text-primary-400"
        >
          <SlidersHorizontal size={18} />
          Filtros
        </button>
        <button
          type="submit"
          className="luxury-button"
        >
          Buscar
        </button>
      </div>
    </form>
  );
};
