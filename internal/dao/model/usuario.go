package model

import (
	"time"
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
