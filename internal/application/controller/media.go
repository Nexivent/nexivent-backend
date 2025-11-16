package controller

import (
	"context"
	"strings"

	"github.com/Loui27/nexivent-backend/errors"
	"github.com/Loui27/nexivent-backend/internal/application/service/storage"
	"github.com/Loui27/nexivent-backend/internal/schemas"
	"github.com/Loui27/nexivent-backend/logging"
)

type MediaController struct {
	Logger  logging.Logger
	Storage *storage.S3Storage
}

func NewMediaController(
	logger logging.Logger,
	storage *storage.S3Storage,
) *MediaController {
	return &MediaController{
		Logger:  logger,
		Storage: storage,
	}
}

func (mc *MediaController) GenerateUploadURL(ctx context.Context, req schemas.PresignUploadRequest) (*schemas.PresignUploadResponse, *errors.Error) {
	if req.FileName == "" || req.ContentType == "" {
		return nil, &errors.UnprocessableEntityError.InvalidRequestBody
	}

	if !strings.HasPrefix(req.ContentType, "image/") && !strings.HasPrefix(req.ContentType, "video/") {
		return nil, &errors.BadRequestError.InvalidUploadMimeType
	}

	key, url, err := mc.Storage.GeneratePresignedPut(ctx, req.FileName, req.ContentType)
	if err != nil {
		return nil, &errors.BadRequestError.UploadURLNotCreated
	}

	return &schemas.PresignUploadResponse{
		UploadURL: url,
		Key:       key,
	}, nil
}
