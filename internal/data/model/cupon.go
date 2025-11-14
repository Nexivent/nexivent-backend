package model

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/model/util"
	"github.com/Nexivent/nexivent-backend/internal/validator"
)

type Cupon struct {
	ID                  uint64      `gorm:"column:cupon_id;primaryKey;autoIncrement" json:"id"`
	Descripcion         string      `gorm:"column:descripcion" json:"descripcion"`
	Tipo                string      `gorm:"column:tipo" json:"tipo"`
	Valor               float64     `gorm:"column:valor" json:"valor"`
	EstadoCupon         util.Estado `gorm:"column:estado_cupon;default:0" json:"estadoCupon"`
	Codigo              string      `gorm:"column:codigo;uniqueIndex" json:"codigo"`
	UsoPorUsuario       uint64      `gorm:"column:uso_por_usuario;default:0" json:"usoPorUsuario"`
	UsoRealizados       uint64      `gorm:"column:uso_realizados;default:0" json:"usoRealizados"`
	UsuarioCreacion     *uint64     `gorm:"column:usuario_creacion" json:"-"`
	FechaCreacion       time.Time   `gorm:"column:fecha_creacion;default:now()" json:"-"`
	UsuarioModificacion *uint64     `gorm:"column:usuario_modificacion" json:"-"`
	FechaModificacion   *time.Time  `gorm:"column:fecha_modificacion" json:"-"`

	// FK al evento (muchos cupones pertenecen a un evento)
	EventoID uint64  `gorm:"column:evento_id" json:"eventoId"`
	// Evento   *Evento `gorm:"foreignKey:EventoID;references:ID"`

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
