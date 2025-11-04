package nexiventpsql

import (
	"github.com/Loui27/nexivent-backend/logging"
	"gorm.io/gorm"
	"github.com/Loui27/nexivent-backend/internal/dao/model"
)

type Cupon struct {
	logger logging.Logger
	PostgresqlDB *gorm.DB
}

func NewCuponController(
	logger logging.Logger,
	postgresqlDB *gorm.DB,
) *Cupon {
	return &Cupon{
		logger:       logger,
		PostgresqlDB: postgresqlDB,
	}
}

func (c *Cupon) ObtenerCupones() ([]*model.Cupon, error) {
	var cupones []*model.Cupon
	respuesta := c.PostgresqlDB.Find(&cupones)

	if respuesta.Error != nil {
		return nil, respuesta.Error
	}

	return cupones, nil
}

func (c *Cupon) CrearCupon(Cupon *model.Cupon) error {
	respuesta := c.PostgresqlDB.Create(Cupon)
	if respuesta.Error != nil {
		return respuesta.Error
	}

	return nil
}

func (c *Cupon) ActualizarCupon(Cupon *model.Cupon) error {
	respuesta := c.PostgresqlDB.Save(Cupon)
	if respuesta.Error != nil {
		return respuesta.Error
	}

	return nil
}

func (c *Cupon) BorrarCupon(CuponId int64) error {
	respuesta := c.PostgresqlDB.Where("id = ?", CuponId).Delete(&model.Cupon{})
	if respuesta.Error != nil {
		return respuesta.Error
	}

	// Check if any rows were affected
	if respuesta.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}