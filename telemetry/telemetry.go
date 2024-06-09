// Copyright 2024 The OpenMail Authors
// SPDX-License-Identifier: Apache-2.0

// package telemetry provides traces, metrics, and logs
package telemetry

import (
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	"github.com/linuxuser586/openmail/logger"
)

type Telemetry interface {
	Tracer() trace.Tracer
	Meter() metric.Meter
	Logger() logger.Logger
}
