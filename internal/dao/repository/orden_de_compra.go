package repository

import (
	"database/sql"
	"time"

	"github.com/Loui27/nexivent-backend/internal/dao/model"
	util "github.com/Loui27/nexivent-backend/internal/dao/model/util"
	"github.com/Loui27/nexivent-backend/logging"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrdenDeCompra struct {
	logger       logging.Logger
	PostgresqlDB *gorm.DB
}

func NewOrdenDeCompraController(
	logger logging.Logger,
	postgresqlDB *gorm.DB,
) *OrdenDeCompra {
	return &OrdenDeCompra{
		logger:       logger,
		PostgresqlDB: postgresqlDB,
	}
}

// CrearOrdenTemporal inserta una orden en estado_de_orden=0 (TEMPORAL).
// Requiere: UsuarioID, Fecha, FechaHoraIni, FechaHoraFin (TTL calculado por el BO), Total.
func (c *OrdenDeCompra) CrearOrdenTemporal(orden *model.OrdenDeCompra) error {
	if orden == nil {
		return gorm.ErrInvalidData
	}
	if orden.FechaHoraFin == nil {
		c.logger.Errorf("CrearOrdenTemporal: FechaHoraFin es nil (el BO debe setear el TTL)")
		return gorm.ErrInvalidData
	}
	// 0 = TEMPORAL (según tu esquema)
	orden.EstadoDeOrden = util.OrdenTemporal.Codigo()

	res := c.PostgresqlDB.Create(orden)
	if res.Error != nil {
		c.logger.Errorf("CrearOrdenTemporal: %v", res.Error)
		return res.Error
	}
	return nil
}

// ObtenerMetaTemporal: trae estado_de_orden, fecha_hora_ini, fecha_hora_fin y total.
func (c *OrdenDeCompra) ObtenerMetaTemporal(orderID int64) (estado int16, ini time.Time, fin *time.Time, total float64, err error) {
	row := c.PostgresqlDB.
		Table("orden_de_compra").
		Where("orden_de_compra_id = ?", orderID).
		Row()

	scanErr := row.Scan(&estado, &ini, &fin, &total)
	if scanErr != nil {
		if scanErr == sql.ErrNoRows {
			return 0, time.Time{}, nil, 0, gorm.ErrRecordNotFound
		}
		c.logger.Errorf("ObtenerMetaTemporal(%d): %v", orderID, scanErr)
		return 0, time.Time{}, nil, 0, scanErr
	}
	return estado, ini, fin, total, nil
}

// MarcarOrdenConfirmada: cambia a estado_de_orden=1 (CONFIRMADA) si sigue TEMPORAL (0) y no está expirada.
func (c *OrdenDeCompra) MarcarOrdenConfirmada(orderID int64) error {
	now := time.Now().UTC()
	res := c.PostgresqlDB.
		Table("orden_de_compra").
		Where(`
			orden_de_compra_id = ?
			AND estado_de_orden = 0
			AND fecha_hora_fin IS NOT NULL
			AND fecha_hora_fin >= ?
		`, orderID, now).
		Update("estado_de_orden", util.OrdenConfirmada.Codigo()) // 1 = CONFIRMADA

	if res.Error != nil {
		c.logger.Errorf("MarcarOrdenConfirmada(%d): %v", orderID, res.Error)
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// CancelarTemporalesVencidas: pone estado_de_orden=2 (CANCELADA) a todas las TEMPORALES (0) vencidas.
func (c *OrdenDeCompra) CancelarTemporalesVencidas() (int64, error) {
	now := time.Now().UTC()
	res := c.PostgresqlDB.
		Table("orden_de_compra").
		Where("estado_de_orden = 0 AND fecha_hora_fin IS NOT NULL AND fecha_hora_fin < ?", now).
		Update("estado_de_orden", util.OrdenCancelada.Codigo()) // 2 = CANCELADA
	if res.Error != nil {
		c.logger.Errorf("CancelarTemporalesVencidas: %v", res.Error)
		return 0, res.Error
	}
	return res.RowsAffected, nil
}

// Ayuda rápida para lecturas básicas por ID (opcional).
func (c *OrdenDeCompra) ObtenerOrdenBasica(orderID int64) (*model.OrdenDeCompra, error) {
	var o model.OrdenDeCompra
	if err := c.PostgresqlDB.First(&o, "orden_de_compra_id = ?", orderID).Error; err != nil {
		return nil, err
	}
	return &o, nil
}

// Lock de una orden temporal (FOR UPDATE), útil si el BO necesita coordinar otras acciones.
func (c *OrdenDeCompra) CerrarOrdenTemporal(orderID int64) (*model.OrdenDeCompra, error) {
	var o model.OrdenDeCompra
	if err := c.PostgresqlDB.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("orden_de_compra_id = ? AND estado_de_orden = ?", orderID, util.OrdenTemporal.Codigo()). // 0 = TEMPORAL
		First(&o).Error; err != nil {
		return nil, err
	}
	return &o, nil
}
