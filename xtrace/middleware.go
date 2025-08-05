package xtrace

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var propagator = otel.GetTextMapPropagator()

func FiberTraceMiddleware(c *fiber.Ctx) error {
	ctx := propagator.Extract(c.Context(), propagation.HeaderCarrier(c.GetReqHeaders()))
	ctx, span := otel.Tracer(defaultTracerName).Start(ctx, c.Path(), trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	span.SetAttributes(
		attribute.String("http.method", c.Method()),
		attribute.String("http.url", c.OriginalURL()),
		attribute.String("http.client_ip", c.IP()),
		attribute.String("http.user_agent", c.Get(fiber.HeaderUserAgent)),
	)

	c.SetUserContext(ctx)
	err := c.Next()

	span.SetAttributes(
		attribute.Int("http.status_code", c.Response().StatusCode()),
		attribute.String("http.response_size", fmt.Sprintf("%d", len(c.Response().Body()))),
	)
	if err != nil {
		span.RecordError(err)
	}
	span.AddEvent("request_completed")

	propagator.Inject(ctx, fasthttpResponseCarrier{h: &c.Response().Header})
	return err
}

const TraceCtxKey = "trace_ctx"

func FasthttpTraceMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctxOT := propagator.Extract(context.Background(), fasthttpRequestCarrier{h: &ctx.Request.Header})
		ctxOT, span := otel.Tracer(defaultTracerName).Start(ctxOT, string(ctx.Path()), trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		span.SetAttributes(
			attribute.String("http.method", string(ctx.Method())),
			attribute.String("http.url", ctx.URI().String()),
			attribute.String("http.client_ip", ctx.RemoteIP().String()),
			attribute.String("http.user_agent", string(ctx.UserAgent())),
		)

		ctx.SetUserValue(TraceCtxKey, ctxOT)
		next(ctx)

		span.SetAttributes(
			attribute.Int("http.status_code", ctx.Response.StatusCode()),
			attribute.Int("http.response_size", len(ctx.Response.Body())),
		)
		span.AddEvent("request_completed")
		propagator.Inject(ctxOT, fasthttpResponseCarrier{h: &ctx.Response.Header})
	}
}

func GRPCUnaryTraceInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	ctxOT := propagator.Extract(ctx, metadataCarrierFromContext(ctx))
	ctxOT, span := otel.Tracer(defaultTracerName).Start(ctxOT, info.FullMethod, trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	span.SetAttributes(
		attribute.String("rpc.system", "grpc"),
		attribute.String("rpc.method", info.FullMethod),
	)

	resp, err := handler(ctxOT, req)
	if err != nil {
		span.RecordError(err)
	}

	span.AddEvent("grpc_response_sent")

	setHeaderErr := grpc.SetHeader(ctx, metadata.Pairs(TraceHeader, span.SpanContext().TraceID().String()))
	if setHeaderErr != nil {
		log.Warn().Err(setHeaderErr).Msg("grpc.SetHeader failed")
	}

	return resp, err
}

func GRPCStreamTraceInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctxOT := propagator.Extract(ss.Context(), metadataCarrierFromContext(ss.Context()))
	ctxOT, span := otel.Tracer(defaultTracerName).Start(ctxOT, info.FullMethod, trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	span.SetAttributes(
		attribute.String("rpc.system", "grpc"),
		attribute.String("rpc.method", info.FullMethod),
		attribute.String("rpc.stream", fmt.Sprintf("%v", info.IsServerStream)),
	)

	wrapper := &wrappedStream{ServerStream: ss, ctx: ctxOT}
	if err := wrapper.SendHeader(metadata.Pairs(TraceHeader, span.SpanContext().TraceID().String())); err != nil {
		log.Warn().Err(err).Msg("stream.SendHeader failed")
	}

	span.AddEvent("grpc_stream_started")
	err := handler(srv, wrapper)
	if err != nil {
		span.RecordError(err)
	}
	span.AddEvent("grpc_stream_completed")
	return err
}

type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedStream) Context() context.Context { return w.ctx }

func metadataCarrierFromContext(ctx context.Context) metadataCarrier {
	md, _ := metadata.FromIncomingContext(ctx)
	return metadataCarrier(md)
}

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
