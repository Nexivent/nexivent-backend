package domain

import (
	"context"
)

type Categoria struct {
	ID          int64  `db:"id_categoria" json:"id"`
	Nombre      string `db:"nombre" json:"nombre"`
	Descripcion string `db:"descripcion" json:"descripcion"`
	Activo      int16  `db:"activo" json:"activo"`
}
type CategoriaRepository interface {
	Save(cont context.Context, c *Categoria) error
	GetById(cont context.Context, id int) (*Categoria, error)
	Delete(cont context.Context, id int) error
}
