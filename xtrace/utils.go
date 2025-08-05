package xtrace

import (
	"context"
	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel"
)

func InjectFastHttp(ctx context.Context, headers *fasthttp.RequestHeader) {
	otel.GetTextMapPropagator().Inject(ctx, fasthttpRequestCarrier{h: headers})
}
