package data

import (
	"time"

	"github.com/google/uuid"
)

type OrdenDeCompra struct {
	ID               uuid.UUID `gorm:"column:orden_de_compra_id;primaryKey"`
	UsuarioID        uuid.UUID
	MetodoDePagoID   uuid.UUID
	Fecha            time.Time `gorm:"default:current_date"`
	FechaHoraIni     time.Time `gorm:"default:now()"`
	FechaHoraFin     *time.Time
	Total            float64
	MontoFeeServicio float64
	EstadoDeOrden    int16 `gorm:"default:0"`

	Usuario      *Usuario      `gorm:"foreignKey:UsuarioID;references:usuario_id"`
	MetodoDePago *MetodoDePago `gorm:"foreignKey:MetodoDePagoID;references:metodo_de_pago_id"`

	Tickets          []Ticket
	ComprobantesPago []ComprobanteDePago
}


