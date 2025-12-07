package logger

import (
	"context"
	"time"
)

type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
)

func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	case PanicLevel:
		return "panic"
	default:
		return "unknown"
	}
}

type Field struct {
	Key   string
	Value interface{}
}

type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	Panic(msg string, fields ...Field)

	WithContext(ctx context.Context) Logger

	With(fields ...Field) Logger

	WithPrefix(prefix string) Logger

	Sync() error

	SetLevel(level Level)

	GetLevel() Level
}

// config holds logger configuration
type Config struct {
	Level        Level  `json:"level" yaml:"level"`
	Encoding     string `json:"encoding" yaml:"encoding"`
	OutputPath   string `json:"output_path" yaml:"output_path"`
	EnableCaller bool   `json:"enable_caller" yaml:"enable_caller"`

	// for file output
	MaxSize    int  `json:"max_size" yaml:"max_size"`
	MaxBackups int  `json:"max_backups" yaml:"max_backups"`
	MaxAge     int  `json:"max_age" yaml:"max_age"` // days
	Compress   bool `json:"compress" yaml:"compress"`
}

func FieldTime(key string, value time.Time) Field {
	return Field{Key: key, Value: value}
}
func FieldString(key, value string) Field {
	return Field{Key: key, Value: value}
}

func FieldInt(key string, value int) Field {
	return Field{Key: key, Value: value}
}

func FieldInt64(key string, value int64) Field {
	return Field{Key: key, Value: value}
}

func FieldBool(key string, value bool) Field {
	return Field{Key: key, Value: value}
}

func FieldError(err error) Field {
	return Field{Key: "error", Value: err}
}

func FieldAny(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}
