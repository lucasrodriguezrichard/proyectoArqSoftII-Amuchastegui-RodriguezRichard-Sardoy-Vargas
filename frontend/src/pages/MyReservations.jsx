import { Navigate, Link } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import { useUserReservations } from '../hooks/useReservations';
import { Loader } from '../components/common/Loader';
import { ErrorMessage } from '../components/common/ErrorMessage';
import { ReservationCard } from '../components/search/ReservationCard';

const MyReservations = () => {
  const { isAuthenticated, user } = useAuth();
  const query = useUserReservations(user?.id, { enabled: isAuthenticated });

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  if (query.isLoading) {
    return <Loader label="Cargando tus reservas..." />;
  }

  if (query.isError) {
    return <ErrorMessage message="No pudimos cargar tus reservas" actionLabel="Reintentar" onAction={() => query.refetch()} />;
  }

  const reservations = query.data ?? [];

  return (
    <div className="mx-auto max-w-6xl px-4 py-8">
      <div className="flex flex-col gap-2">
        <p className="text-sm uppercase tracking-[0.3em] text-slate-400 dark:text-slate-500">Tus reservas</p>
        <h1 className="font-display text-3xl font-semibold text-slate-900 dark:text-slate-50">Hola, {user?.first_name || user?.username}</h1>
        <p className="text-slate-500 dark:text-slate-400">Gestioná tus reservas pendientes, confirmalas o revisá el detalle.</p>
      </div>

      <div className="mt-6 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {reservations.length ? (
          reservations.map((reservation) => <ReservationCard key={reservation.id} reservation={reservation} />)
        ) : (
          <div className="elegant-card px-4 py-6 text-center text-sm text-slate-500 dark:text-slate-400 sm:col-span-2 lg:col-span-3">
            Aún no tenés reservas.{' '}
            <Link to="/" className="font-semibold text-primary-600 dark:text-primary-400">
              Creá la primera
            </Link>
            .
          </div>
        )}
      </div>
    </div>
  );
};

export default MyReservations;
