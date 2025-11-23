package controller

import (
	"github.com/Nexivent/nexivent-backend/errors"
	"github.com/Nexivent/nexivent-backend/internal/application/adapter"
	"github.com/Nexivent/nexivent-backend/internal/schemas"
	"github.com/Nexivent/nexivent-backend/logging"
)

type RolUsuarioController struct {
	Logger       logging.Logger
	RolUsuarioAdapter *adapter.RolUsuario
}

func NewRolUsuarioController(
	logger logging.Logger,
	rolUsuarioAdapter *adapter.RolUsuario,
) *RolUsuarioController {
	return &RolUsuarioController{
		Logger:       logger,
		RolUsuarioAdapter: rolUsuarioAdapter,
	}
}

func (ru *RolUsuarioController) GetUserRoles(idUsuario int64) (*schemas.RolUsuarioResponse, *errors.Error) {
	return ru.RolUsuarioAdapter.GetUserPostgresqlRoles(idUsuario)
}

func (ru *RolUsuarioController) AsignarRolUser(request schemas.RolUsuarioRequest) (*schemas.RolUsuarioResponse, *errors.Error) {
	return ru.RolUsuarioAdapter.AsignarPostgresqlRolUser(request)
}

func (ru *RolUsuarioController) RevokeRolUser(request schemas.RolUsuarioRequest) (string, *errors.Error) {
	return ru.RolUsuarioAdapter.RevokePostgresqlRolUser(request)
}

func (ru *RolUsuarioController) GetUsersByRol(idRol int64) ([]schemas.UsuarioRolResponse, *errors.Error) {
	return ru.RolUsuarioAdapter.GetUserSPostgresqlByRol(idRol)
}