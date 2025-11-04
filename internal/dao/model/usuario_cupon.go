package model

type UsuarioCupon struct {
	CuponID   int64 `gorm:"primaryKey"`
	UsuarioID int64 `gorm:"primaryKey"`
	CantUsada int64

	Cupon   *Cupon   `gorm:"foreignKey:CuponID"`
	Usuario *Usuario `gorm:"foreignKey:UsuarioID"`
}

func (UsuarioCupon) TableName() string { return "usuario_cupon" }
