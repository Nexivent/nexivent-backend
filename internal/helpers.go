package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Nexivent/nexivent-backend/internal/validator"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type Envelope map[string]any

func ReadIDParam(r *http.Request) (uuid.UUID, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := uuid.Parse(params.ByName("id"))

	if err != nil {
		return uuid.Nil, errors.New("invalid id parameter")
	}

	return id, nil
}

func WriteJSON(w http.ResponseWriter, status int, data Envelope, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}
	js = append(js, '\n')

	maps.Copy(w.Header(), headers)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func ReadJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	// Usar http.MaxBytesReader() para limitar el tamaño del cuerpo de la petición a 1MB.
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Decodificar el cuerpo de la petición en el destino objetivo
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		// Si hay un error durante la decodificación, iniciar el diagnóstico...
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {

		// Usar la función errors.As() para verificar si el error tiene el tipo
		// *json.SyntaxError. Si es así, devolver un mensaje de error en lenguaje simple
		// que incluya la ubicación del problema.
		case errors.As(err, &syntaxError):
			return fmt.Errorf("el cuerpo contiene JSON mal formado (en el carácter %d)", syntaxError.Offset)

		// En algunas circunstancias Decode() también puede devolver un error io.ErrUnexpectedEOF
		// para errores de sintaxis en el JSON. Así que verificamos esto usando errors.Is() y
		// devolvemos un mensaje de error genérico. Hay un issue abierto sobre esto en
		// https://github.com/golang/go/issues/25956.
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("el cuerpo contiene JSON mal formado")

		// De igual manera, capturar cualquier error *json.UnmarshalTypeError. Estos ocurren cuando
		// el valor JSON es del tipo incorrecto para el destino objetivo. Si el error se relaciona
		// con un campo específico, entonces lo incluimos en nuestro mensaje de error para hacer más
		// fácil la depuración para el cliente.
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("el cuerpo contiene un tipo JSON incorrecto para el campo %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("el cuerpo contiene un tipo JSON incorrecto (en el carácter %d)", unmarshalTypeError.Offset)

		// Un error io.EOF será devuelto por Decode() si el cuerpo de la petición está vacío.
		// Verificamos esto con errors.Is() y devolvemos un mensaje de error en lenguaje simple
		// en su lugar.
		case errors.Is(err, io.EOF):
			return errors.New("el cuerpo no debe estar vacío")

		// Si el JSON contiene un campo que no puede ser mapeado al destino objetivo,
		// entonces Decode() ahora devolverá un mensaje de error en el formato "json: unknown
		// field "<name>"". Verificamos esto, extraemos el nombre del campo del error,
		// y lo interpolamos en nuestro mensaje de error personalizado. Nota que hay un issue
		// abierto en https://github.com/golang/go/issues/29035 sobre convertir esto
		// en un tipo de error distinto en el futuro.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		// Usar la función errors.As() para verificar si el error tiene el tipo
		// *http.MaxBytesError. Si es así, significa que el cuerpo de la petición excedió nuestro
		// límite de tamaño de 1MB y devolvemos un mensaje de error claro.
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)

		// Un error json.InvalidUnmarshalError será devuelto si pasamos algo
		// que no es un puntero no nulo a Decode(). Capturamos esto y hacemos panic,
		// en lugar de devolver un error a nuestro manejador. Al final de este capítulo
		// hablaremos sobre hacer panic versus devolver errores, y discutiremos por qué es
		// apropiado hacerlo en esta situación específica.
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		// Para cualquier otra cosa, devolver el mensaje de error tal cual.
		default:
			return err
		}
	}

	// Llamar a Decode() nuevamente, usando un puntero a un struct anónimo vacío como
	// destino. Si el cuerpo de la petición solo contenía un único valor JSON, esto devolverá
	// un error io.EOF. Entonces, si obtenemos cualquier otra cosa, sabemos que hay
	// datos adicionales en el cuerpo de la petición y devolvemos nuestro propio mensaje de error
	// personalizado.
	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

// The readString() helper returns a string value from the query string, or the provided
// default value if no matching key could be found.
func ReadString(qs url.Values, key string, defaultValue string) string {
	// Extract the value for a given key from the query string. If no key exists this
	// will return the empty string "".
	s := qs.Get(key)

	// If no key exists (or the value is empty) then return the default value.
	if s == "" {
		return defaultValue
	}

	// Otherwise return the string.
	return s
}

// The readInt() helper reads a string value from the query string and converts it to an
// integer before returning. If no matching key could be found it returns the provided
// default value. If the value couldn't be converted to an integer, then we record an
// error message in the provided Validator instance.
func ReadInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	// Extract the value from the query string.
	s := qs.Get(key)
	// If no key exists (or the value is empty) then return the default value.
	if s == "" {
		return defaultValue
	}
	// Try to convert the value to an int. If this fails, add an error message to the
	// validator instance and return the default value.
	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}
	// Otherwise, return the converted integer value.
	return i
}

func ReadCSV(qs url.Values, key string, defaultValue []string) []string {
	// Extract the value from the query string.
	csv := qs.Get(key)

	// If no key exists (or the value is empty) then return the default value.
	if csv == "" {
		return defaultValue

	}
	// Otherwise parse the value into a []string slice and return it.
	return strings.Split(csv, ",")
}
