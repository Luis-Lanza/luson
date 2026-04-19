# PRD — Sistema de Gestión de Inventario y POS para Tienda de Baterías

**Versión:** 2.0
**Fecha:** 2025-04-18
**Estado:** Borrador actualizado

---

## 1. Descripción del Proyecto

Sistema de gestión de inventario y punto de venta (POS) offline-first para una cadena de tiendas de baterías para vehículos y motocicletas. El sistema centraliza la gestión de stock, ventas, garantías, ácido, chatarra y rendiciones entre un almacén central y múltiples sucursales, funcionando incluso sin conexión a internet.

---

## 2. Contexto del Negocio

### Estructura organizacional

```
PROVEEDORES
    │
    ▼
ALMACÉN (1 solo — pool de stock central)
    ├── Compra baterías y accesorios a proveedores
    ├── Compra ácido por separado
    ├── Distribuye a N sucursales (pedido o push directo)
    ├── Vende a mayoristas (tiendas de terceros)
    ├── Recibe chatarra de sucursales → vende a recicladoras
    ├── Recibe baterías falladas de sucursales → envía a fábrica
    └── Registra gastos operativos

SUCURSALES (N tiendas)
    ├── Reciben stock del almacén o de otras sucursales
    ├── Venden baterías y accesorios a clientes finales
    ├── Generan garantías y recibos en PDF
    ├── Registran clientes
    ├── Gestionan devoluciones (solo fallas → préstamo + envío a almacén)
    ├── Registran ácido consumido por venta
    ├── Compran chatarra a terceros → acumulan → envían a almacén
    ├── Mantienen caja chica + gastos
    └── Realizan cierre de caja diario

ADMIN (1+)
    ├── Ve todo de todas las sucursales (solo lectura)
    ├── Configura precios mínimos y duración de garantías
    ├── Gestiona usuarios, sucursales y proveedores
    └── Ve reportes consolidados
```

---

## 3. Usuarios y Roles

### 3.1 Administrador
- **Cantidad:** 1 o más
- **Acceso:** Global (todas las sucursales + almacén), solo lectura en operación de sucursales
- **Funciones:**
  - Crear/gestionar sucursales
  - Crear/gestionar usuarios (encargados de almacén, cajeros)
  - Configurar precios mínimos de venta por batería (con fecha efectiva opcional)
  - Configurar duración de garantías (secas vs. líquidas)
  - Ver reportes consolidados de todas las sucursales
  - Ver stock de todas las sucursales y almacén en tiempo real
  - Registrar proveedores
  - Definir tipos de chatarra y precios de compra
  - Autorizar inyecciones de efectivo a sucursales

### 3.2 Encargado de Almacén
- **Cantidad:** 1 o más (todos asignados al mismo almacén)
- **Acceso:** Almacén únicamente
- **Funciones:**
  - Crear/gestionar productos (baterías y accesorios)
  - Registrar lotes de compra a proveedores
  - Gestionar stock del almacén (pool central)
  - Comprar y gestionar ácido (stock propio + distribución a sucursales)
  - Distribuir stock a sucursales (en respuesta a pedidos o push directo)
  - Ver stock de todas las sucursales
  - Recibir baterías falladas de sucursales y enviar a fábrica
  - Recibir chatarra de sucursales
  - Vender a mayoristas (terceros)
  - Registrar gastos del almacén
  - Gestionar pedidos entrantes de sucursales (aprobar/rechazar)
  - Registrar reingreso de baterías de garantía al stock

### 3.3 Cajero / Encargado de Tienda
- **Cantidad:** 1 por sucursal
- **Acceso:** Solo su sucursal asignada
- **Funciones:**
  - Realizar ventas (baterías + accesorios)
  - Generar garantías y recibos en PDF
  - Registrar clientes
  - Gestionar devoluciones por falla (préstamo + envío a almacén)
  - Pedir transferencias de stock a almacén u otras sucursales
  - Aceptar/rechazar transferencias entrantes (con justificación)
  - Comprar chatarra a terceros
  - Registrar gastos de la sucursal
  - Realizar cierre de caja diario
  - Ver stock de otras sucursales y almacén
  - Registrar ajustes de inventario (merma, derrame de ácido)

---

## 4. Entidades del Dominio

### 4.1 Usuario
| Campo | Tipo | Descripción |
|---|---|---|
| id | UUID | Identificador único |
| username | string | Nombre de usuario (único) |
| password_hash | string | Contraseña hasheada (bcrypt) |
| role | enum | `admin`, `encargado_almacen`, `cajero` |
| branch_id | UUID? | Sucursal asignada (null para admin y encargado_almacen) |
| active | bool | Si la cuenta está activa |
| created_at | timestamp | Fecha de creación |

### 4.2 Sucursal
| Campo | Tipo | Descripción |
|---|---|---|
| id | UUID | Identificador único |
| name | string | Nombre de la sucursal |
| address | string | Dirección |
| petty_cash_balance | decimal | Saldo de caja chica |
| active | bool | Si la sucursal está activa |
| created_at | timestamp | Fecha de creación |

### 4.3 Producto — Batería
| Campo | Tipo | Descripción |
|---|---|---|
| id | UUID | Identificador único |
| brand | string | Marca |
| model | string | Modelo |
| voltage | decimal | Voltaje |
| amperage | decimal | Amperaje |
| battery_type | enum | `seca`, `liquida` |
| polarity | enum | `izquierda`, `derecha` |
| acid_liters | decimal? | Litros de ácido que admite (solo si es líquida) |
| min_sale_price | decimal | Precio mínimo de venta (configurado por admin) |
| effective_date | timestamp? | Fecha efectiva del precio (para cambios programados) |
| previous_price | decimal? | Precio anterior (cuando hay cambio programado) |
| vehicle_type | enum | `auto`, `moto`, `otro` |
| active | bool | Si el producto está activo |
| created_at | timestamp | Fecha de creación |
| created_by | UUID | Usuario que lo creó |

### 4.4 Producto — Accesorio
| Campo | Tipo | Descripción |
|---|---|---|
| id | UUID | Identificador único |
| name | string | Nombre (ej: bornera, agua destilada) |
| description | string? | Descripción opcional |
| active | bool | Si el producto está activo |
| created_at | timestamp | Fecha de creación |

### 4.5 Stock
| Campo | Tipo | Descripción |
|---|---|---|
| id | UUID | Identificador único |
| product_id | UUID | Producto (batería o accesorio) |
| product_type | enum | `bateria`, `accesorio` |
| location_type | enum | `almacen`, `sucursal` |
| location_id | UUID | ID del almacén o sucursal |
| quantity | int | Cantidad disponible |
| min_stock_alert | int | Umbral de alerta de stock bajo |
| updated_at | timestamp | Última actualización |

### 4.6 Stock de Ácido
| Campo | Tipo | Descripción |
|---|---|---|
| id | UUID | Identificador único |
| location_type | enum | `almacen`, `sucursal` |
| location_id | UUID | ID del almacén o sucursal |
| liters_available | decimal | Litros disponibles (puede ser negativo por derrames/mermas) |
| needs_adjustment | bool | Si el stock está en negativo y requiere ajuste |
| updated_at | timestamp | Última actualización |

### 4.7 Kardex de Ácido
| Campo | Tipo | Descripción |
|---|---|---|
| id | UUID | Identificador único |
| location_type | enum | `almacen`, `sucursal` |
| location_id | UUID | ID del punto |
| movement_type | enum | `entrada`, `salida`, `transferencia_envio`, `transferencia_recepcion`, `ajuste` |
| liters | decimal | Litros del movimiento |
| reference_type | enum | `compra`, `venta`, `transferencia`, `ajuste` |
| reference_id | UUID? | ID de la venta, transferencia o ajuste relacionado |
| notes | string? | Notas adicionales |
| created_at | timestamp | Fecha del movimiento |
| created_by | UUID | Usuario que registró |

### 4.8 Lote de Compra (Almacén)
| Campo | Tipo | Descripción |
|---|---|---|
| id | UUID | Identificador único |
| supplier_id | UUID | Proveedor |
| purchase_date | timestamp | Fecha de compra |
| notes | string? | Notas |
| created_by | UUID | Usuario que registró |

### 4.9 Detalle de Lote
| Campo | Tipo | Descripción |
|---|---|---|
| id | UUID | Identificador único |
| batch_id | UUID | Lote de compra |
| product_id | UUID | Producto |
| product_type | enum | `bateria`, `accesorio` |
| quantity | int | Cantidad comprada |
| unit_cost | decimal | Costo unitario (solo informativo, no visible para cajeros) |

### 4.10 Proveedor
| Campo | Tipo | Descripción |
|---|---|---|
| id | UUID | Identificador único |
| name | string | Nombre del proveedor |
| contact | string? | Teléfono/email |
| address | string? | Dirección |
| active | bool | Si está activo |
| created_at | timestamp | Fecha de registro |

### 4.11 Cliente
| Campo | Tipo | Descripción |
|---|---|---|
| id | UUID | Identificador único |
| name | string | Nombre completo |
| document | string | CI/DNI |
| phone | string | Teléfono |
| address | string? | Dirección |
| registration_branch_id | UUID | Sucursal donde se registró (válido en todas) |
| created_at | timestamp | Fecha de registro |

### 4.12 Venta
| Campo | Tipo | Descripción |
|---|---|---|
| id | UUID | Identificador único |
| branch_id | UUID | Sucursal |
| cashier_id | UUID | Cajero que realizó la venta |
| customer_id | UUID? | Cliente (opcional) |
| payment_method | enum | `qr`, `efectivo` |
| total | decimal | Total de la venta |
| total_profit | decimal | Ganancia total (precio_venta - precio_minimo por cada item) |
| warranty_generated | bool | Si se generó garantía |
| local_timestamp | timestamp | Hora del reloj local de la laptop al momento de la venta |
| synced_at | timestamp? | Hora de sincronización con el backend |
| price_discrepancy | bool | Si la venta se hizo con precio desactualizado |
| created_at | timestamp | Fecha/hora de la venta |

### 4.13 Detalle de Venta
| Campo | Tipo | Descripción |
|---|---|---|
| id | UUID | Identificador único |
| sale_id | UUID | Venta |
| product_id | UUID | Producto |
| product_type | enum | `bateria`, `accesorio` |
| quantity | int | Cantidad vendida |
| sale_price | decimal | Precio de venta unitario |
| min_price_at_sale | decimal | Precio mínimo vigente al momento de la venta (local) |
| min_price_at_sync | decimal? | Precio mínimo vigente al momento de sincronizar |
| local_timestamp | timestamp | Hora local del registro |

### 4.14 Garantía
| Campo | Tipo | Descripción |
|---|---|---|
| id | UUID | Identificador único |
| sale_id | UUID | Venta original |
| customer_id | UUID | Cliente |
| battery_id | UUID | Batería vendida (referencia al producto) |
| status | enum | `pendiente`, `enviada_almacen`, `enviada_fabrica`, `resuelta`, `cerrada` |
| resolution_type | enum? | `reingresa_stock`, `reemplazo_fabrica` (solo al resolver) |
| loan_battery_id | UUID? | Batería de préstamo asignada |
| notes | string? | Notas adicionales |
| warranty_duration_months | int | Duración configurada al momento de la venta |
| created_at | timestamp | Fecha de apertura |
| resolved_at | timestamp? | Fecha de resolución |
| closed_at | timestamp? | Fecha de cierre |

### 4.15 Historial de Garantía
| Campo | Tipo | Descripción |
|---|---|---|
| id | UUID | Identificador único |
| warranty_id | UUID | Garantía |
| previous_status | enum? | Estado anterior |
| new_status | enum | Nuevo estado |
| notes | string? | Notas del cambio |
| changed_by | UUID | Usuario que hizo el cambio |
| changed_at | timestamp | Fecha del cambio |

### 4.16 Batería de Préstamo
| Campo | Tipo | Descripción |
|---|---|---|
| id | UUID | Identificador único |
| branch_id | UUID | Sucursal donde está el pool |
| model | string | Modelo de la batería |
| serial | string? | Número de serie si tiene |
| status | enum | `disponible`, `prestada`, `dañada` |
| assigned_to_warranty_id | UUID? | Garantía asociada (si está prestada) |
| assigned_to_customer_id | UUID? | Cliente que la tiene |
| assigned_at | timestamp? | Fecha de préstamo |
| returned_at | timestamp? | Fecha de devolución |

### 4.17 Tipo de Chatarra
| Campo | Tipo | Descripción |
|---|---|---|
| id | UUID | Identificador único |
| name | string | Nombre (cobre, aluminio, radiador, batería usada, etc.) |
| purchase_price | decimal | Precio de compra por unidad/kg |
| unit_of_measure | enum | `unidad`, `kg` |
| active | bool | Si está activo |
| updated_at | timestamp | Última actualización de precio |

### 4.18 Compra de Chatarra
| Campo | Tipo | Descripción |
|---|---|---|
| id | UUID | Identificador único |
| branch_id | UUID | Sucursal que compró |
| scrap_type_id | UUID | Tipo de chatarra |
| quantity | decimal | Cantidad (unidades o kg) |
| price_paid | decimal | Total pagado |
| created_at | timestamp | Fecha de compra |
| created_by | UUID | Cajero que registró |

### 4.19 Transferencia
| Campo | Tipo | Descripción |
|---|---|---|
| id | UUID | Identificador único |
| origin_type | enum | `almacen`, `sucursal` |
| origin_id | UUID | Punto de origen |
| destination_type | enum | `almacen`, `sucursal` |
| destination_id | UUID | Punto de destino |
| status | enum | `pendiente`, `aprobada`, `rechazada`, `enviada`, `recibida` |
| rejection_reason | string? | Justificación si fue rechazada |
| transfer_type | enum | `stock`, `chatarra`, `bateria_fallada`, `acido` |
| created_by | UUID | Usuario que solicitó |
| created_at | timestamp | Fecha de solicitud |
| approved_by | UUID? | Usuario que aprobó/rechazó |
| approved_at | timestamp? | Fecha de aprobación/rechazo |
| shipped_at | timestamp? | Fecha de envío |
| received_at | timestamp? | Fecha de recepción |
| received_by | UUID? | Usuario que confirmó recepción |

### 4.20 Detalle de Transferencia
| Campo | Tipo | Descripción |
|---|---|---|
| id | UUID | Identificador único |
| transfer_id | UUID | Transferencia |
| product_id | UUID? | Producto transferido (si es stock) |
| product_type | enum | `bateria`, `accesorio`, `chatarra`, `acido` |
| quantity | decimal | Cantidad |
| liters | decimal? | Litros (solo si es ácido) |

### 4.21 Gasto
| Campo | Tipo | Descripción |
|---|---|---|
| id | UUID | Identificador único |
| location_type | enum | `almacen`, `sucursal` |
| location_id | UUID | Punto |
| concept | string | Concepto del gasto |
| amount | decimal | Monto |
| receipt_url | string? | URL del recibo PDF adjunto |
| created_at | timestamp | Fecha del gasto |
| created_by | UUID | Usuario que registró |

### 4.22 Cierre de Caja
| Campo | Tipo | Descripción |
|---|---|---|
| id | UUID | Identificador único |
| branch_id | UUID | Sucursal |
| closed_by | UUID | Cajero que cerró |
| close_date | date | Fecha del cierre |
| total_sales | decimal | Total ventas del día |
| total_expenses | decimal | Total gastos del día |
| total_scrap_purchased | decimal | Total gastado en chatarra |
| total_cash_movements | decimal | Neto de inyecciones/retiros de caja |
| opening_balance | decimal | Saldo inicial (cierre anterior) |
| closing_balance | decimal | Saldo final |
| created_at | timestamp | Hora del cierre |

### 4.23 Venta a Mayorista
| Campo | Tipo | Descripción |
|---|---|---|
| id | UUID | Identificador único |
| customer_name | string | Nombre de la tienda/cliente mayorista |
| items | relación | Productos vendidos |
| total | decimal | Total de la venta |
| notes | string? | Notas |
| created_by | UUID | Encargado que registró |
| created_at | timestamp | Fecha de venta |

### 4.24 Notificación
| Campo | Tipo | Descripción |
|---|---|---|
| id | UUID | Identificador único |
| user_id | UUID | Usuario destinatario |
| type | enum | `stock_bajo`, `stock_negativo_acido`, `cambio_estado_garantia`, `transferencia_rechazada`, `transferencia_recibida`, `precio_desactualizado` |
| title | string | Título corto |
| message | string | Mensaje completo |
| reference_type | string? | Tipo de entidad relacionada |
| reference_id | UUID? | ID de la entidad relacionada |
| read | bool | Si fue leída |
| created_at | timestamp | Fecha de creación |

### 4.25 Mantenimiento (Post-venta)
| Campo | Tipo | Descripción |
|---|---|---|
| id | UUID | Identificador único |
| customer_id | UUID | Cliente |
| warranty_id | UUID? | Garantía asociada (si aplica) |
| branch_id | UUID | Sucursal donde se hizo |
| notes | string | Notas del mantenimiento (relleno de líquido, etc.) |
| next_maintenance_date | date? | Fecha sugerida del próximo mantenimiento |
| created_by | UUID | Cajero que registró |
| created_at | timestamp | Fecha del mantenimiento |

### 4.26 Ajuste de Inventario
| Campo | Tipo | Descripción |
|---|---|---|
| id | UUID | Identificador único |
| location_type | enum | `almacen`, `sucursal` |
| location_id | UUID | Punto |
| product_type | enum | `acido`, `bateria`, `accesorio` |
| adjustment_type | enum | `merma`, `derrame`, `recalibracion`, `otro` |
| quantity_before | decimal | Stock antes del ajuste |
| quantity_after | decimal | Stock después del ajuste |
| reason | string | Motivo del ajuste |
| created_by | UUID | Usuario que registró |
| created_at | timestamp | Fecha del ajuste |

### 4.27 Movimiento de Caja
| Campo | Tipo | Descripción |
|---|---|---|
| id | UUID | Identificador único |
| branch_id | UUID | Sucursal |
| movement_type | enum | `ingreso_manual`, `retiro_manual` |
| amount | decimal | Monto |
| concept | string | Motivo del movimiento |
| authorized_by | UUID | Admin o encargado que autorizó |
| created_at | timestamp | Fecha del movimiento |

---

## 5. Flujos Principales

### 5.1 Venta (POS)

```
┌─────────────┐     ┌──────────────┐     ┌────────────────┐     ┌───────────────┐
│ Cajero abre │────▶│ Agrega items │────▶│ Asigna precio  │────▶│ Elige método  │
│   venta     │     │ (batería/    │     │ (≥ precio      │     │ de pago       │
│             │     │  accesorio)  │     │  mínimo local) │     │ (QR/efectivo) │
└─────────────┘     └──────────────┘     └────────────────┘     └───────┬───────┘
                                                                        │
                                                                        ▼
┌──────────────────┐    ┌──────────────────┐    ┌─────────────────────────────────┐
│ Se descuenta del │◀───│ Genera PDF:      │◀───│ ¿Batería de auto?               │
│ stock LOCAL de   │    │ - Garantía (auto)│    │   SÍ → opcional garantía        │
│ la sucursal      │    │ - Recibo (moto)  │    │   NO → recibo simple            │
└────────┬─────────┘    └──────────────────┘    └─────────────────────────────────┘
         │
         ▼
┌──────────────────────────┐    ┌──────────────────────────┐    ┌────────────────────┐
│ Si batería con líquido:  │───▶│ Se suma a caja local     │───▶│ Guarda en RxDB     │
│ se descuenta ácido       │    │ (si no es cheque)        │    │ estado:            │
│ localmente               │    │                          │    │ "pending_sync"     │
└──────────────────────────┘    └──────────────────────────┘    └──────────┬─────────┘
                                                                           │
                                                                           ▼
                                                              ┌────────────────────┐
                                                              │ Service Worker     │
                                                              │ sincroniza con     │
                                                              │ backend cuando     │
                                                              │ hay internet       │
                                                              └────────────────────┘
```

**Reglas de venta:**
- El precio de venta DEBE ser ≥ precio mínimo vigente según BD local
- Si el effective_date del precio ya pasó, se usa el precio nuevo
- Solo baterías de AUTO generan garantía (opcional, requiere cliente registrado)
- Baterías de MOTO generan recibo simple
- Si método de pago es QR o efectivo → se suma a la caja
- Si batería con líquido → se descuenta ácido según litros configurados en el producto
- El cliente puede ser opcional, pero es REQUERIDO si se genera garantía
- Se registra local_timestamp (reloj de la laptop) para auditoría forense

### 5.2 Garantía

```
ESTADOS: pendiente → enviada_almacen → enviada_fabrica → resuelta → cerrada
```

```
┌────────────────────┐     ┌──────────────────────┐     ┌──────────────────────┐
│ Cliente trae       │────▶│ Cajero registra      │────▶│ Se asigna batería    │
│ batería fallada    │    │ garantía (estado:    │    │ de préstamo al       │
│ (con recibo)       │    │ pendiente)           │    │ cliente              │
└────────────────────┘    └──────────┬───────────┘    └──────────────────────┘
                                      │
                                      ▼
┌────────────────────┐     ┌──────────────────────┐     ┌──────────────────────┐
│ Almacén cambia     │◀────│ Transferencia        │────▶│ Cajero envía batería │
│ estado a           │    │ especial: batería    │    │ fallada a almacén    │
│ "enviada_almacen"  │    │ fallada              │    │                      │
└────────┬───────────┘    └──────────────────────┘    └──────────────────────┘
         │
         ▼
┌────────────────────┐     ┌──────────────────────┐     ┌──────────────────────┐
│ Fábrica resuelve   │◀────│ Almacén registra     │────▶│ Almacén envía a      │
│ (fuera del sistema,│    │ envío a fábrica      │    │ fábrica (se registra │
│ se registra fecha) │    │ (estado:             │    │ fecha de envío)      │
│                    │    │ "enviada_fabrica")   │    │                      │
└────────┬───────────┘    └──────────────────────┘    └──────────────────────┘
         │
         ▼
┌────────────────────┐     ┌──────────────────────┐     ┌──────────────────────┐
│ Almacén marca      │────▶│ ¿Cómo reingresa?     │────▶│ Almacén envía        │
│ garantía como      │    │                      │    │ batería a sucursal   │
│ "resuelta"         │    │ A: Reingresa a stock │    │                      │
│                    │    │ B: Reemplazo fábrica │    │                      │
└────────┬───────────┘    └──────────────────────┘    └──────────┬───────────┘
         │                                                       │
         │ "reingresa_stock"                                     ▼
         │                                          ┌──────────────────────┐
         └──▶ Stock del almacén se incrementa       │ Cliente devuelve    │
             automáticamente con esa batería         │ batería préstamo    │
                                                    │ (vuelve al pool)    │
                                                    └──────────┬───────────┘
                                                               │
                                                               ▼
                                                    ┌──────────────────────┐
                                                    │ Garantía cerrada     │
                                                    └──────────────────────┘
```

**Reglas de garantía:**
- Solo se aplica a baterías de AUTO
- El cliente DEBE tener el recibo original
- Se presta una batería del pool de préstamo de la sucursal
- Cada cambio de estado genera notificación al cajero de la sucursal
- La garantía es válida en CUALQUIER sucursal
- Al resolver: el encargado de almacén indica si la batería reingresa al stock o es un reemplazo de fábrica
- **Reingreso automático:** cuando se resuelve una garantía con `resolution_type = "reingresa_stock"`, el stock del almacén se incrementa automáticamente con esa batería

### 5.3 Transferencia de Stock

```
OPCIÓN A: Sucursal solicita
┌────────────────────┐     ┌──────────────────────┐     ┌──────────────────────┐
│ Sucursal A crea    │────▶│ Sucursal B recibe    │────▶│ ACEPTA: genera envío │
│ solicitud de       │    │ notificación         │    │ descuenta stock      │
│ transferencia      │    │                      │    │                      │
└────────────────────┘    └──────────┬───────────┘    └──────────┬───────────┘
                                      │                           │
                                      │ OPCIÓN B:                 ▼
                                      │ RECHAZA              ┌─────────────────┐
                                      ▼  (con justificación) │ Sucursal A      │
                                 ┌──────────────────┐        │ recibe y        │
                                 │ Notifica a A con │        │ confirma        │
                                 │ motivo del       │        │ recepción       │
                                 │ rechazo          │        └─────────────────┘
                                 └──────────────────┘

OPCIÓN C: Almacén push directo
┌────────────────────┐     ┌──────────────────────┐
│ Almacén ve stock   │────▶│ Envía directamente   │
│ bajo en sucursal y │    │ (genera transferencia │
│ decide enviar      │    │  aprobada de una)     │
└────────────────────┘    └──────────────────────┘
```

**Nota:** Las transferencias requieren internet (operación entre puntos).

### 5.4 Compra a Proveedor (Almacén)

```
┌────────────────────┐     ┌──────────────────────┐     ┌──────────────────────┐
│ Encargado          │────▶│ Registra lote de     │────▶│ Productos entran al  │
│ registra proveedor │    │ compra (items,       │    │ pool de stock del     │
│ (si es nuevo)      │    │ cantidades, costo)   │    │ almacén              │
└────────────────────┘    └──────────────────────┘    └──────────────────────┘
```

### 5.5 Gestión de Ácido

```
COMPRA (separada de baterías):
┌────────────────────┐     ┌──────────────────────┐
│ Almacén compra     │────▶│ Ingresa a stock de   │
│ ácido a proveedor  │    │ ácido del almacén    │
└────────────────────┘    └──────────┬───────────┘
                                      │
DISTRIBUCIÓN:                        ▼
┌────────────────────┐     ┌──────────────────────┐     ┌──────────────────────┐
│ Almacén asigna     │────▶│ Sucursal recibe      │────▶│ Se registra en       │
│ litros a sucursal  │    │ ácido                │    │ kardex de ambas      │
│ (transferencia)    │    │                      │    │ partes               │
└────────────────────┘    └──────────────────────┘    └──────────────────────┘

CONSUMO (automático en venta):
┌────────────────────┐     ┌──────────────────────┐
│ Se vende batería   │────▶│ Se descuentan litros │
│ con líquido        │    │ del stock de ácido   │
│                    │    │ de la sucursal       │
└────────────────────┘    └──────────┬───────────┘
                                      │
                                      ▼
                           ┌──────────────────────┐
                           │ ¿Stock quedó negativo?│
                           │ SÍ → notificación +  │
                           │      needs_adjustment│
                           │ NO → todo OK         │
                           └──────────────────────┘
```

**Regla de stock negativo:** El kardex de ácido PERMITE stock negativo (por derrames/mermas). Al llegar a negativo se genera alerta automática y se marca `needs_adjustment = true`. El encargado puede registrar un ajuste de inventario para corregir.

### 5.6 Chatarra

```
┌────────────────────┐     ┌──────────────────────┐     ┌──────────────────────┐
│ Almacén define     │────▶│ Sucursal ve lista    │────▶│ Sucursal compra      │
│ tipos de chatarra  │    │ actualizada con      │    │ chatarra a terceros  │
│ + precios          │    │ precios              │    │ (descuenta de caja)  │
└────────────────────┘    └──────────────────────┘    └──────────┬───────────┘
                                                                  │
                                                                  ▼
┌────────────────────┐     ┌──────────────────────┐     ┌──────────────────────┐
│ Almacén vende a    │◀────│ Almacén recibe       │◀────│ Sucursal acumula y   │
│ recicladora        │    │ chatarra             │    │ envía a almacén      │
│ (fuera del sistema)│    │                      │    │ (transferencia)      │
└────────────────────┘    └──────────────────────┘    └──────────────────────┘
```

**Reglas de chatarra:**
- Almacén define tipos (cobre, aluminio, radiador, batería usada, etc.) y precios de compra
- Sucursal registra: tipo, cantidad, precio pagado (NO se registra a quién se compró)
- Se descuenta de la caja chica de la sucursal
- La venta a recicladoras se registra fuera del sistema por ahora

### 5.7 Cierre de Caja

```
┌────────────────────┐     ┌──────────────────────┐     ┌──────────────────────┐
│ Cajero presiona    │────▶│ Sistema muestra      │────▶│ Cajero revisa y      │
│ "Cerrar caja"      │    │ resumen del día:     │    │ confirma             │
│ (al final del día) │    │ - Total ventas       │    │                      │
│                    │    │ - Total gastos       │    │                      │
│                    │    │ - Total chatarra     │    │                      │
│                    │    │ - Inyecciones/retiros│    │                      │
│                    │    │ - Saldo de caja      │    │                      │
└────────────────────┘    └──────────────────────┘    └──────────┬───────────┘
                                                                  │
                                                                  ▼
                                                       ┌──────────────────────┐
                                                       │ Se genera registro   │
                                                       │ de cierre            │
                                                       │ Admin puede ver      │
                                                       │ todos los cierres    │
                                                       └──────────────────────┘
```

### 5.8 Gastos

```
Sucursal o Almacén registra gasto:
┌────────────────────┐     ┌──────────────────────┐     ┌──────────────────────┐
│ Concepto del gasto │────▶│ Monto                │────▶│ Adjuntar recibo PDF  │
│ (cuadernos, trapos,│    │                      │    │ (opcional)           │
│  limpieza, etc.)   │    │                      │    │                      │
└────────────────────┘    └──────────────────────┘    └──────────┬───────────┘
                                                                  │
                                                                  ▼
                                                       ┌──────────────────────┐
                                                       │ Se descuenta de      │
                                                       │ caja chica del punto │
                                                       └──────────────────────┘
```

### 5.9 Rendición de Cuentas (Automática)

```
┌────────────────────────────────────────────────────────────────────────────┐
│                    VISTA DEL ADMINISTRADOR                                │
│                                                                            │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │ Ventas por   │  │ Stock de     │  │ Cierres de   │  │ Garantías    │  │
│  │ sucursal/    │  │ todas las    │  │ caja de      │  │ activas/     │  │
│  │ periodo      │  │ sucursales + │  │ todas las    │  │ pendientes   │  │
│  │              │  │ almacén      │  │ sucursales   │  │              │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │ Comisiones   │  │ Ácido        │  │ Chatarra     │  │ Gastos por   │  │
│  │ por cajero   │  │ kardex       │  │ acumulada/   │  │ punto        │  │
│  │              │  │              │  │ vendida      │  │              │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘  │
│  ┌──────────────┐  ┌──────────────┐                                       │
│  │ Ventas con   │  │ Inyecciones/ │                                       │
│  │ precio       │  │ retiros de   │                                       │
│  │ desactualiz. │  │ caja         │                                       │
│  └──────────────┘  └──────────────┘                                       │
└────────────────────────────────────────────────────────────────────────────┘
```

### 5.10 Ajuste de Inventario

```
Se detecta merma/derrame/recalibración:
┌────────────────────┐     ┌──────────────────────┐     ┌──────────────────────┐
│ Encargado/Cajero   │────▶│ Registra tipo de     │────▶│ Stock se ajusta      │
│ selecciona producto│    │ ajuste + motivo      │    │ (positivo o negativo)│
│ a ajustar          │    │                      │    │                      │
└────────────────────┘    └──────────────────────┘    └──────────┬───────────┘
                                                                  │
                                                                  ▼
                                                       ┌──────────────────────┐
                                                       │ Se registra en       │
                                                       │ kardex (si es ácido) │
                                                       │ y en ajustes de      │
                                                       │ inventario           │
                                                       └──────────────────────┘
```

### 5.11 Movimiento de Caja (Inyección/Retiro)

```
Sucursal necesita efectivo:
┌────────────────────┐     ┌──────────────────────┐     ┌──────────────────────┐
│ Cajero solicita    │────▶│ Admin/Encargado      │────▶│ Se registra          │
│ inyección de       │    │ autoriza             │    │ movimiento de caja   │
│ efectivo           │    │                      │    │                      │
└────────────────────┘    └──────────────────────┘    └──────────┬───────────┘
                                                                  │
                                                                  ▼
                                                       ┌──────────────────────┐
                                                       │ Saldo de caja chica  │
                                                       │ se actualiza         │
                                                       │ (+ para ingreso,     │
                                                       │  - para retiro)      │
                                                       └──────────────────────┘
```

---

## 6. Notificaciones

| Disparador | Destinatario | Mensaje |
|---|---|---|
| Stock bajo producto X en sucursal | Cajero de esa sucursal + Almacén | "Stock bajo de [producto] en [sucursal]: [cantidad] unidades" |
| Stock bajo producto X en almacén | Encargado de almacén | "Stock bajo de [producto] en almacén: [cantidad] unidades" |
| Stock negativo de ácido | Cajero/Encargado del punto | "Stock de ácido en negativo en [punto]: [litros] litros. Se requiere ajuste." |
| Cambio de estado en garantía | Cajero de la sucursal del cliente | "Garantía [id] cambió a estado: [nuevo_estado]" |
| Transferencia rechazada | Sucursal solicitante | "Transferencia [id] rechazada: [justificación]" |
| Transferencia recibida | Sucursal de origen | "Transferencia [id] fue recibida por [sucursal destino]" |
| Nueva transferencia entrante | Sucursal destino | "Nueva transferencia pendiente de [origen]" |
| Venta con precio desactualizado | Admin | "Venta [id] en [sucursal] realizada con precio desactualizado: [precio_viejo] vs [precio_actual]" |
| Cambio de precio próximo | Sucursales afectadas | "Precio de [producto] cambiará a [nuevo_precio] el [fecha]" |

**Notificaciones:**
- En tiempo real vía WebSockets cuando hay internet
- Persisten en base de datos para historial
- Al reconectar: se cargan notificaciones pendientes vía REST

---

## 7. Generación de PDFs

| Documento | Contenido | Cuándo se genera | Modo |
|---|---|---|---|
| **Garantía (batería con líquido)** | Datos cliente, datos batería, fecha venta, duración garantía, sucursal | Al vender batería de auto con líquido + cliente registrado | Server (online) / Client (offline) |
| **Garantía (batería seca)** | Datos cliente, datos batería, fecha venta, duración garantía, sucursal | Al vender batería de auto seca + cliente registrado | Server (online) / Client (offline) |
| **Recibo (batería de moto)** | Datos venta, batería, precio, sucursal, cajero | Al vender batería de moto | Server (online) / Client (offline) |
| **Recibo (venta general)** | Datos venta, items, precios, método de pago | Para cualquier venta sin garantía | Server (online) / Client (offline) |
| **Cierre de caja** | Resumen del día: ventas, gastos, chatarra, saldos | Al cerrar caja | Server (online) / Client (offline) |

**Modo offline:** Los PDFs se generan client-side usando jsPDF cuando no hay conexión. El diseño es idéntico al generado por el servidor.
## 8. Stack Tecnológico

| Capa | Tecnología | Justificación |
|---|---|---|
| **Backend** | Go | Elegido por el desarrollador, alto rendimiento para APIs |
| **Frontend** | Vite + React SPA + TypeScript | Ligero, control total del build, ideal para PWA offline-first |
| **Estado UI** | Zustand | Estado temporal (carrito, forms, UI) |
| **BD Local (Offline)** | RxDB con plugin de encriptación AES-256 | IndexedDB wrapper con sync automático, queries tipo SQL, encriptación por clave derivada del password del usuario |
| **Sincronización** | RxDB replication + Service Worker | Background sync, cola de operaciones pendientes |
| **BD Servidor** | PostgreSQL | Fuente de verdad, transacciones fuertes |
| **Auth** | JWT (access + refresh tokens) + bcrypt hash local para offline | Stateless online + desbloqueo offline seguro |
| **Notificaciones** | WebSockets + REST fallback | Tiempo real cuando hay internet, historial vía REST al reconectar |
| **PDFs** | Go server-side (online) / jsPDF client-side (offline) | Server cuando hay internet, local cuando no |
| **PWA** | manifest.json + Service Worker | Se instala como app en escritorio/celular |
| **API** | REST | Para MVP, migrar a GraphQL si se necesita |
| **Web Worker** | Para sync de prioridad 2 | Evita bloquear el Main Thread durante carga masiva |

---

## 9. Arquitectura Backend (Go)

```
├── cmd/
│   └── server/
│       └── main.go              # Entry point
├── internal/
│   ├── domain/                  # Entidades del dominio (sin dependencias)
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
│   │   ├── cash_movement.go     # NUEVO: movimientos de caja
│   │   ├── inventory_adjustment.go  # NUEVO: ajustes de inventario
│   │   ├── notification.go
│   │   └── maintenance.go
│   ├── ports/                   # Interfaces (contratos)
│   │   ├── repositories.go
│   │   ├── services.go
│   │   └── handlers.go
│   ├── application/             # Casos de uso / servicios
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
│   │   ├── cash_movement_service.go    # NUEVO
│   │   ├── inventory_adjustment_service.go  # NUEVO
│   │   ├── notification_service.go
│   │   ├── sync_service.go             # NUEVO: sincronización incremental
│   │   └── report_service.go
│   ├── infrastructure/
│   │   ├── postgres/
│   │   │   ├── migrations/
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
│   │   │   ├── cash_movement_repo.go   # NUEVO
│   │   │   ├── inventory_adjustment_repo.go  # NUEVO
│   │   │   ├── notification_repo.go
│   │   │   └── sync_repo.go            # NUEVO
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
│   │   │   │   ├── cash_movement_handler.go  # NUEVO
│   │   │   │   ├── inventory_adjustment_handler.go  # NUEVO
│   │   │   │   ├── sync_handler.go     # NUEVO
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
│   └── validator/
├── scripts/
│   └── generate-schemas/        # NUEVO: genera JSON Schema desde structs Go
│       └── main.go
├── migrations/
│   └── *.sql
├── go.mod
├── go.sum
└── Makefile
```

### Endpoint de sincronización incremental

```
GET /api/sync/changes?since={timestamp}&branch_id={uuid}

Response:
{
  "sync_timestamp": "2025-04-18T12:30:00Z",
  "products": [...],          // solo los que cambiaron
  "stock_changes": [...],     // movimientos desde last_sync
  "price_changes": [...],     // cambios de precio (incluye effective_date)
  "scrap_types": [...],       // tipos/precios actualizados
  "customers": [...],         // nuevos o modificados
  "config": {...},            // configuración actualizada (garantías, etc.)
  "deactivated_users": [...]  // usuarios desactivados (para invalidar login offline)
}
```

---

## 10. Arquitectura Frontend (Vite + React SPA)

```
├── public/
│   ├── manifest.json          # PWA manifest
│   └── sw.js                  # Service Worker (o generado por vite-plugin-pwa)
├── src/
│   ├── main.tsx
│   ├── App.tsx
│   ├── routes/
│   │   ├── index.tsx          # Router (React Router o TanStack Router)
│   │   ├── login.tsx
│   │   ├── login-offline.tsx  # NUEVO: pantalla de login offline
│   │   ├── pos.tsx            # Pantalla principal del POS
│   │   ├── dashboard.tsx      # Dashboard admin / reportes
│   │   ├── products.tsx
│   │   ├── stock.tsx
│   │   ├── sales.tsx
│   │   ├── warranties.tsx
│   │   ├── transfers.tsx
│   │   ├── acid.tsx
│   │   ├── scrap.tsx
│   │   ├── expenses.tsx
│   │   ├── cash-close.tsx
│   │   ├── customers.tsx
│   │   ├── suppliers.tsx
│   │   ├── branches.tsx
│   │   ├── users.tsx
│   │   ├── reports.tsx
│   │   └── settings.tsx
│   ├── components/
│   │   ├── ui/                # Componentes base (botones, modals, tables)
│   │   ├── layout/            # Sidebar, navbar, connection status indicator
│   │   ├── forms/             # Formularios reutilizables
│   │   ├── pos/               # Carrito, item selector, payment
│   │   ├── stock/             # Tablas de inventario
│   │   ├── warranty/          # Estados de garantía, timeline
│   │   └── reports/           # Charts, tablas de reportes
│   ├── db/
│   │   ├── database.ts        # Inicialización de RxDB
│   │   ├── collections.ts     # Definición de colecciones con JSON Schema
│   │   ├── replication.ts     # Configuración de sync con backend
│   │   └── encryption.ts      # Encriptación AES-256, derivación de clave
│   ├── schemas/               # NUEVO: JSON Schemas auto-generados desde Go
│   │   ├── product.schema.json
│   │   ├── sale.schema.json
│   │   ├── stock.schema.json
│   │   ├── warranty.schema.json
│   │   ├── customer.schema.json
│   │   └── ... (cada entidad tiene su schema)
│   ├── workers/
│   │   └── sync-worker.ts     # NUEVO: Web Worker para sync de prioridad 2
│   ├── lib/
│   │   ├── api.ts             # Cliente API (fetch wrapper)
│   │   ├── auth.ts            # Helpers de autenticación
│   │   ├── pdf.ts             # Generación de PDFs offline (jsPDF)
│   │   └── utils.ts           # Utilidades
│   ├── hooks/                 # Custom hooks
│   │   ├── use-auth.ts
│   │   ├── use-sync-status.ts # NUEVO: estado de sincronización
│   │   ├── use-offline.ts     # NUEVO: detectar modo offline
│   │   └── ...
│   ├── stores/                # Estado global (Zustand)
│   │   ├── auth-store.ts
│   │   ├── cart-store.ts      # Estado del carrito del POS
│   │   ├── sync-store.ts      # NUEVO: cola de sync, estado de conexión
│   │   └── ui-store.ts
│   └── types/                 # Tipos TypeScript (generados o manuales)
├── vite.config.ts
├── tailwind.config.ts
├── tsconfig.json
└── package.json
```

### Generación de esquemas (Go → JSON Schema → TypeScript)

```
FUENTE DE VERDAD: Backend en Go (structs)
│
├── scripts/generate-schemas/main.go
│   Lee structs de internal/domain/
│   Genera JSON Schema con versionado ("version": 0, 1, 2...)
│   Exporta a frontend/src/schemas/
│
▼
frontend/src/schemas/
├── product.schema.json        ← Auto-generado desde Go
│   { "version": 1, "title": "Product", ... }
│
▼
RxDB collections.ts
├── Usa los JSON Schemas para definir colecciones
├── Si hay nueva versión → ejecuta migración automática
└── Si schema no matchea → falla en BUILD (no silenciosamente)
```

### Mapeo de tipos Go ↔ JSON Schema ↔ TypeScript

| Go | JSON Schema | TypeScript |
|---|---|---|
| `string` | `"type": "string"` | `string` |
| `int` | `"type": "integer"` | `number` |
| `float64` | `"type": "number"` | `number` |
| `bool` | `"type": "boolean"` | `boolean` |
| `time.Time` | `"type": "string", "format": "date-time"` | `string` (ISO 8601) |
| `uuid.UUID` | `"type": "string", "format": "uuid"` | `string` |
| `*Type` (nullable) | `"type": ["string", "null"]` | `string \| null` |

### Arquitectura Offline-First

```
┌─────────────────────────────────────────────────────────────┐
│                    NAVEGADOR (PWA instalada)                 │
│                                                              │
│  ┌──────────┐    ┌──────────┐    ┌───────────────────────┐  │
│  │ Zustand  │    │   RxDB   │    │    Service Worker     │  │
│  │ (UI/carro│    │ (IndexedDB│    │  (sincronizador BG)   │  │
│  │  t temp) │    │ + AES-256)│    │                       │  │
│  └────┬─────┘    └────┬─────┘    └───────────┬───────────┘  │
│       │               │                      │              │
│       │          lectura/escritura       sync en BG         │
│       │          LOCAL (0ms)             cuando hay internet│
│       │               │                      │              │
└───────┼───────────────┼──────────────────────┼──────────────┘
        │               │                      │
        │               ▼                      ▼
        │    ┌──────────────────────────────────────────┐
        └───▶│         BACKEND EN GO (API REST)         │
             │                                          │
             │  ┌─────────────┐    ┌────────────────┐   │
             │  │ PostgreSQL  │    │  Cola de sync  │   │
             │  │ (fuente de  │    │  (procesa en   │   │
             │  │  verdad)    │    │   orden FIFO)  │   │
             │  └─────────────┘    └────────────────┘   │
             └──────────────────────────────────────────┘
```

### Autenticación Offline

```
PRIMERA VEZ (con internet):
┌─────────────────────────────────────────────────────────────┐
│ Cajero se loguea con username + password                    │
│ → Backend verifica y devuelve JWT                          │
│ → Frontend guarda en RxDB encriptada:                      │
│   ├── JWT (access + refresh)                               │
│   ├── hash bcrypt del password (NO el pass en texto)       │
│   ├── user_id, role, branch_id                             │
│   └── timestamp del último login exitoso                   │
│ → La clave de encriptación de RxDB se deriva del password  │
│ → RxDB replica: catálogo, stock, clientes, etc.            │
└─────────────────────────────────────────────────────────────┘

SIGUIENTES VECES (sin internet):
┌─────────────────────────────────────────────────────────────┐
│ Cajero abre la PWA                                          │
│ → Service Worker detecta que no hay internet                │
│ → Muestra pantalla de "Login offline"                       │
│ → Cajero ingresa username + password                        │
│ → Frontend compara bcrypt(pass_ingresado) con hash local    │
│ → Si coincide: desencripta RxDB y abre el POS              │
│ → JWT expirado no importa: opera en modo local              │
│ → Máximo 30 días offline antes de forzar re-validación     │
│ → Al revocar acceso de usuario: se borra DB local en sync  │
└─────────────────────────────────────────────────────────────┘
```

### Prioridades de sincronización

```
CARGA INICIAL (primera vez):
PRIORIDAD 1 (bloqueante — sin esto no abre POS):
├── Productos activos (baterías + accesorios)
├── Stock de ESTA sucursal
├── Stock de ácido de ESTA sucursal
├── Tipos de chatarra activos + precios
└── Configuración (duración garantía, etc.)

PRIORIDAD 2 (background vía Web Worker):
├── Clientes registrados
├── Stock de OTRAS sucursales (solo lectura)
├── Stock de almacén
└── Notificaciones pendientes

SYNC INCREMENTAL (las siguientes veces):
GET /api/sync/changes?since={last_sync_timestamp}
→ Solo lo que cambió desde la última sync
→ RxDB aplica los cambios incrementalmente
```

### Operaciones disponibles offline vs online

| Operación | ¿Funciona offline? | Nota |
|---|---|---|
| Ventas (POS) | ✅ Sí | Se sincroniza después |
| Ver stock local | ✅ Sí | Desde BD local |
| Registrar clientes | ✅ Sí | Se sincroniza después |
| Generar PDF | ✅ Sí | jsPDF client-side |
| Cierre de caja | ✅ Sí | Con datos locales |
| Gastos | ✅ Sí | Se sincroniza después |
| Compra de chatarra | ✅ Sí | Se sincroniza después |
| Ajuste de inventario | ✅ Sí | Se sincroniza después |
| Ver stock de OTRAS sucursales | ⚠️ Parcial | Muestra última sync |
| Transferencias entre sucursales | ❌ No | Requiere internet |
| Notificaciones en tiempo real | ❌ No | Se cargan al reconectar |
| Crear productos nuevos | ❌ No | Solo almacén, requiere internet |

### Manejo de cambio de precios offline

```
Estrategia: Fecha efectiva + Aceptación con alerta + Auditoría

1. Admin programa cambio de precio con effective_date
2. Sucursales reciben el cambio en sync incremental
3. La PWA aplica el precio nuevo localmente cuando la fecha llega
4. Si una venta se hizo con precio desactualizado:
   → Backend ACEPTA la venta (nunca impedir que entre dinero)
   → Marca price_discrepancy = true
   → Registra local_timestamp como prueba forense
   → Notifica al admin

DECISIONES DEL BACKEND AL SINCRONIZAR:
| Escenario | Decisión |
|---|---|
| Precio vigente al momento de la venta | ✅ Aceptada sin problema |
| Precio cambiado después (diff < 10%) | ✅ Aceptada con alerta |
| Precio cambiado después (diff > 10%) | ✅ Aceptada con alerta CRÍTICA |
| Producto desactivado | ❌ Rechazada, notifica al cajero |
```

---

## 11. MVP v1 — Alcance

### ✅ Incluido en MVP v1

| Módulo | Funcionalidad |
|---|---|
| **Auth** | Login online + offline, roles, gestión de usuarios por admin |
| **Sucursales** | CRUD por admin |
| **Productos** | CRUD de baterías (marca, modelo, voltaje, amperaje, tipo, polaridad, litros_acido, precio_mínimo con effective_date) + accesorios básicos |
| **Stock** | Vista de stock por punto, alertas de stock bajo |
| **Proveedores** | Registro y gestión |
| **Lotes de compra** | Registro de compras a proveedores (almacén) |
| **Ventas (POS)** | Flujo offline-first: agregar items, asignar precio, QR/efectivo, generar PDF |
| **Clientes** | Registro, búsqueda, asignación a ventas |
| **Garantías** | Flujo completo con tracking de estados + baterías de préstamo + reingreso automático |
| **Transferencias** | Pedidos entre puntos (requiere internet), aprobación/rechazo, push desde almacén |
| **Ácido** | Stock por punto (permite negativo), kardex completo, distribución, descuento automático, ajustes |
| **Chatarra** | Tipos + precios definidos por almacén, compras por sucursal, envío a almacén |
| **Gastos** | Registro por punto con adjunto de recibo PDF |
| **Caja** | Cierre de caja diario + inyecciones/retiros de efectivo |
| **Notificaciones** | Stock bajo, stock negativo ácido, cambios de garantía, transferencias, precios desactualizados |
| **Reportes** | Stock consolidado, ventas, cierres, garantías, kardex ácido, comisiones, precios desactualizados |
| **PDFs** | Garantías (seca/líquida), recibos, cierre de caja (server online / jsPDF offline) |
| **Mantenimiento** | Registro de mantenimientos por cliente |
| **Offline-first** | PWA completa con RxDB encriptada, sync automático, login offline |
| **Ajustes de inventario** | Registro de mermas, derrames, recalibraciones |
| **Generación de schemas** | Script Go → JSON Schema para mantener backend/frontend sincronizados |

### ❌ Fuera del MVP v1 (para versiones futuras)

| Funcionalidad | Razón |
|---|---|
| Facturación electrónica | Requiere integración externa, no es prioridad |
| Cheques | Caso complejo, se puede agregar después |
| App móvil nativa | PWA instalable es suficiente por ahora |
| Descuentos / promociones | No es parte del negocio actual |
| Comisiones automáticas | Se calculan manualmente por ahora |
| Contabilidad integrada | Los reportes básicos cubren la necesidad |
| Venta a recicladoras (registro) | Fuera del sistema por ahora |
| Escáner de código de barras | No se necesita aún |

---

## 12. Orden de Implementación Sugerido

| Fase | Módulos | Dependencias |
|---|---|---|
| **Fase 0: Fundación Técnica** | Setup Go, PostgreSQL, Vite+React, RxDB, Service Worker, PWA, generador de schemas | Ninguna |
| **Fase 1: Auth + Sucursales** | Auth (online + offline), Sucursales, Usuarios, Proveedores | Fase 0 |
| **Fase 2: Productos + Stock** | CRUD productos, Stock por punto, Lotes de compra, Transferencias básicas, Alertas | Fase 1 |
| **Fase 3: POS** | Ventas offline-first, Clientes, PDFs (recibos/garantías), Cierre de caja, Movimientos de caja | Fase 2 |
| **Fase 4: Garantías** | Garantías, Baterías de préstamo, Historial de garantías, Reingreso automático | Fase 3 |
| **Fase 5: Ácido** | Stock de ácido (permite negativo), Kardex, Distribución, Descuento automático, Ajustes | Fase 2 |
| **Fase 6: Chatarra** | Tipos de chatarra, Compras por sucursal, Transferencias a almacén | Fase 2 |
| **Fase 7: Operaciones** | Gastos, Notificaciones (WS + REST), Reportes, Mantenimiento | Fase 3 |

---

## 13. Criterios de Aceptación del MVP

- [ ] Un admin puede crear sucursales y usuarios
- [ ] Un admin puede configurar precios mínimos con fecha efectiva y duración de garantías
- [ ] Un encargado de almacén puede registrar productos, proveedores y lotes de compra
- [ ] Un encargado de almacén puede distribuir stock a sucursales
- [ ] Un cajero puede realizar una venta completa SIN internet (offline-first)
- [ ] Un cajero puede loguearse sin internet usando hash local
- [ ] Las ventas offline se sincronizan automáticamente al reconectar
- [ ] Una venta con precio desactualizado se acepta y se marca con alerta
- [ ] Un cajero puede generar garantía en PDF para baterías de auto (online y offline)
- [ ] Un cajero puede generar recibo en PDF para baterías de moto (online y offline)
- [ ] El flujo de garantía funciona completo (pendiente → cerrada) con tracking
- [ ] El reingreso de batería al stock funciona al resolver garantía
- [ ] Las transferencias entre puntos funcionan con aprobación/rechazo
- [ ] El ácido se descuenta automáticamente al vender baterías con líquido
- [ ] El stock de ácido puede quedar en negativo con alerta
- [ ] Los ajustes de inventario (merma/derrame) se registran correctamente
- [ ] Las inyecciones/retiros de caja se registran y afectan el saldo
- [ ] La chatarra se puede comprar, acumular y enviar a almacén
- [ ] Los gastos se registran con recibo PDF adjunto
- [ ] El cierre de caja genera un resumen correcto del día
- [ ] Las notificaciones se envían por stock bajo, cambios de garantía y precios desactualizados
- [ ] El admin puede ver reportes consolidados de todas las sucursales
- [ ] La PWA se instala en el escritorio como aplicación
- [ ] La BD local está encriptada con AES-256
- [ ] Los JSON Schemas se generan automáticamente desde Go
- [ ] La sync de prioridad 2 usa Web Worker (no bloquea UI)

---

## 14. Preguntas Abiertas (por definir durante desarrollo)

| Tema | Nota |
|---|---|
| Valores iniciales | Litros de ácido por defecto por batería, precios mínimos iniciales |
| Tipos de chatarra iniciales | Lista inicial (cobre, aluminio, radiador, batería usada, etc.) |
| Duración garantía default | Meses para secas vs. líquidas (configurable por admin) |
| Límite stock bajo | A partir de cuántas unidades se dispara la notificación |
| Unidad de medida chatarra | Definido por tipo de chatarra (unidad o kg) |
| Tiempo máximo offline | 30 días antes de forzar re-validación online |
| Tolerancia de precio desactualizado | 10% como umbral para alerta crítica (¿ajustar?) |

---
