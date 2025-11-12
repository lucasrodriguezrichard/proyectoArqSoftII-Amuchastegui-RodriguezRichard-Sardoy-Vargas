# Issue #2: Implementar Reservations API (Entidad Principal)

## Descripción
Crear el microservicio de reservas (entidad principal) con MongoDB, implementando operaciones CRUD, concurrencia con Go Routines, y comunicación con RabbitMQ.

## Objetivo
Este es el microservicio core del sistema que maneja las reservas de restaurante, valida usuarios, implementa cálculo concurrente, y notifica cambios a través de RabbitMQ.

## Tareas

### 1. Setup Inicial del Proyecto
- [ ] Crear carpeta `reservations-api/`
- [ ] Inicializar módulo Go: `go mod init github.com/tu-usuario/restaurant/reservations-api`
- [ ] Estructura de carpetas siguiendo patrón MVC

### 2. Configuración
- [ ] Crear `internal/config/config.go`
  - Variables para MongoDB
  - Variables para RabbitMQ
  - Variables para Users API URL
  - Puerto del servidor

- [ ] Crear `.env.example`
  - MongoDB connection string
  - RabbitMQ connection
  - Users API base URL

### 3. Modelo de Dominio
- [ ] Crear `internal/domain/reservation.go`
  - Estructura Reservation con tags BSON y JSON
  - Métodos de validación
  - Estados: pending, confirmed, cancelled, completed

- [ ] Crear `internal/domain/calculation.go`
  - Estructura para resultado de cálculo concurrente
  - Función de cálculo paralelo con goroutines
  - WaitGroup y Channel para sincronización

### 4. Repository Layer
- [ ] Crear `internal/repository/reservation_repository.go`
  - Interface ReservationRepository
  - Implementación con MongoDB
  - CRUD completo: Create, GetByID, GetAll, Update, Delete, GetByUserID

- [ ] Configurar conexión MongoDB
  - Crear `internal/db/mongodb.go`
  - Connection pooling
  - Ping test

### 5. Service Layer
- [ ] Crear `internal/service/reservation_service.go`
  - Interface ReservationService
  - Implementación con lógica de negocio
  - Validación de owner contra Users API (HTTP call)
  - Cálculo concurrente (goroutines + channels + waitgroup)
  - Publicación en RabbitMQ después de Create/Update/Delete

- [ ] Crear `internal/service/user_client.go`
  - Cliente HTTP para Users API
  - Función ValidateUser(userID) -> error
  - Manejo de errores de red

- [ ] Crear `internal/service/rabbitmq_publisher.go`
  - Cliente RabbitMQ publisher
  - Función Publish(operation, entityID)
  - Manejo de errores de conexión

### 6. Controller Layer
- [ ] Crear `internal/controller/reservation_controller.go`
  - Interface ReservationController
  - Handlers para todos los endpoints
  - Validación de request bodies
  - Códigos HTTP correctos

### 7. Endpoints HTTP
- [ ] **POST /api/reservations** - Crear reserva
  - Validar owner_id contra Users API
  - Ejecutar cálculo concurrente de disponibilidad/precio
  - Guardar en MongoDB
  - Publicar mensaje a RabbitMQ
  - Retornar 201 con reserva creada

- [ ] **GET /api/reservations/:id** - Obtener reserva
  - Validar existencia
  - Retornar 200 con datos

- [ ] **GET /api/reservations** - Listar todas (paginado opcional)
  - Retornar array de reservas
  - Considerar filtros básicos

- [ ] **GET /api/reservations/user/:user_id** - Reservas por usuario
  - Filtrar por owner_id
  - Retornar array

- [ ] **PUT /api/reservations/:id** - Actualizar reserva
  - Validar ownership (owner o admin)
  - Actualizar en MongoDB
  - Publicar a RabbitMQ
  - Retornar 200

- [ ] **DELETE /api/reservations/:id** - Eliminar reserva
  - Validar ownership (owner o admin)
  - Eliminar de MongoDB
  - Publicar a RabbitMQ
  - Retornar 204

- [ ] **POST /api/reservations/:id/confirm** - Confirmar reserva (acción)
  - Cambiar estado a "confirmed"
  - Ejecutar cálculo concurrente (ej: confirmar mesas, calcular descuentos)
  - Actualizar en MongoDB
  - Publicar a RabbitMQ
  - Retornar 200

### 8. Cálculo Concurrente
- [ ] Implementar función de cálculo paralelo
  - Dividir trabajo en múltiples goroutines
  - Ejemplo: calcular disponibilidad de mesas, precios, descuentos
  - Usar channels para comunicar resultados
  - Usar WaitGroup para sincronización
  - Agregar timeout de seguridad

Ejemplo:
```go
func (s *reservationService) CalculateAvailabilityAndPrice(req ReservationRequest) (*CalculationResult, error) {
    results := make(chan PartialResult, 3)
    var wg sync.WaitGroup

    // Goroutine 1: Check table availability
    wg.Add(1)
    go func() {
        defer wg.Done()
        availability := checkTableAvailability(req.TableID, req.DateTime)
        results <- PartialResult{Type: "availability", Data: availability}
    }()

    // Goroutine 2: Calculate base price
    wg.Add(1)
    go func() {
        defer wg.Done()
        price := calculateBasePrice(req.Guests, req.MealType)
        results <- PartialResult{Type: "price", Data: price}
    }()

    // Goroutine 3: Apply discounts
    wg.Add(1)
    go func() {
        defer wg.Done()
        discount := calculateDiscount(req.DateTime, req.UserID)
        results <- PartialResult{Type: "discount", Data: discount}
    }()

    // Close channel when all goroutines finish
    go func() {
        wg.Wait()
        close(results)
    }()

    // Collect results
    finalResult := &CalculationResult{}
    for partial := range results {
        // Merge results
    }

    return finalResult, nil
}
```

### 9. RabbitMQ Integration
- [ ] Configurar conexión RabbitMQ
- [ ] Crear exchange y queue
- [ ] Implementar publisher
- [ ] Formato de mensajes:
```json
{
  "operation": "create|update|delete",
  "entity_id": "uuid",
  "entity_type": "reservation",
  "timestamp": "2025-11-12T10:00:00Z"
}
```

### 10. Router y Middleware
- [ ] Crear `internal/transport/http/router.go`
  - Configurar Gin
  - CORS middleware
  - Logger middleware
  - Health check endpoint: GET /health

### 11. Main Entry Point
- [ ] Crear `cmd/api/main.go`
  - Cargar configuración
  - Conectar MongoDB
  - Conectar RabbitMQ
  - Wire dependencies
  - Iniciar servidor HTTP

### 12. Docker
- [ ] Crear `Dockerfile`
  - Multi-stage build
  - Go 1.21+ base image
  - Exponer puerto 8081

- [ ] Crear scripts
  - `scripts/start.sh`
  - `scripts/start.bat`
  - `scripts/dev.sh` con Air para hot-reload

### 13. Tests
- [ ] Crear `internal/service/reservation_service_test.go`
  - Test de creación de reserva
  - Test de validación de usuario
  - Test de cálculo concurrente
  - Test de publicación RabbitMQ
  - Mocks de repository, user client, rabbitmq publisher

### 14. Documentación
- [ ] Crear `README.md`
  - Descripción del servicio
  - Endpoints disponibles
  - Ejemplos de requests/responses
  - Variables de entorno
  - Cómo ejecutar

## Estructura de Archivos
```
reservations-api/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── domain/
│   │   ├── reservation.go
│   │   └── calculation.go
│   ├── repository/
│   │   └── reservation_repository.go
│   ├── service/
│   │   ├── reservation_service.go
│   │   ├── reservation_service_test.go
│   │   ├── user_client.go
│   │   └── rabbitmq_publisher.go
│   ├── controller/
│   │   └── reservation_controller.go
│   ├── transport/
│   │   └── http/
│   │       └── router.go
│   └── db/
│       └── mongodb.go
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

### Reservation
```go
type Reservation struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    OwnerID     string            `bson:"owner_id" json:"owner_id"`
    TableNumber int               `bson:"table_number" json:"table_number"`
    Guests      int               `bson:"guests" json:"guests"`
    DateTime    time.Time         `bson:"date_time" json:"date_time"`
    MealType    string            `bson:"meal_type" json:"meal_type"` // breakfast, lunch, dinner, event
    Status      string            `bson:"status" json:"status"` // pending, confirmed, cancelled, completed
    TotalPrice  float64           `bson:"total_price" json:"total_price"`
    SpecialRequests string        `bson:"special_requests,omitempty" json:"special_requests,omitempty"`
    CreatedAt   time.Time         `bson:"created_at" json:"created_at"`
    UpdatedAt   time.Time         `bson:"updated_at" json:"updated_at"`
}
```

## Endpoints API

### POST /api/reservations
```json
// Request
{
  "owner_id": "user-uuid",
  "table_number": 5,
  "guests": 4,
  "date_time": "2025-11-20T20:00:00Z",
  "meal_type": "dinner",
  "special_requests": "Window seat please"
}

// Response (201)
{
  "id": "mongodb-object-id",
  "owner_id": "user-uuid",
  "table_number": 5,
  "guests": 4,
  "date_time": "2025-11-20T20:00:00Z",
  "meal_type": "dinner",
  "status": "pending",
  "total_price": 150.00,
  "special_requests": "Window seat please",
  "created_at": "2025-11-12T10:00:00Z",
  "updated_at": "2025-11-12T10:00:00Z"
}
```

### POST /api/reservations/:id/confirm
```json
// Request
{
  "confirmation_notes": "Reservation confirmed by manager"
}

// Response (200)
{
  "id": "mongodb-object-id",
  "status": "confirmed",
  "total_price": 135.00,  // precio recalculado con descuento
  "updated_at": "2025-11-12T10:05:00Z"
}
```

## Variables de Entorno (.env.example)
```env
# MongoDB
MONGO_URI=mongodb://localhost:27017
MONGO_DB=reservations_db
MONGO_COLLECTION=reservations

# RabbitMQ
RABBITMQ_URI=amqp://guest:guest@localhost:5672/
RABBITMQ_EXCHANGE=restaurant_events
RABBITMQ_QUEUE=reservations_updates

# Users API
USERS_API_URL=http://localhost:8080

# Server
PORT=8081
APP_ENV=development
```

## Dependencias Go
```bash
go get go.mongodb.org/mongo-driver/mongo
go get github.com/gin-gonic/gin
go get github.com/rabbitmq/amqp091-go
go get github.com/joho/godotenv
```

## Criterios de Aceptación
- [ ] Todos los endpoints CRUD funcionan
- [ ] Validación de usuario contra Users API funciona
- [ ] Cálculo concurrente implementado con goroutines, channels y waitgroup
- [ ] RabbitMQ publica mensajes correctamente
- [ ] Tests pasan exitosamente
- [ ] MongoDB almacena datos correctamente
- [ ] Códigos HTTP correctos en todas las respuestas
- [ ] Manejo de errores en todas las capas
- [ ] Documentación completa

## Prioridad
🔴 **CRÍTICA** - Entidad principal del sistema

## Estimación
⏱️ 16-20 horas

## Dependencias
- Issue #1 (Users API) debe estar completo para validación de usuarios

## Notas
- El cálculo concurrente puede simular: verificar disponibilidad de mesa, calcular precio base, aplicar descuentos, verificar restricciones
- Los mensajes de RabbitMQ serán consumidos por search-api para sincronizar índices
- Considerar implementar soft delete en lugar de delete físico
- El endpoint de confirmación es la "acción" principal que requiere cálculo concurrente
