// evento.go
package model

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/util"
	"github.com/Nexivent/nexivent-backend/internal/validator"
)

type Evento struct {
	ID                  uint64            `gorm:"column:evento_id" json:"id"`
	OrganizadorID       uint64            `gorm:"column:organizador_id" json:"organizadorId"`
	CategoriaID         uint64            `gorm:"column:categoria_id" json:"categoriaId"`
	Titulo              string            `gorm:"column:titulo" json:"titulo"`
	Descripcion         string            `gorm:"column:descripcion" json:"descripcion"`
	DescripcionArtista  string            `gorm:"column:descripcion_artista" json:"descripcionArtista"`
	Lugar               string            `gorm:"column:lugar" json:"lugar"`
	EventoEstado        util.EstadoEvento `gorm:"column:evento_estado" json:"eventoEstado"`
	CantMeGusta         int               `gorm:"column:cant_me_gusta" json:"cantMeGusta"`
	CantNoInteresa      int               `gorm:"column:cant_no_interesa" json:"cantNoInteresa"`
	CantVendidoTotal    int               `gorm:"column:cant_vendido_total" json:"-"`
	TotalRecaudado      float64           `gorm:"column:total_recaudado" json:"-"`
	Estado              util.Estado       `gorm:"column:estado" json:"-"`
	UsuarioCreacion     *uint64           `gorm:"column:usuario_creacion" json:"-"`
	FechaCreacion       time.Time         `gorm:"column:fecha_creacion" json:"-"`
	UsuarioModificacion *uint64           `gorm:"column:usuario_modificacion" json:"-"`
	FechaModificacion   *time.Time        `gorm:"column:fecha_modificacion" json:"-"`
	ImagenDescripcion   string            `gorm:"column:imagen_descripcion" json:"imagenDescripcion"`
	ImagenPortada       string            `gorm:"column:imagen_portada" json:"imagenPortada"`
	VideoPresentacion   string            `gorm:"column:video_presentacion" json:"videoPresentacion"`
	ImagenEscenario     string            `gorm:"column:imagen_escenario" json:"imagenEscenario"`

	Comentarios []Comentario      `json:"comentarios"`
	Sectores    []Sector          `json:"sectores"`
	TiposTicket []TipoDeTicket    `json:"tiposDeTicket"`
	Perfiles    []PerfilDePersona `json:"perfiles"`
	Fechas      []EventoFecha     `json:"fechas"`
}

func (Evento) TableName() string { return "evento" }

func ValidateEvento(v *validator.Validator, evento *Evento) {
	// Validar Titulo
	v.Check(evento.Titulo != "", "titulo", "el título es obligatorio")
	v.Check(len(evento.Titulo) <= 80, "titulo", "el título no debe exceder 80 caracteres")

	// Validar Descripcion
	v.Check(evento.Descripcion != "", "descripcion", "la descripción es obligatoria")
	v.Check(len(evento.Descripcion) <= 5000, "descripcion", "la descripción no debe exceder 5000 caracteres")

	// Validar DescripcionArtista
	v.Check(evento.DescripcionArtista != "", "descripcion_artista", "la descripción del artista es obligatoria")
	v.Check(len(evento.DescripcionArtista) <= 5000, "descripcion_artista", "la descripción del artista no debe exceder 5000 caracteres")

	// Validar Lugar
	v.Check(evento.Lugar != "", "lugar", "el lugar es obligatorio")
	v.Check(len(evento.Lugar) <= 80, "lugar", "el lugar no debe exceder 80 caracteres")

	// Validar IDs
	v.Check(evento.OrganizadorID != 0, "organizador_id", "el ID del organizador es obligatorio")
	v.Check(evento.CategoriaID != 0, "categoria_id", "el ID de categoría es obligatorio")

	// Validar contadores (no pueden ser negativos)
	v.Check(evento.CantMeGusta >= 0, "cant_me_gusta", "la cantidad de me gusta no puede ser negativa")
	v.Check(evento.CantNoInteresa >= 0, "cant_no_interesa", "la cantidad de no me interesa no puede ser negativa")
	v.Check(evento.CantVendidoTotal >= 0, "cant_vendido_total", "la cantidad vendida no puede ser negativa")
	v.Check(evento.TotalRecaudado >= 0, "total_recaudado", "el total recaudado no puede ser negativo")

	// Validar fechas (obligatorio)
	v.Check(evento.Fechas != nil, "fechas", "las fechas son obligatorias")
	v.Check(len(evento.Fechas) > 0, "fechas", "debe haber al menos una fecha para el evento")

	// Validar URLs de imágenes y video (si están presentes, deben tener longitud razonable)
	if evento.ImagenDescripcion != "" {
		v.Check(len(evento.ImagenDescripcion) <= 255, "imagen_descripcion", "la URL de la imagen de descripción no debe exceder 255 caracteres")
	}
	if evento.ImagenPortada != "" {
		v.Check(len(evento.ImagenPortada) <= 255, "imagen_portada", "la URL de la imagen de portada no debe exceder 255 caracteres")
	}
	if evento.VideoPresentacion != "" {
		v.Check(len(evento.VideoPresentacion) <= 255, "video_presentacion", "la URL del video de presentación no debe exceder 255 caracteres")
	}
	if evento.ImagenEscenario != "" {
		v.Check(len(evento.ImagenEscenario) <= 255, "imagen_escenario", "la URL de la imagen del escenario no debe exceder 255 caracteres")
	}
}
