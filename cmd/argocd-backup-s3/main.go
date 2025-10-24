package main

import (
	"context"
	"log"
	"os"
	"time"

	"argocd-backup/internal/application"
	"argocd-backup/internal/config"
	"argocd-backup/internal/infrastructure/argocd"
	"argocd-backup/internal/infrastructure/datadog"
	"argocd-backup/internal/infrastructure/s3"
)

func main() {
	log.Println("Starting Argo CD Backup Job...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	exporter := argocd.NewExporter(cfg.ArgoNamespace)

	storage, err := s3.NewStorage(ctx, cfg.S3Region, cfg.S3Bucket, cfg.S3Prefix)
	if err != nil {
		log.Fatalf("Failed to create S3 storage adapter: %v", err)
	}

	metrics, err := datadog.NewReporter(cfg.DogStatsdAddr(), cfg.DDEnv, cfg.ArgoNamespace)
	if err != nil {
		log.Fatalf("Failed to create Datadog reporter: %v", err)
	}

	defer metrics.Close()

	backupService := application.NewBackupService(exporter, storage, metrics)

	if err := backupService.Run(ctx); err != nil {
		log.Printf("Job failed: %v", err)
		os.Exit(1)
	}

	log.Println("Argo CD Backup Job finished successfully.")
}
