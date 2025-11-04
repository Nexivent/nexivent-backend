package domain

import "time"

type Comentario struct {
	ID            int64 `gorm:"column:comentario_id;primaryKey;autoIncrement"`
	UsuarioID     int64
	EventoID      int64
	Descripcion   string
	FechaCreacion time.Time
	Estado        int16

	Usuario *Usuario `gorm:"foreignKey:UsuarioID"`
	Evento  *Evento  `gorm:"foreignKey:EventoID"`
}

func (Comentario) TableName() string { return "comentario" }
