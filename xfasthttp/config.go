package xfasthttp

import "github.com/xakepp35/pkg/env"

type ServerConfig struct {
	Addr string
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Addr: ":" + env.String("PORT", "8080"),
	}
}
