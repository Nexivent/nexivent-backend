package adapter

import (
	"github.com/Nexivent/nexivent-backend/errors"
	"github.com/Nexivent/nexivent-backend/internal/dao/model"
	daoPostgresql "github.com/Nexivent/nexivent-backend/internal/dao/repository"
	"github.com/Nexivent/nexivent-backend/internal/schemas"
	"github.com/Nexivent/nexivent-backend/logging"
)

type RolUsuario struct {
	logger        logging.Logger
	DaoPostgresql *daoPostgresql.NexiventPsqlEntidades
}

// Creates Evento adapter
func NewRolUsuarioAdapter(
	logger logging.Logger,
	daoPostgresql *daoPostgresql.NexiventPsqlEntidades,
) *RolUsuario {
	return &RolUsuario{
		logger:        logger,
		DaoPostgresql: daoPostgresql,
	}
}

// FetchPostgresqlEventos retrieves the upcoming events without filters
func (r *RolUsuario) GetUserPostgresqlRoles(idUsuario int64) (*schemas.RolUsuarioResponse, *errors.Error) {
	roles, err := r.DaoPostgresql.RolesUsuario.ListarRolesDeUsuario(idUsuario)
	if err != nil {
		r.logger.Errorf("Failed to fetch roles: %v", err)
		return nil, &errors.BadRequestError.EventoNotFound
	}
	rolesResponse := make([]schemas.RolResponse, len(roles))

	for i, c := range roles {
		rol, err := r.DaoPostgresql.Roles.ObtenerRolPorID(c.RolID)
		if err != nil {
			r.logger.Errorf("Failed to fetch role: %v", err)
			return nil, &errors.BadRequestError.EventoNotFound
		}

		rolesResponse[i] = schemas.RolResponse{
			ID:     c.RolID,
			Nombre: rol.Nombre,
		}
	}
	rolesUsuarioResponse := &schemas.RolUsuarioResponse{
		IDUsuario: idUsuario,
		Roles:     rolesResponse,
	}

	return rolesUsuarioResponse, nil
}

func (r *RolUsuario) AsignarPostgresqlRolUser(rolUser schemas.RolUsuarioRequest) (*schemas.RolUsuarioResponse, *errors.Error) {
	rolUsuario, err := r.DaoPostgresql.RolesUsuario.AsignarRolAUsuario(rolUser.IDUsuario, rolUser.IDRol, 1)
	if err != nil {
		r.logger.Errorf("Failed to asign rol: %v", err)
		return nil, &errors.BadRequestError.EventoNotFound
	}
	rolUserRes := &schemas.RolUsuarioResponse{
		IDUsuario: rolUsuario.UsuarioID,
	}

	return rolUserRes, nil
}

func (ru *RolUsuario) RevokePostgresqlRolUser(rolUser schemas.RolUsuarioRequest) (string, *errors.Error) {
	err := ru.DaoPostgresql.RolesUsuario.BorrarRolDeUsuario(rolUser.IDUsuario, rolUser.IDRol, 1)
	if err != nil {
		ru.logger.Errorf("Failed to revoke roluser: %v", err)
		return "", &errors.BadRequestError.EventoNotFound
	}

	mensaje := "Rol revocado correctamente"
	return mensaje, nil
}

func (ru *RolUsuario) GetUserSPostgresqlByRol(rolId *int64) ([]schemas.UsuarioRolResponse, *errors.Error) {
	var usuarios []model.Usuario
	var err error

	// Si rolId es nil, traer TODOS los usuarios
	if rolId == nil {
		usuarios, err = ru.DaoPostgresql.RolesUsuario.ObtenerTodosLosUsuariosActivos()
	} else {
		// Si rolId tiene valor, traer solo usuarios con ese rol
		usuarios, err = ru.DaoPostgresql.RolesUsuario.ObtenerUsuariosPorRol(*rolId)
	}

	if err != nil {
		ru.logger.Errorf("Failed to list users by Rol: %v", err)
		return nil, &errors.BadRequestError.EventoNotFound
	}

	usuariosByRol := make([]schemas.UsuarioRolResponse, len(usuarios))
	for i, c := range usuarios {
		usuariosByRol[i] = schemas.UsuarioRolResponse{
			IDUsuario: c.ID,
			Nombre:    c.Nombre,
			Correo:    c.Correo,
			Estado:    c.Estado,
		}
	}

	return usuariosByRol, nil
}