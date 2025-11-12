export const Loader = ({ label = 'Cargando...' }) => (
  <div className="flex flex-col items-center justify-center gap-2 py-12 text-slate-500">
    <span className="h-10 w-10 animate-spin rounded-full border-4 border-slate-200 border-t-primary-500" />
    <p className="text-sm font-medium">{label}</p>
  </div>
);
