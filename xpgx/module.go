package xpgx

import (
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xakepp35/pkg/env"
	"go.uber.org/fx"
)

func NewModule(poolName string) fx.Option {
	dsnEnv := "PG_DSN"
	if poolName != "" {
		dsnEnv += strings.ToUpper(poolName)
	}
	tagDsn := `name:"pg_dsn_` + poolName + `"`
	tagCfg := `name:"pg_cfg_` + poolName + `"`
	tagPool := `name:"pg_pool_` + poolName + `"`
	return fx.Module("xpgx-"+poolName,
		fx.Provide(
			fx.Annotate(env.String(dsnEnv, ""), fx.ResultTags(tagDsn)),
			fx.Annotate(pgxpool.ParseConfig, fx.ParamTags(tagDsn), fx.ResultTags(tagCfg)),
			NewTracer,
		),
		fx.Invoke(
			RegisterTracer,
		),
		fx.Provide(
			fx.Annotate(NewPool, fx.ParamTags(tagCfg), fx.ResultTags(tagPool)),
			fx.Annotate(NewLifecycle, fx.ParamTags(tagPool)),
		),
		fx.Invoke(RegisterLifecycle),
	)
}

// var Module = fx.Module("xpgx",
// 	fx.Provide(
// 		fx.Annotate(env.String("PG_DSN", ""), fx.ResultTags(`name:"pg_dsn"`)),
// 		fx.Annotate(pgxpool.ParseConfig, fx.ParamTags(`name:"pg_dsn"`), fx.ResultTags(`name:"pg_cfg"`)),
// 		NewTracer,

// 		fx.Annotate(NewPool, fx.ParamTags(`name:"pg_cfg"`), fx.ResultTags(`name:"pg_pool"`)),
// 		fx.Annotate(NewLifecycle, fx.ParamTags(`name:"pg_pool"`)),
// 	),
// 	fx.Invoke(RegisterLifecycle),
// )
