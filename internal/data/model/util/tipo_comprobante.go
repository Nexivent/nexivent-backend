package util

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type TipoComprobante int16

const (
	ComprobanteBoleta TipoComprobante = iota
	ComprobanteFactura
)

var ErrInvalidTipoComprobante = errors.New("tipo de comprobante inválido")

func (t TipoComprobante) Codigo() int16 { return int16(t) }

func ValueOfTipoComprobanteCodigo(c int16) (TipoComprobante, error) {
	switch c {
	case 0:
		return ComprobanteBoleta, nil
	case 1:
		return ComprobanteFactura, nil
	default:
		return 0, fmt.Errorf("código de tipo de comprobante inválido: %d", c)
	}
}

func (t TipoComprobante) String() string {
	switch t {
	case ComprobanteBoleta:
		return "BOLETA"
	case ComprobanteFactura:
		return "FACTURA"
	default:
		return "DESCONOCIDO"
	}
}

func (t TipoComprobante) IsValid() bool {
	return t == ComprobanteBoleta || t == ComprobanteFactura
}

// ---- Integración con database/sql ----
func (t TipoComprobante) Value() (driver.Value, error) {
	if !t.IsValid() {
		return nil, fmt.Errorf("tipo de comprobante inválido: %d", t)
	}
	return int64(t), nil // Postgres castea a SMALLINT sin problema
}

func (t *TipoComprobante) Scan(src any) error {
	switch v := src.(type) {
	case int64:
		*t = TipoComprobante(v)
	case int32:
		*t = TipoComprobante(v)
	case int16:
		*t = TipoComprobante(v)
	case []byte:
		var n int16
		if _, err := fmt.Sscanf(string(v), "%d", &n); err != nil {
			return fmt.Errorf("scan TipoComprobante: %w", err)
		}
		*t = TipoComprobante(n)
	case string:
		var n int16
		if _, err := fmt.Sscanf(v, "%d", &n); err != nil {
			return fmt.Errorf("scan TipoComprobante: %w", err)
		}
		*t = TipoComprobante(n)
	default:
		return fmt.Errorf("tipo no soportado para TipoComprobante: %T", src)
	}
	if !t.IsValid() {
		return fmt.Errorf("tipo de comprobante inválido: %d", *t)
	}
	return nil
}

func (t *TipoComprobante) UnmarshalJSON(data []byte) error {
	unquotedData, err := strconv.Unquote(string(data))
	if err != nil {
		return ErrInvalidTipoComprobante
	}

	// Split the string to isolate the part containing the number.
	parts := strings.Split(unquotedData, " ")

	// Sanity check the parts of the string to make sure it was in the expected format.
	// If it isn't, we return the ErrInvalidTipoComprobante error again.
	if len(parts) > 1 {
		return ErrInvalidTipoComprobante
	}

	switch strings.ToUpper(parts[0]) {
	case "BOLETA":
		*t = ComprobanteBoleta
	case "FACTURA":
		*t = ComprobanteFactura
	default:
		return ErrInvalidTipoComprobante
	}

	return nil
}

func (t TipoComprobante) MarshalJSON() ([]byte, error) {
	return fmt.Appendf([]byte{}, `"%s"`, t.String()), nil
}
