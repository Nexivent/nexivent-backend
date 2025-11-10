package model

type MetodoDePago struct {
	ID     int64 `gorm:"column:metodo_de_pago_id;primaryKey;autoIncrement"`
	Tipo   string
	Estado int16 `gorm:"default:1"`

	Ordenes []OrdenDeCompra
}

func (MetodoDePago) TableName() string { return "metodo_de_pago" }
