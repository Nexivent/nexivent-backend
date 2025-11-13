package model

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/util"
	"github.com/Nexivent/nexivent-backend/internal/validator"
)

type Notificacion struct {
	ID                 uint64                  `gorm:"column:notificacion_id;primaryKey" json:"id"`
	Mensaje            string                  `gorm:"column:mensaje" json:"mensaje"`
	Canal              string                  `gorm:"column:canal" json:"canal"`
	FechaEnvio         time.Time               `gorm:"column:fecha_envio" json:"fechaEnvio"`
	EstadoNotificacion util.EstadoNotificacion `gorm:"column:estado_notificacion" json:"estadoNotificacion"`
}

func (Notificacion) TableName() string { return "notificacion" }

func ValidateNotificacion(v *validator.Validator, notificacion *Notificacion) {
	// Validar Mensaje
	v.Check(notificacion.Mensaje != "", "mensaje", "el mensaje es obligatorio")
	v.Check(len(notificacion.Mensaje) <= 1000, "mensaje", "el mensaje no debe exceder 1000 caracteres")

	// Validar Canal
	v.Check(notificacion.Canal != "", "canal", "el canal es obligatorio")
	v.Check(len(notificacion.Canal) <= 50, "canal", "el canal no debe exceder 50 caracteres")

	// Validar FechaEnvio
	v.Check(!notificacion.FechaEnvio.IsZero(), "fechaEnvio", "la fecha de envío es obligatoria")
}
