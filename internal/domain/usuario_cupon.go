package domain

type UsuarioCupon struct {
	Cupon     Cupon   `db:"-"          json:"cupon"`   // FK -> cupon (siempre)
	Usuario   Usuario `db:"-"          json:"usuario"` // FK -> usuario (siempre)
	CantUsada int64   `db:"cant_usada" json:"cantUsada"`
}
