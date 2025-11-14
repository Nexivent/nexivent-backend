package adapter

import (
	"github.com/Loui27/nexivent-backend//errors"
	"github.com/Loui27/nexivent-backend/internal/dao/model"
	daoPostgresql "github.com/Loui27/nexivent-backend/internal/dao/repository"
	"github.com/Loui27/nexivent-backend/internal/schemas"
	"github.com/Loui27/nexivent-backend/logging"
	"github.com/google/uuid"
)

type Evento struct {
	logger        logging.Logger
	DaoPostgresql *daoPostgresql.NexiventPsqlEntidades
}

// Creates Evento adapter
func NewEventoAdapter(
	logger logging.Logger,
	daoPostgresql *daoPostgresql.NexiventPsqlEntidades,
) *Evento {
	return &Evento{
		logger:        logger,
		DaoPostgresql: daoPostgresql,
	}
}

// Gets a Local from postgresql DB.
func (l *Evento) GetPostgresqlEventos() (*schemas.EventosPaginados, *errors.Error) {
	eventoModel, err := l.DaoPostgresql.Evento.ObtenerEventosDisponiblesSinFiltros(eventoID)
	if err != nil {
		return nil, &errors.ObjectNotFoundError.LocalNotFound
	}
	eventos := make([]*schemas., len(localsModel))
	for i, localModel := range localsModel {
		locals[i] = &schemas.Local{
			Id:             localModel.Id,
			LocalName:      localModel.LocalName,
			StreetName:     localModel.StreetName,
			BuildingNumber: localModel.BuildingNumber,
			District:       localModel.District,
			Province:       localModel.Province,
			Region:         localModel.Region,
			Reference:      localModel.Reference,
			Capacity:       localModel.Capacity,
			ImageUrl:       localModel.ImageUrl,
		}
	}

	return &schemas.EventosPaginados{
		Eventos:        eventos,
		ImageUrl:       localModel.ImageUrl,
	}, nil
}

// Fetch all locals from postgresql DB.
func (l *Evento) FetchPostgresqlLocals() ([]*schemas.Local, *errors.Error) {
	localsModel, err := l.DaoPostgresql.Local.FetchLocals()
	if err != nil {
		return nil, &errors.ObjectNotFoundError.LocalNotFound
	}

	locals := make([]*schemas.Local, len(localsModel))
	for i, localModel := range localsModel {
		locals[i] = &schemas.Local{
			Id:             localModel.Id,
			LocalName:      localModel.LocalName,
			StreetName:     localModel.StreetName,
			BuildingNumber: localModel.BuildingNumber,
			District:       localModel.District,
			Province:       localModel.Province,
			Region:         localModel.Region,
			Reference:      localModel.Reference,
			Capacity:       localModel.Capacity,
			ImageUrl:       localModel.ImageUrl,
		}
	}
	return locals, nil
}

// Creates a local into postgresql DB and returns it.
func (l *Evento) CreatePostgresqlLocal(
	localName string,
	streetName string,
	buildingNumber string,
	district string,
	province string,
	region string,
	reference string,
	capacity int,
	imageUrl string,
	updatedBy string,
) (*schemas.Local, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	localModel := &model.Local{
		Id:             uuid.New(),
		LocalName:      localName,
		StreetName:     streetName,
		BuildingNumber: buildingNumber,
		District:       district,
		Province:       province,
		Region:         region,
		Reference:      reference,
		Capacity:       capacity,
		ImageUrl:       imageUrl,
		AuditFields: model.AuditFields{
			UpdatedBy: updatedBy,
		},
	}

	if err := l.DaoPostgresql.Local.CreateLocal(localModel); err != nil {
		return nil, &errors.BadRequestError.LocalNotCreated
	}

	return &schemas.Local{
		Id:             localModel.Id,
		LocalName:      localModel.LocalName,
		StreetName:     localModel.StreetName,
		BuildingNumber: localModel.BuildingNumber,
		District:       localModel.District,
		Province:       localModel.Province,
		Region:         localModel.Region,
		Reference:      localModel.Reference,
		Capacity:       localModel.Capacity,
		ImageUrl:       localModel.ImageUrl,
	}, nil
}

// Updates a local given fields in postgresql DB and returns it.
func (l *Evento) UpdatePostgresqlLocal(
	id uuid.UUID,
	localName *string,
	streetName *string,
	buildingNumber *string,
	district *string,
	province *string,
	region *string,
	reference *string,
	capacity *int,
	imageUrl *string,
	updatedBy string,
) (*schemas.Evento, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	localModel, err := l.DaoPostgresql.Local.UpdateLocal(
		id,
		localName,
		streetName,
		buildingNumber,
		district,
		province,
		region,
		reference,
		capacity,
		imageUrl,
		updatedBy,
	)
	if err != nil {
		return nil, &errors.BadRequestError.LocalNotUpdated
	}

	return &schemas.Evento{
		Id:             localModel.Id,
		LocalName:      localModel.LocalName,
		StreetName:     localModel.StreetName,
		BuildingNumber: localModel.BuildingNumber,
		District:       localModel.District,
		Province:       localModel.Province,
		Region:         localModel.Region,
		Reference:      localModel.Reference,
		Capacity:       localModel.Capacity,
		ImageUrl:       localModel.ImageUrl,
	}, nil
}

// Delets a plan from postgresql BD
func (l *Evento) DeletePostgresqlLocal(localId uuid.UUID) *errors.Error {
	if err := l.DaoPostgresql.Local.DeleteLocal(localId); err != nil {
		return &errors.BadRequestError.LocalNotSoftDeleted
	}

	return nil
}
