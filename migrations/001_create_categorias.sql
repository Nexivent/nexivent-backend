CREATE TABLE IF NOT EXISTS categorias (
    id_categoria BIGSERIAL PRIMARY KEY,
    -- int64 en Go
    nombre VARCHAR(80) NOT NULL UNIQUE,
    descripcion TEXT NOT NULL DEFAULT ''
);
-- Índice para búsquedas case-insensitive por nombre (opcional)
CREATE UNIQUE INDEX IF NOT EXISTS idx_categorias_nombre_ci ON categorias (LOWER(nombre));