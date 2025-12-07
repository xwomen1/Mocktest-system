package config

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Loader struct {
	viper *viper.Viper
	path  string
}

type ConfigUpdate struct {
	Config *Config
	Error  error
}

func NewLoader() *Loader {
	v := viper.New()

	v.SetDefault("environment", "development")
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 50051)
	v.SetDefault("server.read_timeout", "30s")
	v.SetDefault("server.write_timeout", "30s")
	v.SetDefault("server.idle_timeout", "60s")
	v.SetDefault("nats.url", "nats://localhost:4222")
	v.SetDefault("nats.cluster_id", "test-cluster")
	v.SetDefault("nats.max_reconnects", -1)
	v.SetDefault("nats.reconnect_wait", "2s")
	v.SetDefault("nats.timeout", "5s")
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "json")
	v.SetDefault("logging.output", "stdout")
	v.SetDefault("logging.enable_json", true)
	v.SetDefault("registry.heartbeat_interval", "30s")
	v.SetDefault("registry.heartbeat_timeout", "90s")
	v.SetDefault("registry.load_balancing_strategy", "round_robin")
	v.SetDefault("registry.cache_ttl", "5m")
	v.SetDefault("features.enable_metrics", true)

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./configs")
	v.AddConfigPath("/etc/upm/")

	v.AutomaticEnv()
	v.SetEnvPrefix("UPM")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	return &Loader{
		viper: v,
	}
}

// load configuration from file and environment
func (l *Loader) Load() (*Config, error) {
	configFound := false

	configPaths := []string{
		"configs/dev/config.yaml",
		"configs/config.yaml",
		"./config.yaml",
	}

	for _, path := range configPaths {
		l.viper.SetConfigFile(path)
		if err := l.viper.ReadInConfig(); err == nil {
			l.path = path
			configFound = true
			fmt.Printf("Config loaded from: %s\n", path)
			break
		}
	}

	if !configFound {
		fmt.Println("No config file found, using defaults and environment variables")
	}

	var config Config
	if err := l.viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

// load configuration from a specific file
func (l *Loader) LoadFromFile(path string) (*Config, error) {
	l.viper.SetConfigFile(path)
	l.path = path
	return l.Load()
}

func (l *Loader) Watch() <-chan ConfigUpdate {
	updates := make(chan ConfigUpdate, 1) // Add buffer

	if l.path == "" {
		// No config file to watch, return closed channel
		go func() {
			updates <- ConfigUpdate{
				Error: fmt.Errorf("no config file to watch"),
			}
			close(updates)
		}()
		return updates
	}

	// Start watching in goroutine
	go l.watchConfig(updates)

	return updates
}

func (l *Loader) watchConfig(updates chan<- ConfigUpdate) {
	defer close(updates)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		updates <- ConfigUpdate{Error: fmt.Errorf("failed to create watcher: %w", err)}
		return
	}
	defer watcher.Close()

	// Watch config file and directory
	configDir := filepath.Dir(l.path)
	if err := watcher.Add(configDir); err != nil {
		updates <- ConfigUpdate{Error: fmt.Errorf("failed to watch directory: %w", err)}
		return
	}

	// Send initial config
	if config, err := l.Load(); err == nil {
		updates <- ConfigUpdate{Config: config}
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			// Check if config file was modified
			if event.Op&fsnotify.Write == fsnotify.Write &&
				filepath.Base(event.Name) == filepath.Base(l.path) {
				time.Sleep(100 * time.Millisecond) // Wait for write to complete

				config, err := l.Load()
				updates <- ConfigUpdate{
					Config: config,
					Error:  err,
				}
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			updates <- ConfigUpdate{Error: err}
		}
	}
}

func (l *Loader) GetString(key string) string {
	return l.viper.GetString(key)
}

func (l *Loader) GetInt(key string) int {
	return l.viper.GetInt(key)
}

func (l *Loader) GetBool(key string) bool {
	return l.viper.GetBool(key)
}

func (l *Loader) GetDuration(key string) time.Duration {
	return l.viper.GetDuration(key)
}

func (l *Loader) GetViper() *viper.Viper {
	return l.viper
}
