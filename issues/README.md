# 📋 Issue Packets - Sistema de Gestión de Restaurante

Este directorio contiene los issue packets detallados para implementar el sistema completo de microservicios del trabajo final de Arquitectura de Software II.

## 🎯 Objetivo del Proyecto

Desarrollar un sistema web completo basado en microservicios para gestión de reservas de restaurante, utilizando:
- **Backend:** Go con arquitectura hexagonal (MVC)
- **Bases de Datos:** MySQL, MongoDB, SolR
- **Infraestructura:** RabbitMQ, Memcached
- **Frontend:** React
- **Containerización:** Docker & Docker Compose

## 📦 Issues Disponibles

### [Issue #1: Users API](./01-users-api-improvements.md)
**Prioridad:** 🔴 ALTA
**Estimación:** 8-12 horas
**Estado:** Parcialmente implementado

**Descripción:**
Completar la API de usuarios con autenticación JWT, hashing de passwords, y soporte para usuarios normales y administradores.

**Tecnologías:**
- Go + Gin
- MySQL + GORM
- JWT
- bcrypt

**Endpoints Principales:**
- `POST /api/users/register` - Crear usuario
- `POST /api/users/login` - Login con JWT
- `GET /api/users/:id` - Obtener usuario
- Tests en service layer

---

### [Issue #2: Reservations API](./02-reservations-api.md)
**Prioridad:** 🔴 CRÍTICA
**Estimación:** 16-20 horas
**Estado:** Por implementar

**Descripción:**
Crear el microservicio de reservas (entidad principal) con MongoDB, implementando CRUD completo, validación de usuarios, cálculo concurrente con goroutines, y notificaciones a RabbitMQ.

**Tecnologías:**
- Go + Gin
- MongoDB
- RabbitMQ Publisher
- Concurrencia (goroutines + channels + waitgroups)

**Endpoints Principales:**
- `POST /api/reservations` - Crear reserva
- `GET /api/reservations/:id` - Obtener reserva
- `PUT /api/reservations/:id` - Actualizar
- `DELETE /api/reservations/:id` - Eliminar
- `POST /api/reservations/:id/confirm` - Confirmar (con cálculo concurrente)
- `GET /api/reservations/user/:user_id` - Reservas por usuario

**Características Especiales:**
- ✅ Validación de owner contra Users API (HTTP)
- ✅ Cálculo concurrente con Go Routines
- ✅ Publicación de eventos a RabbitMQ
- ✅ Tests con mocks

---

### [Issue #3: Search API](./03-search-api.md)
**Prioridad:** 🟠 ALTA
**Estimación:** 20-24 horas
**Estado:** Por implementar

**Descripción:**
Implementar API de búsqueda con SolR, doble capa de caché (CCache local + Memcached distribuida), y consumidor de RabbitMQ para sincronización de índices.

**Tecnologías:**
- Go + Gin
- SolR (motor de búsqueda)
- CCache (caché local)
- Memcached (caché distribuida)
- RabbitMQ Consumer

**Endpoints Principales:**
- `GET /api/search` - Búsqueda paginada con filtros
- `GET /api/search/:id` - Obtener por ID
- `GET /api/search/stats` - Estadísticas de caché
- `POST /api/search/reindex` - Reindexar todo

**Características Especiales:**
- ✅ Búsqueda con filtros dinámicos
- ✅ Paginación y sorting
- ✅ Doble capa de caché
- ✅ Sincronización automática vía RabbitMQ
- ✅ Cache invalidation inteligente

---

### [Issue #4: Frontend React](./04-frontend-react.md)
**Prioridad:** 🟠 ALTA
**Estimación:** 24-30 horas
**Estado:** Por implementar

**Descripción:**
Desarrollar la SPA (Single Page Application) en React con todas las pantallas requeridas, autenticación JWT, búsqueda, y panel de administración.

**Tecnologías:**
- React 18 + Vite
- React Router v6
- Axios + React Query
- React Hook Form
- Tailwind CSS
- Docker + Nginx

**Pantallas:**
1. **Login** - Autenticación
2. **Registro** - Crear cuenta
3. **Home/Búsqueda** - Búsqueda con filtros y paginación
4. **Detalles** - Vista completa de reserva
5. **Mis Reservas** - Reservas del usuario
6. **Admin** - Panel de administración (solo admins)

**Flujo Principal:**
```
Login → Home/Búsqueda → Detalles → Confirmar Reserva → Success
```

---

### [Issue #5: Docker Compose](./05-docker-compose-integration.md)
**Prioridad:** 🔴 CRÍTICA
**Estimación:** 6-8 horas
**Estado:** Por implementar

**Descripción:**
Orquestar todos los servicios con Docker Compose, incluyendo microservicios, bases de datos, infraestructura y herramientas de administración.

**Servicios Incluidos:**

**Backend (3):**
- users-api (8080)
- reservations-api (8081)
- search-api (8082)

**Frontend (1):**
- frontend (3000)

**Bases de Datos (3):**
- MySQL (3306)
- MongoDB (27017)
- SolR (8983)

**Infraestructura (2):**
- RabbitMQ (5672, 15672)
- Memcached (11211)

**Admin Tools (2):**
- Adminer (8084)
- Mongo Express (8085)

**Características:**
- ✅ Healthchecks en todos los servicios
- ✅ Depends_on con conditions
- ✅ Networks para comunicación
- ✅ Volumes persistentes
- ✅ Scripts de inicialización
- ✅ Makefile con comandos útiles

---

## 🗺️ Orden de Implementación Sugerido

### Fase 1: Backend Core (Primera Entrega)
1. **Issue #1** - Users API (completar)
2. **Issue #2** - Reservations API (core sin concurrencia)
3. **Issue #3** - Search API (básica sin caché avanzada)
4. **Issue #5** - Docker Compose (básico)

**Entregable:** Flujo Login → Búsqueda → Detalle → Acción

### Fase 2: Frontend (Primera Entrega)
4. **Issue #4** - Frontend React (pantallas básicas)
   - Login
   - Home/Búsqueda
   - Detalles
   - Confirmación

**Objetivo:** Sistema funcional para primera entrega

### Fase 3: Completar Funcionalidades (Entrega Final)
5. Agregar concurrencia en Reservations API
6. Implementar doble caché en Search API
7. Agregar pantallas de Registro, Mis Reservas, Admin
8. Mejorar Docker Compose con todos los servicios
9. Tests completos
10. Documentación final

---

## 📊 Arquitectura del Sistema

```
┌─────────────────────────────────────────────────────────────┐
│                         Frontend                            │
│                      React (Port 3000)                      │
└────────────────┬────────────────────────┬───────────────────┘
                 │                        │
                 │ HTTP                   │ HTTP
                 │                        │
        ┌────────▼────────┐      ┌────────▼────────┐
        │   Users API     │      │   Search API    │
        │   (Port 8080)   │      │   (Port 8082)   │
        │   MySQL + GORM  │      │   SolR + Cache  │
        └────────┬────────┘      └────────┬────────┘
                 │                        │
                 │ HTTP                   │ HTTP
                 │                        │
        ┌────────▼────────────────────────▼────────┐
        │         Reservations API                 │
        │            (Port 8081)                   │
        │       MongoDB + Concurrency              │
        └────────┬───────────────────────┬─────────┘
                 │                       │
                 │ Publish               │ Consume
                 │                       │
        ┌────────▼───────────────────────▼─────────┐
        │              RabbitMQ                    │
        │         (Exchange: restaurant_events)    │
        └──────────────────────────────────────────┘
```

### Flujo de Datos

**1. Crear Reserva:**
```
Frontend → Reservations API → Validate User (Users API) → MongoDB → RabbitMQ
```

**2. Sincronizar Búsqueda:**
```
RabbitMQ → Search API → Fetch Full Data (Reservations API) → Index in SolR
```

**3. Búsqueda con Caché:**
```
Frontend → Search API → Check CCache → Check Memcached → Query SolR
```

---

## 🧪 Testing

Cada microservicio debe incluir:
- **Unit Tests** en service layer
- **Mocks** de repositories y clients externos
- **Coverage** mínimo 70%

Archivos de tests:
- `users-api/internal/service/user_service_test.go`
- `reservations-api/internal/service/reservation_service_test.go`
- `search-api/internal/service/search_service_test.go`

---

## 📚 Tecnologías por Microservicio

| Servicio | Framework | Database | Otros |
|----------|-----------|----------|-------|
| users-api | Gin | MySQL (GORM) | JWT, bcrypt |
| reservations-api | Gin | MongoDB | RabbitMQ (pub), Goroutines |
| search-api | Gin | SolR | RabbitMQ (sub), CCache, Memcached |
| frontend | React + Vite | - | React Router, Axios, Tailwind |

---

## 🚀 Inicio Rápido

1. **Leer el enunciado completo:** `Enunciado - Arq. Soft II 2025 - Electivo.pdf`

2. **Revisar issue por issue:**
   - Cada issue tiene checklist completa
   - Estructura de archivos detallada
   - Ejemplos de código
   - Endpoints con requests/responses

3. **Seguir orden de implementación sugerido**

4. **Levantar servicios progresivamente:**
   ```bash
   # Primera vez
   docker-compose up mysql mongodb rabbitmq -d

   # Después de implementar users-api
   docker-compose up users-api -d

   # Y así sucesivamente...
   ```

5. **Usar Makefile para comandos comunes:**
   ```bash
   make up      # Levantar todo
   make logs    # Ver logs
   make down    # Detener
   make clean   # Limpiar
   ```

---

## 📞 Equipo

- Amuchastegui, Matias
- Rodriguez Richard, Lucas
- Sardoy, Blas
- Vargas, Santino

---

## 📝 Notas Importantes

### Primera Entrega (7-14 Nov)
Implementar flujo básico:
- ✅ Login
- ✅ Búsqueda
- ✅ Detalle
- ✅ Acción/Confirmación
- ✅ Docker funcional

**NO necesario para primera entrega:**
- ❌ Registro
- ❌ Mis Reservas
- ❌ Admin panel
- ❌ Cálculo concurrente completo

### Entrega Final
- ✅ Todo lo anterior
- ✅ Usuarios admin
- ✅ Pantalla de administración
- ✅ Vista de Registro
- ✅ Vista de Mis Acciones
- ✅ Cálculo concurrente completo
- ✅ Tests completos
- ✅ Documentación final

---

## 🔗 Referencias

- **Repo del curso:** https://github.com/ucc-arqsoft-2/clases2025
- **Repo del proyecto:** https://github.com/lucasrodriguezrichard/proyectoArqSoftII-Amuchastegui-RodriguezRichard-Sardoy-Vargas

---

**¡Buena suerte con el proyecto! 🚀**

Si tienen dudas sobre algún issue específico, revisen el archivo markdown correspondiente que contiene todos los detalles de implementación.
