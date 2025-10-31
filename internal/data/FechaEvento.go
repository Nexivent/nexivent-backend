package data

import "time"

type FechaEvento struct {
	IDFechaEvento int       `db:"id_fecha_evento" json:"idFechaEvento"`
	Evento        Evento    `db:"-" json:"evento"`
	FechaDeEvento time.Time `db:"fecha_de_evento" json:"fechaDeEvento"`
	HoraDeInicio  time.Time `db:"hora_inicio" json:"horaDeInicio"`
	HoraFin       time.Time `db:"hora_fin" json:"horaFin"`
}
