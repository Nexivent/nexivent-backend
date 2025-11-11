package util

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type TipoDocumento string

const (
	TipoDocDNI TipoDocumento = "DNI"
	TipoDocCE  TipoDocumento = "CE"
	TipoDocRUC TipoDocumento = "RUC"
)

var ErrInvalidTipoDocumento = errors.New("tipo de documento inválido")

func (t TipoDocumento) String() string {
	return string(t)
}

func (t TipoDocumento) IsValid() bool {
	return t == TipoDocDNI || t == TipoDocCE || t == TipoDocRUC
}

func (t *TipoDocumento) UnmarshalJSON(data []byte) error {
	unquotedData, err := strconv.Unquote(string(data))
	if err != nil {
		return ErrInvalidTipoDocumento
	}

	// Split the string to isolate the part containing the number.
	parts := strings.Split(unquotedData, " ")

	// Sanity check the parts of the string to make sure it was in the expected format.
	// If it isn't, we return the ErrInvalidTipoDocumento error again.
	if len(parts) > 1 {
		return ErrInvalidTipoDocumento
	}

	switch strings.ToUpper(parts[0]) {
	case "DNI":
		*t = TipoDocDNI
	case "CE":
		*t = TipoDocCE
	case "RUC":
		*t = TipoDocRUC
	default:
		return ErrInvalidTipoDocumento
	}

	return nil
}

func (t TipoDocumento) MarshalJSON() ([]byte, error) {
	return fmt.Appendf([]byte{}, `"%s"`, t.String()), nil
}
