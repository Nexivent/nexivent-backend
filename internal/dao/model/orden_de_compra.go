package model

import (
	"time"
)

type OrdenDeCompra struct {
	ID               int64 `gorm:"column:orden_de_compra_id;primaryKey;autoIncrement"`
	UsuarioID        int64
	MetodoDePagoID   int64
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
	// Campos calculados/virtuales (no se persisten en BD)
    PrecioEntrada      float64    `gorm:"-" json:"precio_entrada,omitempty"`
    TicketID           *int64     `gorm:"-" json:"ticket_id,omitempty"`
}

func (OrdenDeCompra) TableName() string { return "orden_de_compra" }
