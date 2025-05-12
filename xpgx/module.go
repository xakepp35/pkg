package xpgx

import (
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xakepp35/pkg/env"
	"go.uber.org/fx"
)

func NewModule(poolName string) fx.Option {
	dsnKey := genDsnKey(poolName)
	dsnEnv := env.String(dsnKey, DefaultDSN)
	tagDsn := `name:"pg_dsn_` + poolName + `"`
	tagCfg := `name:"pg_cfg_` + poolName + `"`
	tagPool := genPoolTag(poolName)
	tagTrc := `name:"pg_trc_` + poolName + `"`
	tagPlc := `name:"pg_plc_` + poolName + `"`
	tagTxm := genTxmTag(poolName)
	return fx.Module("xpgx-"+poolName,
		fx.Supply(
			fx.Annotate(dsnEnv, fx.ResultTags(tagDsn)),
		),
		fx.Provide(
			fx.Annotate(pgxpool.ParseConfig, fx.ParamTags(tagDsn), fx.ResultTags(tagCfg)),
			fx.Annotate(NewTracer, fx.ResultTags(tagTrc)),
		),
		fx.Invoke(
			fx.Annotate(RegisterTracer, fx.ParamTags(tagTrc, tagCfg)),
		),
		fx.Provide(
			fx.Annotate(NewPool, fx.ParamTags(tagCfg), fx.ResultTags(tagPool)),
			fx.Annotate(NewLifecycle, fx.ParamTags(tagPool), fx.ResultTags(tagPlc)),
		),
		fx.Provide(
			fx.Annotate(NewTxManager, fx.ParamTags(tagPool), fx.ResultTags(tagTxm)),
		),
		fx.Invoke(
			fx.Annotate(RegisterLifecycle, fx.ParamTags("", tagPlc)),
		),
	)
}

func genDsnKey(poolName string) string {
	if poolName == "" {
		return DefaultDSNKey
	}
	return DefaultDSNKey + "_" + strings.ToUpper(poolName)
}

func genPoolTag(name string) string {
	if name == "" {
		return ""
	}
	return `name:"pg_pool_` + name + `"`
}

func genTxmTag(name string) string {
	if name == "" {
		return ""
	}
	return `name:"pg_txm_` + name + `"`
}
