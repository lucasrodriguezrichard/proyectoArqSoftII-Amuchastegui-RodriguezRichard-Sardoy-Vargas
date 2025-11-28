const BASE_BY_MEAL = {
  breakfast: 15,
  lunch: 25,
  dinner: 40,
  event: 75,
};

const isEven = (str) => {
  if (!str || str.length === 0) return false;
  const last = str[str.length - 1];
  const n = Number(last);
  return Number.isInteger(n) && n % 2 === 0;
};

/**
 * Replica simplificada del cálculo de backend para obtener base y descuento.
 * Devuelve base, descuento% y precio final (usa total_price si viene de API).
 */
export const calculateReservationPricing = (reservation) => {
  if (!reservation) {
    return { basePrice: 0, discountPercent: 0, finalPrice: 0 };
  }

  const guests = Number(reservation.guests) || 0;
  const mealType = reservation.meal_type || reservation.mealType || '';
  const dateValue = reservation.date_time || reservation.dateTime;
  const ownerId = reservation.owner_id || reservation.ownerId || '';

  const basePerPerson = BASE_BY_MEAL[mealType] ?? 30;
 const basePrice = guests * basePerPerson;

  let discountPercent = 0;
  const date = dateValue ? new Date(dateValue) : null;
  if (date) {
    // Early bird para cenas antes de 18 hs
    if (mealType === 'dinner' && date.getHours() < 18) {
      discountPercent += 10;
    }
    // Descuento lunes-jueves
    const day = date.getDay(); // 0=Domingo, 1=Lunes
    if (day >= 1 && day <= 4) {
      discountPercent += 5;
    }
  }
  // Cliente "leal": ID termina en número par
  if (isEven(ownerId)) {
    discountPercent += 5;
  }

  // Descuento por grupo: desde 4 comensales, cada +2 suma +5% (4->5%, 6->10%, 8->15%, ...)
  if (guests >= 4) {
    const steps = Math.floor((guests - 4) / 2) + 1;
    discountPercent += steps * 5;
  }

  const computedFinal = basePrice - (basePrice * discountPercent) / 100;
  // Preferimos el valor calculado para reflejar descuentos; si algo falla, usamos el provisto por la API.
  const apiTotal = reservation.total_price ?? reservation.totalPrice;
  const finalPrice = Number.isFinite(apiTotal) ? Math.min(apiTotal, computedFinal) : computedFinal;

  return {
    basePrice,
    discountPercent,
    finalPrice,
  };
};
