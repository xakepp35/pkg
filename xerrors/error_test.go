package xerrors

import (
	"database/sql"
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
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
		err := Err(sql.ErrNoRows).Msg("not found").Str("foo", "bar").Err()

		require.True(t, errors.Is(err, sql.ErrNoRows))
	})
}
