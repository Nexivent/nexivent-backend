package model

import "time"

type Comentario struct {
	ID            int64 `gorm:"column:comentario_id;primaryKey;autoIncrement"`
	UsuarioID     int64
	EventoID      int64
	Descripcion   string
	FechaCreacion time.Time `gorm:"default:now()"`
	Estado        int16     `gorm:"default:1"`

	Usuario *Usuario `gorm:"foreignKey:UsuarioID;references:usuario_id"`
	Evento  *Evento  `gorm:"foreignKey:EventoID;references:evento_id"`
}

func (Comentario) TableName() string { return "comentario" }
