package domain

type UsuarioCupon struct {
	Cupon     Cupon   `db:"-" json:"cupon"`   // PK compuesta
	Usuario   Usuario `db:"-" json:"usuario"` // PK compuesta
	CantUsada int64   `db:"cant_usada" json:"cantUsada"`
}
