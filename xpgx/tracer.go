package xpgx

import (
	"github.com/jackc/pgx/v5"
)

type tracer struct {
	pgx.QueryTracer
	pgx.BatchTracer
	pgx.CopyFromTracer
	pgx.PrepareTracer
	pgx.ConnectTracer
}

func NewTracer() *tracer {
	return &tracer{
		QueryTracer:    NewTracerQuery(),
		BatchTracer:    NewTracerBatch(),
		CopyFromTracer: NewTracerCopyFrom(),
		PrepareTracer:  NewTracerPrepare(),
		ConnectTracer:  NewTracerConnect(),
	}
}

func (s *tracer) Apply(connConfig *pgx.ConnConfig) {
	connConfig.Tracer = s
}
