package env

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Get returns the environment variable specified by key, or the default value specified by def if empty
func Get[T any](key string, def T) T {
	switch tp := (any)(def).(type) {
	case string:
		return (any)(String(key, tp)).(T)
	case []string:
		return (any)(Strings(key, tp)).(T)
	case int:
		return (any)(Int(key, tp)).(T)
	case int8:
		return (any)(int8(Int(key, int(tp)))).(T)
	case int16:
		return (any)(int16(Int(key, int(tp)))).(T)
	case int32:
		return (any)(int32(Int(key, int(tp)))).(T)
	case int64:
		return (any)(Int64(key, tp)).(T)
	case uint:
		return (any)(Uint(key, tp)).(T)
	case uint8:
		return (any)(uint8(Uint8(key, tp))).(T)
	case uint16:
		return (any)(uint16(Uint16(key, tp))).(T)
	case uint32:
		return (any)(uint32(Uint32(key, tp))).(T)
	case uint64:
		return (any)(Uint64(key, tp)).(T)
	case uintptr:
		return (any)(uintptr(Uint(key, uint(tp)))).(T)
	case float32:
		return (any)(float32(Float64(key, float64(tp)))).(T)
	case float64:
		return (any)(Float64(key, tp)).(T)
	case bool:
		return (any)(Bool(key, tp)).(T)
	case time.Duration:
		return (any)(Duration(key, tp)).(T)
	default:
		panic("Env: unsupported parameter type")
	}
}

const envListSeparator = ","

func String(key string, def string) string {
	val, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(val) == "" {
		return def
	}
	return strings.TrimSpace(val)
}

func Strings(key string, def []string) []string {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		return def
	}
	return strings.Split(val, envListSeparator)
}

func Int(key string, def int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	v, err := strconv.Atoi(val)
	if err != nil {
		return def
	}
	return v
}

func Int64(key string, def int64) int64 {
	val, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	v, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return def
	}
	return v
}

func Int32(key string, def int32) int32 {
	val, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	v, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return def
	}
	return int32(v)
}

func Int16(key string, def int16) int16 {
	val, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	v, err := strconv.ParseInt(val, 10, 16)
	if err != nil {
		return def
	}
	return int16(v)
}

func Int8(key string, def int8) int8 {
	val, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	v, err := strconv.ParseInt(val, 10, 8)
	if err != nil {
		return def
	}
	return int8(v)
}

func Uint(key string, def uint) uint {
	val, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	v, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return def
	}
	return uint(v)
}

func Uint64(key string, def uint64) uint64 {
	val, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	v, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return def
	}
	return v
}

func Uint32(key string, def uint32) uint32 {
	val, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	v, err := strconv.ParseUint(val, 10, 32)
	if err != nil {
		return def
	}
	return uint32(v)
}

func Uint16(key string, def uint16) uint16 {
	val, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	v, err := strconv.ParseUint(val, 10, 16)
	if err != nil {
		return def
	}
	return uint16(v)
}

func Uint8(key string, def uint8) uint8 {
	val, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	v, err := strconv.ParseUint(val, 10, 8)
	if err != nil {
		return def
	}
	return uint8(v)
}

func Uintptr(key string, def uintptr) uintptr {
	return uintptr(Uint64(key, uint64(def)))
}

func Float64(key string, def float64) float64 {
	val, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	v, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return def
	}
	return v
}

func Float32(key string, def float32) float32 {
	val, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	v, err := strconv.ParseFloat(val, 32)
	if err != nil {
		return def
	}
	return float32(v)
}

func Bool(key string, def bool) bool {
	val, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	return strings.ToLower(val) == "true"
}

func Duration(key string, def time.Duration) time.Duration {
	val, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	v, err := time.ParseDuration(val)
	if err != nil {
		return def
	}
	return v
}

func Time(key string, def time.Time) time.Time {
	val, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	v, err := time.Parse(time.RFC3339Nano, val)
	if err != nil {
		return def
	}
	return v
}
