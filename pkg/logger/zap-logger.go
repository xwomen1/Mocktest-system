package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type ZapLogger struct {
	logger *zap.Logger
	level  Level
	config Config
}

func NewZapLogger(config Config) (*ZapLogger, error) {

	zapLevel := toZapLevel(config.Level)

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var encoder zapcore.Encoder
	if config.Encoding == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	var writeSyncer zapcore.WriteSyncer
	if config.OutputPath == "" || config.OutputPath == "stdout" {
		writeSyncer = zapcore.AddSync(os.Stdout)
	} else if config.OutputPath == "stderr" {
		writeSyncer = zapcore.AddSync(os.Stderr)
	} else {

		lumberjackLogger := &lumberjack.Logger{
			Filename:   config.OutputPath,
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   config.Compress,
		}
		writeSyncer = zapcore.AddSync(lumberjackLogger)
	}

	core := zapcore.NewCore(encoder, writeSyncer, zapLevel)

	options := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	}

	if config.EnableCaller {
		options = append(options, zap.AddCaller())
	}

	logger := zap.New(core, options...)

	return &ZapLogger{
		logger: logger,
		level:  config.Level,
		config: config,
	}, nil
}

func (z *ZapLogger) Debug(msg string, fields ...Field) {
	z.logger.Debug(msg, toZapFields(fields)...)
}

func (z *ZapLogger) Info(msg string, fields ...Field) {
	z.logger.Info(msg, toZapFields(fields)...)
}

func (z *ZapLogger) Warn(msg string, fields ...Field) {
	z.logger.Warn(msg, toZapFields(fields)...)
}

func (z *ZapLogger) Error(msg string, fields ...Field) {
	z.logger.Error(msg, toZapFields(fields)...)
}

func (z *ZapLogger) Fatal(msg string, fields ...Field) {
	z.logger.Fatal(msg, toZapFields(fields)...)
}

func (z *ZapLogger) Panic(msg string, fields ...Field) {
	z.logger.Panic(msg, toZapFields(fields)...)
}

func (z *ZapLogger) WithContext(ctx context.Context) Logger {

	return z
}

func (z *ZapLogger) With(fields ...Field) Logger {
	newLogger := z.logger.With(toZapFields(fields)...)
	return &ZapLogger{
		logger: newLogger,
		level:  z.level,
		config: z.config,
	}
}

func (z *ZapLogger) WithPrefix(prefix string) Logger {
	return z.With(FieldString("prefix", prefix))
}

func (z *ZapLogger) Sync() error {
	return z.logger.Sync()
}

func (z *ZapLogger) SetLevel(level Level) {
	z.level = level

}

func (z *ZapLogger) GetLevel() Level {
	return z.level
}

func toZapLevel(level Level) zapcore.Level {
	switch level {
	case DebugLevel:
		return zapcore.DebugLevel
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	case FatalLevel:
		return zapcore.FatalLevel
	case PanicLevel:
		return zapcore.PanicLevel
	default:
		return zapcore.InfoLevel
	}
}

func toZapFields(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = zap.Any(f.Key, f.Value)
	}
	return zapFields
}

// creates a default Zap logger
func DefaultZapLogger() (*ZapLogger, error) {
	config := Config{
		Level:        InfoLevel,
		Encoding:     "json",
		OutputPath:   "stdout",
		EnableCaller: true,
	}
	return NewZapLogger(config)
}
