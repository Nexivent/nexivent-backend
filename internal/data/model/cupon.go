package model

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/util"
	"github.com/Nexivent/nexivent-backend/internal/validator"
)

type Cupon struct {
	ID                  uint64      `gorm:"column:cupon_id;primaryKey" json:"id"`
	Descripcion         string      `gorm:"column:descripcion" json:"descripcion"`
	Tipo                string      `gorm:"column:tipo" json:"tipo"`
	Valor               float64     `gorm:"column:valor" json:"valor"`
	EstadoCupon         util.Estado `gorm:"column:estado_cupon" json:"estadoCupon"`
	Codigo              string      `gorm:"column:codigo" json:"codigo"`
	UsoPorUsuario       uint64      `gorm:"column:uso_por_usuario" json:"usoPorUsuario"`
	UsoRealizados       uint64      `gorm:"column:uso_realizados" json:"usoRealizados"`
	UsuarioCreacion     *uint64     `gorm:"column:usuario_creacion" json:"-"`
	FechaCreacion       time.Time   `gorm:"column:fecha_creacion" json:"-"`
	UsuarioModificacion *uint64     `gorm:"column:usuario_modificacion" json:"-"`
	FechaModificacion   *time.Time  `gorm:"column:fecha_modificacion" json:"-"`
	EventoID            uint64      `gorm:"column:evento_id" json:"eventoId"`

	Usuarios []UsuarioCupon `json:"usuarios"`
}

func (Cupon) TableName() string { return "cupon" }

func ValidateCupon(v *validator.Validator, cupon *Cupon) {
	// Validar Descripcion
	v.Check(cupon.Descripcion != "", "descripcion", "la descripción es obligatoria")
	v.Check(len(cupon.Descripcion) <= 500, "descripcion", "la descripción no debe exceder 500 caracteres")

	// Validar Tipo
	v.Check(cupon.Tipo != "", "tipo", "el tipo es obligatorio")
	v.Check(len(cupon.Tipo) <= 50, "tipo", "el tipo no debe exceder 50 caracteres")

	// Validar Valor
	v.Check(cupon.Valor >= 0, "valor", "el valor no puede ser negativo")

	// Validar Codigo
	v.Check(cupon.Codigo != "", "codigo", "el código es obligatorio")
	v.Check(len(cupon.Codigo) <= 50, "codigo", "el código no debe exceder 50 caracteres")

	v.Check(cupon.EventoID > 0, "eventoId", "el ID del evento asociado debe ser mayor que cero")
}
