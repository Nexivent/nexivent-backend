package adapter

import (
	"time"
	"fmt"
	"github.com/Nexivent/nexivent-backend/errors"
	"github.com/Nexivent/nexivent-backend/internal/dao/model"
	util "github.com/Nexivent/nexivent-backend/internal/dao/model/util"
	daoPostgresql "github.com/Nexivent/nexivent-backend/internal/dao/repository"
	schemas "github.com/Nexivent/nexivent-backend/internal/schemas"
	"github.com/Nexivent/nexivent-backend/logging"
	"gorm.io/gorm"
)

const ttlReservaSegundos int64 = 600 // 10 minutos de hold

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

	// Validación mínima según contrato
	if req.IdUsuario == 0 || req.IdEvento == 0 || req.IdFechaEvento == 0 || len(req.Entradas) == 0 {
		// 400: datos inválidos
		return nil, &errors.UnprocessableEntityError.InvalidRequestBody
	}

	// TODO: validar stock de cada tarifa / cantidad
	// TODO: calcular total = sum(precioTarifa * cantidad)
	var total float64 = 0

	now := time.Now()
	expiresAt := now.Add(time.Duration(ttlReservaSegundos) * time.Second)

	orden := &model.OrdenDeCompra{
		UsuarioID:        req.IdUsuario,
		Fecha:            now,
		FechaHoraIni:     now,
		FechaHoraFin:     &expiresAt,
		Total:            total,
		MontoFeeServicio: 0,
		EstadoDeOrden:    util.OrdenTemporal.Codigo(), // el DAO igual lo refuerza
	}

	if err := a.DaoPostgresql.OrdenDeCompra.CrearOrdenTemporal(orden); err != nil {
		a.logger.Errorf("CrearSesionOrdenTemporal: %v", err)
		// Ajusta este error a algo tipo OrdenNotCreated si lo creas
		return nil, &errors.BadRequestError.EventoNotCreated
	}

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

func (a *OrdenDeCompra) ObtenerEstadoHold(
	orderID int64,
) (*schemas.ObtenerHoldResponse, *errors.Error) {

	estadoEnum, ini, fin, total, err := a.DaoPostgresql.OrdenDeCompra.ObtenerMetaTemporal(orderID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 404: reserva no encontrada
			return nil, &errors.ObjectNotFoundError.EventoNotFound
		}
		a.logger.Errorf("ObtenerEstadoHold(%d): %v", orderID, err)
		return nil, &errors.BadRequestError.EventoNotFound
	}

	now := time.Now()
	if fin == nil || now.After(*fin) {
		// opcional: marcar como cancelada
		_ = a.DaoPostgresql.OrdenDeCompra.ActualizarEstadoOrden(orderID, util.OrdenCancelada)
		// 410: reserva expirada
		// (usa un error más específico si lo defines)
		return nil, &errors.BadRequestError.EventoNotFound
	}

	remaining := int64(fin.Sub(now).Seconds())
	if remaining < 0 {
		remaining = 0
	}

	// El contrato dice estado: "BORRADOR" aunque en BD esté TEMPORAL
	_ = estadoEnum // por si luego mapeas a otros strings
	

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

	// 1) Verificar que exista y esté en estado TEMPORAL
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

	// 2) Lock de la fila (FOR UPDATE) sobre la orden temporal
	orden, err := a.DaoPostgresql.OrdenDeCompra.CerrarOrdenTemporal(orderID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &errors.ObjectNotFoundError.EventoNotFound
		}
		a.logger.Errorf("ConfirmarOrden.CerrarTemporal(%d): %v", orderID, err)
		return nil, &errors.BadRequestError.EventoNotFound
	}

	// 3) Verificar expiración
	now := time.Now()
	if orden.FechaHoraFin != nil && now.After(*orden.FechaHoraFin) {
		_ = a.DaoPostgresql.OrdenDeCompra.ActualizarEstadoOrden(orderID, util.OrdenCancelada)
		return nil, &errors.BadRequestError.EventoNotFound
	}

	// 4) Verificar pago
	if req.PaymentID == "" {
		return nil, &errors.BadRequestError.EventoNotFound
	}

	// 5) Parsear método de pago desde paymentId
	metodoPagoID := int64(1) // Default: Tarjeta
	
	if len(req.PaymentID) > 0 {
		var tmpID int64
		if _, scanErr := fmt.Sscanf(req.PaymentID, "%d", &tmpID); scanErr == nil && tmpID > 0 {
			metodoPagoID = tmpID
		}
	}

	a.logger.Infof("ConfirmarOrden: orderID=%d, paymentId=%s, metodoPagoID=%d", 
		orderID, req.PaymentID, metodoPagoID)

	// 6) Actualizar estado Y método de pago
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

	resp := &schemas.ConfirmarOrdenResponse{
		OrderID: orderID,
		Estado:  "CONFIRMADA",
		Mensaje: "Compra confirmada",
	}
	return resp, nil
}
