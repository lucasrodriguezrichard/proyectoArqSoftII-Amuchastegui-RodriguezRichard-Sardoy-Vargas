import { Link, Navigate } from 'react-router-dom';
import { RegisterForm } from '../components/auth/RegisterForm';
import { useAuth } from '../hooks/useAuth';

const Register = () => {
  const { isAuthenticated } = useAuth();

  if (isAuthenticated) {
    return <Navigate to="/" replace />;
  }

  return (
    <div className="mx-auto flex min-h-[80vh] max-w-5xl flex-col items-center justify-center gap-8 px-4 py-12">
      <div className="text-center">
        <p className="text-sm uppercase tracking-[0.3em] text-slate-400">Crear cuenta</p>
        <h1 className="font-display text-4xl font-semibold text-slate-900">Sumate a la experiencia</h1>
        <p className="mt-2 text-slate-500">Registrate para gestionar reservas y acceder al panel.</p>
      </div>
      <div className="w-full max-w-2xl rounded-3xl border border-slate-100 bg-white/80 p-8 shadow-soft">
        <RegisterForm />
        <p className="mt-4 text-center text-sm text-slate-500">
          ¿Ya tenés cuenta?{' '}
          <Link to="/login" className="font-semibold text-primary-600">
            Iniciá sesión
          </Link>
        </p>
      </div>
    </div>
  );
};

export default Register;
