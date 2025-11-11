-- =========================================================
-- RESET: DROP tables (hijas/asociativas primero) y tipos
-- =========================================================
DROP TABLE IF EXISTS rol_usuario;
DROP TABLE IF EXISTS usuario_cupon;
DROP TABLE IF EXISTS evento_cupon;
DROP TABLE IF EXISTS ticket;
DROP TABLE IF EXISTS comprobante_de_pago;
DROP TABLE IF EXISTS evento_fecha;
DROP TABLE IF EXISTS fecha;
DROP TABLE IF EXISTS tarifa;
DROP TABLE IF EXISTS sector;
DROP TABLE IF EXISTS tipo_de_ticket;
DROP TABLE IF EXISTS perfil_de_persona;
DROP TABLE IF EXISTS comentario;
DROP TABLE IF EXISTS orden_de_compra;
DROP TABLE IF EXISTS metodo_de_pago;
DROP TABLE IF EXISTS evento;
DROP TABLE IF EXISTS cupon;
DROP TABLE IF EXISTS rol;
DROP TABLE IF EXISTS notificacion;
DROP TABLE IF EXISTS categoria;
DROP TABLE IF EXISTS usuario;
DROP TYPE IF EXISTS tipo_metodo_pago_enum;
DROP TYPE IF EXISTS tipo_documento_enum;

-- =========================================================
-- EXTENSIÓN UUID
-- =========================================================
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- =========================================================
-- TIPOS
-- =========================================================
CREATE TYPE tipo_documento_enum AS ENUM ('DNI', 'CE', 'RUC');
CREATE TYPE tipo_metodo_pago_enum AS ENUM ('Tarjeta', 'Yape');
-- =========================================================
-- TABLAS BASE
-- =========================================================
CREATE TABLE usuario (
    usuario_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nombre VARCHAR(80) NOT NULL,
    tipo_documento tipo_documento_enum NOT NULL,
    num_documento VARCHAR(25) NOT NULL,
    correo VARCHAR(80) NOT NULL UNIQUE,
    contrasenha VARCHAR(255) NOT NULL,
    telefono VARCHAR(15),
    estado_de_cuenta SMALLINT NOT NULL DEFAULT 0,
    codigo_verificacion VARCHAR(10),
    fecha_expiracion_codigo TIMESTAMPTZ,
    usuario_creacion UUID,
    fecha_creacion TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    usuario_modificacion UUID,
    fecha_modificacion TIMESTAMPTZ,
    estado SMALLINT NOT NULL DEFAULT 1,
    CONSTRAINT uq_usuario_doc UNIQUE (tipo_documento, num_documento),
    CONSTRAINT chk_usuario_estado_cta CHECK (estado_de_cuenta IN (0, 1, 2)),
    CONSTRAINT chk_usuario_estado CHECK (estado IN (0, 1)),
    -- patrón simple (no perfecto) para evitar errores groseros
    CONSTRAINT chk_usuario_correo_fmt CHECK (correo LIKE '%_@_%.__%')
);

CREATE TABLE categoria (
    id_categoria UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nombre VARCHAR(80) NOT NULL UNIQUE,
    descripcion TEXT NOT NULL DEFAULT '',
    estado SMALLINT NOT NULL DEFAULT 1,
    CONSTRAINT chk_categoria_estado CHECK (estado IN (0, 1))
);

CREATE TABLE evento (
    evento_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organizador_id UUID NOT NULL,
    categoria_id UUID NOT NULL,
    titulo VARCHAR(80) NOT NULL,
    descripcion TEXT NOT NULL,
    lugar VARCHAR(80) NOT NULL,
    evento_estado SMALLINT NOT NULL DEFAULT 0,
    cant_me_gusta INT NOT NULL DEFAULT 0,
    cant_no_interesa INT NOT NULL DEFAULT 0,
    cant_vendido_total INT NOT NULL DEFAULT 0,
    total_recaudado NUMERIC(12, 2) NOT NULL DEFAULT 0,
    estado SMALLINT NOT NULL DEFAULT 1,
    usuario_creacion UUID,
    fecha_creacion TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    usuario_modificacion UUID,
    fecha_modificacion TIMESTAMPTZ,
    imagen_descripcion VARCHAR(255),
    imagen_portada VARCHAR(255),
    video_presentacion VARCHAR(255),
    imagen_escenario VARCHAR(255)
    -- CONSTRAINT fk_evento_organizador FOREIGN KEY (organizador_id) REFERENCES usuario(usuario_id) ON DELETE RESTRICT,
    -- CONSTRAINT fk_evento_categoria FOREIGN KEY (categoria_id) REFERENCES categoria(id_categoria) ON DELETE RESTRICT,
    -- CONSTRAINT chk_evento_estado CHECK (evento_estado IN (0, 1, 2)),
    -- CONSTRAINT chk_evento_estado_flag CHECK (estado IN (0, 1)),
    -- CONSTRAINT chk_evento_contadores_nn CHECK (
    --     cant_me_gusta >= 0
    --     AND cant_no_interesa >= 0
    --     AND cant_vendido_total >= 0
    --     AND total_recaudado >= 0
    -- )
);

CREATE TABLE comentario (
    comentario_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    usuario_id UUID NOT NULL,
    evento_id UUID NOT NULL,
    descripcion TEXT NOT NULL,
    fecha_creacion TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    estado SMALLINT NOT NULL DEFAULT 1,
    CONSTRAINT fk_comentario_usuario FOREIGN KEY (usuario_id) REFERENCES usuario(usuario_id) ON DELETE RESTRICT,
    CONSTRAINT fk_comentario_evento FOREIGN KEY (evento_id) REFERENCES evento(evento_id) ON DELETE RESTRICT,
    CONSTRAINT chk_comentario_estado CHECK (estado IN (0, 1))
);

CREATE TABLE sector (
    sector_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    evento_id UUID NOT NULL,
    sector_tipo VARCHAR(30) NOT NULL,
    total_entradas INT NOT NULL,
    cant_vendidas INT NOT NULL DEFAULT 0,
    estado SMALLINT NOT NULL DEFAULT 1,
    usuario_creacion UUID,
    fecha_creacion TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    usuario_modificacion UUID,
    fecha_modificacion TIMESTAMPTZ,
    CONSTRAINT fk_sector_evento FOREIGN KEY (evento_id) REFERENCES evento(evento_id) ON DELETE RESTRICT,
    CONSTRAINT uq_sector_tipo UNIQUE (evento_id, sector_tipo),
    CONSTRAINT chk_sector_estado CHECK (estado IN (0, 1)),
    CONSTRAINT chk_sector_capacidad CHECK (
        total_entradas > 0
        AND cant_vendidas >= 0
        AND cant_vendidas <= total_entradas
    )
);

CREATE TABLE tipo_de_ticket (
    tipo_de_ticket_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    evento_id UUID NOT NULL,
    nombre VARCHAR(25) NOT NULL,
    fecha_ini DATE NOT NULL,
    fecha_fin DATE NOT NULL,
    estado SMALLINT NOT NULL DEFAULT 1,
    usuario_creacion UUID,
    fecha_creacion TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    usuario_modificacion UUID,
    fecha_modificacion TIMESTAMPTZ,
    CONSTRAINT fk_tipo_de_ticket_evento FOREIGN KEY (evento_id) REFERENCES evento(evento_id) ON DELETE RESTRICT,
    CONSTRAINT uq_tipo_ticket_nombre UNIQUE (evento_id, nombre),
    CONSTRAINT chk_tipo_de_ticket_estado CHECK (estado IN (0, 1)),
    CONSTRAINT chk_tipo_de_ticket_rango CHECK (fecha_fin >= fecha_ini)
);

CREATE TABLE perfil_de_persona (
    perfil_de_persona_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    evento_id UUID NOT NULL,
    nombre VARCHAR(25) NOT NULL,
    estado SMALLINT NOT NULL DEFAULT 1,
    usuario_creacion UUID,
    fecha_creacion TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    usuario_modificacion UUID,
    fecha_modificacion TIMESTAMPTZ,
    CONSTRAINT fk_perfil_de_persona_evento FOREIGN KEY (evento_id) REFERENCES evento(evento_id) ON DELETE RESTRICT,
    CONSTRAINT uq_perfil_de_persona_nombre UNIQUE (evento_id, nombre),
    CONSTRAINT chk_perfil_de_persona_estado CHECK (estado IN (0, 1))
);

CREATE TABLE tarifa (
    tarifa_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sector_id UUID NOT NULL,
    tipo_de_ticket_id UUID NOT NULL,
    perfil_de_persona_id UUID,
    precio NUMERIC(10, 2) NOT NULL DEFAULT 0,
    estado SMALLINT NOT NULL DEFAULT 1,
    usuario_creacion UUID,
    fecha_creacion TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    usuario_modificacion UUID,
    fecha_modificacion TIMESTAMPTZ,
    CONSTRAINT fk_tarifa_sector FOREIGN KEY (sector_id) REFERENCES sector(sector_id) ON DELETE RESTRICT,
    CONSTRAINT fk_tarifa_tipo_de_ticket FOREIGN KEY (tipo_de_ticket_id) REFERENCES tipo_de_ticket(tipo_de_ticket_id) ON DELETE RESTRICT,
    CONSTRAINT fk_tarifa_perfil_de_persona FOREIGN KEY (perfil_de_persona_id) REFERENCES perfil_de_persona(perfil_de_persona_id) ON DELETE RESTRICT,
    CONSTRAINT chk_tarifa_estado CHECK (estado IN (0, 1)),
    CONSTRAINT chk_tarifa_precio_nn CHECK (precio >= 0)
);

CREATE TABLE metodo_de_pago (
    metodo_de_pago_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tipo tipo_metodo_pago_enum NOT NULL,
    estado SMALLINT NOT NULL DEFAULT 1,
    CONSTRAINT chk_metodo_de_pago_estado CHECK (estado IN (0, 1))
);

CREATE TABLE orden_de_compra(
    orden_de_compra_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    usuario_id UUID NOT NULL,
    metodo_de_pago_id UUID NOT NULL,
    fecha DATE NOT NULL DEFAULT CURRENT_DATE,
    fecha_hora_ini TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    fecha_hora_fin TIMESTAMPTZ,
    total NUMERIC(12, 2) NOT NULL,
    monto_fee_servicio NUMERIC(12, 2) NOT NULL,
    estado_de_orden SMALLINT NOT NULL DEFAULT 0,
    CONSTRAINT fk_orden_de_compra_usuario FOREIGN KEY (usuario_id) REFERENCES usuario(usuario_id),
    CONSTRAINT fk_orden_de_compra_pago FOREIGN KEY (metodo_de_pago_id) REFERENCES metodo_de_pago(metodo_de_pago_id),
    CONSTRAINT chk_orden_de_compra_estado CHECK (estado_de_orden IN (0, 1, 2)),
    CONSTRAINT chk_orden_de_compra_rango CHECK (
        fecha_hora_fin IS NULL
        OR fecha_hora_fin >= fecha_hora_ini
    ),
    CONSTRAINT chk_orden_de_compra_montos CHECK (
        total >= 0
        AND monto_fee_servicio >= 0
    )
);

CREATE TABLE fecha (
    fecha_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fecha_evento DATE NOT NULL
);

CREATE TABLE evento_fecha (
    evento_fecha_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    evento_id UUID NOT NULL,
    fecha_id UUID NOT NULL,
    hora_inicio TIMESTAMPTZ NOT NULL,
    estado SMALLINT NOT NULL DEFAULT 1,
    usuario_creacion UUID,
    fecha_creacion TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    usuario_modificacion UUID,
    fecha_modificacion TIMESTAMPTZ,
    CONSTRAINT fk_evento_fecha_evento FOREIGN KEY (evento_id) REFERENCES evento(evento_id) ON DELETE RESTRICT,
    CONSTRAINT fk_evento_fecha_fecha FOREIGN KEY (fecha_id) REFERENCES fecha(fecha_id) ON DELETE RESTRICT,
    CONSTRAINT uq_evento_fecha UNIQUE (evento_id, fecha_id, hora_inicio),
    CONSTRAINT chk_evento_fecha_estado CHECK (estado IN (0, 1))
);

CREATE TABLE ticket (
    ticket_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    orden_de_compra_id UUID,
    evento_fecha_id UUID NOT NULL,
    tarifa_id UUID NOT NULL,
    codigo_qr VARCHAR(50) NOT NULL,
    estado_de_ticket SMALLINT NOT NULL DEFAULT 0,
    CONSTRAINT fk_ticket_orden FOREIGN KEY (orden_de_compra_id) REFERENCES orden_de_compra(orden_de_compra_id),
    CONSTRAINT fk_ticket_fecha FOREIGN KEY (evento_fecha_id) REFERENCES evento_fecha(evento_fecha_id),
    CONSTRAINT fk_ticket_tarifa FOREIGN KEY (tarifa_id) REFERENCES tarifa(tarifa_id),
    CONSTRAINT chk_ticket_estado CHECK (estado_de_ticket IN (0, 1, 2, 3)),
    CONSTRAINT uq_ticket_qr UNIQUE (codigo_qr)
);

CREATE TABLE cupon (
    cupon_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    descripcion TEXT NOT NULL,
    tipo VARCHAR(20) NOT NULL,
    valor NUMERIC(10, 2) NOT NULL,
    estado_cupon SMALLINT NOT NULL DEFAULT 0,
    codigo VARCHAR(20) NOT NULL,
    uso_por_usuario BIGINT NOT NULL DEFAULT 0,
    uso_realizados BIGINT NOT NULL DEFAULT 0,
    usuario_creacion UUID,
    fecha_creacion TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    usuario_modificacion UUID,
    fecha_modificacion TIMESTAMPTZ,
    CONSTRAINT uq_cupon_codigo UNIQUE (codigo),
    CONSTRAINT chk_cupon_estado CHECK (estado_cupon IN (0, 1)),
    CONSTRAINT chk_cupon_valor_nn CHECK (valor >= 0),
    CONSTRAINT chk_cupon_usos_nn CHECK (
        uso_por_usuario >= 0
        AND uso_realizados >= 0
    )
);

CREATE TABLE evento_cupon (
    evento_id UUID NOT NULL,
    cupon_id UUID NOT NULL,
    cant_cupones BIGINT NOT NULL,
    fecha_ini DATE NOT NULL,
    fecha_fin DATE NOT NULL,
    usuario_creacion UUID,
    fecha_creacion TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    usuario_modificacion UUID,
    fecha_modificacion TIMESTAMPTZ,
    CONSTRAINT pk_evento_cupon PRIMARY KEY (evento_id, cupon_id),
    CONSTRAINT fk_evento_cupon_evento FOREIGN KEY (evento_id) REFERENCES evento(evento_id) ON DELETE RESTRICT,
    CONSTRAINT fk_evento_cupon_cupon FOREIGN KEY (cupon_id) REFERENCES cupon(cupon_id) ON DELETE RESTRICT,
    CONSTRAINT chk_evento_cupon_rango CHECK (fecha_fin >= fecha_ini),
    CONSTRAINT chk_evento_cupon_cnt CHECK (cant_cupones > 0)
);

CREATE TABLE usuario_cupon (
    cupon_id UUID NOT NULL,
    usuario_id UUID NOT NULL,
    cant_usada BIGINT NOT NULL DEFAULT 0,
    CONSTRAINT pk_usuario_cupon PRIMARY KEY (cupon_id, usuario_id),
    CONSTRAINT fk_usuario_cupon_usuario FOREIGN KEY (usuario_id) REFERENCES usuario(usuario_id) ON DELETE RESTRICT,
    CONSTRAINT fk_usuario_cupon_cupon FOREIGN KEY (cupon_id) REFERENCES cupon(cupon_id) ON DELETE RESTRICT,
    CONSTRAINT chk_usuario_cupon_usada CHECK (cant_usada >= 0)
);

CREATE TABLE comprobante_de_pago (
    comprobante_de_pago_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    orden_de_compra_id UUID NOT NULL,
    tipo_de_comprobante SMALLINT NOT NULL DEFAULT 0,
    -- 0=boleta,1=factura (ej.)
    numero VARCHAR(20) NOT NULL,
    fecha_emision TIMESTAMPTZ NOT NULL,
    ruc VARCHAR(20),
    direccion_fiscal VARCHAR(80),
    CONSTRAINT fk_comprobante_de_pago_orden FOREIGN KEY (orden_de_compra_id) REFERENCES orden_de_compra(orden_de_compra_id) ON DELETE RESTRICT,
    CONSTRAINT chk_comprobante_de_pago_tipo CHECK (tipo_de_comprobante IN (0, 1))
);

CREATE TABLE rol (
    rol_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nombre VARCHAR(20) NOT NULL,
    usuario_creacion UUID,
    fecha_creacion TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    usuario_modificacion UUID,
    fecha_modificacion TIMESTAMPTZ,
    CONSTRAINT uq_rol_nombre UNIQUE (nombre)
);

-- Relación usuario-rol con soft revoke
CREATE TABLE rol_usuario (
    rol_usuario_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rol_id UUID NOT NULL,
    usuario_id UUID NOT NULL,
    usuario_creacion UUID,
    fecha_creacion TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    usuario_modificacion UUID,
    fecha_modificacion TIMESTAMPTZ,
    estado SMALLINT NOT NULL DEFAULT 1,
    CONSTRAINT fk_rol_usuario_rol FOREIGN KEY (rol_id) REFERENCES rol(rol_id),
    CONSTRAINT fk_rol_usuario_usuario FOREIGN KEY (usuario_id) REFERENCES usuario(usuario_id),
    CONSTRAINT uq_rol_usuario UNIQUE (usuario_id, rol_id),
    CONSTRAINT chk_rol_usuario_estado CHECK (estado IN (0, 1))
);

CREATE TABLE notificacion (
    notificacion_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    mensaje TEXT NOT NULL,
    canal VARCHAR(40) NOT NULL,
    fecha_envio TIMESTAMPTZ NOT NULL,
    estado_notificacion SMALLINT NOT NULL,
    CONSTRAINT chk_notificacion CHECK (estado_notificacion IN (0, 1, 2))
);