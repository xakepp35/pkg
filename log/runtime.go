package log

import (
	"runtime"
	"runtime/debug"

	"github.com/xakepp35/pkg/fslice"
)

func RuntimeStackTrace() [][]byte {
	stackBytes := debug.Stack()
	return fslice.SplitBytes(stackBytes, '\n')
}

const RuntimeFunctionSkip = 2

func RuntimeFunctionName(addSkip int) string {
	res := RuntimeFunction(addSkip + 1).Name()
	lastIndex := fslice.LastIndexByteString(res, '/')
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
