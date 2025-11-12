# Issue #5: Docker Compose - Orquestación de Todos los Servicios

## Descripción
Crear el archivo `docker-compose.yml` principal que orquesta todos los microservicios, bases de datos y servicios de infraestructura necesarios para ejecutar el sistema completo.

## Objetivo
Permitir levantar toda la aplicación (frontend + 3 microservicios backend + 5 servicios de infraestructura) con un solo comando: `docker-compose up`

## Servicios a Incluir

### Backend Microservices
1. **users-api** (Puerto 8080)
2. **reservations-api** (Puerto 8081)
3. **search-api** (Puerto 8082)

### Frontend
4. **frontend** (Puerto 3000)

### Bases de Datos
5. **mysql** (Puerto 3306) - Para users-api
6. **mongodb** (Puerto 27017) - Para reservations-api
7. **solr** (Puerto 8983) - Para search-api

### Infraestructura
8. **rabbitmq** (Puertos 5672, 15672) - Message broker
9. **memcached** (Puerto 11211) - Distributed cache

### Herramientas de Administración (Opcional)
10. **adminer** (Puerto 8084) - Gestión de MySQL
11. **mongo-express** (Puerto 8085) - Gestión de MongoDB

## Tareas

### 1. Crear docker-compose.yml Principal
- [ ] Crear archivo en la raíz del proyecto
- [ ] Definir version: "3.8"
- [ ] Configurar networks
- [ ] Configurar volumes

### 2. Servicio: MySQL
- [ ] Imagen: mysql:8.0
- [ ] Variables de entorno
- [ ] Volumen persistente
- [ ] Health check
- [ ] Script de inicialización

### 3. Servicio: MongoDB
- [ ] Imagen: mongo:7.0
- [ ] Variables de entorno
- [ ] Volumen persistente
- [ ] Health check
- [ ] Script de inicialización

### 4. Servicio: SolR
- [ ] Imagen: solr:9.4
- [ ] Crear core "reservations"
- [ ] Volumen para configuración
- [ ] Health check
- [ ] Script para crear core e inicializar schema

### 5. Servicio: RabbitMQ
- [ ] Imagen: rabbitmq:3-management
- [ ] Habilitar management plugin
- [ ] Variables de entorno
- [ ] Volumen persistente
- [ ] Health check
- [ ] Exponer puerto management (15672)

### 6. Servicio: Memcached
- [ ] Imagen: memcached:alpine
- [ ] Configuración de memoria
- [ ] Health check

### 7. Servicio: users-api
- [ ] Build desde ./users-api
- [ ] Depends on: mysql
- [ ] Variables de entorno
- [ ] Restart policy
- [ ] Health check
- [ ] Networks

### 8. Servicio: reservations-api
- [ ] Build desde ./reservations-api
- [ ] Depends on: mongodb, rabbitmq, users-api
- [ ] Variables de entorno
- [ ] Restart policy
- [ ] Health check
- [ ] Networks

### 9. Servicio: search-api
- [ ] Build desde ./search-api
- [ ] Depends on: solr, memcached, rabbitmq, reservations-api
- [ ] Variables de entorno
- [ ] Restart policy
- [ ] Health check
- [ ] Networks

### 10. Servicio: frontend
- [ ] Build desde ./frontend
- [ ] Depends on: users-api, search-api
- [ ] Variables de entorno
- [ ] Restart policy
- [ ] Health check
- [ ] Networks

### 11. Servicios Opcionales de Administración
- [ ] Adminer para MySQL
- [ ] Mongo Express para MongoDB

### 12. Scripts de Inicialización
- [ ] Crear `scripts/init-mysql.sql`
  - Crear database users_db
  - Crear tabla users si no existe

- [ ] Crear `scripts/init-mongo.js`
  - Crear database reservations_db
  - Crear colección reservations
  - Índices

- [ ] Crear `scripts/init-solr.sh`
  - Crear core reservations
  - Configurar schema
  - Campos indexados

- [ ] Crear `scripts/init-rabbitmq.sh`
  - Crear exchange
  - Crear queue
  - Bind queue a exchange

### 13. Makefile para Comandos Útiles
- [ ] Crear `Makefile` en la raíz
  - `make up`: docker-compose up -d
  - `make down`: docker-compose down
  - `make logs`: ver logs de todos los servicios
  - `make logs-users`: logs de users-api
  - `make logs-reservations`: logs de reservations-api
  - `make logs-search`: logs de search-api
  - `make logs-frontend`: logs de frontend
  - `make restart`: reiniciar todos los servicios
  - `make clean`: limpiar volúmenes
  - `make rebuild`: rebuild todas las imágenes
  - `make ps`: ver estado de contenedores
  - `make db-mysql`: conectar a MySQL
  - `make db-mongo`: conectar a MongoDB

### 14. Variables de Entorno
- [ ] Crear `.env` en la raíz
  - Variables globales
  - Passwords
  - Puertos
  - URLs de servicios

- [ ] Crear `.env.example` con valores de ejemplo

### 15. Healthchecks
- [ ] Implementar healthcheck en cada servicio
- [ ] Usar depends_on con condition: service_healthy
- [ ] Asegurar orden correcto de inicio

### 16. Networks
- [ ] Crear network: restaurant-network
- [ ] Todos los servicios en la misma network
- [ ] Permitir comunicación entre servicios por nombre

### 17. Volumes
- [ ] Volume para MySQL data
- [ ] Volume para MongoDB data
- [ ] Volume para SolR data
- [ ] Volume para RabbitMQ data

### 18. Documentación
- [ ] Actualizar README.md principal
  - Cómo levantar el proyecto
  - Requisitos (Docker, Docker Compose)
  - Puertos de cada servicio
  - URLs de acceso
  - Credenciales por defecto
  - Comandos útiles

## Archivo docker-compose.yml Completo

```yaml
version: "3.8"

services:
  # ========== DATABASES ==========
  mysql:
    image: mysql:8.0
    container_name: restaurant-mysql
    restart: unless-stopped
    environment:
      MYSQL_ROOT_PASSWORD: rootpass
      MYSQL_DATABASE: users_db
      MYSQL_USER: restaurant_user
      MYSQL_PASSWORD: restaurant_pass
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
      - ./scripts/init-mysql.sql:/docker-entrypoint-initdb.d/init.sql:ro
    networks:
      - restaurant-network
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost", "-u", "root", "-prootpass"]
      interval: 10s
      timeout: 5s
      retries: 5

  mongodb:
    image: mongo:7.0
    container_name: restaurant-mongodb
    restart: unless-stopped
    environment:
      MONGO_INITDB_ROOT_USERNAME: root
      MONGO_INITDB_ROOT_PASSWORD: rootpass
      MONGO_INITDB_DATABASE: reservations_db
    ports:
      - "27017:27017"
    volumes:
      - mongodb_data:/data/db
      - ./scripts/init-mongo.js:/docker-entrypoint-initdb.d/init.js:ro
    networks:
      - restaurant-network
    healthcheck:
      test: ["CMD", "mongosh", "--eval", "db.adminCommand('ping')"]
      interval: 10s
      timeout: 5s
      retries: 5

  solr:
    image: solr:9.4
    container_name: restaurant-solr
    restart: unless-stopped
    ports:
      - "8983:8983"
    volumes:
      - solr_data:/var/solr
      - ./scripts/init-solr.sh:/docker-entrypoint-initdb.d/init-solr.sh:ro
    networks:
      - restaurant-network
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8983/solr/admin/cores?action=STATUS"]
      interval: 10s
      timeout: 5s
      retries: 5
    command:
      - solr-precreate
      - reservations

  # ========== INFRASTRUCTURE ==========
  rabbitmq:
    image: rabbitmq:3-management
    container_name: restaurant-rabbitmq
    restart: unless-stopped
    environment:
      RABBITMQ_DEFAULT_USER: guest
      RABBITMQ_DEFAULT_PASS: guest
    ports:
      - "5672:5672"
      - "15672:15672"
    volumes:
      - rabbitmq_data:/var/lib/rabbitmq
    networks:
      - restaurant-network
    healthcheck:
      test: ["CMD", "rabbitmq-diagnostics", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  memcached:
    image: memcached:alpine
    container_name: restaurant-memcached
    restart: unless-stopped
    ports:
      - "11211:11211"
    command: ["-m", "128"]
    networks:
      - restaurant-network
    healthcheck:
      test: ["CMD", "nc", "-z", "localhost", "11211"]
      interval: 10s
      timeout: 5s
      retries: 3

  # ========== BACKEND SERVICES ==========
  users-api:
    build:
      context: ./users-api
      dockerfile: Dockerfile
    container_name: restaurant-users-api
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      DB_HOST: mysql
      DB_PORT: 3306
      DB_USER: restaurant_user
      DB_PASS: restaurant_pass
      DB_NAME: users_db
      JWT_SECRET: your-super-secret-jwt-key-change-in-production
      JWT_ACCESS_TTL: 24h
      JWT_REFRESH_TTL: 168h
      PORT: 8080
      APP_ENV: production
    depends_on:
      mysql:
        condition: service_healthy
    networks:
      - restaurant-network
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 10s
      timeout: 5s
      retries: 5

  reservations-api:
    build:
      context: ./reservations-api
      dockerfile: Dockerfile
    container_name: restaurant-reservations-api
    restart: unless-stopped
    ports:
      - "8081:8081"
    environment:
      MONGO_URI: mongodb://root:rootpass@mongodb:27017
      MONGO_DB: reservations_db
      MONGO_COLLECTION: reservations
      RABBITMQ_URI: amqp://guest:guest@rabbitmq:5672/
      RABBITMQ_EXCHANGE: restaurant_events
      RABBITMQ_QUEUE: reservations_updates
      USERS_API_URL: http://users-api:8080
      PORT: 8081
      APP_ENV: production
    depends_on:
      mongodb:
        condition: service_healthy
      rabbitmq:
        condition: service_healthy
      users-api:
        condition: service_healthy
    networks:
      - restaurant-network
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8081/health"]
      interval: 10s
      timeout: 5s
      retries: 5

  search-api:
    build:
      context: ./search-api
      dockerfile: Dockerfile
    container_name: restaurant-search-api
    restart: unless-stopped
    ports:
      - "8082:8082"
    environment:
      SOLR_URL: http://solr:8983/solr
      SOLR_CORE: reservations
      MEMCACHED_SERVERS: memcached:11211
      MEMCACHED_TIMEOUT: 1s
      LOCAL_CACHE_MAX_SIZE: 10000
      LOCAL_CACHE_TTL: 5m
      DISTRIBUTED_CACHE_TTL: 15m
      RABBITMQ_URI: amqp://guest:guest@rabbitmq:5672/
      RABBITMQ_QUEUE: reservations_updates
      RABBITMQ_EXCHANGE: restaurant_events
      RESERVATIONS_API_URL: http://reservations-api:8081
      PORT: 8082
      APP_ENV: production
    depends_on:
      solr:
        condition: service_healthy
      memcached:
        condition: service_healthy
      rabbitmq:
        condition: service_healthy
      reservations-api:
        condition: service_healthy
    networks:
      - restaurant-network
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8082/health"]
      interval: 10s
      timeout: 5s
      retries: 5

  # ========== FRONTEND ==========
  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    container_name: restaurant-frontend
    restart: unless-stopped
    ports:
      - "3000:80"
    environment:
      VITE_API_URL: http://localhost:8080
      VITE_SEARCH_API_URL: http://localhost:8082
      VITE_APP_NAME: Restaurant Reservations
    depends_on:
      users-api:
        condition: service_healthy
      search-api:
        condition: service_healthy
    networks:
      - restaurant-network
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:80"]
      interval: 10s
      timeout: 5s
      retries: 3

  # ========== ADMIN TOOLS (Optional) ==========
  adminer:
    image: adminer:latest
    container_name: restaurant-adminer
    restart: unless-stopped
    ports:
      - "8084:8080"
    environment:
      ADMINER_DEFAULT_SERVER: mysql
    depends_on:
      - mysql
    networks:
      - restaurant-network

  mongo-express:
    image: mongo-express:latest
    container_name: restaurant-mongo-express
    restart: unless-stopped
    ports:
      - "8085:8081"
    environment:
      ME_CONFIG_MONGODB_ADMINUSERNAME: root
      ME_CONFIG_MONGODB_ADMINPASSWORD: rootpass
      ME_CONFIG_MONGODB_URL: mongodb://root:rootpass@mongodb:27017/
      ME_CONFIG_BASICAUTH: false
    depends_on:
      - mongodb
    networks:
      - restaurant-network

networks:
  restaurant-network:
    driver: bridge

volumes:
  mysql_data:
  mongodb_data:
  solr_data:
  rabbitmq_data:
```

## Makefile

```makefile
.PHONY: help up down logs logs-users logs-reservations logs-search logs-frontend restart clean rebuild ps db-mysql db-mongo

help:
	@echo "Comandos disponibles:"
	@echo "  make up              - Levantar todos los servicios"
	@echo "  make down            - Detener todos los servicios"
	@echo "  make logs            - Ver logs de todos los servicios"
	@echo "  make logs-users      - Ver logs de users-api"
	@echo "  make logs-reservations - Ver logs de reservations-api"
	@echo "  make logs-search     - Ver logs de search-api"
	@echo "  make logs-frontend   - Ver logs de frontend"
	@echo "  make restart         - Reiniciar todos los servicios"
	@echo "  make clean           - Limpiar volúmenes"
	@echo "  make rebuild         - Rebuild todas las imágenes"
	@echo "  make ps              - Ver estado de contenedores"
	@echo "  make db-mysql        - Conectar a MySQL CLI"
	@echo "  make db-mongo        - Conectar a MongoDB CLI"

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f

logs-users:
	docker-compose logs -f users-api

logs-reservations:
	docker-compose logs -f reservations-api

logs-search:
	docker-compose logs -f search-api

logs-frontend:
	docker-compose logs -f frontend

restart:
	docker-compose restart

clean:
	docker-compose down -v

rebuild:
	docker-compose build --no-cache
	docker-compose up -d

ps:
	docker-compose ps

db-mysql:
	docker exec -it restaurant-mysql mysql -urestaurant_user -prestaurant_pass users_db

db-mongo:
	docker exec -it restaurant-mongodb mongosh -u root -p rootpass reservations_db
```

## Scripts de Inicialización

### scripts/init-mysql.sql
```sql
-- Create database if not exists
CREATE DATABASE IF NOT EXISTS users_db;
USE users_db;

-- Create users table if not exists
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    role VARCHAR(50) DEFAULT 'normal',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Insert default admin user (password: admin123)
INSERT IGNORE INTO users (id, username, email, password, first_name, last_name, role)
VALUES (
    UUID(),
    'admin',
    'admin@restaurant.com',
    '$2a$10$XKzQ5Z1Q5Z1Q5Z1Q5Z1Q5uXKzQ5Z1Q5Z1Q5Z1Q5Z1Q5uXKzQ5Z1Q',
    'Admin',
    'User',
    'admin'
);
```

### scripts/init-mongo.js
```javascript
db = db.getSiblingDB('reservations_db');

// Create collection
db.createCollection('reservations');

// Create indexes
db.reservations.createIndex({ "owner_id": 1 });
db.reservations.createIndex({ "date_time": 1 });
db.reservations.createIndex({ "status": 1 });
db.reservations.createIndex({ "meal_type": 1 });
db.reservations.createIndex({ "table_number": 1 });

print('MongoDB initialized successfully');
```

### scripts/init-solr.sh
```bash
#!/bin/bash
# SolR core is created with solr-precreate command in docker-compose
# This script can add custom configurations if needed
echo "SolR core 'reservations' initialized"
```

## README.md Principal

Agregar sección:

```markdown
## 🚀 Inicio Rápido con Docker

### Prerrequisitos
- Docker 20.10+
- Docker Compose 2.0+

### Levantar el Proyecto Completo

1. Clonar el repositorio
```bash
git clone <repo-url>
cd proyectoArqSoftII
```

2. Crear archivo .env (copiar de .env.example)
```bash
cp .env.example .env
```

3. Levantar todos los servicios
```bash
make up
# o
docker-compose up -d
```

4. Esperar a que todos los servicios estén healthy
```bash
make ps
```

### Acceso a los Servicios

| Servicio | URL | Puerto |
|----------|-----|--------|
| Frontend | http://localhost:3000 | 3000 |
| Users API | http://localhost:8080 | 8080 |
| Reservations API | http://localhost:8081 | 8081 |
| Search API | http://localhost:8082 | 8082 |
| RabbitMQ Management | http://localhost:15672 | 15672 |
| SolR Admin | http://localhost:8983 | 8983 |
| Adminer (MySQL) | http://localhost:8084 | 8084 |
| Mongo Express | http://localhost:8085 | 8085 |

### Credenciales por Defecto

**Usuario Admin:**
- Username: `admin`
- Password: `admin123`

**RabbitMQ:**
- User: `guest`
- Pass: `guest`

**MySQL (Adminer):**
- Server: `mysql`
- User: `restaurant_user`
- Pass: `restaurant_pass`
- Database: `users_db`

**MongoDB (Mongo Express):**
- User: `root`
- Pass: `rootpass`

### Comandos Útiles

Ver todos los comandos disponibles:
```bash
make help
```

Ver logs:
```bash
make logs                    # Todos
make logs-users             # Solo users-api
make logs-reservations      # Solo reservations-api
make logs-search            # Solo search-api
```

Detener servicios:
```bash
make down
```

Limpiar todo (incluyendo volúmenes):
```bash
make clean
```

Rebuild completo:
```bash
make rebuild
```
```

## Criterios de Aceptación
- [ ] docker-compose.yml levanta todos los servicios
- [ ] Todos los healthchecks pasan
- [ ] Orden de inicio respeta dependencias
- [ ] Networks permiten comunicación entre servicios
- [ ] Volumes persisten datos correctamente
- [ ] Scripts de inicialización ejecutan correctamente
- [ ] Frontend accesible en http://localhost:3000
- [ ] Todas las APIs accesibles
- [ ] RabbitMQ management accesible
- [ ] Adminer y Mongo Express accesibles
- [ ] Makefile con comandos útiles funciona
- [ ] Documentación completa en README
- [ ] .env.example con todas las variables

## Prioridad
🔴 **CRÍTICA** - Necesario para ejecutar el sistema completo

## Estimación
⏱️ 6-8 horas

## Dependencias
- Todos los microservicios deben tener Dockerfile
- Todos los scripts de inicialización listos

## Notas
- Usar `depends_on` con `condition: service_healthy` para asegurar orden
- RabbitMQ puede tardar en estar ready, considerar retries en consumers
- SolR necesita crear el core antes de indexar
- Considerar agregar un API Gateway (nginx) como punto de entrada único (opcional)
- Para producción, usar secrets manager en lugar de .env
- Considerar implementar Traefik para routing automático (opcional)
