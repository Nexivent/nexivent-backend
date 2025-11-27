package controller

import (
	"github.com/Nexivent/nexivent-backend/errors"
	"github.com/Nexivent/nexivent-backend/internal/application/adapter"
	"github.com/Nexivent/nexivent-backend/internal/schemas"
	"github.com/Nexivent/nexivent-backend/logging"
)

type RolController struct {
	Logger       logging.Logger
	RolAdapter *adapter.Rol
}

func NewRolController(
	logger logging.Logger,
	rolAdapter *adapter.Rol,
) *RolController {
	return &RolController{
		Logger:       logger,
		RolAdapter: rolAdapter,
	}
}

func (ec *RolController) FetchRoles() ([]schemas.RolResponse, *errors.Error) {
	return ec.RolAdapter.FetchPostgresqlRoles()
}

func (ec *RolController) GetRolPorNombre(nombre string) (*schemas.RolResponse, *errors.Error) {
	return ec.RolAdapter.GetPostgresqlRolPorNombre(nombre)
}

func (ec *RolController) GetRolPorUsuario(usuarioID int64) ([]*schemas.RolResponse, *errors.Error) {
	ec.Logger.Infof("🎯 [CONTROLLER] GetRolPorUsuario usuarioID=%d", usuarioID)
    roles, err := ec.RolAdapter.GetPostgresqlRolPorUsuario(usuarioID)
    if err != nil {
        ec.Logger.Errorf("❌ [CONTROLLER] error adapter: %v", err)
        return nil, err
    }
    ec.Logger.Infof("✅ [CONTROLLER] roles obtenidos: %d", len(roles))
    return roles, nil
}

func (r *RolController) ActualizarRol(request *schemas.RolRequest,rolID int64) (*schemas.RolResponse, *errors.Error) {
		
	return r.RolAdapter.ActualizarPostgresqlRol(rolID ,
	request.Nombre ,
*request.UsuarioModificacion)}