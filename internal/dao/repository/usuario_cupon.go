package repository

import (
	"github.com/Nexivent/nexivent-backend/internal/dao/model"
	"github.com/Nexivent/nexivent-backend/logging"
	"gorm.io/gorm"
)

type UsuarioCupon struct {
	logger       logging.Logger
	PostgresqlDB *gorm.DB
}

func NewUsuarioCuponController(
	logger logging.Logger,
	postgresqlDB *gorm.DB,
) *UsuarioCupon {
	return &UsuarioCupon{
		logger:       logger,
		PostgresqlDB: postgresqlDB,
	}
}

func (uc *UsuarioCupon) ObtenerUsuarioCuponPorId(usuarioId int64, cuponId int64) (*model.UsuarioCupon, error) {
	var usuarioCupon model.UsuarioCupon  // ← Cambio 1: valor, no puntero

	resp := uc.PostgresqlDB.
		Where("usuario_id = ? AND cupon_id = ?", usuarioId, cuponId).
		First(&usuarioCupon)  // ← Cambio 2: First, no Find

	if resp.Error != nil {  // ← Cambio 3: resp.Error, no resp
		return nil, resp.Error
	}

	return &usuarioCupon, nil
}

func (uc *UsuarioCupon) CrearUsuarioCupon(usuarioCupon *model.UsuarioCupon) error {
	respuesta := uc.PostgresqlDB.Create(usuarioCupon)
	if respuesta.Error != nil {
		return respuesta.Error
	}
	return nil
}

func (uc *UsuarioCupon) ActualizarUsuarioCupon(usuarioCupon *model.UsuarioCupon) error {
	respuesta := uc.PostgresqlDB.Save(usuarioCupon)
	if respuesta.Error != nil {
		return respuesta.Error
	}

	return nil
}