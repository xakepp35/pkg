package xpgx

import (
	"runtime"
	"strings"
)

func callerFunc(skip int) string {
	pc := make([]uintptr, 10)
	n := runtime.Callers(skip, pc)
	if n == 0 {
		return "unknown"
	}
	frames := runtime.CallersFrames(pc[:n])
	for {
		frame, more := frames.Next()
		// фильтруем внутренние вызовы pgx
		if !strings.Contains(frame.Function, "pgx") &&
			!strings.Contains(frame.Function, "pgconn") &&
			!strings.Contains(frame.Function, "xpgx") {
			return shortFuncName(frame.Function)
		}
		if !more {
			break
		}
	}
	return "unknown"
}

func shortFuncName(f string) string {
	parts := strings.Split(f, "/")
	return parts[len(parts)-1]
}
