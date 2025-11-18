package model

import "time"

type Cupon struct {
	ID                  int64 `gorm:"column:cupon_id;primaryKey;autoIncrement"`
	Descripcion         string
	Tipo                int16
	Valor               float64
	EstadoCupon         int16  `gorm:"default:0"`
	Codigo              string `gorm:"uniqueIndex:uq_cupon_evento"` // único global; si quieres por evento, usa uniqueIndex combinado
	UsoPorUsuario       int64  `gorm:"default:0"`
	UsoRealizados       int64  `gorm:"default:0"`
	FechaInicio         time.Time
	FechaFin            time.Time
	UsuarioCreacion     *int64
	FechaCreacion       time.Time `gorm:"default:now()"`
	UsuarioModificacion *int64
	FechaModificacion   *time.Time

	// FK al evento (muchos cupones pertenecen a un evento)
	EventoID int64   `gorm:"uniqueIndex:uq_cupon_evento"`
	Evento   *Evento `gorm:"foreignKey:EventoID;references:ID"`

	// Mantienes tu relación con usuarios
	Usuarios []UsuarioCupon
}

func (Cupon) TableName() string { return "cupon" }
