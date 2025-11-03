package domain

import (
	"database/sql"
	"time"
)

type Rol struct {
	ID                    int64         `db:"rol_id" json:"id"`
	Nombre                string        `db:"nombre" json:"nombre"`
	UsuarioCreacionID     sql.NullInt64 `db:"usuario_creacion" json:"usuarioCreacionId,omitempty"`
	FechaCreacion         time.Time     `db:"fecha_creacion" json:"fechaCreacion"`
	UsuarioModificacionID sql.NullInt64 `db:"usuario_modificacion" json:"usuarioModificacionId,omitempty"`
	FechaModificacion     sql.NullTime  `db:"fecha_modificacion" json:"fechaModificacion,omitempty"`
}
