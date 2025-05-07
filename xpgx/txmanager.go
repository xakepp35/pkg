package xpgx

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
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
func NewTxManager(pool Transactional) TxManager {
	return &txManager{
		pool: pool,
	}
}

func (t txManager) Do(ctx context.Context, exec TxFunc) error {
	return t.DoWithOptions(ctx, pgx.TxOptions{}, exec)
}

func (t txManager) DoWithOptions(ctx context.Context, opt pgx.TxOptions, exec TxFunc) error {
	tx, err := t.pool.BeginTx(ctx, opt)
	if err != nil {
		log.Error().Err(err).Msg("begin transaction")
		return fmt.Errorf("transaction begin: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			log.Warn().Err(err).Msg("rollback transaction")
		}
	}()

	if err := exec(tx); err != nil {
		log.Error().Err(err).Msg("exec transaction")
		return fmt.Errorf("transaction exec: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error().Err(err).Msg("commit transaction")
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
