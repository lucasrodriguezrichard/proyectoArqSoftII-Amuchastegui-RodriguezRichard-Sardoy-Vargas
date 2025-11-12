import { useState } from 'react';
import { useParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import toast from 'react-hot-toast';

import { getReservationDocument } from '../api/search';
import { ReservationDetails as Details } from '../components/reservation/ReservationDetails';
import { ConfirmModal } from '../components/reservation/ConfirmModal';
import { Loader } from '../components/common/Loader';
import { ErrorMessage } from '../components/common/ErrorMessage';
import { useAuth } from '../hooks/useAuth';
import { useConfirmReservation } from '../hooks/useReservations';

const ReservationDetailsPage = () => {
  const { id } = useParams();
  const { isAuthenticated } = useAuth();
  const [modalOpen, setModalOpen] = useState(false);

  const query = useQuery({
    queryKey: ['reservation-detail', id],
    queryFn: () => getReservationDocument(id),
    enabled: Boolean(id),
  });

  const confirmMutation = useConfirmReservation();

  const handleConfirm = async (payload) => {
    try {
      await confirmMutation.mutateAsync({ reservationId: id, payload });
      setModalOpen(false);
      query.refetch();
    } catch (error) {
      toast.error(error?.response?.data?.error ?? 'Error confirmando la reserva');
    }
  };

  if (query.isLoading) {
    return <Loader label="Cargando reserva..." />;
  }

  if (query.isError || !query.data) {
    return <ErrorMessage message="No encontramos la reserva solicitada" actionLabel="Volver a intentar" onAction={() => query.refetch()} />;
  }

  return (
    <div className="mx-auto max-w-5xl px-4 py-8">
      <Details reservation={query.data} />
      {isAuthenticated ? (
        <div className="mt-6 rounded-3xl border border-slate-100 bg-white/80 p-4 shadow-soft">
          <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <p className="text-sm font-semibold text-slate-700">¿Listo para confirmar?</p>
              <p className="text-xs text-slate-500">
                Ejecutará el cálculo concurrente en el backend y notificará vía RabbitMQ.
              </p>
            </div>
            <button
              type="button"
              onClick={() => setModalOpen(true)}
              className="rounded-full bg-primary-600 px-5 py-2 text-sm font-semibold text-white shadow-soft hover:bg-primary-700"
            >
              Confirmar reserva
            </button>
          </div>
        </div>
      ) : (
        <p className="mt-4 rounded-2xl border border-dashed border-slate-200 px-4 py-3 text-center text-sm text-slate-500">
          Iniciá sesión para confirmar esta reserva.
        </p>
      )}
      <ConfirmModal
        open={modalOpen}
        loading={confirmMutation.isPending}
        onClose={() => setModalOpen(false)}
        onConfirm={handleConfirm}
      />
    </div>
  );
};

export default ReservationDetailsPage;
