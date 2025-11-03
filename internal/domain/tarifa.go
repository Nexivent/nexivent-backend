package domain

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

type Tarifa struct {
	ID                int64            `db:"tarifa_id" json:"id"`
	Sector            Sector           `db:"-" json:"sector"`                    // FK -> sector
	TipoDeTicket      TipoDeTicket     `db:"-" json:"tipoDeTicket"`              // FK -> tipo_de_ticket
	PerfilDePersona   *PerfilDePersona `db:"-" json:"perfilDePersona,omitempty"` // FK opcional
	Precio            decimal.Decimal  `db:"precio" json:"precio"`
	Activo            int16            `db:"activo" json:"activo"`
	UsuarioCreacionID sql.NullInt64    `db:"usuario_creacion" json:"usuarioCreacionId,omitempty"`
	FechaCreacion     time.Time        `db:"fecha_creacion" json:"fechaCreacion"`
	UsuarioModID      sql.NullInt64    `db:"usuario_modificacion" json:"usuarioModificacionId,omitempty"`
	FechaModificacion sql.NullTime     `db:"fecha_modificacion" json:"fechaModificacion,omitempty"`
}
