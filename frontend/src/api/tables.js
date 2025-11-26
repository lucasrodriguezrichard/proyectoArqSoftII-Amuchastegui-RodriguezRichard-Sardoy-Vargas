import { reservationsApi } from './axios';

export const getTables = async (mealType = '') => {
  const params = mealType ? { meal_type: mealType } : {};
  const response = await reservationsApi.get('/api/tables', { params });
  return response.data;
};

export const getTableById = async (id) => {
  const response = await reservationsApi.get(`/api/tables/${id}`);
  return response.data;
};

export const createTable = async (tableData) => {
  const response = await reservationsApi.post('/api/admin/tables', tableData);
  return response.data;
};

export const updateTable = async (id, tableData) => {
  const response = await reservationsApi.put(`/api/admin/tables/${id}`, tableData);
  return response.data;
};

export const deleteTable = async (id) => {
  await reservationsApi.delete(`/api/admin/tables/${id}`);
};
