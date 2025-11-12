package model

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/util"
)

type EventoFecha struct {
	ID                  uint64      `gorm:"column:evento_fecha_id;primaryKey" json:"id"`
	EventoID            uint64      `gorm:"column:evento_id" json:"eventoId"`
	FechaID             uint64      `gorm:"column:fecha_id" json:"fechaId"`
	HoraInicio          time.Time   `gorm:"column:hora_inicio" json:"horaInicio"`
	Estado              util.Estado `gorm:"column:estado" json:"-"`
	UsuarioCreacion     *uint64     `gorm:"column:usuario_creacion" json:"-"`
	FechaCreacion       time.Time   `gorm:"column:fecha_creacion" json:"-"`
	UsuarioModificacion *uint64     `gorm:"column:usuario_modificacion" json:"-"`
	FechaModificacion   *time.Time  `gorm:"column:fecha_modificacion" json:"-"`

	Tickets []Ticket `json:"-"`
}

func (EventoFecha) TableName() string { return "evento_fecha" }
