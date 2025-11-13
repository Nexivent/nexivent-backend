package model

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/util"
	"github.com/Nexivent/nexivent-backend/internal/validator"
)

type Comentario struct {
	ID            uint64      `gorm:"column:comentario_id;primaryKey" json:"id"`
	UsuarioID     uint64      `gorm:"column:usuario_id" json:"usuarioId"`
	EventoID      uint64      `gorm:"column:evento_id" json:"eventoId"`
	Descripcion   string      `gorm:"column:descripcion" json:"descripcion"`
	FechaCreacion time.Time   `gorm:"column:fecha_creacion" json:"fechaCreacion"`
	Estado        util.Estado `gorm:"column:estado" json:"-"`
}

func (Comentario) TableName() string { return "comentario" }

func ValidateComentario(v *validator.Validator, comentario *Comentario) {
	// Validar Descripcion
	v.Check(comentario.Descripcion != "", "descripcion", "la descripción es obligatoria")
	v.Check(len(comentario.Descripcion) <= 1000, "descripcion", "la descripción no debe exceder 1000 caracteres")

	// Validar IDs
	v.Check(comentario.UsuarioID != 0, "usuarioId", "el ID del usuario es obligatorio")
	v.Check(comentario.EventoID != 0, "eventoId", "el ID del evento es obligatorio")
}
