package xpgx

import (
	"context"
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
	traceData := s.pool.Get().(*tracerQueryData)
	traceData.Sql = data.SQL
	traceData.Args = data.Args
	traceData.StartedAt = time.Now()
	return context.WithValue(ctx, s, traceData)
}

func (s *tracerQuery) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
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
