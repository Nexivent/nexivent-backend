package domain

import (
	"database/sql"
	"time"
)

type RolUsuario struct {
	ID                int64         `db:"rol_usuario_id" json:"id"`
	Rol               Rol           `db:"-" json:"rol"`     // FK -> rol
	Usuario           Usuario       `db:"-" json:"usuario"` // FK -> usuario
	UsuarioCreacionID sql.NullInt64 `db:"usuario_creacion" json:"usuarioCreacionId,omitempty"`
	FechaCreacion     time.Time     `db:"fecha_creacion" json:"fechaCreacion"`
	UsuarioModID      sql.NullInt64 `db:"usuario_modificacion" json:"usuarioModificacionId,omitempty"`
	FechaModificacion sql.NullTime  `db:"fecha_modificacion" json:"fechaModificacion,omitempty"`
}
