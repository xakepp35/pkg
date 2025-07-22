package utils

import (
	"encoding/json"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	fasthttp "github.com/valyala/fasthttp"
)

type ErrorMessage struct {
	Message string `json:"message"`
}

var (
	FastHttpHandleUnmarshalError  = FastHTTPHandleNonGrpcError
	FastHttpHandleValidationError = FastHTTPHandleNonGrpcError
)

func writeJSON(c *fasthttp.RequestCtx, statusCode int, v any) {
	data, err := json.Marshal(v)
	if err != nil {
		c.Error("unable to marshal JSON", fasthttp.StatusInternalServerError)
		return
	}

	c.Response.Header.SetContentType("application/json")
	c.SetStatusCode(statusCode)

	if _, err = c.Write(data); err != nil {
		c.Error("internal server error", fasthttp.StatusInternalServerError)
	}
}

func FastHTTPHandleGRPCStatusError(c *fasthttp.RequestCtx, err error) {
	if err == nil {
		c.Response.Header.SetContentType("application/json")
		c.SetStatusCode(fasthttp.StatusOK)
		return
	}

	s, ok := status.FromError(err)
	if !ok {
		writeJSON(c, fasthttp.StatusInternalServerError, ErrorMessage{Message: err.Error()})
		return
	}

	var httpStatus int
	switch s.Code() {
	case codes.InvalidArgument:
		httpStatus = fasthttp.StatusBadRequest
	case codes.NotFound:
		httpStatus = fasthttp.StatusNotFound
	case codes.AlreadyExists:
		httpStatus = fasthttp.StatusConflict
	case codes.PermissionDenied:
		httpStatus = fasthttp.StatusForbidden
	case codes.Unauthenticated:
		httpStatus = fasthttp.StatusUnauthorized
	case codes.ResourceExhausted:
		httpStatus = fasthttp.StatusTooManyRequests
	case codes.Unimplemented:
		httpStatus = fasthttp.StatusNotImplemented
	case codes.Internal:
		httpStatus = fasthttp.StatusInternalServerError
	case codes.Unavailable:
		httpStatus = fasthttp.StatusServiceUnavailable
	default:
		httpStatus = fasthttp.StatusInternalServerError
	}

	writeJSON(c, httpStatus, s.Proto())
}

func FastHTTPHandleNonGrpcError(c *fasthttp.RequestCtx, err error) {
	if err == nil {
		c.Response.Header.SetContentType("application/json")
		c.SetStatusCode(fasthttp.StatusOK)
		c.Write([]byte("{}"))
		return
	}

	writeJSON(c, fasthttp.StatusBadRequest, ErrorMessage{Message: err.Error()})
}
