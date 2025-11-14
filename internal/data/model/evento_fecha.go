package model

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/model/util"
	"github.com/Nexivent/nexivent-backend/internal/validator"
)

type EventoFecha struct {
	ID                  uint64      `gorm:"column:evento_fecha_id;primaryKey;autoIncrement" json:"id"`
	EventoID            uint64      `gorm:"column:evento_id;uniqueIndex:uq_evento_fecha" json:"eventoId"`
	FechaID             uint64      `gorm:"column:fecha_id;uniqueIndex:uq_evento_fecha" json:"fechaId"`
	HoraInicio          time.Time   `gorm:"column:hora_inicio;uniqueIndex:uq_evento_fecha" json:"horaInicio"`
	Estado              util.Estado `gorm:"column:estado" json:"-"`
	UsuarioCreacion     *uint64     `gorm:"column:usuario_creacion" json:"-"`
	FechaCreacion       time.Time   `gorm:"column:fecha_creacion;default:now()" json:"-"`
	UsuarioModificacion *uint64     `gorm:"column:usuario_modificacion" json:"-"`
	FechaModificacion   *time.Time  `gorm:"column:fecha_modificacion" json:"-"`

	// Evento *Evento `gorm:"foreignKey:EventoID;references:evento_id"`
	// Fecha  *Fecha  `gorm:"foreignKey:FechaID;references:fecha_id"`

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
