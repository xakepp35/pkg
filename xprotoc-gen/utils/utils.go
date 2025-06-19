package utils

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	fasthttp "github.com/valyala/fasthttp"
)

type ErrorMessage struct {
	Message string `json:"message"`
}

/*-------------------------------------------------------- FIBER --------------------------------------------------------*/

var (
	HandleUnmarshalError  = HandleNonGrpcError
	HandleValidationError = HandleNonGrpcError
)

// HandleGRPCStatusError default handler for grpc errors
func HandleGRPCStatusError(c *fiber.Ctx, err error) error {
	if err == nil {
		c.Response().Header.SetContentType(fiber.MIMEApplicationJSON)
		return c.Status(fiber.StatusOK).SendString("{}")
	}

	s, ok := status.FromError(err)
	if !ok {
		// non GRPC
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorMessage{Message: err.Error()})
	}

	var httpStatus int
	switch s.Code() {
	case codes.InvalidArgument:
		httpStatus = fiber.StatusBadRequest
	case codes.NotFound:
		httpStatus = fiber.StatusNotFound
	case codes.AlreadyExists:
		httpStatus = fiber.StatusConflict
	case codes.PermissionDenied:
		httpStatus = fiber.StatusForbidden
	case codes.Unauthenticated:
		httpStatus = fiber.StatusUnauthorized
	case codes.ResourceExhausted:
		httpStatus = fiber.StatusTooManyRequests
	case codes.Unimplemented:
		httpStatus = fiber.StatusNotImplemented
	case codes.Internal:
		httpStatus = fiber.StatusInternalServerError
	case codes.Unavailable:
		httpStatus = fiber.StatusServiceUnavailable
	default:
		httpStatus = fiber.StatusInternalServerError
	}

	return c.Status(httpStatus).JSON(s.Proto())
}

// HandleNonGrpcError default handler for non grpc errors
func HandleNonGrpcError(c *fiber.Ctx, err error) error {
	if err == nil {
		c.Response().Header.SetContentType(fiber.MIMEApplicationJSON)
		return c.Status(fiber.StatusOK).SendString("{}")
	}

	return c.Status(fiber.StatusBadRequest).JSON(ErrorMessage{Message: err.Error()})
}

/*-------------------------------------------------------- FASTHTTP --------------------------------------------------------*/

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
