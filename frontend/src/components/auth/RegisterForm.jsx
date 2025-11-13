import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { UserPlus } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';

export const RegisterForm = () => {
  const { register: registerUser, loading } = useAuth();
  const navigate = useNavigate();
  const [errorMessage, setErrorMessage] = useState(null);
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm({
    defaultValues: {
      username: '',
      email: '',
      first_name: '',
      last_name: '',
      password: '',
    },
  });

  const onSubmit = handleSubmit(async (values) => {
    setErrorMessage(null);
    try {
      await registerUser(values);
      navigate('/');
    } catch (error) {
      const apiError = error?.response?.data?.error;
      const friendlyMessage =
        apiError === 'invalid_input'
          ? 'Datos inválidos. Revisá el usuario, email y que la contraseña tenga al menos 8 caracteres.'
          : apiError === 'user_already_exists'
          ? 'Ese usuario o email ya está registrado.'
          : 'No pudimos registrar tu usuario';
      setErrorMessage(friendlyMessage);
    }
  });

  const renderField = (name, label, type = 'text', rules = { required: 'Campo obligatorio' }) => (
    <div>
      <label className="text-sm font-medium text-slate-600 dark:text-slate-300">{label}</label>
      <input
        type={type}
        className="luxury-input mt-1"
        {...register(name, rules)}
      />
      {errors[name] && <p className="mt-1 text-xs text-rose-500 dark:text-rose-400">{errors[name].message}</p>}
    </div>
  );

  return (
    <form onSubmit={onSubmit} className="space-y-4">
      <div className="grid gap-4 sm:grid-cols-2">
        {renderField('username', 'Usuario', 'text', {
          required: 'Campo obligatorio',
          pattern: {
            value: /^[a-zA-Z0-9._-]{3,32}$/,
            message: 'Usá 3-32 caracteres sin espacios',
          },
        })}
        {renderField('email', 'Email', 'email', {
          required: 'Campo obligatorio',
          pattern: {
            value: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
            message: 'Ingresá un email válido',
          },
        })}
        {renderField('first_name', 'Nombre')}
        {renderField('last_name', 'Apellido')}
      </div>
      {renderField('password', 'Contraseña', 'password', {
        required: 'Campo obligatorio',
        minLength: { value: 8, message: 'Al menos 8 caracteres' },
      })}
      {errorMessage && (
        <p className="rounded-xl border border-rose-100 bg-rose-50 px-3 py-2 text-sm text-rose-600 dark:border-rose-800 dark:bg-rose-950/30 dark:text-rose-400">
          {errorMessage}
        </p>
      )}
      <button
        type="submit"
        disabled={loading}
        className="luxury-button flex w-full items-center justify-center gap-2"
      >
        <UserPlus size={18} />
        {loading ? 'Creando cuenta...' : 'Crear cuenta'}
      </button>
    </form>
  );
};
