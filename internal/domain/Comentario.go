package domain

import "time"

type Comentario struct {
	ID            int64     `db:"comentario_id" json:"id"`
	Usuario       Usuario   `db:"-" json:"usuario"` // FK -> usuario
	Evento        Evento    `db:"-" json:"evento"`  // FK -> evento
	Descripcion   string    `db:"descripcion" json:"descripcion"`
	FechaCreacion time.Time `db:"fecha_creacion" json:"fechaCreacion"`
	Activo        int16     `db:"activo" json:"activo"`
}
