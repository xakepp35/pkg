package xfasthttp

import (
	"context"
	"net"

	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"github.com/xakepp35/pkg/xlog"
	"go.uber.org/fx"
)

func NewServer(handler fasthttp.RequestHandler) *fasthttp.Server {
	return &fasthttp.Server{
		Handler: handler,
	}
}

func RunServer(lc fx.Lifecycle, slc *Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: slc.OnStart,
		OnStop:  slc.OnStop,
	})
}

type Lifecycle struct {
	srv *fasthttp.Server
	cfg *ServerConfig
}

func NewLifecycle(srv *fasthttp.Server, cfg *ServerConfig) *Lifecycle {
	return &Lifecycle{
		srv: srv,
		cfg: cfg,
	}
}

func (s *Lifecycle) OnStart(ctx context.Context) error {
	ln, err := net.Listen("tcp4", s.cfg.Addr)
	if err != nil {
		log.Error().Err(err).Msg("net.Listen")
		return err
	}
	go func() {
		err := s.srv.Serve(ln)
		xlog.FatalInfo(err).
			Msg("srv.Serve")
	}()
	return nil
}

func (s *Lifecycle) OnStop(ctx context.Context) error {
	err := s.srv.ShutdownWithContext(ctx)
	xlog.ErrInfo(err).
		Msg("srv.Shutdown")
	return err
}
