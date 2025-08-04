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

type tracerCopyFrom struct {
	pool sync.Pool
}

func NewTracerCopyFrom() pgx.CopyFromTracer {
	return &tracerCopyFrom{
		pool: sync.Pool{
			New: func() any { return new(tracerCopyFromData) },
		},
	}
}

type tracerCopyFromData struct {
	Tables    []string
	Columns   []string
	StartedAt time.Time
}

func (s *tracerCopyFrom) TraceCopyFromStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceCopyFromStartData) context.Context {
	method := callerFunc(5)
	ctx, span := xtrace.Trace(ctx, method+" (copy_from)", trace.WithSpanKind(trace.SpanKindClient))
	span.SetAttributes(
		attribute.String("db.table", data.TableName.Sanitize()),
		attribute.String("repo.method", method),
	)

	traceData := s.pool.Get().(*tracerCopyFromData)
	traceData.Tables = data.TableName
	traceData.Columns = data.ColumnNames
	traceData.StartedAt = time.Now()
	return context.WithValue(ctx, s, traceData)
}

func (s *tracerCopyFrom) TraceCopyFromEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceCopyFromEndData) {
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

	traceData, _ := ctx.Value(s).(*tracerCopyFromData)
	duration := time.Since(traceData.StartedAt)
	xlog.ErrDebug(data.Err).
		Strs("tables", traceData.Tables).
		Any("columns", traceData.Columns).
		Dur("cost", duration).
		Str("command_tag", data.CommandTag.String()).
		Int64("rows", data.CommandTag.RowsAffected()).
		Uint32("pid", conn.PgConn().PID()).
		Msg("end")
}
