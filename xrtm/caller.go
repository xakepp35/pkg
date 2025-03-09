package xrtm

import (
	"runtime"

	"github.com/xakepp35/pkg/xslice"
)

const PathSeparator = '/'

const CallerFnDefaultSkip = 2

// RuntimeFunctionName возвращает имя вызывающей функции
func CallerFnName(addSkip int) string {
	fn := CallerFn(addSkip + 1)
	res := fn.Name()
	lastIndex := xslice.LastIndexByteString(res, PathSeparator)
	return res[lastIndex+1:]
}

// RuntimeFunction получает *runtime.Func для указанного уровня вызова
func CallerFn(skip int) *runtime.Func {
	pc, _, _, ok := runtime.Caller(skip + 1)
	if !ok {
		return nil
	}
	return runtime.FuncForPC(pc)
}
