package repository

import (
	"github.com/Nexivent/nexivent-backend/internal/dao/model"
	"github.com/Nexivent/nexivent-backend/logging"
	"gorm.io/gorm"
)

type Cupon struct {
	logger       logging.Logger
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

func (c *Cupon) ObtenerCuponPorIdYIdEvento(id int64, eventoId int64) (*model.Cupon, error) {
	var cupon model.Cupon
	respuesta := c.PostgresqlDB.
		Where("cupon_id = ? AND evento_id = ?", id, eventoId). 
		First(&cupon)

	if respuesta.Error != nil {
		return nil, respuesta.Error
	}

	// verificar que se encontró el registro
	if respuesta.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &cupon, nil
}

func (c *Cupon) ObtenerCuponesPorOrganizador(organizadorId int64) ([]*model.Cupon, error) {
	var cupones []*model.Cupon
	respuesta := c.PostgresqlDB.
		Table("cupon").
		Joins("JOIN evento e ON e.evento_id = cupon.evento_id").
		Where("e.organizador_id  = ?", organizadorId).
		Find(&cupones)

	if respuesta.Error != nil {
		return nil, respuesta.Error
	}

	return cupones, nil
}

func (c *Cupon) ObtenerCuponPorCodYIdEvento(eventoId int64, codigo string) (*model.Cupon, error) {
	var cupon model.Cupon
	respuesta := c.PostgresqlDB.
		Where("codigo = ? AND evento_id = ?", codigo, eventoId). 
		First(&cupon)

	if respuesta.Error != nil {
		return nil, respuesta.Error
	}

	// verificar que se encontró el registro
	if respuesta.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &cupon, nil
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