package main

import (
	"fmt"
	"log"
	"time"

	config "github.com/Loui27/nexivent-backend/internal/config"
	"github.com/Loui27/nexivent-backend/internal/dao/model"
	"github.com/Loui27/nexivent-backend/internal/dao/repository"
	"github.com/Loui27/nexivent-backend/logging"
	setupDB "github.com/Loui27/nexivent-backend/utils"
)

func main() {
	testLogger := logging.NewLoggerMock()
	envSettings := config.NuevoConfigEnv(testLogger)

	entidad, nexiventPsqlDB := repository.NewNexiventPsqlEntidades(testLogger, envSettings)

	setupDB.ClearPostgresqlDatabase(testLogger, nexiventPsqlDB, envSettings, nil)
	usuario := &model.Usuario{
		Nombre:         "Juan Pérez",
		TipoDocumento:  "DNI",
		NumDocumento:   "12345678",
		Correo:         "juanperez@example.com",
		Contrasenha:    "12345",
		EstadoDeCuenta: 1,
		FechaCreacion:  time.Now(),
		Estado:         1,
	}
	entidad.Usuario.CrearUsuario(usuario)

	categoria := &model.Categoria{
		Nombre:      "Conciertos",
		Descripcion: "Eventos musicales en vivo",
		Estado:      1,
	}

	entidad.Categoria.CrearCategoria(categoria)

	// 🔹 Crear un evento de prueba con relaciones mínimas para que GET /evento funcione
	createdBy := int64(1)
	now := time.Now()
	evento := model.Evento{
		Titulo:          "Concierto de POP",
		OrganizadorID:   usuario.ID,
		CategoriaID:     categoria.ID,
		FechaCreacion:   now,
		UsuarioCreacion: &usuario.ID,
		Estado:          1,
		EventoEstado:    1,
		Lugar:           "Estadio Nacional",
		Descripcion:     "Evento semilla para pruebas",

		Sectores: []model.Sector{
			{SectorTipo: "VIP", TotalEntradas: 1000, Estado: 1, UsuarioCreacion: &usuario.ID, FechaCreacion: now},
			{SectorTipo: "General", TotalEntradas: 2000, Estado: 1, UsuarioCreacion: &usuario.ID, FechaCreacion: now},
		},

		TiposTicket: []model.TipoDeTicket{
			{Nombre: "Entrada General", FechaIni: now, FechaFin: now.AddDate(0, 0, 30), Estado: 1, UsuarioCreacion: &createdBy, FechaCreacion: now},
		},

		Perfiles: []model.PerfilDePersona{
			{Nombre: "Estudiante", Estado: 1, UsuarioCreacion: &createdBy, FechaCreacion: now},
			{Nombre: "Adulto", Estado: 1, UsuarioCreacion: &createdBy, FechaCreacion: now},
		},
	}

	if err := entidad.Evento.CrearEvento(&evento); err != nil {
		log.Fatalf("❌ Error creando evento seed: %v", err)
	}

	// Crear una fecha futura y vincularla con el evento para que el listado encuentre registros
	fechaEvento := now.AddDate(0, 0, 7)
	fecha := &model.Fecha{FechaEvento: fechaEvento}
	if err := nexiventPsqlDB.Create(fecha).Error; err != nil {
		log.Fatalf("❌ Error creando fecha seed: %v", err)
	}

	horaInicio := time.Date(fechaEvento.Year(), fechaEvento.Month(), fechaEvento.Day(), 20, 0, 0, 0, fechaEvento.Location())
	eventoFecha := &model.EventoFecha{
		EventoID:        evento.ID,
		FechaID:         fecha.ID,
		HoraInicio:      horaInicio,
		Estado:          1,
		UsuarioCreacion: &createdBy,
		FechaCreacion:   now,
	}
	if err := nexiventPsqlDB.Create(eventoFecha).Error; err != nil {
		log.Fatalf("❌ Error creando evento_fecha seed: %v", err)
	}

	eventos, err := entidad.Evento.ObtenerEventosDisponiblesSinFiltros()
	if err != nil {
		log.Fatalf("❌ Error al obtener eventos: %v", err)
	}

	fmt.Println("✅ Eventos disponibles:")

	if eventos == nil {
		fmt.Println("No se encontraron eventos disponibles.")
		return
	}

	// for _, e := range eventos.Eventos {
	// 	fmt.Printf("ID: %d - Titulo: %s - Estado: %v\n", e.ID, e.Titulo, e.Estado)
	// }
}
