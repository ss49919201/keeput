package appotel

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
)

var shutdownTraceProvider func(context.Context) error
var errInitOtelProvider error
var initOtelProviderOnce sync.Once

func InitTraceProvider(ctx context.Context) (func(context.Context) error, error) {
	initOtelProviderOnce.Do(func() {
		exporter, err := otlptracehttp.New(ctx)
		if err != nil {
			errInitOtelProvider = err
			return
		}
		resource, err := resource.New(ctx)
		if err != nil {
			errInitOtelProvider = err
			return
		}
		tp := sdkTrace.NewTracerProvider(
			sdkTrace.WithBatcher(exporter),
			sdkTrace.WithResource(resource),
		)
		otel.SetTracerProvider(tp)
		otel.SetTextMapPropagator(propagation.TraceContext{})
		shutdownTraceProvider = tp.Shutdown
	})
	return shutdownTraceProvider, errInitOtelProvider
}
