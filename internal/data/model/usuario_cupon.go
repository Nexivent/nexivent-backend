package model

type UsuarioCupon struct {
	CuponID   uint64
	UsuarioID uint64
	CantUsada int64
}

func (UsuarioCupon) TableName() string { return "usuario_cupon" }
