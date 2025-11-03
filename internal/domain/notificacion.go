package domain

import "time"

type Notificacion struct {
	ID                 int64     `db:"notificacion_id" json:"id"`
	Mensaje            string    `db:"mensaje" json:"mensaje"`
	Canal              string    `db:"canal" json:"canal"`
	FechaEnvio         time.Time `db:"fecha_envio" json:"fechaEnvio"`
	EstadoNotificacion int16     `db:"estado_notificación" json:"estadoNotificacion"` // columna con tilde en DDL
}
