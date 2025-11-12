export const APP_NAME = import.meta.env.VITE_APP_NAME || 'Restaurant Reservations';

export const MEAL_TYPES = [
  { value: 'breakfast', label: 'Desayuno' },
  { value: 'lunch', label: 'Almuerzo' },
  { value: 'dinner', label: 'Cena' },
  { value: 'event', label: 'Evento' },
];

export const RESERVATION_STATUSES = [
  { value: 'pending', label: 'Pendiente' },
  { value: 'confirmed', label: 'Confirmada' },
  { value: 'cancelled', label: 'Cancelada' },
  { value: 'completed', label: 'Completada' },
];

export const DEFAULT_PAGE_SIZE = 6;
