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
    <div className="mx-auto max-w-3xl px-4 py-12">
      <section className="rounded-3xl border border-primary-100 bg-gradient-to-br from-primary-50 to-white p-8 dark:border-primary-800/50 dark:bg-slate-800/80">
        <p className="flex items-center gap-2 text-sm font-light uppercase tracking-[0.35em] text-primary-700 dark:text-primary-400">
          <PlusCircle size={18} />
          Nueva reserva
        </p>
        <h1 className="section-title mt-3 text-3xl">Crear reserva</h1>
        <p className="subtitle mt-3">
          Completá el formulario para crear una nueva reserva.
        </p>
      </section>

      <div className="mt-6 elegant-card p-8">
        <CreateReservationForm onSubmit={handleSubmit} loading={createMutation.isPending} userId={String(user?.id)} />
      </div>
    </div>
  );
};

export default CreateReservation;
