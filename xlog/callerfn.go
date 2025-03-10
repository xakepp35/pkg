package xlog

import (
	"github.com/rs/zerolog"

	"github.com/xakepp35/pkg/xrtm"
)

// HookCallerFunc implements zerolog.Hook to add caller function name
type HookCallerFunc struct{}

var CallerFnFieldName = "func"

// Run adds the function name of the caller to the log event
func (h HookCallerFunc) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	e.Str(CallerFnFieldName, xrtm.CallerFnName(xrtm.CallerFnDefaultSkip+1)) // 3, чтобы пропустить сам хук
}

func CallerFn(addSkip int) func(e *zerolog.Event) {
	return func(e *zerolog.Event) {
		e.Str(CallerFnFieldName, xrtm.CallerFnName(addSkip+1))
	}
}
