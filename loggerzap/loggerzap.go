// Copyright 2024 The OpenMail Authors
// SPDX-License-Identifier: Apache-2.0

// Package loggerzap is a thin wrapper around zap that implements Logger.
package loggerzap

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/linuxuser586/openmail/logger"
)

// Environment is used to provide reasonable configuration defaults
type Environment int

const (
	Development Environment = iota
	Production
)

// Level is the verbosity for logging
type Level int

const (
	Debug Level = iota
	Info
	Warn
	Error
	Panic
	Fatal
)

// Config provides options for creating a new Logger instance
type Config interface {
	GetEnvironment() Environment
	GetLevel() Level
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

func (z ZapLogger) Sync() {
	z.logger.Sync()
}

func (z ZapLogger) WithField(key string, val any) logger.Logger {
	z.fields = append(z.fields, zap.Any(key, val))
	return z
}
func New(config Config) logger.Logger {
	cfg := zap.NewDevelopmentEncoderConfig()
	level := zap.DebugLevel
	dev := true
	encoding := "console"
	var sampling *zap.SamplingConfig
	if config.GetEnvironment() == Production {
		cfg = zap.NewProductionEncoderConfig()
		level = zap.InfoLevel
		dev = false
		encoding = "json"
		sampling = &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		}
	}
	cfg.MessageKey = "short_message"

	c := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Development:      dev,
		Sampling:         sampling,
		Encoding:         encoding,
		EncoderConfig:    cfg,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := c.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic((err))
	}
	return &ZapLogger{logger: logger}
}
