package model

type Categoria struct {
	ID          int64 `gorm:"column:categoria_id;primaryKey;autoIncrement"`
	Nombre      string 
	Descripcion string
	Estado      int16
}

func (Categoria) TableName() string { return "categoria" }
