package xpgx

import (
	"context"
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

func (s *tracerCopyFrom) TraceCopyFromStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceCopyFromStartData) context.Context {
	traceData := s.pool.Get().(*tracerCopyFromData)
	traceData.Tables = []string(data.TableName)
	traceData.Columns = data.ColumnNames
	traceData.StartedAt = time.Now()
	return context.WithValue(ctx, s, traceData)
}

func (s *tracerCopyFrom) TraceCopyFromEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceCopyFromEndData) {
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
