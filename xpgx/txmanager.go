package xpgx

import (
	"context"
	"errors"
	"fmt"
	"github.com/xakepp35/pkg/xtrace"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type TxFunc func(tx pgx.Tx) error

// TxManager tx wrapper interface
type TxManager interface {
	// Do run tx with base pgx options
	Do(ctx context.Context, exec TxFunc) error
	// DoWithOptions run tx with users options
	DoWithOptions(ctx context.Context, opt pgx.TxOptions, exec TxFunc) error
}

type Transactional interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

type txManager struct {
	pool Transactional
}

// NewTxManager TxManager constructor
func NewTxManager(pool *pgxpool.Pool) TxManager {
	return &txManager{
		pool: pool,
	}
}

func (t txManager) Do(ctx context.Context, exec TxFunc) error {
	return t.DoWithOptions(ctx, pgx.TxOptions{}, exec)
}

func (t txManager) DoWithOptions(ctx context.Context, opt pgx.TxOptions, exec TxFunc) error {
	ctx, span := xtrace.Trace(ctx, "txManager.Transaction", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	span.SetAttributes(
		attribute.String("db.transaction.isolation_level", string(opt.IsoLevel)),
		attribute.Bool("db.transaction.read_only", opt.AccessMode == pgx.ReadOnly),
	)

	tx, err := t.pool.BeginTx(ctx, opt)
	if err != nil {
		span.RecordError(err)
		log.Error().Err(err).Msg("begin transaction")
		return fmt.Errorf("transaction begin: %w", err)
	}
	defer func() {
		if rerr := tx.Rollback(ctx); rerr != nil && !errors.Is(rerr, pgx.ErrTxClosed) {
			span.RecordError(rerr)
			log.Warn().Err(rerr).Msg("rollback transaction")
		}
	}()

	if err := exec(tx); err != nil {
		span.RecordError(err)
		log.Error().Err(err).Msg("exec transaction")
		return fmt.Errorf("transaction exec: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		span.RecordError(err)
		log.Error().Err(err).Msg("commit transaction")
		return fmt.Errorf("commit transaction: %w", err)
	}

	span.AddEvent("transaction committed")
	return nil
}
