import { useMemo, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Sparkles } from 'lucide-react';

import { searchReservations } from '../api/search';
import { DEFAULT_PAGE_SIZE } from '../utils/constants';
import { SearchBar } from '../components/search/SearchBar';
import { FilterPanel } from '../components/search/FilterPanel';
import { ReservationCard } from '../components/search/ReservationCard';
import { Pagination } from '../components/search/Pagination';
import { Loader } from '../components/common/Loader';
import { ErrorMessage } from '../components/common/ErrorMessage';
import { useAuth } from '../hooks/useAuth';

const Home = () => {
  const { isAuthenticated } = useAuth();
  const [query, setQuery] = useState('');
  const [page, setPage] = useState(1);
  const [filters, setFilters] = useState({ meal_type: '', status: '', guests: '' });
  const [showFilters, setShowFilters] = useState(false);

  const params = useMemo(
    () => ({
      q: query || '*:*',
      page,
      size: DEFAULT_PAGE_SIZE,
      ...filters,
    }),
    [filters, page, query],
  );

  const { data, isLoading, isError, refetch, isFetching } = useQuery({
    queryKey: ['search', params],
    queryFn: () => searchReservations(params),
    staleTime: 1000 * 30,
  });

  const handleSearch = (value) => {
    setQuery(value);
    setPage(1);
  };

  const handleFilterChange = (nextFilters) => {
    setFilters(nextFilters);
    setPage(1);
  };

  const results = data?.results ?? [];

  return (
    <div className="mx-auto max-w-6xl px-4 py-8">
      <section className="rounded-3xl bg-gradient-to-br from-primary-600 via-primary-500 to-primary-400 p-6 text-white shadow-soft">
        <p className="flex items-center gap-2 text-sm font-semibold uppercase tracking-[0.3em]">
          <Sparkles size={16} />
          Reservas Inteligentes
        </p>
        <h1 className="mt-2 font-display text-4xl font-semibold">Encontrá la mesa ideal</h1>
        <p className="mt-2 max-w-2xl text-white/80">
          Consultá disponibilidades en tiempo real, confirmá reservas y mantené el control del salón desde una sola
          interfaz.
        </p>
      </section>

      <div className="mt-6 space-y-4">
        <SearchBar
          initialQuery={query}
          onSearch={handleSearch}
          onToggleFilters={() => setShowFilters((prev) => !prev)}
        />
        <FilterPanel filters={filters} onChange={handleFilterChange} visible={showFilters} />
      </div>

      {isLoading || isFetching ? (
        <Loader label="Buscando reservas..." />
      ) : isError ? (
        <ErrorMessage message="No pudimos obtener los resultados" actionLabel="Reintentar" onAction={() => refetch()} />
      ) : (
        <>
          <div className="mt-6 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {results.length ? (
              results.map((reservation) => <ReservationCard key={reservation.id} reservation={reservation} />)
            ) : (
              <p className="rounded-2xl border border-slate-100 bg-white px-4 py-6 text-center text-sm text-slate-500 sm:col-span-2 lg:col-span-3">
                No encontramos reservas con esos filtros.
              </p>
            )}
          </div>
          <Pagination page={data?.page ?? 1} pages={data?.pages ?? 1} onChange={setPage} />
        </>
      )}

      {!isAuthenticated && (
        <div className="mt-10 rounded-3xl border border-dashed border-primary-300 bg-white/60 p-6 text-center text-slate-600">
          <p className="text-lg font-semibold text-slate-900">¿Querés confirmar una reserva?</p>
          <p className="mt-2 text-sm">Creá una cuenta o iniciá sesión para acceder a tus reservas y panel admin.</p>
        </div>
      )}
    </div>
  );
};

export default Home;
