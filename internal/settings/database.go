package settings

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

func OpenDB(cfg Config) (*sql.DB, error) {
	// Usa sql.Open() para crear un pool de conexiones vacío, usando el DSN de la estructura
	// de configuración.
	db, err := sql.Open("postgres", cfg.DB.URL)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)
	db.SetConnMaxIdleTime(cfg.DB.MaxIdleTime)

	// Crea un contexto con un tiempo límite de 5 segundos.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Usa PingContext() para establecer una nueva conexión a la base de datos, pasando el
	// contexto que creamos arriba como parámetro. Si la conexión no se puede
	// establecer exitosamente dentro del tiempo límite de 5 segundos, esto devolverá un
	// error. Si obtenemos este error, o cualquier otro, cerramos el pool de conexiones y
	// devolvemos el error.
	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	// Devuelve el pool de conexiones sql.DB.
	return db, nil
}
