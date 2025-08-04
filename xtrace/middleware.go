package xtrace

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var propagator = otel.GetTextMapPropagator()

func FiberTraceMiddleware(tracer trace.Tracer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := propagator.Extract(c.Context(), propagation.HeaderCarrier(c.GetReqHeaders()))
		ctx, span := tracer.Start(ctx, c.Path(), trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		// при необходимости span.SpanContext().TraceID().String()

		c.SetUserContext(ctx)
		propagator.Inject(ctx, fasthttpResponseCarrier{h: &c.Response().Header})
		return c.Next()
	}
}

const TraceCtxKey = "trace_ctx"

func FasthttpTraceMiddleware(tracer trace.Tracer) func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			ctxOT := propagator.Extract(context.Background(), fasthttpRequestCarrier{h: &ctx.Request.Header})
			ctxOT, span := tracer.Start(ctxOT, string(ctx.Path()), trace.WithSpanKind(trace.SpanKindServer))
			defer span.End()

			ctx.SetUserValue(TraceCtxKey, ctxOT)
			propagator.Inject(ctxOT, fasthttpResponseCarrier{h: &ctx.Response.Header})
			next(ctx)
		}
	}
}

func GRPCUnaryTraceInterceptor(tracer trace.Tracer) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		carrier := metadataCarrier(md)
		ctxOT := propagator.Extract(ctx, carrier)
		ctxOT, span := tracer.Start(ctxOT, info.FullMethod, trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		resp, err := handler(ctxOT, req)

		setHeaderErr := grpc.SetHeader(ctx, metadata.Pairs(TraceHeader, span.SpanContext().TraceID().String()))
		if setHeaderErr != nil {
			log.Warn().Err(setHeaderErr).Msg("grpc.SetHeader failed")
		}

		return resp, err
	}
}

func GRPCStreamTraceInterceptor(tracer trace.Tracer) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		md, _ := metadata.FromIncomingContext(ss.Context())
		carrier := metadataCarrier(md)
		ctxOT := propagator.Extract(ss.Context(), carrier)
		ctxOT, span := tracer.Start(ctxOT, info.FullMethod, trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		wrapper := &wrappedStream{ServerStream: ss, ctx: ctxOT}

		// Отправляем trace-id в заголовках
		header := metadata.Pairs(TraceHeader, span.SpanContext().TraceID().String())
		if err := wrapper.SendHeader(header); err != nil {
			log.Warn().Err(err).Msg("stream.SendHeader failed")
		}

		return handler(srv, wrapper)
	}
}

type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedStream) Context() context.Context { return w.ctx }

type metadataCarrier metadata.MD

func (mc metadataCarrier) Get(key string) string {
	vals := metadata.MD(mc).Get(key)
	if len(vals) > 0 {
		return vals[0]
	}
	return ""
}
func (mc metadataCarrier) Set(key string, value string) {
	// no-op
}
func (mc metadataCarrier) Keys() []string {
	out := make([]string, 0, len(mc))
	for k := range mc {
		out = append(out, k)
	}
	return out
}

type fasthttpResponseCarrier struct {
	h *fasthttp.ResponseHeader
}

func (c fasthttpResponseCarrier) Get(key string) string {
	return string(c.h.Peek(key))
}

func (c fasthttpResponseCarrier) Set(key string, value string) {
	c.h.Set(key, value)
}

func (c fasthttpResponseCarrier) Keys() []string {
	keys := []string{}
	c.h.VisitAll(func(k, _ []byte) {
		keys = append(keys, string(k))
	})
	return keys
}

type fasthttpRequestCarrier struct {
	h *fasthttp.RequestHeader
}

func (c fasthttpRequestCarrier) Get(key string) string {
	return string(c.h.Peek(key))
}

func (c fasthttpRequestCarrier) Set(key string, value string) {
	c.h.Set(key, value)
}

func (c fasthttpRequestCarrier) Keys() []string {
	var keys []string
	c.h.VisitAll(func(k, _ []byte) {
		keys = append(keys, string(k))
	})
	return keys
}
