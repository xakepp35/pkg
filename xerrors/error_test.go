package xerrors

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestNew(t *testing.T) {
	err := New(New(errors.New("000"), "msg"), "bar")

	require.Equal(t, "bar: msg: 000", err.Error())
}

func TestIs(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		wrappedErr := New(sql.ErrNoRows, "123")

		err := New(wrappedErr, "msg")

		require.True(t, errors.Is(err, sql.ErrNoRows))
	})

	t.Run("builder", func(t *testing.T) {
		err := Err(sql.ErrNoRows).Str("foo", "bar").Msg("not found")

		require.True(t, errors.Is(err, sql.ErrNoRows))
	})

	t.Run("err by proto", func(t *testing.T) {
		err := Err(sql.ErrNoRows).Proto(codes.NotFound)

		require.True(t, errors.Is(err, sql.ErrNoRows))
	})

	t.Run("err by proto msg", func(t *testing.T) {
		err := Err(sql.ErrNoRows).MsgProto(codes.NotFound, "not found")

		require.True(t, errors.Is(err, sql.ErrNoRows))
	})
}

func TestNewProto(t *testing.T) {
	t.Run("check unwrapped error", func(t *testing.T) {
		err := NewProto(codes.AlreadyExists, errors.New("foo"), "bar")
		require.Equal(t, "bar: foo", err.Error())
	})
	t.Run("check code by error", func(t *testing.T) {
		err := NewProto(codes.AlreadyExists, errors.New("foo"), "bar")
		status, ok := status.FromError(err)
		if !ok {
			t.Error("got absent grpc status")
			return
		}
		if status.Code() != codes.AlreadyExists {
			t.Errorf("got grpc status code: %v, want %v", status.Code().String(), codes.AlreadyExists.String())
			return
		}
		if status.Message() != "bar" {
			t.Errorf("got grpc status message: %v, want bar", status.Message())
			return
		}
	})
}

func TestNewProto(t *testing.T) {
	t.Run("check unwrapped error", func(t *testing.T) {
		err := NewProto(codes.AlreadyExists, errors.New("foo"), "bar")
		require.Equal(t, "bar: foo", err.Error())
	})
	t.Run("check code by error", func(t *testing.T) {
		err := NewProto(codes.AlreadyExists, errors.New("foo"), "bar")
		status, ok := status.FromError(err)
		if !ok {
			t.Error("got absent grpc status")
			return
		}
		if status.Code() != codes.AlreadyExists {
			t.Errorf("got grpc status code: %v, want %v", status.Code().String(), codes.AlreadyExists.String())
			return
		}
		if status.Message() != "bar" {
			t.Errorf("got grpc status message: %v, want bar", status.Message())
			return
		}
	})
}
