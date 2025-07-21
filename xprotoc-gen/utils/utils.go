package utils

import (
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrorMessage struct {
	Message string `json:"message"`
}

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

var (
	HandleUnmarshalError  = HandleNonGrpcError
	HandleValidationError = HandleNonGrpcError
)

// HandleNonGrpcError default handler for non grpc errors
func HandleNonGrpcError(c *fiber.Ctx, err error) error {
	if err == nil {
		c.Response().Header.SetContentType(fiber.MIMEApplicationJSON)
		return c.Status(fiber.StatusOK).SendString("{}")
	}

	return c.Status(fiber.StatusBadRequest).JSON(ErrorMessage{Message: err.Error()})
}
