// Copyright 2024 The OpenMail Authors
// SPDX-License-Identifier: Apache-2.0

// Package loggerzap is a thin wrapper around [go.uber.org/zap] that implements [logger.Logger].
package loggerzap

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/linuxuser586/openmail/logger"
)

// Encoding enum
type Encoding string

const (
	Console Encoding = "console"
	Json    Encoding = "json"
)

func (e Encoding) String() string {
	return string(e)
}

// Config data
type Config struct {
	Env        logger.Environment
	Level      logger.Level
	Encoding   Encoding
	MessageKey string
}

func (c Config) GetEnvironment() logger.Environment {
	return c.Env
}

func (c Config) GetLevel() logger.Level {
	return c.Level
}

// NewDevConfig creates a reasonable development config
func NewDevConfig() Config {
	return Config{
		Env:        logger.Development,
		Level:      logger.Info,
		Encoding:   Console,
		MessageKey: "msg",
	}
}

// NewProdConfig creates a reasonable production config
func NewProdConfig() Config {
	return Config{
		Env:        logger.Production,
		Level:      logger.Info,
		Encoding:   Json,
		MessageKey: "short_message",
	}
}

func (c Config) toLevel() zapcore.Level {
	switch c.Level {
	case logger.Debug:
		return zapcore.DebugLevel
	case logger.Info:
		return zapcore.InfoLevel
	case logger.Warn:
		return zapcore.WarnLevel
	case logger.Error:
		return zapcore.ErrorLevel
	case logger.Panic:
		return zapcore.PanicLevel
	case logger.Fatal:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// ZapLogger is the Logger implementation using zap
type ZapLogger struct {
	logger *zap.Logger
	fields []zap.Field
}

func (z ZapLogger) Debug(msg string) {
	z.logger.Debug(msg, z.fields...)
}

func (z ZapLogger) Debugf(format string, a any) {
	z.logger.Debug(fmt.Sprintf(format, a), z.fields...)
}

func (z ZapLogger) Info(msg string) {
	z.logger.Info(msg, z.fields...)
}

func (z ZapLogger) Infof(format string, a any) {
	z.logger.Info(fmt.Sprintf(format, a), z.fields...)
}

func (z ZapLogger) Warn(msg string) {
	z.logger.Warn(msg, z.fields...)
}

func (z ZapLogger) Warnf(format string, a any) {
	z.logger.Warn(fmt.Sprintf(format, a), z.fields...)
}

func (z ZapLogger) Error(msg string) {
	z.logger.Error(msg, z.fields...)
}

func (z ZapLogger) Errorf(format string, a any) {
	z.logger.Error(fmt.Sprintf(format, a), z.fields...)
}

func (z ZapLogger) Panic(msg string) {
	z.logger.Panic(msg, z.fields...)
}

func (z ZapLogger) Panicf(format string, a any) {
	z.logger.Panic(fmt.Sprintf(format, a), z.fields...)
}

func (z ZapLogger) Fatal(msg string) {
	z.logger.Fatal(msg, z.fields...)
}

func (z ZapLogger) Fatalf(format string, a any) {
	z.logger.Fatal(fmt.Sprintf(format, a), z.fields...)
}

func (z ZapLogger) WithField(key string, val any) logger.Logger {
	z.fields = append(z.fields, zap.Any(key, val))
	return z
}

func (z ZapLogger) Sync() {
	z.logger.Sync()
}

func (z ZapLogger) GetLevel() logger.Level {
	l := z.logger.Level()
	switch l {
	case zapcore.DebugLevel:
		return logger.Debug
	case zapcore.InfoLevel:
		return logger.Info
	case zapcore.WarnLevel:
		return logger.Warn
	case zap.ErrorLevel:
		return logger.Error
	case zapcore.DPanicLevel, zap.PanicLevel:
		return logger.Panic
	case zapcore.FatalLevel:
		return logger.Fatal
	default:
		z.Panicf("unknown log level %v", l)
	}

	panic("does not reach")
}

// New creates a [go.uber.org/zap] instance
func New(config Config) logger.Logger {
	cfg := zap.NewDevelopmentEncoderConfig()
	var sampling *zap.SamplingConfig
	if config.Env == logger.Production {
		cfg = zap.NewProductionEncoderConfig()
		sampling = &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		}
	}
	cfg.MessageKey = config.MessageKey

	c := zap.Config{
		Level:            zap.NewAtomicLevelAt(config.toLevel()),
		Development:      config.Env == logger.Development,
		Sampling:         sampling,
		Encoding:         config.Encoding.String(),
		EncoderConfig:    cfg,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := c.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
	return &ZapLogger{logger: logger}
}
