package xtrace

import (
	"context"
	"github.com/xakepp35/pkg/env"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"time"
)

const defaultTracerName = "xtrace"

var defaultTracerOptions = []otlptracegrpc.Option{
	otlptracegrpc.WithInsecure(),
}

func Trace(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return otel.Tracer(defaultTracerName).Start(ctx, name, opts...)
}

func InitTracer(serviceName string, opts ...otlptracegrpc.Option) error {
	if len(opts) == 0 {
		opts = defaultTracerOptions
	}

	endpoint := env.Get("OTEL_DEFAULT_ENDPOINT", "localhost:4317")

	sampler := sdktrace.ParentBased(sdktrace.TraceIDRatioBased(env.Float64("OTEL_TRACES_SAMPLER_RATIO", 1)))

	opts = append(opts, otlptracegrpc.WithEndpoint(endpoint))
	exporter, err := otlptracegrpc.New(context.Background(), opts...)
	if err != nil {
		return err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sampler),
		sdktrace.WithBatcher(exporter, sdktrace.WithBatchTimeout(time.Second)),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			xTextMapPropagator{},
			propagation.TraceContext{},
		),
	)

	otel.SetTracerProvider(tp)

	return nil
}
