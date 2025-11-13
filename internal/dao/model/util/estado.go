package model

import (
	"database/sql/driver"
	"fmt"
)

type Estado int16

const (
	Inactivo Estado = iota // 0
	Activo                 // 1
)

func (t Estado) Codigo() int16 { return int16(t) }

func ValueOfEstadoCodigo(c int16) (Estado, error) {
	switch c {
	case 0:
		return Activo, nil
	case 1:
		return Inactivo, nil
	default:
		return 0, fmt.Errorf("código de estado de tabla inválido: %d", c)
	}
}

func (t Estado) String() string {
	switch t {
	case Activo:
		return "Activo"
	case Inactivo:
		return "Inactivo"
	default:
		return "DESCONOCIDO"
	}
}

func (t Estado) IsValid() bool {
	return t == Activo || t == Inactivo
}

// ---- Integración con database/sql ----
func (t Estado) Value() (driver.Value, error) {
	if !t.IsValid() {
		return nil, fmt.Errorf("estado de tabla inválido: %d", t)
	}
	return int64(t), nil // Postgres castea a SMALLINT sin problema
}

func (t *Estado) Scan(src any) error {
	switch v := src.(type) {
	case int64:
		*t = Estado(v)
	case int32:
		*t = Estado(v)
	case int16:
		*t = Estado(v)
	case []byte:
		var n int16
		if _, err := fmt.Sscanf(string(v), "%d", &n); err != nil {
			return fmt.Errorf("scan Estado: %w", err)
		}
		*t = Estado(n)
	case string:
		var n int16
		if _, err := fmt.Sscanf(v, "%d", &n); err != nil {
			return fmt.Errorf("scan Estado: %w", err)
		}
		*t = Estado(n)
	default:
		return fmt.Errorf("tipo no soportado para Estado: %T", src)
	}
	if !t.IsValid() {
		return fmt.Errorf("tipo de comprobante inválido: %d", *t)
	}
	return nil
}
