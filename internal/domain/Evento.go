package domain

import "time"

type Evento struct {
	IDEvento        int       `db:"id_evento" json:"idEvento"`
	Organizador     Usuario   `db:"-" json:"organizador"`
	Categoria       Categoria `db:"-" json:"categoria"`
	Titulo          string    `db:"titulo" json:"titulo"`
	Descripcion     string    `db:"descripcion" json:"descripcion"`
	Lugar           string    `db:"lugar" json:"lugar"`
	FechaHoraInicio time.Time `db:"fecha_hora_inicio" json:"fechaHoraInicio"`
	FechaHoraFin    time.Time `db:"fecha_hora_fin" json:"fechaHoraFin"`
	EstadoEvento    Estado    `db:"estado_evento" json:"estadoEvento"`
	CantLikes       int       `db:"cant_likes" json:"cantLikes"`
}
