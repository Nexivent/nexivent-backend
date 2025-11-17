package repository

import (
	"github.com/Nexivent/nexivent-backend/logging"
	"gorm.io/gorm"
	// "github.com/Loui27/nexivent-backend/internal/dao/model"
)

type ComprobanteDePago struct {
	logger       logging.Logger
	PostgresqlDB *gorm.DB
}
