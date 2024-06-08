// Copyright The OpenMail Authors
// SPDX-License-Identifier: Apache-2.0

// Package logger provides the logger API.
package logger

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

// Logger is the API for logging
type Logger interface {
	// Debug logs a message at the most verbose level
	Debug(msg string)
	// Debugf is similar to combining Debug and Printf
	Debugf(format string, a any)
	// Info logs a message one verbose level higher than Debug
	Info(msg string)
	// Infof is similar to combining Info and Printf
	Infof(format string, a any)
	// Warn logs a message one verbose level higher than Info
	Warn(msg string)
	// Warnf is similar to combining Warn and Printf
	Warnf(format string, a any)
	// Error logs a message one verbose level higher than Warn
	Error(msg string)
	// Errorf is similar to combining Error and Printf
	Errorf(format string, a any)
	// Panic logs a message one verbose level higher than Error and panics
	Panic(msg string)
	// Panicf is similar to combining Panic and Printf
	Panicf(format string, a any)
	// Fatal logs a message one verbose level higher than Panic and
	// calls os.Exit(1) after logging the message
	Fatal(msg string)
	// Fatalf is similar to combining Panic and Printf
	Fatalf(format string, a any)
	// WithField adds key value pairs to the log entry
	WithField(key string, val any) Logger
	// Sync flushes the buffer
	Sync()
	// GetLevel gets the current [Level]
	GetLevel() Level
}
