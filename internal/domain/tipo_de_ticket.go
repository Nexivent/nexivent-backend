package domain

import (
	"database/sql"
	"time"
)

type TipoDeTicket struct {
	ID                int64         `db:"tipo_de_ticket_id" json:"id"`
	Evento            Evento        `db:"-" json:"evento"` // FK -> evento
	Nombre            string        `db:"nombre" json:"nombre"`
	FechaIni          time.Time     `db:"fecha_ini" json:"fechaIni"`
	FechaFin          time.Time     `db:"fecha_fin" json:"fechaFin"`
	Activo            int16         `db:"activo" json:"activo"`
	UsuarioCreacionID sql.NullInt64 `db:"usuario_creacion" json:"usuarioCreacionId,omitempty"`
	FechaCreacion     time.Time     `db:"fecha_creacion" json:"fechaCreacion"`
	UsuarioModID      sql.NullInt64 `db:"usuario_modificacion" json:"usuarioModificacionId,omitempty"`
	FechaModificacion sql.NullTime  `db:"fecha_modificacion" json:"fechaModificacion,omitempty"`
}
