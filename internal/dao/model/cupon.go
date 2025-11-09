package model

import (
	"time"
)

type Cupon struct {
	ID                  int64 `gorm:"column:cupon_id;primaryKey;autoIncrement"`
	EventoID            int64
	Descripcion         string
	Tipo                string
	Valor               float64
	EstadoCupon         int16
	Codigo              string `gorm:"uniqueIndex:idx_codigo_unico"`
	FechaIni            time.Time
	FechaFin            time.Time
	UsoPorUsuario       int64
	UsoRealizados       int64
	UsuarioCreacion     int64
	FechaCreacion       time.Time
	UsuarioModificacion *int64
	FechaModificacion   *time.Time

	Evento   *Evento `gorm:"foreignKey:EventoID"`
	Usuarios []UsuarioCupon
}

func (Cupon) TableName() string { return "cupon" }
