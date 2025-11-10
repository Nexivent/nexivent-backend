package model

type UsuarioCupon struct {
	CuponID   int64 `gorm:"primaryKey"`
	UsuarioID int64 `gorm:"primaryKey"`
	CantUsada int64 `gorm:"default:0"`

	Cupon   *Cupon   `gorm:"foreignKey:CuponID;references:cupon_id"`
	Usuario *Usuario `gorm:"foreignKey:UsuarioID;references:usuario_id"`
}

func (UsuarioCupon) TableName() string { return "usuario_cupon" }
