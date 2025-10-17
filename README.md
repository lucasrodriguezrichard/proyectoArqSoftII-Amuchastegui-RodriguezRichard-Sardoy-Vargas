# 🍽️ Sistema de Gestión de Restaurant

Sistema completo de gestión de restaurante desarrollado con Go, siguiendo arquitectura hexagonal y buenas prácticas de desarrollo.

## 📋 Descripción del Proyecto

Sistema integral que permite:
- ✅ **Reservas de mesas** para diferentes tipos de comida (desayuno, almuerzo, merienda, cena, eventos privados)
- 📝 **Gestión de pedidos** con carga directa en el sistema
- 🧾 **Generación automática de tickets** al finalizar la comida/evento
- ⭐ **Sistema de reseñas** automático para calificar comida, servicio y ambiente

## 🏗️ Arquitectura

El proyecto sigue **Arquitectura Hexagonal (Clean Architecture)** con las siguientes capas:

```
restaurant-system/
├── cmd/
│   └── api/                    # Punto de entrada
│       └── main.go
├── internal/
│   ├── config/                 # Configuración
│   ├── domain/                 # Entidades de dominio
│   ├── dao/                    # Data Access Objects
│   ├── repository/             # Repositorios (interfaces)
│   ├── services/               # Lógica de negocio
│   ├── controllers/            # Controladores HTTP
│   └── middleware/             # Middleware (CORS, Auth, Logger)
├── docker-compose.yml          # Docker Compose
├── Dockerfile                  # Imagen Docker
├── Makefile                    # Comandos útiles
└── README.md
```

## 🛠️ Tecnologías

- **Lenguaje:** Go 1.21
- **Framework Web:** Gin
- **Base de Datos:** PostgreSQL 15
- **Autenticación:** JWT
- **Containerización:** Docker & Docker Compose
- **ORM:** database/sql con lib/pq

## 📦 Entidades del Sistema

### 1. **Restaurant**
Información del restaurante (nombre, dirección, contacto)

### 2. **Table**
Mesas del restaurante con capacidad y disponibilidad

### 3. **Reservation**
Reservas de mesas con estados:
- `pending` - Pendiente
- `confirmed` - Confirmada
- `cancelled` - Cancelada
- `completed` - Completada

### 4. **Order**
Pedidos de comida con estados:
- `pending`, `confirmed`, `preparing`, `ready`, `served`, `completed`, `cancelled`

### 5. **Ticket**
Tickets de pago generados automáticamente

### 6. **Review**
Reseñas con calificaciones detalladas (comida, servicio, ambiente)

### 7. **MenuItem**
Items del menú con categorías

### 8. **User**
Usuarios del sistema (admin, manager, waiter, chef, customer)

## 🚀 Instalación y Uso

### Prerrequisitos
- Docker & Docker Compose
- Go 1.21+ (para desarrollo local)
- Make (opcional, pero recomendado)

### Opción 1: Con Docker (Recomendado)

1. **Clonar el repositorio**
```bash
git clone <repository-url>
cd restaurant-system
```

2. **Iniciar los servicios**
```bash
make docker-up
# o
docker-compose up -d
```

3. **Verificar que todo funciona**
```bash
make health
# o
curl http://localhost:8080/health
```

La API estará disponible en: `http://localhost:8080`
Adminer (DB Manager) en: `http://localhost:8081`

### Opción 2: Desarrollo Local

1. **Instalar dependencias**
```bash
make install-deps
# o
go mod download
```

2. **Configurar variables de entorno**
```bash
cp .env.example .env
# Editar .env con tus configuraciones
```

3. **Iniciar PostgreSQL** (con Docker)
```bash
docker-compose up -d db
```

4. **Ejecutar la aplicación**
```bash
make run
# o
go run ./cmd/api/main.go
```

## 📚 API Endpoints

### Reservas
- `POST /api/v1/reservations` - Crear reserva
- `GET /api/v1/reservations/:id` - Obtener reserva
- `GET /api/v1/reservations` - Listar reservas
- `PUT /api/v1/reservations/:id/confirm` - Confirmar
- `PUT /api/v1/reservations/:id/cancel` - Cancelar
- `PUT /api/v1/reservations/:id/complete` - Completar

### Pedidos
- `POST /api/v1/orders` - Crear pedido
- `GET /api/v1/orders/:id` - Obtener pedido
- `GET /api/v1/orders/reservation/:id` - Pedidos por reserva
- `PUT /api/v1/orders/:id/status` - Actualizar estado

### Tickets
- `POST /api/v1/tickets` - Generar ticket
- `GET /api/v1/tickets/:id` - Obtener ticket
- `GET /api/v1/tickets/order/:id` - Ticket por pedido
- `GET /api/v1/tickets/reports/sales` - Reporte de ventas

### Reseñas
- `POST /api/v1/reviews` - Crear reseña
- `GET /api/v1/reviews/:id` - Obtener reseña
- `GET /api/v1/reviews` - Listar reseñas
- `GET /api/v1/reviews/stats/average` - Promedios
- `PUT /api/v1/reviews/:id` - Actualizar
- `DELETE /api/v1/reviews/:id` - Eliminar

### Admin (Requiere JWT)
- `GET /api/v1/admin/reservations` - Todas las reservas
- `GET /api/v1/admin/reports/sales` - Reporte de ventas
- `GET /api/v1/admin/reports/reviews` - Estadísticas de reseñas

Ver [API_DOCUMENTATION.md](API_DOCUMENTATION.md) para más detalles.

## 🔧 Comandos Útiles (Makefile)

```bash
make help              # Ver todos los comandos disponibles
make docker-up         # Iniciar servicios con Docker
make docker-down       # Detener servicios
make docker-rebuild    # Reconstruir imágenes
make logs              # Ver logs de la API
make logs-db           # Ver logs de la base de datos
make test              # Ejecutar tests
make test-coverage     # Tests con cobertura
make db-shell          # Abrir shell de PostgreSQL
make db-backup         # Backup de la base de datos
make health            # Verificar salud de la API
```

## 🗄️ Base de Datos

### Conexión
```
Host: localhost
Port: 5432
User: restaurant_user
Password: restaurant_pass
Database: restaurant_db
```

### Administrador Web (Adminer)
Accede a `http://localhost:8081` para gestionar la base de datos visualmente.

### Datos de Prueba
El sistema incluye datos de prueba:
- 1 Restaurant
- 10 Mesas
- 17 Items de menú
- 1 Usuario admin (username: `admin`, password: `admin123`)

## 🧪 Testing

```bash
# Ejecutar todos los tests
make test

# Tests con cobertura
make test-coverage

# Ver reporte de cobertura
open coverage.html
```

## 📝 Flujo del Sistema

### 1️⃣ Reserva de Mesa
```
Cliente solicita reserva → Sistema valida disponibilidad → 
Reserva creada (pending) → Admin confirma → Reserva confirmada
```

### 2️⃣ Pedido de Comida
```
Cliente hace pedido → Sistema calcula totales (subtotal + impuesto) → 
Pedido asociado a reserva → Estado: pending → preparing → ready → served
```

### 3️⃣ Generación de Ticket
```
Pedido completado → Sistema genera ticket automáticamente → 
Incluye todos los items + totales → Registra método de pago
```

### 4️⃣ Sistema de Reseñas
```
Experiencia completada → Cliente recibe invitación → 
Califica comida, servicio y ambiente → Reseña almacenada
```

## 🔐 Autenticación

Los endpoints de administración requieren JWT Token:

```http
Authorization: Bearer <your_jwt_token>
```

## 🌐 Variables de Entorno

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=restaurant_user
DB_PASSWORD=restaurant_pass
DB_NAME=restaurant_db
DB_SSLMODE=disable

# Server
PORT=8080

# JWT
JWT_SECRET=your-super-secret-key

# Environment
APP_ENV=development
```

## 📊 Estructura de Respuestas

### Success (200/201)
```json
{
  "id": "uuid",
  "data": { ... }
}
```

### Error (4xx/5xx)
```json
{
  "error": "Mensaje de error descriptivo"
}
```

## 🤝 Contribución

Este proyecto es parte del curso de **Arquitectura de Software II - UCC 2025**.

### Equipo
- Amuchastegui, Matias
- Rodriguez Richard, Lucas
- Sardoy, Blas
- Vargas, Santino

## 📄 Licencia

Este proyecto es educativo y forma parte del curso de Arquitectura de Software II en la Universidad Católica de Córdoba (UCC).

## 🐛 Troubleshooting

### Puerto 8080 ocupado
```bash
# Cambiar el puerto en docker-compose.yml
ports:
  - "8081:8080"  # Usar puerto 8081 en vez de 8080
```

### Error de conexión a la base de datos
```bash
# Verificar que PostgreSQL esté corriendo
docker-compose ps

# Ver logs de la base de datos
make logs-db

# Reiniciar servicios
make docker-rebuild
```

### Limpiar y reiniciar
```bash
# Detener servicios
make docker-down

# Limpiar volúmenes
docker-compose down -v

# Reiniciar
make docker-up
```

## 📞 Soporte

Para preguntas o problemas, contactar al equipo de desarrollo o revisar la documentación en el repositorio.

---

**Desarrollado con ❤️ por el equipo de Arquitectura de Software II - UCC 2025**