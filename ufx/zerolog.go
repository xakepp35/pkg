package ufx

import (
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"github.com/xakepp35/pkg/xlog"
)

type ZeroLogger struct{}

func WithZeroLogger() fx.Option {
	return fx.WithLogger(func() fxevent.Logger {
		return &ZeroLogger{}
	})
}

func NewZeroLogger() *ZeroLogger {
	return &ZeroLogger{}
}

func (l *ZeroLogger) Use() fx.Option {
	return fx.WithLogger(func() fxevent.Logger {
		return l
	})
}

var _ fxevent.Logger = (*ZeroLogger)(nil)

func (l *ZeroLogger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		log.Info().
			Str("function", e.FunctionName).
			Str("caller", e.CallerName).
			Msg("OnStartExecuting")

	case *fxevent.OnStartExecuted:
		xlog.ErrInfo(e.Err).
			Str("function", e.FunctionName).
			Str("caller", e.CallerName).
			Str("method", e.Method).
			Dur("runtime", e.Runtime).
			Msg("OnStartExecuted")

	case *fxevent.OnStopExecuting:
		log.Info().
			Str("function", e.FunctionName).
			Str("caller", e.CallerName).
			Msg("OnStopExecuting")

	case *fxevent.OnStopExecuted:
		xlog.ErrInfo(e.Err).
			Str("function", e.FunctionName).
			Str("caller", e.CallerName).
			Dur("runtime", e.Runtime).
			Msg("OnStopExecuted")

	case *fxevent.Supplied:
		xlog.ErrInfo(e.Err).
			Str("type", e.TypeName).
			Strs("stack", e.StackTrace).
			Strs("module_trace", e.ModuleTrace).
			Str("module", e.ModuleName).
			Msg("Supplied")

	case *fxevent.Provided:
		xlog.ErrInfo(e.Err).
			Str("constructor", e.ConstructorName).
			Strs("stack", e.StackTrace).
			Strs("module_trace", e.ModuleTrace).
			Strs("output_types", e.OutputTypeNames).
			Str("module", e.ModuleName).
			Bool("private", e.Private).
			Msg("Provided")

	case *fxevent.Replaced:
		xlog.ErrInfo(e.Err).
			Strs("output_types", e.OutputTypeNames).
			Strs("stack", e.StackTrace).
			Strs("module_trace", e.ModuleTrace).
			Str("module", e.ModuleName).
			Msg("Replaced")

	case *fxevent.Decorated:
		xlog.ErrInfo(e.Err).
			Str("decorator", e.DecoratorName).
			Strs("stack", e.StackTrace).
			Strs("module_trace", e.ModuleTrace).
			Strs("output_types", e.OutputTypeNames).
			Str("module", e.ModuleName).
			Msg("Decorated")

	case *fxevent.Run:
		xlog.ErrInfo(e.Err).
			Str("name", e.Name).
			Str("kind", e.Kind).
			Str("module", e.ModuleName).
			Dur("runtime", e.Runtime).
			Msg("Run")

	case *fxevent.Invoking:
		log.Info().
			Str("function", e.FunctionName).
			Str("module", e.ModuleName).
			Msg("Invoking")

	case *fxevent.Invoked:
		xlog.ErrInfo(e.Err).
			Str("function", e.FunctionName).
			Str("module", e.ModuleName).
			Str("trace", e.Trace).
			Msg("Invoked")

	case *fxevent.Started:
		xlog.ErrInfo(e.Err).
			Msg("Started")

	case *fxevent.Stopping:
		log.Info().
			Str("signal", e.Signal.String()).
			Msg("Stopping")

	case *fxevent.Stopped:
		xlog.ErrInfo(e.Err).
			Msg("Stopped")

	case *fxevent.RollingBack:
		xlog.ErrInfo(e.StartErr).
			Msg("RollingBack")

	case *fxevent.RolledBack:
		xlog.ErrInfo(e.Err).
			Msg("RolledBack")

	case *fxevent.LoggerInitialized:
		xlog.ErrInfo(e.Err).
			Str("constructor", e.ConstructorName).
			Msg("LoggerInitialized")

	}
}
