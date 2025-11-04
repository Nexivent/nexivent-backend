package domain

import "time"

type Fecha struct {
	ID          int64 `gorm:"column:fecha_id;primaryKey;autoIncrement"`
	FechaEvento time.Time
}

func (Fecha) TableName() string { return "fecha" }
