package appotel

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type typeShutdownProvider = func(context.Context) error

var (
	shutdownTraceProvider typeShutdownProvider
	shutdownMeterProvider typeShutdownProvider
	errInitTraceProvider  error
	initTraceProviderOnce sync.Once
	errInitMeterProvider  error
	initMeterProviderOnce sync.Once
)

func InitTraceProvider(ctx context.Context) (typeShutdownProvider, error) {
	initTraceProviderOnce.Do(func() {
		exporter, err := otlptracehttp.New(
			ctx,
			otlptracehttp.WithInsecure(),
		)
		if err != nil {
			errInitTraceProvider = err
			return
		}
		resource, err := resource.New(ctx)
		if err != nil {
			errInitTraceProvider = err
			return
		}
		tp := sdktrace.NewTracerProvider(
			sdktrace.WithSyncer(exporter),
			sdktrace.WithResource(resource),
		)
		otel.SetTracerProvider(tp)
		otel.SetTextMapPropagator(propagation.TraceContext{})
		shutdownTraceProvider = tp.Shutdown
	})
	return shutdownTraceProvider, errInitTraceProvider
}

func InitMeterProvider(ctx context.Context) (typeShutdownProvider, error) {
	initMeterProviderOnce.Do(func() {
		exporter, err := otlpmetrichttp.New(
			ctx,
			otlpmetrichttp.WithInsecure(),
		)
		if err != nil {
			errInitMeterProvider = err
			return
		}
		mp := metric.NewMeterProvider(metric.WithReader(metric.NewPeriodicReader(exporter)))
		otel.SetMeterProvider(mp)
		shutdownMeterProvider = mp.Shutdown
	})
	return shutdownMeterProvider, errInitMeterProvider
}

func RecordError(ctx context.Context, err error) {
	trace.SpanFromContext(ctx).RecordError(err, trace.WithStackTrace(true))
}
