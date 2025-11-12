import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { useNavigate } from 'react-router-dom';
import { LogIn } from 'lucide-react';
import { useAuth } from '../../hooks/useAuth';

export const LoginForm = () => {
  const { login, loading } = useAuth();
  const navigate = useNavigate();
  const [formError, setFormError] = useState(null);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm({
    defaultValues: {
      identifier: '',
      password: '',
    },
  });

  const onSubmit = handleSubmit(async (values) => {
    setFormError(null);
    try {
      await login(values);
      navigate('/');
    } catch (error) {
      setFormError(error?.response?.data?.error ?? 'Credenciales inválidas');
    }
  });

  return (
    <form onSubmit={onSubmit} className="space-y-4">
      <div>
        <label className="text-sm font-medium text-slate-600">Email o usuario</label>
        <input
          type="text"
          placeholder="usuario@ejemplo.com"
          className="mt-1 w-full rounded-xl border border-slate-200 bg-white px-4 py-3 text-slate-700 focus:border-primary-400 focus:outline-none focus:ring-2 focus:ring-primary-100"
          {...register('identifier', { required: 'Campo obligatorio' })}
        />
        {errors.identifier && <p className="mt-1 text-xs text-rose-500">{errors.identifier.message}</p>}
      </div>
      <div>
        <label className="text-sm font-medium text-slate-600">Contraseña</label>
        <input
          type="password"
          placeholder="********"
          className="mt-1 w-full rounded-xl border border-slate-200 bg-white px-4 py-3 text-slate-700 focus:border-primary-400 focus:outline-none focus:ring-2 focus:ring-primary-100"
          {...register('password', { required: 'Campo obligatorio' })}
        />
        {errors.password && <p className="mt-1 text-xs text-rose-500">{errors.password.message}</p>}
      </div>
      {formError && <p className="rounded-xl border border-rose-100 bg-rose-50 px-3 py-2 text-sm text-rose-600">{formError}</p>}
      <button
        type="submit"
        disabled={loading}
        className="flex w-full items-center justify-center gap-2 rounded-xl bg-primary-600 px-4 py-3 font-semibold text-white transition hover:bg-primary-700 disabled:cursor-not-allowed disabled:bg-slate-300"
      >
        <LogIn size={18} />
        {loading ? 'Ingresando...' : 'Ingresar'}
      </button>
    </form>
  );
};
