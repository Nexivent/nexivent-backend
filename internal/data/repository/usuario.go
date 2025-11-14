package repository

import (
	"errors"
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Usuario struct {
	DB *gorm.DB
}

func NewUsuariosController(
	postgresqlDB *gorm.DB,
) *Usuario {
	return &Usuario{
		DB: postgresqlDB,
	}
}

func (c *Usuario) ObtenerUsuarios() ([]*model.Usuario, error) {
	usuarios := []*model.Usuario{}

	result := c.DB.Find(&usuarios)
	// Si requieres relaciones:
	// Preload("Comentarios").
	// Preload("Ordenes").
	// Preload("RolesAsignados").
	// Preload("Cupones").
	// Posiblemente se cree otras funciones para las relaciones más complejas.

	if result.Error != nil {
		return nil, result.Error
	}

	return usuarios, nil
}

func (c *Usuario) CrearUsuario(usuario *model.Usuario) error {
	if usuario == nil {
		return errors.New("CrearUsuario: usuario es nil")
	}

	result := c.DB.Create(usuario)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
func (u *Usuario) ActualizarUsuario(
	id int64,
	nombre *string,
	tipoDocumento *string,
	numDocumento *string,
	correo *string,
	contrasenha *string,
	telefono *string,
	estadoDeCuenta *int16,
	codigoVerificacion *string,
	fechaExpiracionCodigo *time.Time,
	estado *int16,
	updatedBy int64,
) (*model.Usuario, error) {

	updateFields := map[string]any{
		"usuario_modificacion": updatedBy,
		"fecha_modificacion":   time.Now(),
	}

	// Solo agregamos lo que llega no-nil (forzando incluso valores cero/vacíos)
	if nombre != nil {
		updateFields["nombre"] = *nombre
	}
	if tipoDocumento != nil {
		updateFields["tipo_documento"] = *tipoDocumento
	}
	if numDocumento != nil {
		updateFields["num_documento"] = *numDocumento
	}
	if correo != nil {
		updateFields["correo"] = *correo
	}
	if contrasenha != nil {
		updateFields["contrasenha"] = *contrasenha
	}
	if telefono != nil {
		updateFields["telefono"] = *telefono
	}
	if estadoDeCuenta != nil {
		updateFields["estado_de_cuenta"] = *estadoDeCuenta
	}
	if codigoVerificacion != nil {
		updateFields["codigo_verificacion"] = *codigoVerificacion
	}
	if fechaExpiracionCodigo != nil {
		updateFields["fecha_expiracion_codigo"] = *fechaExpiracionCodigo
	}
	if estado != nil {
		updateFields["estado"] = *estado
	}

	const baseAuditFields = 2 // usuario_modificacion + fecha_modificacion
	var user model.Usuario
	if len(updateFields) == baseAuditFields {
		if err := u.DB.First(&user, "usuario_id = ?", id).Error; err != nil {
			return nil, err
		}
		return &user, nil
	}

	result := u.DB.Model(&user).
		Clauses(clause.Returning{}).
		Where("usuario_id = ?", id).
		Updates(updateFields)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &user, nil
}

func (u *Usuario) DesactivarUsuario(id int64, updatedBy int64) error {
	result := u.DB.
		Model(&model.Usuario{}).
		Where("usuario_id = ? AND estado = 1", id).
		Updates(map[string]any{
			"estado":               int16(0),
			"usuario_modificacion": updatedBy,  //Quien modifico el registro
			"fecha_modificacion":   time.Now(), //Cuando se modifico el registro
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		// No existía o ya estaba en estado=0
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (u *Usuario) ObtenerUsuarioPorCorreo(correo string) (*model.Usuario, error) {
	var user model.Usuario

	result := u.DB.
		Where("correo = ?", correo).
		First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// func (u *Usuario) ExisteUsuarioPorCorreo(correo string) (bool, error) {
// 	var count int64

// 	result := u.DB.
// 		Model(&model.Usuario{}).
// 		Where("correo = ?", correo).
// 		Count(&count)
// 	if result.Error != nil {
// 		return false, result.Error
// 	}
// 	return count > 0, nil
// }

func (u *Usuario) ObtenerUsuarioBasicoPorID(id uint64) (*model.Usuario, error) {
	var user model.Usuario
	result := u.DB.
		First(&user, "usuario_id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

//Funcion para obtener los usuarios que tienen un rol específico por id de rol

func (c *Usuario) ObtenerUsuariosPorRolID(rolID uint64) ([]*model.Usuario, error) {
	usuarios := []*model.Usuario{}

	result := c.DB.
		Select("usuario.*").
		Joins("JOIN rol_usuario ru ON ru.usuario_id = usuario.usuario_id AND ru.estado = 1").
		Where("ru.rol_id = ?", rolID).
		Distinct().
		Find(&usuarios)

	if result.Error != nil {
		return nil, result.Error
	}
	return usuarios, nil
}
