package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
	"path/filepath"
)

var (
	ErrorLoggerConfigNotValidated = fmt.Errorf("logger config not validated")
	ll                            *zap.Logger
)

type (
	ZapLogger struct{}
)

func NewLoggerDefault() (*zap.Logger, error) {
	ll, err := initLoggerDefault()
	if err != nil {
		return nil, err
	}

	return ll, nil
}

func NewLogger(cfg *Config) (*zap.Logger, error) {
	ll, err := initLogger(cfg)
	if err != nil {
		return nil, err
	}
	return ll, nil
}

func initLoggerDefault() (*zap.Logger, error) {
	cfg := NewDefaultConfig()
	return initLogger(cfg)
}

func initLogger(cfg *Config) (*zap.Logger, error) {
	defer func() {
		cfg = nil // clean global config instance
	}()

	writeSyncer := zapcore.AddSync(os.Stdout)

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	if !cfg.Stdout {
		if !cfg.validated {
			return nil, ErrorLoggerConfigNotValidated
		}

		if _, err := os.Stat(cfg.LogDir); os.IsNotExist(err) {
			if err := os.MkdirAll(cfg.LogDir, 0755); err != nil {
				log.Printf("Failed to create log directory: %s", cfg.LogDir)
				os.Exit(1)
			}
		}

		writeSyncer = zapcore.AddSync(&lumberjack.Logger{
			Filename:   filepath.Join(cfg.LogDir, cfg.Filename),
			MaxSize:    cfg.MaxFileSize,
			MaxBackups: cfg.MaxBackups,
			LocalTime:  cfg.LocalTime,
			Compress:   cfg.Compress,
			MaxAge:     cfg.MaxAge,
		})
	}
	log.Printf("using logger to stdout: %v\n", cfg.Stdout)
	return zap.New(zapcore.NewCore(zapcore.NewConsoleEncoder(encoderCfg), writeSyncer, zapcore.InfoLevel)), nil
}

func Info(format string, args ...any) {
	ll.Info(fmt.Sprintf(format, args...))
}

func Error(format string, args ...any) {
	ll.Error(fmt.Sprintf(format, args...))
}

func Debug(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	ll.Debug(msg)
	fmt.Println(msg)
}

func Warn(format string, args ...any) {
	ll.Warn(fmt.Sprintf(format, args...))
}
