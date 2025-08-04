package xtrace

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"github.com/xakepp35/pkg/xpgx"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"runtime"
	"strings"
)

const dbTracerName = "pgx-tracer"

type OTLPTracer struct{}

func NewOTLPTracer() *xpgx.Tracer {
	otlp := &OTLPTracer{}
	return &xpgx.Tracer{
		QueryTracer:    otlp,
		BatchTracer:    otlp,
		CopyFromTracer: otlp,
		PrepareTracer:  otlp,
		ConnectTracer:  otlp,
	}
}

// ==== вспомогательное: извлечение имени вызывающего метода ====
func callerFunc(skip int) string {
	pc := make([]uintptr, 10)
	n := runtime.Callers(skip, pc)
	if n == 0 {
		return "unknown"
	}
	frames := runtime.CallersFrames(pc[:n])
	for {
		frame, more := frames.Next()
		// фильтруем внутренние вызовы pgx
		if !strings.Contains(frame.Function, "pgx") && !strings.Contains(frame.Function, "pgconn") {
			return shortFuncName(frame.Function)
		}
		if !more {
			break
		}
	}
	return "unknown"
}

func shortFuncName(f string) string {
	parts := strings.Split(f, "/")
	return parts[len(parts)-1]
}

// ==== QueryTracer ====
func (t *OTLPTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	method := callerFunc(5)
	ctx, span := Trace(ctx, method, trace.WithSpanKind(trace.SpanKindClient))
	span.SetAttributes(
		attribute.String("db.system", "postgresql"),
		attribute.String("repo.method", method),
		attribute.String("db.user", conn.Config().User),
		attribute.String("db.host", conn.Config().Host),
		attribute.String("db.database", conn.Config().Database),
	)
	return ctx
}

func (t *OTLPTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return
	}
	if data.Err != nil {
		span.RecordError(data.Err)
	} else {
		span.AddEvent("DB query executed")
	}
	span.End()
}

// ==== BatchTracer ====
func (t *OTLPTracer) TraceBatchStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceBatchStartData) context.Context {
	method := callerFunc(5)
	ctx, span := Trace(ctx, method+" (batch)", trace.WithSpanKind(trace.SpanKindClient))
	span.SetAttributes(
		attribute.Int("db.batch_size", len(data.Batch.QueuedQueries)),
		attribute.String("repo.method", method),
	)
	return ctx
}

func (t *OTLPTracer) TraceBatchQuery(ctx context.Context, conn *pgx.Conn, data pgx.TraceBatchQueryData) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		span.AddEvent("Batch query")
	}
}

func (t *OTLPTracer) TraceBatchEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceBatchEndData) {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return
	}
	if data.Err != nil {
		span.RecordError(data.Err)
	} else {
		span.AddEvent("Batch completed")
	}
	span.End()
}

// ==== CopyFromTracer ====
func (t *OTLPTracer) TraceCopyFromStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceCopyFromStartData) context.Context {
	method := callerFunc(5)
	ctx, span := Trace(ctx, method+" (copy_from)", trace.WithSpanKind(trace.SpanKindClient))
	span.SetAttributes(
		attribute.String("db.table", data.TableName.Sanitize()),
		attribute.String("repo.method", method),
	)
	return ctx
}

func (t *OTLPTracer) TraceCopyFromEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceCopyFromEndData) {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return
	}
	if data.Err != nil {
		span.RecordError(data.Err)
	} else {
		span.SetAttributes(attribute.Int64("db.rows", data.CommandTag.RowsAffected()))
		span.AddEvent("CopyFrom completed")
	}
	span.End()
}

// ==== PrepareTracer ====
func (t *OTLPTracer) TracePrepareStart(ctx context.Context, conn *pgx.Conn, data pgx.TracePrepareStartData) context.Context {
	method := callerFunc(5)
	ctx, span := Trace(ctx, method+" (prepare)", trace.WithSpanKind(trace.SpanKindClient))
	span.SetAttributes(attribute.String("repo.method", method))
	return ctx
}

func (t *OTLPTracer) TracePrepareEnd(ctx context.Context, conn *pgx.Conn, data pgx.TracePrepareEndData) {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return
	}
	if data.Err != nil {
		span.RecordError(data.Err)
		log.Error().Err(data.Err).Msg("Prepare failed")
	} else {
		span.AddEvent("Prepare completed")
	}
	span.End()
}

// ==== ConnectTracer ====
func (t *OTLPTracer) TraceConnectStart(ctx context.Context, data pgx.TraceConnectStartData) context.Context {
	ctx, span := Trace(ctx, "db.connect", trace.WithSpanKind(trace.SpanKindClient))
	span.SetAttributes(
		attribute.String("db.host", data.ConnConfig.Host),
		attribute.String("db.user", data.ConnConfig.User),
		attribute.String("db.database", data.ConnConfig.Database),
	)
	return ctx
}

func (t *OTLPTracer) TraceConnectEnd(ctx context.Context, data pgx.TraceConnectEndData) {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return
	}
	if data.Err != nil {
		span.RecordError(data.Err)
	} else {
		span.AddEvent("Connect successful")
	}
	span.End()
}
