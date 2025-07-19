package xerrors

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
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

	e1 := Err(orig).Str("x", "1").Send()
	e2 := Err(orig).Str("y", "2").Send()

	if e1 == nil || e2 == nil {
		t.Fatal("expected non-nil errors")
	}
}

func TestErrBuilder_Empty(t *testing.T) {
	orig := errors.New("something failed")
	err := Err(orig).Send()

	if got := err.Error(); !strings.Contains(got, "something failed") {
		t.Errorf("unexpected error string: %s", got)
	}
}

func TestFielder_Bool(t *testing.T) {
	t.Run("true value", func(t *testing.T) {
		err := Err(errors.New("base error")).Bool("flag", true).Msg("operation failed")

		expected := "operation failed flag=true: base error"
		require.Equal(t, expected, err.Error())
	})

	t.Run("false value", func(t *testing.T) {
		err := Err(errors.New("base error")).Bool("flag", false).Msg("operation failed")

		expected := "operation failed flag=false: base error"
		require.Equal(t, expected, err.Error())
	})
}

func TestFielder_Time(t *testing.T) {
	testTime := time.Date(2025, 7, 19, 15, 30, 0, 0, time.UTC)
	err := Err(errors.New("database error")).Time("timestamp", testTime).Msg("fetch failed")

	estr := err.Error()
	require.Contains(t, estr, "fetch failed timestamp=\"2025-07-19T15:30:00Z\": database error")
}

func TestFielder_Bools(t *testing.T) {
	boolSlice := []bool{true, false, true}
	err := Err(errors.New("validation error")).Bools("flags", boolSlice).Msg("multiple issues")

	require.Contains(t, err.Error(), "multiple issues flags=[true,false,true]: validation error")
}

func TestFielder_Int64s(t *testing.T) {
	t.Run("int64 slice default", func(t *testing.T) {
		intSlice := []int64{1, 2, 3}
		err := Err(errors.New("validation error")).Int64s("ids", intSlice).Msg("multiple issues")

		require.Contains(t, err.Error(), "multiple issues ids=[1,2,3]: validation error")
	})

	t.Run("int64 slice empty", func(t *testing.T) {
		intSlice := []int64{}
		err := Err(errors.New("validation error")).Int64s("ids", intSlice).Msg("multiple issues")

		require.Contains(t, err.Error(), "multiple issues ids=[]: validation error")
	})
}
