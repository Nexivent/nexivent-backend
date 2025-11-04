package domain

import (
	"time"
)

type PerfilDePersona struct {
	ID                  int64 `gorm:"column:perfil_de_persona_id;primaryKey;autoIncrement"`
	EventoID            int64
	Nombre              string
	Estado              int16
	UsuarioCreacion     *int64
	FechaCreacion       time.Time
	UsuarioModificacion *int64
	FechaModificacion   *time.Time

	Evento *Evento `gorm:"foreignKey:EventoID"`
}

func (PerfilDePersona) TableName() string { return "perfil_de_persona" }
