package xerrors

import (
	"errors"
	"github.com/stretchr/testify/require"
	"github.com/xakepp35/pkg/src/pkg/types"
	"strings"
	"testing"
	"time"
)

func TestFielder_Str(t *testing.T) {
	err := Err(errors.New("base error")).Str("key", "value").Msg("message")
	require.Contains(t, err.Error(), "message key=value: base error")
}

func TestFielder_Float64(t *testing.T) {
	err := Err(errors.New("base error")).Float64("key", 123.456).Msg("message")
	require.Contains(t, err.Error(), "message key=123.456: base error")
}

func TestFielder_Float32(t *testing.T) {
	err := Err(errors.New("base error")).Float32("key", 123.45).Msg("message")
	require.Contains(t, err.Error(), "message key=123.45: base error")
}

func TestFielder_Int64(t *testing.T) {
	err := Err(errors.New("base error")).Int64("key", 123).Msg("message")
	require.Contains(t, err.Error(), "message key=123: base error")
}

func TestFielder_Uint64(t *testing.T) {
	err := Err(errors.New("base error")).Uint64("key", 123).Msg("message")
	require.Contains(t, err.Error(), "message key=123: base error")
}

func TestFielder_Int(t *testing.T) {
	err := Err(errors.New("base error")).Int("key", 123).Msg("message")
	require.Contains(t, err.Error(), "message key=123: base error")
}

func TestFielder_Int8(t *testing.T) {
	err := Err(errors.New("base error")).Int8("key", 123).Msg("message")
	require.Contains(t, err.Error(), "message key=123: base error")
}

func TestFielder_Int16(t *testing.T) {
	err := Err(errors.New("base error")).Int16("key", 123).Msg("message")
	require.Contains(t, err.Error(), "message key=123: base error")
}

func TestFielder_Int32(t *testing.T) {
	err := Err(errors.New("base error")).Int32("key", 123).Msg("message")
	require.Contains(t, err.Error(), "message key=123: base error")
}

func TestFielder_Uint(t *testing.T) {
	err := Err(errors.New("base error")).Uint("key", 123).Msg("message")
	require.Contains(t, err.Error(), "message key=123: base error")
}

func TestFielder_Uint8(t *testing.T) {
	err := Err(errors.New("base error")).Uint8("key", 123).Msg("message")
	require.Contains(t, err.Error(), "message key=123: base error")
}

func TestFielder_Uint16(t *testing.T) {
	err := Err(errors.New("base error")).Uint16("key", 123).Msg("message")
	require.Contains(t, err.Error(), "message key=123: base error")
}

func TestFielder_Uint32(t *testing.T) {
	err := Err(errors.New("base error")).Uint32("key", 123).Msg("message")
	require.Contains(t, err.Error(), "message key=123: base error")
}

func TestFielder_Bool(t *testing.T) {
	t.Run("true value", func(t *testing.T) {
		err := Err(errors.New("base error")).Bool("flag", true).Msg("message")
		require.Contains(t, err.Error(), "message flag=true: base error")
	})

	t.Run("false value", func(t *testing.T) {
		err := Err(errors.New("base error")).Bool("flag", false).Msg("message")
		require.Contains(t, err.Error(), "message flag=false: base error")
	})
}

func TestFielder_Strs(t *testing.T) {
	t.Run("non-empty slice", func(t *testing.T) {
		values := []string{"one", "two", "three"}
		err := Err(errors.New("base error")).Strs("keys", values).Msg("message")
		require.Contains(t, err.Error(), "message keys=[\"one\",\"two\",\"three\"]: base error")
	})

	t.Run("empty slice", func(t *testing.T) {
		values := []string{}
		err := Err(errors.New("base error")).Strs("keys", values).Msg("message")
		require.Contains(t, err.Error(), "message keys=[]: base error")
	})
}

func TestFielder_Float64s(t *testing.T) {
	t.Run("non-empty slice", func(t *testing.T) {
		values := []float64{1.1, 2.2, 3.3}
		err := Err(errors.New("base error")).Float64s("values", values).Msg("message")
		require.Contains(t, err.Error(), "message values=[1.1,2.2,3.3]: base error")
	})

	t.Run("empty slice", func(t *testing.T) {
		values := []float64{}
		err := Err(errors.New("base error")).Float64s("values", values).Msg("message")
		require.Contains(t, err.Error(), "message values=[]: base error")
	})
}

func TestFielder_Float32s(t *testing.T) {
	t.Run("non-empty slice", func(t *testing.T) {
		values := []float32{1.1, 2.2, 3.3}
		err := Err(errors.New("base error")).Float32s("values", values).Msg("message")
		require.Contains(t, err.Error(), "message values=[1.1,2.2,3.3]: base error")
	})

	t.Run("empty slice", func(t *testing.T) {
		values := []float32{}
		err := Err(errors.New("base error")).Float32s("values", values).Msg("message")
		require.Contains(t, err.Error(), "message values=[]: base error")
	})
}

func TestFielder_Int64s(t *testing.T) {
	t.Run("non-empty slice", func(t *testing.T) {
		values := []int64{1, 2, 3}
		err := Err(errors.New("base error")).Int64s("values", values).Msg("message")
		require.Contains(t, err.Error(), "message values=[1,2,3]: base error")
	})

	t.Run("empty slice", func(t *testing.T) {
		values := []int64{}
		err := Err(errors.New("base error")).Int64s("values", values).Msg("message")
		require.Contains(t, err.Error(), "message values=[]: base error")
	})
}

func TestFielder_Uint64s(t *testing.T) {
	t.Run("non-empty slice", func(t *testing.T) {
		values := []uint64{1, 2, 3}
		err := Err(errors.New("base error")).Uint64s("values", values).Msg("message")
		require.Contains(t, err.Error(), "message values=[1,2,3]: base error")
	})

	t.Run("empty slice", func(t *testing.T) {
		values := []uint64{}
		err := Err(errors.New("base error")).Uint64s("values", values).Msg("message")
		require.Contains(t, err.Error(), "message values=[]: base error")
	})
}

func TestFielder_Ints(t *testing.T) {
	t.Run("non-empty slice", func(t *testing.T) {
		values := []int{1, 2, 3}
		err := Err(errors.New("base error")).Ints("values", values).Msg("message")
		require.Contains(t, err.Error(), "message values=[1,2,3]: base error")
	})

	t.Run("empty slice", func(t *testing.T) {
		values := []int{}
		err := Err(errors.New("base error")).Ints("values", values).Msg("message")
		require.Contains(t, err.Error(), "message values=[]: base error")
	})
}

func TestFielder_Int8s(t *testing.T) {
	t.Run("non-empty slice", func(t *testing.T) {
		values := []int8{1, 2, 3}
		err := Err(errors.New("base error")).Int8s("values", values).Msg("message")
		require.Contains(t, err.Error(), "message values=[1,2,3]: base error")
	})

	t.Run("empty slice", func(t *testing.T) {
		values := []int8{}
		err := Err(errors.New("base error")).Int8s("values", values).Msg("message")
		require.Contains(t, err.Error(), "message values=[]: base error")
	})
}

func TestFielder_Int16s(t *testing.T) {
	t.Run("non-empty slice", func(t *testing.T) {
		values := []int16{1, 2, 3}
		err := Err(errors.New("base error")).Int16s("values", values).Msg("message")
		require.Contains(t, err.Error(), "message values=[1,2,3]: base error")
	})

	t.Run("empty slice", func(t *testing.T) {
		values := []int16{}
		err := Err(errors.New("base error")).Int16s("values", values).Msg("message")
		require.Contains(t, err.Error(), "message values=[]: base error")
	})
}

func TestFielder_Int32s(t *testing.T) {
	t.Run("non-empty slice", func(t *testing.T) {
		values := []int32{1, 2, 3}
		err := Err(errors.New("base error")).Int32s("values", values).Msg("message")
		require.Contains(t, err.Error(), "message values=[1,2,3]: base error")
	})

	t.Run("empty slice", func(t *testing.T) {
		values := []int32{}
		err := Err(errors.New("base error")).Int32s("values", values).Msg("message")
		require.Contains(t, err.Error(), "message values=[]: base error")
	})
}

func TestFielder_Uints(t *testing.T) {
	t.Run("non-empty slice", func(t *testing.T) {
		values := []uint{1, 2, 3}
		err := Err(errors.New("base error")).Uints("values", values).Msg("message")
		require.Contains(t, err.Error(), "message values=[1,2,3]: base error")
	})

	t.Run("empty slice", func(t *testing.T) {
		values := []uint{}
		err := Err(errors.New("base error")).Uints("values", values).Msg("message")
		require.Contains(t, err.Error(), "message values=[]: base error")
	})
}

func TestFielder_Uint8s(t *testing.T) {
	t.Run("non-empty slice", func(t *testing.T) {
		values := []uint8{1, 2, 3}
		err := Err(errors.New("base error")).Uint8s("values", values).Msg("message")
		require.Contains(t, err.Error(), "message values=[1,2,3]: base error")
	})

	t.Run("empty slice", func(t *testing.T) {
		values := []uint8{}
		err := Err(errors.New("base error")).Uint8s("values", values).Msg("message")
		require.Contains(t, err.Error(), "message values=[]: base error")
	})
}

func TestFielder_Uint16s(t *testing.T) {
	t.Run("non-empty slice", func(t *testing.T) {
		values := []uint16{1, 2, 3}
		err := Err(errors.New("base error")).Uint16s("values", values).Msg("message")
		require.Contains(t, err.Error(), "message values=[1,2,3]: base error")
	})

	t.Run("empty slice", func(t *testing.T) {
		values := []uint16{}
		err := Err(errors.New("base error")).Uint16s("values", values).Msg("message")
		require.Contains(t, err.Error(), "message values=[]: base error")
	})
}

func TestFielder_Uint32s(t *testing.T) {
	t.Run("non-empty slice", func(t *testing.T) {
		values := []uint32{1, 2, 3}
		err := Err(errors.New("base error")).Uint32s("values", values).Msg("message")
		require.Contains(t, err.Error(), "message values=[1,2,3]: base error")
	})

	t.Run("empty slice", func(t *testing.T) {
		values := []uint32{}
		err := Err(errors.New("base error")).Uint32s("values", values).Msg("message")
		require.Contains(t, err.Error(), "message values=[]: base error")
	})
}

func TestFielder_Bools(t *testing.T) {
	t.Run("non-empty slice", func(t *testing.T) {
		values := []bool{true, false, true}
		err := Err(errors.New("base error")).Bools("flags", values).Msg("message")
		require.Contains(t, err.Error(), "message flags=[true,false,true]: base error")
	})

	t.Run("empty slice", func(t *testing.T) {
		values := []bool{}
		err := Err(errors.New("base error")).Bools("flags", values).Msg("message")
		require.Contains(t, err.Error(), "message flags=[]: base error")
	})
}

func TestFielder_Time(t *testing.T) {
	testTime := time.Date(2025, 7, 19, 15, 30, 0, 0, time.UTC)
	err := Err(errors.New("base error")).Time("timestamp", testTime).Msg("message")
	require.Contains(t, err.Error(), "message timestamp=\"2025-07-19T15:30:00Z\": base error")
}

func TestFielder_Times(t *testing.T) {
	t.Run("non-empty slice", func(t *testing.T) {
		time1 := time.Date(2025, 7, 19, 15, 30, 0, 0, time.UTC)
		time2 := time.Date(2025, 7, 20, 15, 30, 0, 0, time.UTC)
		values := []time.Time{time1, time2}
		err := Err(errors.New("base error")).Times("timestamps", values).Msg("message")
		require.Contains(t, err.Error(), "message timestamps=[\"2025-07-19T15:30:00Z\",\"2025-07-20T15:30:00Z\"]: base error")
	})

	t.Run("empty slice", func(t *testing.T) {
		values := []time.Time{}
		err := Err(errors.New("base error")).Times("timestamps", values).Msg("message")
		require.Contains(t, err.Error(), "message timestamps=[]: base error")
	})
}

func TestFielder_XTime(t *testing.T) {
	testTime := time.Date(2025, 7, 19, 15, 30, 0, 0, time.UTC)
	xTime := types.NewTime(testTime)
	err := Err(errors.New("base error")).XTime("timestamp", xTime).Msg("message")
	require.Contains(t, err.Error(), "message timestamp=\"2025-07-19T15:30:00Z\": base error")
}

func TestFielder_XTimes(t *testing.T) {
	t.Run("non-empty slice", func(t *testing.T) {
		time1 := time.Date(2025, 7, 19, 15, 30, 0, 0, time.UTC)
		time2 := time.Date(2025, 7, 20, 15, 30, 0, 0, time.UTC)
		xTime1 := types.NewTime(time1)
		xTime2 := types.NewTime(time2)
		values := []*types.Time{xTime1, xTime2}
		err := Err(errors.New("base error")).XTimes("timestamps", values).Msg("message")
		require.Contains(t, err.Error(), "message timestamps=[\"2025-07-19T15:30:00Z\",\"2025-07-20T15:30:00Z\"]: base error")
	})

	t.Run("empty slice", func(t *testing.T) {
		values := []*types.Time{}
		err := Err(errors.New("base error")).XTimes("timestamps", values).Msg("message")
		require.Contains(t, err.Error(), "message timestamps=[]: base error")
	})
}

func TestFielder_MultipleFields(t *testing.T) {
	err := Err(errors.New("base error")).
		Str("str_field", "value").
		Int("int_field", 42).
		Bool("bool_field", true).
		Msg("complex message")

	require.True(t, strings.Contains(err.Error(), "str_field=value"))
	require.True(t, strings.Contains(err.Error(), "int_field=42"))
	require.True(t, strings.Contains(err.Error(), "bool_field=true"))
	require.True(t, strings.Contains(err.Error(), "complex message"))
	require.True(t, strings.Contains(err.Error(), "base error"))
}
