package xtrace

import (
	"context"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const TraceHeader = "x-trace-id"

type xTextMapPropagator struct{}

func (x xTextMapPropagator) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	spanCtx := trace.SpanContextFromContext(ctx)
	if !spanCtx.IsValid() {
		return
	}
	// Устанавливаем trace-id в заголовок
	carrier.Set(TraceHeader, spanCtx.TraceID().String())
}

func (x xTextMapPropagator) Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	traceIDHex := carrier.Get(TraceHeader)
	if traceIDHex == "" {
		return ctx
	}

	traceID, err := trace.TraceIDFromHex(traceIDHex)
	if err != nil || !traceID.IsValid() {
		return ctx
	}

	// Создаём новый SpanContext только с TraceID (SpanID пустой, т.к. он не передается)
	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     trace.SpanID{}, // новый span создаст свой spanID
		TraceFlags: trace.FlagsSampled,
		Remote:     true,
	})
	return trace.ContextWithSpanContext(ctx, spanCtx)
}

func (x xTextMapPropagator) Fields() []string {
	return []string{TraceHeader}
}
