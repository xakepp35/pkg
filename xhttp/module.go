package xhttp

import (
	"context"

	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
)

var Module = fx.Module("xhttp",
	fx.Supply(
		NewServerConfig(),
	),
	fx.Invoke(Run),
)

func Run(lc fx.Lifecycle, mux chi.Router, cfg *ServerConfig) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go RunServer(cfg, mux)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}
