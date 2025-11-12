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
      className="flex flex-col gap-3 rounded-3xl border border-slate-200 bg-white/80 p-4 shadow-soft sm:flex-row sm:items-center"
    >
      <div className="flex flex-1 items-center gap-2 rounded-2xl bg-slate-50 px-4 py-2">
        <Search className="text-primary-500" size={20} />
        <input
          value={query}
          onChange={(event) => setQuery(event.target.value)}
          placeholder="Buscar por nombre, cliente, mesa..."
          className="w-full border-none bg-transparent text-base text-slate-700 outline-none"
        />
      </div>
      <div className="flex gap-2">
        <button
          type="button"
          onClick={onToggleFilters}
          className="inline-flex items-center gap-2 rounded-2xl border border-slate-200 px-4 py-2 text-sm font-medium text-slate-600 hover:border-primary-200 hover:text-primary-700"
        >
          <SlidersHorizontal size={18} />
          Filtros
        </button>
        <button
          type="submit"
          className="inline-flex items-center gap-2 rounded-2xl bg-primary-600 px-6 py-2 text-sm font-semibold text-white shadow-primary-200 hover:bg-primary-700"
        >
          Buscar
        </button>
      </div>
    </form>
  );
};
