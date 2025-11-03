package domain

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

type Cupon struct {
	ID                    int64           `db:"cupon_id" json:"id"`
	Descripcion           string          `db:"descripcion" json:"descripcion"`
	Tipo                  string          `db:"tipo" json:"tipo"`
	Valor                 decimal.Decimal `db:"valor" json:"valor"`
	EstadoCupon           int16           `db:"estado_cupon" json:"estadoCupon"`
	Codigo                string          `db:"codigo" json:"codigo"`
	UsoPorUsuario         int64           `db:"uso_por_usuario" json:"usoPorUsuario"`
	UsoRealizados         int64           `db:"uso_realizados" json:"usoRealizados"`
	UsuarioCreacionID     sql.NullInt64   `db:"usuario_creacion" json:"usuarioCreacionId,omitempty"`
	FechaCreacion         time.Time       `db:"fecha_creacion" json:"fechaCreacion"`
	UsuarioModificacionID sql.NullInt64   `db:"usuario_modificacion" json:"usuarioModificacionId,omitempty"`
	FechaModificacion     sql.NullTime    `db:"fecha_modificacion" json:"fechaModificacion,omitempty"`
}
