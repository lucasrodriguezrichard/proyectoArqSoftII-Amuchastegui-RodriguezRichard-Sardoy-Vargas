import { reservationsApi } from './axios';

const BASE_PATH = '/api/reservations';

export const listReservations = async (params = {}) => {
  const { data } = await reservationsApi.get(BASE_PATH, { params });
  return data;
};

export const listUserReservations = async (userId) => {
  if (!userId) {
    return [];
  }
  const { data } = await reservationsApi.get(`${BASE_PATH}/user/${userId}`);
  return data;
};

export const getReservationById = async (reservationId) => {
  const { data } = await reservationsApi.get(`${BASE_PATH}/${reservationId}`);
  return data;
};

export const createReservation = async (payload) => {
  const { data } = await reservationsApi.post(BASE_PATH, payload);
  return data;
};

export const updateReservation = async ({ reservationId, payload }) => {
  const { data } = await reservationsApi.put(`${BASE_PATH}/${reservationId}`, payload);
  return data;
};

export const deleteReservation = async (reservationId) => {
  await reservationsApi.delete(`${BASE_PATH}/${reservationId}`);
};

export const confirmReservation = async ({ reservationId, payload }) => {
  const { data } = await reservationsApi.post(`${BASE_PATH}/${reservationId}/confirm`, payload);
  return data;
};
