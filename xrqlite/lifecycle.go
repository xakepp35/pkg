package xrqlite

import (
	"context"

	"github.com/rqlite/gorqlite"
	"github.com/rs/zerolog/log"
	"github.com/xakepp35/pkg/xlog"
	"go.uber.org/fx"
)

type Lifecycle struct {
	conn *gorqlite.Connection
}

func NewConnection(connURL string) (*gorqlite.Connection, error) {
	res, err := gorqlite.OpenWithClient(connURL, gorqlite.DefaultHTTPClient)
	xlog.ErrInfo(err).Str("connURL", connURL).Msg("gorqlite.OpenWithClient")
	if err != nil {
		return nil, err
	}
	return res, nil
}

func NewLifecycle(conn *gorqlite.Connection) *Lifecycle {
	return &Lifecycle{
		conn: conn,
	}
}

func RegisterLifecycle(lc fx.Lifecycle, rqlc *Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: rqlc.OnStart,
		OnStop:  rqlc.OnStop,
	})
}

// OnStart pings the database
func (s *Lifecycle) OnStart(ctx context.Context) error {
	_, err := s.conn.QueryOne("SELECT 1")
	xlog.ErrInfo(err).Msg("conn.QueryOne")
	if err != nil {

		return err
	}
	return nil
}

// OnStop gracefully shuts down the rqlite connection
func (s *Lifecycle) OnStop(ctx context.Context) error {
	s.conn.Close()
	log.Info().Msg("conn.Close")
	return nil
}
