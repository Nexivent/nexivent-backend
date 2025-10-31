package data

type MetodoPago struct {
	IDMetodoPago int             `db:"id_metodo_pago" json:"idMetodoPago"`
	Tipo         TipoMetodoPago  `db:"tipo" json:"tipo"`
	Monto        float64         `db:"monto" json:"monto"`
	Estado       EstadoMetodoPago `db:"estado" json:"estado"`
}
