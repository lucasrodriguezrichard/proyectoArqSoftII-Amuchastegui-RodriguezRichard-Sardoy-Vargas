import { Navigate } from 'react-router-dom';
import { useState } from 'react';
import { ShieldCheck } from 'lucide-react';

import { useAuth } from '../hooks/useAuth';
import { useReservations, useUpdateReservation, useDeleteReservation } from '../hooks/useReservations';
import { Loader } from '../components/common/Loader';
import { ErrorMessage } from '../components/common/ErrorMessage';
import { ReservationTable } from '../components/admin/ReservationTable';
import { EditModal } from '../components/admin/EditModal';

const Admin = () => {
  const { isAuthenticated, isAdmin } = useAuth();
  const [modalOpen, setModalOpen] = useState(false);
  const [selectedReservation, setSelectedReservation] = useState(null);

  const reservationsQuery = useReservations();
  const updateMutation = useUpdateReservation();
  const deleteMutation = useDeleteReservation();

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  if (!isAdmin) {
    return <Navigate to="/" replace />;
  }

  const handleEdit = (reservation) => {
    setSelectedReservation(reservation);
    setModalOpen(true);
  };

  const handleDelete = async (reservation) => {
    const confirmed = window.confirm(`Eliminar reserva ${reservation.id}?`);
    if (!confirmed) return;
    await deleteMutation.mutateAsync(reservation.id);
    reservationsQuery.refetch();
  };

  const handleSave = async ({ reservationId, payload }) => {
    await updateMutation.mutateAsync({ reservationId, payload });
    setModalOpen(false);
    reservationsQuery.refetch();
  };

  if (reservationsQuery.isLoading) {
    return <Loader label="Cargando reservas..." />;
  }

  if (reservationsQuery.isError) {
    return <ErrorMessage message="No pudimos cargar las reservas" actionLabel="Reintentar" onAction={() => reservationsQuery.refetch()} />;
  }

  const reservations = reservationsQuery.data ?? [];

  return (
    <div className="mx-auto max-w-6xl px-4 py-8">
      <section className="rounded-3xl border border-primary-100 bg-primary-50 p-6 text-primary-900">
        <p className="flex items-center gap-2 text-sm font-semibold uppercase tracking-[0.3em]">
          <ShieldCheck size={18} />
          Panel admin
        </p>
        <h1 className="mt-2 font-display text-3xl font-semibold">Control total del salón</h1>
        <p className="mt-2 text-primary-900/70">Editá, confirmá o eliminá reservas. Todos los cambios se replican en MongoDB y Solr.</p>
      </section>

      <div className="mt-6">
        <ReservationTable reservations={reservations} onEdit={handleEdit} onDelete={handleDelete} />
      </div>

      <EditModal
        open={modalOpen}
        reservation={selectedReservation}
        loading={updateMutation.isPending}
        onClose={() => setModalOpen(false)}
        onSave={handleSave}
      />
    </div>
  );
};

export default Admin;
