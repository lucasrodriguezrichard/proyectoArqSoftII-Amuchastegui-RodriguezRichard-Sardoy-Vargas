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
    } catch {
      // handled by toast inside the mutation hook
    }
  };

  return (
    <div className="mx-auto max-w-3xl px-4 py-12">
      <section className="rounded-3xl border border-white/50 bg-white/95 p-8 text-slate-900 shadow-xl shadow-primary-500/10 backdrop-blur-xl dark:border-white/10 dark:bg-gradient-to-br dark:from-slate-950 dark:via-slate-900 dark:to-slate-800 dark:text-white">
        <p className="flex items-center gap-2 text-xs font-semibold uppercase tracking-[0.35em] text-slate-500 dark:text-white/70">
          <PlusCircle size={18} />
          Nuevas Reservas
        </p>
       
        <p className="mt-3 max-w-2xl text-sm text-slate-500 dark:text-white/80">
          Completá el formulario para crear una nueva reserva y mantener tu salón bajo control en segundos.
        </p>
      </section>

      <div className="mt-6 rounded-3xl border border-slate-100 bg-white/95 p-8 shadow-2xl shadow-slate-200/70 backdrop-blur-sm dark:border-white/10 dark:bg-slate-900/80">
        <CreateReservationForm onSubmit={handleSubmit} loading={createMutation.isPending} userId={String(user?.id)} />
      </div>
    </div>
  );
};

export default CreateReservation;
