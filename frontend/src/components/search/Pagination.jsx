export const Pagination = ({ page, pages, onChange }) => {
  if (!pages || pages <= 1) return null;

  const canPrev = page > 1;
  const canNext = page < pages;

  const handleChange = (next) => {
    if (next >= 1 && next <= pages) {
      onChange?.(next);
    }
  };

  return (
    <div className="mt-6 flex items-center justify-between rounded-2xl border border-slate-200 bg-white px-4 py-2 text-sm text-slate-600">
      <span>
        Página {page} de {pages}
      </span>
      <div className="flex items-center gap-2">
        <button
          type="button"
          disabled={!canPrev}
          onClick={() => handleChange(page - 1)}
          className="rounded-full border border-slate-200 px-3 py-1 disabled:cursor-not-allowed disabled:opacity-50"
        >
          Anterior
        </button>
        <button
          type="button"
          disabled={!canNext}
          onClick={() => handleChange(page + 1)}
          className="rounded-full border border-slate-200 px-3 py-1 disabled:cursor-not-allowed disabled:opacity-50"
        >
          Siguiente
        </button>
      </div>
    </div>
  );
};
