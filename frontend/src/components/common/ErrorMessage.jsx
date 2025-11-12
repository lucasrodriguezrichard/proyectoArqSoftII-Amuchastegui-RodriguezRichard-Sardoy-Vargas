export const ErrorMessage = ({ message = 'Ocurrió un error inesperado', actionLabel, onAction }) => (
  <div className="rounded-xl border border-rose-100 bg-rose-50/70 px-4 py-3 text-sm text-rose-700">
    <div className="flex items-center justify-between gap-4">
      <span>{message}</span>
      {actionLabel && onAction ? (
        <button
          type="button"
          onClick={onAction}
          className="rounded-full border border-rose-200 px-3 py-1 text-xs font-semibold text-rose-700 hover:bg-white"
        >
          {actionLabel}
        </button>
      ) : null}
    </div>
  </div>
);
