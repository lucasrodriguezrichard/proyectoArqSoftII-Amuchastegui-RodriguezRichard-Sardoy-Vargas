# Frontend SPA – Restaurant Reservations

Aplicación React (Vite) que consume los microservicios de Users, Reservations y Search para gestionar el flujo completo de reservas del restaurante.

## Características
- Login/registro con JWT e interceptores de Axios.
- Home con buscador + filtros conectados a Search API (Solr).
- Vista de detalles con confirmación de reservas (llama al cálculo concurrente).
- Sección “Mis reservas” protegida para cada usuario.
- Panel Admin con edición y eliminación (solo rol `admin`).
- React Router v6, React Query, React Hook Form, TailwindCSS + lucide-react.
- Notifications con `react-hot-toast` y estados de carga/errores cubiertos.

## Requisitos
- Node.js 20+
- Microservicios corriendo localmente:
  - Users API `http://localhost:8080`
  - Reservations API `http://localhost:8081`
  - Search API `http://localhost:8082`

## Configuración
1. Copiar variables:
   ```bash
   cp .env.example .env
   ```
2. Ajustar las URLs si corrés los servicios en otros puertos u hostnames.

## Scripts
```bash
npm install          # dependencias
npm run dev          # entorno local (http://localhost:5173)
npm run build        # genera dist/
npm run preview      # sirve build local
npm run lint         # reglas ESLint recomendadas por Vite
```

## Estructura principal
```
src/
 ├─ api/              # axios instances + clients (auth, reservations, search)
 ├─ components/       # UI reutilizable (auth, search, admin, etc.)
 ├─ context/          # AuthProvider (maneja sesión + localStorage)
 ├─ hooks/            # React Query hooks para reservas
 ├─ pages/            # Login, Register, Home, MyReservations, Admin
 ├─ routes/           # AppRoutes + guards (Private/Admin)
 └─ utils/            # constantes + formateadores
```

## Docker
Build multi-stage + Nginx listo para servir la SPA y hacer proxy a los microservicios:
```bash
docker build -t reservations-frontend .
docker run --rm -p 3000:80 --env-file .env reservations-frontend
```
> Ajustá `nginx.conf` si tus servicios viven con otros nombres (`users-api`, `reservations-api`, `search-api`) cuando orquestes con Docker Compose.

## Endpoints Consumidos
- `POST /api/users/login`, `POST /api/users/register`
- `GET/PUT/DELETE /api/reservations`, `POST /api/reservations/:id/confirm`, `GET /api/reservations/user/:user_id`
- `GET /api/search`, `GET /api/search/:id`

## Próximos pasos sugeridos
- Añadir pruebas con Vitest + React Testing Library para formularios críticos.
- Implementar dark mode (Tailwind) y animaciones con Framer Motion si se requiere.
- Integrar el frontend en el `docker-compose` global del proyecto.
