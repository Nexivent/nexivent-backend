package util

import (
	"database/sql/driver"
	"fmt"
)

type EstadoOrdenCompra int16

const (
	OrdenExitoso EstadoOrdenCompra = iota // 0
	OrdenFallido                          // 1
)

func (e EstadoOrdenCompra) Codigo() int16 { return int16(e) }

func ValueOfEstadoOrdenCompraCodigo(c int16) (EstadoOrdenCompra, error) {
	switch c {
	case 0:
		return OrdenExitoso, nil
	case 1:
		return OrdenFallido, nil
	default:
		return 0, fmt.Errorf("código de estado de orden inválido: %d", c)
	}
}

func (e EstadoOrdenCompra) String() string {
	switch e {
	case OrdenExitoso:
		return "EXITOSO"
	case OrdenFallido:
		return "FALLIDO"
	default:
		return "DESCONOCIDO"
	}
}

func (e EstadoOrdenCompra) IsValid() bool {
	return e == OrdenExitoso || e == OrdenFallido
}

/* ---- Integración con database/sql (columna SMALLINT) ---- */

func (e EstadoOrdenCompra) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, fmt.Errorf("estado de orden inválido: %d", e)
	}
	return int64(e), nil
}

func (e *EstadoOrdenCompra) Scan(src any) error {
	switch v := src.(type) {
	case int64:
		*e = EstadoOrdenCompra(v)
	case int32:
		*e = EstadoOrdenCompra(v)
	case int16:
		*e = EstadoOrdenCompra(v)
	case []byte:
		var n int16
		if _, err := fmt.Sscanf(string(v), "%d", &n); err != nil {
			return fmt.Errorf("scan EstadoOrdenCompra: %w", err)
		}
		*e = EstadoOrdenCompra(n)
	case string:
		var n int16
		if _, err := fmt.Sscanf(v, "%d", &n); err != nil {
			return fmt.Errorf("scan EstadoOrdenCompra: %w", err)
		}
		*e = EstadoOrdenCompra(n)
	default:
		return fmt.Errorf("tipo no soportado para EstadoOrdenCompra: %T", src)
	}
	if !e.IsValid() {
		return fmt.Errorf("estado de orden inválido: %d", *e)
	}
	return nil
}
