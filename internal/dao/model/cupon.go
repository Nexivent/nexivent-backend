package model

import (
	"time"
)

type Cupon struct {
	ID                  int64 `gorm:"column:cupon_id;primaryKey;autoIncrement"`
	Descripcion         string
	Tipo                string
	Valor               float64
	EstadoCupon         int16
	Codigo              string
	UsoPorUsuario       int64
	UsoRealizados       int64
	UsuarioCreacion     int64
	FechaCreacion       time.Time
	UsuarioModificacion *int64
	FechaModificacion   *time.Time

	EventoCupones []EventoCupon
	Usuarios      []UsuarioCupon
}

func (Cupon) TableName() string { return "cupon" }
