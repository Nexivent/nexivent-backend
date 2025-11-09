package repository

import (
	"github.com/Loui27/nexivent-backend/internal/dao/model"
	"github.com/Loui27/nexivent-backend/logging"
	"gorm.io/gorm"
)

type Categoria struct {
	logger       logging.Logger
	PostgresqlDB *gorm.DB
}

func NewCategoriaController(
	logger logging.Logger,
	postgresqlDB *gorm.DB,
) *Categoria {
	return &Categoria{
		logger:       logger,
		PostgresqlDB: postgresqlDB,
	}
}

func (c *Categoria) CrearCategoria(Categoria *model.Categoria) error {
	respuesta := c.PostgresqlDB.Create(Categoria)
	if respuesta.Error != nil {
		return respuesta.Error
	}
	return nil
}

func (c *Categoria) ObtenerCategorias() ([]*model.Categoria, error) {
	var categorias []*model.Categoria
	respuesta := c.PostgresqlDB.Find(&categorias)

	if respuesta.Error != nil {
		return nil, respuesta.Error
	}

	return categorias, nil
}
