package storage

import (
	"context"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/Nexivent/nexivent-backend/internal/config"
	"github.com/Nexivent/nexivent-backend/logging"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Storage handles S3 presigned URL generation.
type S3Storage struct {
	logger         logging.Logger
	client         *s3.Client
	presignClient  *s3.PresignClient
	bucket         string
	prefix         string
	presignExpires time.Duration
}

func NewS3Storage(logger logging.Logger, cfg *config.ConfigEnv) (*S3Storage, error) {
	if cfg.AwsRegion == "" || cfg.AwsS3Bucket == "" {
		return nil, fmt.Errorf("missing AWS S3 configuration")
	}

	awsCfg, err := awsConfig.LoadDefaultConfig(context.Background(), awsConfig.WithRegion(cfg.AwsRegion))
	if err != nil {
		return nil, fmt.Errorf("loading AWS config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg)
	presignClient := s3.NewPresignClient(client)

	exp := time.Duration(cfg.AwsS3UploadDuration) * time.Second
	if exp <= 0 {
		exp = 15 * time.Minute
	}

	return &S3Storage{
		logger:         logger,
		client:         client,
		presignClient:  presignClient,
		bucket:         cfg.AwsS3Bucket,
		prefix:         strings.Trim(cfg.AwsS3Prefix, "/"),
		presignExpires: exp,
	}, nil
}

// GeneratePresignedPut creates a presigned PUT URL for direct upload from the frontend.
func (s *S3Storage) GeneratePresignedPut(ctx context.Context, fileName, contentType string) (string, string, error) {
	safeName := sanitizeFileName(fileName)
	key := fmt.Sprintf("%d-%s", time.Now().UnixNano(), safeName)
	if s.prefix != "" {
		key = path.Join(s.prefix, key)
	}

	req, err := s.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(s.presignExpires))
	if err != nil {
		s.logger.Errorf("failed to presign S3 upload: %v", err)
		return "", "", err
	}

	return key, req.URL, nil
}

func sanitizeFileName(name string) string {
	name = path.Base(name)
	name = strings.ReplaceAll(name, " ", "_")
	return name
}
