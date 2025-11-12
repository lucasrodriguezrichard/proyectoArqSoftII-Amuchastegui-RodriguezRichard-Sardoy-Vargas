# Issue #4: Implementar Frontend con React

## Descripción
Desarrollar la aplicación web frontend en React que consume los microservicios del backend, implementando todas las pantallas requeridas: Login, Registro, Búsqueda, Detalles, Mis Reservas y Admin.

## Objetivo
Crear una SPA (Single Page Application) completa con autenticación JWT, búsqueda de reservas, gestión de reservas propias y panel de administración para usuarios admin.

## Tareas

### 1. Setup Inicial del Proyecto
- [ ] Crear proyecto React con Vite
```bash
npm create vite@latest frontend -- --template react
cd frontend
npm install
```

- [ ] Instalar dependencias necesarias
```bash
npm install react-router-dom axios
npm install @tanstack/react-query
npm install react-hook-form
npm install tailwindcss postcss autoprefixer
npm install lucide-react  # iconos
```

- [ ] Configurar Tailwind CSS
```bash
npx tailwindcss init -p
```

### 2. Estructura del Proyecto
- [ ] Crear estructura de carpetas
```
frontend/
├── src/
│   ├── api/
│   │   ├── axios.js
│   │   ├── auth.js
│   │   ├── reservations.js
│   │   └── search.js
│   ├── components/
│   │   ├── common/
│   │   │   ├── Navbar.jsx
│   │   │   ├── Footer.jsx
│   │   │   ├── Loader.jsx
│   │   │   └── ErrorMessage.jsx
│   │   ├── auth/
│   │   │   ├── LoginForm.jsx
│   │   │   └── RegisterForm.jsx
│   │   ├── search/
│   │   │   ├── SearchBar.jsx
│   │   │   ├── FilterPanel.jsx
│   │   │   ├── ReservationCard.jsx
│   │   │   └── Pagination.jsx
│   │   ├── reservation/
│   │   │   ├── ReservationDetails.jsx
│   │   │   └── ConfirmModal.jsx
│   │   └── admin/
│   │       ├── ReservationTable.jsx
│   │       └── EditModal.jsx
│   ├── pages/
│   │   ├── Login.jsx
│   │   ├── Register.jsx
│   │   ├── Home.jsx
│   │   ├── ReservationDetails.jsx
│   │   ├── MyReservations.jsx
│   │   └── Admin.jsx
│   ├── hooks/
│   │   ├── useAuth.js
│   │   └── useReservations.js
│   ├── context/
│   │   └── AuthContext.jsx
│   ├── utils/
│   │   ├── constants.js
│   │   └── formatters.js
│   ├── routes/
│   │   ├── AppRoutes.jsx
│   │   ├── PrivateRoute.jsx
│   │   └── AdminRoute.jsx
│   ├── App.jsx
│   └── main.jsx
├── public/
├── .env.example
├── Dockerfile
├── nginx.conf
└── package.json
```

### 3. Configuración de Axios
- [ ] Crear `src/api/axios.js`
  - Instancia de axios con baseURL
  - Interceptor para agregar JWT token
  - Interceptor para manejar errores 401 (logout)
  - Timeout configuration

```javascript
import axios from 'axios';

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080',
  timeout: 10000,
});

// Request interceptor
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export default api;
```

### 4. Context de Autenticación
- [ ] Crear `src/context/AuthContext.jsx`
  - Estado global de autenticación
  - Funciones: login, logout, register
  - Persistir token y usuario en localStorage
  - Hook useAuth para consumir el contexto

### 5. API Clients
- [ ] Crear `src/api/auth.js`
  - `login(username, password)` → POST /api/users/login
  - `register(userData)` → POST /api/users/register
  - `getUserById(id)` → GET /api/users/:id

- [ ] Crear `src/api/search.js`
  - `searchReservations(params)` → GET /api/search
  - `getReservationById(id)` → GET /api/search/:id

- [ ] Crear `src/api/reservations.js`
  - `createReservation(data)` → POST /api/reservations
  - `confirmReservation(id)` → POST /api/reservations/:id/confirm
  - `getMyReservations(userId)` → GET /api/reservations/user/:user_id
  - `getAllReservations()` → GET /api/reservations (admin)
  - `updateReservation(id, data)` → PUT /api/reservations/:id (admin)
  - `deleteReservation(id)` → DELETE /api/reservations/:id (admin)

### 6. Páginas

#### 6.1 Login Page
- [ ] Crear `src/pages/Login.jsx`
  - Formulario con campos: username/email y password
  - Validación con react-hook-form
  - Llamar a API de login
  - Guardar token y usuario
  - Redireccionar a /home
  - Link a página de registro
  - Mostrar errores de autenticación

#### 6.2 Register Page
- [ ] Crear `src/pages/Register.jsx`
  - Formulario con campos: username, email, first_name, last_name, password, confirm_password
  - Validación de campos
  - Validación de password match
  - Llamar a API de registro
  - Auto-login después de registro exitoso
  - Redireccionar a /home
  - Link a página de login

#### 6.3 Home/Search Page
- [ ] Crear `src/pages/Home.jsx`
  - Barra de búsqueda
  - Panel de filtros:
    - meal_type (breakfast, lunch, dinner, event)
    - date_from, date_to
    - guests_min, guests_max
    - status
  - Lista de resultados con ReservationCard
  - Paginación
  - Empty state para "no results"
  - Loader mientras carga
  - Botón "Ver Detalles" en cada card

- [ ] Crear `src/components/search/SearchBar.jsx`
  - Input de búsqueda
  - Botón de buscar
  - Clear filters

- [ ] Crear `src/components/search/FilterPanel.jsx`
  - Filtros colapsables
  - Inputs de fecha
  - Select de meal_type
  - Range de guests

- [ ] Crear `src/components/search/ReservationCard.jsx`
  - Vista compacta de reserva
  - Información: fecha, hora, meal_type, guests, table, precio
  - Badge de estado (pending/confirmed/cancelled/completed)
  - Botón "Ver Detalles"

- [ ] Crear `src/components/search/Pagination.jsx`
  - Botones prev/next
  - Números de página
  - Info de "Mostrando X de Y resultados"

#### 6.4 Reservation Details Page
- [ ] Crear `src/pages/ReservationDetails.jsx`
  - Vista completa de una reserva
  - Toda la información detallada
  - Botón "Confirmar Reserva" si status === "pending"
  - Modal de confirmación
  - Mostrar precio final después de confirmación
  - Success message después de acción

- [ ] Crear `src/components/reservation/ConfirmModal.jsx`
  - Modal de confirmación
  - Campo opcional para notas
  - Botón "Confirmar"
  - Botón "Cancelar"

#### 6.5 My Reservations Page
- [ ] Crear `src/pages/MyReservations.jsx`
  - Lista de reservas del usuario logueado
  - Filtros por estado
  - Tabs: Todas | Pendientes | Confirmadas | Pasadas
  - Cada reserva clickeable para ver detalles
  - Empty state: "No tienes reservas aún"

#### 6.6 Admin Page
- [ ] Crear `src/pages/Admin.jsx`
  - Protegida con AdminRoute
  - Tabla de todas las reservas
  - Columnas: ID, Usuario, Mesa, Guests, Fecha, Tipo, Estado, Precio, Acciones
  - Botón "Editar" por cada fila
  - Botón "Eliminar" por cada fila
  - Confirmación antes de eliminar
  - Paginación
  - Filtros por estado

- [ ] Crear `src/components/admin/ReservationTable.jsx`
  - Tabla responsive
  - Sorting por columnas
  - Acciones en línea

- [ ] Crear `src/components/admin/EditModal.jsx`
  - Formulario de edición
  - Todos los campos editables
  - Validación
  - Guardar cambios

### 7. Componentes Comunes
- [ ] Crear `src/components/common/Navbar.jsx`
  - Logo del restaurante
  - Links: Home | Mis Reservas | Admin (si es admin)
  - Usuario logueado con dropdown
  - Botón Logout

- [ ] Crear `src/components/common/Loader.jsx`
  - Spinner de carga
  - Reutilizable

- [ ] Crear `src/components/common/ErrorMessage.jsx`
  - Mensaje de error con estilo
  - Reutilizable

### 8. Rutas y Protección
- [ ] Crear `src/routes/AppRoutes.jsx`
  - Configurar React Router
  - Rutas públicas: /login, /register
  - Rutas privadas: /home, /reservations/:id, /my-reservations
  - Rutas admin: /admin

- [ ] Crear `src/routes/PrivateRoute.jsx`
  - HOC para proteger rutas
  - Verificar token
  - Redireccionar a /login si no autenticado

- [ ] Crear `src/routes/AdminRoute.jsx`
  - HOC para rutas de admin
  - Verificar role === "admin"
  - Redireccionar a /home si no es admin

### 9. Hooks Personalizados
- [ ] Crear `src/hooks/useAuth.js`
  - Hook para consumir AuthContext
  - Exponer: user, token, login, logout, register, isAuthenticated, isAdmin

- [ ] Crear `src/hooks/useReservations.js`
  - Hook con React Query
  - Queries: searchReservations, getReservationById, getMyReservations
  - Mutations: createReservation, confirmReservation, updateReservation, deleteReservation

### 10. Estilos con Tailwind
- [ ] Configurar `tailwind.config.js`
  - Colores personalizados del restaurante
  - Fonts
  - Breakpoints

- [ ] Crear tema consistente
  - Botones primarios y secundarios
  - Cards con sombras
  - Formularios con validación visual
  - Estados hover y focus

### 11. Docker
- [ ] Crear `Dockerfile`
  - Multi-stage build
  - Build stage con Node
  - Production stage con nginx
  - Copiar build a nginx

```dockerfile
# Build stage
FROM node:18-alpine AS build
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# Production stage
FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

- [ ] Crear `nginx.conf`
  - Configurar SPA routing
  - Proxy para APIs

```nginx
server {
    listen 80;
    server_name localhost;
    root /usr/share/nginx/html;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api/ {
        proxy_pass http://localhost:8080/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### 12. Variables de Entorno
- [ ] Crear `.env.example`
```env
VITE_API_URL=http://localhost:8080
VITE_SEARCH_API_URL=http://localhost:8082
VITE_APP_NAME=Restaurant Reservations
```

### 13. Testing (Opcional)
- [ ] Instalar Vitest y React Testing Library
- [ ] Tests de componentes críticos
- [ ] Tests de hooks

### 14. Documentación
- [ ] Crear `README.md`
  - Descripción del frontend
  - Cómo ejecutar en desarrollo
  - Cómo buildear
  - Estructura del proyecto
  - Screenshots

## Flujos Principales

### Flujo de Login → Búsqueda → Detalle → Confirmación
```
1. Usuario abre /login
2. Ingresa credenciales
3. Frontend llama POST /api/users/login
4. Recibe token y user
5. Guarda en localStorage
6. Redirecciona a /home
7. /home carga búsqueda vacía (GET /api/search)
8. Usuario aplica filtros
9. Frontend llama GET /api/search?filters...
10. Muestra resultados
11. Usuario click "Ver Detalles"
12. Navega a /reservations/:id
13. Frontend llama GET /api/search/:id
14. Muestra detalle completo
15. Usuario click "Confirmar Reserva"
16. Frontend llama POST /api/reservations/:id/confirm
17. Muestra mensaje de éxito
```

## Criterios de Aceptación
- [ ] Todas las páginas requeridas están implementadas
- [ ] Login y registro funcionan correctamente
- [ ] JWT se envía en todas las requests autenticadas
- [ ] Búsqueda con filtros funciona
- [ ] Paginación funciona
- [ ] Crear y confirmar reserva funciona
- [ ] Mis Reservas muestra reservas del usuario
- [ ] Panel admin solo visible para admins
- [ ] Admin puede editar y eliminar reservas
- [ ] UI responsive (móvil, tablet, desktop)
- [ ] Manejo de errores en todas las acciones
- [ ] Loading states en todas las requests
- [ ] Logout funciona correctamente
- [ ] Rutas protegidas funcionan
- [ ] Docker build exitoso
- [ ] Documentación completa

## Prioridad
🟠 **ALTA** - Interfaz de usuario del sistema

## Estimación
⏱️ 24-30 horas

## Dependencias
- Issues #1, #2, #3 deben estar completas
- APIs del backend funcionando

## Tecnologías
- React 18
- Vite
- React Router v6
- Axios
- React Query (TanStack Query)
- React Hook Form
- Tailwind CSS
- Lucide React (iconos)
- Docker + Nginx

## Notas
- Considerar implementar Dark Mode
- Agregar animaciones con Framer Motion (opcional)
- Implementar toast notifications con react-hot-toast
- Considerar lazy loading de componentes con React.lazy
- Implementar debounce en búsqueda en tiempo real
- Validar formularios en cliente antes de enviar
- Manejar estados de carga, error y éxito en todas las mutaciones
