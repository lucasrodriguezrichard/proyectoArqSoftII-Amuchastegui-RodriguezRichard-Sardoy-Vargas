import { searchApi } from './axios';

const SEARCH_PATH = '/api/search';

export const searchReservations = async (params = {}) => {
  const cleanParams = Object.entries(params).reduce((acc, [key, value]) => {
    if (value !== undefined && value !== null && value !== '') {
      acc[key] = value;
    }
    return acc;
  }, {});

  const { data } = await searchApi.get(SEARCH_PATH, { params: cleanParams });
  return data;
};

export const getReservationDocument = async (id) => {
  const { data } = await searchApi.get(`${SEARCH_PATH}/${id}`);
  return data;
};

export const fetchSearchStats = async () => {
  const { data } = await searchApi.get(`${SEARCH_PATH}/stats`);
  return data;
};

export const triggerReindex = async () => {
  const { data } = await searchApi.post(`${SEARCH_PATH}/reindex`);
  return data;
};
