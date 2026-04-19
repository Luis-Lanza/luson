# Exploración de Fase 0: Fundación Técnica - Battery-POS

## 1. Orden Recomendado de Setup y Comandos

### Backend (Go)
```bash
# Inicializar módulo Go
go mod init github.com/Luis-Lanza/luson

# Crear estructura de directorios hexagonal/clean architecture
mkdir -p cmd/server
mkdir -p internal/{domain,ports,application,infrastructure,config}
mkdir -p internal/infrastructure/{postgres,http,middleware,handlers,dto,websocket,pdf,jwt}
mkdir -p internal/application/service
mkdir -p internal/infrastructure/postgres/{migrations,repo}
mkdir -p internal/infrastructure/http/{handlers,middleware,dto}
mkdir -p internal/domain
mkdir -p pkg/{logger,validator}
mkdir -p scripts
mkdir -p migrations

# Dependencias clave (versiones actuales al 2026-04-18)
go get github.com/gin-gonic/gin@v1.10.0  # Router lightweight y rápido
go get github.com/jackc/pgx/v5@v5.8.0     # Driver PostgreSQL nativo para Go
go get github.com/golang-migrate/migrate/v4@v4.18.1  # Migraciones
go get github.com/golang-jwt/jwt/v5@v5.4.2  # JWT tokens
go get golang.org/x/crypto/bcrypt@v0.20.0  # Hash de contraseñas
go get github.com/google/uuid@v1.6.0      # UUID generation
go get github.com/joho/godotenv@v1.5.1    # Carga de variables de entorno
go get github.com/signalfx/golang/v3/sdktrace@v0.40.0  # OpenTelemetry tracing
go get github.com/prometheus/client_golang@v1.20.0  # Métricas
```

### Frontend (Vite + React)
```bash
# Crear proyecto Vite con React y TypeScript
npm create vite@latest client -- --template react-ts
cd client

# Dependencias clave
npm install zustand@5.0.0
npm install tailwindcss@4.0.0 postcss@8.4.0 autoprefixer@10.4.0
npm install zod@4.0.0
npm install react-router-dom@7.0.0
npm install rxdb@15.0.0 rxdb-plugin-encrypt@15.0.0
npm install jspdf@2.5.0
npm install vitest@2.0.0 @vitest/coverage-v8@2.0.0 @testing-library/react@14.0.0 @testing-library/jest-dom@6.4.0 @testing-library/user-event@14.5.0
npm install playwright@1.45.0
npm install -D @types/node@22.0.0

# Inicializar Tailwind 4
npx tailwindcss init -p
```

### Infrastructure
```bash
# Docker setup para PostgreSQL (opcional pero recomendado para desarrollo)
mkdir -p docker
cat > docker/docker-compose.yml << 'EOF'
version: '3.8'
services:
  postgres:
    image: postgres:16-alpine
    container_name: luson_postgres
    environment:
      POSTGRES_DB: luson
      POSTGRES_USER: luson_user
      POSTGRES_PASSWORD: luson_pass
      POSTGRES_HOST_AUTH_METHOD: trust
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U $$POSTGRES_USER"]
      interval: 5s
      timeout: 5s
      retries: 5
volumes:
  postgres_data:
EOF
```

## 2. Listas de Dependencias con Versiones

### Backend Go Dependencies
| Dependencia | Versión | Propósito |
|-------------|---------|-----------|
| gin-gonic/gin | v1.10.0 | Router web ligero y performante |
| jackc/pgx/v5 | v5.8.0 | Driver PostgreSQL nativo |
| golang-migrate/migrate/v4 | v4.18.1 | Sistema de migraciones |
| golang-jwt/jwt/v5 | v5.4.2 | JWT para autenticación |
| golang.org/x/crypto/bcrypt | v0.20.0 | Hash de contraseñas |
| google/uuid | v1.6.0 | Generación de UUIDs |
| joho/godotenv | v1.5.1 | Carga de variables de entorno |
| signalfx/golang/v3/sdktrace | v0.40.0 | Tracing distribuido |
| prometheus/client_golang | v1.20.0 | Métricas para monitoreo |

### Frontend Dependencies
| Dependencia | Versión | Propósito |
|-------------|---------|-----------|
| vite | latest | Bundler rápido |
| react | 19.0.0 | Biblioteca UI |
| react-dom | 19.0.0 | Renderizado DOM |
| zustand | 5.0.0 | Estado global |
| tailwindcss | 4.0.0 | Styling utility-first |
| zod | 4.0.0 | Validación de esquemas |
| react-router-dom | 7.0.0 | Routing cliente |
| rxdb | 15.0.0 | Base de datos local offline |
| rxdb-plugin-encrypt | 15.0.0 | Encriptación AES-256 para RxDB |
| jspdf | 2.5.0 | Generación de PDFs client-side |
| vitest | 2.0.0 | Framework de testing |
| @testing-library/react | 14.0.0 | Utilidades de testing React |
| playwright | 1.45.0 | Testing E2E |

## 3. Estructura de Archivos Después de Fase 0

```
/home/luis/velay/luson_2
├── PRD.md
├── .gitignore
├── .atl/
│   └── skill-registry.md
├── go.mod
├── go.sum
├── Makefile
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── domain/
│   │   ├── user.go
│   │   ├── branch.go
│   │   ├── product.go
│   │   ├── stock.go
│   │   ├── acid.go
│   │   ├── sale.go
│   │   ├── warranty.go
│   │   ├── transfer.go
│   │   ├── scrap.go
│   │   ├── expense.go
│   │   ├── cash_close.go
│   │   ├── cash_movement.go
│   │   ├── inventory_adjustment.go
│   │   ├── notification.go
│   │   └── maintenance.go
│   ├── ports/
│   │   ├── repositories.go
│   │   ├── services.go
│   │   └── handlers.go
│   ├── application/
│   │   ├── auth_service.go
│   │   ├── product_service.go
│   │   ├── stock_service.go
│   │   ├── sale_service.go
│   │   ├── warranty_service.go
│   │   ├── transfer_service.go
│   │   ├── acid_service.go
│   │   ├── scrap_service.go
│   │   ├── expense_service.go
│   │   ├── cash_close_service.go
│   │   ├── cash_movement_service.go
│   │   ├── inventory_adjustment_service.go
│   │   ├── notification_service.go
│   │   ├── sync_service.go
│   │   └── report_service.go
│   ├── infrastructure/
│   │   ├── postgres/
│   │   │   ├── migrations/
│   │   │   │   └── 0001_initial_schema.sql
│   │   │   ├── user_repo.go
│   │   │   ├── branch_repo.go
│   │   │   ├── product_repo.go
│   │   │   ├── stock_repo.go
│   │   │   ├── sale_repo.go
│   │   │   ├── warranty_repo.go
│   │   │   ├── transfer_repo.go
│   │   │   ├── acid_repo.go
│   │   │   ├── scrap_repo.go
│   │   │   ├── expense_repo.go
│   │   │   ├── cash_close_repo.go
│   │   │   ├── cash_movement_repo.go
│   │   │   ├── inventory_adjustment_repo.go
│   │   │   ├── notification_repo.go
│   │   │   └── sync_repo.go
│   │   ├── http/
│   │   │   ├── router.go
│   │   │   ├── middleware/
│   │   │   │   ├── auth.go
│   │   │   │   ├── cors.go
│   │   │   │   └── logging.go
│   │   │   ├── handlers/
│   │   │   │   ├── auth_handler.go
│   │   │   │   ├── product_handler.go
│   │   │   │   ├── stock_handler.go
│   │   │   │   ├── sale_handler.go
│   │   │   │   ├── warranty_handler.go
│   │   │   │   ├── transfer_handler.go
│   │   │   │   ├── acid_handler.go
│   │   │   │   ├── scrap_handler.go
│   │   │   │   ├── expense_handler.go
│   │   │   │   ├── cash_close_handler.go
│   │   │   │   ├── cash_movement_handler.go
│   │   │   │   ├── inventory_adjustment_handler.go
│   │   │   │   ├── sync_handler.go
│   │   │   │   ├── report_handler.go
│   │   │   │   └── notification_handler.go
│   │   │   └── dto/
│   │   │       ├── requests.go
│   │   │       └── responses.go
│   │   ├── websocket/
│   │   │   └── hub.go
│   │   ├── pdf/
│   │   │   ├── warranty_pdf.go
│   │   │   ├── receipt_pdf.go
│   │   │   └── cash_close_pdf.go
│   │   └── jwt/
│   │       └── token.go
│   └── config/
│       └── config.go
├── pkg/
│   ├── logger/
│   │   └── logger.go
│   └── validator/
│       └── validator.go
├── scripts/
│   └── generate-schemas/
│       └── main.go
├── migrations/
│   └── *.sql
├── client/ (Frontend)
│   ├── index.html
│   ├── package.json
│   ├── vite.config.ts
│   ├── tsconfig.json
│   ├── tailwind.config.cjs
│   ├── postcss.config.cjs
│   ├── src/
│   │   ├── main.tsx
│   │   ├── App.tsx
│   │   ├── index.css
│   │   ├── lib/
│   │   │   ├── rxdb.ts
│   │   │   ├── zustand/
│   │   │   │   ├── store.ts
│   │   │   │   ├── slices/
│   │   │   │   │   ├── authSlice.ts
│   │   │   │   │   ├── cartSlice.ts
│   │   │   │   │   ├── uiSlice.ts
│   │   │   │   │   └── notificationSlice.ts
│   │   │   ├── utils/
│   │   │   │   ├── zod/
│   │   │   │   │   ├── authSchema.ts
│   │   │   │   │   ├── productSchema.ts
│   │   │   │   │   └── saleSchema.ts
│   │   │   │   ├── api.ts
│   │   │   │   ├── websocket.ts
│   │   │   │   └── helpers.ts
│   │   │   ├── components/
│   │   │   │   ├── layout/
│   │   │   │   ├── pages/
│   │   │   │   └── ui/
│   │   │   ├── hooks/
│   │   │   ├── services/
│   │   │   └── tests/
│   │   │       ├── unit/
│   │   │       └── e2e/
│   ├── public/
│   │   ├── manifest.json
│   │   └── service-worker.js
│   └── vitest.config.ts
├── docker/
│   └── docker-compose.yml
└── README.md (generado después)
```

## 4. Decisiones y Tradeoffs

### Backend Architecture Decisions
1. **Hexagonal/Clean Architecture**: Elegido por separación clara de preocupaciones y facilidad de testing
   - Tradeoff: Más archivos inicialmente pero mejor mantenibilidad a largo plazo

2. **Gin como Router**: En lugar de el estándar net/http
   - Pros: Middleware built-in, routing express-like, buen rendimiento
   - Cons: Dependencia externa adicional

3. **PGX sobre database/sql**: Driver PostgreSQL nativo
   - Pros: Mejor performance, acceso a características específicas de PostgreSQL
   - Cons: Menos portable a otras bases de datos

4. **Migraciones con golang-migrate**: En lugar de embedir SQL en Go
   - Pros: Separación clara, rollbacks fáciles, herramientas CLI
   - Cons: Otra dependencia y proceso de build

### Frontend Architecture Decisions
1. **Vite + React 19**: En lugar de Create React App o Next.js
   - Pros: Build extremadamente rápido, HMR excelente, control total
   - Cons: Necesita configuración manual para SSR si se requiere en futuro

2. **Zustand 5**: En lugar de Redux o Context API
   - Pros: API minimalista, excelente performance, boilerplate mínimo
   - Cons: Menos herramientas de dev integradas que Redux

3. **Tailwind 4**: En lugar de CSS tradicional o CSS-in-JS
   - Pros: Utility-first rápido, tamaño de bundle pequeño, consistency
   - Cons: Curva de aprendizaje, clases largas en JSX

4. **RxDB con encriptación**: En lugar de localStorage o IndexedDB directo
   - Pros: API tipo SQL, sincronización automática, encriptación built-in
   - Cons: Bundle size mayor, complejidad adicional

### Infraestructura Decisions
1. **Docker para PostgreSQL**: En lugar de instalación local
   - Pros: Consistencia entre entornos, fácil setup/teardown, versionado
   - Cons: Requiere Docker, ligeramente más complejo inicialmente

2. **Variables de Entorno con .env**: En lugar de config hardcodeada o flags
   - Pros: 12-factor app compliant, diferente config por entorno
   - Cons: Riesgo de commit accidental de .env (mitigado con .gitignore)

## 5. Riesgos y Gotchas

### Backend Gotchas
1. **Go Modules**: Asegurarse de ejecutar `go mod tidy` después de agregar/quitar dependencias
2. **Circular Dependencies**: La arquitectura hexagonal ayuda pero hay que tener cuidado con dependencias entre paquetes internos
3. **Migraciones**: Siempre testear migraciones en ambiente de staging antes de producción
4. **Manejo de Errores**: Go no tiene excepciones, hay que ser explícito con error handling
5. **Testing de Base de Datos**: Necesitar testcontainers o similar para tests de integración reales

### Frontend Gotchas
1. **Tailwind 4 Breaking Changes**: Si se actualiza desde v3, hay cambios en configuración y sintaxis
2. **RxDB Complejidad**: La curva de aprendizaje puede ser alta inicialmente
3. **Service Worker Debugging**: Los service workers pueden ser difíciles de depurar debido al caching
4. **Hydration Mismatch**: Con React 18+ y rendering del servidor, hay que tener cuidado con useEffect
5. **Bundle Size**: RxDB y zod pueden aumentar el tamaño del bundle, considerar code-splitting

### Infraestructura Gotchas
1. **Docker en Desarrollo vs Producción**: Las configuraciones pueden diferir significativamente
2. **Migraciones de Base de Datos**: Riesgo de migraciones fallidas en producción
3. **Variables de Entorno**: Diferentes entornos (dev, staging, prod) requieren manejo cuidadoso
4. **Backup y Recovery**: Necesitar estrategia para backups de PostgreSQL desde el día uno

### TDD Strict Mode Considerations
1. **Cobertura de Tests**: El modo TDD estricto requerirá alta cobertura (>80%)
2. **Tests de Unidad vs Integración**: Balance entre rapidez de tests unitarios y realismo de tests de integración
3. **Mocks vs Reales**: En Go, usar interfaces para mocking; en frontend, vitest con mocking integrado
4. **Testing de WebSockets**: Más complejo que testing REST, puede requerir bibliotecas especiales
5. **Testing de Service Workers**: Requiere enfoques especiales o evitarlos en tests unitarios

## Próximos Pasos Recomendados

1. Ejecutar los comandos de setup en el orden recomendado
2. Crear archivos básicos de entrada (main.go, main.tsx)
3. Implementar salud de los endpoints básicos (health check)
4. Configurar variables de entorno iniciales
5. Crear primera migración de base de datos (schema inicial)
6. Configurar vitest y Playwright para testing
7. Establecer pipeline de CI básico (incluso si es local inicialmente)