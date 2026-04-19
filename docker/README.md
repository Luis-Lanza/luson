# Battery POS - Sistema de Gestión de Inventario y Punto de Venta

## Docker Setup

### 1. Configurar Variables de Entorno

Copia el archivo de ejemplo y configura tus variables:

```bash
cp .env.example .env
```

Edita `.env` con tus valores de producción (o deja los defaults para desarrollo local).

Para desarrollo local, también puedes usar `.env.local` (ya está incluido con valores de desarrollo).

### 2. Iniciar PostgreSQL

```bash
cd docker
docker compose up -d
```

Esto iniciará PostgreSQL 16 en el puerto 5432.

### 3. Verificar que PostgreSQL está funcionando

```bash
# Verificar estado del contenedor
docker compose ps

# Verificar que PostgreSQL está aceptando conexiones
docker compose exec postgres pg_isready -U postgres -d battery_pos

# Conectar con psql (opcional)
docker compose exec postgres psql -U postgres -d battery_pos
```

### 4. Detener PostgreSQL

```bash
docker compose down
```

Para eliminar también los datos persistentes:

```bash
docker compose down -v
```

## Comandos Make

Una vez que PostgreSQL esté corriendo, puedes usar los comandos de Make:

```bash
# Ejecutar migraciones
make migrate-up

# Revertir última migración
make migrate-down

# Iniciar servidor de desarrollo
make run

# Construir binario
make build

# Ejecutar tests
make test
```

## Troubleshooting

### Puerto 5432 ya está en uso

Si tienes otro PostgreSQL corriendo localmente, puedes:
1. Detener el otro servicio, o
2. Cambiar el puerto en `docker-compose.yml` (ej: `"5433:5432"`)

### Permisos denegados en volumen

En Linux/macOS, si tienes problemas de permisos:

```bash
sudo chown -R $USER:$USER docker/
```

### No se puede conectar desde Go

Asegúrate de que:
1. El contenedor está corriendo: `docker compose ps`
2. La variable `DATABASE_URL` en `.env` apunta a `localhost:5432`
3. El healthcheck pasa: `docker compose exec postgres pg_isready ...`

## Estructura del Proyecto

```
battery-pos/
├── docker/              # Docker compose y configuración
├── cmd/server/          # Entry point del servidor Go
├── internal/            # Código del backend Go
├── migrations/          # Migraciones SQL
├── client/             # Frontend React (Vite)
└── Makefile            # Comandos de desarrollo
```
