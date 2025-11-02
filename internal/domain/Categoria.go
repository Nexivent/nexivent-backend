package domain

import (
	"context"
)

type Categoria struct {
	IDCategoria int    `db:"id_categoria" json:"idCategoria"`
	Nombre      string `db:"nombre"        json:"nombre"`
	Descripcion string `db:"descripcion"   json:"descripcion"`
}

type CategoriaRepository interface {
	Save(cont context.Context, c *Categoria) error
	GetById(cont context.Context, id int64) (*Categoria, error)
	Delete(cont context.Context, id int64) error
}
