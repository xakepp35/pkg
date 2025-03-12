package xhttp

import (
	"encoding/json"
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
	log.Info().
		Any("cfg", cfg).
		Msg("xhttp.StartServer")
	if err := http.ListenAndServe(cfg.Addr, mux); err != nil {
		log.Fatal().
			Err(err).
			Msg("http.ListenAndServe")
	}
	log.Info().
		Msg("xhttp.StopServer")
}

func RespondJSON(w http.ResponseWriter, val any) {
	w.Header().Set(HeaderContentType, ContentTypeJson)
	if err := json.NewEncoder(w).Encode(val); err != nil {
		log.Error().
			Err(err).
			Any("val", val).
			Msg("json.NewEncoder(w).Encode(val)")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
