package repository

import (
	"database/sql"
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/model"
	"github.com/Nexivent/nexivent-backend/internal/data/model/util"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrdenDeCompra struct {
	DB *gorm.DB
}

// CrearOrdenTemporal inserta una nueva orden temporal (estado = TEMPORAL).
func (c *OrdenDeCompra) CrearOrdenTemporal(orden *model.OrdenDeCompra) error {
	if orden == nil {
		return gorm.ErrInvalidData
	}
	orden.EstadoDeOrden = util.OrdenTemporal

	if err := c.DB.Create(orden).Error; err != nil {
		return err
	}
	return nil
}

// ObtenerOrdenBasica trae una orden completa por ID.
func (c *OrdenDeCompra) ObtenerOrdenBasica(orderID int64) (*model.OrdenDeCompra, error) {
	var o model.OrdenDeCompra
	if err := c.DB.First(&o, "orden_de_compra_id = ?", orderID).Error; err != nil {
		return nil, err
	}
	return &o, nil
}

// CerrarOrdenTemporal aplica un lock de escritura (FOR UPDATE) sobre una orden temporal.
func (c *OrdenDeCompra) CerrarOrdenTemporal(orderID int64) (*model.OrdenDeCompra, error) {
	var o model.OrdenDeCompra
	if err := c.DB.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("orden_de_compra_id = ? AND estado_de_orden = ?", orderID, util.OrdenTemporal.Codigo()).
		First(&o).Error; err != nil {
		return nil, err
	}
	return &o, nil
}

// ObtenerMetaTemporal devuelve (estado como enum), ini, fin, total.
func (c *OrdenDeCompra) ObtenerMetaTemporal(orderID int64) (estado util.EstadoOrden, ini time.Time, fin *time.Time, total float64, err error) {
	var estInt int16
	row := c.DB.
		Table("orden_de_compra").
		Select("estado_de_orden, fecha_hora_ini, fecha_hora_fin, total").
		Where("orden_de_compra_id = ?", orderID).
		Row()

	scanErr := row.Scan(&estInt, &ini, &fin, &total)
	if scanErr != nil {
		if scanErr == sql.ErrNoRows {
			return 0, time.Time{}, nil, 0, gorm.ErrRecordNotFound
		}
		return 0, time.Time{}, nil, 0, scanErr
	}
	estado = util.EstadoOrden(estInt) // casteo al enum
	return estado, ini, fin, total, nil
}

// VerificarOrdenExisteYEstado valida si la orden existe y está en un estado específico (enum).
func (c *OrdenDeCompra) VerificarOrdenExisteYEstado(orderID int64, estadoEsperado util.EstadoOrden) (bool, error) {
	var estInt int16
	row := c.DB.
		Table("orden_de_compra").
		Select("estado_de_orden").
		Where("orden_de_compra_id = ?", orderID).
		Row()

	if err := row.Scan(&estInt); err != nil {
		if err == sql.ErrNoRows {
			return false, gorm.ErrRecordNotFound
		}
		return false, err
	}
	return util.EstadoOrden(estInt) == estadoEsperado, nil
}

// ActualizarEstadoOrden cambia el estado de la orden usando enum.
func (c *OrdenDeCompra) ActualizarEstadoOrden(orderID int64, nuevo util.EstadoOrden) error {
	res := c.DB.
		Table("orden_de_compra").
		Where("orden_de_compra_id = ?", orderID).
		Update("estado_de_orden", nuevo.Codigo())

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
