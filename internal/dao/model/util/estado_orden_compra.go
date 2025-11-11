package model

import (
	"database/sql/driver"
	"fmt"
)

// EstadoOrden modela el WORKFLOW de la orden (columna: estado_de_orden)
// 0=TEMPORAL, 1=CONFIRMADA, 2=CANCELADA
type EstadoOrden int16

const (
	OrdenTemporal   EstadoOrden = iota // 0
	OrdenConfirmada                    // 1
	OrdenCancelada                     // 2
)

func (e EstadoOrden) Codigo() int16 { return int16(e) }

func (e EstadoOrden) String() string {
	switch e {
	case OrdenTemporal:
		return "TEMPORAL"
	case OrdenConfirmada:
		return "CONFIRMADA"
	case OrdenCancelada:
		return "CANCELADA"
	default:
		return "DESCONOCIDO"
	}
}

func (e EstadoOrden) IsValid() bool {
	return e == OrdenTemporal || e == OrdenConfirmada || e == OrdenCancelada
}

/* ---- Integración con database/sql (columna SMALLINT) ---- */

func (e EstadoOrden) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, fmt.Errorf("estado_de_orden inválido: %d", e)
	}
	return int64(e), nil
}

func (e *EstadoOrden) Scan(src any) error {
	switch v := src.(type) {
	case int64:
		*e = EstadoOrden(v)
	case int32:
		*e = EstadoOrden(v)
	case int16:
		*e = EstadoOrden(v)
	case []byte:
		var n int16
		if _, err := fmt.Sscanf(string(v), "%d", &n); err != nil {
			return fmt.Errorf("scan EstadoOrden: %w", err)
		}
		*e = EstadoOrden(n)
	case string:
		var n int16
		if _, err := fmt.Sscanf(v, "%d", &n); err != nil {
			return fmt.Errorf("scan EstadoOrden: %w", err)
		}
		*e = EstadoOrden(n)
	default:
		return fmt.Errorf("tipo no soportado para EstadoOrden: %T", src)
	}
	if !e.IsValid() {
		return fmt.Errorf("estado_de_orden inválido: %d", *e)
	}
	return nil
}
