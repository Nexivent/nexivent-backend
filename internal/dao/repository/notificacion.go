package repository

import (
	"github.com/Loui27/nexivent-backend/logging"
	"gorm.io/gorm"
	"github.com/Loui27/nexivent-backend/internal/dao/model"
)

type Notificacion struct{
	logger logging.Logger
	PostgresqlDB *gorm.DB
}