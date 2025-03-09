package xlog

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/xakepp35/pkg/xrtm"
)

// Init should be called once, from init() before main()
func Init() {
	// Set custom time function for UTC time in RFC3339Nano format
	zerolog.TimestampFunc = xrtm.TimeUTC

	// Set the format for the time field
	zerolog.TimeFieldFormat = time.RFC3339Nano

	// Create custom logger with caller hook
	log.Logger = zerolog.
		New(os.Stdout).
		With().
		Timestamp().
		Logger().
		Hook(HookCallerFunc{})

	zerolog.DefaultContextLogger = &log.Logger
}

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
