package model

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/util"
	"github.com/Nexivent/nexivent-backend/internal/validator"
)

type ComprobanteDePago struct {
	ID                uint64               `gorm:"column:comprobante_de_pago_id;primaryKey" json:"id"`
	OrdenDeCompraID   uint64               `gorm:"column:orden_de_compra_id" json:"ordenDeCompraId"`
	TipoDeComprobante util.TipoComprobante `gorm:"column:tipo_de_comprobante" json:"tipoDeComprobante"`
	Numero            string               `gorm:"column:numero" json:"numero"`
	FechaEmision      time.Time            `gorm:"column:fecha_emision" json:"fechaEmision"`
	RUC               *string              `gorm:"column:ruc" json:"ruc,omitempty"`
	DireccionFiscal   *string              `gorm:"column:direccion_fiscal" json:"direccionFiscal,omitempty"`
}

func (ComprobanteDePago) TableName() string { return "comprobante_de_pago" }

func ValidateComprobanteDePago(v *validator.Validator, comprobante *ComprobanteDePago) {
	// Validar OrdenDeCompraID
	v.Check(comprobante.OrdenDeCompraID != 0, "ordenDeCompraId", "el ID de la orden de compra es obligatorio")

	// Validar Numero
	v.Check(comprobante.Numero != "", "numero", "el número es obligatorio")
	v.Check(len(comprobante.Numero) <= 50, "numero", "el número no debe exceder 50 caracteres")

	// Validar FechaEmision
	v.Check(!comprobante.FechaEmision.IsZero(), "fechaEmision", "la fecha de emisión es obligatoria")

	// Validar RUC y DireccionFiscal si el tipo es factura (1)
	if comprobante.TipoDeComprobante == 1 {
		v.Check(comprobante.RUC != nil && *comprobante.RUC != "", "ruc", "el RUC es obligatorio para facturas")
		v.Check(comprobante.DireccionFiscal != nil && *comprobante.DireccionFiscal != "", "direccionFiscal", "la dirección fiscal es obligatoria para facturas")
	}

	// Validar longitud de RUC si está presente
	if comprobante.RUC != nil && *comprobante.RUC != "" {
		v.Check(len(*comprobante.RUC) == 11, "ruc", "el RUC debe tener 11 caracteres")
	}
}
