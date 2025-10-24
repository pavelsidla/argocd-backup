package datadog

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
)

const (
	metricPrefix     = "argocd.backup"
	serviceCheckName = "argocd.backup.status"
	successGauge     = "argocd.backup.success"
	durationMetric   = "argocd.backup.duration_seconds"
)

type Reporter struct {
	client statsd.ClientInterface
}

func NewReporter(statsdAddr, env, argoNs string) (*Reporter, error) {
	if statsdAddr == "" {
		log.Println("DD_AGENT_HOST not set. Metrics will not be sent.")
		return &Reporter{client: &statsd.NoOpClient{}}, nil
	}

	log.Printf("Sending metrics to DogStatsD at %s", statsdAddr)

	defaultTags := []string{
		"app:argocd-backup",
		fmt.Sprintf("argo_namespace:%s", argoNs),
	}

	if env != "" {
		defaultTags = append(defaultTags, fmt.Sprintf("env:%s", env))
	}

	if customTags := os.Getenv("DD_TAGS"); customTags != "" {
		tags := strings.Split(customTags, ",")
		defaultTags = append(defaultTags, tags...)
	}

	client, err := statsd.New(statsdAddr, statsd.WithTags(defaultTags))
	if err != nil {
		return nil, fmt.Errorf("could not create statsd client: %w", err)
	}

	return &Reporter{client: client}, nil
}

func (r *Reporter) ReportSuccess(duration time.Duration) {
	log.Println("Reporting success to Datadog")
	r.client.Timing(durationMetric, duration, nil, 1)
	r.client.Gauge(successGauge, 1, nil, 1)

	sc := &statsd.ServiceCheck{
		Name:    serviceCheckName,
		Status:  statsd.Ok,
		Message: "Backup successful",
	}
	r.client.ServiceCheck(sc)
}

func (r *Reporter) ReportFailure(duration time.Duration, reportedErr error) {
	log.Printf("Reporting failure to Datadog: %v", reportedErr)
	r.client.Timing(durationMetric, duration, nil, 1)
	r.client.Gauge(successGauge, 0, nil, 1)

	sc := &statsd.ServiceCheck{
		Name:    serviceCheckName,
		Status:  statsd.Critical,
		Message: reportedErr.Error(),
	}
	r.client.ServiceCheck(sc)
}

func (r *Reporter) Close() {
	if err := r.client.Close(); err != nil {
		log.Printf("Error closing Datadog client: %v", err)
	}
	if err := r.client.Flush(); err != nil {
		log.Printf("Error flushing Datadog metrics: %v", err)
	}
	log.Println("Datadog client closed and flushed.")
}
