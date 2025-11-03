package domain

import "time"

type EventoCupon struct {
	Evento      Evento    `db:"-"          json:"evento"` // FK -> evento (siempre)
	Cupon       Cupon     `db:"-"          json:"cupon"`  // FK -> cupon  (siempre)
	CantCupones int64     `db:"cant_cupones" json:"cantCupones"`
	FechaIni    time.Time `db:"fecha_ini"  json:"fechaIni"` // DATE
	FechaFin    time.Time `db:"fecha_fin"  json:"fechaFin"` // DATE
}
