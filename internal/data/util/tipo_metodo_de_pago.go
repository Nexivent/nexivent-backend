package util

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type TipoMetodoPago string

const (
	MetodoTarjeta TipoMetodoPago = "Tarjeta"
	MetodoYape    TipoMetodoPago = "Yape"
)

var ErrInvalidTipoMetodoPago = errors.New("tipo de método de pago inválido")

func (t TipoMetodoPago) String() string {
	return string(t)
}

func (t TipoMetodoPago) IsValid() bool {
	return t == MetodoTarjeta || t == MetodoYape
}

func (t *TipoMetodoPago) UnmarshalJSON(data []byte) error {
	unquotedData, err := strconv.Unquote(string(data))
	if err != nil {
		return ErrInvalidTipoMetodoPago
	}

	// Split the string to isolate the part containing the number.
	parts := strings.Split(unquotedData, " ")

	// Sanity check the parts of the string to make sure it was in the expected format.
	// If it isn't, we return the ErrInvalidTipoMetodoPago error again.
	if len(parts) > 1 {
		return ErrInvalidTipoMetodoPago
	}

	switch strings.ToUpper(parts[0]) {
	case "TARJETA":
		*t = MetodoTarjeta
	case "YAPE":
		*t = MetodoYape
	default:
		return ErrInvalidTipoMetodoPago
	}

	return nil
}

func (t TipoMetodoPago) MarshalJSON() ([]byte, error) {
	return fmt.Appendf([]byte{}, `"%s"`, t.String()), nil
}
