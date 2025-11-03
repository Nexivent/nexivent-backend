package domain

import "github.com/shopspring/decimal"

type Cupon struct {
	ID            int64           `db:"cupon_id"        json:"cuponId"`
	Descripcion   string          `db:"descripcion"     json:"descripcion"`
	Tipo          string          `db:"tipo"            json:"tipo"` // (varchar en DDL)
	Valor         decimal.Decimal `db:"valor"           json:"valor"`
	EstadoCupon   int16           `db:"estado_cupon"    json:"estadoCupon"`
	Codigo        string          `db:"codigo"          json:"codigo"`
	UsoPorUsuario int64           `db:"uso_por_usuario" json:"usoPorUsuario"`
	UsoRealizados int64           `db:"uso_realizados"  json:"usoRealizados"`
}
