# Issue #3: Implementar Search API con SolR, Caché y RabbitMQ

## Descripción
Crear el microservicio de búsqueda que indexa reservas en SolR, implementa doble capa de caché (CCache local + Memcached distribuida), y consume mensajes de RabbitMQ para sincronizar índices.

## Objetivo
Proporcionar búsqueda rápida, filtrada y paginada sobre reservas, con sistema de caché para optimizar queries frecuentes y sincronización automática mediante eventos de RabbitMQ.

## Tareas

### 1. Setup Inicial del Proyecto
- [ ] Crear carpeta `search-api/`
- [ ] Inicializar módulo Go: `go mod init github.com/tu-usuario/restaurant/search-api`
- [ ] Estructura de carpetas siguiendo patrón MVC

### 2. Configuración
- [ ] Crear `internal/config/config.go`
  - Variables para SolR
  - Variables para Memcached
  - Variables para RabbitMQ
  - Variables para Reservations API URL
  - Configuración de caché (TTL)

- [ ] Crear `.env.example`
  - SolR connection
  - Memcached servers
  - RabbitMQ connection
  - Reservations API base URL

### 3. Modelo de Dominio
- [ ] Crear `internal/domain/reservation_document.go`
  - Estructura para documento SolR
  - Tags para SolR schema
  - Métodos de conversión desde Reservation

### 4. SolR Integration
- [ ] Crear `internal/solr/client.go`
  - Cliente HTTP para SolR
  - Funciones: Index, Update, Delete, Search
  - Manejo de errores de conexión

- [ ] Crear `internal/solr/schema.go`
  - Definición de schema SolR
  - Campos indexados y stored
  - Configuración de filtros y sorting

- [ ] Crear script de inicialización SolR
  - `scripts/init-solr.sh`
  - Crear core "reservations"
  - Configurar schema

### 5. Cache Layer
- [ ] Crear `internal/cache/local_cache.go`
  - Implementación con CCache
  - TTL configurable
  - Funciones: Get, Set, Delete, Clear
  - Estadísticas de hit/miss

- [ ] Crear `internal/cache/distributed_cache.go`
  - Implementación con Memcached
  - Connection pooling
  - Funciones: Get, Set, Delete
  - Fallback si Memcached falla

- [ ] Crear `internal/cache/dual_cache.go`
  - Wrapper que combina local + distributed
  - Estrategia: Check local → Check distributed → Query SolR
  - Populate both caches on miss
  - Invalidar ambas en updates

### 6. Repository Layer
- [ ] Crear `internal/repository/search_repository.go`
  - Interface SearchRepository
  - Implementación con SolR client
  - Métodos: Search, GetByID, Index, Update, Delete

### 7. Service Layer
- [ ] Crear `internal/service/search_service.go`
  - Interface SearchService
  - Implementación con lógica de búsqueda
  - Integración de dual cache
  - Paginación
  - Filtros dinámicos
  - Sorting por múltiples campos

- [ ] Crear `internal/service/reservation_client.go`
  - Cliente HTTP para Reservations API
  - Función GetReservationByID(id) -> Reservation
  - Usado para sincronización completa desde RabbitMQ

- [ ] Crear `internal/service/sync_service.go`
  - Maneja sincronización de índices
  - Recibe eventos de RabbitMQ
  - Llama a Reservations API para datos completos
  - Indexa/actualiza/elimina en SolR
  - Invalida caché

### 8. RabbitMQ Consumer
- [ ] Crear `internal/rabbitmq/consumer.go`
  - Configurar consumidor de queue
  - Auto-reconnect en caso de fallo
  - Procesar mensajes:
    - create → Index en SolR
    - update → Update en SolR
    - delete → Delete en SolR
  - ACK después de procesar exitosamente
  - NACK y requeue en caso de error

- [ ] Crear `internal/rabbitmq/message_handler.go`
  - Parser de mensajes
  - Validación de formato
  - Routing a sync_service

### 9. Controller Layer
- [ ] Crear `internal/controller/search_controller.go`
  - Handlers para endpoints de búsqueda
  - Validación de query params
  - Paginación
  - Códigos HTTP correctos

### 10. Endpoints HTTP
- [ ] **GET /api/search** - Búsqueda paginada
  - Query params:
    - `q`: query string (opcional, default: *)
    - `page`: número de página (default: 1)
    - `size`: items por página (default: 10, max: 100)
    - `sort`: campo de ordenamiento (default: created_at)
    - `order`: asc|desc (default: desc)
    - Filtros: `meal_type`, `status`, `date_from`, `date_to`, `guests_min`, `guests_max`
  - Response:
```json
{
  "results": [...],
  "total": 150,
  "page": 1,
  "size": 10,
  "pages": 15
}
```

- [ ] **GET /api/search/:id** - Obtener reserva por ID
  - Usar caché
  - Fallback a SolR si no está en caché
  - Retornar 200 con documento

- [ ] **GET /api/search/stats** - Estadísticas de búsqueda
  - Total de documentos indexados
  - Cache hit rate
  - Queries más frecuentes

- [ ] **POST /api/search/reindex** - Reindexar todo (admin)
  - Obtener todas las reservas de Reservations API
  - Indexar en SolR
  - Limpiar caché
  - Retornar progreso

### 11. Cache Strategy
Implementar estrategia de caché en capas:

```go
func (s *searchService) Search(query SearchQuery) (*SearchResult, error) {
    // 1. Generate cache key
    cacheKey := generateCacheKey(query)

    // 2. Check local cache (CCache)
    if result, found := s.localCache.Get(cacheKey); found {
        return result, nil
    }

    // 3. Check distributed cache (Memcached)
    if result, err := s.distributedCache.Get(cacheKey); err == nil {
        // Populate local cache
        s.localCache.Set(cacheKey, result, 5*time.Minute)
        return result, nil
    }

    // 4. Query SolR
    result, err := s.repository.Search(query)
    if err != nil {
        return nil, err
    }

    // 5. Populate both caches
    s.localCache.Set(cacheKey, result, 5*time.Minute)
    s.distributedCache.Set(cacheKey, result, 15*time.Minute)

    return result, nil
}
```

### 12. Router y Middleware
- [ ] Crear `internal/transport/http/router.go`
  - Configurar Gin
  - CORS middleware
  - Logger middleware
  - Health check: GET /health
  - Metrics endpoint: GET /metrics (cache stats)

### 13. Main Entry Point
- [ ] Crear `cmd/api/main.go`
  - Cargar configuración
  - Conectar SolR
  - Conectar Memcached
  - Inicializar local cache
  - Iniciar consumer RabbitMQ en goroutine
  - Wire dependencies
  - Iniciar servidor HTTP

### 14. Docker
- [ ] Crear `Dockerfile`
  - Multi-stage build
  - Go 1.21+ base image
  - Exponer puerto 8082

- [ ] Crear scripts
  - `scripts/start.sh`
  - `scripts/start.bat`
  - `scripts/dev.sh` con Air

### 15. Tests
- [ ] Crear `internal/service/search_service_test.go`
  - Test de búsqueda con caché
  - Test de cache hit/miss
  - Test de paginación
  - Test de filtros
  - Mocks de repository y caches

### 16. Documentación
- [ ] Crear `README.md`
  - Descripción del servicio
  - Endpoints disponibles
  - Ejemplos de búsqueda con filtros
  - Arquitectura de caché
  - Variables de entorno
  - Schema de SolR

## Estructura de Archivos
```
search-api/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── domain/
│   │   └── reservation_document.go
│   ├── solr/
│   │   ├── client.go
│   │   └── schema.go
│   ├── cache/
│   │   ├── local_cache.go
│   │   ├── distributed_cache.go
│   │   └── dual_cache.go
│   ├── repository/
│   │   └── search_repository.go
│   ├── service/
│   │   ├── search_service.go
│   │   ├── search_service_test.go
│   │   ├── sync_service.go
│   │   └── reservation_client.go
│   ├── rabbitmq/
│   │   ├── consumer.go
│   │   └── message_handler.go
│   ├── controller/
│   │   └── search_controller.go
│   └── transport/
│       └── http/
│           └── router.go
├── scripts/
│   ├── start.sh
│   ├── start.bat
│   ├── dev.sh
│   └── init-solr.sh
├── Dockerfile
├── .env.example
├── go.mod
└── README.md
```

## Schema SolR (managed-schema.xml)
```xml
<schema name="reservations" version="1.6">
  <field name="id" type="string" indexed="true" stored="true" required="true" multiValued="false" />
  <field name="owner_id" type="string" indexed="true" stored="true" />
  <field name="table_number" type="int" indexed="true" stored="true" />
  <field name="guests" type="int" indexed="true" stored="true" />
  <field name="date_time" type="pdate" indexed="true" stored="true" />
  <field name="meal_type" type="string" indexed="true" stored="true" />
  <field name="status" type="string" indexed="true" stored="true" />
  <field name="total_price" type="pfloat" indexed="true" stored="true" />
  <field name="special_requests" type="text_general" indexed="true" stored="true" />
  <field name="created_at" type="pdate" indexed="true" stored="true" />
  <field name="updated_at" type="pdate" indexed="true" stored="true" />

  <uniqueKey>id</uniqueKey>
</schema>
```

## Endpoints API

### GET /api/search
```
GET /api/search?q=dinner&meal_type=dinner&date_from=2025-11-15&guests_min=2&page=1&size=10&sort=date_time&order=asc
```

Response:
```json
{
  "results": [
    {
      "id": "mongodb-id",
      "owner_id": "user-uuid",
      "table_number": 5,
      "guests": 4,
      "date_time": "2025-11-20T20:00:00Z",
      "meal_type": "dinner",
      "status": "confirmed",
      "total_price": 150.00,
      "special_requests": "Window seat",
      "created_at": "2025-11-12T10:00:00Z"
    }
  ],
  "total": 45,
  "page": 1,
  "size": 10,
  "pages": 5,
  "cache_hit": false
}
```

### GET /api/search/stats
```json
{
  "total_documents": 1250,
  "cache_stats": {
    "local_cache": {
      "hit_rate": 0.85,
      "total_hits": 8500,
      "total_misses": 1500
    },
    "distributed_cache": {
      "hit_rate": 0.65,
      "total_hits": 975,
      "total_misses": 525
    }
  }
}
```

## Variables de Entorno (.env.example)
```env
# SolR
SOLR_URL=http://localhost:8983/solr
SOLR_CORE=reservations

# Memcached
MEMCACHED_SERVERS=localhost:11211
MEMCACHED_TIMEOUT=1s

# Local Cache (CCache)
LOCAL_CACHE_MAX_SIZE=10000
LOCAL_CACHE_TTL=5m

# Distributed Cache
DISTRIBUTED_CACHE_TTL=15m

# RabbitMQ
RABBITMQ_URI=amqp://guest:guest@localhost:5672/
RABBITMQ_QUEUE=reservations_updates
RABBITMQ_EXCHANGE=restaurant_events

# Reservations API
RESERVATIONS_API_URL=http://localhost:8081

# Server
PORT=8082
APP_ENV=development
```

## Dependencias Go
```bash
go get github.com/gin-gonic/gin
go get github.com/vanng822/go-solr/solr
go get github.com/bradfitz/gomemcache/memcache
go get github.com/karlseguin/ccache/v3
go get github.com/rabbitmq/amqp091-go
go get github.com/joho/godotenv
```

## Flujo de Sincronización

```
1. Reservations API crea/actualiza/elimina reserva
   ↓
2. Publica mensaje a RabbitMQ
   ↓
3. Search API consume mensaje
   ↓
4. Search API obtiene datos completos de Reservations API
   ↓
5. Search API indexa/actualiza/elimina en SolR
   ↓
6. Search API invalida caché (local + distributed)
```

## Criterios de Aceptación
- [ ] SolR indexa documentos correctamente
- [ ] Búsqueda con filtros y paginación funciona
- [ ] Doble capa de caché funciona (local + distributed)
- [ ] Consumer de RabbitMQ sincroniza cambios automáticamente
- [ ] Cache invalidation funciona correctamente
- [ ] Sorting por múltiples campos funciona
- [ ] Manejo de errores en todas las capas
- [ ] Tests pasan exitosamente
- [ ] Documentación completa

## Prioridad
🟠 **ALTA** - Componente crítico para búsqueda

## Estimación
⏱️ 20-24 horas

## Dependencias
- Issue #2 (Reservations API) debe estar completo
- RabbitMQ debe estar configurado

## Notas
- El caché local (CCache) es más rápido pero limitado en memoria
- El caché distribuido (Memcached) es compartido entre instancias
- Considerar implementar circuit breaker para llamadas a Reservations API
- Los queries más frecuentes deben tener mayor hit rate en caché
- Implementar logging detallado para debugging de sincronización
- La estrategia de caché debe considerar el TTL según la naturaleza de los datos (reservas pasadas vs futuras)
