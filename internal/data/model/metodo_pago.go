package model

import (
	"github.com/Nexivent/nexivent-backend/internal/data/util"
	"github.com/Nexivent/nexivent-backend/internal/validator"
)

type MetodoDePago struct {
	ID     uint64      `gorm:"column:metodo_de_pago_id;primaryKey" json:"id"`
	Tipo   string      `gorm:"column:tipo" json:"tipo"`
	Estado util.Estado `gorm:"column:estado" json:"-"`

	Ordenes []OrdenDeCompra `json:"ordenes,omitempty"`
}

func (MetodoDePago) TableName() string { return "metodo_de_pago" }

func ValidateMetodoDePago(v *validator.Validator, metodo *MetodoDePago) {
	// Validar Tipo
	v.Check(metodo.Tipo != "", "tipo", "el tipo de método de pago es obligatorio")
	v.Check(metodo.Tipo == "Tarjeta" || metodo.Tipo == "Yape", "tipo", "el tipo debe ser 'Tarjeta' o 'Yape'")
}
