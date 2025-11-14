import { searchApi } from './axios';

const SEARCH_PATH = '/api/search';

// Search for table availability (main entity)
export const searchTables = async (params = {}) => {
  const cleanParams = Object.entries(params).reduce((acc, [key, value]) => {
    if (value !== undefined && value !== null && value !== '') {
      acc[key] = value;
    }
    return acc;
  }, {});

  const { data } = await searchApi.get(SEARCH_PATH, { params: cleanParams });
  return data;
};

// Legacy name for backwards compatibility
export const searchReservations = searchTables;

export const getTableAvailability = async (id) => {
  const { data } = await searchApi.get(`${SEARCH_PATH}/${id}`);
  return data;
};

// Legacy name for backwards compatibility
export const getReservationDocument = getTableAvailability;

export const fetchSearchStats = async () => {
  const { data } = await searchApi.get(`${SEARCH_PATH}/stats`);
  return data;
};

export const triggerReindex = async () => {
  const { data } = await searchApi.post(`${SEARCH_PATH}/reindex`);
  return data;
};
