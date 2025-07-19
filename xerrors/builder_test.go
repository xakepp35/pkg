package xerrors

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
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
