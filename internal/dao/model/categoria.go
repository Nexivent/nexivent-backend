package model

type Categoria struct {
	ID          int64  `gorm:"column:id_categoria;primaryKey;autoIncrement"`
	Nombre      string `gorm:"uniqueIndex"`
	Descripcion string `gorm:"default:''"`
	Estado      int16  `gorm:"default:1"`

	Eventos []Evento
}

func (Categoria) TableName() string { return "categoria" }
