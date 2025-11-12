# Users API - Sistema de Gestión de Restaurante

API de usuarios con autenticación JWT para el sistema de gestión de restaurante. Maneja registro, login, y gestión de usuarios (normales y administradores).

## 🛠️ Tecnologías

- **Go 1.21+**
- **Framework**: Chi Router
- **Base de Datos**: MySQL 8.0
- **ORM**: GORM
- **Autenticación**: JWT (golang-jwt/jwt)
- **Hash de Passwords**: bcrypt
- **Testing**: Go testing + Mocks
- **Containerización**: Docker

## 📦 Características

- ✅ Registro de usuarios normales
- ✅ Login con email o username
- ✅ Autenticación JWT (access + refresh tokens)
- ✅ Hashing de passwords con bcrypt
- ✅ Creación de usuarios admin (protegido)
- ✅ Consulta de usuarios por ID (para otros microservicios)
- ✅ Middleware de autenticación
- ✅ Middleware de verificación de admin
- ✅ Tests unitarios con mocks
- ✅ Health check endpoint
- ✅ CORS configurado

## 🚀 Inicio Rápido

### Opción 1: Con Docker

```bash
# Build y run
docker build -t users-api .
docker run -p 8080:8080 \
  -e DB_HOST=host.docker.internal \
  -e DB_PORT=3306 \
  -e DB_USER=restaurant_user \
  -e DB_PASS=restaurant_pass \
  -e DB_NAME=users_db \
  -e JWT_SECRET=your-secret-key \
  users-api
```

### Opción 2: Local

```bash
# 1. Copiar variables de entorno
cp .env.example .env
# Editar .env con tus configuraciones

# 2. Instalar dependencias
go mod download

# 3. Ejecutar en modo desarrollo (hot-reload)
chmod +x scripts/dev.sh
./scripts/dev.sh

# O ejecutar en modo producción
chmod +x scripts/start.sh
./scripts/start.sh
```

### Windows

```cmd
# Copiar variables de entorno
copy .env.example .env

# Instalar dependencias
go mod download

# Ejecutar
scripts\start.bat
```

## 📚 Endpoints API

### Base URL
```
http://localhost:8080
```

### Públicos (sin autenticación)

#### Health Check
```http
GET /health
```

**Response:**
```json
{
  "status": "healthy"
}
```

#### Registro de Usuario
```http
POST /api/users/register
Content-Type: application/json

{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe"
}
```

**Response (201):**
```json
{
  "id": 1,
  "username": "johndoe",
  "email": "john@example.com",
  "first_name": "John",
  "last_name": "Doe",
  "role": "user",
  "created_at": "2025-11-12T10:00:00Z",
  "updated_at": "2025-11-12T10:00:00Z"
}
```

#### Login
```http
POST /api/users/login
Content-Type: application/json

{
  "identifier": "john@example.com",
  "password": "password123"
}
```

**Response (200):**
```json
{
  "user": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "role": "user",
    "created_at": "2025-11-12T10:00:00Z",
    "updated_at": "2025-11-12T10:00:00Z"
  },
  "tokens": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2025-11-13T10:00:00Z"
  }
}
```

#### Obtener Usuario por ID
```http
GET /api/users/{id}
```

**Response (200):**
```json
{
  "id": 1,
  "username": "johndoe",
  "email": "john@example.com",
  "first_name": "John",
  "last_name": "Doe",
  "role": "user",
  "created_at": "2025-11-12T10:00:00Z",
  "updated_at": "2025-11-12T10:00:00Z"
}
```

### Protegidos (requieren autenticación y rol admin)

#### Crear Usuario Admin
```http
POST /api/admin/users
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "username": "admin",
  "email": "admin@example.com",
  "password": "adminpass123",
  "first_name": "Admin",
  "last_name": "User"
}
```

**Response (201):**
```json
{
  "id": 2,
  "username": "admin",
  "email": "admin@example.com",
  "first_name": "Admin",
  "last_name": "User",
  "role": "admin",
  "created_at": "2025-11-12T10:00:00Z",
  "updated_at": "2025-11-12T10:00:00Z"
}
```

## 🔐 Autenticación JWT

Para endpoints protegidos, incluir el token JWT en el header:

```http
Authorization: Bearer <your_jwt_token>
```

Los tokens contienen los siguientes claims:
- `sub`: User ID
- `username`: Username
- `role`: user | admin
- `exp`: Expiration timestamp
- `iat`: Issued at timestamp

## ⚠️ Códigos de Estado HTTP

| Código | Significado |
|--------|-------------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request (validación) |
| 401 | Unauthorized (sin token o token inválido) |
| 403 | Forbidden (sin permisos de admin) |
| 404 | Not Found |
| 409 | Conflict (usuario ya existe) |
| 422 | Unprocessable Entity (datos inválidos) |
| 500 | Internal Server Error |

## 🔧 Variables de Entorno

Ver [`.env.example`](.env.example) para todas las variables disponibles.

Variables críticas:
- `DB_HOST`: Host de MySQL (default: localhost)
- `DB_PORT`: Puerto de MySQL (default: 3306)
- `DB_USER`: Usuario de MySQL
- `DB_PASS`: Contraseña de MySQL
- `DB_NAME`: Nombre de la base de datos
- `JWT_SECRET`: Secret key para JWT (min 32 caracteres en producción)
- `JWT_ACCESS_TTL`: Tiempo de vida del access token (default: 24h)
- `JWT_REFRESH_TTL`: Tiempo de vida del refresh token (default: 168h)

## 🧪 Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test ./... -cover

# Run verbose
go test ./... -v

# Run only service tests
go test ./internal/service/... -v

# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

Resultado esperado:
```
ok  	github.com/blassardoy/restaurant-reservas/users-api/internal/service	0.848s
```

## 📁 Estructura del Proyecto

```
users-api/
├── cmd/
│   └── server/
│       └── main.go           # Entry point
├── internal/
│   ├── auth/
│   │   └── jwt.go            # JWT issuer
│   ├── config/
│   │   └── config.go         # Configuration
│   ├── crypto/
│   │   └── password.go       # Password hasher (bcrypt)
│   ├── db/
│   │   └── mysql.go          # MySQL connection
│   ├── domain/
│   │   └── user.go           # Domain entities, DTOs, interfaces
│   ├── repository/
│   │   └── user_repository.go # Database operations
│   ├── service/
│   │   ├── user_service.go   # Business logic
│   │   └── user_service_test.go # Unit tests
│   └── transport/
│       └── http/
│           ├── router.go     # HTTP router
│           ├── middleware.go # Auth middlewares
│           └── users_endpoints.go # Legacy endpoints
├── scripts/
│   ├── start.sh              # Production start script (Linux/Mac)
│   ├── start.bat             # Production start script (Windows)
│   └── dev.sh                # Development with hot-reload (Air)
├── Dockerfile                # Docker image
├── .env.example              # Environment variables template
├── go.mod                    # Go dependencies
└── README.md
```

## 🏗️ Arquitectura

El proyecto sigue **Arquitectura Hexagonal (Clean Architecture)**:

```
┌─────────────────────────────────────────┐
│         HTTP Transport Layer            │
│  (Router, Handlers, Middleware)         │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│         Service Layer                   │
│  (Business Logic, Orchestration)        │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│         Repository Layer                │
│  (Database Access with GORM)            │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│             MySQL                       │
└─────────────────────────────────────────┘
```

**Principios aplicados:**
- Dependency Inversion (interfaces en domain)
- Separation of Concerns
- Clean Architecture
- Repository Pattern
- Dependency Injection

## 🔗 Integración con Otros Microservicios

Este microservicio es consultado por:
- **reservations-api**: Para validar usuarios antes de crear reservas
- **search-api**: Para obtener información de usuarios

Endpoint de validación:
```http
GET /api/users/{id}
```

## 📝 Notas de Desarrollo

### Hot Reload con Air

El script `dev.sh` usa [Air](https://github.com/cosmtrek/air) para hot-reload automático. Se instala automáticamente la primera vez.

### Base de Datos

El servicio auto-migra la tabla `users` al iniciar. Schema:

```sql
CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(32) UNIQUE NOT NULL,
    email VARCHAR(191) UNIQUE NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'user' NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### Seguridad

- ✅ Passwords hasheados con bcrypt (cost 10)
- ✅ JWT con HMAC-SHA256
- ✅ Passwords nunca expuestos en responses
- ✅ Validación de inputs
- ✅ CORS configurado
- ⚠️ Cambiar `JWT_SECRET` en producción

## 🐛 Troubleshooting

### Error: "connection refused"
- Verificar que MySQL esté corriendo
- Verificar host y puerto en `.env`

### Error: "invalid JWT token"
- Verificar que el token no haya expirado
- Verificar que `JWT_SECRET` sea el mismo usado para generar el token

### Tests fallan
```bash
# Limpiar cache y re-run
go clean -testcache
go test ./...
```

## 👥 Equipo

- Amuchastegui, Matias
- Rodriguez Richard, Lucas
- Sardoy, Blas
- Vargas, Santino

## 📄 Licencia

Proyecto educativo - Arquitectura de Software II - UCC 2025
