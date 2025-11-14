package repository

import (
	"github.com/Nexivent/nexivent-backend/internal/data/model"
	"gorm.io/gorm"
)

type Cupon struct {
	DB *gorm.DB
}

func NewCuponController(
	postgresqlDB *gorm.DB,
) *Cupon {
	return &Cupon{
		DB: postgresqlDB,
	}
}

func (c *Cupon) ObtenerCupones() ([]*model.Cupon, error) {
	var cupones []*model.Cupon
	respuesta := c.DB.Find(&cupones)

	if respuesta.Error != nil {
		return nil, respuesta.Error
	}

	return cupones, nil
}

func (c *Cupon) CrearCupon(Cupon *model.Cupon) error {
	respuesta := c.DB.Create(Cupon)
	if respuesta.Error != nil {
		return respuesta.Error
	}

	return nil
}

func (c *Cupon) ActualizarCupon(Cupon *model.Cupon) error {
	respuesta := c.DB.Save(Cupon)
	if respuesta.Error != nil {
		return respuesta.Error
	}

	return nil
}

/*
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
}*/
