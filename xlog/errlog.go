package xlog

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ErrLevel(level zerolog.Level, err error) *zerolog.Event {
	if err != nil {
		return log.Error().Err(err)
	}
	return log.WithLevel(level)
}

func ErrDebug(err error) *zerolog.Event {
	if err != nil {
		return log.Error().Err(err)
	}
	return log.Debug()
}

func ErrInfo(err error) *zerolog.Event {
	if err != nil {
		return log.Error().Err(err)
	}
	return log.Info()
}

func ErrWarn(err error) *zerolog.Event {
	if err != nil {
		return log.Error().Err(err)
	}
	return log.Warn()
}

func FatalInfo(err error) *zerolog.Event {
	if err != nil {
		return log.Fatal().Err(err)
	}
	return log.Info()
}
