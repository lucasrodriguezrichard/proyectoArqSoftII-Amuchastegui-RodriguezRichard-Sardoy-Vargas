import { createContext, useCallback, useMemo, useState } from 'react';
import PropTypes from 'prop-types';
import toast from 'react-hot-toast';

import { login as loginRequest, registerUser as registerRequest } from '../api/auth';
import { STORAGE_KEYS } from '../api/axios';

export const AuthContext = createContext(undefined);

const safeParse = (value) => {
  try {
    return JSON.parse(value);
  } catch {
    return null;
  }
};

export const AuthProvider = ({ children }) => {
  const [token, setToken] = useState(() => localStorage.getItem(STORAGE_KEYS.token));
  const [user, setUser] = useState(() => safeParse(localStorage.getItem(STORAGE_KEYS.user)));
  const [loading, setLoading] = useState(false);

  const persistSession = useCallback((nextToken, nextUser) => {
    if (nextToken) {
      localStorage.setItem(STORAGE_KEYS.token, nextToken);
    }
    if (nextUser) {
      localStorage.setItem(STORAGE_KEYS.user, JSON.stringify(nextUser));
    }
    setToken(nextToken);
    setUser(nextUser);
  }, []);

  const clearSession = useCallback(() => {
    localStorage.removeItem(STORAGE_KEYS.token);
    localStorage.removeItem(STORAGE_KEYS.user);
    setToken(null);
    setUser(null);
  }, []);

  const handleLogin = useCallback(
    async (credentials) => {
      setLoading(true);
      try {
        const data = await loginRequest(credentials);
        const accessToken =
          data?.tokens?.access_token ||
          data?.tokens?.token ||
          data?.token ||
          data?.access_token;
        if (!accessToken || !data?.user) {
          throw new Error('invalid_login_response');
        }
        persistSession(accessToken, data.user);
        toast.success(`Bienvenido ${data.user.first_name ?? data.user.username ?? ''}`.trim());
        return data.user;
      } catch (error) {
        toast.error('No pudimos iniciar sesión. Verifica tus credenciales.');
        throw error;
      } finally {
        setLoading(false);
      }
    },
    [persistSession],
  );

  const handleRegister = useCallback(
    async (payload) => {
      setLoading(true);
      try {
        await registerRequest(payload);
        toast.success('Registro exitoso. Iniciando sesión...');
        await handleLogin({
          identifier: payload.email || payload.username,
          password: payload.password,
        });
      } catch (error) {
        const apiError = error?.response?.data?.error;
        if (apiError === 'invalid_input') {
          toast.error('Datos inválidos: revisá usuario/email y contraseña (mínimo 8 caracteres).');
        } else if (apiError === 'user_already_exists') {
          toast.error('Ese usuario o email ya existe.');
        } else {
          toast.error('No pudimos completar el registro.');
        }
        throw error;
      } finally {
        setLoading(false);
      }
    },
    [handleLogin],
  );

  const handleLogout = useCallback(() => {
    clearSession();
    toast.success('Sesión cerrada');
  }, [clearSession]);

  const value = useMemo(
    () => ({
      token,
      user,
      loading,
      isAuthenticated: Boolean(token && user),
      isAdmin: user?.role === 'admin',
      login: handleLogin,
      register: handleRegister,
      logout: handleLogout,
      setUser,
    }),
    [handleLogin, handleLogout, handleRegister, token, user, loading],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

AuthProvider.propTypes = {
  children: PropTypes.node.isRequired,
};
