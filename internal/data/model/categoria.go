package model

import (
	"github.com/Nexivent/nexivent-backend/internal/data/model/util"
	"github.com/Nexivent/nexivent-backend/internal/validator"
)

type Categoria struct {
	ID          uint64      `gorm:"column:categoria_id;primaryKey;autoIncrement" json:"id"`
	Nombre      string      `gorm:"column:nombre;uniqueIndex" json:"nombre"`
	Descripcion string      `gorm:"column:descripcion;default:''" json:"descripcion"`
	Estado      util.Estado `gorm:"column:estado;default:1" json:"-"`
	
	Eventos []Evento `json:"eventos"`
}

func (Categoria) TableName() string { return "categoria" }

func ValidateCategoria(v *validator.Validator, categoria *Categoria) {
	// Validar Nombre
	v.Check(categoria.Nombre != "", "nombre", "el nombre es obligatorio")
	v.Check(len(categoria.Nombre) <= 100, "nombre", "el nombre no debe exceder 100 caracteres")

	// Validar Descripción
	v.Check(len(categoria.Descripcion) <= 500, "descripcion", "la descripción no debe exceder 500 caracteres")
}
