package s3

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"argocd-backup/internal/domain"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Storage struct {
	client *s3.Client
	bucket string
	prefix string
}

func NewStorage(ctx context.Context, region, bucket, prefix string) (*Storage, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &Storage{
		client: s3.NewFromConfig(cfg),
		bucket: bucket,
		prefix: prefix,
	}, nil
}

func (s *Storage) Upload(ctx context.Context, backup *domain.Backup) (string, error) {
	s3Key := backup.Filename
	if s.prefix != "" {
		s3Key = fmt.Sprintf("%s/%s", s.prefix, backup.Filename)
	}

	s3Path := fmt.Sprintf("s3://%s/%s", s.bucket, s3Key)

	log.Printf("Uploading backup to %s", s3Path)

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s3Key),
		Body:   bytes.NewReader(backup.Data),
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	return s3Path, nil
}
