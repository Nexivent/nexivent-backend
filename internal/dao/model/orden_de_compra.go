package model

import (
	"time"
)

type OrdenDeCompra struct {
	ID               int64 `gorm:"column:orden_de_compra_id;primaryKey;autoIncrement"`
	UsuarioID        int64
	MetodoDePagoID   int64
	Fecha            time.Time
	FechaHoraIni     time.Time
	FechaHoraFin     time.Time
	Total            float64
	MontoFeeServicio float64
	EstadoDeOrden    int16

	Usuario      *Usuario      `gorm:"foreignKey:UsuarioID"`
	MetodoDePago *MetodoDePago `gorm:"foreignKey:MetodoDePagoID"`

	Tickets          []Ticket
	ComprobantesPago []ComprobanteDePago
}

func (OrdenDeCompra) TableName() string { return "orden_de_compra" }
