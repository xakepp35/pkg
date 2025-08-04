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

type tracerPrepare struct {
	pool sync.Pool
}

func NewTracerPrepare() pgx.PrepareTracer {
	return &tracerPrepare{
		pool: sync.Pool{
			New: func() any { return new(tracerPrepareData) },
		},
	}
}

type tracerPrepareData struct {
	Name      string
	SQL       string
	StartedAt time.Time
}

func (s *tracerPrepare) TracePrepareStart(ctx context.Context, _ *pgx.Conn, data pgx.TracePrepareStartData) context.Context {
	method := callerFunc(5)
	ctx, span := xtrace.Trace(ctx, method+" (prepare)", trace.WithSpanKind(trace.SpanKindClient))
	span.SetAttributes(attribute.String("repo.method", method))

	traceData := s.pool.Get().(*tracerPrepareData)
	traceData.Name = data.Name
	traceData.SQL = data.SQL
	traceData.StartedAt = time.Now()
	return context.WithValue(ctx, s, traceData)
}

func (s *tracerPrepare) TracePrepareEnd(ctx context.Context, conn *pgx.Conn, data pgx.TracePrepareEndData) {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return
	}
	if data.Err != nil {
		span.RecordError(data.Err)
	} else {
		span.AddEvent("Prepare completed")
	}

	traceData, _ := ctx.Value(s).(*tracerPrepareData)
	duration := time.Since(traceData.StartedAt)
	xlog.ErrDebug(data.Err).
		Str("name", traceData.Name).
		Str("sql", traceData.SQL).
		Dur("cost", duration).
		Bool("already_prepared", data.AlreadyPrepared).
		Uint32("pid", conn.PgConn().PID()).
		Msg("end")
}
