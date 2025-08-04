package xpgx

import (
	"context"
	"github.com/xakepp35/pkg/xtrace"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/xakepp35/pkg/xlog"
)

type tracerBatch struct {
	pool sync.Pool
}

func NewTracerBatch() pgx.BatchTracer {
	return &tracerBatch{
		pool: sync.Pool{
			New: func() any { return new(tracerBatchData) },
		},
	}
}

type tracerBatchData struct {
	Batch     *pgx.Batch
	StartedAt time.Time
}

func (s *tracerBatch) TraceBatchStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceBatchStartData) context.Context {
	traceData := s.pool.Get().(*tracerBatchData)
	traceData.Batch = data.Batch
	traceData.StartedAt = time.Now()

	method := callerFunc(5)
	ctx, span := xtrace.Trace(ctx, method+" (batch)", trace.WithSpanKind(trace.SpanKindClient))
	span.SetAttributes(
		attribute.Int("db.batch_size", len(data.Batch.QueuedQueries)),
		attribute.String("repo.method", method),
	)

	return context.WithValue(ctx, s, traceData)
}

func (s *tracerBatch) TraceBatchQuery(ctx context.Context, conn *pgx.Conn, data pgx.TraceBatchQueryData) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		span.AddEvent("Batch query")
	}

	traceData, _ := ctx.Value(s).(*tracerBatchData)
	duration := time.Since(traceData.StartedAt)
	xlog.ErrDebug(data.Err).
		Str("sql", data.SQL).
		Any("args", data.Args).
		Str("command_tag", data.CommandTag.String()).
		Dur("cost", duration).
		Uint32("pid", conn.PgConn().PID()).
		Msg("query")
}

func (s *tracerBatch) TraceBatchEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceBatchEndData) {
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

	traceData, _ := ctx.Value(s).(*tracerBatchData)
	duration := time.Since(traceData.StartedAt)
	xlog.ErrDebug(data.Err).
		Int("len", traceData.Batch.Len()).
		Dur("cost", duration).
		Uint32("pid", conn.PgConn().PID()).
		Msg("end")
}
