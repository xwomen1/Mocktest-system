package config

import (
	"fmt"
	"time"
)

// server configuration
type ServerConfig struct {
	Host string `yaml:"host" env:"SERVER_HOST" default:"0.0.0.0"`
	Port int    `yaml:"port" env:"SERVER_PORT" default:"50051"`

	// timeouts
	ReadTimeout  time.Duration `yaml:"read_timeout" env:"READ_TIMEOUT" default:"30s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env:"WRITE_TIMEOUT" default:"30s"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" env:"IDLE_TIMEOUT" default:"60s"`

	// security
	EnableTLS   bool   `yaml:"enable_tls" env:"ENABLE_TLS" default:"false"`
	TLSCertPath string `yaml:"tls_cert_path" env:"TLS_CERT_PATH"`
	TLSKeyPath  string `yaml:"tls_key_path" env:"TLS_KEY_PATH"`
}

// NATS message queue configuration
type NATSConfig struct {
	URL       string `yaml:"url" env:"NATS_URL" default:"nats://localhost:4222"`
	ClusterID string `yaml:"cluster_id" env:"NATS_CLUSTER_ID" default:"test-cluster"`
	ClientID  string `yaml:"client_id" env:"NATS_CLIENT_ID"`

	// connection
	MaxReconnects int           `yaml:"max_reconnects" env:"NATS_MAX_RECONNECTS" default:"-1"`
	ReconnectWait time.Duration `yaml:"reconnect_wait" env:"NATS_RECONNECT_WAIT" default:"2s"`
	Timeout       time.Duration `yaml:"timeout" env:"NATS_TIMEOUT" default:"5s"`
}

// logging configuration
type LoggingConfig struct {
	Level      string `yaml:"level" env:"LOG_LEVEL" default:"info"`
	Format     string `yaml:"format" env:"LOG_FORMAT" default:"json"`   // json or text
	Output     string `yaml:"output" env:"LOG_OUTPUT" default:"stdout"` // stdout, stderr, or file path
	EnableJSON bool   `yaml:"enable_json" env:"LOG_ENABLE_JSON" default:"true"`

	// file losging (if output is file)
	MaxSize    int `yaml:"max_size" env:"LOG_MAX_SIZE" default:"100"` // MB
	MaxBackups int `yaml:"max_backups" env:"LOG_MAX_BACKUPS" default:"10"`
	MaxAge     int `yaml:"max_age" env:"LOG_MAX_AGE" default:"30"` // days
}

type RegistryConfig struct {
	HeartbeatInterval time.Duration `yaml:"heartbeat_interval" env:"HEARTBEAT_INTERVAL" default:"30s"`
	HeartbeatTimeout  time.Duration `yaml:"heartbeat_timeout" env:"HEARTBEAT_TIMEOUT" default:"90s"`

	LoadBalancingStrategy string `yaml:"load_balancing_strategy" env:"LB_STRATEGY" default:"round_robin"` // round_robin, least_connections, random

	CacheTTL time.Duration `yaml:"cache_ttl" env:"CACHE_TTL" default:"5m"`
}

type Config struct {
	Environment string         `yaml:"environment" env:"ENVIRONMENT" default:"development"`
	Server      ServerConfig   `yaml:"server"`
	NATS        NATSConfig     `yaml:"nats"`
	Logging     LoggingConfig  `yaml:"logging"`
	Registry    RegistryConfig `yaml:"registry"`

	Features struct {
		EnableMetrics   bool `yaml:"enable_metrics" env:"ENABLE_METRICS" default:"true"`
		EnableTracing   bool `yaml:"enable_tracing" env:"ENABLE_TRACING" default:"false"`
		EnableProfiling bool `yaml:"enable_profiling" env:"ENABLE_PROFILING" default:"false"`
	} `yaml:"features"`
}

func (c *Config) Validate() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	if c.NATS.URL == "" {
		return fmt.Errorf("NATS URL is required")
	}

	if c.Logging.Level != "debug" && c.Logging.Level != "info" &&
		c.Logging.Level != "warn" && c.Logging.Level != "error" {
		return fmt.Errorf("invalid log level: %s", c.Logging.Level)
	}

	return nil
}
