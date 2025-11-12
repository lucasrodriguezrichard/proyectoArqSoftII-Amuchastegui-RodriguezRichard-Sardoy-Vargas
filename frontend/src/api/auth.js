import { usersApi } from './axios';

export const login = async (credentials) => {
  const { data } = await usersApi.post('/api/users/login', credentials);
  return data;
};

export const registerUser = async (payload) => {
  const { data } = await usersApi.post('/api/users/register', payload);
  return data;
};

export const getUserById = async (userId) => {
  const { data } = await usersApi.get(`/api/users/${userId}`);
  return data;
};
