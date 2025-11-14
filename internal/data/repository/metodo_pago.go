package repository

import (
	"database/sql"

	"github.com/Nexivent/nexivent-backend/internal/data/model"
	util "github.com/Nexivent/nexivent-backend/internal/data/model/util"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MetodoDePago struct {
	DB *gorm.DB
}

func (r *MetodoDePago) CrearMetodoDePago(m *model.MetodoDePago) error {
	if m == nil {
		return gorm.ErrInvalidData
	}
	if err := r.DB.Create(m).Error; err != nil {
		return err
	}
	return nil
}

// Actualización **parcial** (recomendado): solo cambia los campos no-nil y retorna el registro.
func (r *MetodoDePago) ActualizarMetodoDePagoParcial(
	id int64,
	tipo *string,
	estado *int16,
) (*model.MetodoDePago, error) {

	updateFields := map[string]any{}
	if tipo != nil {
		updateFields["tipo"] = *tipo
	}
	if estado != nil {
		updateFields["estado"] = *estado
	}

	var m model.MetodoDePago

	// Sin cambios: retorna el registro actual
	if len(updateFields) == 0 {
		if err := r.DB.First(&m, "metodo_de_pago_id = ?", id).Error; err != nil {
			return nil, err
		}
		return &m, nil
	}

	res := r.DB.
		Model(&m).
		Clauses(clause.Returning{}).
		Where("metodo_de_pago_id = ?", id).
		Updates(updateFields)

	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &m, nil
}

// Desactivar (estado=0)
func (r *MetodoDePago) DesactivarMetodoDePago(id int64) error {
	res := r.DB.
		Table("metodo_de_pago").
		Where("metodo_de_pago_id = ? AND estado = 1", id).
		Update("estado", util.Inactivo)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *MetodoDePago) ObtenerMetodoDePagoBasicoPorID(id int64) (*model.MetodoDePago, error) {
	var m model.MetodoDePago
	if err := r.DB.First(&m, "metodo_de_pago_id = ?", id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *MetodoDePago) ListarMetodosActivos() ([]model.MetodoDePago, error) {
	var list []model.MetodoDePago
	if err := r.DB.
		Where("estado = 1").
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

//
// =====================================
//  Verificaciones para /payments/intent
// =====================================
//

// true si existe y está activo (estado=1)
func (r *MetodoDePago) VerificarMetodoDePagoActivo(id int64) (bool, error) {
	var count int64
	res := r.DB.
		Table("metodo_de_pago").
		Where("metodo_de_pago_id = ? AND estado = 1", id).
		Count(&count)
	if res.Error != nil {
		return false, res.Error
	}
	return count == 1, nil
}

// Obtiene 'tipo' (p.ej., 'Tarjeta' | 'Yape') para responder el campo "metodoPago"
func (r *MetodoDePago) ObtenerTipoDeMetodoPago(id int64) (string, error) {
	var tipo string
	err := r.DB.
		Table("metodo_de_pago").
		Select("tipo").
		Where("metodo_de_pago_id = ?", id).
		Row().
		Scan(&tipo)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", gorm.ErrRecordNotFound
		}
		return "", err
	}
	return tipo, nil
}

// true si la orden existe y no está CANCELADA (el BO puede endurecer: TEMPORAL vigente)
func (r *MetodoDePago) VerificarOrdenPermitePago(orderID int64) (bool, error) {
	var estInt int16
	err := r.DB.
		Table("orden_de_compra").
		Select("estado_de_orden").
		Where("orden_de_compra_id = ?", orderID).
		Row().
		Scan(&estInt)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, gorm.ErrRecordNotFound
		}
		return false, err
	}
	estado := util.EstadoOrden(estInt)
	// Bloqueamos solo CANCELADA aquí; el BO decide si exige TEMPORAL+vigente antes del intent.
	if estado == util.OrdenCancelada {
		return false, nil
	}
	return true, nil
}
