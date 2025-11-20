package main

import (
	"fmt"
	"log"
	"math"
	"time"

	config "github.com/Nexivent/nexivent-backend/internal/config"
	"github.com/Nexivent/nexivent-backend/internal/dao/model"
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
	if err := seedEventos(logger, db, entidad, usuarios, categorias); err != nil {
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
) error {
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
		Titulo         string
		Descripcion    string
		Lugar          string
		Categoria      string
		OrganizadorIdx int
		DiasHastaFecha int
		HoraInicio     int
		MinutoInicio   int
		Sectores       []sectorSeed
		Perfiles       []string
		Tickets        []ticketSeed
		Cupon          *couponSeed
	}

	eventos := []eventoSeed{
		{
			Titulo:         "Festival de Tecnología",
			Descripcion:    "Charlas, workshops y demo day con startups de IA y cloud.",
			Lugar:          "Centro de Convenciones Costa Verde",
			Categoria:      "Tecnología",
			OrganizadorIdx: 0,
			DiasHastaFecha: 10,
			HoraInicio:     9,
			MinutoInicio:   0,
			Sectores: []sectorSeed{
				{Nombre: "VIP", Capacidad: 150, PrecioBase: 220},
				{Nombre: "General", Capacidad: 800, PrecioBase: 120},
				{Nombre: "Talleres", Capacidad: 300, PrecioBase: 160},
			},
			Perfiles: []string{"Profesional", "Estudiante", "Founder"},
			Tickets: []ticketSeed{
				{Nombre: "Preventa", InicioDiasAntes: 60, FinDiasAntes: 20, MultiplicadorPrecio: 0.9},
				{Nombre: "General", InicioDiasAntes: 30, FinDiasAntes: 1, MultiplicadorPrecio: 1.0},
				{Nombre: "Último minuto", InicioDiasAntes: 7, FinDiasAntes: 0, MultiplicadorPrecio: 1.15},
			},
			Cupon: &couponSeed{Codigo: "TECH10", Valor: 10, Tipo: 1},
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
			Titulo:         "Partido de Exhibición",
			Descripcion:    "Equipos históricos se enfrentan en un duelo amistoso.",
			Lugar:          "Estadio Nacional",
			Categoria:      "Deportes",
			OrganizadorIdx: 2,
			DiasHastaFecha: 5,
			HoraInicio:     18,
			MinutoInicio:   0,
			Sectores: []sectorSeed{
				{Nombre: "Palco", Capacidad: 200, PrecioBase: 300},
				{Nombre: "Occidente", Capacidad: 1800, PrecioBase: 180},
				{Nombre: "Oriente", Capacidad: 2200, PrecioBase: 150},
				{Nombre: "Popular", Capacidad: 4000, PrecioBase: 90},
			},
			Perfiles: []string{"Adulto", "Niño"},
			Tickets: []ticketSeed{
				{Nombre: "General", InicioDiasAntes: 20, FinDiasAntes: 0, MultiplicadorPrecio: 1.0},
			},
			Cupon: &couponSeed{Codigo: "GOLES15", Valor: 15, Tipo: 1},
		},
		{
			Titulo:         "Obra de Teatro Urbano",
			Descripcion:    "Dramaturgia contemporánea con elenco joven y música en vivo.",
			Lugar:          "Teatro Municipal",
			Categoria:      "Teatro",
			OrganizadorIdx: 0,
			DiasHastaFecha: 8,
			HoraInicio:     19,
			MinutoInicio:   30,
			Sectores: []sectorSeed{
				{Nombre: "Platea", Capacidad: 300, PrecioBase: 95},
				{Nombre: "Mezzanine", Capacidad: 200, PrecioBase: 75},
				{Nombre: "Galería", Capacidad: 150, PrecioBase: 55},
			},
			Perfiles: []string{"Adulto", "Estudiante"},
			Tickets: []ticketSeed{
				{Nombre: "General", InicioDiasAntes: 25, FinDiasAntes: 0, MultiplicadorPrecio: 1.0},
			},
			Cupon: &couponSeed{Codigo: "TEATRO8", Valor: 8, Tipo: 1},
		},
		{
			Titulo:         "Feria Gastronómica de Verano",
			Descripcion:    "Food trucks, cerveza artesanal y shows acústicos.",
			Lugar:          "Parque de la Exposición",
			Categoria:      "Gastronomía",
			OrganizadorIdx: 3,
			DiasHastaFecha: 12,
			HoraInicio:     12,
			MinutoInicio:   0,
			Sectores: []sectorSeed{
				{Nombre: "Degustación", Capacidad: 400, PrecioBase: 80},
				{Nombre: "General", Capacidad: 1800, PrecioBase: 45},
			},
			Perfiles: []string{"Adulto", "Niño"},
			Tickets: []ticketSeed{
				{Nombre: "Pase día", InicioDiasAntes: 40, FinDiasAntes: 0, MultiplicadorPrecio: 1.0},
			},
			Cupon: &couponSeed{Codigo: "FOOD12", Valor: 12, Tipo: 1},
		},
		{
			Titulo:         "Cumbre de Startups",
			Descripcion:    "Rondas de pitch, VC office hours y paneles de inversión.",
			Lugar:          "WeWork San Isidro",
			Categoria:      "Negocios",
			OrganizadorIdx: 1,
			DiasHastaFecha: 20,
			HoraInicio:     10,
			MinutoInicio:   0,
			Sectores: []sectorSeed{
				{Nombre: "Founder Pass", Capacidad: 120, PrecioBase: 260},
				{Nombre: "General", Capacidad: 500, PrecioBase: 140},
			},
			Perfiles: []string{"Founder", "Inversionista", "Asistente"},
			Tickets: []ticketSeed{
				{Nombre: "Preventa", InicioDiasAntes: 50, FinDiasAntes: 15, MultiplicadorPrecio: 0.9},
				{Nombre: "General", InicioDiasAntes: 30, FinDiasAntes: 0, MultiplicadorPrecio: 1.05},
			},
			Cupon: &couponSeed{Codigo: "VC20", Valor: 20, Tipo: 1},
		},
		{
			Titulo:         "Festival Indie",
			Descripcion:    "Bandas emergentes, arte urbano y zona de food trucks.",
			Lugar:          "Campo Mar",
			Categoria:      "Conciertos",
			OrganizadorIdx: 2,
			DiasHastaFecha: 25,
			HoraInicio:     17,
			MinutoInicio:   0,
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
			Titulo:         "Carrera 10K Ciudad",
			Descripcion:    "Circuito urbano con chip de cronometraje y medallas finisher.",
			Lugar:          "Circuito de Playas",
			Categoria:      "Deportes",
			OrganizadorIdx: 3,
			DiasHastaFecha: 18,
			HoraInicio:     7,
			MinutoInicio:   30,
			Sectores: []sectorSeed{
				{Nombre: "Competitivo", Capacidad: 800, PrecioBase: 85},
				{Nombre: "Recreativo", Capacidad: 1200, PrecioBase: 60},
			},
			Perfiles: []string{"Adulto", "Estudiante"},
			Tickets: []ticketSeed{
				{Nombre: "General", InicioDiasAntes: 40, FinDiasAntes: 0, MultiplicadorPrecio: 1.0},
			},
			Cupon: &couponSeed{Codigo: "RUNNER7", Valor: 7, Tipo: 1},
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

	for _, seed := range eventos {
		if seed.OrganizadorIdx >= len(usuarios) {
			return fmt.Errorf("organizador fuera de rango para evento %s", seed.Titulo)
		}
		organizador := usuarios[seed.OrganizadorIdx]
		categoriaID := categoriaPorNombre[seed.Categoria]
		if categoriaID == 0 {
			return fmt.Errorf("no se encontró categoría %s", seed.Categoria)
		}

		fechaEvento := now.AddDate(0, 0, seed.DiasHastaFecha)
		usuarioCreacion := organizador.ID
		evento := model.Evento{
			Titulo:          seed.Titulo,
			OrganizadorID:   organizador.ID,
			CategoriaID:     categoriaID,
			Descripcion:     seed.Descripcion,
			Lugar:           seed.Lugar,
			EventoEstado:    1,
			Estado:          1,
			UsuarioCreacion: &usuarioCreacion,
			FechaCreacion:   now,
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
			return fmt.Errorf("no se pudo crear evento %s: %w", seed.Titulo, err)
		}

		fecha := &model.Fecha{FechaEvento: fechaEvento}
		if err := db.Create(fecha).Error; err != nil {
			return fmt.Errorf("no se pudo crear fecha para %s: %w", seed.Titulo, err)
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
			return fmt.Errorf("no se pudo crear evento_fecha para %s: %w", seed.Titulo, err)
		}

		var sectores []model.Sector
		if err := db.Where("evento_id = ?", evento.ID).Find(&sectores).Error; err != nil {
			return fmt.Errorf("no se pudo obtener sectores de %s: %w", seed.Titulo, err)
		}
		var perfiles []model.PerfilDePersona
		if err := db.Where("evento_id = ?", evento.ID).Find(&perfiles).Error; err != nil {
			return fmt.Errorf("no se pudo obtener perfiles de %s: %w", seed.Titulo, err)
		}
		var tickets []model.TipoDeTicket
		if err := db.Where("evento_id = ?", evento.ID).Find(&tickets).Error; err != nil {
			return fmt.Errorf("no se pudo obtener tickets de %s: %w", seed.Titulo, err)
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
						return fmt.Errorf("no se pudo crear tarifa de %s: %w", seed.Titulo, err)
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
				return fmt.Errorf("no se pudo crear cupón %s: %w", seed.Cupon.Codigo, err)
			}
		}

		logger.Infof("Evento %s creado con %d sectores y %d tickets", seed.Titulo, len(sectores), len(tickets))
	}

	return nil
}

func isLocalhost(host string) bool {
	return host == "localhost" || host == "127.0.0.1"
}
