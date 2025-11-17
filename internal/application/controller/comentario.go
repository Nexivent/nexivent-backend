package controller

import (
	"github.com/Nexivent/nexivent-backend/logging"
	"github.com/Nexivent/nexivent-backend/internal/dao/repository"
)

type ComentarioController struct {
	Logger logging.Logger
	DB     *repository.NexiventPsqlEntidades
}
