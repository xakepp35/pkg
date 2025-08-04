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

// ConnectTracer traces Connect and ConnectConfig.
type tracerConnect struct {
	pool sync.Pool
}

// ConnectTracer traces Connect and ConnectConfig.
func NewTracerConnect() pgx.ConnectTracer {
	return &tracerConnect{
		pool: sync.Pool{
			New: func() any { return new(tracerConnectData) },
		},
	}
}

type tracerConnectData struct {
	ConnConfig *pgx.ConnConfig
	StartedAt  time.Time
}

func (s *tracerConnect) TraceConnectStart(ctx context.Context, data pgx.TraceConnectStartData) context.Context {
	ctx, span := xtrace.Trace(ctx, "db.connect", trace.WithSpanKind(trace.SpanKindClient))
	span.SetAttributes(
		attribute.String("db.host", data.ConnConfig.Host),
		attribute.String("db.user", data.ConnConfig.User),
		attribute.String("db.database", data.ConnConfig.Database),
	)

	traceData := s.pool.Get().(*tracerConnectData)
	traceData.ConnConfig = data.ConnConfig
	traceData.StartedAt = time.Now()
	return context.WithValue(ctx, s, traceData)
}

func (s *tracerConnect) TraceConnectEnd(ctx context.Context, data pgx.TraceConnectEndData) {
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
	traceData, _ := ctx.Value(s).(*tracerConnectData)
	duration := time.Since(traceData.StartedAt)
	xlog.ErrDebug(data.Err).
		Str("dsn", traceData.ConnConfig.Host).
		Dur("cost", duration).
		Msg("end")
}
