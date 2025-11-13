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
    <form onSubmit={handleSubmit} className="rounded-3xl border border-white/60 bg-white/70 p-5 shadow-2xl shadow-primary-500/5 backdrop-blur-xl dark:border-white/10 dark:bg-white/5 sm:flex sm:items-center sm:gap-4">
      <div className="flex flex-1 items-center gap-3 rounded-2xl border border-white/60 bg-white/90 px-4 py-3 text-slate-700 shadow-inner shadow-white/70 dark:border-white/10 dark:bg-white/10 dark:text-slate-100">
        <Search className="text-primary-500 dark:text-primary-300" size={20} />
        <input
          value={query}
          onChange={(event) => setQuery(event.target.value)}
          placeholder="Buscar por nombre, cliente, mesa..."
          className="w-full border-none bg-transparent text-base text-slate-700 outline-none placeholder:text-slate-400 dark:text-slate-100 dark:placeholder:text-slate-500"
        />
      </div>
      <div className="mt-3 flex gap-2 sm:mt-0">
        <button
          type="button"
          onClick={onToggleFilters}
          className="inline-flex items-center gap-2 rounded-2xl border border-white/40 px-5 py-3 text-sm font-semibold text-slate-700 transition hover:-translate-y-0.5 hover:border-primary-300 hover:bg-white/70 dark:border-white/10 dark:text-slate-200 dark:hover:bg-white/10"
        >
          <SlidersHorizontal size={18} />
          Filtros
        </button>
        <button type="submit" className="rounded-2xl bg-gradient-to-r from-primary-500 via-sky-500 to-emerald-400 px-6 py-3 text-sm font-semibold uppercase tracking-wider text-white shadow-lg shadow-primary-500/40 transition hover:-translate-y-0.5">
          Buscar
        </button>
      </div>
    </form>
  );
};
