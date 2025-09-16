package xerrors

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestEqual(t *testing.T) {

	t.Run("str", func(t *testing.T) {
		err := Err(sql.ErrNoRows).Str("foo", "bar").Msg("not found")
		require.Equal(t, err.Error(), fmt.Errorf("not found foo=bar: %w", sql.ErrNoRows).Error())
	})

	t.Run("int64", func(t *testing.T) {
		err := Err(sql.ErrNoRows).Int64("foo", 123).Msg("not found")

		require.Equal(t, err.Error(), fmt.Errorf("not found foo=123: %w", sql.ErrNoRows).Error())
	})
}

func TestErrBuilder_Basic(t *testing.T) {
	orig := errors.New("original error")

	err := Err(orig).Str("field1", "value1").Str("field2", "value2").Msg("message 123")
	if err == nil {
		t.Fatal("expected non-nil error")
	}

	if !strings.Contains(err.Error(), "field1") {
		t.Error("missing 'field1' in error string")
	}
	if !strings.Contains(err.Error(), "original error") {
		t.Error("missing original error in result")
	}
}

func TestErrBuilder_ReusesBuffer(t *testing.T) {
	orig := errors.New("base")

	builder1 := Err(orig).(*errorBuilder)
	e1 := builder1.Str("x", "1").Send()

	builder2 := Err(orig).(*errorBuilder)
	e2 := builder2.Str("y", "2").Send()

	if e1 == nil || e2 == nil {
		t.Fatal("expected non-nil errors")
	}

	// Fill buffers before comparing
	builder1.errBuffer = append(builder1.errBuffer, []byte("test")...)
	builder1.argsBuffer = append(builder1.argsBuffer, []byte("test")...)

	require.Equal(t, &builder1.errBuffer[0], &builder2.errBuffer[0], "error buffers should have same address")
	require.Equal(t, &builder1.argsBuffer[0], &builder2.argsBuffer[0], "args buffers should have same address")
	require.Equal(t, cap(builder1.errBuffer), cap(builder2.errBuffer), "error buffers should have same capacity")
	require.Equal(t, cap(builder1.argsBuffer), cap(builder2.argsBuffer), "args buffers should have same capacity")
	require.Equal(t, builder1, builder2, "expected same builder instance to be reused")
}

func TestErrBuilder_Empty(t *testing.T) {
	orig := errors.New("something failed")
	err := Err(orig).Send()

	if got := err.Error(); !strings.Contains(got, "something failed") {
		t.Errorf("unexpected error string: %s", got)
	}
}

func TestErrBuilder_ProtoMsg(t *testing.T) {
	orig := errors.New("foo")
	t.Run("check unwrapped error", func(t *testing.T) {
		err := Err(orig).MsgProto(codes.AlreadyExists, "bar")
		require.Equal(t, "bar: foo", err.Error())
	})
	t.Run("check code by error", func(t *testing.T) {
		err := Err(orig).MsgProto(codes.AlreadyExists, "bar")
		grpcStatus, ok := status.FromError(err)
		fmt.Println(grpcStatus, ok)
		if !ok {
			t.Error("got absent grpc status")
			return
		}
		if grpcStatus.Code() != codes.AlreadyExists {
			t.Errorf("got grpc status code: %v, want %v", grpcStatus.Code().String(), codes.AlreadyExists.String())
			return
		}
		if grpcStatus.Message() != "bar" {
			t.Errorf("got grpc status message: %v, want bar", grpcStatus.Message())
			return
		}
	})
}

func TestErrBuilder_Proto(t *testing.T) {
	orig := errors.New("foo")
	t.Run("check unwrapped error", func(t *testing.T) {
		err := Err(orig).Proto(codes.AlreadyExists)
		require.Equal(t, "foo", err.Error())
	})
	t.Run("check code by error", func(t *testing.T) {
		err := Err(orig).Proto(codes.AlreadyExists)
		grpcStatus, ok := status.FromError(err)
		if !ok {
			t.Error("got absent grpc status")
			return
		}
		if grpcStatus.Code() != codes.AlreadyExists {
			t.Errorf("got grpc status code: %v, want %v", grpcStatus.Code().String(), codes.AlreadyExists.String())
			return
		}
		if grpcStatus.Message() != "" {
			t.Errorf("got grpc status message: %v, want bar", grpcStatus.Message())
			return
		}
	})
}
