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

// QueryTracer traces Query, QueryRow, and Exec.
type tracerQuery struct {
	pool sync.Pool
}

// QueryTracer traces Query, QueryRow, and Exec.
func NewTracerQuery() pgx.QueryTracer {
	return &tracerQuery{
		pool: sync.Pool{
			New: func() any { return new(tracerQueryData) },
		},
	}
}

type tracerQueryData struct {
	Sql       string
	Args      []any
	StartedAt time.Time
}

func (s *tracerQuery) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	method := callerFunc(5)
	ctx, span := xtrace.Trace(ctx, method, trace.WithSpanKind(trace.SpanKindClient))
	span.SetAttributes(
		attribute.String("db.system", "postgresql"),
		attribute.String("repo.method", method),
		attribute.String("db.user", conn.Config().User),
		attribute.String("db.host", conn.Config().Host),
		attribute.String("db.database", conn.Config().Database),
	)

	traceData := s.pool.Get().(*tracerQueryData)
	traceData.Sql = data.SQL
	traceData.Args = data.Args
	traceData.StartedAt = time.Now()
	return context.WithValue(ctx, s, traceData)
}

func (s *tracerQuery) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
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

	traceData, _ := ctx.Value(s).(*tracerQueryData)
	duration := time.Since(traceData.StartedAt)
	xlog.ErrDebug(data.Err).
		Str("sql", traceData.Sql).
		Any("args", traceData.Args).
		Dur("cost", duration).
		Str("tag", data.CommandTag.String()).
		Int64("rows", data.CommandTag.RowsAffected()).
		Uint32("pid", conn.PgConn().PID()).
		Msg("end")
}
