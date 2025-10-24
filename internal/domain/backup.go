package domain

import (
	"context"
	"fmt"
	"time"
)

type Backup struct {
	Data     []byte
	Filename string
}

func NewBackup(data []byte) *Backup {
	filename := fmt.Sprintf(
		"argocd-backup-%s.yml",
		time.Now().UTC().Format("2006-01-02_15-04-05"),
	)
	return &Backup{
		Data:     data,
		Filename: filename,
	}
}

type Exporter interface {
	Export(ctx context.Context) (*Backup, error)
}

type Storage interface {
	Upload(ctx context.Context, backup *Backup) (string, error)
}

type Metrics interface {
	ReportSuccess(duration time.Duration)
	ReportFailure(duration time.Duration, reportedErr error)
	Close()
}
