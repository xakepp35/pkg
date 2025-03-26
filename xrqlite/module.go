package xrqlite

import (
	"strings"

	"github.com/xakepp35/pkg/env"
	"go.uber.org/fx"
)

// Module provides the RQLite dependencies to Uber FX
func NewModule(name string) fx.Option {
	dsnKey := genConnURL(name)
	dsnEnv := env.String(dsnKey, DefaultConnURL)
	tagDsn := `name:"rq_dsn_` + name + `"`
	tagConn := genConnTag(name)
	tagRqlc := `name:"rq_lc_` + name + `"`
	return fx.Module("rqlite",
		fx.Supply(
			fx.Annotate(dsnEnv, fx.ResultTags(tagDsn)),
		),
		fx.Provide(
			fx.Annotate(NewConnection, fx.ParamTags(tagDsn), fx.ResultTags(tagConn)),
			fx.Annotate(NewLifecycle, fx.ParamTags(tagConn), fx.ResultTags(tagRqlc)),
		),
		fx.Invoke(
			fx.Annotate(RegisterLifecycle, fx.ParamTags("", tagRqlc)),
		),
	)
}

func genConnURL(name string) string {
	if name == "" {
		return DefaultConnURLKey
	}
	return DefaultConnURLKey + "_" + strings.ToUpper(name)
}

func genConnTag(name string) string {
	if name == "" {
		return ""
	}
	return `name:"rq_conn_` + name + `"`
}

const (
	DefaultConnURLKey = "DB_DSN"
	DefaultConnURL    = "http://localhost"
)
