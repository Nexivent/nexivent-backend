package model

import (
	"time"
)

type OrdenDeCompra struct {
	ID               uint64 
	UsuarioID        uint64
	Fecha            time.Time 
	FechaHoraIni     time.Time 
	FechaHoraFin     *time.Time
	Total            float64
	MontoFeeServicio float64
	EstadoDeOrden    int16 

	Tickets          []Ticket
	ComprobantesPago []ComprobanteDePago
}

func (OrdenDeCompra) TableName() string { return "orden_de_compra" }