package model

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/model/util"
	"github.com/Nexivent/nexivent-backend/internal/validator"
	"golang.org/x/crypto/argon2"
)

const (
	argonTime    uint32 = 1         // Number of iterations
	argonMemory  uint32 = 64 * 1024 // 64 MB
	argonThreads uint8  = 4         // Number of threads
	argonKeyLen  uint32 = 32        // 32-byte hash
	saltLength          = 16        // 16-byte salt
)

type Usuario struct {
	ID                    uint64             `gorm:"column:usuario_id;primaryKey;autoIncrement" json:"id"`
	Nombre                string             `gorm:"column:nombre" json:"nombre"`
	TipoDocumento         util.TipoDocumento `gorm:"column:tipo_documento;uniqueIndex:uq_usuario_doc" json:"tipoDocumento"`
	NumDocumento          string             `gorm:"column:num_documento;uniqueIndex:uq_usuario_doc" json:"numDocumento"`
	Correo                string             `gorm:"column:correo;uniqueIndex" json:"correo"`
	Password              []byte             `gorm:"column:password" json:"-"`
	Telefono              *string            `gorm:"column:telefono" json:"telefono,omitempty"`
	EstadoDeCuenta        util.Estado        `gorm:"column:estado_de_cuenta;default:0" json:"estadoDeCuenta"`
	CodigoVerificacion    *string            `gorm:"column:codigo_verificacion" json:"-"`
	FechaExpiracionCodigo *time.Time         `gorm:"column:fecha_expiracion_codigo" json:"-"`
	UsuarioCreacion       *uint64            `gorm:"column:usuario_creacion" json:"-"`
	FechaCreacion         time.Time          `gorm:"column:fecha_creacion;default:now()" json:"-"`
	UsuarioModificacion   *uint64            `gorm:"column:usuario_modificacion" json:"-"`
	FechaModificacion     *time.Time         `gorm:"column:fecha_modificacion" json:"-"`
	Estado                util.Estado        `gorm:"column:estado;default:0" json:"-"`

	Comentarios    []Comentario    `json:"comentarios,"`
	Ordenes        []OrdenDeCompra `json:"ordenes"`
	RolesAsignados []RolUsuario    `json:"roles"`
	Cupones        []UsuarioCupon  `json:"cupones"`
}

func (Usuario) TableName() string { return "usuario" }

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUsuario(v *validator.Validator, usuario *Usuario, plaintext *string) {
	// Validar Nombre
	v.Check(usuario.Nombre != "", "nombre", "el nombre es obligatorio")
	v.Check(len(usuario.Nombre) <= 100, "nombre", "el nombre no debe exceder 100 caracteres")

	// Validar TipoDocumento
	v.Check(usuario.TipoDocumento != "", "tipoDocumento", "el tipo de documento es obligatorio")
	v.Check(usuario.TipoDocumento == "DNI" || usuario.TipoDocumento == "CE" || usuario.TipoDocumento == "RUC",
		"tipoDocumento", "el tipo de documento debe ser DNI, CE o RUC")

	// Validar NumDocumento
	v.Check(usuario.NumDocumento != "", "numDocumento", "el número de documento es obligatorio")
	v.Check(len(usuario.NumDocumento) <= 20, "numDocumento", "el número de documento no debe exceder 20 caracteres")

	// Validar Correo
	ValidateEmail(v, usuario.Correo)

	if plaintext != nil {
		ValidatePasswordPlaintext(v, *plaintext)
	}

	// Validar Telefono (si está presente)
	if usuario.Telefono != nil && *usuario.Telefono != "" {
		v.Check(len(*usuario.Telefono) <= 15, "telefono", "el teléfono no debe exceder 15 caracteres")
	}
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func HashPassword(password string) (string, error) {
	salt, err := generateRandomBytes(saltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)

	// Encode parameters, salt and hash into a single string
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		argonMemory, argonTime, argonThreads, b64Salt, b64Hash)

	return encodedHash, nil
}

func VerifyPassword(password, encodedHash string) (bool, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, fmt.Errorf("invalid hash format")
	}

	// parts[0] = "" (before first $)
	// parts[1] = "argon2id"
	// parts[2] = "v=19"
	// parts[3] = "m=65536,t=1,p=4"
	// parts[4] = salt (b64)
	// parts[5] = hash (b64)

	var memory uint32
	var time uint32
	var threads uint8

	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	// Recompute the hash with the same parameters
	computedHash := argon2.IDKey([]byte(password), salt, time, memory, threads, uint32(len(hash)))

	// Constant-time comparison
	if subtle.ConstantTimeCompare(hash, computedHash) == 1 {
		return true, nil
	}

	return false, nil
}
