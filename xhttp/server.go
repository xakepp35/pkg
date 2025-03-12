package xhttp

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/xakepp35/pkg/env"
)

type ServerConfig struct {
	Addr string `json:"addr"`
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Addr: ":" + env.String("PORT", "8080"),
	}
}

// startHTTPServer запускает HTTP-сервер.
func RunServer(cfg *ServerConfig, mux chi.Router) {
	log.Info().Str("addr", cfg.Addr).Msg("Starting HTTP server")
	if err := http.ListenAndServe(cfg.Addr, mux); err != nil {
		log.Fatal().Err(err).Msg("http.ListenAndServe")
	}
	log.Info().Msg("Stopping HTTP server")
}
