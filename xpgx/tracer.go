package xpgx

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Tracer struct {
	pgx.QueryTracer
	pgx.BatchTracer
	pgx.CopyFromTracer
	pgx.PrepareTracer
	pgx.ConnectTracer
}

func NewTracer() *Tracer {
	return &Tracer{
		QueryTracer:    NewTracerQuery(),
		BatchTracer:    NewTracerBatch(),
		CopyFromTracer: NewTracerCopyFrom(),
		PrepareTracer:  NewTracerPrepare(),
		ConnectTracer:  NewTracerConnect(),
	}
}

func RegisterTracer(tr *Tracer, cfg *pgxpool.Config) {
	tr.Apply(cfg.ConnConfig)
}

func (s *Tracer) Apply(connConfig *pgx.ConnConfig) {
	connConfig.Tracer = s
}
