import { useMemo, useState } from 'react';
import { useParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import toast from 'react-hot-toast';

import { getReservationDocument } from '../api/search';
import { getUserById } from '../api/auth';
import { ReservationDetails as Details } from '../components/reservation/ReservationDetails';
import { ConfirmModal } from '../components/reservation/ConfirmModal';
import { Loader } from '../components/common/Loader';
import { ErrorMessage } from '../components/common/ErrorMessage';
import { useAuth } from '../hooks/useAuth';
import { useConfirmReservation, useReservation } from '../hooks/useReservations';

const mongoIdRegex = /^[a-f\d]{24}$/i;

const ReservationDetailsPage = () => {
  const { id } = useParams();
  const { isAuthenticated, isAdmin, user } = useAuth();
  const [modalOpen, setModalOpen] = useState(false);

  const isMongoId = mongoIdRegex.test(id ?? '');

  const reservationQuery = useReservation(id, { enabled: Boolean(id) && isMongoId });
  const tableQuery = useQuery({
    queryKey: ['reservation-detail-solr', id],
    queryFn: () => getReservationDocument(id),
    enabled: Boolean(id) && !isMongoId,
  });

  const activeQuery = isMongoId ? reservationQuery : tableQuery;

  const confirmMutation = useConfirmReservation();

  const ownerId = useMemo(() => {
    const entity = activeQuery.data;
    if (!entity) return undefined;
    return entity.owner_id || entity.ownerId;
  }, [activeQuery.data]);

  const ownerQuery = useQuery({
    queryKey: ['reservation-owner', ownerId],
    queryFn: () => getUserById(ownerId),
    enabled: Boolean(ownerId),
    staleTime: 1000 * 60,
  });

  const ownerName = useMemo(() => {
    const owner = ownerQuery.data;
    if (!owner) return undefined;
    const composed = `${owner.first_name ?? ''} ${owner.last_name ?? ''}`.trim();
    return composed || owner.username || owner.email;
  }, [ownerQuery.data]);

  const currentUserId = user?.id ? String(user.id) : undefined;
  const isOwner = ownerId && currentUserId ? String(ownerId) === currentUserId : false;
  const canConfirm = isAuthenticated && (isOwner || isAdmin);

  const handleConfirm = async (payload) => {
    if (!canConfirm) {
      toast.error('No tenés permisos para confirmar esta reserva');
      return;
    }
    try {
      await confirmMutation.mutateAsync({ reservationId: id, payload });
      setModalOpen(false);
      if (isMongoId) {
        reservationQuery.refetch();
      } else {
        tableQuery.refetch();
      }
    } catch (error) {
      toast.error(error?.response?.data?.error ?? 'Error confirmando la reserva');
    }
  };

  if (activeQuery.isLoading) {
    return <Loader label="Cargando reserva..." />;
  }

  if (activeQuery.isError || !activeQuery.data) {
    return (
      <ErrorMessage
        message="No encontramos la reserva solicitada"
        actionLabel="Volver a intentar"
        onAction={() => activeQuery.refetch()}
      />
    );
  }

  return (
    <div className="mx-auto max-w-5xl px-4 py-8">
      <Details reservation={activeQuery.data} requesterName={ownerName} />
      {canConfirm ? (
        <div className="mt-6 elegant-card p-4">
          <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <p className="text-sm font-semibold text-slate-700 dark:text-slate-300">¿Listo para confirmar?</p>
              <p className="text-xs text-slate-500 dark:text-slate-400">
                Confirmá tu reserva y recibí la confirmación final.
              </p>
            </div>
            <button
              type="button"
              onClick={() => setModalOpen(true)}
              className="luxury-button text-sm"
            >
              Confirmar reserva
            </button>
          </div>
        </div>
      ) : isAuthenticated ? (
        <p className="mt-4 elegant-card px-4 py-3 text-center text-sm text-slate-500 dark:text-slate-400">
          Solo el solicitante o un administrador pueden confirmar esta reserva.
        </p>
      ) : (
        <p className="mt-4 elegant-card px-4 py-3 text-center text-sm text-slate-500 dark:text-slate-400">
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
