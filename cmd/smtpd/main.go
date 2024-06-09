// Copyright 2024 The OpenMail Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"log"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	"github.com/linuxuser586/openmail/logger"
	"github.com/linuxuser586/openmail/loggerzap"
	"github.com/linuxuser586/openmail/smtpd"
)

const (
	defaultTimeout = 5 * 60
	name           = "smtpd"
)

func main() {
	s := smtpd.New(conf{}, tel{}, repo{})
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

// TODO create configuration
type conf struct{}

func (conf) Network() smtpd.Network {
	return smtpd.Tcp
}

func (conf) Host() string {
	return "localhost"
}

func (conf) Port() int {
	return 2525
}

func (conf) InitialTimeout() int {
	return defaultTimeout
}

func (conf) MailCmdTimeout() int {
	return defaultTimeout
}

func (conf) RecipientCmdTimeout() int {
	return defaultTimeout
}

func (conf) DataInitTimeout() int {
	return defaultTimeout
}

func (conf) DataBlockTimeout() int {
	return defaultTimeout
}

func (conf) DataTerminationTimeout() int {
	return defaultTimeout
}

func (conf) HostName() string {
	return "smtp.example.com"
}

type tel struct{}

func (tel) Tracer() trace.Tracer {
	return otel.Tracer(name)
}

func (tel) Meter() metric.Meter {
	return otel.Meter(name)
}

func (tel) Logger() logger.Logger {
	return loggerzap.New(loggerzap.NewDevConfig())
}

type repo struct{}
