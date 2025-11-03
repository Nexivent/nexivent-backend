package domain

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

type Evento struct {
	ID                int64           `db:"evento_id" json:"id"`
	Organizador       Usuario         `db:"-" json:"organizador"` // FK -> usuario
	Categoria         Categoria       `db:"-" json:"categoria"`   // FK -> categoria
	Titulo            string          `db:"titulo" json:"titulo"`
	Descripcion       string          `db:"descripcion" json:"descripcion"`
	Lugar             string          `db:"lugar" json:"lugar"`
	EventoEstado      EstadoEvento    `db:"evento_estado"`
	CantMeGusta       int             `db:"cant_me_gusta" json:"cantMeGusta"`
	CantNoInteresa    int             `db:"cant_no_interesa" json:"cantNoInteresa"`
	CantVendidoTotal  int             `db:"cant_vendido_total" json:"cantVendidoTotal"`
	TotalRecaudado    decimal.Decimal `db:"total_recaudado" json:"totalRecaudado"`
	Activo            int16           `db:"activo" json:"activo"`
	UsuarioCreacionID sql.NullInt64   `db:"usuario_creacion" json:"usuarioCreacionId,omitempty"`
	FechaCreacion     time.Time       `db:"fecha_creacion" json:"fechaCreacion"`
	UsuarioModID      sql.NullInt64   `db:"usuario_modificacion" json:"usuarioModificacionId,omitempty"`
	FechaModificacion sql.NullTime    `db:"fecha_modificacion" json:"fechaModificacion,omitempty"`
}
