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

type Rol struct {
	logger        logging.Logger
	DaoPostgresql *daoPostgresql.NexiventPsqlEntidades
}

// Creates Evento adapter
func NewRolAdapter(
	logger logging.Logger,
	daoPostgresql *daoPostgresql.NexiventPsqlEntidades,
) *Rol {
	return &Rol{
		logger:        logger,
		DaoPostgresql: daoPostgresql,
	}
}


// FetchPostgresqlEventos retrieves the upcoming events without filters
func (r *Rol) FetchPostgresqlRoles() ([] schemas.RolResponse, *errors.Error) {
	roles, err := r.DaoPostgresql.Roles.ObtenerRoles()
	if err != nil {
		r.logger.Errorf("Failed to fetch roles: %v", err)
		return nil, &errors.BadRequestError.EventoNotFound
	}
	rolesResponse := make([]schemas.RolResponse, len(roles))

	for i, c := range roles{
		rolesResponse[i] = schemas.RolResponse{
			ID: c.ID,
			Nombre: c.Nombre,
		}
	}

	return rolesResponse, nil
}

// GetPostgresqlEventoById gets an event by ID

func (e *Rol) GetPostgresqlRolPorNombre(nombre string) (*schemas.RolResponse, *errors.Error) {
	var rolModel *model.Rol

	rolModel, err := e.DaoPostgresql.Roles.ObtenerRolPorNombre(nombre)
	if err != nil {
		e.logger.Errorf("Failed to fetch roles: %v", err)
		return nil, &errors.BadRequestError.EventoNotFound
	}

	response := &schemas.RolResponse{
		ID:          rolModel.ID,
		Nombre:     rolModel.Nombre,
	}

	return response, nil
}

func (r *Rol) GetPostgresqlRolPorUsuario(usuarioID int64) ([]*schemas.RolResponse, *errors.Error) {
	//var roles []*model.Rol
 	r.logger.Infof("🔍 [ADAPTER] Obteniendo roles para usuario ID: %d", usuarioID)
	roles, err := r.DaoPostgresql.Roles.ObtenerRolesDeUsuario(usuarioID)
	if err != nil {
        r.logger.Errorf("❌ [ADAPTER] Error obteniendo roles: %v", err)
        return nil, &errors.BadRequestError.EventoNotFound
    }
	r.logger.Infof("📦 [ADAPTER] Roles obtenidos del repo: %d", len(roles))
	response := make([]*schemas.RolResponse, 0, len(roles))

	for _, mr := range roles {
        // mapear model.Rol -> schemas.RolResponse
        response = append(response, &schemas.RolResponse{
            ID:     mr.ID,
            Nombre: mr.Nombre,
        })
        r.logger.Infof("   [ADAPTER] rol encontrado ID=%d Nombre=%s", mr.ID, mr.Nombre)
    }

    r.logger.Infof("✅ [ADAPTER] roles procesados: %d", len(response))
	return response, nil
}

//Actualizar Rol
func (r *Rol) ActualizarPostgresqlRol(id int64,
	nombre string,
	updatedBy int64 ) ( *schemas.RolResponse, *errors.Error) {
	rolModel, err := r.DaoPostgresql.Roles.ActualizarRol(id, &nombre, updatedBy,)

	if err != nil {
		r.logger.Errorf("Failed to update rol: %v", err)
		return nil, &errors.BadRequestError.EventoNotFound
	}

	rolResponse := &schemas.RolResponse {
		ID: rolModel.ID,
		Nombre: rolModel.Nombre,
		//UsuarioModificacion: rolModel.UsuarioModificacion,
	}

	return rolResponse, nil
}
