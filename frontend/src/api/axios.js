import axios from 'axios';

export const STORAGE_KEYS = Object.freeze({
  token: 'token',
  user: 'user',
});

const API_TIMEOUT = 10000;
const AUTH_SAFE_PATHS = ['/api/users/login', '/api/users/register'];

const shouldSkipAuthHandling = (config) =>
  AUTH_SAFE_PATHS.some((path) => config?.url?.includes(path));

const attachInterceptors = (instance) => {
  instance.interceptors.request.use(
    (config) => {
      const token = localStorage.getItem(STORAGE_KEYS.token);
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    },
    (error) => Promise.reject(error),
  );

  instance.interceptors.response.use(
    (response) => response,
    (error) => {
      const status = error.response?.status;
      if (status === 401 && !shouldSkipAuthHandling(error.config)) {
        localStorage.removeItem(STORAGE_KEYS.token);
        localStorage.removeItem(STORAGE_KEYS.user);
        if (typeof window !== 'undefined' && window.location.pathname !== '/login') {
          window.location.href = '/login';
        }
      }
      return Promise.reject(error);
    },
  );
  return instance;
};

const createClient = (baseURL) =>
  attachInterceptors(
    axios.create({
      baseURL,
      timeout: API_TIMEOUT,
      headers: { 'Content-Type': 'application/json' },
    }),
  );

export const usersApi = createClient(import.meta.env.VITE_API_URL || 'http://localhost:8080');
export const reservationsApi = createClient(
  import.meta.env.VITE_RESERVATIONS_API_URL || 'http://localhost:8081',
);
export const searchApi = createClient(
  import.meta.env.VITE_SEARCH_API_URL || 'http://localhost:8082',
);
