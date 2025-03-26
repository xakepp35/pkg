package xpgx

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/xakepp35/pkg/xlog"
	"go.uber.org/fx"
)

// NewPool создаёт новый пул соединений с базой данных
func NewPool(cfg *pgxpool.Config) (*pgxpool.Pool, error) {
	NewTracer().Apply(cfg.ConnConfig)
	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	xlog.ErrInfo(err).
		Msg("pgxpool.NewWithConfig")
	if err != nil {
		return nil, err
	}
	return pool, nil
}

// RegisterLifecycle
func RegisterLifecycle(lc fx.Lifecycle, plc *Lifecycle) {
	plc.Register(lc)
}

type Lifecycle struct {
	pool *pgxpool.Pool
}

func NewLifecycle(pool *pgxpool.Pool) *Lifecycle {
	return &Lifecycle{
		pool: pool,
	}
}

func (s *Lifecycle) Register(lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: s.OnStart,
		OnStop:  s.OnStop,
	})
}

// onStart проверяет соединение при запуске приложения
func (s *Lifecycle) OnStart(ctx context.Context) error {
	err := s.pool.Ping(ctx)
	xlog.ErrInfo(err).
		Msg("pgxpool.Ping")
	return err
}

// OnStop closes connection pool correctly
func (s *Lifecycle) OnStop(ctx context.Context) error {
	s.pool.Close()
	log.Info().
		Msg("pgxpool.Close")
	return nil
}
