package utils

import (
	"strconv"
	"strings"
	"time"
)

func FirstArgPtr[T any, E any](f T, s E) (*T, E) {
	return &f, s
}

func ParseRepeated[T any](arr string, parser func(string) (T, error)) ([]T, error) {
	j := 0
	out := make([]T, 0, strings.Count(arr, ","))
	for i := 0; i < len(arr); i++ {
		if arr[i] != ',' {
			continue
		}

		t, err := parser(arr[j:i])
		if err != nil {
			return nil, err
		}

		out = append(out, t)
		j = i + 1
	}

	t, err := parser(arr[j:])
	if err != nil {
		return nil, err
	}

	out = append(out, t)

	return out, nil
}

// ParseInt32 Integer types int32
func ParseInt32(s string) (int32, error) {
	v, err := strconv.ParseInt(s, 10, 32)
	return int32(v), err
}

// ParseInt64 Integer types int64
func ParseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// ParseUint32 Integer types uint32
func ParseUint32(s string) (uint32, error) {
	v, err := strconv.ParseUint(s, 10, 32)
	return uint32(v), err
}

// ParseUint64 Integer types uint64
func ParseUint64(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}

// ParseFloat32 Floating point types
func ParseFloat32(s string) (float32, error) {
	v, err := strconv.ParseFloat(s, 32)
	return float32(v), err
}

// ParseFloat64 Floating64 point types
func ParseFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// ParseBool Boolean
func ParseBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

// ParseBytes Bytes (base64 or hex manually if needed; here raw conversion)
func ParseBytes(s string) ([]byte, error) {
	return []byte(s), nil
}

// ParseTimestamp Timestamp (RFC3339)
func ParseTimestamp(s string) (*time.Time, error) {
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// ParseEnum Wrapper for enums (assumes numeric)
func ParseEnum[T ~int32](s string) (T, error) {
	v, err := strconv.ParseInt(s, 10, 32)
	return T(v), err
}
