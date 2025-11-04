package domain

type Categoria struct {
	ID          int64 `gorm:"column:id_categoria;primaryKey;autoIncrement"`
	Nombre      string
	Descripcion string
	Estado      int16
}

func (Categoria) TableName() string { return "categoria" }
