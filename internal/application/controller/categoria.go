package controller

import (
	"github.com/Loui27/nexivent-backend/errors"
	"github.com/Loui27/nexivent-backend/internal/application/adapter"
	"github.com/Loui27/nexivent-backend/internal/schemas"
	"github.com/Loui27/nexivent-backend/logging"
)

type CategoriaController struct {
	Logger        logging.Logger
	CategoriaAdapter *adapter.Categoria
}

// NewEventoController creates a new controller for event operations
func NewCategoriaController(
	logger logging.Logger,
	categoriaAdapter *adapter.Categoria,
) *CategoriaController {
	return &CategoriaController{
		Logger:        logger,
		CategoriaAdapter: categoriaAdapter,
	}
}

// CreateEvento creates a new event with all related entities
func (ec *CategoriaController) CreateCategoria(
	categoriaReq schemas.CategoriaRequest,
) (*schemas.CategoriaResponse, *errors.Error) {
	return ec.CategoriaAdapter.CreatePostgresqlCategoria(&categoriaReq)
}

// FetchEventos retrieves the list of available events
func (ec *CategoriaController) FetchCategorias() ([]schemas.CategoriaResponse, *errors.Error) {
	return ec.CategoriaAdapter.FetchPostgresqlCategorias()
}

// GetEventoById retrieves an event by its ID with all related entities

func (ec *CategoriaController) GetCategoriaById(eventoID int64) (*schemas.CategoriaResponse, *errors.Error) {
	return ec.CategoriaAdapter.GetPostgresqlCategoriaById(eventoID)
}
