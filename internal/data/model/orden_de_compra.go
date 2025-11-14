package model

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/model/util"
	"github.com/Nexivent/nexivent-backend/internal/validator"
)

type OrdenDeCompra struct {
	ID               uint64           `gorm:"column:orden_de_compra_id;primaryKey;autoIncrement" json:"id"`
	UsuarioID        uint64           `gorm:"column:usuario_id" json:"usuarioId"`
	MetodoDePagoID   uint64           `gorm:"column:metodo_de_pago_id" json:"metodoDePagoId"`
	Fecha            time.Time        `gorm:"column:fecha;default:current_date" json:"fecha"`
	FechaHoraIni     time.Time        `gorm:"column:fecha_hora_ini;default:now()" json:"fechaHoraIni"`
	FechaHoraFin     *time.Time       `gorm:"column:fecha_hora_fin" json:"fechaHoraFin,omitempty"`
	Total            float64          `gorm:"column:total" json:"total"`
	MontoFeeServicio float64          `gorm:"column:monto_fee_servicio" json:"montoFeeServicio"`
	EstadoDeOrden    util.EstadoOrden `gorm:"column:estado_de_orden;default:0" json:"estadoDeOrden"`

	// Usuario      *Usuario      `gorm:"foreignKey:UsuarioID;references:usuario_id"`
	// MetodoDePago *MetodoDePago `gorm:"foreignKey:MetodoDePagoID;references:metodo_de_pago_id"`

	Tickets          []Ticket            `json:"tickets"`
	ComprobantesPago []ComprobanteDePago `json:"comprobantesPago"`
}

func (OrdenDeCompra) TableName() string { return "orden_de_compra" }

func ValidateOrdenDeCompra(v *validator.Validator, orden *OrdenDeCompra) {
	// Validar IDs
	v.Check(orden.UsuarioID != 0, "usuarioId", "el ID del usuario es obligatorio")
	v.Check(orden.MetodoDePagoID != 0, "metodoDePagoId", "el ID del método de pago es obligatorio")

	// Validar montos
	v.Check(orden.Total >= 0, "total", "el total no puede ser negativo")
	v.Check(orden.MontoFeeServicio >= 0, "montoFeeServicio", "el monto del fee de servicio no puede ser negativo")

	// Validar rango de fechas
	if orden.FechaHoraFin != nil {
		v.Check(!orden.FechaHoraFin.Before(orden.FechaHoraIni), "fechaHoraFin", "la fecha de fin debe ser mayor o igual a la fecha de inicio")
	}
}
