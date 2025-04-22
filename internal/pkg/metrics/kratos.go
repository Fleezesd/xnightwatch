package metrics

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	Name = "xmetrics"

	meter = otel.Meter(Name)

	AttributeKind      = attribute.Key("kind")
	AttributeOperation = attribute.Key("operation")
	AttributeCode      = attribute.Key("code")
	AttributeReason    = attribute.Key("reason")

	// OpenTelemetry metrics
	KratosServerMetricRequests metric.Int64Counter
	KratosMetricSeconds        metric.Float64Histogram
)

func init() {
	var err error

	KratosServerMetricRequests, err = meter.Int64Counter(
		"requests_total",
		metric.WithDescription("The total number of processed requests"),
		metric.WithUnit("1"),
	)
	if err != nil {
		panic(err)
	}

	KratosMetricSeconds, err = meter.Float64Histogram(
		"request_duration_ms",
		metric.WithDescription("server requests duration(ms)"),
		metric.WithUnit("ms"),
		metric.WithExplicitBucketBoundaries(5, 10, 25, 50, 100, 250, 500, 1000),
	)
	if err != nil {
		panic(err)
	}
}

func RecordRequestDuration(ctx context.Context, duration float64, kind, operation string) {
	KratosMetricSeconds.Record(ctx, duration,
		metric.WithAttributes(
			AttributeKind.String(kind),
			AttributeOperation.String(operation),
		),
	)
}

func RecordRequest(ctx context.Context, kind, operation, code, reason string) {
	KratosServerMetricRequests.Add(ctx, 1,
		metric.WithAttributes(
			AttributeKind.String(kind),
			AttributeOperation.String(operation),
			AttributeCode.String(code),
			AttributeReason.String(reason),
		),
	)
}
