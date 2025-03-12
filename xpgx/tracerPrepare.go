package xpgx

import (
	"context"
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

func (s *tracerPrepare) TracePrepareStart(ctx context.Context, conn *pgx.Conn, data pgx.TracePrepareStartData) context.Context {
	traceData := s.pool.Get().(*tracerPrepareData)
	traceData.Name = data.Name
	traceData.SQL = data.SQL
	traceData.StartedAt = time.Now()
	return context.WithValue(ctx, s, traceData)
}

func (s *tracerPrepare) TracePrepareEnd(ctx context.Context, conn *pgx.Conn, data pgx.TracePrepareEndData) {
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
