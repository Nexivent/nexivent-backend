package model

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/util"
	"github.com/Nexivent/nexivent-backend/internal/validator"
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

	Tickets []Ticket `json:"tickets,omitempty"`
}

func (EventoFecha) TableName() string { return "evento_fecha" }

func ValidateEventoFecha(v *validator.Validator, eventoFecha *EventoFecha) {
	// Validar IDs
	v.Check(eventoFecha.EventoID != 0, "eventoId", "el ID del evento es obligatorio")
	v.Check(eventoFecha.FechaID != 0, "fechaId", "el ID de la fecha es obligatorio")

	// Validar HoraInicio
	v.Check(!eventoFecha.HoraInicio.IsZero(), "horaInicio", "la hora de inicio es obligatoria")
}
