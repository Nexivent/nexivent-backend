package model

type Ticket struct {
	ID              int64 `gorm:"column:ticket_id;primaryKey;autoIncrement"`
	OrdenDeCompraID *int64
	EventoFechaID   int64
	TarifaID        int64
	CodigoQR        string
	EstadoDeTicket  int16

	OrdenDeCompra *OrdenDeCompra `gorm:"foreignKey:OrdenDeCompraID"`
	EventoFecha   *EventoFecha   `gorm:"foreignKey:EventoFechaID"`
	Tarifa        *Tarifa        `gorm:"foreignKey:TarifaID"`
}

func (Ticket) TableName() string { return "ticket" }
