package xlog

import (
	"github.com/xakepp35/pkg/xerrors"
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

// AddStackHookKey add hook to global zerolog logger
func AddStackHookKey(key string) error {
	hook, err := RegisterHook(key, nil)
	if err != nil {
		return xerrors.Err(err).Str("key", key).Msg("register stack hook failed")
	}

	log.Logger = log.Hook(hook)

	zerolog.DefaultContextLogger = &log.Logger

	return nil
}
