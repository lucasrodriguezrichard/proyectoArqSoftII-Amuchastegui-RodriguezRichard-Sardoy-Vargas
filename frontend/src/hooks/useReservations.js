import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import toast from 'react-hot-toast';

import {
  confirmReservation,
  createReservation,
  deleteReservation,
  getReservationById,
  listReservations,
  listUserReservations,
  updateReservation,
} from '../api/reservations';

const invalidateReservationQueries = (queryClient) => {
  queryClient.invalidateQueries({ queryKey: ['reservations'] });
  queryClient.invalidateQueries({ queryKey: ['reservation'] });
};

const invalidateSearchQueries = (queryClient) => {
  queryClient.invalidateQueries({ queryKey: ['search-tables'] });
};

export const useReservations = (params = {}, options = {}) =>
  useQuery({
    queryKey: ['reservations', params],
    queryFn: () => listReservations(params),
    ...options,
  });

export const useUserReservations = (userId, options = {}) =>
  useQuery({
    queryKey: ['reservations', 'user', userId],
    queryFn: () => listUserReservations(userId),
    enabled: Boolean(userId) && (options?.enabled ?? true),
    ...options,
  });

export const useReservation = (reservationId, options = {}) =>
  useQuery({
    queryKey: ['reservation', reservationId],
    queryFn: () => getReservationById(reservationId),
    enabled: Boolean(reservationId) && (options?.enabled ?? true),
    ...options,
  });

export const useCreateReservation = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: createReservation,
    onSuccess: () => {
      invalidateReservationQueries(queryClient);
      invalidateSearchQueries(queryClient);
      toast.success('Reserva creada');
    },
    onError: () => toast.error('No pudimos crear la reserva'),
  });
};

export const useUpdateReservation = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: updateReservation,
    onSuccess: () => {
      invalidateReservationQueries(queryClient);
      invalidateSearchQueries(queryClient);
      toast.success('Reserva actualizada');
    },
    onError: () => toast.error('No pudimos actualizar la reserva'),
  });
};

export const useDeleteReservation = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: deleteReservation,
    onSuccess: () => {
      invalidateReservationQueries(queryClient);
      invalidateSearchQueries(queryClient);
      toast.success('Reserva eliminada');
    },
    onError: () => toast.error('No pudimos eliminar la reserva'),
  });
};

export const useConfirmReservation = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: confirmReservation,
    onSuccess: () => {
      invalidateReservationQueries(queryClient);
      invalidateSearchQueries(queryClient);
      toast.success('Reserva confirmada');
    },
    onError: () => toast.error('No pudimos confirmar la reserva'),
  });
};
