package xfasthttp

import "go.uber.org/fx"

var Module = fx.Module("xfasthttp",
	fx.Supply(
		NewServerConfig(),
	),
	fx.Provide(
		NewServer,
		NewLifecycle,
	),
	fx.Invoke(
		RunServer,
	),
)
