# Issue #1: Completar y Mejorar Users API

## Descripción
Completar la implementación de users-api con todas las funcionalidades requeridas para el trabajo final, siguiendo el patrón MVC y utilizando MySQL con GORM.

## Estado Actual
- ✅ Conexión MySQL con GORM
- ✅ Modelo User con GORM
- ✅ Repositorio implementado
- ✅ JWT implementado
- ✅ Hashing de passwords con bcrypt

## Tareas Pendientes

### 1. Completar Endpoints Faltantes
- [ ] **POST /api/users/register** - Crear usuario normal
  - Validar datos de entrada (username, email, password, first_name, last_name)
  - Hash de password
  - Role por defecto: "normal"
  - Retornar usuario creado (sin password)

- [ ] **GET /api/users/:id** - Obtener usuario por ID
  - Validar existencia del usuario
  - Retornar datos sin password
  - Este endpoint será usado por otros microservicios para validación

- [ ] **POST /api/users/login** - Login de usuario
  - Validar credenciales (username/email + password)
  - Generar JWT token con claims (user_id, role)
  - Retornar: `{ "token": "...", "user": {...} }`

### 2. Agregar Soporte para Usuarios Admin
- [ ] Agregar campo `role` al modelo User
  - Valores: "normal", "admin"
  - Validación en creación y actualización

- [ ] **POST /api/users/admin** - Crear usuario admin (solo para admin)
  - Middleware de autenticación JWT
  - Validar que el usuario que crea sea admin

### 3. Middleware de Autenticación
- [ ] Crear middleware `AuthMiddleware`
  - Validar JWT token en header `Authorization: Bearer <token>`
  - Extraer user_id y role del token
  - Guardar en contexto de Gin

- [ ] Crear middleware `AdminMiddleware`
  - Verificar que role == "admin"
  - Retornar 403 si no es admin

### 4. Validación y Manejo de Errores
- [ ] Implementar validaciones en Service layer
  - Email válido y único
  - Username único
  - Password mínimo 8 caracteres

- [ ] Códigos de estado HTTP correctos
  - 200: Success
  - 201: Created
  - 400: Bad Request (validación)
  - 401: Unauthorized
  - 403: Forbidden
  - 404: Not Found
  - 500: Internal Server Error

### 5. Tests
- [ ] Crear `user_service_test.go`
  - Test de creación de usuario
  - Test de login exitoso
  - Test de login con credenciales inválidas
  - Test de obtener usuario por ID
  - Mock del repository

### 6. Configuración Docker
- [ ] Verificar Dockerfile
- [ ] Agregar scripts de inicio
  - `start.sh` (Linux/Mac)
  - `start.bat` (Windows)
  - `dev.sh` con hot-reload usando Air

### 7. Documentación
- [ ] Actualizar README con endpoints disponibles
- [ ] Ejemplos de requests/responses en formato JSON
- [ ] Variables de entorno necesarias

## Estructura de Archivos
```
users-api/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── domain/
│   │   └── user.go
│   ├── repository/
│   │   └── user_repository.go
│   ├── service/
│   │   ├── user_service.go
│   │   └── user_service_test.go
│   ├── transport/
│   │   └── http/
│   │       ├── router.go
│   │       ├── users_handler.go
│   │       └── middleware.go
│   ├── auth/
│   │   └── jwt.go
│   ├── crypto/
│   │   └── password.go
│   └── db/
│       └── mysql.go
├── scripts/
│   ├── start.sh
│   ├── start.bat
│   └── dev.sh
├── Dockerfile
├── .env.example
├── go.mod
└── README.md
```

## Modelo de Datos

### User
```go
type User struct {
    ID        string    `gorm:"primaryKey" json:"id"`
    Username  string    `gorm:"unique;not null" json:"username"`
    Email     string    `gorm:"unique;not null" json:"email"`
    Password  string    `gorm:"not null" json:"-"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    Role      string    `gorm:"default:'normal'" json:"role"` // normal, admin
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

## Endpoints API

### POST /api/users/register
```json
// Request
{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "securepass123",
  "first_name": "John",
  "last_name": "Doe"
}

// Response (201)
{
  "id": "uuid",
  "username": "johndoe",
  "email": "john@example.com",
  "first_name": "John",
  "last_name": "Doe",
  "role": "normal",
  "created_at": "2025-11-12T10:00:00Z"
}
```

### POST /api/users/login
```json
// Request
{
  "username": "johndoe",  // o "email": "john@example.com"
  "password": "securepass123"
}

// Response (200)
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "uuid",
    "username": "johndoe",
    "email": "john@example.com",
    "role": "normal"
  }
}
```

### GET /api/users/:id
```json
// Response (200)
{
  "id": "uuid",
  "username": "johndoe",
  "email": "john@example.com",
  "first_name": "John",
  "last_name": "Doe",
  "role": "normal",
  "created_at": "2025-11-12T10:00:00Z"
}
```

## Variables de Entorno (.env.example)
```env
# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=restaurant_user
DB_PASS=restaurant_pass
DB_NAME=users_db

# JWT
JWT_SECRET=your-super-secret-key-change-this
JWT_ACCESS_TTL=24h
JWT_REFRESH_TTL=168h

# Server
PORT=8080
APP_ENV=development
```

## Tecnologías
- Go 1.21+
- Gin Web Framework
- GORM
- MySQL 8.0
- JWT (golang-jwt/jwt)
- bcrypt
- Docker

## Criterios de Aceptación
- [ ] Todos los endpoints funcionan correctamente
- [ ] Tests pasan exitosamente
- [ ] Validaciones implementadas en todas las capas
- [ ] JWT funciona correctamente
- [ ] Passwords hasheados con bcrypt
- [ ] Códigos HTTP correctos
- [ ] Documentación completa
- [ ] Docker funciona correctamente
- [ ] Hot-reload funciona en modo dev

## Prioridad
🔴 **ALTA** - Base fundamental para los otros microservicios

## Estimación
⏱️ 8-12 horas

## Notas
- Este microservicio será consultado por `reservations-api` y `search-api` para validar usuarios
- El endpoint GET /api/users/:id debe ser público (sin auth) para permitir validaciones entre servicios
- Considerar rate limiting para el endpoint de login
