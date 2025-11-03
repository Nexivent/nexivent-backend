package domain

import "time"

type TipoDeTicket struct {
	ID       int64     `db:"tipo_de_ticket_id" json:"tipoDeTicketId"`
	Evento   Evento    `db:"-" json:"evento"`
	Nombre   string    `db:"nombre" json:"nombre"`
	FechaIni time.Time `db:"fecha_ini" json:"fechaIni"`
	FechaFin time.Time `db:"fecha_fin" json:"fechaFin"`
	Activo   int16     `db:"activo" json:"activo"`
}
