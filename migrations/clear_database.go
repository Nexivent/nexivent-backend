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

	// 🔹 Crear un evento de prueba
	evento := model.Evento{
		Titulo:          "Concierto de Rock",
		OrganizadorID:   1,
		CategoriaID:     1,
		FechaCreacion:   time.Now(),
		UsuarioCreacion: 1,
		Estado:          1,

		Sectores: []model.Sector{
			{SectorTipo: "VIP", TotalEntradas: 1000, Estado: 1, UsuarioCreacion: 1, FechaCreacion: time.Now()},
			{SectorTipo: "General", TotalEntradas: 2000, Estado: 1, UsuarioCreacion: 1, FechaCreacion: time.Now()},
		},

		TiposTicket: []model.TipoDeTicket{
			{Nombre: "Entrada General", FechaIni: time.Now(), FechaFin: time.Now().AddDate(0, 0, 30), Estado: 1, UsuarioCreacion: 1, FechaCreacion: time.Now()},
		},

		Perfiles: []model.PerfilDePersona{
			{Nombre: "Estudiante", Estado: 1, UsuarioCreacion: 1, FechaCreacion: time.Now()},
			{Nombre: "Adulto", Estado: 1, UsuarioCreacion: 1, FechaCreacion: time.Now()},
		},
	}

	entidad.Evento.CrearEvento(&evento)
	//falta crear tabla de fechas y + para que funcione
	eventos, err := entidad.Evento.ObtenerEventosDisponiblesSinFiltros(10, 1)
	if err != nil {
		log.Fatalf("❌ Error al obtener eventos: %v", err)
	}

	fmt.Println("✅ Eventos disponibles:")

	if eventos == nil || len(eventos.Eventos) == 0 {
		fmt.Println("No se encontraron eventos disponibles.")
		return
	}

	for _, e := range eventos.Eventos {
		fmt.Printf("ID: %d - Titulo: %s - Estado: %v\n", e.ID, e.Titulo, e.Estado)
	}
}
