package repository

import (
	"gorm.io/gorm"
)

type Comentario struct {
	DB *gorm.DB
}
