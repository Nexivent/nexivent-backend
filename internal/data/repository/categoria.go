package repository

import (
	"github.com/Nexivent/nexivent-backend/internal/data/model"
	"gorm.io/gorm"
)

type Categoria struct {
	DB *gorm.DB
}

func (c *Categoria) CrearCategoria(Categoria *model.Categoria) error {
	respuesta := c.DB.Create(Categoria)
	if respuesta.Error != nil {
		return respuesta.Error
	}
	return nil
}

func (c *Categoria) ObtenerCategorias() ([]*model.Categoria, error) {
	var categorias []*model.Categoria
	respuesta := c.DB.Find(&categorias)

	if respuesta.Error != nil {
		return nil, respuesta.Error
	}

	return categorias, nil
}
