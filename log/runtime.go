package log

import (
	"runtime"
	"runtime/debug"

	"github.com/xakepp35/pkg/xslice"
)

func RuntimeStackTrace() [][]byte {
	stackBytes := debug.Stack()
	return xslice.SplitBytes(stackBytes, '\n')
}

const RuntimeFunctionSkip = 2

func RuntimeFunctionName(addSkip int) string {
	res := RuntimeFunction(addSkip + 1).Name()
	lastIndex := xslice.LastIndexByteString(res, '/')
	res = res[lastIndex+1:]
	return res
}

func RuntimeFunction(skip int) *runtime.Func {
	pc, _, _, ok := runtime.Caller(skip + 1)
	if !ok {
		return nil
	}
	return runtime.FuncForPC(pc)
}
