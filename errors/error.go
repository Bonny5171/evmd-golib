package errors

import (
	"fmt"
	"runtime"
)

func Trace(err error) error {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	return fmt.Errorf("[%s:%d %s]: %w", file, line-1, f.Name(), err)
}
