package domain

import (
	"github.com/shopspring/decimal"
)

type Tarifa struct {
	ID              int64            `db:"tarifa_id" json:"tarifaId"`
	Sector          Sector           `db:"-"         json:"sector"`                    // ← objeto, sin ID duplicado
	TipoDeTicket    TipoDeTicket     `db:"-"         json:"tipoDeTicket"`              // ← objeto, sin ID duplicado
	PerfilDePersona *PerfilDePersona `db:"-"         json:"perfilDePersona,omitempty"` // ← opcional (puede ser nil)
	Precio          decimal.Decimal  `db:"precio" json:"precio"`                       // o int64 en centavos
	Activo          int16            `db:"activo" json:"activo"`
}
