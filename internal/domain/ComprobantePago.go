package domain

import (
	"database/sql"
	"time"
)

type ComprobantePago struct {
	IDComprobante   int             `db:"id_comprobante" json:"idComprobante"`
	Orden           OrdenDeCompra   `db:"-" json:"orden"`
	Tipo            TipoComprobante `db:"tipo" json:"tipo"`
	Numero          string          `db:"numero" json:"numero"`
	FechaEmision    time.Time       `db:"fecha_emision" json:"fechaEmision"`
	Total           float64         `db:"total" json:"total"`
	NombreCliente   sql.NullString  `db:"nombre_cliente" json:"nombreCliente,omitempty"`
	RucCliente      sql.NullString  `db:"ruc_cliente" json:"rucCliente,omitempty"`
	DireccionFiscal sql.NullString  `db:"direccion_fiscal" json:"direccionFiscal,omitempty"`
}
