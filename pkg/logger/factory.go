package logger

import (
	"context"
	"fmt"
	"sync"
)

var (
	defaultLogger Logger
	loggerOnce    sync.Once
)

// creates a new logger based on config
func New(config Config) (Logger, error) {
	return NewZapLogger(config)
}

func Default() Logger {
	loggerOnce.Do(func() {
		var err error
		defaultLogger, err = DefaultZapLogger()
		if err != nil {
			// Fallback to a basic logger
			fmt.Printf("Failed to create default logger: %v\n", err)
			defaultLogger = &NoopLogger{}
		}
	})
	return defaultLogger
}

// setDefault sets the default logger
func SetDefault(logger Logger) {
	defaultLogger = logger
}

// fromConfig creates logger from config struct
func FromConfig(cfg Config) (Logger, error) {
	return New(cfg)
}

// noopLogger is a logger that does nothing (for fallback)
type NoopLogger struct{}

func (n *NoopLogger) Debug(msg string, fields ...Field)      {}
func (n *NoopLogger) Info(msg string, fields ...Field)       {}
func (n *NoopLogger) Warn(msg string, fields ...Field)       {}
func (n *NoopLogger) Error(msg string, fields ...Field)      {}
func (n *NoopLogger) Fatal(msg string, fields ...Field)      {}
func (n *NoopLogger) Panic(msg string, fields ...Field)      {}
func (n *NoopLogger) WithContext(ctx context.Context) Logger { return n }
func (n *NoopLogger) With(fields ...Field) Logger            { return n }
func (n *NoopLogger) WithPrefix(prefix string) Logger        { return n }
func (n *NoopLogger) Sync() error                            { return nil }
func (n *NoopLogger) SetLevel(level Level)                   {}
func (n *NoopLogger) GetLevel() Level                        { return InfoLevel }
