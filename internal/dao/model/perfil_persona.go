package model

import (
	"time"
)

type PerfilDePersona struct {
	ID                  int64  `gorm:"column:perfil_de_persona_id;primaryKey;autoIncrement"`
	EventoID            int64  `gorm:"uniqueIndex:uq_perfil_de_persona_nombre"`
	Nombre              string `gorm:"uniqueIndex:uq_perfil_de_persona_nombre"`
	Estado              int16  `gorm:"default:1"`
	UsuarioCreacion     *int64
	FechaCreacion       time.Time `gorm:"default:now()"`
	UsuarioModificacion *int64
	FechaModificacion   *time.Time

	Evento *Evento `gorm:"foreignKey:EventoID;references:evento_id"`
}

func (PerfilDePersona) TableName() string { return "perfil_de_persona" }
