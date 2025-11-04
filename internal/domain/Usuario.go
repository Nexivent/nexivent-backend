package domain

import (
	"time"
)

type Usuario struct {
	ID                    int64 `gorm:"column:usuario_id;primaryKey;autoIncrement"`
	Nombre                string
	TipoDocumento         string
	NumDocumento          string
	Correo                string
	Contrasenha           string
	Telefono              *string
	EstadoDeCuenta        int16
	CodigoVerificacion    *string
	FechaExpiracionCodigo *time.Time
	UsuarioCreacion       *int64
	FechaCreacion         time.Time
	UsuarioModificacion   *int64
	FechaModificacion     *time.Time
	Estado                int16

	Comentarios    []Comentario
	Ordenes        []OrdenDeCompra
	RolesAsignados []RolUsuario
	Cupones        []UsuarioCupon
}

func (Usuario) TableName() string { return "usuario" }
