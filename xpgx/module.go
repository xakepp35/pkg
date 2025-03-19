package xpgx

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xakepp35/pkg/env"
	"go.uber.org/fx"
)

var Module = fx.Module("xpgx",
	fx.Provide(
		fx.Annotate(env.String("PG_DSN", ""), fx.ResultTags(`name:"pg_dsn"`)),
		fx.Annotate(pgxpool.ParseConfig, fx.ParamTags(`name:"pg_dsn"`), fx.ResultTags(`name:"pg_cfg"`)),
		fx.Annotate(NewPool, fx.ParamTags(`name:"pg_cfg"`)),
	),
)
