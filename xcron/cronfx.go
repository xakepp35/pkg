package xcron

import (
	"context"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

// Job — функция задачи. Возвращаемая ошибка попадёт в логи.
type Job func(ctx context.Context) error

// Scheduler — интерфейс планировщика для DI.
type Scheduler interface {
	// Add регистрирует задачу по cron-спеку и возвращает её ID.
	Add(spec string, fn Job, opts ...JobOption) (cron.EntryID, error)
	// Now запускает задачу немедленно (без расписания).
	Now(fn Job, opts ...JobOption)
}

type scheduler struct {
	c      *cron.Cron
	ctx    context.Context
	cancel context.CancelFunc
	mu     sync.RWMutex
	names  map[cron.EntryID]string
}

// --- опции планировщика ---

type Option func(*schedulerConfig)

type schedulerConfig struct {
	location *time.Location
	recover  bool
}

func WithLocation(loc *time.Location) Option { return func(c *schedulerConfig) { c.location = loc } }
func WithRecover() Option                    { return func(c *schedulerConfig) { c.recover = true } }

// --- опции задач ---

type JobOption func(*jobOptions)

type jobOptions struct {
	name               string
	timeout            time.Duration
	skipIfStillRunning bool
}

func WithName(name string) JobOption        { return func(o *jobOptions) { o.name = name } }
func WithTimeout(d time.Duration) JobOption { return func(o *jobOptions) { o.timeout = d } }
func SkipIfStillRunning() JobOption         { return func(o *jobOptions) { o.skipIfStillRunning = true } }

// New — провайдер планировщика. Добавляет lifecycle hooks для старта/остановки.
func New(lc fx.Lifecycle, opts ...Option) Scheduler {
	cfg := schedulerConfig{location: time.Local, recover: true}
	for _, o := range opts {
		o(&cfg)
	}

	parser := cron.NewParser(
		cron.SecondOptional |
			cron.Second |
			cron.Minute |
			cron.Hour |
			cron.Dom |
			cron.Month |
			cron.Dow |
			cron.Descriptor,
	)

	var wrappers []cron.JobWrapper
	if cfg.recover {
		wrappers = append(wrappers, cron.Recover(cronLoggerAdapter{}))
	}

	c := cron.New(
		cron.WithParser(parser),
		cron.WithLocation(cfg.location),
		cron.WithChain(wrappers...),
		cron.WithLogger(cronLoggerAdapter{}),
	)

	ctx, cancel := context.WithCancel(context.Background())
	s := &scheduler{c: c, ctx: ctx, cancel: cancel, names: make(map[cron.EntryID]string)}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			log.Info().Msg("cron scheduler starting")
			s.c.Start()
			log.Info().Msg("cron scheduler started")
			return nil
		},
		OnStop: func(context.Context) error {
			log.Info().Msg("cron scheduler stopping")
			s.cancel()
			// Останавливаем, дожидаемся завершения активных задач
			stopCtx := s.c.Stop()
			<-stopCtx.Done()
			log.Info().Msg("cron scheduler stopped")
			return nil
		},
	})

	return s
}

func (s *scheduler) wrap(fn Job, o jobOptions) cron.Job {
	base := cron.FuncJob(func() {
		start := time.Now()
		ctx := s.ctx
		if o.timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, o.timeout)
			defer cancel()
		}

		if o.name != "" {
			log.Info().Str("job", o.name).Msg("cron job start")
		} else {
			log.Info().Msg("cron job start")
		}

		if err := fn(ctx); err != nil {
			if o.name != "" {
				log.Error().Str("job", o.name).Err(err).Dur("duration", time.Since(start)).Msg("cron job error")
			} else {
				log.Error().Err(err).Dur("duration", time.Since(start)).Msg("cron job error")
			}
			return
		}

		if o.name != "" {
			log.Info().Str("job", o.name).Dur("duration", time.Since(start)).Msg("cron job done")
		} else {
			log.Info().Dur("duration", time.Since(start)).Msg("cron job done")
		}
	})

	if o.skipIfStillRunning {
		return cron.NewChain(cron.SkipIfStillRunning(cronLoggerAdapter{})).Then(base)
	}
	return base
}

func (s *scheduler) Add(spec string, fn Job, opts ...JobOption) (cron.EntryID, error) {
	var o jobOptions
	for _, opt := range opts {
		opt(&o)
	}

	id, err := s.c.AddJob(spec, s.wrap(fn, o))
	if err == nil {
		if o.name != "" {
			s.mu.Lock()
			s.names[id] = o.name
			s.mu.Unlock()
			log.Info().Str("job", o.name).Str("spec", spec).Int64("entry_id", int64(id)).Msg("cron job registered")
		} else {
			log.Info().Str("spec", spec).Int64("entry_id", int64(id)).Msg("cron job registered")
		}
	}
	return id, err
}

func (s *scheduler) Now(fn Job, opts ...JobOption) {
	var o jobOptions
	for _, opt := range opts {
		opt(&o)
	}
	if o.name != "" {
		log.Info().Str("job", o.name).Msg("cron job immediate run")
	} else {
		log.Info().Msg("cron job immediate run")
	}
	go s.wrap(fn, o).Run()
}

// Register — хелпер для декларативной регистрации задач из модулей.
func Register(spec string, fn Job, opts ...JobOption) fx.Option {
	return fx.Invoke(func(s Scheduler) {
		if _, err := s.Add(spec, fn, opts...); err != nil {
			panic(err)
		}
	})
}

// Module — готовый fx-модуль планировщика.
var Module = fx.Module("xcron",
	fx.Provide(New),
)

// --- zerolog adapters for robfig/cron v3 ---
// cronLoggerAdapter satisfies cron.Logger (used by cron.WithLogger and wrappers).
type cronLoggerAdapter struct{}

func (cronLoggerAdapter) Info(msg string, keysAndValues ...interface{}) {
	e := log.Info()
	for i := 0; i+1 < len(keysAndValues); i += 2 {
		if k, ok := keysAndValues[i].(string); ok {
			e = e.Interface(k, keysAndValues[i+1])
		}
	}
	e.Msg(msg)
}

func (cronLoggerAdapter) Error(err error, msg string, keysAndValues ...interface{}) {
	e := log.Error().Err(err)
	for i := 0; i+1 < len(keysAndValues); i += 2 {
		if k, ok := keysAndValues[i].(string); ok {
			e = e.Interface(k, keysAndValues[i+1])
		}
	}
	e.Msg(msg)
}
