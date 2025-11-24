package repository

import (
	"errors"
	"time"
	//"github.com/Nexivent/nexivent-backend/internal/schemas"
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

func NewTicketController(
	logger logging.Logger,
	postgresqlDB *gorm.DB,
) *Ticket {
	return &Ticket{
		logger:       logger,
		PostgresqlDB: postgresqlDB,
	}
}

// / VerificarOrdenConfirmada: true si la orden existe y está CONFIRMADA (enum).
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

// CrearTicket: inserta un solo ticket.
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

// CrearTicketsBatch: inserta varios tickets.
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

// CambiarEstadoTicket: actualiza el estado de un ticket usando el enum util.EstadoDeTicket.
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

// ObtenerTicketsPorOrden: devuelve tickets (modelo puro) de una orden.
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

// Helpers extra para BO (solo structs internos, NO schemas
// DetalleOrden representa una fila de la tabla orden_de_compra_detalle.
type DetalleOrden struct {
	TarifaID      int64 `gorm:"column:tarifa_id"`
	Cantidad      int64 `gorm:"column:cantidad"`
	EventoFechaID int64 `gorm:"column:evento_fecha_id"`
}

func (c *Ticket) ObtenerDetallesOrden(orderID int64) ([]DetalleOrden, error) {
	var detalles []DetalleOrden
	res := c.PostgresqlDB.
		Table("orden_de_compra_detalle").
		Select("tarifa_id, cantidad, evento_fecha_id").
		Where("orden_de_compra_id = ?", orderID).
		Find(&detalles)

	if res.Error != nil {
		c.logger.Errorf("ObtenerDetallesOrden(orderID=%d): %v", orderID, res.Error)
		return nil, res.Error
	}
	return detalles, nil
}

// TicketInfo: datos enriquecidos para mostrar tickets (JOIN con evento, sector, fecha, etc.).
type TicketInfo struct {
	ID          int64     `gorm:"column:ticket_id"`
	CodigoQR    string    `gorm:"column:codigo_qr"`
	Estado      int16     `gorm:"column:estado_de_ticket"`
	Titulo      string    `gorm:"column:titulo"`
	FechaEvento time.Time `gorm:"column:fecha_evento"`
	HoraInicio  time.Time `gorm:"column:hora_inicio"`
	SectorTipo  string    `gorm:"column:sector_tipo"`
}

func (c *Ticket) ObtenerTicketsInfoPorOrden(orderID int64) ([]TicketInfo, error) {
	var info []TicketInfo

	res := c.PostgresqlDB.
		Table("ticket t").
		Select(`
			t.ticket_id,
			t.codigo_qr,
			t.estado_de_ticket,
			e.titulo,
			f.fecha_evento,
			ef.hora_inicio,
			s.sector_tipo
		`).
		Joins("JOIN evento_fecha ef ON ef.evento_fecha_id = t.evento_fecha_id").
		Joins("JOIN evento e ON e.evento_id = ef.evento_id").
		Joins("JOIN tarifa tf ON tf.tarifa_id = t.tarifa_id").
		Joins("JOIN sector s ON s.sector_id = tf.sector_id").
		Joins("JOIN fecha f ON f.fecha_id = ef.fecha_id").
		Where("t.orden_de_compra_id = ?", orderID).
		Find(&info)

	if res.Error != nil {
		c.logger.Errorf("ObtenerTicketsInfoPorOrden(orderID=%d): %v", orderID, res.Error)
		return nil, res.Error
	}
	return info, nil
}

// TicketEstadoTarifa: helper para cancelación (estado actual + tarifa).
type TicketEstadoTarifa struct {
	ID             int64 `gorm:"column:ticket_id"`
	TarifaID       int64 `gorm:"column:tarifa_id"`
	EstadoDeTicket int16 `gorm:"column:estado_de_ticket"`
}

func (c *Ticket) ObtenerTicketsEstadoTarifaPorIDs(ids []int64) ([]TicketEstadoTarifa, error) {
	if len(ids) == 0 {
		return []TicketEstadoTarifa{}, nil
	}

	var rows []TicketEstadoTarifa
	res := c.PostgresqlDB.
		Table("ticket").
		Select("ticket_id, tarifa_id, estado_de_ticket").
		Where("ticket_id IN ?", ids).
		Find(&rows)

	if res.Error != nil {
		c.logger.Errorf("ObtenerTicketsEstadoTarifaPorIDs: %v", res.Error)
		return nil, res.Error
	}
	return rows, nil
}

// Crear un nuevo ticket
func (c *Ticket) Crear(ticket *model.Ticket) error {
	if ticket == nil {
		return gorm.ErrInvalidData
	}
	
	if err := c.PostgresqlDB.Create(ticket).Error; err != nil {
		c.logger.Errorf("Ticket.Crear: %v", err)
		return err
	}
	return nil
}

type TicketRaro struct {
	IDTicket    int64     //`json:"ticket_id"`
	TipoSector  string    //`json:"sector_tipo"`
	//Evento      EventoMini    `json:"evento"`
	IDEvento      int64  //`json:"evento_id"`
	Titulo        string //`json:"titulo"`
	Lugar         string //`json:"lugar"`
	ImagenPortada string //`json:"imagenPortada"`
	FechaInicio string //`json:"hora_inicio"`
}


func (t *Ticket) ObternerTicketsPorUsuario (idUser int64) ([]*TicketRaro, error) {

    var tickets []*TicketRaro
	res := t.PostgresqlDB.
		Table("ticket t").
		Select(
			"t.ticket_id as id_ticket",
			"s.sector_tipo as tipo_sector",
			"e.evento_id as id_evento",
			"e.titulo",
			"e.lugar",
			"e.imagen_portada",
			"ef.hora_inicio as fecha_inicio",
		).
		Joins("INNER JOIN orden_de_compra oc ON t.orden_de_compra_id = oc.orden_de_compra_id").
		Joins("INNER JOIN tarifa tar ON t.tarifa_id = tar.tarifa_id").
		Joins("INNER JOIN sector s ON tar.sector_id = s.sector_id").
		Joins("INNER JOIN evento e ON s.evento_id = e.evento_id").
		Joins("INNER JOIN evento_fecha ef ON t.evento_fecha_id = ef.evento_fecha_id").
		Joins("LEFT JOIN perfil_de_persona pp ON tar.perfil_de_persona_id = pp.perfil_de_persona_id").
		Where("oc.usuario_id = ?", idUser).
		Where("e.estado = 1").
		Where("ef.estado = 1").
		Find(&tickets)

	if res.Error != nil {
		t.logger.Errorf("ObtenerTickets de usuario: %v", res.Error)
		return nil, res.Error
	}
	return tickets, nil
}