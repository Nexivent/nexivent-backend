# API de Eventos - Nexivent Backend

## Descripción General
API REST para la gestión de eventos en la plataforma Nexivent.

Base URL: `/v1/eventos`

---

## Endpoints Disponibles

### 1. Obtener un Evento por ID

**Método:** `GET`  
**Ruta:** `/v1/eventos/:id`

#### Parámetros
- `id` (UUID, requerido): ID del evento a consultar

#### Respuesta Exitosa (200 OK)
```json
{
  "evento": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "organizadorId": "223e4567-e89b-12d3-a456-426614174001",
    "categoriaId": "323e4567-e89b-12d3-a456-426614174002",
    "titulo": "Concierto de Rock",
    "descripcion": "Un increíble concierto de rock con las mejores bandas locales",
    "lugar": "Estadio Municipal",
    "eventoEstado": 1,
    "cantMeGusta": 150,
    "cantNoInteresa": 10,
    "cantVendidoTotal": 500,
    "totalRecaudado": 25000.50,
    "estado": 1,
    "usuarioCreacion": "423e4567-e89b-12d3-a456-426614174003",
    "fechaCreacion": "2025-11-01T10:00:00Z",
    "usuarioModificacion": null,
    "fechaModificacion": null
  }
}
```

#### Respuestas de Error
- `404 Not Found`: Evento no encontrado
- `500 Internal Server Error`: Error del servidor

---

### 2. Crear un Nuevo Evento

**Método:** `POST`  
**Ruta:** `/v1/eventos/`

#### Headers
```
Content-Type: application/json
```

#### Cuerpo de la Solicitud (Request Body)
```json
{
  "organizadorId": "223e4567-e89b-12d3-a456-426614174001",
  "categoriaId": "323e4567-e89b-12d3-a456-426614174002",
  "titulo": "Concierto de Rock",
  "descripcion": "Un increíble concierto de rock con las mejores bandas locales",
  "lugar": "Estadio Municipal",
  "eventoEstado": 1,
  "cantMeGusta": 0,
  "cantNoInteresa": 0,
  "cantVendidoTotal": 0,
  "totalRecaudado": 0,
  "estado": 1,
  "usuarioCreacion": "423e4567-e89b-12d3-a456-426614174003"
}
```

#### Campos Requeridos
| Campo | Tipo | Descripción | Validación |
|-------|------|-------------|------------|
| `organizadorId` | UUID | ID del organizador del evento | Requerido, no puede ser vacío |
| `categoriaId` | UUID | ID de la categoría del evento | Requerido, no puede ser vacío |
| `titulo` | string | Título del evento | Requerido, máximo 80 caracteres |
| `descripcion` | string | Descripción detallada del evento | Requerido |
| `lugar` | string | Ubicación del evento | Requerido, máximo 80 caracteres |
| `eventoEstado` | int | Estado del evento | Requerido |
| `cantMeGusta` | int | Cantidad de "me gusta" | Requerido, debe ser >= 0 |
| `cantNoInteresa` | int | Cantidad de "no me interesa" | Requerido, debe ser >= 0 |
| `cantVendidoTotal` | int | Cantidad total de tickets vendidos | Requerido, debe ser >= 0 |
| `totalRecaudado` | float64 | Total recaudado en ventas | Requerido, debe ser >= 0 |
| `estado` | int | Estado del registro (1: activo, 0: inactivo) | Requerido |
| `usuarioCreacion` | UUID | ID del usuario que crea el evento | Requerido |

#### Respuesta Exitosa (201 Created)
```json
{
  "evento": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "organizadorId": "223e4567-e89b-12d3-a456-426614174001",
    "categoriaId": "323e4567-e89b-12d3-a456-426614174002",
    "titulo": "Concierto de Rock",
    "descripcion": "Un increíble concierto de rock con las mejores bandas locales",
    "lugar": "Estadio Municipal",
    "eventoEstado": 1,
    "cantMeGusta": 0,
    "cantNoInteresa": 0,
    "cantVendidoTotal": 0,
    "totalRecaudado": 0,
    "estado": 1,
    "usuarioCreacion": "423e4567-e89b-12d3-a456-426614174003",
    "fechaCreacion": "2025-11-10T14:30:00Z",
    "usuarioModificacion": null,
    "fechaModificacion": null
  }
}
```

#### Headers de Respuesta
```
Location: /v1/eventos/123e4567-e89b-12d3-a456-426614174000
```

#### Respuestas de Error
- `400 Bad Request`: Datos inválidos o falta de campos requeridos
- `422 Unprocessable Entity`: Validación fallida
- `500 Internal Server Error`: Error del servidor

---

### 3. Actualizar un Evento

**Método:** `PUT`  
**Ruta:** `/v1/eventos/:id`

#### Parámetros
- `id` (UUID, requerido): ID del evento a actualizar

#### Headers
```
Content-Type: application/json
```

#### Cuerpo de la Solicitud (Request Body)
Todos los campos son opcionales. Solo se actualizarán los campos enviados.

```json
{
  "categoriaId": "323e4567-e89b-12d3-a456-426614174002",
  "titulo": "Concierto de Rock - Edición Especial",
  "descripcion": "Un increíble concierto de rock con las mejores bandas locales - Ahora con artistas internacionales",
  "lugar": "Estadio Nacional",
  "eventoEstado": 2,
  "cantMeGusta": 200,
  "cantNoInteresa": 15,
  "cantVendidoTotal": 750,
  "totalRecaudado": 37500.75,
  "estado": 1
}
```

#### Campos Actualizables
| Campo | Tipo | Descripción |
|-------|------|-------------|
| `categoriaId` | UUID | ID de la categoría del evento |
| `titulo` | string | Título del evento (máximo 80 caracteres) |
| `descripcion` | string | Descripción del evento |
| `lugar` | string | Ubicación del evento (máximo 80 caracteres) |
| `eventoEstado` | int | Estado del evento |
| `cantMeGusta` | int | Cantidad de "me gusta" (>= 0) |
| `cantNoInteresa` | int | Cantidad de "no me interesa" (>= 0) |
| `cantVendidoTotal` | int | Total de tickets vendidos (>= 0) |
| `totalRecaudado` | float64 | Total recaudado (>= 0) |
| `estado` | int | Estado del registro |

**Nota:** El campo `organizadorId` no se puede actualizar.

#### Respuesta Exitosa (200 OK)
```json
{
  "evento": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "organizadorId": "223e4567-e89b-12d3-a456-426614174001",
    "categoriaId": "323e4567-e89b-12d3-a456-426614174002",
    "titulo": "Concierto de Rock - Edición Especial",
    "descripcion": "Un increíble concierto de rock con las mejores bandas locales - Ahora con artistas internacionales",
    "lugar": "Estadio Nacional",
    "eventoEstado": 2,
    "cantMeGusta": 200,
    "cantNoInteresa": 15,
    "cantVendidoTotal": 750,
    "totalRecaudado": 37500.75,
    "estado": 1,
    "usuarioCreacion": "423e4567-e89b-12d3-a456-426614174003",
    "fechaCreacion": "2025-11-10T14:30:00Z",
    "usuarioModificacion": "523e4567-e89b-12d3-a456-426614174004",
    "fechaModificacion": "2025-11-10T16:45:00Z"
  }
}
```

#### Respuestas de Error
- `400 Bad Request`: Datos inválidos
- `404 Not Found`: Evento no encontrado
- `422 Unprocessable Entity`: Validación fallida
- `500 Internal Server Error`: Error del servidor

---

## Ejemplo Completo con cURL

### Crear un Evento
```bash
curl -X POST http://localhost:4000/v1/eventos/ \
  -H "Content-Type: application/json" \
  -d '{
    "organizadorId": "223e4567-e89b-12d3-a456-426614174001",
    "categoriaId": "323e4567-e89b-12d3-a456-426614174002",
    "titulo": "Festival de Música Electrónica",
    "descripcion": "El mejor festival de música electrónica de la región",
    "lugar": "Parque Central",
    "eventoEstado": 1,
    "cantMeGusta": 0,
    "cantNoInteresa": 0,
    "cantVendidoTotal": 0,
    "totalRecaudado": 0,
    "estado": 1,
    "usuarioCreacion": "423e4567-e89b-12d3-a456-426614174003"
  }'
```

### Obtener un Evento
```bash
curl -X GET http://localhost:4000/v1/eventos/123e4567-e89b-12d3-a456-426614174000
```

### Actualizar un Evento
```bash
curl -X PUT http://localhost:4000/v1/eventos/123e4567-e89b-12d3-a456-426614174000 \
  -H "Content-Type: application/json" \
  -d '{
    "titulo": "Festival de Música Electrónica - 2025",
    "cantMeGusta": 350,
    "cantVendidoTotal": 1200,
    "totalRecaudado": 60000.00
  }'
```

---

## Códigos de Estado HTTP

| Código | Descripción |
|--------|-------------|
| 200 | OK - Solicitud exitosa |
| 201 | Created - Recurso creado exitosamente |
| 400 | Bad Request - Solicitud malformada |
| 404 | Not Found - Recurso no encontrado |
| 422 | Unprocessable Entity - Error de validación |
| 500 | Internal Server Error - Error del servidor |

---

## Notas Adicionales

- Todos los IDs son UUID v4
- Los timestamps se devuelven en formato ISO 8601 (RFC3339)
- El sistema implementa "soft delete" (el campo `estado` se pone en 0 en vez de eliminar físicamente)
- Las fechas de creación y modificación se generan automáticamente en el servidor
- Los campos de contadores (`cantMeGusta`, `cantNoInteresa`, `cantVendidoTotal`) y `totalRecaudado` no pueden ser negativos
