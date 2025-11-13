# 🍽️ Sistema de Gestión de Reservas de Restaurante

Sistema de gestión de reservas desarrollado con arquitectura de microservicios en Go, siguiendo patrones de Clean Architecture, event-driven architecture y caching strategies.

## 📋 Descripción del Proyecto

Sistema integral que permite:

- **Gestión de usuarios** con autenticación JWT (access + refresh tokens)
- **Reservas de mesas** para diferentes tipos de comida (desayuno, almuerzo, merienda, cena, eventos privados)
- **Búsqueda avanzada** de reservas con Apache Solr y caché en capas
- **Comunicación asíncrona** entre microservicios con RabbitMQ
- **Frontend SPA** con React para gestión visual

### Servicios

#### 1. **Users API** (Puerto 8080)

- **Base de datos**: MySQL 8.0
- **Framework**: Chi Router + GORM
- **Funcionalidad**:
  - Registro y autenticación de usuarios
  - Emisión de JWT tokens (access + refresh)
  - Gestión de roles (user, admin)
  - Endpoint de validación para otros servicios

#### 2. **Reservations API** (Puerto 8081)

- **Base de datos**: MongoDB 7.0
- **Message Broker**: RabbitMQ (Publisher)
- **Framework**: Gin
- **Funcionalidad**:
  - CRUD de reservas
  - Validación de usuarios (integración con Users API)
  - Cálculo de precios con descuentos
  - Publicación de eventos (create, update, delete, confirm) a RabbitMQ

#### 3. **Search API** (Puerto 8082)

- **Motor de búsqueda**: Apache Solr 9
- **Caché**: Memcached + CCache (dual-layer)
- **Message Broker**: RabbitMQ (Consumer)
- **Framework**: Gin
- **Funcionalidad**:
  - Búsqueda full-text de reservas
  - Indexación automática desde eventos RabbitMQ
  - Caché dual (local + distribuida)
  - Re-indexación manual
  - Introspección de caché

#### 4. **Frontend** (Puerto 3000)

- **Framework**: React + Vite
- **Funcionalidad**:
  - SPA para gestión visual de reservas
  - Integración con las 3 APIs

## 🛠️ Tecnologías

### Backend

- **Lenguaje**: Go 1.21+
- **Frameworks**: Gin, Chi Router
- **Bases de Datos**: MySQL 8.0, MongoDB 7.0
- **ORMs**: GORM (MySQL), MongoDB Driver oficial
- **Autenticación**: JWT (golang-jwt/jwt)
- **Message Broker**: RabbitMQ 3.13
- **Motor de Búsqueda**: Apache Solr 9
- **Caché**: Memcached 1.6 + CCache (in-memory)
- **Containerización**: Docker & Docker Compose

### Frontend

- **Framework**: React 18
- **Build Tool**: Vite
- **Servidor**: Nginx (en contenedor)

### Infraestructura

- **Orchestration**: Docker Compose
- **Database Admin**: Adminer (puerto 18080)
- **RabbitMQ Management**: Puerto 15672

## Instalación y Uso

### Prerrequisitos

- Docker & Docker Compose
- Go 1.21+ (opcional, para desarrollo local)
- Node.js 18+ (opcional, para desarrollo frontend)

## 📚 API Endpoints

### Users API (http://localhost:8080)

#### Públicos

- `POST /api/users/register` - Registrar usuario
- `POST /api/users/login` - Login (retorna JWT)
- `GET /api/users/{id}` - Obtener usuario por ID

#### Protegidos (requieren JWT admin)

- `POST /api/admin/users` - Crear usuario admin

### Reservations API (http://localhost:8081)

- `POST /api/reservations` - Crear reserva
- `GET /api/reservations` - Listar todas las reservas
- `GET /api/reservations/:id` - Obtener reserva por ID
- `GET /api/reservations/user/:user_id` - Reservas de un usuario
- `PUT /api/reservations/:id` - Actualizar reserva
- `DELETE /api/reservations/:id` - Eliminar reserva
- `POST /api/reservations/:id/confirm` - Confirmar reserva

### Search API (http://localhost:8082)

#### Búsqueda

- `GET /api/search?q=...&page=1&size=10` - Búsqueda full-text
  - Parámetros: `q`, `page`, `size`, `sort`, `order`, `meal_type`, `status`, `guests`
- `GET /api/search/:id` - Buscar por ID
- `GET /api/search/stats` - Estadísticas (documentos + caché)
- `POST /api/search/reindex` - Re-indexar Solr desde Reservations API

#### Introspección de Caché

- `GET /__cache/stats` - Estadísticas de caché (local entries, distributed hits/misses)
- `GET /__cache/get?key=<key>` - Obtener valor de caché por key
- `POST /__cache/invalidate` - Invalidar toda la caché

## 🔄 Flujo Event-Driven

```
1. Cliente crea/actualiza/elimina reserva en Reservations API
                    ↓
2. Reservations API publica evento a RabbitMQ
   Exchange: restaurant_events
   Routing Key: reservation.create | reservation.update | reservation.delete
                    ↓
3. Search API consume evento desde RabbitMQ
                    ↓
4. Search API sincroniza con Solr (index/update/delete)
                    ↓
5. Search API invalida caché afectada
```

### Eventos RabbitMQ

**Estructura de mensaje:**

```json
{
  "operation": "create|update|delete|confirm",
  "entity_id": "reservation_id",
  "entity_type": "reservation",
  "timestamp": "2025-11-12T10:00:00Z"
}
```

## 🗄️ Bases de Datos

### MySQL (Users)

```
Host: localhost
Port: 3307
User: root
Password: bla2ucc
Database: users
```

**Tabla `users`:**

- id, username, email, first_name, last_name, password_hash, role, created_at, updated_at

### MongoDB (Reservations)

```
Host: localhost
Port: 27017
Database: reservations_db
Collection: reservations
```

**Documento `reservation`:**

- \_id, owner_id, table_number, guests, meal_type, date_time, total_price, status, created_at, updated_at

### Solr (Search Index)

```
URL: http://localhost:8983/solr
Core: reservations
```

**Acceso a UI Solr:**
http://localhost:8983/solr/#/reservations

### Adminer (Web DB Manager)

```
URL: http://localhost:18080
Server: users-db
User: root
Password: bla2ucc
```

## 🎯 Caché Strategy (Search API)

**Dual-Layer Cache:**

1. **Local Cache (CCache)** - In-process memory cache

   - TTL: 300 segundos (configurable)
   - Primera capa, más rápida

2. **Distributed Cache (Memcached)** - Shared cache
   - TTL: 900 segundos (configurable)
   - Segunda capa, compartida entre instancias
   - Host: localhost:11211

**Cache Keys:**

- Búsquedas: `SHA1(q=...&page=...&size=...)`
- Documentos: `doc:<reservation_id>`

**Invalidación:**

- Automática en eventos RabbitMQ (create, update, delete, confirm)
- Manual vía `POST /__cache/invalidate`

## 🔐 Autenticación

Los endpoints de Reservations API validan usuarios contra Users API.

**JWT Token (Users API):**

```http
POST /api/users/login
{
  "identifier": "user@example.com",
  "password": "password123"
}
```

**Response:**

```json
{
  "user": {...},
  "tokens": {
    "access_token": "eyJhbGci...",
    "refresh_token": "eyJhbGci...",
    "expires_at": "2025-11-13T10:00:00Z"
  }
}
```

**Uso (ejemplo):**

```http
Authorization: Bearer <access_token>
```

## 🌐 Variables de Entorno

### Users API

```env
DB_HOST=users-db
DB_PORT=3306
DB_USER=root
DB_PASS=bla2ucc
DB_NAME=users
JWT_SECRET=supersecreto-docker-key-min-32-chars
JWT_ACCESS_TTL=24h
JWT_REFRESH_TTL=168h
PORT=8080
```

### Reservations API

```env
PORT=8081
MONGO_URI=mongodb://reservations-mongodb:27017
MONGO_DB=reservations_db
MONGO_COLLECTION=reservations
RABBITMQ_URI=amqp://admin:admin@reservations-rabbitmq:5672/
RABBITMQ_EXCHANGE=restaurant_events
RABBITMQ_QUEUE=reservations_updates
USERS_API_URL=http://users-api:8080
```

### Search API

```env
PORT=8082
SOLR_URL=http://search-solr:8983/solr
SOLR_CORE=reservations
CACHE_TTL_SECONDS=300
DIST_CACHE_TTL_SECONDS=900
MEMCACHED_ADDRS=search-memcached:11211
RABBITMQ_URI=amqp://admin:admin@reservations-rabbitmq:5672/
RABBITMQ_QUEUE=reservations_updates
RABBITMQ_EXCHANGE=restaurant_events
RESERVATIONS_API_URL=http://reservations-api:8081
```

## 📝 Flujo de Uso Típico

### 1️⃣ Registro de Usuario

```bash
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

### 2️⃣ Login

```bash
curl -X POST http://localhost:8080/api/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "john@example.com",
    "password": "password123"
  }'
```

### 3️⃣ Crear Reserva

```bash
curl -X POST http://localhost:8081/api/reservations \
  -H "Content-Type: application/json" \
  -d '{
    "owner_id": "1",
    "table_number": 5,
    "guests": 4,
    "meal_type": "cena",
    "date_time": "2025-11-15T20:00:00Z"
  }'
```

### 4️⃣ Buscar Reservas

```bash
# Búsqueda simple
curl "http://localhost:8082/api/search?q=*:*&page=1&size=10"

# Filtrar por meal_type
curl "http://localhost:8082/api/search?meal_type=cena&page=1&size=10"

# Ver estadísticas
curl http://localhost:8082/api/search/stats
```

### 5️⃣ Inspeccionar Caché

```bash
# Ver estadísticas de caché
curl http://localhost:8082/__cache/stats

# Obtener valor específico
curl "http://localhost:8082/__cache/get?key=doc:reservation_id"

# Invalidar caché
curl -X POST http://localhost:8082/__cache/invalidate
```

## 🧪 Testing

Cada servicio tiene sus propios tests:

```bash
# Users API
cd users-api && go test ./...

# Reservations API
cd reservations-api && go test ./...

# Search API
cd search-api && go test ./...
```

## 🐛 Troubleshooting

### Servicios no inician

```bash
# Ver logs de todos los servicios
docker-compose logs -f

# Ver logs de un servicio específico
docker-compose logs -f users-api
docker-compose logs -f reservations-api
docker-compose logs -f search-api
```

### Reiniciar servicios

```bash
# Detener todo
docker-compose down

# Limpiar volúmenes (¡cuidado! borra datos)
docker-compose down -v

# Reiniciar
docker-compose up --build
```

### Verificar conectividad entre servicios

```bash
# Entrar a un contenedor
docker exec -it users-api sh

# Hacer ping a otro servicio
curl http://reservations-api:8081/health
```

### RabbitMQ no consume mensajes

1. Acceder a RabbitMQ Management: http://localhost:15672
2. Credenciales: `admin` / `admin`
3. Verificar queues en pestaña "Queues"
4. Verificar bindings en pestaña "Exchanges"

### Solr no indexa

1. Acceder a Solr UI: http://localhost:8983/solr
2. Verificar core "reservations" existe
3. Ejecutar reindex manual: `curl -X POST http://localhost:8082/api/search/reindex`

## Monitoreo

- **RabbitMQ Management UI**: http://localhost:15672 (admin/admin)
- **Solr Admin UI**: http://localhost:8983/solr
- **Adminer (DB)**: http://localhost:18080
- **Frontend**: http://localhost:3000

## Grupo

- Amuchastegui, Matias
- Rodriguez Richard, Lucas
- Sardoy, Blas
- Vargas, Santino
