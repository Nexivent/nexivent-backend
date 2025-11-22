package main

import (
	"fmt"
	"log"
	"math"
	"time"

	config "github.com/Nexivent/nexivent-backend/internal/config"
	"github.com/Nexivent/nexivent-backend/internal/dao/model"
	util "github.com/Nexivent/nexivent-backend/internal/dao/model/util"
	"github.com/Nexivent/nexivent-backend/internal/dao/repository"
	"github.com/Nexivent/nexivent-backend/logging"
	"gorm.io/gorm"
)

func main() {
	//logger := logging.NewLoggerMock()
	logger := logging.NewLogger("Migrations", "Version 1.0", logging.FormatText, 4)
	envSettings := config.NuevoConfigEnv(logger)

	entidad, nexiventPsqlDB := repository.NewNexiventPsqlEntidades(logger, envSettings)

	//if isLocalhost(envSettings.PostgresHost) {
	//	setupDB.ClearPostgresqlDatabase(logger, nexiventPsqlDB, envSettings, nil)
	//} else {
	//	log.Printf("ℹ️ Saltando borrado completo porque el host no es local (%s)\n", envSettings.PostgresHost)
	//}
	//
	//if err := seedDatabase(logger, nexiventPsqlDB, entidad); err != nil {
	//	log.Fatalf("❌ Error sembrando datos: %v", err)
	//}
	//
	//fmt.Println("✅ Base de datos inicializada con datos semilla.")
	// 3. SIEMPRE ejecutar seeds independientemente de si es local o no
	logger.Info("🌱 Iniciando proceso de seeds...")

	if err := seedDatabase(logger, nexiventPsqlDB, entidad); err != nil {
		log.Fatalf("❌ Error al ejecutar seeds: %v", err)
	}

	logger.Info("✅ Base de datos inicializada con datos semilla.")
}

func seedDatabase(
	logger logging.Logger,
	db *gorm.DB,
	entidad *repository.NexiventPsqlEntidades,
) error {
	//var eventosExistentes int64
	//if err := db.Model(&model.Evento{}).Count(&eventosExistentes).Error; err != nil {
	//	return fmt.Errorf("no se pudo contar eventos existentes: %w", err)
	//}
	//if eventosExistentes > 0 {
	//	logger.Infof("Seed omitido: la BD ya tiene %d eventos", eventosExistentes)
	//	return nil
	//}
	// Crear roles primero usando el repositorio
	logger.Info("Iniciando seed de roles...")
	if err := seedRoles(logger, entidad); err != nil {
		logger.Errorf("Error al crear roles: %v", err)
		return fmt.Errorf("error en seedRoles: %w", err)
	}
	logger.Info("Roles creados exitosamente")

	usuarios, err := seedUsuarios(entidad)
	if err != nil {
		return err
	}
	categorias, err := seedCategorias(entidad)
	if err != nil {
		return err
	}
	if err := seedMetodosPago(db); err != nil {
		return err
	}
	eventos, err := seedEventos(logger, db, entidad, usuarios, categorias)
	if err != nil {
		return err
	}
	if err := seedTicketsComprados(logger, db, entidad, eventos, usuarios); err != nil {
		return err
	}

	return nil
}

func seedRoles(logger logging.Logger, entidad *repository.NexiventPsqlEntidades) error {
	now := time.Now()
	roles := []struct {
		nombre string
		rol    *model.Rol
	}{
		{
			nombre: "ASISTENTE",
			rol: &model.Rol{
				Nombre:        "ASISTENTE",
				FechaCreacion: now,
			},
		},
		{
			nombre: "ADMINISTRADOR",
			rol: &model.Rol{
				Nombre:        "ADMINISTRADOR",
				FechaCreacion: now,
			},
		},
		{
			nombre: "ORGANIZADOR",
			rol: &model.Rol{
				Nombre:        "ORGANIZADOR",
				FechaCreacion: now,
			},
		},
	}

	for _, r := range roles {
		logger.Infof("Verificando rol: %s", r.nombre)

		// Verificar si el rol ya existe
		existente, err := entidad.Roles.ObtenerRolPorNombre(r.nombre)

		if err != nil && err.Error() != "record not found" {
			logger.Errorf("Error al buscar rol %s: %v", r.nombre, err)
			return fmt.Errorf("error al verificar rol %s: %w", r.nombre, err)
		}

		if existente != nil {
			logger.Infof("✅ Rol %s ya existe con ID: %d", r.nombre, existente.ID)
			continue
		}

		// Crear el rol usando el repositorio
		logger.Infof("Creando rol: %s", r.nombre)
		if err := entidad.Roles.CrearRol(r.rol); err != nil {
			logger.Errorf("Error al crear rol %s: %v", r.nombre, err)
			return fmt.Errorf("no se pudo crear rol %s: %w", r.nombre, err)
		}
		logger.Infof("✅ Rol %s creado exitosamente con ID: %d", r.nombre, r.rol.ID)
	}

	return nil
}

func seedUsuarios(entidad *repository.NexiventPsqlEntidades) ([]model.Usuario, error) {
	usuarios := []model.Usuario{
		{
			Nombre:         "Ana Rojas",
			TipoDocumento:  "DNI",
			NumDocumento:   "45812345",
			Correo:         "ana.rojas@nexivent.com",
			Contrasenha:    "admin123",
			EstadoDeCuenta: 1,
			Estado:         1,
		},
		{
			Nombre:         "Luis Salazar",
			TipoDocumento:  "DNI",
			NumDocumento:   "78451236",
			Correo:         "luis.salazar@nexivent.com",
			Contrasenha:    "organizer123",
			EstadoDeCuenta: 1,
			Estado:         1,
		},
		{
			Nombre:         "María Castillo",
			TipoDocumento:  "CE",
			NumDocumento:   "X8945123",
			Correo:         "maria.castillo@nexivent.com",
			Contrasenha:    "ventas123",
			EstadoDeCuenta: 1,
			Estado:         1,
		},
		{
			Nombre:         "Carlos Ruiz",
			TipoDocumento:  "DNI",
			NumDocumento:   "70112584",
			Correo:         "carlos.ruiz@nexivent.com",
			Contrasenha:    "cliente123",
			EstadoDeCuenta: 1,
			Estado:         1,
		},
	}

	for i := range usuarios {
		if err := entidad.Usuario.CrearUsuario(&usuarios[i]); err != nil {
			return nil, fmt.Errorf("no se pudo crear usuario %s: %w", usuarios[i].Nombre, err)
		}
	}

	return usuarios, nil
}

func seedCategorias(entidad *repository.NexiventPsqlEntidades) ([]model.Categoria, error) {
	categorias := []model.Categoria{
		{Nombre: "Conciertos", Descripcion: "Música en vivo y festivales", Estado: 1},
		{Nombre: "Deportes", Descripcion: "Eventos deportivos y competencias", Estado: 1},
		{Nombre: "Teatro", Descripcion: "Obras y shows escénicos", Estado: 1},
		{Nombre: "Tecnología", Descripcion: "Meetups, conferencias y hackathons", Estado: 1},
		{Nombre: "Negocios", Descripcion: "Cumbres, ferias y networking", Estado: 1},
		{Nombre: "Gastronomía", Descripcion: "Ferias y experiencias culinarias", Estado: 1},
	}

	for i := range categorias {
		if err := entidad.Categoria.CrearCategoria(&categorias[i]); err != nil {
			return nil, fmt.Errorf("no se pudo crear categoría %s: %w", categorias[i].Nombre, err)
		}
	}

	return categorias, nil
}

func seedMetodosPago(db *gorm.DB) error {
	metodos := []model.MetodoDePago{
		{Tipo: "Tarjeta", Estado: 1},
		{Tipo: "Yape", Estado: 1},
	}

	for _, metodo := range metodos {
		var existente model.MetodoDePago
		if err := db.Where("tipo = ?", metodo.Tipo).FirstOrCreate(&existente, metodo).Error; err != nil {
			return fmt.Errorf("no se pudo crear método de pago %s: %w", metodo.Tipo, err)
		}
	}

	return nil
}

func seedEventos(
	logger logging.Logger,
	db *gorm.DB,
	entidad *repository.NexiventPsqlEntidades,
	usuarios []model.Usuario,
	categorias []model.Categoria,
) ([]model.Evento, error) {
	now := time.Now()
	categoriaPorNombre := make(map[string]int64)
	for i := range categorias {
		categoriaPorNombre[categorias[i].Nombre] = categorias[i].ID
	}

	type sectorSeed struct {
		Nombre     string
		Capacidad  int
		PrecioBase float64
	}

	type ticketSeed struct {
		Nombre              string
		InicioDiasAntes     int
		FinDiasAntes        int
		MultiplicadorPrecio float64
	}

	type couponSeed struct {
		Codigo string
		Valor  float64
		Tipo   int16
	}

	type eventoSeed struct {
		Titulo            string
		Descripcion       string
		Lugar             string
		Categoria         string
		OrganizadorIdx    int
		DiasHastaFecha    int
		HoraInicio        int
		MinutoInicio      int
		Sectores          []sectorSeed
		Perfiles          []string
		Tickets           []ticketSeed
		ImagenDescripcion string
		ImagenPortada     string
		VideoPresentacion string
		ImagenEscenario   string
		Cupon             *couponSeed
	}

	eventos := []eventoSeed{
		{
			Titulo:            "Shakira - Estoy aqui World Tour",
			Descripcion:       "Tras agotar entradas en tiempo récord para sus conciertos del 15 y 16 de noviembre, la superestrella global Shakira anuncia una nueva y última fecha en Lima: El 18 de noviembre en el Estadio Nacional, como parte de su histórica gira mundial Las Mujeres Ya No Lloran World Tour",
			Lugar:             "Estado Nacional",
			Categoria:         "Conciertos",
			OrganizadorIdx:    2,
			DiasHastaFecha:    270,
			HoraInicio:        9,
			MinutoInicio:      0,
			ImagenDescripcion: "https://cdn.teleticket.com.pe/images/eventos/csi006_rs.jpg",
			ImagenPortada:     "https://cdn.teleticket.com.pe/images/eventos/csi006_rs.jpg",
			VideoPresentacion: "https://www.youtube.com/watch?v=2Ndra-1Pwug",
			ImagenEscenario:   "https://cdn.teleticket.com.pe/especiales/shakira-estoy-aqui-2025/images/mapa.png",
			Sectores: []sectorSeed{
				{Nombre: "VIP", Capacidad: 500, PrecioBase: 280},
				{Nombre: "Platea", Capacidad: 1200, PrecioBase: 180},
				{Nombre: "General", Capacidad: 2500, PrecioBase: 110},
			},
			Perfiles: []string{"Adulto", "Estudiante", "Fan Club"},
			Tickets: []ticketSeed{
				{Nombre: "Preventa", InicioDiasAntes: 45, FinDiasAntes: 10, MultiplicadorPrecio: 0.92},
				{Nombre: "General", InicioDiasAntes: 25, FinDiasAntes: 0, MultiplicadorPrecio: 1.0},
			},
			Cupon: &couponSeed{Codigo: "BICHOTA", Valor: 25, Tipo: 1},
		},
		{
			Titulo:         "Concierto Latin Pop",
			Descripcion:    "Noche de pop latino con artistas invitados y DJ set.",
			Lugar:          "Arena San Miguel",
			Categoria:      "Conciertos",
			OrganizadorIdx: 1,
			DiasHastaFecha: 15,
			HoraInicio:     20,
			MinutoInicio:   30,
			Sectores: []sectorSeed{
				{Nombre: "VIP", Capacidad: 500, PrecioBase: 280},
				{Nombre: "Platea", Capacidad: 1200, PrecioBase: 180},
				{Nombre: "General", Capacidad: 2500, PrecioBase: 110},
			},
			Perfiles: []string{"Adulto", "Estudiante", "Fan Club"},
			Tickets: []ticketSeed{
				{Nombre: "Preventa", InicioDiasAntes: 45, FinDiasAntes: 10, MultiplicadorPrecio: 0.92},
				{Nombre: "General", InicioDiasAntes: 25, FinDiasAntes: 0, MultiplicadorPrecio: 1.0},
			},
			Cupon: &couponSeed{Codigo: "POP5", Valor: 5, Tipo: 1},
		},
		{
			Titulo:            "Las cazadoras KPOP",
			Descripcion:       "¡El espectáculo más esperado llega al Canout! Las Cazadoras del Kpop el musical es una historia original inspirada en los musicales Kpop más vistos de los últimos tiempos.",
			Lugar:             "Teatro Canout",
			Categoria:         "Teatro",
			OrganizadorIdx:    2,
			DiasHastaFecha:    23,
			HoraInicio:        9,
			MinutoInicio:      0,
			ImagenDescripcion: "https://cdn.teleticket.com.pe/especiales/lascazadoraskpop/images/disco-imagen-2.jpg",
			ImagenPortada:     "https://cdn.teleticket.com.pe/especiales/lascazadoraskpop/images/disco-imagen-2.jpg",
			VideoPresentacion: "https://nexivent-multimedia.s3.us-east-2.amazonaws.com/kpopdemon.mp4",
			ImagenEscenario:   "https://cdn.teleticket.com.pe/especiales/lascazadoraskpop/images/mapa.png",
			Sectores: []sectorSeed{
				{Nombre: "VIP", Capacidad: 500, PrecioBase: 280},
				{Nombre: "Platea", Capacidad: 1200, PrecioBase: 180},
				{Nombre: "General", Capacidad: 2500, PrecioBase: 110},
			},
			Perfiles: []string{"Adulto", "Estudiante", "Fan Club"},
			Tickets: []ticketSeed{
				{Nombre: "Preventa", InicioDiasAntes: 45, FinDiasAntes: 10, MultiplicadorPrecio: 0.92},
				{Nombre: "General", InicioDiasAntes: 25, FinDiasAntes: 0, MultiplicadorPrecio: 1.0},
			},
			Cupon: &couponSeed{Codigo: "BICHOTA", Valor: 25, Tipo: 1},
		},
		{
			Titulo:            "Linkin Park - From Zero World Tour",
			Descripcion:       "Linkin Park regresa a Perú con su nueva vocalista Emily Armstrong y Colin Brittain en batería. Disfruta de los clásicos como Numb, In The End y nuevo material del álbum From Zero.",
			Lugar:             "Estadio San Marcos, Lima",
			Categoria:         "Conciertos",
			OrganizadorIdx:    0,
			DiasHastaFecha:    200,
			HoraInicio:        20,
			MinutoInicio:      0,
			ImagenDescripcion: "https://cdn.getcrowder.com/images/d23c610b-dd42-4082-be88-d11d7e477838-tmbanner-mobile.jpg",
			ImagenPortada:     "https://cdn.getcrowder.com/images/d23c610b-dd42-4082-be88-d11d7e477838-tmbanner-mobile.jpg",
			VideoPresentacion: "https://nexivent-multimedia.s3.us-east-2.amazonaws.com/linkin+park.mp4",
			ImagenEscenario:   "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSH6bmEdwdMH0i3GtHliDu7YZ8DJHAg8ZrogA&s",
			Sectores: []sectorSeed{
				{Nombre: "Campo A", Capacidad: 5000, PrecioBase: 750},
				{Nombre: "Campo B", Capacidad: 10000, PrecioBase: 480},
				{Nombre: "Tribuna Norte", Capacidad: 6000, PrecioBase: 320},
				{Nombre: "Tribuna Sur", Capacidad: 6000, PrecioBase: 320},
			},
			Perfiles: []string{"Rock Fan", "Millennials", "Nostálgico"},
			Tickets: []ticketSeed{
				{Nombre: "Preventa BBVA", InicioDiasAntes: 90, FinDiasAntes: 85, MultiplicadorPrecio: 0.85},
				{Nombre: "Venta General", InicioDiasAntes: 85, FinDiasAntes: 1, MultiplicadorPrecio: 1.0},
			},
			Cupon: &couponSeed{Codigo: "LP15OFF", Valor: 15, Tipo: 1},
		},
		{
			Titulo:            "Geek Festival 2025 – Lima",
			Descripcion:       "Festival cultural puede incluir cómics, gaming, cultura pop y tecnología, con expositores, charlas y concursos.",
			Lugar:             "Parque de la Exposición, Lima",
			Categoria:         "Tecnología",
			OrganizadorIdx:    0,
			DiasHastaFecha:    63,
			HoraInicio:        20,
			MinutoInicio:      0,
			ImagenDescripcion: "https://manoalzada.pe/wp-content/uploads/2025/10/geek-fetsival-ofoicial.jpg",
			ImagenPortada:     "https://manoalzada.pe/wp-content/uploads/2025/10/geek-fetsival-ofoicial.jpg",
			VideoPresentacion: "https://nexivent-multimedia.s3.us-east-2.amazonaws.com/geek+festival.mp4",
			ImagenEscenario:   "https://cdn.teleticket.com.pe/especiales/lascazadoraskpop/images/mapa.png",
			Sectores: []sectorSeed{
				{Nombre: "VIP", Capacidad: 200, PrecioBase: 300},
				{Nombre: "General", Capacidad: 1000, PrecioBase: 150},
				{Nombre: "Estudiante", Capacidad: 500, PrecioBase: 100},
			},
			Perfiles: []string{"Fan", "Estudiante", "Profesional"},
			Tickets: []ticketSeed{
				{Nombre: "Preventa", InicioDiasAntes: 60, FinDiasAntes: 20, MultiplicadorPrecio: 0.9},
				{Nombre: "General", InicioDiasAntes: 19, FinDiasAntes: 1, MultiplicadorPrecio: 1.0},
				{Nombre: "Día evento", InicioDiasAntes: 0, FinDiasAntes: 0, MultiplicadorPrecio: 1.25},
			},
			Cupon: &couponSeed{Codigo: "GEEK15", Valor: 15, Tipo: 1},
		},
		{
			Titulo:            "The Weeknd en Lima 2025",
			Descripcion:       "El cantante canadiense The Weeknd regresa a Lima como parte de su gira mundial After Hours Til Dawn Tour.",
			Lugar:             "Estadio San Marcos, Lima",
			Categoria:         "Conciertos",
			OrganizadorIdx:    3,
			DiasHastaFecha:    25,
			HoraInicio:        17,
			MinutoInicio:      0,
			ImagenDescripcion: "https://i.ytimg.com/vi/zLsR8-iOd-E/maxresdefault.jpg",
			ImagenPortada:     "https://i.ytimg.com/vi/zLsR8-iOd-E/maxresdefault.jpg",
			VideoPresentacion: "https://nexivent-multimedia.s3.us-east-2.amazonaws.com/The+weelend.mp4",
			ImagenEscenario:   "https://cdn.teleticket.com.pe/especiales/the-weeknd-2023/images/mapa-v1.png",
			Sectores: []sectorSeed{
				{Nombre: "VIP", Capacidad: 300, PrecioBase: 210},
				{Nombre: "General", Capacidad: 3200, PrecioBase: 95},
			},
			Perfiles: []string{"Adulto", "Fan Club"},
			Tickets: []ticketSeed{
				{Nombre: "Preventa", InicioDiasAntes: 60, FinDiasAntes: 20, MultiplicadorPrecio: 0.9},
				{Nombre: "General", InicioDiasAntes: 35, FinDiasAntes: 0, MultiplicadorPrecio: 1.0},
			},
		},
		{
			Titulo:            "Aria Bela - World Tour 2025",
			Descripcion:       "La superestrella británica Aria Bela llega a Lima con su gira mundial Radical Optimism Tour, presentando sus más grandes éxitos y nuevo material de su tercer álbum.",
			Lugar:             "Estadio San Marcos, Lima",
			Categoria:         "Conciertos",
			OrganizadorIdx:    0,
			DiasHastaFecha:    53,
			HoraInicio:        20,
			MinutoInicio:      0,
			ImagenDescripcion: "https://imagendelgolfo.mx/img/2025/06/11/20250611_051250039_Ariatopxa_de_Aria_Belax_xCxmo_se_llaman_sus_nuevas_canciones_y_cuxndo_se_estrenanx.jpg",
			ImagenPortada:     "https://imagendelgolfo.mx/img/2025/06/11/20250611_051250039_Ariatopxa_de_Aria_Belax_xCxmo_se_llaman_sus_nuevas_canciones_y_cuxndo_se_estrenanx.jpg",
			VideoPresentacion: "https://nexivent-multimedia.s3.us-east-2.amazonaws.com/aria+bela.mp4",
			ImagenEscenario:   "https://cdn.teleticket.com.pe/especiales/lascazadoraskpop/images/mapa.png",
			Sectores: []sectorSeed{
				{Nombre: "Campo A", Capacidad: 5000, PrecioBase: 750},
				{Nombre: "Campo B", Capacidad: 10000, PrecioBase: 480},
				{Nombre: "Tribuna Norte", Capacidad: 6000, PrecioBase: 320},
				{Nombre: "Tribuna Sur", Capacidad: 6000, PrecioBase: 320},
			},
			Perfiles: []string{"Rock Fan", "Millennials", "Nostálgico"},
			Tickets: []ticketSeed{
				{Nombre: "Preventa BBVA", InicioDiasAntes: 90, FinDiasAntes: 85, MultiplicadorPrecio: 0.85},
				{Nombre: "Venta General", InicioDiasAntes: 85, FinDiasAntes: 1, MultiplicadorPrecio: 1.0},
			},
			Cupon: &couponSeed{Codigo: "LP15OFF", Valor: 15, Tipo: 1},
		},
		{
			Titulo:            "Grupo 5 en Concierto 2025",
			Descripcion:       "El Grupo 5 regresa a los escenarios con su gira 2025, presentando sus más grandes éxitos de cumbia y nuevas canciones.",
			Lugar:             "Circuito de Playas",
			Categoria:         "Conciertos",
			OrganizadorIdx:    3,
			DiasHastaFecha:    18,
			HoraInicio:        7,
			MinutoInicio:      30,
			ImagenDescripcion: "https://depor.com/resizer/v2/7I2UK6V4CBDZVD2PXLNIZH35KI.png?auth=49fc3c4f69ad73ef06573f33bae353f78e2ff12338f674f13c575ec7d81cace2&width=1600&height=900&quality=90&smart=true",
			ImagenPortada:     "https://depor.com/resizer/v2/7I2UK6V4CBDZVD2PXLNIZH35KI.png?auth=49fc3c4f69ad73ef06573f33bae353f78e2ff12338f674f13c575ec7d81cace2&width=1600&height=900&quality=90&smart=true",
			VideoPresentacion: "https://nexivent-multimedia.s3.us-east-2.amazonaws.com/grupo+5.mp4",
			ImagenEscenario:   "https://cdn.teleticket.com.pe/especiales/grupo5-lima-2024/images/mapa-2.png",
			Sectores: []sectorSeed{
				{Nombre: "VIP", Capacidad: 800, PrecioBase: 85},
				{Nombre: "General", Capacidad: 1200, PrecioBase: 60},
			},
			Perfiles: []string{"Adulto", "Estudiante"},
			Tickets: []ticketSeed{
				{Nombre: "General", InicioDiasAntes: 40, FinDiasAntes: 0, MultiplicadorPrecio: 1.0},
			},
			Cupon: &couponSeed{Codigo: "MOTORYMOTIVO", Valor: 7, Tipo: 1},
		},
	}

	perfilMultiplicador := map[string]float64{
		"Adulto":        1.0,
		"Niño":          0.7,
		"Estudiante":    0.8,
		"Profesional":   1.05,
		"Founder":       1.1,
		"Inversionista": 1.2,
		"Fan Club":      1.05,
		"Asistente":     1.0,
	}

	var eventosCreados []model.Evento

	for _, seed := range eventos {
		if seed.OrganizadorIdx >= len(usuarios) {
			return nil, fmt.Errorf("organizador fuera de rango para evento %s", seed.Titulo)
		}
		organizador := usuarios[seed.OrganizadorIdx]
		categoriaID := categoriaPorNombre[seed.Categoria]
		if categoriaID == 0 {
			return nil, fmt.Errorf("no se encontró categoría %s", seed.Categoria)
		}

		fechaEvento := now.AddDate(0, 0, seed.DiasHastaFecha)
		usuarioCreacion := organizador.ID
		evento := model.Evento{
			Titulo:            seed.Titulo,
			OrganizadorID:     organizador.ID,
			CategoriaID:       categoriaID,
			Descripcion:       seed.Descripcion,
			Lugar:             seed.Lugar,
			EventoEstado:      1,
			Estado:            1,
			UsuarioCreacion:   &usuarioCreacion,
			FechaCreacion:     now,
			ImagenDescripcion: seed.ImagenDescripcion,
			ImagenPortada:     seed.ImagenPortada,
			VideoPresentacion: seed.VideoPresentacion,
			ImagenEscenario:   seed.ImagenEscenario,
		}

		for _, sector := range seed.Sectores {
			evento.Sectores = append(evento.Sectores, model.Sector{
				SectorTipo:      sector.Nombre,
				TotalEntradas:   sector.Capacidad,
				Estado:          1,
				UsuarioCreacion: &usuarioCreacion,
				FechaCreacion:   now,
			})
		}

		for _, perfil := range seed.Perfiles {
			evento.Perfiles = append(evento.Perfiles, model.PerfilDePersona{
				Nombre:          perfil,
				Estado:          1,
				UsuarioCreacion: &usuarioCreacion,
				FechaCreacion:   now,
			})
		}

		for _, ticket := range seed.Tickets {
			evento.TiposTicket = append(evento.TiposTicket, model.TipoDeTicket{
				Nombre:          ticket.Nombre,
				FechaIni:        fechaEvento.AddDate(0, 0, -ticket.InicioDiasAntes),
				FechaFin:        fechaEvento.AddDate(0, 0, -ticket.FinDiasAntes),
				Estado:          1,
				UsuarioCreacion: &usuarioCreacion,
				FechaCreacion:   now,
			})
		}

		if err := entidad.Evento.CrearEvento(&evento); err != nil {
			return nil, fmt.Errorf("no se pudo crear evento %s: %w", seed.Titulo, err)
		}

		fecha := &model.Fecha{FechaEvento: fechaEvento}
		if err := db.Create(fecha).Error; err != nil {
			return nil, fmt.Errorf("no se pudo crear fecha para %s: %w", seed.Titulo, err)
		}

		horaInicio := time.Date(
			fechaEvento.Year(),
			fechaEvento.Month(),
			fechaEvento.Day(),
			seed.HoraInicio,
			seed.MinutoInicio,
			0,
			0,
			fechaEvento.Location(),
		)
		eventoFecha := &model.EventoFecha{
			EventoID:        evento.ID,
			FechaID:         fecha.ID,
			HoraInicio:      horaInicio,
			Estado:          1,
			UsuarioCreacion: &usuarioCreacion,
			FechaCreacion:   now,
		}
		if err := db.Create(eventoFecha).Error; err != nil {
			return nil, fmt.Errorf("no se pudo crear evento_fecha para %s: %w", seed.Titulo, err)
		}

		var sectores []model.Sector
		if err := db.Where("evento_id = ?", evento.ID).Find(&sectores).Error; err != nil {
			return nil, fmt.Errorf("no se pudo obtener sectores de %s: %w", seed.Titulo, err)
		}
		var perfiles []model.PerfilDePersona
		if err := db.Where("evento_id = ?", evento.ID).Find(&perfiles).Error; err != nil {
			return nil, fmt.Errorf("no se pudo obtener perfiles de %s: %w", seed.Titulo, err)
		}
		var tickets []model.TipoDeTicket
		if err := db.Where("evento_id = ?", evento.ID).Find(&tickets).Error; err != nil {
			return nil, fmt.Errorf("no se pudo obtener tickets de %s: %w", seed.Titulo, err)
		}

		sectorPrecio := make(map[string]float64)
		for _, sector := range seed.Sectores {
			sectorPrecio[sector.Nombre] = sector.PrecioBase
		}

		for _, sector := range sectores {
			base := sectorPrecio[sector.SectorTipo]
			for _, ticket := range tickets {
				multiplicadorTicket := 1.0
				for _, t := range seed.Tickets {
					if t.Nombre == ticket.Nombre {
						multiplicadorTicket = t.MultiplicadorPrecio
						break
					}
				}
				for _, perfil := range perfiles {
					multiplicadorPerfil := perfilMultiplicador[perfil.Nombre]
					if multiplicadorPerfil == 0 {
						multiplicadorPerfil = 1.0
					}
					precio := math.Round(base*multiplicadorTicket*multiplicadorPerfil*100) / 100
					perfilID := perfil.ID
					tarifa := &model.Tarifa{
						SectorID:          sector.ID,
						TipoDeTicketID:    ticket.ID,
						PerfilDePersonaID: &perfilID,
						Precio:            precio,
						Estado:            1,
						UsuarioCreacion:   &usuarioCreacion,
						FechaCreacion:     now,
					}
					if err := entidad.Tarifa.CrearTarifa(tarifa); err != nil {
						return nil, fmt.Errorf("no se pudo crear tarifa de %s: %w", seed.Titulo, err)
					}
				}
			}
		}

		if seed.Cupon != nil {
			cupon := &model.Cupon{
				Descripcion:     fmt.Sprintf("Cupón %s para %s", seed.Cupon.Codigo, seed.Titulo),
				Tipo:            seed.Cupon.Tipo,
				Valor:           seed.Cupon.Valor,
				EstadoCupon:     1,
				Codigo:          seed.Cupon.Codigo,
				UsoPorUsuario:   2,
				FechaInicio:     now,
				FechaFin:        fechaEvento,
				UsuarioCreacion: &usuarioCreacion,
				FechaCreacion:   now,
				EventoID:        evento.ID,
			}
			if err := entidad.Cupon.CrearCupon(cupon); err != nil {
				return nil, fmt.Errorf("no se pudo crear cupón %s: %w", seed.Cupon.Codigo, err)
			}
		}

		eventosCreados = append(eventosCreados, evento)
		logger.Infof("Evento %s creado con %d sectores y %d tickets", seed.Titulo, len(sectores), len(tickets))
	}

	return eventosCreados, nil
}

func isLocalhost(host string) bool {
	return host == "localhost" || host == "127.0.0.1"
}

func seedTicketsComprados(
	logger logging.Logger,
	db *gorm.DB,
	entidad *repository.NexiventPsqlEntidades,
	eventos []model.Evento,
	usuarios []model.Usuario,
) error {
	if len(usuarios) == 0 {
		return fmt.Errorf("no hay usuarios para asignar compras seed")
	}

	var ticketsExistentes int64
	if err := db.Model(&model.Ticket{}).Count(&ticketsExistentes).Error; err != nil {
		return fmt.Errorf("no se pudo contar tickets existentes: %w", err)
	}
	if ticketsExistentes > 0 {
		logger.Infof("Seed de tickets comprados omitido: ya hay %d tickets", ticketsExistentes)
		return nil
	}

	var metodoPago model.MetodoDePago
	if err := db.First(&metodoPago).Error; err != nil {
		return fmt.Errorf("no se pudo obtener método de pago para seeds: %w", err)
	}

	horaCompra := time.Now()
	maxEventos := len(eventos)
	if maxEventos > 3 {
		maxEventos = 3
	}

	for i := 0; i < maxEventos; i++ {
		ev := eventos[i]
		comprador := usuarios[i%len(usuarios)]

		var eventoFecha model.EventoFecha
		if err := db.Where("evento_id = ?", ev.ID).First(&eventoFecha).Error; err != nil {
			logger.Warnf("No se encontró evento_fecha para evento %s: %v", ev.Titulo, err)
			continue
		}

		var tarifas []model.Tarifa
		if err := db.
			Joins("JOIN tipo_de_ticket tdt ON tdt.tipo_de_ticket_id = tarifa.tipo_de_ticket_id").
			Where("tdt.evento_id = ?", ev.ID).
			Order("tarifa_id").
			Find(&tarifas).Error; err != nil {
			return fmt.Errorf("no se pudieron obtener tarifas para %s: %w", ev.Titulo, err)
		}
		if len(tarifas) == 0 {
			logger.Warnf("Evento %s no tiene tarifas para seed de tickets", ev.Titulo)
			continue
		}

		ticketsPorOrden := 2
		if len(tarifas) < ticketsPorOrden {
			ticketsPorOrden = len(tarifas)
		}

		seleccion := tarifas[:ticketsPorOrden]
		var total float64
		for _, tf := range seleccion {
			total += tf.Precio
		}

		orden := model.OrdenDeCompra{
			UsuarioID:        comprador.ID,
			MetodoDePagoID:   metodoPago.ID,
			Fecha:            horaCompra,
			FechaHoraIni:     horaCompra,
			Total:            math.Round(total*100) / 100,
			MontoFeeServicio: math.Round(total*0.05*100) / 100,
			EstadoDeOrden:    util.OrdenConfirmada.Codigo(),
		}
		if err := db.Create(&orden).Error; err != nil {
			return fmt.Errorf("no se pudo crear orden seed para %s: %w", ev.Titulo, err)
		}

		var tickets []model.Ticket
		for idxTf, tf := range seleccion {
			qr := fmt.Sprintf("SEED-%d-%d-%d", ev.ID, orden.ID, idxTf+1)
			tickets = append(tickets, model.Ticket{
				OrdenDeCompraID: &orden.ID,
				EventoFechaID:   eventoFecha.ID,
				TarifaID:        tf.ID,
				CodigoQR:        qr,
				EstadoDeTicket:  util.TicketVendido.Codigo(),
			})
		}

		if err := entidad.Ticket.CrearTicketsBatch(tickets); err != nil {
			return fmt.Errorf("no se pudieron crear tickets seed para %s: %w", ev.Titulo, err)
		}

		for _, tf := range seleccion {
			if err := entidad.Ticket.IncrementarVendidasPorSector(tf.SectorID, 1); err != nil {
				logger.Warnf("No se pudo incrementar vendidas para sector %d: %v", tf.SectorID, err)
			}
		}

		logger.Infof("Seed: orden %d con %d tickets vendidos para evento %s", orden.ID, len(tickets), ev.Titulo)
	}

	return nil
}
