package model

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/argon2"
)

type Usuario struct {
	ID                    int64 `gorm:"column:usuario_id;primaryKey;autoIncrement"`
	Nombre                string
	TipoDocumento         string `gorm:"uniqueIndex:uq_usuario_doc"`
	NumDocumento          string `gorm:"uniqueIndex:uq_usuario_doc"`
	Correo                string `gorm:"uniqueIndex"`
	Contrasenha           string
	Telefono              *string
	EstadoDeCuenta        int16 `gorm:"default:0"`
	CodigoVerificacion    *string
	FechaExpiracionCodigo *time.Time
	UsuarioCreacion       *int64
	FechaCreacion         time.Time `gorm:"default:now()"`
	UsuarioModificacion   *int64
	FechaModificacion     *time.Time
	Estado                int16 `gorm:"default:1"`

	Comentarios    []Comentario
	Ordenes        []OrdenDeCompra
	RolesAsignados []RolUsuario
	Cupones        []UsuarioCupon
}

func (Usuario) TableName() string { return "usuario" }

const (
	argonTime    uint32 = 1         // Number of iterations
	argonMemory  uint32 = 64 * 1024 // 64 MB
	argonThreads uint8  = 4         // Number of threads
	argonKeyLen  uint32 = 32        // 32-byte hash
	saltLength          = 16        // 16-byte salt
)

var AnonymousUser = &Usuario{}

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

func (u *Usuario) IsAnonymous() bool {
	return u == AnonymousUser
}
