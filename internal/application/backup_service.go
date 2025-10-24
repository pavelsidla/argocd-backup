package application

import (
	"context"
	"fmt"
	"log"
	"time"

	"argocd-backup/internal/domain"
)

type BackupService struct {
	exporter domain.Exporter
	storage  domain.Storage
	metrics  domain.Metrics
}

func NewBackupService(
	exporter domain.Exporter,
	storage domain.Storage,
	metrics domain.Metrics,
) *BackupService {
	return &BackupService{
		exporter: exporter,
		storage:  storage,
		metrics:  metrics,
	}
}

func (s *BackupService) Run(ctx context.Context) (err error) {
	log.Println("Backup service starting...")
	startTime := time.Now()

	defer func() {
		duration := time.Since(startTime)
		if err != nil {
			s.metrics.ReportFailure(duration, err)
		} else {
			s.metrics.ReportSuccess(duration)
		}
	}()

	backup, err := s.exporter.Export(ctx)
	if err != nil {
		log.Printf("Export failed: %v", err)
		return fmt.Errorf("export failed: %w", err)
	}

	path, err := s.storage.Upload(ctx, backup)
	if err != nil {
		log.Printf("Upload failed: %v", err)
		return fmt.Errorf("storage upload failed: %w", err)
	}

	log.Printf("Backup service completed successfully. Path: %s", path)
	return nil
}
