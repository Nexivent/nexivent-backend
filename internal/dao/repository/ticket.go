package repository

import (
	"errors"

	"github.com/Nexivent/nexivent-backend/internal/dao/model"
	util "github.com/Nexivent/nexivent-backend/internal/dao/model/util"
	"github.com/Nexivent/nexivent-backend/logging"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Ticket struct {
	logger       logging.Logger
	PostgresqlDB *gorm.DB
}

func NewTicketsController(
	logger logging.Logger,
	postgresqlDB *gorm.DB,
) *Ticket {
	return &Ticket{
		logger:       logger,
		PostgresqlDB: postgresqlDB,
	}
}

// VerificarOrdenConfirmada: true si la orden existe y está CONFIRMADA (enum).
func (c *Ticket) VerificarOrdenConfirmada(orderID int64) (bool, error) {
	var estInt int16
	err := c.PostgresqlDB.
		Table("orden_de_compra").
		Select("estado_de_orden").
		Where("orden_de_compra_id = ?", orderID).
		Row().
		Scan(&estInt)

	if err != nil {
		c.logger.Errorf("VerificarOrdenConfirmada(%d): %v", orderID, err)
		return false, err
	}
	return util.EstadoOrden(estInt) == util.OrdenConfirmada, nil
}

// VerificarTicketsExistentes: true si la orden ya tiene tickets creados.
func (c *Ticket) VerificarTicketsExistentes(orderID int64) (bool, error) {
	var count int64
	res := c.PostgresqlDB.
		Table("ticket").
		Where("orden_de_compra_id = ?", orderID).
		Count(&count)
	if res.Error != nil {
		c.logger.Errorf("VerificarTicketsExistentes(%d): %v", orderID, res.Error)
		return false, res.Error
	}
	return count > 0, nil
}

// VerificarStockDisponible: revisa capacidad del sector asociado a la tarifa.
func (c *Ticket) VerificarStockDisponible(tarifaID int64, cantidad int64) (bool, error) {
	if tarifaID <= 0 || cantidad <= 0 {
		return false, gorm.ErrInvalidData
	}

	var total, vendidas int64
	row := c.PostgresqlDB.
		Table("sector s").
		Select("s.total_entradas, s.cant_vendidas").
		Joins("JOIN tarifa t ON t.sector_id = s.sector_id").
		Where("t.tarifa_id = ?", tarifaID).
		Row()

	if err := row.Scan(&total, &vendidas); err != nil {
		c.logger.Errorf("VerificarStockDisponible(tarifa=%d): %v", tarifaID, err)
		return false, err
	}

	return vendidas+cantidad <= total, nil
}

// VerificarTicketPerteneceAOrden: true si ticket_id pertenece a la orden dada.
func (c *Ticket) VerificarTicketPerteneceAOrden(ticketID, orderID int64) (bool, error) {
	var count int64
	res := c.PostgresqlDB.
		Table("ticket").
		Where("ticket_id = ? AND orden_de_compra_id = ?", ticketID, orderID).
		Count(&count)
	if res.Error != nil {
		c.logger.Errorf("VerificarTicketPerteneceAOrden(%d,%d): %v", ticketID, orderID, res.Error)
		return false, res.Error
	}
	return count == 1, nil
}

// CrearTicket: inserta un solo ticket (el BO define QR y estado).
func (c *Ticket) CrearTicket(ticket *model.Ticket) error {
	if ticket == nil {
		return gorm.ErrInvalidData
	}
	if err := c.PostgresqlDB.Create(ticket).Error; err != nil {
		c.logger.Errorf("CrearTicket: %v", err)
		return err
	}
	return nil
}

// CrearTicketsBatch: inserta varios tickets (BO validó todo antes).
func (c *Ticket) CrearTicketsBatch(tickets []model.Ticket) error {
	if len(tickets) == 0 {
		return nil
	}
	if err := c.PostgresqlDB.Create(&tickets).Error; err != nil {
		c.logger.Errorf("CrearTicketsBatch: %v", err)
		return err
	}
	return nil
}

// IncrementarVendidasPorSector: suma cantidad a cant_vendidas (sector).
func (c *Ticket) IncrementarVendidasPorSector(sectorID int64, cantidad int64) error {
	if sectorID <= 0 || cantidad <= 0 {
		return gorm.ErrInvalidData
	}
	res := c.PostgresqlDB.
		Table("sector").
		Where("sector_id = ?", sectorID).
		UpdateColumn("cant_vendidas", gorm.Expr("cant_vendidas + ?", cantidad))
	if res.Error != nil {
		c.logger.Errorf("IncrementarVendidasPorSector(%d): %v", sectorID, res.Error)
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// DecrementarVendidasPorSector: resta cantidad a cant_vendidas (sector).
func (c *Ticket) DecrementarVendidasPorSector(sectorID int64, cantidad int64) error {
	if sectorID <= 0 || cantidad <= 0 {
		return gorm.ErrInvalidData
	}
	res := c.PostgresqlDB.
		Table("sector").
		Where("sector_id = ?", sectorID).
		UpdateColumn("cant_vendidas", gorm.Expr("cant_vendidas - ?", cantidad))
	if res.Error != nil {
		c.logger.Errorf("DecrementarVendidasPorSector(%d): %v", sectorID, res.Error)
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// ObtenerSectorPorTarifa: retorna sector_id para una tarifa.
func (c *Ticket) ObtenerSectorPorTarifa(tarifaID int64) (int64, error) {
	var sectorID int64
	row := c.PostgresqlDB.
		Table("tarifa").
		Select("sector_id").
		Where("tarifa_id = ?", tarifaID).
		Row()
	if err := row.Scan(&sectorID); err != nil {
		c.logger.Errorf("ObtenerSectorPorTarifa(%d): %v", tarifaID, err)
		return 0, err
	}
	return sectorID, nil
}

// CambiarEstadoTicket: actualiza el estado de un ticket usando tu enum model.EstadoDeTicket.
func (c *Ticket) CambiarEstadoTicket(ticketID int64, nuevo util.EstadoDeTicket) error {
	if !nuevo.IsValid() {
		return errors.New("estado de ticket inválido")
	}
	res := c.PostgresqlDB.
		Table("ticket").
		Where("ticket_id = ?", ticketID).
		Update("estado_de_ticket", nuevo.Codigo())
	if res.Error != nil {
		c.logger.Errorf("CambiarEstadoTicket(%d): %v", ticketID, res.Error)
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// CambiarEstadoTicketsDeOrden: actualiza el estado de todos los tickets de una orden.
func (c *Ticket) CambiarEstadoTicketsDeOrden(orderID int64, nuevo util.EstadoDeTicket) (int64, error) {
	if !nuevo.IsValid() {
		return 0, errors.New("estado de ticket inválido")
	}
	res := c.PostgresqlDB.
		Table("ticket").
		Where("orden_de_compra_id = ?", orderID).
		Update("estado_de_ticket", nuevo.Codigo())
	if res.Error != nil {
		c.logger.Errorf("CambiarEstadoTicketsDeOrden(%d): %v", orderID, res.Error)
		return 0, res.Error
	}
	return res.RowsAffected, nil
}

// ObtenerTicketsPorOrden: devuelve tickets (modelo puro) de una orden (para el BO).
func (c *Ticket) ObtenerTicketsPorOrden(orderID int64) ([]model.Ticket, error) {
	var ts []model.Ticket
	if err := c.PostgresqlDB.
		Where("orden_de_compra_id = ?", orderID).
		Find(&ts).Error; err != nil {
		c.logger.Errorf("ObtenerTicketsPorOrden(%d): %v", orderID, err)
		return nil, err
	}
	return ts, nil
}

// LockTicketsDeOrden: bloquea filas de tickets por orden (si el BO necesita transacción estricta).
func (c *Ticket) LockTicketsDeOrden(orderID int64) ([]model.Ticket, error) {
	var ts []model.Ticket
	if err := c.PostgresqlDB.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("orden_de_compra_id = ?", orderID).
		Find(&ts).Error; err != nil {
		c.logger.Errorf("LockTicketsDeOrden(%d): %v", orderID, err)
		return nil, err
	}
	return ts, nil
}
