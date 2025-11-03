package domain

import (
	"database/sql"
	"time"
)

type Sector struct {
	ID                int64         `db:"sector_id" json:"id"`
	Evento            Evento        `db:"-" json:"evento"` // FK -> evento
	SectorTipo        string        `db:"sector_tipo" json:"sectorTipo"`
	TotalEntradas     int           `db:"total_entradas" json:"totalEntradas"`
	CantVendidas      int           `db:"cant_vendidas" json:"cantVendidas"`
	Activo            int16         `db:"activo" json:"activo"`
	UsuarioCreacionID sql.NullInt64 `db:"usuario_creacion" json:"usuarioCreacionId,omitempty"`
	FechaCreacion     time.Time     `db:"fecha_creacion" json:"fechaCreacion"`
	UsuarioModID      sql.NullInt64 `db:"usuario_modificacion" json:"usuarioModificacionId,omitempty"`
	FechaModificacion sql.NullTime  `db:"fecha_modificacion" json:"fechaModificacion,omitempty"`
}
