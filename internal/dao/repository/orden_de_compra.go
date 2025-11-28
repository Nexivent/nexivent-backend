package repository

import (
	"database/sql"
	"time"

	"github.com/Nexivent/nexivent-backend/internal/dao/model"
	util "github.com/Nexivent/nexivent-backend/internal/dao/model/util"
	"github.com/Nexivent/nexivent-backend/logging"
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

// CrearOrdenTemporal inserta una nueva orden temporal (estado = TEMPORAL).
func (c *OrdenDeCompra) CrearOrdenTemporal(orden *model.OrdenDeCompra) error {
	if orden == nil {
		return gorm.ErrInvalidData
	}
	orden.EstadoDeOrden = util.OrdenTemporal.Codigo()

	if err := c.PostgresqlDB.Create(orden).Error; err != nil {
		c.logger.Errorf("CrearOrdenTemporal: %v", err)
		return err
	}
	return nil
}

// ObtenerOrdenBasica trae una orden completa por ID.
func (c *OrdenDeCompra) ObtenerOrdenBasica(orderID int64) (*model.OrdenDeCompra, error) {
	var o model.OrdenDeCompra
	if err := c.PostgresqlDB.First(&o, "orden_de_compra_id = ?", orderID).Error; err != nil {
		return nil, err
	}
	return &o, nil
}

// CerrarOrdenTemporal aplica un lock de escritura (FOR UPDATE) sobre una orden temporal.
func (c *OrdenDeCompra) CerrarOrdenTemporal(orderID int64) (*model.OrdenDeCompra, error) {
	var o model.OrdenDeCompra

	if err := c.PostgresqlDB.
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
	row := c.PostgresqlDB.
		Table("orden_de_compra").
		Select("estado_de_orden, fecha_hora_ini, fecha_hora_fin, total").
		Where("orden_de_compra_id = ?", orderID).
		Row()

	scanErr := row.Scan(&estInt, &ini, &fin, &total)
	if scanErr != nil {
		if scanErr == sql.ErrNoRows {
			return 0, time.Time{}, nil, 0, gorm.ErrRecordNotFound
		}
		c.logger.Errorf("ObtenerMetaTemporal(%d): %v", orderID, scanErr)
		return 0, time.Time{}, nil, 0, scanErr
	}
	estado = util.EstadoOrden(estInt) // casteo al enum
	return estado, ini, fin, total, nil
}

// VerificarOrdenExisteYEstado valida si la orden existe y está en un estado específico (enum).
func (c *OrdenDeCompra) VerificarOrdenExisteYEstado(orderID int64, estadoEsperado util.EstadoOrden) (bool, error) {
	var estInt int16
	row := c.PostgresqlDB.
		Table("orden_de_compra").
		Select("estado_de_orden").
		Where("orden_de_compra_id = ?", orderID).
		Row()

	if err := row.Scan(&estInt); err != nil {
		if err == sql.ErrNoRows {
			return false, gorm.ErrRecordNotFound
		}
		c.logger.Errorf("VerificarOrdenExisteYEstado(%d): %v", orderID, err)
		return false, err
	}
	return estInt == estadoEsperado.Codigo(), nil
}

// ActualizarEstadoOrden cambia el estado de la orden usando enum.
func (c *OrdenDeCompra) ActualizarEstadoOrden(orderID int64, nuevo util.EstadoOrden) error {
	res := c.PostgresqlDB.
		Table("orden_de_compra").
		Where("orden_de_compra_id = ?", orderID).
		Update("estado_de_orden", nuevo.Codigo())

	if res.Error != nil {
		c.logger.Errorf("ActualizarEstadoOrden(%d): %v", orderID, res.Error)
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (o *OrdenDeCompra) ObtenerIngresoCargoPorFecha(eventoID int64, fechaDesde *time.Time, fechaHasta *time.Time) (float64, float64, int64) {
	if fechaHasta == nil {
		fecha := time.Now()
		fechaHasta = &fecha
	}

	type IngresoCargoDTO struct {
		IngresoTotal    float64 `gorm:"column:ingreso_total"`
		CargoServ       float64 `gorm:"column:cargo_serv"`
		TicketsVendidos int64   `gorm:"column:tickets_vendidos"`
	}
	var data IngresoCargoDTO

	query := o.PostgresqlDB.Model(&model.OrdenDeCompra{}).
		Select(`
            oc.total,
            oc.monto_fee_servicio,
            COUNT(t.ticket_id) AS tickets_vendidos
        `).Table("orden_de_compra oc").
		Joins("JOIN ticket t ON t.orden_de_compra_id = oc.orden_de_compra_id").
		Joins("JOIN evento_fecha ef ON ef.evento_fecha_id = t.evento_fecha_id").
		Where("ef.evento_id = ?", eventoID).
		Where("oc.estado_de_orden = ?", util.OrdenConfirmada.Codigo()).
		Where("t.estado_de_ticket = ?", util.TicketVendido.Codigo())

	if fechaDesde != nil {
		query = query.Where("oc.fecha BETWEEN ? AND ?", fechaDesde, fechaHasta)
	} else {
		query = query.Where("oc.fecha <= ?", fechaHasta)
	}
	query = query.Group("oc.orden_de_compra_id, oc.total, oc.monto_fee_servicio")

	err := o.PostgresqlDB.Table("(?) as ordenes", query).
		Select("COALESCE(SUM(total), 0) AS ingreso_total, COALESCE(SUM(monto_fee_servicio), 0) AS cargo_serv, COALESCE(SUM(tickets_vendidos), 0) AS tickets_vendidos").
		Scan(&data).Error

	if err != nil {
		o.logger.Errorf("ObtenerIngresoCargoPorFecha evento_id=%d: %v", eventoID, err)
		return 0, 0, 0
	}

	return data.IngresoTotal, data.CargoServ, data.TicketsVendidos
}

// VentaPorSectorDTO resume ventas por sector para un evento.
type VentaPorSectorDTO struct {
	Sector          string  `gorm:"column:sector"`
	TicketsVendidos int64   `gorm:"column:tickets_vendidos"`
	Ingresos        float64 `gorm:"column:ingresos"`
}

func (o *OrdenDeCompra) ObtenerVentasPorSector(eventoID int64, fechaDesde *time.Time, fechaHasta *time.Time) ([]VentaPorSectorDTO, error) {
	if fechaHasta == nil {
		fecha := time.Now()
		fechaHasta = &fecha
	}

	var data []VentaPorSectorDTO

	query := o.PostgresqlDB.Table("ticket t").
		Select(`
			s.sector_tipo AS sector,
			COUNT(t.ticket_id) AS tickets_vendidos,
			COALESCE(SUM(tf.precio), 0) AS ingresos
		`).
		Joins("JOIN orden_de_compra oc ON oc.orden_de_compra_id = t.orden_de_compra_id").
		Joins("JOIN evento_fecha ef ON ef.evento_fecha_id = t.evento_fecha_id").
		Joins("JOIN tarifa tf ON tf.tarifa_id = t.tarifa_id").
		Joins("JOIN sector s ON s.sector_id = tf.sector_id").
		Where("ef.evento_id = ?", eventoID).
		Where("oc.estado_de_orden = ?", util.OrdenConfirmada.Codigo()).
		Where("t.estado_de_ticket <> ?", util.TicketCancelado.Codigo()).
		Group("s.sector_tipo")

	if fechaDesde != nil {
		query = query.Where("oc.fecha BETWEEN ? AND ?", fechaDesde, fechaHasta)
	} else {
		query = query.Where("oc.fecha <= ?", fechaHasta)
	}

	if err := query.Find(&data).Error; err != nil {
		o.logger.Errorf("ObtenerVentasPorTipo evento_id=%d: %v", eventoID, err)
		return nil, err
	}

	return data, nil
}

func (c *OrdenDeCompra) ConfirmarOrdenConPago(orderID int64, metodoPagoID int64, paymentReference string) error {
	updates := map[string]interface{}{
		"estado_de_orden":   util.OrdenConfirmada.Codigo(),
		"metodo_de_pago_id": metodoPagoID,
	}

	res := c.PostgresqlDB.
		Table("orden_de_compra").
		Where("orden_de_compra_id = ?", orderID).
		Updates(updates)

	if res.Error != nil {
		c.logger.Errorf("ConfirmarOrdenConPago(%d): %v", orderID, res.Error)
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// ObtenerTransaccionesPorEvento obtiene todas las órdenes de compra (transacciones) asociadas a un evento específico.
func (c *OrdenDeCompra) ListarTransaccionesPorEvento(eventoId string) ([]model.OrdenDeCompra, error) {
	var ordenes []model.OrdenDeCompra
	err := c.PostgresqlDB.
		Joins("JOIN ticket t ON t.orden_de_compra_id = orden_de_compra.orden_de_compra_id").
		Joins("JOIN evento_fecha ef ON ef.evento_fecha_id = t.evento_fecha_id").
		Where("ef.evento_id = ?", eventoId).
		Find(&ordenes).Error
	if err != nil {
		c.logger.Errorf("ListarTransaccionesPorEvento(%s): %v", eventoId, err)
		return nil, err
	}

	return ordenes, nil
}
