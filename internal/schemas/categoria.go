package schemas

import (
	"github.com/Nexivent/nexivent-backend/internal/dao/model"
)

type ListaCategorias struct{
	
	categorias   []*model.Categoria `json:"categorias"`
}

type CategoriaRequest struct {
	//ID          int64  `gorm:"column:id_categoria;primaryKey;autoIncrement"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	Estado      int16  `json:"estado"`

	//Eventos []Evento
}




type CategoriaResponse struct {
	ID          int64  `json:"id"`
	Nombre      string `json:"nombre"`
	//Descripcion string `gorm:"default:''"`
	//Estado      int16  `gorm:"default:1"`

	//Eventos []Evento
}