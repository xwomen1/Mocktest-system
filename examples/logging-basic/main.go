package main

import (
	"fmt"
	"time"

	"upm-simple/pkg/config"
	"upm-simple/pkg/logger"
)

func main() {
	fmt.Println("=== Logging Framework Example ===")

	// load config
	cfg, err := config.LoadConfig("development")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		fmt.Println("Using default logger config")
	}

	logConfig := logger.Config{
		Level:        logger.InfoLevel,
		Encoding:     "json",
		OutputPath:   "stdout",
		EnableCaller: true,
	}

	if cfg != nil {
		// map config levels
		switch cfg.Logging.Level {
		case "debug":
			logConfig.Level = logger.DebugLevel
		case "info":
			logConfig.Level = logger.InfoLevel
		case "warn":
			logConfig.Level = logger.WarnLevel
		case "error":
			logConfig.Level = logger.ErrorLevel
		}

		logConfig.Encoding = cfg.Logging.Format
		logConfig.OutputPath = cfg.Logging.Output
	}

	log, err := logger.FromConfig(logConfig)
	if err != nil {
		fmt.Printf("Error creating logger: %v\n", err)
		return
	}

	defer log.Sync()

	// test different log levels
	fmt.Println("\nTesting log levels:")

	log.Debug("This is a debug message",
		logger.FieldString("user", "john"),
		logger.FieldInt("attempt", 1))

	log.Info("Service started",
		logger.FieldString("service", "registry"),
		logger.FieldInt("port", 50051),
		logger.FieldString("version", "1.0.0"))

	log.Warn("High memory usage detected",
		logger.FieldString("component", "memory"),
		logger.FieldInt64("usage_mb", 850),
		logger.FieldInt64("limit_mb", 1024))

	log.Error("Failed to connect to database",
		logger.FieldString("database", "postgres"),
		logger.FieldInt("port", 5432),
		logger.FieldError(fmt.Errorf("connection timeout")))

	// test with structured logging
	fmt.Println("\nTesting structured logging:")

	serviceLogger := log.With(
		logger.FieldString("service", "api-gateway"),
		logger.FieldString("instance", "gateway-1"),
	)

	serviceLogger.Info("Processing request",
		logger.FieldString("method", "POST"),
		logger.FieldString("path", "/api/v1/users"),
		logger.FieldInt("status", 200),
		FieldDuration("duration", 150*time.Millisecond))

	// test context (simulated)
	fmt.Println("\nTesting different scenarios:")

	for i := 1; i <= 3; i++ {
		requestLogger := log.With(
			logger.FieldString("request_id", fmt.Sprintf("req-%d", i)),
			logger.FieldString("client_ip", fmt.Sprintf("192.168.1.%d", i)),
		)

		requestLogger.Info("Request processed",
			logger.FieldInt("iterations", i*10),
			logger.FieldBool("success", true))
	}

	fmt.Println("\n Logging test completed!")
	fmt.Println("Check the structured JSON output above.")
}

// create duration field
func FieldDuration(key string, value time.Duration) logger.Field {
	return logger.Field{Key: key, Value: value.String()}
}
