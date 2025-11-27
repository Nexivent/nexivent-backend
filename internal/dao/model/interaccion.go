package model

import "time"

type Interaccion struct {
	ID            int64     `gorm:"column:interaccion_id;primaryKey;autoIncrement"`
	UsuarioID     int64     `gorm:"index:idx_usuario_evento,unique"`
	EventoID      int64     `gorm:"index:idx_usuario_evento,unique"`
	Tipo          int64     `gorm:"default:0"`
	FechaCreacion time.Time `gorm:"default:now()"`
	Estado        int16     `gorm:"default:1"`

	Usuario *Usuario `gorm:"foreignKey:UsuarioID;references:usuario_id"`
	Evento  *Evento  `gorm:"foreignKey:EventoID;references:evento_id"`
}

func (Interaccion) TableName() string { return "interaccion" }
