package domain

import "time"

type Fecha struct {
	ID          int64     `db:"fecha_id"     json:"fechaId"`
	FechaEvento time.Time `db:"fecha_evento" json:"fechaEvento"` // DATE
}
