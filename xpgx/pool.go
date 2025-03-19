package xpgx

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/xakepp35/pkg/xlog"
	"go.uber.org/fx"
)

func NewPool(lc fx.Lifecycle, cfg *pgxpool.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	xlog.ErrInfo(err).
		Bool("starting", true).
		Msg("pgxpool.NewWithConfig")
	if err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Info().
				Bool("starting", true).
				Msg("pgxpool.Close")
			pool.Close()
			log.Info().
				Bool("starting", false).
				Msg("pgxpool.Close")
			return nil
		},
	})
	return pool, nil
}
