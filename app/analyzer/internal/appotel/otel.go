package appotel

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type typeShutdownProvider = func(context.Context) error
type typeFlushMetrics = func(context.Context) error

var (
	// TraceProvider
	shutdownTraceProvider typeShutdownProvider
	errInitTraceProvider  error
	initTraceProviderOnce sync.Once

	// MeterProvider
	shutdownMeterProvider typeShutdownProvider
	errInitMeterProvider  error
	initMeterProviderOnce sync.Once
	flushMetrics          typeFlushMetrics
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
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resource),
		)
		otel.SetTracerProvider(tp)
		otel.SetTextMapPropagator(propagation.TraceContext{})
		shutdownTraceProvider = tp.Shutdown
	})
	return shutdownTraceProvider, errInitTraceProvider
}

// NOTE: MetricReader は ManualReader なので FlushMetrics を実行してオンデマンドに収集、送出が必要
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
		reader := metric.NewManualReader()
		mp := metric.NewMeterProvider(metric.WithReader(reader))
		otel.SetMeterProvider(mp)
		shutdownMeterProvider = mp.Shutdown
		flushMetrics = func(ctx context.Context) error {
			var resource metricdata.ResourceMetrics
			if err := reader.Collect(ctx, &resource); err != nil {
				return err
			}
			return exporter.Export(ctx, &resource)
		}
	})
	return shutdownMeterProvider, errInitMeterProvider
}

func FlushMetrics(ctx context.Context) error {
	_, err := InitMeterProvider(ctx)
	if err != nil {
		return err
	}
	return flushMetrics(ctx)
}

func RecordError(ctx context.Context, err error) {
	trace.SpanFromContext(ctx).RecordError(err, trace.WithStackTrace(true))
}
