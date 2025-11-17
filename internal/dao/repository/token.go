package repository

import (
	"gorm.io/gorm"

	"github.com/Nexivent/nexivent-backend/logging"
)

type Token struct {
	logger logging.Logger
	DB     *gorm.DB
}
