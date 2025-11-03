package domain

type MetodoDePago struct {
	ID     int64          `db:"metodo_de_pago_id" json:"id"`
	Tipo   TipoMetodoPago `db:"tipo" json:"tipo"`
	Activo int16          `db:"activo" json:"activo"`
}
