package argocd

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"

	"argocd-backup/internal/domain"
)

type Exporter struct {
	namespace string
}

func NewExporter(namespace string) *Exporter {
	return &Exporter{namespace: namespace}
}

// Export executes the 'argocd admin export' command.
func (e *Exporter) Export(ctx context.Context) (*domain.Backup, error) {
	log.Printf("Starting Argo CD export for namespace '%s'", e.namespace)

	cmd := exec.CommandContext(ctx, "argocd", "admin", "export", "-n", e.namespace)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()
	if err != nil {
		log.Printf("Stderr: %s", errBuf.String())
		return nil, fmt.Errorf("argocd export command failed: %w", err)
	}

	backupData := outBuf.Bytes()
	log.Printf("Successfully exported %d bytes from Argo CD", len(backupData))

	backup := domain.NewBackup(backupData)

	return backup, nil
}
