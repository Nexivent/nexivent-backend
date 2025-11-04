package model

import "time"

type Notificacion struct {
	ID                 int64 `gorm:"column:notificacion_id;primaryKey;autoIncrement"`
	Mensaje            string
	Canal              string
	FechaEnvio         time.Time
	EstadoNotificacion int16
}

func (Notificacion) TableName() string { return "notificacion" }
