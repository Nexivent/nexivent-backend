package model

type Ticket struct {
	ID              int64 `gorm:"column:ticket_id;primaryKey;autoIncrement"`
	OrdenDeCompraID *int64
	EventoFechaID   int64
	TarifaID        int64
	CodigoQR        string `gorm:"uniqueIndex"`
	EstadoDeTicket  int16  `gorm:"default:0"`

	OrdenDeCompra *OrdenDeCompra `gorm:"foreignKey:OrdenDeCompraID;references:orden_de_compra_id"`
	EventoFecha   *EventoFecha   `gorm:"foreignKey:EventoFechaID;references:evento_fecha_id"`
	Tarifa        *Tarifa        `gorm:"foreignKey:TarifaID;references:tarifa_id"`
}

func (Ticket) TableName() string { return "ticket" }
