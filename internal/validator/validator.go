package validator

import (
	"regexp"
	"slices"
)

var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// Define un nuevo tipo Validator que contiene un mapa de errores de validación.
type Validator struct {
	Errors map[string]string
}

// New es un helper que crea una nueva instancia de Validator con un mapa de errores vacío.
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Valid devuelve true si el mapa de errores no contiene ninguna entrada.
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError agrega un mensaje de error al mapa (siempre y cuando no exista ya una entrada
// para la clave dada).
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check agrega un mensaje de error al mapa solo si una validación no es 'ok'.
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// PermittedValue es una función genérica que devuelve true si un valor específico está en una lista
// de valores permitidos.
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

// Matches devuelve true si un valor de cadena coincide con un patrón de expresión regular
// específico.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// Unique es una función genérica que devuelve true si todos los valores en un slice son únicos.
func Unique[T comparable](values []T) bool {
	uniqueValues := make(map[T]bool)
	for _, value := range values {
		uniqueValues[value] = true
	}
	return len(values) == len(uniqueValues)
}
