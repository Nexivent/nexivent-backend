package domain

import "time"

type EventoFecha struct {
	ID         int64     `db:"evento_fecha_id" json:"eventoFechaId"`
	Evento     Evento    `db:"-"               json:"evento"` // FK -> evento (siempre)
	Fecha      Fecha     `db:"-"               json:"fecha"`  // FK -> fecha  (siempre)
	HoraInicio time.Time `db:"hora_inicio"     json:"horaInicio"`
	Activo     int16     `db:"activo" json:"activo"`
}
