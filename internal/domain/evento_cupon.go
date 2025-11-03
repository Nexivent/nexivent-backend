package domain

import (
	"database/sql"
	"time"
)

type EventoCupon struct {
	Evento            Evento        `db:"-" json:"evento"` // PK compuesta: (evento, cupon)
	Cupon             Cupon         `db:"-" json:"cupon"`
	CantCupones       int64         `db:"cant_cupones" json:"cantCupones"`
	FechaIni          time.Time     `db:"fecha_ini" json:"fechaIni"`
	FechaFin          time.Time     `db:"fecha_fin" json:"fechaFin"`
	UsuarioCreacionID sql.NullInt64 `db:"usuario_creacion" json:"usuarioCreacionId,omitempty"`
	FechaCreacion     time.Time     `db:"fecha_creacion" json:"fechaCreacion"`
	UsuarioModID      sql.NullInt64 `db:"usuario_modificacion" json:"usuarioModificacionId,omitempty"`
	FechaModificacion sql.NullTime  `db:"fecha_modificacion" json:"fechaModificacion,omitempty"`
}
