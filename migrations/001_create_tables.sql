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
    CONSTRAINT uq_usuario_doc UNIQUE (tipo_documento, num_documento),
    CONSTRAINT chk_estado_cuenta CHECK (estado_de_cuenta IN (0, 1, 2))
);
DROP TABLE IF EXISTS categoria;
CREATE TABLE categoria (
    id_categoria BIGSERIAL PRIMARY KEY,
    nombre VARCHAR(80) NOT NULL UNIQUE,
    descripcion TEXT NOT NULL DEFAULT ''
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
    CONSTRAINT fk_evento_organizador FOREIGN KEY (organizador_id) REFERENCES usuario(usuario_id) ON DELETE RESTRICT,
    CONSTRAINT fk_evento_categoria FOREIGN KEY (categoria_id) REFERENCES categoria(categoria_id) ON DELETE RESTRICT,
    CONSTRAINT chk_evento_estado CHECK (evento_estado IN (0, 1, 2))
);
DROP TABLE IF EXISTS comentario;
CREATE TABLE comentario (
    comentario_id BIGSERIAL PRIMARY KEY,
    usuario_id BIGINT NOT NULL,
    evento_id BIGINT NOT NULL,
    descripcion TEXT NOT NULL,
    fecha_creacion TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_comentario_usuario FOREIGN KEY (usuario_id) REFERENCES usuario(usuario_id) ON DELETE RESTRICT,
    CONSTRAINT fk_comentario_evento FOREIGN KEY (evento_id) REFERENCES evento(evento_id) ON DELETE RESTRICT
);
DROP TABLE IF EXISTS sector;
CREATE TABLE sector (
    sector_id BIGSERIAL PRIMARY KEY,
    evento_id BIGINT NOT NULL,
    sector_tipo VARCHAR(30) NOT NULL,
    total_entradas INT NOT NULL,
    cant_vendidas INT NOT NULL DEFAULT 0,
    CONSTRAINT fk_sector_evento FOREIGN KEY (evento_id) REFERENCES evento(evento_id) ON DELETE RESTRICT,
    CONSTRAINT uq_sector_tipo UNIQUE (evento_id, sector_tipo)
);
DROP TABLE IF EXISTS tipo_de_ticket;
CREATE TABLE tipo_de_ticket (
    tipo_de_ticket_id BIGSERIAL PRIMARY KEY,
    evento_id BIGINT NOT NULL,
    nombre VARCHAR (25) NOT NULL,
    fecha_ini TIMESTAMP NOT NULL,
    fecha_fin TIMESTAMP NOT NULL,
    CONSTRAINT fk_tipo_de_ticket_evento FOREIGN KEY (evento_id) REFERENCES evento(evento_id) ON DELETE RESTRICT CONSTRAINT uq_tipo_ticket_nombre UNIQUE (evento_id, nombre)
);
DROP TABLE IF EXISTS perfil_de_persona;
CREATE TABLE perfil_de_persona (
    perfil_de_persona_id BIGSERIAL PRIMARY KEY,
    evento_id BIGINT NOT NULL,
    nombre VARCHAR (25) NOT NULL,
    CONSTRAINT fk_perfil_de_persona_evento FOREIGN KEY (evento_id) REFERENCES evento(evento_id) ON DELETE RESTRICT,
    CONSTRAINT uq_perfil_de_persona_nombre UNIQUE (evento_id, nombre)
) DROP TABLE IF EXISTS tarifa;
CREATE TABLE tarifa (
    tarifa_id BIGSERIAL PRIMARY KEY,
    sector_id BIGINT NOT NULL,
    tipo_de_ticket_id BIGINT NOT NULL,
    perfil_de_persona_id BIGINT,
    precio NUMERIC(4, 2) NOT NULL DEFAULT 0 CONSTRAINT fk_tarifa_sector FOREIGN KEY (sector_id) REFERENCES sector(sector_id) ON DELETE RESTRICT,
    CONSTRAINT fk_tarifa_tipo_de_ticket FOREIGN KEY (tipo_de_ticket_id) REFERENCES tipo_de_ticket(tipo_de_ticket_id) ON DELETE RESTRICT,
    CONSTRAINT fk_tarifa_perfil_de_persona FOREIGN KEY (perfil_de_persona_id) REFERENCES perfil_de_persona(perfil_de_persona_id) ON DELETE RESTRICT,
);
DROP TABLE IF EXISTS ticket;
CREATE TABLE ticket (ticket_id) -- Índice para búsquedas case-insensitive por nombre (opcional)
CREATE UNIQUE INDEX IF NOT EXISTS idx_categorias_nombre_ci ON categorias (LOWER(nombre));