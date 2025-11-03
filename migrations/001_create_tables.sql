CREATE TYPE tipo_documento_enum AS ENUM ('DNI', 'CE', 'RUC');
DROP TABLE IF EXISTS usuario;
CREATE TABLE usuario (
    usuario_id BIGSERIAL PRIMARY KEY,
    nombre VARCHAR(80) NOT NULL,
    tipo_documento tipo_documento_enum NOT NULL,
    num_documento VARCHAR (25) NOT NULL,
    correo VARCHAR (80) NOT NULL UNIQUE,
    contrasenha VARCHAR (255) NOT NULL,
    --hash
    telefono VARCHAR (15),
    estado_de_cuenta SMALLINT NOT NULL DEFAULT 0,
    codigo_verificacion VARCHAR (10),
    fecha_expiracion_codigo TIMESTAMP,
    usuario_creacion BIGINT,
    fecha_creacion TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    usuario_modificacion BIGINT,
    fecha_modificacion TIMESTAMPTZ,
    activo SMALLINT NOT NULL DEFAULT 1,
    CONSTRAINT uq_usuario_doc UNIQUE (tipo_documento, num_documento),
    CONSTRAINT chk_estado_cuenta CHECK (estado_de_cuenta IN (0, 1, 2))
);
DROP TABLE IF EXISTS categoria;
CREATE TABLE categoria (
    id_categoria BIGSERIAL PRIMARY KEY,
    nombre VARCHAR(80) NOT NULL UNIQUE,
    descripcion TEXT NOT NULL DEFAULT '',
    activo SMALLINT NOT NULL DEFAULT 1,
    CONSTRAINT chk_categoria_activo CHECK (activo IN (0, 1))
);
DROP TABLE IF EXISTS evento;
CREATE TABLE evento (
    evento_id BIGSERIAL PRIMARY KEY,
    organizador_id BIGINT NOT NULL,
    categoria_id BIGINT NOT NULL,
    titulo VARCHAR(80) NOT NULL,
    descripcion TEXT NOT NULL,
    lugar VARCHAR(80) NOT NULL,
    evento_estado SMALLINT NOT NULL DEFAULT 0,
    cant_me_gusta INT NOT NULL DEFAULT 0,
    cant_no_interesa INT NOT NULL DEFAULT 0,
    cant_vendido_total INT NOT NULL DEFAULT 0,
    total_recaudado NUMERIC(12, 2) NOT NULL DEFAULT 0,
    activo SMALLINT NOT NULL DEFAULT 1,
    usuario_creacion BIGINT,
    fecha_creacion TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    usuario_modificacion BIGINT,
    fecha_modificacion TIMESTAMPTZ,
    CONSTRAINT fk_evento_organizador FOREIGN KEY (organizador_id) REFERENCES usuario(usuario_id) ON DELETE RESTRICT,
    CONSTRAINT fk_evento_categoria FOREIGN KEY (categoria_id) REFERENCES categoria(categoria_id) ON DELETE RESTRICT,
    CONSTRAINT chk_evento_estado CHECK (evento_estado IN (0, 1, 2)),
    CONSTRAINT chk_evento_activo CHECK (activo IN (0, 1))
);
DROP TABLE IF EXISTS comentario;
CREATE TABLE comentario (
    comentario_id BIGSERIAL PRIMARY KEY,
    usuario_id BIGINT NOT NULL,
    evento_id BIGINT NOT NULL,
    descripcion TEXT NOT NULL,
    fecha_creacion TIMESTAMP NOT NULL DEFAULT NOW(),
    activo SMALLINT NOT NULL DEFAULT 1,
    CONSTRAINT fk_comentario_usuario FOREIGN KEY (usuario_id) REFERENCES usuario(usuario_id) ON DELETE RESTRICT,
    CONSTRAINT fk_comentario_evento FOREIGN KEY (evento_id) REFERENCES evento(evento_id) ON DELETE RESTRICT,
    CONSTRAINT chk_comentario_activo CHECK (activo IN (0, 1))
);
DROP TABLE IF EXISTS sector;
CREATE TABLE sector (
    sector_id BIGSERIAL PRIMARY KEY,
    evento_id BIGINT NOT NULL,
    sector_tipo VARCHAR(30) NOT NULL,
    total_entradas INT NOT NULL,
    cant_vendidas INT NOT NULL DEFAULT 0,
    activo SMALLINT NOT NULL DEFAULT 1,
    usuario_creacion BIGINT,
    fecha_creacion TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    usuario_modificacion BIGINT,
    fecha_modificacion TIMESTAMPTZ,
    CONSTRAINT fk_sector_evento FOREIGN KEY (evento_id) REFERENCES evento(evento_id) ON DELETE RESTRICT,
    CONSTRAINT uq_sector_tipo UNIQUE (evento_id, sector_tipo),
    CONSTRAINT chk_sector_activo CHECK (activo IN (0, 1))
);
DROP TABLE IF EXISTS tipo_de_ticket;
CREATE TABLE tipo_de_ticket (
    tipo_de_ticket_id BIGSERIAL PRIMARY KEY,
    evento_id BIGINT NOT NULL,
    nombre VARCHAR (25) NOT NULL,
    fecha_ini DATE NOT NULL,
    fecha_fin DATE NOT NULL,
    activo SMALLINT NOT NULL DEFAULT 1,
    usuario_creacion BIGINT,
    fecha_creacion TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    usuario_modificacion BIGINT,
    fecha_modificacion TIMESTAMPTZ,
    CONSTRAINT fk_tipo_de_ticket_evento FOREIGN KEY (evento_id) REFERENCES evento(evento_id) ON DELETE RESTRICT CONSTRAINT uq_tipo_ticket_nombre UNIQUE (evento_id, nombre),
    CONSTRAINT chk_tipo_de_ticket_activo CHECK (activo IN (0, 1))
);
DROP TABLE IF EXISTS perfil_de_persona;
CREATE TABLE perfil_de_persona (
    perfil_de_persona_id BIGSERIAL PRIMARY KEY,
    evento_id BIGINT NOT NULL,
    nombre VARCHAR (25) NOT NULL,
    activo SMALLINT NOT NULL DEFAULT 1,
    usuario_creacion BIGINT,
    fecha_creacion TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    usuario_modificacion BIGINT,
    fecha_modificacion TIMESTAMPTZ,
    CONSTRAINT fk_perfil_de_persona_evento FOREIGN KEY (evento_id) REFERENCES evento(evento_id) ON DELETE RESTRICT,
    CONSTRAINT uq_perfil_de_persona_nombre UNIQUE (evento_id, nombre),
    CONSTRAINT chk_perfil_de_persona_activo CHECK (activo IN (0, 1))
) DROP TABLE IF EXISTS tarifa;
CREATE TABLE tarifa (
    tarifa_id BIGSERIAL PRIMARY KEY,
    sector_id BIGINT NOT NULL,
    tipo_de_ticket_id BIGINT NOT NULL,
    perfil_de_persona_id BIGINT,
    precio NUMERIC(4, 2) NOT NULL DEFAULT 0,
    activo SMALLINT NOT NULL DEFAULT 1,
    usuario_creacion BIGINT,
    fecha_creacion TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    usuario_modificacion BIGINT,
    fecha_modificacion TIMESTAMPTZ,
    CONSTRAINT fk_tarifa_sector FOREIGN KEY (sector_id) REFERENCES sector(sector_id) ON DELETE RESTRICT,
    CONSTRAINT fk_tarifa_tipo_de_ticket FOREIGN KEY (tipo_de_ticket_id) REFERENCES tipo_de_ticket(tipo_de_ticket_id) ON DELETE RESTRICT,
    CONSTRAINT fk_tarifa_perfil_de_persona FOREIGN KEY (perfil_de_persona_id) REFERENCES perfil_de_persona(perfil_de_persona_id) ON DELETE RESTRICT,
    CONSTRAINT chk_tarifa_activo CHECK (activo IN (0, 1))
);
CREATE TYPE tipo_metodo_pago_enum AS ENUM ('Tarjeta', 'Yape');
DROP TABLE IF EXISTS metodo_de_pago;
CREATE TABLE IF EXISTS metodo_de_pago(
    metodo_de_pago_id BIGSERIAL PRIMARY KEY,
    tipo tipo_metodo_pago_enum NOT NULL,
    activo SMALLINT NOT NULL DEFAULT 1,
    CONSTRAINT chk_metodo_de_pago_activo CHECK (activo IN (0, 1))
);
DROP TABLE IF EXISTS orden_de_compra;
CREATE TABLE orden_de_compra(
    orden_de_compra_id BIGSERIAL PRIMARY KEY,
    usuario_id BIGINT NOT NULL,
    metodo_de_pago_id BIGINT NOT NULL,
    fecha DATE NOT NULL DEFAULT NOW(),
    total NUMERIC(4, 2) NOT NULL,
    monto_fee_servicio NUMERIC(4, 2) NOT NULL,
    estado_de_orden SMALLINT NOT NULL DEFAULT 0,
    CONSTRAINT fk_orden_de_compra_usuario FOREIGN KEY (usuario_id) REFERENCES usuario(usuario_id),
    CONSTRAINT fk_orden_de_compra_pago FOREIGN KEY (metodo_de_pago_id) REFERENCES metodo_de_pago(metodo_de_pago_id),
    CONSTRAINT chk_orden_de_compra_estado CHECK (estado_de_orden IN (0, 1, 2))
);
DROP TABLE IF EXISTS fecha;
CREATE TABLE fecha(
    fecha_id BIGSERIAL PRIMARY KEY,
    fecha_evento DATE NOT NULL,
);
DROP TABLE IF EXISTS evento_fecha;
CREATE TABLE evento_fecha(
    evento_fecha_id BIGSERIAL PRIMARY KEY,
    evento_id BIGINT NOT NULL,
    fecha_id BIGINT NOT NULL,
    hora_inicio TIMESTAMP NOT NULL,
    activo SMALLINT NOT NULL DEFAULT 1,
    usuario_creacion BIGINT,
    fecha_creacion TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    usuario_modificacion BIGINT,
    fecha_modificacion TIMESTAMPTZ,
    CONSTRAINT fk_evento_fecha_evento FOREIGN KEY (evento_id) REFERENCES evento(evento_id) ON DELETE RESTRICT,
    CONSTRAINT fk_evento_fecha_fecha FOREIGN KEY (fecha_id) REFERENCES fecha(fecha_id) ON DELETE RESTRICT,
    CONSTRAINT uq_evento_fecha UNIQUE (evento_id, fecha_id, horaInicio),
    CONSTRAINT chk_evento_fecha_activo CHECK (activo IN (0, 1))
);
DROP TABLE IF EXISTS ticket;
CREATE TABLE ticket (
    ticket_id BIGSERIAL PRIMARY KEY,
    orden_de_compra_id BIGINT,
    evento_fecha_id BIGINT NOT NULL,
    tarifa_id BIGINT NOT NULL,
    codigo_qr VARCHAR(50) NOT NULL,
    estado_de_ticket SMALLINT NOT NULL DEFAULT 0,
    CONSTRAINT fk_ticket_orden FOREIGN KEY (orden_de_compra_id) REFERENCES orden_de_compra(orden_de_compra_id),
    CONSTRAINT fk_ticket_fecha FOREIGN KEY (evento_fecha_id) REFERENCES evento_fecha(evento_fecha_id),
    CONSTRAINT fk_ticket_tarifa FOREIGN KEY (tarifa_id) REFERENCES tarifa(tarifa_id),
    CONSTRAINT chk_ticket_estado CHECK (estado_de_ticket IN (0, 1, 2, 3))
) DROP TABLE IF EXISTS cupon;
CREATE TABLE cupon(
    cupon_id BIGSERIAL PRIMARY KEY,
    descripcion TEXT NOT NULL,
    tipo VARCHAR(20) NOT NULL,
    valor NUMERIC(10, 2) NOT NULL,
    estado_cupon SMALLINT NOT NULL DEFAULT 0,
    codigo VARCHAR(20) NOT NULL,
    uso_por_usuario BIGINT NOT NULL DEFAULT 0,
    uso_realizados BIGINT NOT NULL DEFAULT 0,
    usuario_creacion BIGINT,
    fecha_creacion TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    usuario_modificacion BIGINT,
    fecha_modificacion TIMESTAMPTZ,
    CONSTRAINT uq_cupon_codigo UNIQUE (codigo),
    CONSTRAINT chk_cupon_estado CHECK (estado_cupon IN (0, 1))
);
DROP TABLE IF EXISTS evento_cupon;
CREATE TABLE evento_cupon(
    evento_id BIGINT NOT NULL,
    cupon_id BIGINT NOT NULL,
    cant_cupones BIGINT NOT NULL,
    fecha_ini DATE NOT NULL,
    fecha_fin DATE NOT NULL,
    usuario_creacion BIGINT,
    fecha_creacion TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    usuario_modificacion BIGINT,
    fecha_modificacion TIMESTAMPTZ,
    CONSTRAINT pk_evento_cupon PRIMARY KEY (evento_id, cupon_id),
    CONSTRAINT fk_evento_cupon_evento FOREIGN KEY (evento_id) REFERENCES evento(evento_id) ON DELETE RESTRICT,
    CONSTRAINT fk_evento_cupon_cupon FOREIGN KEY (cupon_id) REFERENCES cupon(cupon_id) ON DELETE RESTRICT,
    CONSTRAINT chk_evento_cupon_rango CHECK (fecha_fin >= fecha_ini)
);
DROP TABLE IF EXISTS usuario_cupon;
CREATE TABLE usuario_cupon(
    cupon_id BIGINT NOT NULL,
    usuario_id BIGINT NOT NULL,
    cant_usada BIGINT NOT NULL DEFAULT 0,
    CONSTRAINT pk_usuario_cupon PRIMARY KEY (cupon_id, usuario_id),
    CONSTRAINT fk_usuario_cupon_usuario FOREIGN KEY (usuario_id) REFERENCES usuario(usuario_id) ON DELETE RESTRICT,
    CONSTRAINT fk_evento_cupon_cupon FOREIGN KEY (cupon_id) REFERENCES cupon(cupon_id) ON DELETE RESTRICT,
);
DROP TABLE IF EXISTS comprobante_de_pago;
CREATE TABLE comprobante_de_pago(
    comprobante_de_pago_id BIGSERIAL PRIMARY KEY,
    orden_de_compra_id BIGINT NOT NULL,
    tipo_de_comprobante SMALLINT NOT NULL DEFAULT 0,
    numero VARCHAR(20) NOT NULL,
    fecha_emision TIMESTAMP NOT NULL,
    ruc VARCHAR(20),
    direccion_fiscal VARCHAR(80),
    CONSTRAINT fk_comprobante_de_pago_orden FOREIGN KEY (orden_de_compra_id) REFERENCES orden_de_compra(orden_de_compra_id) ON DELETE RESTRICT,
    CONSTRAINT chk_comprobante_de_pago_tipo CHECK (tipo_de_comprobante IN (0, 1))
);
DROP TABLE IF EXISTS rol;
CREATE TABLE rol(
    rol_id BIGSERIAL PRIMARY KEY,
    nombre VARCHAR(20) NOT NULL,
    usuario_creacion BIGINT,
    fecha_creacion TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    usuario_modificacion BIGINT,
    fecha_modificacion TIMESTAMPTZ,
);
DROP TABLE IF EXISTS rol_usuario;
CREATE TABLE rol_usuario(
    rol_usuario_id BIGSERIAL PRIMARY KEY,
    rol_id BIGINT NOT NULL,
    usuario_id BIGINT NOT NULL,
    usuario_creacion BIGINT,
    fecha_creacion TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    usuario_modificacion BIGINT,
    fecha_modificacion TIMESTAMPTZ,
    CONSTRAINT fk_rol_usuario_rol FOREIGN KEY (rol_id) REFERENCES rol(rol_id),
    CONSTRAINT fk_rol_usuario_usuario FOREIGN KEY (usuario_id) REFERENCES usuario(usuario_id)
);
DROP TABLE IF EXISTS notificacion;
CREATE TABLE notificacion(
    notificacion_id BIGSERIAL PRIMARY KEY,
    mensaje TEXT NOT NULL,
    canal VARCHAR(40) NOT NULL,
    fecha_envio TIMESTAMP NOT NULL,
    estado_notificación SMALLINT NOT NULL,
    CONSTRAINT chk_notificacion CHECK (notificacion IN (0, 1, 2))
);
-- Índice para búsquedas case-insensitive por nombre (opcional)
CREATE UNIQUE INDEX IF NOT EXISTS idx_categorias_nombre_ci ON categorias (LOWER(nombre));