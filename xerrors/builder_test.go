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
		err := Err(sql.ErrNoRows).Msg("not found").Str("foo", "bar").Err()

		require.Equal(t, err.Error(), fmt.Errorf("not found foo=bar: %w", sql.ErrNoRows).Error())
	})

	t.Run("int64", func(t *testing.T) {
		err := Err(sql.ErrNoRows).Msg("not found").Int64("foo", 123).Err()

		require.Equal(t, err.Error(), fmt.Errorf("not found foo=123: %w", sql.ErrNoRows).Error())
	})
}

func TestErrBuilder_Basic(t *testing.T) {
	orig := errors.New("original error")

	err := Err(orig).Msg("message 123").Str("field1", "value1").Str("field2", "value2").Err()
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

	e1 := Err(orig).Str("x", "1").Err()
	e2 := Err(orig).Str("y", "2").Err()

	if e1 == nil || e2 == nil {
		t.Fatal("expected non-nil errors")
	}
}

func TestErrBuilder_Empty(t *testing.T) {
	orig := errors.New("something failed")
	err := Err(orig).Err()

	if got := err.Error(); !strings.Contains(got, "something failed") {
		t.Errorf("unexpected error string: %s", got)
	}
}
