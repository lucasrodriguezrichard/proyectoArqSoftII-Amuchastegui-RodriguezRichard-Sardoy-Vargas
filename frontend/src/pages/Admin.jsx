import { Navigate } from 'react-router-dom';
import { useState } from 'react';
import { ShieldCheck, Calendar, Utensils } from 'lucide-react';

import { useAuth } from '../hooks/useAuth';
import { useReservations, useUpdateReservation, useDeleteReservation } from '../hooks/useReservations';
import { Loader } from '../components/common/Loader';
import { ErrorMessage } from '../components/common/ErrorMessage';
import { ReservationTable } from '../components/admin/ReservationTable';
import { EditModal } from '../components/admin/EditModal';
import { TablesManagement } from '../components/admin/TablesManagement';

const Admin = () => {
  const { isAuthenticated, isAdmin } = useAuth();
  const [activeTab, setActiveTab] = useState('reservations');
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
        <p className="mt-2 text-primary-900/70">Gestioná reservas y mesas. Todos los cambios se replican en MongoDB y Solr.</p>
      </section>

      {/* Tabs */}
      <div className="mt-6 flex gap-2 border-b border-slate-200">
        <button
          onClick={() => setActiveTab('reservations')}
          className={`flex items-center gap-2 px-4 py-3 text-sm font-semibold transition border-b-2 ${
            activeTab === 'reservations'
              ? 'border-primary-500 text-primary-600'
              : 'border-transparent text-slate-600 hover:text-slate-900'
          }`}
        >
          <Calendar size={18} />
          Reservas
        </button>
        <button
          onClick={() => setActiveTab('tables')}
          className={`flex items-center gap-2 px-4 py-3 text-sm font-semibold transition border-b-2 ${
            activeTab === 'tables'
              ? 'border-primary-500 text-primary-600'
              : 'border-transparent text-slate-600 hover:text-slate-900'
          }`}
        >
          <Utensils size={18} />
          Mesas
        </button>
      </div>

      {/* Tab Content */}
      <div className="mt-6">
        {activeTab === 'reservations' ? (
          <ReservationTable reservations={reservations} onEdit={handleEdit} onDelete={handleDelete} />
        ) : (
          <TablesManagement />
        )}
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
