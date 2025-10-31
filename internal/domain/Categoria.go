package domain

type Categoria struct {
	IDCategoria int64  `db:"id_categoria" json:"idCategoria"`
	Nombre      string `db:"nombre"        json:"nombre"`
	Descripcion string `db:"descripcion"   json:"descripcion"`
}
