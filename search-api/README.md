# Search API

Microservicio de búsqueda sobre reservas con Solr y caché en capas.

## Ejecutar

1) Levantar Solr y Memcached

```
cd search-api
docker-compose up -d
```

2) Variables de entorno (opcional)

```
SOLR_URL=http://localhost:8983/solr
SOLR_CORE=reservations
RABBITMQ_URI=amqp://admin:admin@localhost:5672/
RABBITMQ_EXCHANGE=restaurant_events
RABBITMQ_QUEUE=reservations_updates
RESERVATIONS_API_URL=http://localhost:8081
PORT=8082
```

3) Ejecutar API

```
go run ./cmd/api/main.go
```

4) Endpoints

- GET /health
- GET /api/search
- GET /api/search/:id
- GET /api/search/stats
- POST /api/search/reindex

Nota: la implementación de Solr, Memcached y consumer está preparada para ser completada; este es un esqueleto funcional para continuar Issue #3.

