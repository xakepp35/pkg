package env

import (
	"os"
	"strconv"
	"strings"
	"time"
)

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
		return (any)(uint8(Uint(key, uint(tp)))).(T)
	case uint16:
		return (any)(uint16(Uint(key, uint(tp)))).(T)
	case uint32:
		return (any)(uint32(Uint(key, uint(tp)))).(T)
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

// Env returns the environment variable specified by key, or the default value specified by def if empty
func String(key string, def string) string {
	result := os.Getenv(key)
	if result == "" {
		return def
	}
	return strings.TrimSpace(result)
}

const envListSeparator = ","

func Strings(key string, def []string) []string {
	return strings.Split(Get(key, strings.Join(def, envListSeparator)), envListSeparator)
}

// IntFromEnv Check environment variable and return the value as int
func Int(key string, def int) int {
	envStr := os.Getenv(key)
	if envStr == "" {
		return def
	}
	result, err := strconv.Atoi(envStr)
	if err != nil {
		return def
	}
	return result
}

func Int64(key string, def int64) int64 {
	envStr := os.Getenv(key)
	if envStr == "" {
		return def
	}
	result, err := strconv.ParseInt(envStr, 10, 64)
	if err != nil {
		return def
	}
	return result
}

func Uint(key string, def uint) uint {
	envStr := os.Getenv(key)
	if envStr == "" {
		return def
	}
	result, err := strconv.ParseUint(envStr, 10, 64)
	if err != nil {
		return def
	}
	return uint(result)
}

func Uint64(key string, def uint64) uint64 {
	envStr := os.Getenv(key)
	if envStr == "" {
		return def
	}
	result, err := strconv.ParseUint(envStr, 10, 64)
	if err != nil {
		return def
	}
	return result
}

func Float64(key string, def float64) float64 {
	envStr := os.Getenv(key)
	if envStr == "" {
		return def
	}
	result, err := strconv.ParseFloat(envStr, 64)
	if err != nil {
		return def
	}
	return result
}

// Bool return the value of env variable as bool
func Bool(key string, def bool) bool {
	envStr := os.Getenv(key)
	if envStr == "" {
		return def
	}
	return strings.ToLower(envStr) == "true"
}

func Duration(key string, def time.Duration) time.Duration {
	envStr := os.Getenv(key)
	if envStr == "" {
		return def
	}
	result, err := time.ParseDuration(envStr)
	if err != nil {
		return def
	}
	return result
}
