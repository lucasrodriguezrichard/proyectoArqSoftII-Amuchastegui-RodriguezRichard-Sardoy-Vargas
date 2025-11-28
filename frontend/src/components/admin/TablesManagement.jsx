import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Plus, Pencil, Trash2 } from 'lucide-react';
import toast from 'react-hot-toast';

import { getTables, createTable, updateTable, deleteTable } from '../../api/tables';

export const TablesManagement = () => {
  const queryClient = useQueryClient();
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingTable, setEditingTable] = useState(null);
  const [formData, setFormData] = useState({
    table_number: '',
    capacity: '',
    meal_type: 'breakfast',
  });

  const { data: tables = [], isLoading, refetch } = useQuery({
    queryKey: ['tables'],
    queryFn: () => getTables(),
  });

  const createMutation = useMutation({
    mutationFn: createTable,
    onSuccess: () => {
      toast.success('Mesa creada exitosamente');
      queryClient.invalidateQueries({ queryKey: ['tables'] });
      queryClient.invalidateQueries({ queryKey: ['search-tables'], exact: false });
      resetForm();
    },
    onError: (error) => {
      toast.error(error.response?.data?.error || 'Error al crear la mesa');
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }) => updateTable(id, data),
    onSuccess: () => {
      toast.success('Mesa actualizada exitosamente');
      queryClient.invalidateQueries({ queryKey: ['tables'] });
      queryClient.invalidateQueries({ queryKey: ['search-tables'], exact: false });
      resetForm();
    },
    onError: (error) => {
      toast.error(error.response?.data?.error || 'Error al actualizar la mesa');
    },
  });

  const deleteMutation = useMutation({
    mutationFn: deleteTable,
    onSuccess: () => {
      toast.success('Mesa eliminada exitosamente');
      queryClient.invalidateQueries({ queryKey: ['tables'] });
      queryClient.invalidateQueries({ queryKey: ['search-tables'], exact: false });
    },
    onError: (error) => {
      toast.error(error.response?.data?.error || 'Error al eliminar la mesa');
    },
  });

  const resetForm = () => {
    setFormData({ table_number: '', capacity: '', meal_type: 'breakfast' });
    setEditingTable(null);
    setIsFormOpen(false);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    const payload = {
      table_number: parseInt(formData.table_number),
      capacity: parseInt(formData.capacity),
      meal_type: formData.meal_type,
    };

    if (editingTable) {
      await updateMutation.mutateAsync({ id: editingTable.id, data: payload });
    } else {
      await createMutation.mutateAsync(payload);
    }
  };

  const handleEdit = (table) => {
    setEditingTable(table);
    setFormData({
      table_number: table.table_number.toString(),
      capacity: table.capacity.toString(),
      meal_type: table.meal_type,
    });
    setIsFormOpen(true);
  };

  const handleDelete = async (table) => {
    if (window.confirm(`¿Eliminar mesa ${table.table_number} (${table.meal_type})?`)) {
      await deleteMutation.mutateAsync(table.id);
    }
  };

  const mealTypeLabels = {
    breakfast: 'Desayuno',
    lunch: 'Almuerzo',
    dinner: 'Cena',
    event: 'Evento',
  };

  if (isLoading) {
    return <div className="text-center py-8">Cargando mesas...</div>;
  }

  return (
    <div className="space-y-6">
      {/* Header with Add Button */}
      <div className="flex justify-between items-center">
        <h2 className="text-2xl font-semibold text-slate-800">Gestión de Mesas</h2>
        <button
          onClick={() => setIsFormOpen(true)}
          className="flex items-center gap-2 rounded-xl bg-primary-500 px-4 py-2 text-sm font-semibold text-white hover:bg-primary-600 transition"
        >
          <Plus size={18} />
          Nueva Mesa
        </button>
      </div>

      {/* Form */}
      {isFormOpen && (
        <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-lg">
          <h3 className="text-lg font-semibold text-slate-800 mb-4">
            {editingTable ? 'Editar Mesa' : 'Nueva Mesa'}
          </h3>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">
                  Número de Mesa
                </label>
                <input
                  type="number"
                  required
                  min="1"
                  value={formData.table_number}
                  onChange={(e) => setFormData({ ...formData, table_number: e.target.value })}
                  className="w-full rounded-lg border border-slate-300 px-3 py-2 text-slate-800 placeholder-slate-400 focus:border-primary-500 focus:ring-2 focus:ring-primary-200 outline-none dark:bg-slate-900 dark:text-slate-100 dark:placeholder-slate-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">
                  Capacidad
                </label>
                <input
                  type="number"
                  required
                  min="1"
                  max="50"
                  value={formData.capacity}
                  onChange={(e) => setFormData({ ...formData, capacity: e.target.value })}
                  className="w-full rounded-lg border border-slate-300 px-3 py-2 text-slate-800 placeholder-slate-400 focus:border-primary-500 focus:ring-2 focus:ring-primary-200 outline-none dark:bg-slate-900 dark:text-slate-100 dark:placeholder-slate-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">
                  Tipo de Comida
                </label>
                <select
                  value={formData.meal_type}
                  onChange={(e) => setFormData({ ...formData, meal_type: e.target.value })}
                  className="w-full rounded-lg border border-slate-300 px-3 py-2 focus:border-primary-500 focus:ring-2 focus:ring-primary-200 outline-none"
                >
                  <option value="breakfast">Desayuno</option>
                  <option value="lunch">Almuerzo</option>
                  <option value="dinner">Cena</option>
                  <option value="event">Evento</option>
                </select>
              </div>
            </div>

            <div className="flex gap-3">
              <button
                type="submit"
                disabled={createMutation.isPending || updateMutation.isPending}
                className="flex-1 rounded-lg bg-primary-500 px-4 py-2 text-sm font-semibold text-white hover:bg-primary-600 transition disabled:opacity-50"
              >
                {editingTable ? 'Actualizar' : 'Crear'}
              </button>
              <button
                type="button"
                onClick={resetForm}
                className="rounded-lg border border-slate-300 px-4 py-2 text-sm font-semibold text-slate-700 hover:bg-slate-50 transition"
              >
                Cancelar
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Tables List */}
      <div className="rounded-2xl border border-slate-200 bg-white overflow-hidden shadow-lg">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-slate-50 border-b border-slate-200">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-semibold text-slate-600 uppercase tracking-wider">
                  Número
                </th>
                <th className="px-6 py-3 text-left text-xs font-semibold text-slate-600 uppercase tracking-wider">
                  Capacidad
                </th>
                <th className="px-6 py-3 text-left text-xs font-semibold text-slate-600 uppercase tracking-wider">
                  Tipo de Comida
                </th>
                <th className="px-6 py-3 text-right text-xs font-semibold text-slate-600 uppercase tracking-wider">
                  Acciones
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-200">
              {tables.length === 0 ? (
                <tr>
                  <td colSpan="4" className="px-6 py-8 text-center text-slate-500">
                    No hay mesas creadas. Haz clic en "Nueva Mesa" para crear una.
                  </td>
                </tr>
              ) : (
                tables.map((table) => (
                  <tr key={table.id} className="hover:bg-slate-50 transition">
                    <td className="px-6 py-4 text-sm font-medium text-slate-900">
                      {table.table_number}
                    </td>
                    <td className="px-6 py-4 text-sm text-slate-700">
                      {table.capacity} personas
                    </td>
                    <td className="px-6 py-4 text-sm text-slate-700">
                      <span className="inline-flex rounded-full bg-primary-100 px-3 py-1 text-xs font-medium text-primary-800">
                        {mealTypeLabels[table.meal_type]}
                      </span>
                    </td>
                    <td className="px-6 py-4 text-right text-sm">
                      <div className="flex justify-end gap-2">
                        <button
                          onClick={() => handleEdit(table)}
                          className="rounded-lg p-2 text-blue-600 hover:bg-blue-50 transition"
                          title="Editar"
                        >
                          <Pencil size={16} />
                        </button>
                        <button
                          onClick={() => handleDelete(table)}
                          disabled={deleteMutation.isPending}
                          className="rounded-lg p-2 text-red-600 hover:bg-red-50 transition disabled:opacity-50"
                          title="Eliminar"
                        >
                          <Trash2 size={16} />
                        </button>
                      </div>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};
