package xfasthttp

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"

	"github.com/xakepp35/pkg/xlog"
)

// FasthttpZerolog — логгирующая мидлвара для fasthttp
func MiddlewareZerolog(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		start := time.Now()
		defer func() {
			duration := time.Since(start)
			requestBody := ctx.PostBody()
			responseBody := ctx.Response.Body()
			log.Debug().
				Str("method", string(ctx.Method())).
				Str("route", string(ctx.Path())).
				Int("status", ctx.Response.StatusCode()).
				Dur("cost", duration).
				Str("headers", string(ctx.Request.Header.RawHeaders())).
				Str("uri", string(ctx.Request.RequestURI())).
				Int("req_size", len(requestBody)).
				Int("res_size", len(responseBody)).
				Func(xlog.RawJSON("req_body", requestBody)).
				Func(xlog.RawJSON("res_body", responseBody)).
				Msg("next")
		}()
		next(ctx)
	}
}

// мидлвара для отлова паники и логирования через zerolog
func MiddlewarePanicRecovery(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if err := recover(); err != nil {
				log.Error().
					Interface("panic", err).
					Stack().
					Msg("recover")
				ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
			}
		}()
		next(ctx)
	}
}
