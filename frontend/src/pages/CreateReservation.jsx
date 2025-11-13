import { Navigate, useNavigate } from 'react-router-dom';
import { PlusCircle } from 'lucide-react';
import { useAuth } from '../hooks/useAuth';
import { useCreateReservation } from '../hooks/useReservations';
import { CreateReservationForm } from '../components/reservation/CreateReservationForm';

const CreateReservation = () => {
  const { isAuthenticated, user } = useAuth();
  const navigate = useNavigate();
  const createMutation = useCreateReservation();

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  const handleSubmit = async (payload) => {
    try {
      await createMutation.mutateAsync(payload);
      navigate('/my-reservations');
    } catch (error) {
      // Error handling is done by the mutation hook with toast
    }
  };

  return (
    <div className="mx-auto max-w-3xl px-4 py-8">
      <section className="rounded-3xl border border-primary-100 bg-gradient-to-br from-primary-50 to-white p-6">
        <p className="flex items-center gap-2 text-sm font-semibold uppercase tracking-[0.3em] text-primary-700">
          <PlusCircle size={18} />
          Nueva reserva
        </p>
        <h1 className="mt-2 font-display text-3xl font-semibold text-slate-900">Crear reserva</h1>
        <p className="mt-2 text-slate-600">
          Completá el formulario para crear una nueva reserva. Los cambios se sincronizarán automáticamente con MongoDB, Solr y la caché.
        </p>
      </section>

      <div className="mt-6 rounded-3xl border border-slate-100 bg-white p-6 shadow-sm">
        <CreateReservationForm onSubmit={handleSubmit} loading={createMutation.isPending} userId={String(user?.id)} />
      </div>

      <div className="mt-6 rounded-2xl border border-dashed border-slate-200 bg-slate-50 p-4 text-sm text-slate-600">
        <p className="font-semibold text-slate-900">💡 Flujo event-driven</p>
        <p className="mt-1">
          Al crear la reserva, se guardará en MongoDB y se publicará un evento en RabbitMQ. La Search API consumirá ese evento y sincronizará Solr automáticamente.
        </p>
      </div>
    </div>
  );
};

export default CreateReservation;
