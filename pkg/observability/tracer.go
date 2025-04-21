package observability

import (
	"context"
	"go.opentelemetry.io/otel"
)

import "go.opentelemetry.io/otel/trace"

var tracer = otel.Tracer("")

func StartSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	return tracer.Start(ctx, name)
}
