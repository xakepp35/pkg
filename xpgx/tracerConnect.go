package xpgx

import (
	"context"
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
	traceData := s.pool.Get().(*tracerConnectData)
	traceData.ConnConfig = data.ConnConfig
	traceData.StartedAt = time.Now()
	return context.WithValue(ctx, s, traceData)
}

func (s *tracerConnect) TraceConnectEnd(ctx context.Context, data pgx.TraceConnectEndData) {
	traceData, _ := ctx.Value(s).(*tracerConnectData)
	duration := time.Since(traceData.StartedAt)
	xlog.ErrDebug(data.Err).
		Str("dsn", traceData.ConnConfig.Host).
		Dur("cost", duration).
		Msg("end")
}
