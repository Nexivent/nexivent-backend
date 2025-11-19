package adapter

import (
	//"fmt"
	//"time"

	"github.com/Nexivent/nexivent-backend/errors"
	"github.com/Nexivent/nexivent-backend/internal/dao/model"
	daoPostgresql "github.com/Nexivent/nexivent-backend/internal/dao/repository"
	"github.com/Nexivent/nexivent-backend/internal/schemas"
	"github.com/Nexivent/nexivent-backend/logging"
	//"github.com/Loui27/nexivent-backend/utils/convert"
	//"gorm.io/gorm"
)

type Categoria struct {
	logger        logging.Logger
	DaoPostgresql *daoPostgresql.NexiventPsqlEntidades
}

// Creates Evento adapter
func NewCategoriaAdapter(
	logger logging.Logger,
	daoPostgresql *daoPostgresql.NexiventPsqlEntidades,
) *Categoria {
	return &Categoria{
		logger:        logger,
		DaoPostgresql: daoPostgresql,
	}
}

// CreatePostgresqlEvento creates a new event with all related entities
func (e *Categoria) CreatePostgresqlCategoria(categoriaReq *schemas.CategoriaRequest) (*schemas.CategoriaResponse, *errors.Error) {
	// Start a transaction
	tx := e.DaoPostgresql.Categoria.PostgresqlDB.Begin()
	if tx.Error != nil {
		e.logger.Errorf("Failed to begin transaction: %v", tx.Error)
		return nil, &errors.BadRequestError.CategoriaNotCreated
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create the main event model
	categoriaModel := &model.Categoria{
		Descripcion:     categoriaReq.Descripcion,
		Estado:          categoriaReq.Estado,
		Nombre:          categoriaReq.Nombre,
	}

	// Create the event
	if err := tx.Create(categoriaModel).Error; err != nil {
		tx.Rollback()
		e.logger.Errorf("Failed to create event: %v", err)
		return nil, &errors.BadRequestError.CategoriaNotCreated
	}

	if err := tx.Commit().Error; err!=nil{
		e.logger.Errorf("Failed to commit categoria: %v", err)
		return nil, &errors.BadRequestError.CategoriaNotCreated
	}
	
	// Build the response
	response := &schemas.CategoriaResponse{
		ID: categoriaModel.ID,
		Nombre: categoriaModel.Nombre,
		//Descripcion: categoriaModel.Descripcion,

		
	}

	return response, nil
}

// FetchPostgresqlEventos retrieves the upcoming events without filters
func (e *Categoria) FetchPostgresqlCategorias() ([] schemas.CategoriaResponse, *errors.Error) {
	categorias, err := e.DaoPostgresql.Categoria.ObtenerCategorias()
	if err != nil {
		e.logger.Errorf("Failed to fetch categorias: %v", err)
		return nil, &errors.BadRequestError.CategoriaNotFound
	}
	categoriasResponse := make([]schemas.CategoriaResponse, len(categorias))

	for i, c := range categorias{
		categoriasResponse[i] = schemas.CategoriaResponse{
			ID: c.ID,
			Nombre: c.Nombre,
		}
	}

	return categoriasResponse, nil
}

// GetPostgresqlEventoById gets an event by ID

func (e *Categoria) GetPostgresqlCategoriaById(categoriaID int64) (*schemas.CategoriaResponse, *errors.Error) {
	var categoriaModel *model.Categoria

	categoriaModel, err := e.DaoPostgresql.Categoria.ObtenerCategoriaPorId(categoriaID)
	if err != nil {
		e.logger.Errorf("Failed to fetch categorias: %v", err)
		return nil, &errors.BadRequestError.EventoNotFound
	}

	response := &schemas.CategoriaResponse{
		ID:          categoriaModel.ID,
		Nombre:     categoriaModel.Nombre,
	}

	return response, nil
}