package config

import (
	"fmt"
	"net"
	"os"
)

type Config struct {
	ArgoNamespace string
	S3Bucket      string
	S3Region      string
	S3Prefix      string
	DDAgentHost   string
	DDEnv         string
}

func Load() (*Config, error) {
	cfg := &Config{
		ArgoNamespace: os.Getenv("ARGO_NAMESPACE"),
		S3Bucket:      os.Getenv("S3_BUCKET"),
		S3Region:      os.Getenv("S3_REGION"),
		S3Prefix:      os.Getenv("S3_PATH_PREFIX"),
		DDAgentHost:   os.Getenv("DD_AGENT_HOST"),
		DDEnv:         os.Getenv("DD_ENV"),
	}

	if cfg.ArgoNamespace == "" {
		return nil, fmt.Errorf("ARGO_NAMESPACE must be set")
	}
	if cfg.S3Bucket == "" {
		return nil, fmt.Errorf("S3_BUCKET must be set")
	}
	if cfg.S3Region == "" {
		return nil, fmt.Errorf("S3_REGION must be set")
	}

	return cfg, nil
}

func (c *Config) DogStatsdAddr() string {
	if c.DDAgentHost == "" {
		return ""
	}
	return net.JoinHostPort(c.DDAgentHost, "8125")
}
