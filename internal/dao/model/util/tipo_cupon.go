package model

import "fmt"

type TipoCupon int16

const (
	TipoPorcentaje TipoCupon = iota // 0
	TipoMonto                       // 1
)

// Retorna el código numérico
func (t TipoCupon) Codigo() int16 {
	return int16(t)
}

// Convierte código numérico a enum
func ValueOfTipoCuponCodigo(c int16) (TipoCupon, error) {
	switch c {
	case 0:
		return TipoPorcentaje, nil
	case 1:
		return TipoMonto, nil
	default:
		return 0, fmt.Errorf("código de tipo de cupón inválido: %d", c)
	}
}

// Nombre legible
func (t TipoCupon) String() string {
	switch t {
	case TipoPorcentaje:
		return "PORCENTAJE"
	case TipoMonto:
		return "MONTO"
	default:
		return "DESCONOCIDO"
	}
}

// Validación
func (t TipoCupon) IsValid() bool {
	return t == TipoPorcentaje || t == TipoMonto
}
