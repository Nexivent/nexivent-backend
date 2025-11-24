package adapter

import (
	"fmt"
	"time"

	"github.com/Nexivent/nexivent-backend/errors"
	"github.com/Nexivent/nexivent-backend/internal/dao/model"
	util "github.com/Nexivent/nexivent-backend/internal/dao/model/util"
	daoPostgresql "github.com/Nexivent/nexivent-backend/internal/dao/repository"
	schemas "github.com/Nexivent/nexivent-backend/internal/schemas"
	"github.com/Nexivent/nexivent-backend/logging"
	"gorm.io/gorm"
)

const ttlReservaSegundos int64 = 600 // 10 minutos de hold

type StockReservado struct {
	SectorID int64
	Cantidad int64
}

type OrdenDeCompra struct {
	logger        logging.Logger
	DaoPostgresql *daoPostgresql.NexiventPsqlEntidades
}

func NewOrdenDeCompraAdapter(
	logger logging.Logger,
	daoPostgresql *daoPostgresql.NexiventPsqlEntidades,
) *OrdenDeCompra {
	return &OrdenDeCompra{
		logger:        logger,
		DaoPostgresql: daoPostgresql,
	}
}

func (a *OrdenDeCompra) CrearSesionOrdenTemporal(
	req *schemas.CrearOrdenTemporalRequest,
) (*schemas.CrearOrdenTemporalResponse, *errors.Error) {

	if req.IdUsuario == 0 || req.IdEvento == 0 || req.IdFechaEvento == 0 || len(req.Entradas) == 0 {
		return nil, &errors.UnprocessableEntityError.InvalidRequestBody
	}

	// ============================================================================
	// Verificar y reservar stock ANTES de crear la orden
	// ============================================================================

	var total float64 = 0
	stocksReservados := []StockReservado{}

	// 1. Validar stock disponible para cada entrada
	for _, entrada := range req.Entradas {
		sectorID := entrada.IdSector

		a.logger.Infof("🔍 Validando stock: Sector %d, Cantidad solicitada %d", sectorID, entrada.Cantidad)

		// Verificar que haya stock disponible
		var totalEntradas, cantVendidas int64
		row := a.DaoPostgresql.OrdenDeCompra.PostgresqlDB.
			Table("sector").
			Select("total_Entradas, cant_vendidas").
			Where("sector_id = ?", sectorID).
			Row()

		if err := row.Scan(&totalEntradas, &cantVendidas); err != nil {
			a.logger.Errorf("Hold.ObtenerStockSector(sector=%d): %v", sectorID, err)
			a.rollbackStockReservado(stocksReservados)
			return nil, &errors.InternalServerError.Default
		}

		disponible := (cantVendidas + entrada.Cantidad) <= totalEntradas

		if !disponible {
			a.logger.Warnf("❌ Stock insuficiente para sector %d (solicitado: %d, disponible: %d)",
				sectorID, entrada.Cantidad, totalEntradas-cantVendidas)
			a.rollbackStockReservado(stocksReservados)
			return nil, &errors.BadRequestError.EventoNotFound
		}

		// Reservar stock (incrementar cant_vendidas)
		res := a.DaoPostgresql.OrdenDeCompra.PostgresqlDB.
			Table("sector").
			Where("sector_id = ?", sectorID).
			UpdateColumn("cant_vendidas", gorm.Expr("cant_vendidas + ?", entrada.Cantidad))

		if res.Error != nil {
			a.logger.Errorf("Hold.IncrementarVendidasPorSector(sector=%d): %v", sectorID, res.Error)
			a.rollbackStockReservado(stocksReservados)
			return nil, &errors.InternalServerError.Default
		}

		// Guardar para posible rollback
		stocksReservados = append(stocksReservados, StockReservado{
			SectorID: sectorID,
			Cantidad: entrada.Cantidad,
		})

		a.logger.Infof("📉 Stock reservado: Sector %d, Cantidad %d", sectorID, entrada.Cantidad)
	}

	// ============================================================================
	// Crear la orden temporal
	// ============================================================================

	now := time.Now()
	expiresAt := now.Add(time.Duration(ttlReservaSegundos) * time.Second)

	orden := &model.OrdenDeCompra{
		UsuarioID:        req.IdUsuario,
		Fecha:            now,
		FechaHoraIni:     now,
		FechaHoraFin:     &expiresAt,
		Total:            total,
		MontoFeeServicio: 0,
		EstadoDeOrden:    util.OrdenTemporal.Codigo(),
	}

	if err := a.DaoPostgresql.OrdenDeCompra.CrearOrdenTemporal(orden); err != nil {
		a.logger.Errorf("CrearSesionOrdenTemporal: %v", err)
		a.rollbackStockReservado(stocksReservados)
		return nil, &errors.BadRequestError.EventoNotCreated
	}

	a.logger.Infof("✅ Orden temporal %d creada con stock reservado", orden.ID)

	resp := &schemas.CrearOrdenTemporalResponse{
		OrderID:    orden.ID,
		Estado:     "TEMPORAL",
		Total:      orden.Total,
		StartedAt:  orden.FechaHoraIni.Format(time.RFC3339),
		ExpiresAt:  expiresAt.Format(time.RFC3339),
		TTLSeconds: ttlReservaSegundos,
	}
	return resp, nil
}

func (a *OrdenDeCompra) rollbackStockReservado(stocks []StockReservado) {
	for _, stock := range stocks {
		res := a.DaoPostgresql.OrdenDeCompra.PostgresqlDB.
			Table("sector").
			Where("id = ?", stock.SectorID).
			UpdateColumn("cant_vendidas", gorm.Expr("cant_vendidas - ?", stock.Cantidad))

		if res.Error != nil {
			a.logger.Errorf("⚠️ Error al hacer rollback de stock: Sector %d, Cantidad %d: %v",
				stock.SectorID, stock.Cantidad, res.Error)
		} else {
			a.logger.Infof("📈 Rollback stock: Sector %d, Cantidad %d",
				stock.SectorID, stock.Cantidad)
		}
	}
}

func (a *OrdenDeCompra) ObtenerEstadoHold(
	orderID int64,
) (*schemas.ObtenerHoldResponse, *errors.Error) {

	estadoEnum, ini, fin, total, err := a.DaoPostgresql.OrdenDeCompra.ObtenerMetaTemporal(orderID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &errors.ObjectNotFoundError.EventoNotFound
		}
		a.logger.Errorf("ObtenerEstadoHold(%d): %v", orderID, err)
		return nil, &errors.BadRequestError.EventoNotFound
	}

	now := time.Now()
	if fin == nil || now.After(*fin) {
		_ = a.CancelarOrdenYLiberarStock(orderID)
		return nil, &errors.BadRequestError.EventoNotFound
	}

	remaining := int64(fin.Sub(now).Seconds())
	if remaining < 0 {
		remaining = 0
	}

	_ = estadoEnum

	resp := &schemas.ObtenerHoldResponse{
		OrderID:       orderID,
		Estado:        "TEMPORAL",
		RemainingSecs: remaining,
		StartedAt:     ini.Format(time.RFC3339),
		ExpiresAt:     fin.Format(time.RFC3339),
		Total:         total,
	}
	return resp, nil
}

func (a *OrdenDeCompra) ConfirmarOrden(
	orderID int64,
	req *schemas.ConfirmarOrdenRequest,
) (*schemas.ConfirmarOrdenResponse, *errors.Error) {

	ok, err := a.DaoPostgresql.OrdenDeCompra.VerificarOrdenExisteYEstado(orderID, util.OrdenTemporal)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &errors.ObjectNotFoundError.EventoNotFound
		}
		a.logger.Errorf("ConfirmarOrden.Verificar(%d): %v", orderID, err)
		return nil, &errors.BadRequestError.EventoNotFound
	}
	if !ok {
		return nil, &errors.BadRequestError.EventoNotFound
	}

	orden, err := a.DaoPostgresql.OrdenDeCompra.CerrarOrdenTemporal(orderID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &errors.ObjectNotFoundError.EventoNotFound
		}
		a.logger.Errorf("ConfirmarOrden.CerrarTemporal(%d): %v", orderID, err)
		return nil, &errors.BadRequestError.EventoNotFound
	}

	now := time.Now()
	if orden.FechaHoraFin != nil && now.After(*orden.FechaHoraFin) {
		_ = a.CancelarOrdenYLiberarStock(orderID)
		return nil, &errors.BadRequestError.EventoNotFound
	}

	if req.PaymentID == "" {
		return nil, &errors.BadRequestError.EventoNotFound
	}

	metodoPagoID := int64(1)

	if len(req.PaymentID) > 0 {
		var tmpID int64
		if _, scanErr := fmt.Sscanf(req.PaymentID, "%d", &tmpID); scanErr == nil && tmpID > 0 {
			metodoPagoID = tmpID
		}
	}

	a.logger.Infof("ConfirmarOrden: orderID=%d, paymentId=%s, metodoPagoID=%d",
		orderID, req.PaymentID, metodoPagoID)

	if errUpd := a.DaoPostgresql.OrdenDeCompra.ConfirmarOrdenConPago(
		orderID,
		metodoPagoID,
		req.PaymentID,
	); errUpd != nil {
		if errUpd == gorm.ErrRecordNotFound {
			return nil, &errors.ObjectNotFoundError.EventoNotFound
		}
		a.logger.Errorf("ConfirmarOrden.ConfirmarConPago(%d): %v", orderID, errUpd)
		return nil, &errors.BadRequestError.EventoNotFound
	}

	a.logger.Infof("✅ Orden %d confirmada exitosamente con método de pago %d",
		orderID, metodoPagoID)

	montoBruto := orden.Total
	montoFee := orden.MontoFeeServicio
	gananciaNeta := montoBruto - montoFee
	if gananciaNeta < 0 {
		gananciaNeta = 0
	}
	a.logger.Infof(
		"Orden %d: Total=%.2f, FeeServicio=%.2f, GananciaNetaOrganizador=%.2f",
		orderID, montoBruto, montoFee, gananciaNeta,
	)
	// usar evento + fecha del request
	if req.IdEvento > 0 && req.FechaEvento != "" {

		// Parsear fecha que manda el front (YYYY-MM-DD)
		fechaParsed, parseErr := time.Parse("2006-01-02", req.FechaEvento)
		if parseErr != nil {
			a.logger.Errorf("ConfirmarOrden: fechaEvento inválida '%s': %v", req.FechaEvento, parseErr)
			// aquí tú decides si devuelves error o solo logueas
		} else {
			// Llamar al DAO que hace:
			// 1) buscar fecha_evento en TABLA fecha
			// 2) usar (evento_id, fecha_id) en TABLA evento_fecha
			if errSum := a.DaoPostgresql.EventoFecha.SumarGananciaNetaPorEventoYFecha(
				req.IdEvento,
				fechaParsed,
				gananciaNeta,
			); errSum != nil {
				a.logger.Errorf(
					"ConfirmarOrden.SumarGananciaNetaPorEventoYFecha(order=%d, evento=%d, fecha=%s, monto=%.2f): %v",
					orderID, req.IdEvento, fechaParsed.Format("2006-01-02"), gananciaNeta, errSum,
				)
				// igual, aquí decides si solo log o error
			} else {
				a.logger.Infof(
					"Ganancia acumulada para evento=%d, fecha=%s: +%.2f (orden=%d)",
					req.IdEvento, fechaParsed.Format("2006-01-02"), gananciaNeta, orderID,
				)
			}
		}
	} else {
		a.logger.Warnf(
			"ConfirmarOrden: IdEvento o FechaEvento no enviados (evento=%d, fecha='%s'), no se actualiza ganancia.",
			req.IdEvento, req.FechaEvento,
		)
	}

	resp := &schemas.ConfirmarOrdenResponse{
		OrderID: orderID,
		Estado:  "CONFIRMADA",
		Mensaje: "Compra confirmada",
	}

	return resp, nil
}

func (a *OrdenDeCompra) CancelarOrdenYLiberarStock(orderID int64) *errors.Error {
	ok, err := a.DaoPostgresql.OrdenDeCompra.VerificarOrdenExisteYEstado(orderID, util.OrdenTemporal)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &errors.ObjectNotFoundError.EventoNotFound
		}
		a.logger.Errorf("CancelarOrden.Verificar(%d): %v", orderID, err)
		return &errors.InternalServerError.Default
	}
	if !ok {
		a.logger.Warnf("Orden %d no está en estado TEMPORAL", orderID)
		return &errors.BadRequestError.EventoNotFound
	}

	// Obtener detalles de la orden para liberar stock
	type DetalleConSector struct {
		SectorID int64
		Cantidad int64
	}
	var detalles []DetalleConSector

	err = a.DaoPostgresql.OrdenDeCompra.PostgresqlDB.
		Table("orden_de_compra_detalle").
		Select("id_sector as sector_id, cantidad").
		Where("orden_de_compra_id = ?", orderID).
		Find(&detalles).Error

	if err != nil {
		a.logger.Errorf("CancelarOrden.ObtenerDetalles(%d): %v", orderID, err)
	} else {
		for _, d := range detalles {
			res := a.DaoPostgresql.OrdenDeCompra.PostgresqlDB.
				Table("sector").
				Where("id = ?", d.SectorID).
				UpdateColumn("cant_vendidas", gorm.Expr("cant_vendidas - ?", d.Cantidad))

			if res.Error != nil {
				a.logger.Errorf("CancelarOrden.DecrementarStock(sector=%d): %v", d.SectorID, res.Error)
			} else {
				a.logger.Infof("📈 Stock liberado: Sector %d, Cantidad %d (Orden %d cancelada)",
					d.SectorID, d.Cantidad, orderID)
			}
		}
	}

	if err := a.DaoPostgresql.OrdenDeCompra.ActualizarEstadoOrden(orderID, util.OrdenCancelada); err != nil {
		a.logger.Errorf("CancelarOrden.ActualizarEstado(%d): %v", orderID, err)
		return &errors.InternalServerError.Default
	}

	a.logger.Infof("✅ Orden %d cancelada y stock liberado", orderID)
	return nil
}
